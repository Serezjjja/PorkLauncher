package bootstrap

import (
	"crypto/ed25519"
	"encoding/hex"
	"fmt"
	"runtime"
)

const (
	// BootstrapVersion is the version of the bootstrap binary.
	// Bump this when bootstrap protocol changes.
	BootstrapVersion = "1.0.0"

	// PayloadDirName is the directory name inside the app data folder
	// where downloaded payloads are stored.
	PayloadDirName = "payload"

	// MetadataFileName is the name of the metadata JSON file on the update server.
	MetadataFileName = "metadata.json"

	// SignatureFileName is the name of the Ed25519 signature file.
	SignatureFileName = "metadata.json.sig"
)

// Build-time variables injected via -ldflags.
// Example: -ldflags="-X 'HyLauncher/internal/bootstrap.UpdateServerURL=https://example.com/updates'"
var (
	// UpdateServerURL is the base URL for the update server.
	// Must be set at build time. Example: https://example.com/updates
	// The bootstrap will fetch {UpdateServerURL}/metadata.json and {UpdateServerURL}/metadata.json.sig
	UpdateServerURL = ""

	// Ed25519PublicKeyHex is the hex-encoded Ed25519 public key (64 hex chars = 32 bytes).
	// Used to verify the digital signature of metadata.json.
	// Generate a keypair with: go run ./tools/sign-payload -generate
	Ed25519PublicKeyHex = ""

	// FallbackDownloadURL is shown to the user if automatic download fails.
	// Typically a GitHub releases page.
	FallbackDownloadURL = "https://github.com/Serezjjja/PorkLauncher/releases"
)

// PayloadExecutableName returns the platform-specific payload executable name.
func PayloadExecutableName() string {
	if runtime.GOOS == "windows" {
		return "HyLauncher.exe"
	}
	return "HyLauncher"
}

// PlatformKey returns the platform key used in metadata.json (e.g. "windows/amd64").
func PlatformKey() string {
	return runtime.GOOS + "/" + runtime.GOARCH
}

// GetPublicKey parses and returns the Ed25519 public key from the hex-encoded build variable.
func GetPublicKey() (ed25519.PublicKey, error) {
	if Ed25519PublicKeyHex == "" {
		return nil, fmt.Errorf("Ed25519 public key not embedded (build without -ldflags?)")
	}

	keyBytes, err := hex.DecodeString(Ed25519PublicKeyHex)
	if err != nil {
		return nil, fmt.Errorf("invalid Ed25519 public key hex: %w", err)
	}

	if len(keyBytes) != ed25519.PublicKeySize {
		return nil, fmt.Errorf("Ed25519 public key wrong size: got %d, want %d", len(keyBytes), ed25519.PublicKeySize)
	}

	return ed25519.PublicKey(keyBytes), nil
}

// GetUpdateServerURL returns the update server URL, validating it is set.
func GetUpdateServerURL() (string, error) {
	if UpdateServerURL == "" {
		return "", fmt.Errorf("update server URL not configured (build without -ldflags?)")
	}
	return UpdateServerURL, nil
}
