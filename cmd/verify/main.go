// Command verify is a CLI tool for Hytale file verification
// Usage: go run cmd/verify/main.go [command] [flags]
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"HyLauncher/internal/verify"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "verify":
		verifyCmd(os.Args[2:])
	case "generate":
		generateCmd(os.Args[2:])
	case "restore":
		restoreCmd(os.Args[2:])
	case "backup":
		backupCmd(os.Args[2:])
	case "help", "-h", "--help":
		printUsage()
	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", command)
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println(`Hytale File Verification Tool

Usage: verify <command> [flags]

Commands:
  verify    Verify game files against a manifest
  generate  Generate a manifest from an existing installation
  restore   Restore a file from backup
  backup    Create a backup of specified files
  help      Show this help message

Examples:
  # Verify game files
  verify verify --game-dir="C:\Users\User\AppData\Roaming\Hytale\install\release\package\game\latest" --version=2026.02.001

  # Generate manifest from clean installation
  verify generate --game-dir="C:\HytaleClean" --version=2026.02.001 --output=manifest_2026.02.001.json

  # Restore a modified file
  verify restore --file="Client/HytaleClient.jar" --game-dir="..."

Flags for 'verify':`)
	flag.PrintDefaults()
}

func verifyCmd(args []string) {
	fs := flag.NewFlagSet("verify", flag.ExitOnError)
	gameDir := fs.String("game-dir", getDefaultGameDir(), "Path to game installation directory")
	version := fs.String("version", "", "Game version to verify against (required)")
	manifestPath := fs.String("manifest", "", "Path to manifest file (optional, auto-detected if not specified)")
	skipVerify := fs.Bool("skip-verify", false, "Skip verification (or set HYTALE_SKIP_VERIFY=1)")
	noBackup := fs.Bool("no-backup", false, "Disable automatic backup of modified files")
	backupDir := fs.String("backup-dir", "", "Directory for backups (default: gameDir/.backups)")
	ignore := fs.String("ignore", "", "Comma-separated list of file patterns to ignore")
	jsonOutput := fs.Bool("json", false, "Output results as JSON")
	verbose := fs.Bool("v", false, "Verbose output with progress")

	if err := fs.Parse(args); err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing flags: %v\n", err)
		os.Exit(1)
	}

	if *version == "" {
		fmt.Fprintf(os.Stderr, "Error: --version is required\n")
		fs.Usage()
		os.Exit(1)
	}

	// Expand environment variables in paths
	*gameDir = os.ExpandEnv(*gameDir)
	if *manifestPath != "" {
		*manifestPath = os.ExpandEnv(*manifestPath)
	}
	if *backupDir != "" {
		*backupDir = os.ExpandEnv(*backupDir)
	}

	// Parse ignore list
	var ignoreList []string
	if *ignore != "" {
		ignoreList = strings.Split(*ignore, ",")
		for i := range ignoreList {
			ignoreList[i] = strings.TrimSpace(ignoreList[i])
		}
	}

	// Setup progress callback
	var progressCallback func(current, total int64, fileName string)
	if *verbose {
		progressCallback = func(current, total int64, fileName string) {
			percent := float64(current) / float64(total) * 100
			fmt.Printf("\r  Hashing %s: %.1f%% (%s / %s)",
				fileName,
				percent,
				formatBytes(current),
				formatBytes(total))
			if current == total {
				fmt.Println() // New line when complete
			}
		}
	}

	options := verify.Options{
		GameDir:          *gameDir,
		Version:          *version,
		ManifestPath:     *manifestPath,
		SkipVerify:       *skipVerify,
		CreateBackups:    !*noBackup,
		BackupDir:        *backupDir,
		IgnoreList:       ignoreList,
		ProgressCallback: progressCallback,
		ProgressInterval: 100 * 1024 * 1024, // 100 MB
	}

	if !*jsonOutput {
		fmt.Printf("Verifying game files...\n")
		fmt.Printf("Game directory: %s\n", *gameDir)
		fmt.Printf("Version: %s\n", *version)
		if *verbose {
			fmt.Println()
		}
	}

	report, err := verify.VerifyWithOptions(options)
	if err != nil {
		if *jsonOutput {
			errorJSON, _ := json.Marshal(map[string]string{"error": err.Error()})
			fmt.Println(string(errorJSON))
		} else {
			fmt.Fprintf(os.Stderr, "Verification failed: %v\n", err)
		}
		os.Exit(1)
	}

	if *jsonOutput {
		output, _ := json.MarshalIndent(report, "", "  ")
		fmt.Println(string(output))
	} else {
		printReport(report)
	}

	// Exit with error code if verification failed
	if report.OverallStatus == verify.StatusFailed {
		os.Exit(2)
	}
	if report.OverallStatus == verify.StatusWarning {
		os.Exit(3)
	}
}

func generateCmd(args []string) {
	fs := flag.NewFlagSet("generate", flag.ExitOnError)
	gameDir := fs.String("game-dir", "", "Path to clean game installation (required)")
	version := fs.String("version", "", "Game version (required)")
	output := fs.String("output", "", "Output manifest file path (default: launcher_dir/manifests/manifest_<version>.json)")
	ignore := fs.String("ignore", "", "Comma-separated list of file patterns to ignore")

	if err := fs.Parse(args); err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing flags: %v\n", err)
		os.Exit(1)
	}

	if *gameDir == "" || *version == "" {
		fmt.Fprintf(os.Stderr, "Error: --game-dir and --version are required\n")
		fs.Usage()
		os.Exit(1)
	}

	*gameDir = os.ExpandEnv(*gameDir)
	if *output != "" {
		*output = os.ExpandEnv(*output)
	}

	// Parse ignore list
	var ignoreList []string
	if *ignore != "" {
		ignoreList = strings.Split(*ignore, ",")
		for i := range ignoreList {
			ignoreList[i] = strings.TrimSpace(ignoreList[i])
		}
	}

	fmt.Printf("Generating manifest from: %s\n", *gameDir)
	fmt.Printf("Version: %s\n", *version)
	fmt.Println("Scanning files...")

	manifest, err := verify.GenerateManifest(*gameDir, *version, ignoreList)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to generate manifest: %v\n", err)
		os.Exit(1)
	}

	// Determine output path
	outputPath := *output
	if outputPath == "" {
		outputPath = verify.GetDefaultManifestPath(*gameDir, *version)
	}

	// Ensure output directory exists
	if err := os.MkdirAll(filepath.Dir(outputPath), 0755); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create output directory: %v\n", err)
		os.Exit(1)
	}

	if err := manifest.Save(outputPath); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to save manifest: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Manifest saved to: %s\n", outputPath)
	fmt.Printf("Files included: %d\n", len(manifest.Files))
}

func restoreCmd(args []string) {
	fs := flag.NewFlagSet("restore", flag.ExitOnError)
	file := fs.String("file", "", "Relative path to file to restore (required)")
	gameDir := fs.String("game-dir", getDefaultGameDir(), "Path to game installation directory")
	backupDir := fs.String("backup-dir", "", "Backup directory (default: gameDir/.backups)")

	if err := fs.Parse(args); err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing flags: %v\n", err)
		os.Exit(1)
	}

	if *file == "" {
		fmt.Fprintf(os.Stderr, "Error: --file is required\n")
		fs.Usage()
		os.Exit(1)
	}

	*gameDir = os.ExpandEnv(*gameDir)
	if *backupDir != "" {
		*backupDir = os.ExpandEnv(*backupDir)
	}

	fullPath := filepath.Join(*gameDir, *file)
	if *backupDir == "" {
		*backupDir = filepath.Join(*gameDir, ".backups")
	}

	fmt.Printf("Restoring: %s\n", *file)
	fmt.Printf("From backup directory: %s\n", *backupDir)

	backupPath, err := verify.RestoreFile(fullPath, *backupDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Restore failed: %v\n", err)
		os.Exit(1)
	}

	if backupPath == "" {
		fmt.Println("No backup found for this file")
	} else {
		fmt.Printf("Successfully restored from: %s\n", backupPath)
	}
}

func backupCmd(args []string) {
	fs := flag.NewFlagSet("backup", flag.ExitOnError)
	file := fs.String("file", "", "Relative path to file to backup (required)")
	gameDir := fs.String("game-dir", getDefaultGameDir(), "Path to game installation directory")
	backupDir := fs.String("backup-dir", "", "Backup directory (default: gameDir/.backups)")

	if err := fs.Parse(args); err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing flags: %v\n", err)
		os.Exit(1)
	}

	if *file == "" {
		fmt.Fprintf(os.Stderr, "Error: --file is required\n")
		fs.Usage()
		os.Exit(1)
	}

	*gameDir = os.ExpandEnv(*gameDir)
	if *backupDir != "" {
		*backupDir = os.ExpandEnv(*backupDir)
	}

	fullPath := filepath.Join(*gameDir, *file)
	if *backupDir == "" {
		*backupDir = filepath.Join(*gameDir, ".backups")
	}

	backupManager := verify.NewBackupManager(*backupDir, true)

	fmt.Printf("Creating backup of: %s\n", *file)

	backupPath, err := backupManager.CreateBackup(fullPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Backup failed: %v\n", err)
		os.Exit(1)
	}

	if backupPath == "" {
		fmt.Println("File does not exist, nothing to backup")
	} else {
		fmt.Printf("Backup created: %s\n", backupPath)
	}
}

func printReport(report *verify.Report) {
	fmt.Println()
	fmt.Println("=" + strings.Repeat("=", 60))
	fmt.Println("VERIFICATION REPORT")
	fmt.Println("=" + strings.Repeat("=", 60))
	fmt.Printf("Version: %s\n", report.Version)
	fmt.Printf("Game Directory: %s\n", report.GameDir)
	fmt.Printf("Timestamp: %s\n", report.Timestamp.Format("2006-01-02 15:04:05"))
	fmt.Printf("Overall Status: %s\n", report.OverallStatus)
	fmt.Println()

	fmt.Println("SUMMARY:")
	fmt.Printf("  Total Files: %d\n", report.Summary.TotalFiles)
	fmt.Printf("  Passed: %d\n", report.Summary.Passed)
	fmt.Printf("  Failed: %d\n", report.Summary.Failed)
	fmt.Printf("  Warnings: %d\n", report.Summary.Warnings)
	fmt.Printf("  Missing Files: %d\n", report.Summary.MissingFiles)
	fmt.Printf("  Modified Files: %d\n", report.Summary.ModifiedFiles)
	fmt.Println()

	// Show failed files
	if report.Summary.Failed > 0 {
		fmt.Println("FAILED FILES:")
		for _, f := range report.Files {
			if f.Status == verify.StatusFailed {
				fmt.Printf("  [FAIL] %s\n", f.Path)
				fmt.Printf("         %s\n", f.Message)
			}
		}
		fmt.Println()
	}

	// Show warnings
	if report.Summary.Warnings > 0 {
		fmt.Println("WARNINGS:")
		for _, f := range report.Files {
			if f.Status == verify.StatusWarning {
				fmt.Printf("  [WARN] %s\n", f.Path)
				fmt.Printf("         %s\n", f.Message)
			}
		}
		fmt.Println()
	}

	if report.OverallStatus == verify.StatusOK {
		fmt.Println("All files verified successfully!")
	} else if report.OverallStatus == verify.StatusWarning {
		fmt.Println("Some files have warnings. Review the list above.")
	} else {
		fmt.Println("Verification FAILED. Some files are missing or corrupted.")
		fmt.Println("Recommendation: Reinstall the game or restore from backup.")
	}
	fmt.Println(strings.Repeat("=", 61))
}

func formatBytes(bytes int64) string {
	const (
		KB = 1024
		MB = 1024 * KB
		GB = 1024 * MB
	)

	switch {
	case bytes >= GB:
		return fmt.Sprintf("%.2f GB", float64(bytes)/GB)
	case bytes >= MB:
		return fmt.Sprintf("%.2f MB", float64(bytes)/MB)
	case bytes >= KB:
		return fmt.Sprintf("%.2f KB", float64(bytes)/KB)
	default:
		return fmt.Sprintf("%d B", bytes)
	}
}

func getDefaultGameDir() string {
	home, _ := os.UserHomeDir()
	switch runtime.GOOS {
	case "windows":
		return filepath.Join(home, "AppData", "Roaming", "Hytale", "install", "release", "package", "game", "latest")
	case "darwin":
		return filepath.Join(home, "Library", "Application Support", "Hytale", "install", "release", "package", "game", "latest")
	default:
		return filepath.Join(home, ".local", "share", "Hytale", "install", "release", "package", "game", "latest")
	}
}
