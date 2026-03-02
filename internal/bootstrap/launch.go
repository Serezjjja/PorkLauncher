package bootstrap

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

// LaunchPayload starts the payload executable from the extracted directory.
// On macOS, if an .app bundle is found, it launches via "open" command.
// On other platforms, it executes the binary directly.
// The bootstrap process exits after launching the payload.
func LaunchPayload(payloadDir string) error {
	if runtime.GOOS == "darwin" {
		return launchDarwin(payloadDir)
	}
	return launchDefault(payloadDir)
}

func launchDefault(payloadDir string) error {
	exePath := filepath.Join(payloadDir, PayloadExecutableName())

	info, err := os.Stat(exePath)
	if err != nil {
		return fmt.Errorf("payload executable not found: %w", err)
	}

	// Ensure executable permission on Unix
	if runtime.GOOS != "windows" {
		if info.Mode()&0111 == 0 {
			if err := os.Chmod(exePath, info.Mode()|0755); err != nil {
				return fmt.Errorf("set executable permission: %w", err)
			}
		}
	}

	cmd := exec.Command(exePath)
	cmd.Dir = payloadDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	// Detach the child process so bootstrap can exit
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("start payload: %w", err)
	}

	// Release the child — don't wait for it
	if err := cmd.Process.Release(); err != nil {
		// Non-fatal: payload is already running
		fmt.Fprintf(os.Stderr, "warning: release payload process: %v\n", err)
	}

	return nil
}

func launchDarwin(payloadDir string) error {
	// Look for .app bundle first
	entries, err := os.ReadDir(payloadDir)
	if err != nil {
		return fmt.Errorf("read payload dir: %w", err)
	}

	for _, entry := range entries {
		if entry.IsDir() && filepath.Ext(entry.Name()) == ".app" {
			appPath := filepath.Join(payloadDir, entry.Name())
			cmd := exec.Command("open", appPath)
			if err := cmd.Start(); err != nil {
				return fmt.Errorf("open .app bundle: %w", err)
			}
			if err := cmd.Process.Release(); err != nil {
				fmt.Fprintf(os.Stderr, "warning: release payload process: %v\n", err)
			}
			return nil
		}
	}

	// Fallback to direct binary launch
	return launchDefault(payloadDir)
}
