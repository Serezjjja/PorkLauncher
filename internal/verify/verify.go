package verify

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// Verifier handles file verification operations
type Verifier struct {
	options Options
}

// NewVerifier creates a new verifier with the given options
func NewVerifier(options Options) *Verifier {
	if options.ProgressInterval == 0 {
		options.ProgressInterval = 100 * 1024 * 1024 // 100 MB
	}
	return &Verifier{
		options: options,
	}
}

// VerifyGameFiles performs a complete verification of game files
func VerifyGameFiles(gameDir string, version string) (*Report, error) {
	options := DefaultOptions()
	options.GameDir = gameDir
	options.Version = version
	return VerifyWithOptions(options)
}

// VerifyWithOptions performs verification with custom options
func VerifyWithOptions(options Options) (*Report, error) {
	// Check if verification should be skipped
	if shouldSkipVerification(options) {
		return createSkippedReport(options), nil
	}

	verifier := NewVerifier(options)
	return verifier.runVerification()
}

// shouldSkipVerification checks if verification should be skipped
func shouldSkipVerification(options Options) bool {
	if options.SkipVerify {
		return true
	}
	// Check environment variable
	if os.Getenv("HYTALE_SKIP_VERIFY") == "1" || strings.ToLower(os.Getenv("HYTALE_SKIP_VERIFY")) == "true" {
		return true
	}
	return false
}

// createSkippedReport creates a report indicating verification was skipped
func createSkippedReport(options Options) *Report {
	return &Report{
		Version:       options.Version,
		GameDir:       options.GameDir,
		Timestamp:     time.Now(),
		Files:         []FileStatus{},
		OverallStatus: StatusSkipped,
		Summary: Summary{
			TotalFiles: 0,
			Skipped:    1,
		},
	}
}

// runVerification performs the actual verification
func (v *Verifier) runVerification() (*Report, error) {
	report := &Report{
		Version:   v.options.Version,
		GameDir:   v.options.GameDir,
		Timestamp: time.Now(),
		Files:     []FileStatus{},
	}

	// Starting file verification

	// Load or generate manifest
	manifest, err := v.loadManifest()
	if err != nil {
		return nil, err
	}

	// Setup backup manager
	backupDir := v.options.BackupDir
	if backupDir == "" {
		backupDir = filepath.Join(v.options.GameDir, ".backups")
	}
	backupManager := NewBackupManager(backupDir, v.options.CreateBackups)

	// Create hash calculator with progress callback
	hashCalculator := NewHashCalculator(v.options.ProgressCallback, v.options.ProgressInterval)

	// Verify each file in the manifest
	for relPath, expectedInfo := range manifest.Files {
		// Check if file should be ignored
		if manifest.IsIgnored(relPath) {
			// Skipping ignored file
			continue
		}

		// Check custom ignore list
		if v.isIgnored(relPath) {
			// Skipping file (custom ignore list)
			continue
		}

		fullPath := filepath.Join(v.options.GameDir, relPath)
		status := v.verifyFile(fullPath, relPath, expectedInfo, hashCalculator, backupManager)
		report.Files = append(report.Files, status)

		// Verified file logged
	}

	// Calculate summary
	report.Summary = v.calculateSummary(report.Files)
	report.OverallStatus = v.determineOverallStatus(report.Files)

	// Verification complete logged

	return report, nil
}

// loadManifest loads the verification manifest
func (v *Verifier) loadManifest() (*Manifest, error) {
	// Try to load from specified path first
	if v.options.ManifestPath != "" {
		manifest, err := LoadManifest(v.options.ManifestPath)
		if err == nil {
			return manifest, nil
		}
		// Failed to load specified manifest, trying default location
	}

	// Try default location
	defaultPath := GetDefaultManifestPath(v.options.GameDir, v.options.Version)
	manifest, err := LoadManifest(defaultPath)
	if err == nil {
		return manifest, nil
	}

	// Try to find any manifest for this version
	launcherDir := getLauncherDir()
	manifestDir := filepath.Join(launcherDir, "manifests")
	if entries, err := os.ReadDir(manifestDir); err == nil {
		for _, entry := range entries {
			if !entry.IsDir() && strings.HasPrefix(entry.Name(), "manifest_") {
				path := filepath.Join(manifestDir, entry.Name())
				if m, err := LoadManifest(path); err == nil {
					// Found manifest
					return m, nil
				}
			}
		}
	}

	return nil, &VerificationError{
		Op:      "load manifest",
		Message: fmt.Sprintf("no manifest found for version %s", v.options.Version),
	}
}

// isIgnored checks if a file matches any pattern in the ignore list
func (v *Verifier) isIgnored(relPath string) bool {
	for _, pattern := range v.options.IgnoreList {
		matched, _ := filepath.Match(pattern, relPath)
		if matched {
			return true
		}
		matched, _ = filepath.Match(pattern, filepath.ToSlash(relPath))
		if matched {
			return true
		}
	}
	return false
}

// verifyFile verifies a single file
func (v *Verifier) verifyFile(fullPath, relPath string, expected FileInfo, calculator *HashCalculator, backupManager *BackupManager) FileStatus {
	status := FileStatus{
		Path:         relPath,
		ExpectedSize: expected.Size,
		ExpectedHash: expected.SHA256,
	}

	// Check if file exists
	stat, err := os.Stat(fullPath)
	if err != nil {
		if os.IsNotExist(err) {
			status.Exists = false
			status.Status = StatusFailed
			status.Message = "File missing - reinstallation recommended"
			return status
		}
		status.Status = StatusFailed
		status.Message = fmt.Sprintf("Cannot access file: %v", err)
		return status
	}

	status.Exists = true
	status.Size = stat.Size()

	// Check size
	if stat.Size() != expected.Size {
		status.Status = StatusFailed
		status.Message = fmt.Sprintf("Size mismatch: expected %d, got %d", expected.Size, stat.Size())
		return status
	}

	// Calculate hash
	hash, _, err := calculator.CalculateHash(fullPath)
	if err != nil {
		status.Status = StatusFailed
		status.Message = fmt.Sprintf("Failed to calculate hash: %v", err)
		return status
	}

	status.Hash = hash
	status.Match = hash == expected.SHA256

	if status.Match {
		status.Status = StatusOK
		status.Message = "File verified successfully"
	} else {
		status.Status = StatusWarning
		status.Message = "File modified or corrupted"

		// Create backup of modified file
		if v.options.CreateBackups {
			backupPath, err := backupManager.CreateBackup(fullPath)
			if err != nil {
				// Failed to create backup
			} else if backupPath != "" {
				// Created backup
			}
		}
	}

	return status
}

// calculateSummary calculates the verification summary
func (v *Verifier) calculateSummary(files []FileStatus) Summary {
	summary := Summary{
		TotalFiles: len(files),
	}

	for _, f := range files {
		switch f.Status {
		case StatusOK:
			summary.Passed++
		case StatusFailed:
			summary.Failed++
			if !f.Exists {
				summary.MissingFiles++
			}
		case StatusWarning:
			summary.Warnings++
			if f.Exists && !f.Match {
				summary.ModifiedFiles++
			}
		case StatusSkipped:
			summary.Skipped++
		}
	}

	return summary
}

// determineOverallStatus determines the overall verification status
func (v *Verifier) determineOverallStatus(files []FileStatus) VerificationStatus {
	hasFailed := false
	hasWarning := false

	for _, f := range files {
		switch f.Status {
		case StatusFailed:
			hasFailed = true
		case StatusWarning:
			hasWarning = true
		}
	}

	if hasFailed {
		return StatusFailed
	}
	if hasWarning {
		return StatusWarning
	}
	return StatusOK
}

// RestoreFile restores a file from backup
func RestoreFile(filePath string, backupDir string) (string, error) {
	if backupDir == "" {
		backupDir = filepath.Join(filepath.Dir(filePath), ".backups")
	}

	backupManager := NewBackupManager(backupDir, true)
	backupPath, err := backupManager.RestoreFile(filePath)
	if err != nil {
		return "", err
	}

	return backupPath, nil
}

// GenerateAndSaveManifest creates a manifest from an existing installation and saves it
func GenerateAndSaveManifest(gameDir string, version string, outputPath string, ignoreList []string) error {
	manifest, err := GenerateManifest(gameDir, version, ignoreList)
	if err != nil {
		return err
	}

	if outputPath == "" {
		outputPath = GetDefaultManifestPath(gameDir, version)
	}

	return manifest.Save(outputPath)
}
