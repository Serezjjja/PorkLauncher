package verify

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
)

// HashCalculator provides progress-aware file hashing
type HashCalculator struct {
	progressCallback func(current, total int64, fileName string)
	progressInterval int64
}

// NewHashCalculator creates a new hash calculator with the given callback
func NewHashCalculator(callback func(current, total int64, fileName string), interval int64) *HashCalculator {
	if interval <= 0 {
		interval = 100 * 1024 * 1024 // 100 MB default
	}
	return &HashCalculator{
		progressCallback: callback,
		progressInterval: interval,
	}
}

// CalculateHash computes the SHA-256 hash of a file with progress reporting
func (hc *HashCalculator) CalculateHash(filePath string) (string, int64, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", 0, &VerificationError{
			Op:      "open file for hashing",
			Path:    filePath,
			Err:     err,
			Message: fmt.Sprintf("failed to open file %s: %v", filePath, err),
		}
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		return "", 0, &VerificationError{
			Op:      "stat file",
			Path:    filePath,
			Err:     err,
			Message: fmt.Sprintf("failed to stat file %s: %v", filePath, err),
		}
	}

	size := stat.Size()
	hasher := sha256.New()

	if hc.progressCallback != nil && size > hc.progressInterval {
		// Use progress-aware reading for large files
		reader := &progressReader{
			reader:           file,
			total:            size,
			progressCallback: hc.progressCallback,
			progressInterval: hc.progressInterval,
			fileName:         stat.Name(),
		}
		if _, err := io.Copy(hasher, reader); err != nil {
			return "", 0, &VerificationError{
				Op:      "hash file",
				Path:    filePath,
				Err:     err,
				Message: fmt.Sprintf("failed to hash file %s: %v", filePath, err),
			}
		}
	} else {
		// Simple hash for small files
		if _, err := io.Copy(hasher, file); err != nil {
			return "", 0, &VerificationError{
				Op:      "hash file",
				Path:    filePath,
				Err:     err,
				Message: fmt.Sprintf("failed to hash file %s: %v", filePath, err),
			}
		}
	}

	hash := hex.EncodeToString(hasher.Sum(nil))
	return hash, size, nil
}

// progressReader wraps an io.Reader to report progress
type progressReader struct {
	reader           io.Reader
	total            int64
	progressCallback func(current, total int64, fileName string)
	progressInterval int64
	fileName         string
	readSoFar        int64
	lastReport       int64
}

func (pr *progressReader) Read(p []byte) (n int, err error) {
	n, err = pr.reader.Read(p)
	if n > 0 {
		pr.readSoFar += int64(n)
		// Report progress every progressInterval bytes
		if pr.readSoFar-pr.lastReport >= pr.progressInterval {
			pr.progressCallback(pr.readSoFar, pr.total, pr.fileName)
			pr.lastReport = pr.readSoFar
		}
	}
	// Report final progress when done
	if err == io.EOF && pr.lastReport < pr.total {
		pr.progressCallback(pr.total, pr.total, pr.fileName)
	}
	return n, err
}

// QuickVerify performs a quick verification using file size only
func QuickVerify(filePath string, expectedSize int64) (bool, int64, error) {
	stat, err := os.Stat(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return false, 0, nil
		}
		return false, 0, err
	}
	return stat.Size() == expectedSize, stat.Size(), nil
}

// CalculateHashSimple computes SHA-256 hash without progress reporting
func CalculateHashSimple(filePath string) (string, int64, error) {
	calculator := NewHashCalculator(nil, 0)
	return calculator.CalculateHash(filePath)
}
