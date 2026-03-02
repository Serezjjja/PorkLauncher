package bootstrap

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// ExtractZip extracts a zip archive into destDir.
// It validates paths to prevent zip-slip attacks.
// Returns a list of extracted file paths (relative to destDir).
func ExtractZip(zipPath, destDir string) ([]string, error) {
	r, err := zip.OpenReader(zipPath)
	if err != nil {
		return nil, fmt.Errorf("open zip: %w", err)
	}
	defer r.Close()

	if err := os.MkdirAll(destDir, 0755); err != nil {
		return nil, fmt.Errorf("create destination dir: %w", err)
	}

	// Resolve destDir to absolute to detect zip-slip
	absDestDir, err := filepath.Abs(destDir)
	if err != nil {
		return nil, fmt.Errorf("resolve destination path: %w", err)
	}

	var extracted []string

	for _, f := range r.File {
		// Skip macOS resource forks
		if strings.Contains(f.Name, "__MACOSX") || strings.HasPrefix(filepath.Base(f.Name), "._") {
			continue
		}

		targetPath := filepath.Join(absDestDir, filepath.FromSlash(f.Name))

		// Zip-slip protection
		if !strings.HasPrefix(filepath.Clean(targetPath), absDestDir+string(os.PathSeparator)) && filepath.Clean(targetPath) != absDestDir {
			return nil, fmt.Errorf("zip-slip detected: %s escapes %s", f.Name, destDir)
		}

		if f.FileInfo().IsDir() {
			if err := os.MkdirAll(targetPath, 0755); err != nil {
				return nil, fmt.Errorf("create dir %s: %w", f.Name, err)
			}
			continue
		}

		// Create parent directories
		if err := os.MkdirAll(filepath.Dir(targetPath), 0755); err != nil {
			return nil, fmt.Errorf("create parent dir for %s: %w", f.Name, err)
		}

		if err := extractFile(f, targetPath); err != nil {
			return nil, fmt.Errorf("extract %s: %w", f.Name, err)
		}

		extracted = append(extracted, f.Name)
	}

	return extracted, nil
}

func extractFile(f *zip.File, targetPath string) error {
	src, err := f.Open()
	if err != nil {
		return fmt.Errorf("open zip entry: %w", err)
	}
	defer src.Close()

	dst, err := os.OpenFile(targetPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, f.Mode()|0644)
	if err != nil {
		return fmt.Errorf("create file: %w", err)
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		return fmt.Errorf("write file: %w", err)
	}

	return nil
}
