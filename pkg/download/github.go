package download

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"HyLauncher/internal/progress"
	"HyLauncher/pkg/logger"
)

const (
	defaultRepoOwner = "Serezjjja"
	defaultRepoName  = "PorkLauncher"
)

type GitHubReleaseAsset struct {
	Name               string `json:"name"`
	BrowserDownloadURL string `json:"browser_download_url"`
	Size               int64  `json:"size"`
}

type GitHubRelease struct {
	TagName string               `json:"tag_name"`
	Name    string               `json:"name"`
	Assets  []GitHubReleaseAsset `json:"assets"`
}

// DownloadLatestReleaseAsset downloads an asset from the latest GitHub release
// If reporter and scaler are nil, downloads silently without progress updates
func DownloadLatestReleaseAsset(
	ctx context.Context,
	assetName string,
	destPath string,
	stage progress.Stage,
	reporter *progress.Reporter,
	scaler *progress.Scaler,
) error {
	apiURL := fmt.Sprintf("https://api.github.com/repos/%s/%s/releases/latest", defaultRepoOwner, defaultRepoName)

	// Create HTTP request with context
	req, err := http.NewRequestWithContext(ctx, "GET", apiURL, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers for GitHub API
	req.Header.Set("Accept", "application/vnd.github.v3+json")
	req.Header.Set("User-Agent", "HyLauncher")

	// Make the request
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to query GitHub API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		logger.Error("GitHub API error", "status", resp.StatusCode, "status_text", resp.Status)
		return fmt.Errorf("GitHub API returned status %d: %s", resp.StatusCode, resp.Status)
	}

	// Decode the release information
	var release GitHubRelease
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		logger.Error("Failed to decode GitHub release", "error", err)
		return fmt.Errorf("failed to decode GitHub release JSON: %w", err)
	}

	logger.Info("GitHub release found", "tag", release.TagName, "assets_count", len(release.Assets))

	// Find the requested asset
	var downloadURL string
	for _, asset := range release.Assets {
		if asset.Name == assetName {
			downloadURL = asset.BrowserDownloadURL
			break
		}
	}

	if downloadURL == "" {
		logger.Error("Asset not found in release", "asset", assetName, "tag", release.TagName, "available_assets", release.Assets)
		return fmt.Errorf("asset '%s' not found in latest release (tag: %s)", assetName, release.TagName)
	}

	logger.Info("Asset found", "asset", assetName, "url", downloadURL)

	// Ensure destination directory exists
	destDir := filepath.Dir(destPath)
	if err := os.MkdirAll(destDir, 0755); err != nil {
		return fmt.Errorf("failed to create destination directory: %w", err)
	}

	// Download the file with progress if reporter is provided
	if reporter != nil {
		reporter.Report(stage, 0, fmt.Sprintf("Downloading %s from release %s...", assetName, release.TagName))
	}

	if err := DownloadWithReporter(ctx, destPath, downloadURL, assetName, reporter, stage, scaler); err != nil {
		// Clean up partial download on error
		_ = os.Remove(destPath)
		return fmt.Errorf("failed to download %s: %w", assetName, err)
	}

	if reporter != nil {
		reporter.Report(stage, 100, fmt.Sprintf("Downloaded %s successfully", assetName))
	}

	return nil
}

func GetLatestReleaseInfo(ctx context.Context) (*GitHubRelease, error) {
	apiURL := fmt.Sprintf("https://api.github.com/repos/%s/%s/releases/latest", defaultRepoOwner, defaultRepoName)

	req, err := http.NewRequestWithContext(ctx, "GET", apiURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Accept", "application/vnd.github.v3+json")
	req.Header.Set("User-Agent", "HyLauncher")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to query GitHub API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		logger.Error("GitHub API error in GetLatestReleaseInfo", "status", resp.StatusCode)
		return nil, fmt.Errorf("GitHub API returned status %d: %s", resp.StatusCode, resp.Status)
	}

	var release GitHubRelease
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		logger.Error("Failed to decode release info", "error", err)
		return nil, fmt.Errorf("failed to decode GitHub release JSON: %w", err)
	}

	logger.Info("Latest release info", "tag", release.TagName, "name", release.Name, "assets", len(release.Assets))
	return &release, nil
}

func ListLatestReleaseAssets(ctx context.Context) ([]GitHubReleaseAsset, error) {
	release, err := GetLatestReleaseInfo(ctx)
	if err != nil {
		return nil, err
	}
	return release.Assets, nil
}
