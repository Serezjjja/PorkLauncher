# File Verification Module

This package provides file integrity verification for Hytale game files.

## Overview

The verification module checks game files for:
- **Existence**: Files must be present
- **Size**: Files must match expected size
- **SHA-256 Hash**: Files must not be corrupted or modified

## File Structure

```
internal/verify/
├── types.go      # Core types: FileStatus, Report, Options
├── manifest.go   # Manifest loading and generation
├── hash.go       # Progress-aware hash calculation
├── backup.go     # Backup and restore functionality
├── verify.go     # Main verification logic
├── logger.go     # Verification logging
└── README.md     # This file
```

## Usage

### Basic Verification

```go
import "HyLauncher/internal/verify"

// Simple verification
report, err := verify.VerifyGameFiles(gameDir, version)
if err != nil {
    log.Fatal(err)
}

// Check results
if report.OverallStatus == verify.StatusOK {
    fmt.Println("All files verified!")
} else if report.OverallStatus == verify.StatusWarning {
    fmt.Println("Some files modified")
} else {
    fmt.Println("Verification failed!")
}
```

### Advanced Verification with Options

```go
options := verify.Options{
    GameDir:       "/path/to/game",
    Version:       "2026.02.001",
    ManifestPath:  "/path/to/manifest.json",
    CreateBackups: true,
    IgnoreList:    []string{"*.log", "cache/*"},
    ProgressCallback: func(current, total int64, fileName string) {
        percent := float64(current) / float64(total) * 100
        fmt.Printf("\r%s: %.1f%%", fileName, percent)
    },
}

report, err := verify.VerifyWithOptions(options)
```

### Generate Manifest

```go
manifest, err := verify.GenerateManifest(gameDir, version, ignoreList)
if err != nil {
    log.Fatal(err)
}

// Save manifest
err = manifest.Save("/path/to/manifest.json")
```

### Restore File from Backup

```go
backupPath, err := verify.RestoreFile(filePath, backupDir)
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Restored from: %s\n", backupPath)
```

## CLI Tool

The `cmd/verify` directory contains a CLI tool for verification operations:

```bash
# Verify game files
go run cmd/verify/main.go verify \
    --game-dir="C:\Users\User\AppData\Roaming\Hytale\install\release\package\game\latest" \
    --version=2026.02.001 \
    --v

# Generate manifest from clean installation
go run cmd/verify/main.go generate \
    --game-dir="/path/to/clean/install" \
    --version=2026.02.001 \
    --output=manifest.json

# Restore a file
go run cmd/verify/main.go restore \
    --file="Client/HytaleClient.jar" \
    --game-dir="/path/to/game"

# Create backup
go run cmd/verify/main.go backup \
    --file="Assets.zip" \
    --game-dir="/path/to/game"
```

## Manifest Format

```json
{
  "version": "2026.02.001",
  "created_at": "2026-02-17T20:30:00Z",
  "files": {
    "Assets.zip": {
      "size": 4123456789,
      "sha256": "a1b2c3d4e5f6789012345678901234567890abcdef1234567890abcdef123456"
    },
    "Client/HytaleClient.jar": {
      "size": 84567890,
      "sha256": "b2c3d4e5f6789012345678901234567890abcdef1234567890abcdef123456a1"
    }
  },
  "ignore": [
    "*.log",
    "logs/*"
  ]
}
```

## Environment Variables

- `HYTALE_SKIP_VERIFY=1` - Skip verification entirely

## Exit Codes

- `0` - Success (all files verified)
- `1` - Error (manifest not found, etc.)
- `2` - Verification failed (missing files)
- `3` - Verification warnings (modified files)
