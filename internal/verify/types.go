// Package verify provides file integrity verification for Hytale game files.
// It checks files for existence, size, and SHA-256 hash to detect corruption or modifications.
package verify

import (
	"time"
)

// VerificationStatus represents the overall status of a verification operation
type VerificationStatus string

const (
	StatusOK      VerificationStatus = "OK"
	StatusWarning VerificationStatus = "Warning"
	StatusFailed  VerificationStatus = "Failed"
	StatusSkipped VerificationStatus = "Skipped"
)

// FileStatus represents the verification result for a single file
type FileStatus struct {
	Path         string             `json:"path"`
	Exists       bool               `json:"exists"`
	Size         int64              `json:"size"`
	ExpectedSize int64              `json:"expected_size"`
	Hash         string             `json:"hash"`
	ExpectedHash string             `json:"expected_hash"`
	Match        bool               `json:"match"`
	Status       VerificationStatus `json:"status"`
	Message      string             `json:"message"`
}

// IsCritical returns true if the file verification failed critically
func (fs FileStatus) IsCritical() bool {
	return fs.Status == StatusFailed
}

// NeedsAttention returns true if the file has issues that need to be addressed
func (fs FileStatus) NeedsAttention() bool {
	return fs.Status == StatusFailed || fs.Status == StatusWarning
}

// Report contains the complete verification results
type Report struct {
	Version       string             `json:"version"`
	GameDir       string             `json:"game_dir"`
	Timestamp     time.Time          `json:"timestamp"`
	Files         []FileStatus       `json:"files"`
	OverallStatus VerificationStatus `json:"overall_status"`
	Summary       Summary            `json:"summary"`
	LogFile       string             `json:"log_file,omitempty"`
}

// Summary provides a quick overview of verification results
type Summary struct {
	TotalFiles    int `json:"total_files"`
	Passed        int `json:"passed"`
	Failed        int `json:"failed"`
	Warnings      int `json:"warnings"`
	Skipped       int `json:"skipped"`
	MissingFiles  int `json:"missing_files"`
	ModifiedFiles int `json:"modified_files"`
}

// Options configures the verification behavior
type Options struct {
	// GameDir is the path to the game installation directory
	GameDir string
	// Version is the game version to verify against
	Version string
	// ManifestPath is the path to the verification manifest (optional)
	ManifestPath string
	// SkipVerify disables verification when set to true
	SkipVerify bool
	// CreateBackups creates backups of modified files before any action
	CreateBackups bool
	// BackupDir is the directory to store backups (default: gameDir/.backups)
	BackupDir string
	// IgnoreList contains file patterns to ignore during verification
	IgnoreList []string
	// ProgressCallback is called periodically during hash calculation
	ProgressCallback func(current, total int64, fileName string)
	// ProgressInterval is the minimum bytes between progress updates (default: 100MB)
	ProgressInterval int64
}

// DefaultOptions returns options with sensible defaults
func DefaultOptions() Options {
	return Options{
		CreateBackups:    true,
		ProgressInterval: 100 * 1024 * 1024, // 100 MB
	}
}

// VerificationError represents an error during verification
type VerificationError struct {
	Op      string // operation
	Path    string // file path
	Err     error  // underlying error
	Message string // human-readable message
}

func (e *VerificationError) Error() string {
	if e.Message != "" {
		return e.Message
	}
	if e.Path != "" {
		return e.Op + " " + e.Path + ": " + e.Err.Error()
	}
	return e.Op + ": " + e.Err.Error()
}

func (e *VerificationError) Unwrap() error {
	return e.Err
}
