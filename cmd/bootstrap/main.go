// Bootstrap loader for PorkLauncher.
//
// This binary contains NO business logic, no game URLs, no auth tokens,
// no frontend assets. It only:
//  1. Fetches metadata.json + metadata.json.sig from a hardcoded update server
//  2. Verifies Ed25519 signature of the metadata
//  3. Downloads the payload archive (Wails app) for the current platform
//  4. Verifies SHA256 hash of the downloaded payload
//  5. Extracts and launches the payload
//
// Build:
//
//	go build -ldflags="-s -w \
//	  -X 'HyLauncher/internal/bootstrap.UpdateServerURL=https://example.com/updates' \
//	  -X 'HyLauncher/internal/bootstrap.Ed25519PublicKeyHex=<64-hex-chars>' \
//	  -X 'HyLauncher/internal/bootstrap.FallbackDownloadURL=https://github.com/Serezjjja/PorkLauncher/releases'" \
//	  -o PorkLauncher ./cmd/bootstrap
package main

import (
	"HyLauncher/internal/bootstrap"
	"context"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"syscall"
)

func main() {
	os.Exit(run())
}

func run() int {
	ui := bootstrap.NewUI()

	ui.Status(fmt.Sprintf("PorkLauncher Bootstrap v%s (%s/%s)",
		bootstrap.BootstrapVersion, runtime.GOOS, runtime.GOARCH))

	// Set up cancellation on interrupt
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	// --- Step 1: Validate embedded configuration ---

	serverURL, err := bootstrap.GetUpdateServerURL()
	if err != nil {
		ui.Error(err.Error())
		ui.FallbackMessage()
		return 1
	}

	pubKey, err := bootstrap.GetPublicKey()
	if err != nil {
		ui.Error(err.Error())
		ui.FallbackMessage()
		return 1
	}

	// --- Step 2: Check if cached payload exists ---

	appDir := getAppDataDir()
	payloadBaseDir := filepath.Join(appDir, bootstrap.PayloadDirName)
	versionFilePath := filepath.Join(payloadBaseDir, "current_version")

	// Try to use cached version first, update check happens in parallel
	cachedVersion := readCachedVersion(versionFilePath)
	if cachedVersion != "" {
		cachedDir := filepath.Join(payloadBaseDir, cachedVersion)
		exePath := filepath.Join(cachedDir, bootstrap.PayloadExecutableName())
		if _, err := os.Stat(exePath); err == nil {
			ui.Status(fmt.Sprintf("Found cached payload v%s, checking for updates...", cachedVersion))

			// Try to check for updates, but don't block on failure
			meta, rawMeta, fetchErr := bootstrap.FetchMetadata(ctx, serverURL)
			if fetchErr == nil {
				sig, sigErr := bootstrap.FetchSignature(ctx, serverURL)
				if sigErr == nil {
					if verifyErr := bootstrap.VerifyMetadataSignature(pubKey, rawMeta, sig); verifyErr == nil {
						if meta.Version != cachedVersion {
							ui.Status(fmt.Sprintf("New version available: v%s -> v%s", cachedVersion, meta.Version))
							// Proceed to download new version below
							return downloadAndLaunch(ctx, ui, meta, pubKey, payloadBaseDir, versionFilePath)
						}
					}
				}
			}

			// No update or failed to check — launch cached version
			ui.Status("Launching cached payload...")
			if err := bootstrap.LaunchPayload(cachedDir); err != nil {
				ui.Error(fmt.Sprintf("Failed to launch cached payload: %v", err))
				ui.FallbackMessage()
				return 1
			}
			return 0
		}
	}

	// --- Step 3: No cache — must download ---

	ui.Status("Fetching update metadata...")

	meta, rawMeta, err := bootstrap.FetchMetadata(ctx, serverURL)
	if err != nil {
		ui.Error(fmt.Sprintf("Failed to fetch metadata: %v", err))
		ui.FallbackMessage()
		return 1
	}

	ui.Status("Verifying metadata signature...")

	sig, err := bootstrap.FetchSignature(ctx, serverURL)
	if err != nil {
		ui.Error(fmt.Sprintf("Failed to fetch signature: %v", err))
		ui.FallbackMessage()
		return 1
	}

	if err := bootstrap.VerifyMetadataSignature(pubKey, rawMeta, sig); err != nil {
		ui.Error(fmt.Sprintf("Signature verification FAILED: %v", err))
		ui.Error("This may indicate a tampered update. Aborting for safety.")
		ui.FallbackMessage()
		return 1
	}

	ui.Status("Signature verified OK")

	return downloadAndLaunch(ctx, ui, meta, pubKey, payloadBaseDir, versionFilePath)
}

func downloadAndLaunch(ctx context.Context, ui *bootstrap.UI, meta *bootstrap.Metadata, pubKey []byte, payloadBaseDir, versionFilePath string) int {
	platform := bootstrap.PlatformKey()
	artifact, ok := meta.Payload[platform]
	if !ok {
		ui.Error(fmt.Sprintf("No payload available for platform %s", platform))
		ui.FallbackMessage()
		return 1
	}

	ui.Status(fmt.Sprintf("Downloading payload v%s for %s (%s)...",
		meta.Version, platform, formatSize(artifact.Size)))

	// Download to a temp directory first
	tmpDir := filepath.Join(payloadBaseDir, ".download-tmp")
	os.RemoveAll(tmpDir) // clean up previous failed downloads
	if err := os.MkdirAll(tmpDir, 0755); err != nil {
		ui.Error(fmt.Sprintf("Failed to create temp dir: %v", err))
		ui.FallbackMessage()
		return 1
	}
	defer os.RemoveAll(tmpDir)

	zipPath, err := bootstrap.DownloadPayload(ctx, artifact, tmpDir, ui.Progress)
	if err != nil {
		ui.ProgressDone()
		ui.Error(fmt.Sprintf("Download failed: %v", err))
		ui.FallbackMessage()
		return 1
	}
	ui.ProgressDone()

	// --- Verify SHA256 ---

	ui.Status("Verifying payload integrity...")

	if err := bootstrap.VerifyFileSHA256(zipPath, artifact.SHA256); err != nil {
		ui.Error(fmt.Sprintf("Integrity check FAILED: %v", err))
		ui.Error("The downloaded file may be corrupted or tampered with.")
		ui.FallbackMessage()
		return 1
	}

	ui.Status("SHA256 verified OK")

	// --- Extract ---

	destDir := filepath.Join(payloadBaseDir, meta.Version)
	// Clean destination if it exists (partial previous install)
	os.RemoveAll(destDir)

	ui.Status(fmt.Sprintf("Extracting to %s...", destDir))

	files, err := bootstrap.ExtractZip(zipPath, destDir)
	if err != nil {
		ui.Error(fmt.Sprintf("Extraction failed: %v", err))
		os.RemoveAll(destDir)
		ui.FallbackMessage()
		return 1
	}

	ui.Status(fmt.Sprintf("Extracted %d files", len(files)))

	// --- Save version ---

	if err := os.WriteFile(versionFilePath, []byte(meta.Version), 0644); err != nil {
		// Non-fatal: next boot will just re-download
		ui.Error(fmt.Sprintf("Warning: could not save version cache: %v", err))
	}

	// --- Cleanup old versions ---

	cleanupOldVersions(payloadBaseDir, meta.Version)

	// --- Launch ---

	ui.Status("Launching PorkLauncher...")

	if err := bootstrap.LaunchPayload(destDir); err != nil {
		ui.Error(fmt.Sprintf("Failed to launch: %v", err))
		ui.FallbackMessage()
		return 1
	}

	return 0
}

// getAppDataDir returns the platform-specific app data directory.
// Mirrors internal/env.GetDefaultAppDir() but without importing it.
func getAppDataDir() string {
	switch runtime.GOOS {
	case "windows":
		if local := os.Getenv("LOCALAPPDATA"); local != "" {
			return filepath.Join(local, "HyLauncher")
		}
		home, _ := os.UserHomeDir()
		return filepath.Join(home, "AppData", "Local", "HyLauncher")
	case "darwin":
		home, _ := os.UserHomeDir()
		return filepath.Join(home, "Library", "Application Support", "HyLauncher")
	default:
		home, _ := os.UserHomeDir()
		return filepath.Join(home, ".hylauncher")
	}
}

// readCachedVersion reads the last successfully launched payload version.
func readCachedVersion(path string) string {
	data, err := os.ReadFile(path)
	if err != nil {
		return ""
	}
	v := string(data)
	if v == "" || len(v) > 64 {
		return ""
	}
	return v
}

// cleanupOldVersions removes payload directories that are not the current version.
func cleanupOldVersions(baseDir, currentVersion string) {
	entries, err := os.ReadDir(baseDir)
	if err != nil {
		return
	}
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		name := entry.Name()
		// Skip special dirs and current version
		if name == currentVersion || name == ".download-tmp" {
			continue
		}
		// If it looks like a version dir (starts with digit), remove it
		if len(name) > 0 && name[0] >= '0' && name[0] <= '9' {
			os.RemoveAll(filepath.Join(baseDir, name))
		}
	}
}

func formatSize(bytes int64) string {
	if bytes <= 0 {
		return "unknown size"
	}
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMG"[exp])
}
