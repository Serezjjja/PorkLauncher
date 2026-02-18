package verify

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// LogEntry represents a single log entry
type LogEntry struct {
	Timestamp time.Time `json:"timestamp"`
	Level     string    `json:"level"`
	Message   string    `json:"message"`
}

// VerifyLogger handles logging for verification operations
type VerifyLogger struct {
	mu         sync.Mutex
	logFile    *os.File
	entries    []LogEntry
	logPath    string
	maxEntries int
}

// NewVerifyLogger creates a new verification logger
func NewVerifyLogger(logDir string) (*VerifyLogger, error) {
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create log directory: %w", err)
	}

	timestamp := time.Now().Format("20060102_150405")
	logPath := filepath.Join(logDir, fmt.Sprintf("verify_%s.log", timestamp))

	file, err := os.Create(logPath)
	if err != nil {
		return nil, fmt.Errorf("failed to create log file: %w", err)
	}

	return &VerifyLogger{
		logFile:    file,
		logPath:    logPath,
		entries:    make([]LogEntry, 0),
		maxEntries: 10000,
	}, nil
}

// Log writes a log entry
func (vl *VerifyLogger) Log(level, message string) {
	vl.mu.Lock()
	defer vl.mu.Unlock()

	entry := LogEntry{
		Timestamp: time.Now(),
		Level:     level,
		Message:   message,
	}

	vl.entries = append(vl.entries, entry)

	// Write to file immediately
	fmt.Fprintf(vl.logFile, "[%s] %s: %s\n", entry.Timestamp.Format("2006-01-02 15:04:05"), level, message)

	// Flush to ensure it's written
	vl.logFile.Sync()
}

// Info logs an info message
func (vl *VerifyLogger) Info(format string, args ...interface{}) {
	vl.Log("INFO", fmt.Sprintf(format, args...))
}

// Warn logs a warning message
func (vl *VerifyLogger) Warn(format string, args ...interface{}) {
	vl.Log("WARN", fmt.Sprintf(format, args...))
}

// Error logs an error message
func (vl *VerifyLogger) Error(format string, args ...interface{}) {
	vl.Log("ERROR", fmt.Sprintf(format, args...))
}

// Debug logs a debug message
func (vl *VerifyLogger) Debug(format string, args ...interface{}) {
	vl.Log("DEBUG", fmt.Sprintf(format, args...))
}

// GetLogPath returns the path to the log file
func (vl *VerifyLogger) GetLogPath() string {
	return vl.logPath
}

// Close closes the log file
func (vl *VerifyLogger) Close() error {
	vl.mu.Lock()
	defer vl.mu.Unlock()
	return vl.logFile.Close()
}

// GetEntries returns all log entries
func (vl *VerifyLogger) GetEntries() []LogEntry {
	vl.mu.Lock()
	defer vl.mu.Unlock()
	result := make([]LogEntry, len(vl.entries))
	copy(result, vl.entries)
	return result
}

// GetDefaultLogDir returns the default log directory
func GetDefaultLogDir() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, "AppData", "Local", "HyLauncher", "logs")
}

// CleanupOldLogs removes log files older than the specified duration
func CleanupOldLogs(logDir string, maxAge time.Duration) (int, error) {
	entries, err := os.ReadDir(logDir)
	if err != nil {
		if os.IsNotExist(err) {
			return 0, nil
		}
		return 0, err
	}

	count := 0
	cutoff := time.Now().Add(-maxAge)

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		if !entry.IsDir() && len(entry.Name()) > 4 && entry.Name()[len(entry.Name())-4:] == ".log" {
			info, err := entry.Info()
			if err != nil {
				continue
			}
			if info.ModTime().Before(cutoff) {
				path := filepath.Join(logDir, entry.Name())
				if err := os.Remove(path); err == nil {
					count++
				}
			}
		}
	}

	return count, nil
}
