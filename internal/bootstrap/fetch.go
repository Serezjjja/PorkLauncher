package bootstrap

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

// Metadata represents the payload metadata downloaded from the update server.
type Metadata struct {
	Version             string                      `json:"version"`
	Payload             map[string]PlatformArtifact `json:"payload"`
	MinBootstrapVersion string                      `json:"min_bootstrap_version,omitempty"`
}

// PlatformArtifact describes a downloadable payload for a specific platform.
type PlatformArtifact struct {
	URL    string `json:"url"`
	SHA256 string `json:"sha256"`
	Size   int64  `json:"size"`
}

var httpClient = &http.Client{
	Timeout: 30 * time.Second,
}

// FetchMetadata downloads and parses metadata.json from the update server.
func FetchMetadata(ctx context.Context, baseURL string) (*Metadata, []byte, error) {
	metadataURL := baseURL + "/" + MetadataFileName

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, metadataURL, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("create metadata request: %w", err)
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, nil, fmt.Errorf("fetch metadata: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, nil, fmt.Errorf("fetch metadata: HTTP %d", resp.StatusCode)
	}

	// Read raw bytes for signature verification
	rawBody, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20)) // 1 MiB limit
	if err != nil {
		return nil, nil, fmt.Errorf("read metadata body: %w", err)
	}

	var meta Metadata
	if err := json.Unmarshal(rawBody, &meta); err != nil {
		return nil, nil, fmt.Errorf("parse metadata JSON: %w", err)
	}

	if meta.Version == "" {
		return nil, nil, fmt.Errorf("metadata has empty version")
	}

	return &meta, rawBody, nil
}

// FetchSignature downloads the metadata.json.sig file from the update server.
func FetchSignature(ctx context.Context, baseURL string) ([]byte, error) {
	sigURL := baseURL + "/" + SignatureFileName

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, sigURL, nil)
	if err != nil {
		return nil, fmt.Errorf("create signature request: %w", err)
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fetch signature: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("fetch signature: HTTP %d", resp.StatusCode)
	}

	// Ed25519 signature is 64 bytes, allow some overhead for encoding
	sig, err := io.ReadAll(io.LimitReader(resp.Body, 512))
	if err != nil {
		return nil, fmt.Errorf("read signature: %w", err)
	}

	return sig, nil
}

// DownloadPayload downloads the payload archive to a temporary file with progress reporting.
// Returns the path to the downloaded file.
func DownloadPayload(ctx context.Context, artifact PlatformArtifact, destDir string, progress func(downloaded, total int64)) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, artifact.URL, nil)
	if err != nil {
		return "", fmt.Errorf("create payload request: %w", err)
	}

	// Use a longer timeout for payload downloads
	client := &http.Client{
		Timeout: 10 * time.Minute,
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("download payload: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("download payload: HTTP %d", resp.StatusCode)
	}

	if err := os.MkdirAll(destDir, 0755); err != nil {
		return "", fmt.Errorf("create temp dir: %w", err)
	}

	tmpFile, err := os.CreateTemp(destDir, "payload-*.zip.tmp")
	if err != nil {
		return "", fmt.Errorf("create temp file: %w", err)
	}
	tmpPath := tmpFile.Name()

	// Ensure cleanup on error
	success := false
	defer func() {
		tmpFile.Close()
		if !success {
			os.Remove(tmpPath)
		}
	}()

	total := resp.ContentLength
	if total <= 0 && artifact.Size > 0 {
		total = artifact.Size
	}

	var downloaded int64
	buf := make([]byte, 32*1024)

	for {
		n, readErr := resp.Body.Read(buf)
		if n > 0 {
			if _, writeErr := tmpFile.Write(buf[:n]); writeErr != nil {
				return "", fmt.Errorf("write payload to disk: %w", writeErr)
			}
			downloaded += int64(n)
			if progress != nil {
				progress(downloaded, total)
			}
		}
		if readErr != nil {
			if readErr == io.EOF {
				break
			}
			return "", fmt.Errorf("read payload stream: %w", readErr)
		}
	}

	if err := tmpFile.Sync(); err != nil {
		return "", fmt.Errorf("sync payload file: %w", err)
	}

	// Rename to final name
	finalPath := filepath.Join(destDir, "payload.zip")
	if err := tmpFile.Close(); err != nil {
		return "", fmt.Errorf("close temp file: %w", err)
	}
	if err := os.Rename(tmpPath, finalPath); err != nil {
		return "", fmt.Errorf("rename temp file: %w", err)
	}

	success = true
	return finalPath, nil
}
