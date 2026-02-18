package verify

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

// FileInfo represents expected file metadata in the manifest
type FileInfo struct {
	Size   int64  `json:"size"`
	SHA256 string `json:"sha256"`
}

// Manifest contains the expected file states for verification
type Manifest struct {
	Version   string              `json:"version"`
	CreatedAt time.Time           `json:"created_at"`
	Files     map[string]FileInfo `json:"files"`
	Ignore    []string            `json:"ignore,omitempty"`
}

// ManifestEntry represents a single file entry with its relative path
type ManifestEntry struct {
	RelativePath string
	Info         FileInfo
}

// LoadManifest loads a manifest from a JSON file
func LoadManifest(path string) (*Manifest, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, &VerificationError{
			Op:      "load manifest",
			Path:    path,
			Err:     err,
			Message: fmt.Sprintf("failed to load manifest from %s: %v", path, err),
		}
	}

	var manifest Manifest
	if err := json.Unmarshal(data, &manifest); err != nil {
		return nil, &VerificationError{
			Op:      "parse manifest",
			Path:    path,
			Err:     err,
			Message: fmt.Sprintf("failed to parse manifest JSON: %v", err),
		}
	}

	return &manifest, nil
}

// SaveManifest saves the manifest to a JSON file
func (m *Manifest) Save(path string) error {
	data, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		return &VerificationError{
			Op:      "marshal manifest",
			Err:     err,
			Message: fmt.Sprintf("failed to marshal manifest: %v", err),
		}
	}

	// Ensure directory exists
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return &VerificationError{
			Op:      "create manifest directory",
			Path:    dir,
			Err:     err,
			Message: fmt.Sprintf("failed to create manifest directory: %v", err),
		}
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		return &VerificationError{
			Op:      "save manifest",
			Path:    path,
			Err:     err,
			Message: fmt.Sprintf("failed to save manifest to %s: %v", path, err),
		}
	}

	return nil
}

// GetExpectedInfo returns the expected file info for a given relative path
func (m *Manifest) GetExpectedInfo(relativePath string) (FileInfo, bool) {
	// Normalize path separators for cross-platform compatibility
	normalizedPath := filepath.ToSlash(relativePath)
	info, exists := m.Files[normalizedPath]
	return info, exists
}

// GenerateManifest creates a new manifest by scanning a directory
func GenerateManifest(gameDir string, version string, ignoreList []string) (*Manifest, error) {
	manifest := &Manifest{
		Version:   version,
		CreatedAt: time.Now(),
		Files:     make(map[string]FileInfo),
		Ignore:    ignoreList,
	}

	// Key files to include in priority order
	keyFiles := []string{
		"Assets.zip",
		filepath.Join("Client", "HytaleClient.jar"),
		filepath.Join("Client", "HytaleClient.exe"),
		filepath.Join("config", "client.json"),
	}

	// Add key files first
	for _, relPath := range keyFiles {
		fullPath := filepath.Join(gameDir, relPath)
		if _, err := os.Stat(fullPath); err == nil {
			info, err := calculateFileInfo(fullPath)
			if err != nil {
				continue // Skip files we can't read
			}
			manifest.Files[filepath.ToSlash(relPath)] = info
		}
	}

	// Scan Client/lib for JAR files
	libDir := filepath.Join(gameDir, "Client", "lib")
	if entries, err := os.ReadDir(libDir); err == nil {
		for _, entry := range entries {
			if entry.IsDir() {
				continue
			}
			if filepath.Ext(entry.Name()) == ".jar" {
				relPath := filepath.Join("Client", "lib", entry.Name())
				fullPath := filepath.Join(gameDir, relPath)
				info, err := calculateFileInfo(fullPath)
				if err != nil {
					continue
				}
				manifest.Files[filepath.ToSlash(relPath)] = info
			}
		}
	}

	// Scan patches directory for .pw files
	patchesDir := filepath.Join(gameDir, "patches")
	if entries, err := os.ReadDir(patchesDir); err == nil {
		for _, entry := range entries {
			if entry.IsDir() {
				continue
			}
			if filepath.Ext(entry.Name()) == ".pw" {
				relPath := filepath.Join("patches", entry.Name())
				fullPath := filepath.Join(gameDir, relPath)
				info, err := calculateFileInfo(fullPath)
				if err != nil {
					continue
				}
				manifest.Files[filepath.ToSlash(relPath)] = info
			}
		}
	}

	return manifest, nil
}

// calculateFileInfo computes size and SHA-256 hash for a file
func calculateFileInfo(filePath string) (FileInfo, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return FileInfo{}, err
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		return FileInfo{}, err
	}

	hasher := sha256.New()
	if _, err := io.Copy(hasher, file); err != nil {
		return FileInfo{}, err
	}

	return FileInfo{
		Size:   stat.Size(),
		SHA256: hex.EncodeToString(hasher.Sum(nil)),
	}, nil
}

// GetDefaultManifestPath returns the default path for a manifest file
func GetDefaultManifestPath(gameDir string, version string) string {
	// Store manifest in the launcher directory, not game directory
	// This prevents it from being overwritten during updates
	launcherDir := getLauncherDir()
	return filepath.Join(launcherDir, "manifests", fmt.Sprintf("manifest_%s.json", version))
}

// getLauncherDir returns the launcher application directory
func getLauncherDir() string {
	home, _ := os.UserHomeDir()
	switch runtime.GOOS {
	case "windows":
		return filepath.Join(home, "AppData", "Local", "HyLauncher")
	case "darwin":
		return filepath.Join(home, "Library", "Application Support", "HyLauncher")
	default:
		return filepath.Join(home, ".hylauncher")
	}
}

// IsIgnored checks if a file path matches any pattern in the ignore list
func (m *Manifest) IsIgnored(relativePath string) bool {
	for _, pattern := range m.Ignore {
		matched, _ := filepath.Match(pattern, relativePath)
		if matched {
			return true
		}
		// Also check with ToSlash for cross-platform compatibility
		matched, _ = filepath.Match(pattern, filepath.ToSlash(relativePath))
		if matched {
			return true
		}
	}
	return false
}
