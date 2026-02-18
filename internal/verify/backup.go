package verify

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"HyLauncher/pkg/fileutil"
)

// BackupManager handles file backup and restoration
type BackupManager struct {
	backupDir string
	enabled   bool
}

// NewBackupManager creates a new backup manager
func NewBackupManager(backupDir string, enabled bool) *BackupManager {
	return &BackupManager{
		backupDir: backupDir,
		enabled:   enabled,
	}
}

// EnsureBackupDir creates the backup directory if it doesn't exist
func (bm *BackupManager) EnsureBackupDir() error {
	if !bm.enabled {
		return nil
	}
	return os.MkdirAll(bm.backupDir, 0755)
}

// CreateBackup creates a backup of the specified file
func (bm *BackupManager) CreateBackup(filePath string) (string, error) {
	if !bm.enabled {
		return "", nil
	}

	if err := bm.EnsureBackupDir(); err != nil {
		return "", &VerificationError{
			Op:      "create backup directory",
			Path:    bm.backupDir,
			Err:     err,
			Message: fmt.Sprintf("failed to create backup directory: %v", err),
		}
	}

	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return "", nil // Nothing to backup
	}

	// Generate backup filename with timestamp
	fileName := filepath.Base(filePath)
	timestamp := time.Now().Format("20060102_150405")
	backupName := fmt.Sprintf("%s.%s.bak", fileName, timestamp)
	backupPath := filepath.Join(bm.backupDir, backupName)

	// Handle subdirectories in backup
	subDir := filepath.Dir(filePath)
	if subDir != "." && subDir != "/" {
		backupPath = filepath.Join(bm.backupDir, subDir, backupName)
		if err := os.MkdirAll(filepath.Dir(backupPath), 0755); err != nil {
			return "", &VerificationError{
				Op:      "create backup subdirectory",
				Path:    filepath.Dir(backupPath),
				Err:     err,
				Message: fmt.Sprintf("failed to create backup subdirectory: %v", err),
			}
		}
	}

	// Copy file to backup location
	if err := fileutil.CopyFile(filePath, backupPath); err != nil {
		return "", &VerificationError{
			Op:      "backup file",
			Path:    filePath,
			Err:     err,
			Message: fmt.Sprintf("failed to backup file %s: %v", filePath, err),
		}
	}

	return backupPath, nil
}

// RestoreFile restores a file from its most recent backup
func (bm *BackupManager) RestoreFile(originalPath string) (string, error) {
	if !bm.enabled {
		return "", &VerificationError{
			Op:      "restore",
			Message: "backup manager is disabled",
		}
	}

	fileName := filepath.Base(originalPath)
	subDir := filepath.Dir(originalPath)

	// Determine backup search directory
	searchDir := bm.backupDir
	if subDir != "." && subDir != "/" {
		searchDir = filepath.Join(bm.backupDir, subDir)
	}

	// Find the most recent backup
	backups, err := bm.listBackups(searchDir, fileName)
	if err != nil {
		return "", &VerificationError{
			Op:      "list backups",
			Path:    searchDir,
			Err:     err,
			Message: fmt.Sprintf("failed to list backups: %v", err),
		}
	}

	if len(backups) == 0 {
		return "", &VerificationError{
			Op:      "restore",
			Path:    originalPath,
			Message: fmt.Sprintf("no backup found for %s", originalPath),
		}
	}

	// Sort by modification time (newest first) and pick the most recent
	mostRecent := backups[0]
	for _, backup := range backups {
		if backup.ModTime().After(mostRecent.ModTime()) {
			mostRecent = backup
		}
	}

	backupPath := filepath.Join(searchDir, mostRecent.Name())

	// Ensure target directory exists
	targetDir := filepath.Dir(originalPath)
	if err := os.MkdirAll(targetDir, 0755); err != nil {
		return "", &VerificationError{
			Op:      "create target directory",
			Path:    targetDir,
			Err:     err,
			Message: fmt.Sprintf("failed to create target directory: %v", err),
		}
	}

	// Restore the file
	if err := fileutil.CopyFile(backupPath, originalPath); err != nil {
		return "", &VerificationError{
			Op:      "restore file",
			Path:    originalPath,
			Err:     err,
			Message: fmt.Sprintf("failed to restore file from %s: %v", backupPath, err),
		}
	}

	return backupPath, nil
}

// listBackups returns all backup files for a given original filename
func (bm *BackupManager) listBackups(searchDir string, originalName string) ([]os.FileInfo, error) {
	entries, err := os.ReadDir(searchDir)
	if err != nil {
		if os.IsNotExist(err) {
			return []os.FileInfo{}, nil
		}
		return nil, err
	}

	var backups []os.FileInfo
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		// Check if this is a backup of the original file
		name := entry.Name()
		if len(name) > len(originalName)+4 &&
			name[:len(originalName)] == originalName &&
			name[len(name)-4:] == ".bak" {
			info, err := entry.Info()
			if err == nil {
				backups = append(backups, info)
			}
		}
	}

	return backups, nil
}

// ListAllBackups returns all backup files in the backup directory
func (bm *BackupManager) ListAllBackups() ([]string, error) {
	if !bm.enabled {
		return []string{}, nil
	}

	var backups []string
	err := filepath.Walk(bm.backupDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil // Continue walking
		}
		if !info.IsDir() && filepath.Ext(path) == ".bak" {
			relPath, _ := filepath.Rel(bm.backupDir, path)
			backups = append(backups, relPath)
		}
		return nil
	})

	return backups, err
}

// CleanupOldBackups removes backups older than the specified duration
func (bm *BackupManager) CleanupOldBackups(maxAge time.Duration) (int, error) {
	if !bm.enabled {
		return 0, nil
	}

	count := 0
	cutoff := time.Now().Add(-maxAge)

	err := filepath.Walk(bm.backupDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil // Continue walking
		}
		if info.IsDir() {
			return nil
		}
		if filepath.Ext(path) == ".bak" && info.ModTime().Before(cutoff) {
			if err := os.Remove(path); err == nil {
				count++
			}
		}
		return nil
	})

	return count, err
}

// GetBackupPath returns the configured backup directory
func (bm *BackupManager) GetBackupPath() string {
	return bm.backupDir
}
