package patch

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	"HyLauncher/internal/config"
	"HyLauncher/internal/env"
	"HyLauncher/pkg/logger"
)

type VersionCheckResult struct {
	LatestVersion int
	Error         error
}

type AllVersionsResult struct {
	Versions []int
	Error    error
}

// PatchesConfigResponse represents the response from patches-config endpoint
// Returns { patches_url: "https://..." }
type PatchesConfigResponse struct {
	PatchesURL string `json:"patches_url"`
}

// ManifestFile represents a file entry in the manifest
type ManifestFile struct {
	Size int64 `json:"size"`
}

// Manifest represents the manifest.json structure
// { files: { "windows/amd64/release/0_to_11.pwr": { size: 12345 }, ... } }
type Manifest struct {
	Files map[string]ManifestFile `json:"files"`
}

// PatchInfo represents a parsed patch entry
type PatchInfo struct {
	From int
	To   int
	Key  string
	Size int64
}

// patchesState holds the cached patches URL and manifest
var (
	patchesStateMu    sync.RWMutex
	patchesBaseURL    string
	patchesBaseURLSet time.Time
	manifestCache     *Manifest
	manifestCacheSet  time.Time
	fallbackBuild     = 11 // Fallback latest build if manifest unreachable
)

type cache struct {
	mu            sync.RWMutex
	latestVersion map[string]*VersionCheckResult
	versionExists map[string]bool
	allVersions   map[string]*AllVersionsResult
	lastSet       map[string]time.Time
	ttl           time.Duration
}

var versionCache = &cache{
	latestVersion: make(map[string]*VersionCheckResult),
	versionExists: make(map[string]bool),
	allVersions:   make(map[string]*AllVersionsResult),
	lastSet:       make(map[string]time.Time),
	ttl:           5 * time.Minute,
}

func (c *cache) getLatest(key string) (*VersionCheckResult, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	result, exists := c.latestVersion[key]
	if !exists || time.Since(c.lastSet[key]) >= c.ttl {
		return nil, false
	}
	return result, true
}

func (c *cache) setLatest(key string, result *VersionCheckResult) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.latestVersion[key] = result
	c.lastSet[key] = time.Now()
}

func (c *cache) checkVersion(key string) (bool, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	exists, cached := c.versionExists[key]
	return exists, cached
}

func (c *cache) setVersion(key string, exists bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.versionExists[key] = exists
}

func (c *cache) getAllVersions(key string) (*AllVersionsResult, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	result, exists := c.allVersions[key]
	if !exists || time.Since(c.lastSet[key]) >= c.ttl {
		return nil, false
	}
	return result, true
}

func (c *cache) setAllVersions(key string, result *AllVersionsResult) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.allVersions[key] = result
	c.lastSet[key] = time.Now()
}

func (c *cache) clear() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.latestVersion = make(map[string]*VersionCheckResult)
	c.versionExists = make(map[string]bool)
	c.allVersions = make(map[string]*AllVersionsResult)
	c.lastSet = make(map[string]time.Time)
}

// fetchPatchesConfigWithFallback fetches patches URL from config endpoints with fallback
func fetchPatchesConfigWithFallback() (string, error) {
	// Check memory cache (5 min TTL)
	patchesStateMu.RLock()
	if patchesBaseURL != "" && time.Since(patchesBaseURLSet) < 5*time.Minute {
		url := patchesBaseURL
		patchesStateMu.RUnlock()
		return url, nil
	}
	patchesStateMu.RUnlock()

	// Try all config sources
	sources := config.GetPatchesConfigSources()
	client := &http.Client{Timeout: 8 * time.Second}

	for _, source := range sources {
		logger.Info("Trying patches config source", "url", source)
		resp, err := client.Get(source)
		if err != nil {
			logger.Warn("Patches config source failed", "url", source, "error", err)
			continue
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			logger.Warn("Patches config returned non-200", "url", source, "status", resp.StatusCode)
			continue
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			logger.Warn("Failed to read patches config response", "url", source, "error", err)
			continue
		}

		var cfg PatchesConfigResponse
		if err := json.Unmarshal(body, &cfg); err != nil {
			logger.Warn("Failed to parse patches config", "url", source, "error", err)
			continue
		}

		if cfg.PatchesURL != "" {
			url := strings.TrimRight(cfg.PatchesURL, "/")
			patchesStateMu.Lock()
			patchesBaseURL = url
			patchesBaseURLSet = time.Now()
			patchesStateMu.Unlock()
			logger.Info("Got patches URL from config", "url", url)
			return url, nil
		}
	}

	// Use hardcoded fallback
	fallback := config.GetPatchesFallbackURL()
	logger.Warn("All patches config sources failed, using fallback", "url", fallback)
	patchesStateMu.Lock()
	patchesBaseURL = fallback
	patchesBaseURLSet = time.Now()
	patchesStateMu.Unlock()
	return fallback, nil
}

// fetchManifest fetches the manifest.json from the patches base URL
func fetchManifest() (*Manifest, error) {
	// Check memory cache (1 min TTL)
	patchesStateMu.RLock()
	if manifestCache != nil && time.Since(manifestCacheSet) < 1*time.Minute {
		m := manifestCache
		patchesStateMu.RUnlock()
		return m, nil
	}
	patchesStateMu.RUnlock()

	baseURL, err := fetchPatchesConfigWithFallback()
	if err != nil {
		return nil, err
	}

	manifestURL := baseURL + "/manifest.json"
	logger.Info("Fetching manifest", "url", manifestURL)

	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Get(manifestURL)
	if err != nil {
		return nil, fmt.Errorf("fetch manifest: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("manifest returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read manifest: %w", err)
	}

	var manifest Manifest
	if err := json.Unmarshal(body, &manifest); err != nil {
		return nil, fmt.Errorf("parse manifest: %w", err)
	}

	if manifest.Files == nil {
		return nil, fmt.Errorf("invalid manifest: no files")
	}

	patchesStateMu.Lock()
	manifestCache = &manifest
	manifestCacheSet = time.Now()
	patchesStateMu.Unlock()

	logger.Info("Manifest fetched successfully", "files", len(manifest.Files))
	return &manifest, nil
}

// getPlatformPatches extracts patches for current platform and branch
func getPlatformPatches(manifest *Manifest, branch string) []PatchInfo {
	os := env.GetOS()
	arch := env.GetArchForAPI()
	prefix := fmt.Sprintf("%s/%s/%s/", os, arch, branch)

	var patches []PatchInfo
	for key, info := range manifest.Files {
		if !strings.HasPrefix(key, prefix) || !strings.HasSuffix(key, ".pwr") {
			continue
		}

		// Parse filename: "0_to_11.pwr" -> from=0, to=11
		filename := key[len(prefix) : len(key)-4] // Remove prefix and .pwr
		parts := strings.Split(filename, "_to_")
		if len(parts) != 2 {
			continue
		}

		from, err1 := strconv.Atoi(parts[0])
		to, err2 := strconv.Atoi(parts[1])
		if err1 != nil || err2 != nil {
			continue
		}

		patches = append(patches, PatchInfo{
			From: from,
			To:   to,
			Key:  key,
			Size: info.Size,
		})
	}

	return patches
}

func FindLatestVersion(branch string) (int, error) {
	key := cacheKey(branch)

	if cached, ok := versionCache.getLatest(key); ok {
		return cached.LatestVersion, cached.Error
	}

	// Use coalescer to prevent duplicate in-flight requests
	result, err := versionCoalescer.Do("latest:"+key, func() (interface{}, error) {
		r := findLatestVersion(branch)
		versionCache.setLatest(key, &r)
		return r.LatestVersion, r.Error
	})

	if err != nil {
		return 0, err
	}
	return result.(int), nil
}

func ListAllVersions(branch string) ([]int, error) {
	key := cacheKey(branch)

	if cached, ok := versionCache.getAllVersions(key); ok {
		return cached.Versions, cached.Error
	}

	// Use coalescer to prevent duplicate in-flight requests
	result, err := versionCoalescer.Do("all:"+key, func() (interface{}, error) {
		r := listAllVersions(branch)
		versionCache.setAllVersions(key, &r)
		return r.Versions, r.Error
	})

	if err != nil {
		return nil, err
	}
	return result.([]int), nil
}

func ClearVersionCache() {
	versionCache.clear()
}

func VerifyVersionExists(branch string, version int) error {
	client := createClient()

	if exists := checkVersionExists(client, branch, version); exists {
		return nil
	}

	return fmt.Errorf("version %d not found", version)
}

func findLatestVersion(branch string) VersionCheckResult {
	manifest, err := fetchManifest()
	if err != nil {
		logger.Warn("Failed to fetch manifest, using fallback build", "error", err, "fallback", fallbackBuild)
		return VersionCheckResult{LatestVersion: fallbackBuild}
	}

	patches := getPlatformPatches(manifest, branch)
	if len(patches) == 0 {
		logger.Warn("No patches found for platform, using fallback build", "os", runtime.GOOS, "arch", runtime.GOARCH, "fallback", fallbackBuild)
		return VersionCheckResult{LatestVersion: fallbackBuild}
	}

	// Find the highest "to" version
	latest := 0
	for _, p := range patches {
		if p.To > latest {
			latest = p.To
		}
	}

	logger.Info("Found latest version", "branch", branch, "version", latest)
	return VersionCheckResult{LatestVersion: latest}
}

func listAllVersions(branch string) AllVersionsResult {
	manifest, err := fetchManifest()
	if err != nil {
		return AllVersionsResult{
			Error: fmt.Errorf("cannot reach API or no patches available for %s/%s: %w", runtime.GOOS, runtime.GOARCH, err),
		}
	}

	patches := getPlatformPatches(manifest, branch)
	if len(patches) == 0 {
		return AllVersionsResult{
			Error: fmt.Errorf("no patches available for %s/%s", runtime.GOOS, runtime.GOARCH),
		}
	}

	versionMap := make(map[int]bool)
	for _, p := range patches {
		versionMap[p.From] = true
		versionMap[p.To] = true
	}

	var versions []int
	for v := range versionMap {
		versions = append(versions, v)
	}

	// Sort ascending
	for i := 0; i < len(versions); i++ {
		for j := i + 1; j < len(versions); j++ {
			if versions[i] > versions[j] {
				versions[i], versions[j] = versions[j], versions[i]
			}
		}
	}

	return AllVersionsResult{Versions: versions}
}

func fetchPatchStepsFromAPI(branch string, currentVer int) ([]PatchStep, error) {
	manifest, err := fetchManifest()
	if err != nil {
		return nil, err
	}

	patches := getPlatformPatches(manifest, branch)
	if len(patches) == 0 {
		return nil, fmt.Errorf("no patches available for %s/%s/%s", runtime.GOOS, runtime.GOARCH, branch)
	}

	// For now, we just return a full patch (0 -> currentVer)
	// The actual patch selection logic is in the pwr_patcher.go
	var steps []PatchStep
	baseURL, _ := fetchPatchesConfigWithFallback()

	for _, p := range patches {
		if p.From == currentVer {
			steps = append(steps, PatchStep{
				From:    p.From,
				To:      p.To,
				PWR:     fmt.Sprintf("%s/%s", baseURL, p.Key),
				PWRHead: "", // Not used in new API
				Sig:     "", // Signature not used in new API
			})
		}
	}

	return steps, nil
}

func checkVersionExists(client *http.Client, branch string, version int) bool {
	key := versionCacheKey(branch, version)

	if exists, cached := versionCache.checkVersion(key); cached {
		return exists
	}

	url := buildPatchURL(branch, version)
	resp, err := client.Head(url)
	if resp != nil {
		defer resp.Body.Close()
	}

	exists := err == nil && resp.StatusCode == http.StatusOK
	versionCache.setVersion(key, exists)

	time.Sleep(100 * time.Millisecond)

	return exists
}

func createClient() *http.Client {
	return &http.Client{
		Timeout: 10 * time.Second,
		Transport: &http.Transport{
			Proxy:                 http.ProxyFromEnvironment,
			MaxIdleConns:          3,
			MaxIdleConnsPerHost:   2,
			MaxConnsPerHost:       3,
			IdleConnTimeout:       30 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
		},
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
}

func buildPatchURL(branch string, version int) string {
	baseURL, _ := fetchPatchesConfigWithFallback()
	return fmt.Sprintf("%s/%s/%s/%s/0_to_%d.pwr",
		baseURL, runtime.GOOS, runtime.GOARCH, branch, version)
}

func cacheKey(branch string) string {
	return fmt.Sprintf("%s-%s-%s", runtime.GOOS, runtime.GOARCH, branch)
}

func versionCacheKey(branch string, version int) string {
	return fmt.Sprintf("%s-%d", cacheKey(branch), version)
}
