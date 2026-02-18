package app

import (
	"HyLauncher/internal/config"
	"HyLauncher/internal/platform"
	"HyLauncher/internal/progress"
	"HyLauncher/internal/updater"
	"HyLauncher/pkg/fileutil"
	"HyLauncher/pkg/hyerrors"
	"HyLauncher/pkg/logger"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	wailsRuntime "github.com/wailsapp/wails/v2/pkg/runtime"
)

func (a *App) CheckUpdate() (*updater.Asset, error) {
	asset, latestVersion, err := updater.CheckUpdate(a.ctx, config.LauncherVersion)
	if err != nil {
		// Log the error for debugging but don't report to user (expected when offline)
		logger.Warn("Update check failed", "error", err.Error(), "current", config.LauncherVersion)
		return nil, nil
	}

	if asset == nil {
		logger.Info("No update available", "current", config.LauncherVersion, "latest", latestVersion)
	} else {
		logger.Info("Update available", "current", config.LauncherVersion, "latest", latestVersion, "url", asset.URL)
	}

	return asset, nil
}

func (a *App) Update() error {
	asset, newVersion, err := updater.CheckUpdate(a.ctx, config.LauncherVersion)
	if err != nil {
		appErr := hyerrors.WrapNetwork(err, "failed to check for updates").
			WithContext("current_version", config.LauncherVersion)
		hyerrors.Report(appErr)
		return appErr
	}

	if asset == nil {
		return nil
	}

	reporter := progress.New(a.ctx)

	tmp, err := updater.DownloadTemp(a.ctx, asset.URL, reporter)
	if err != nil {
		appErr := hyerrors.WrapNetwork(err, "failed to download update").
			WithContext("url", asset.URL).
			WithContext("version", newVersion)
		hyerrors.Report(appErr)
		return appErr
	}

	if asset.Sha256 != "" {
		reporter.Report(progress.StageUpdate, 100, "Verifying checksum...")

		if err := fileutil.VerifySHA256(tmp, asset.Sha256); err != nil {
			os.Remove(tmp)
			appErr := hyerrors.WrapFileSystem(err, "update file verification failed").
				WithContext("expected_sha256", asset.Sha256).
				WithContext("file", tmp)
			hyerrors.Report(appErr)
			return appErr
		}
	}

	helperPath, err := updater.EnsureUpdateHelper(a.ctx)
	if err != nil {
		appErr := hyerrors.WrapFileSystem(err, "failed to prepare update helper")
		hyerrors.Report(appErr)
		return appErr
	}

	exe, err := os.Executable()
	if err != nil {
		appErr := hyerrors.WrapFileSystem(err, "failed to get executable path")
		hyerrors.Report(appErr)
		return appErr
	}

	var cmd *exec.Cmd
	if runtime.GOOS == "darwin" {
		// On macOS, find the .app bundle and use open command
		appPath := findAppBundle(exe)
		if appPath != "" {
			// Use update helper to replace then open the app
			cmd = exec.Command(helperPath, exe, tmp, appPath)
		} else {
			cmd = exec.Command(helperPath, exe, tmp)
		}
	} else {
		cmd = exec.Command(helperPath, exe, tmp)
	}
	platform.HideConsoleWindow(cmd)

	cmd.Stdin = nil
	cmd.Stdout = nil
	cmd.Stderr = nil

	if err := cmd.Start(); err != nil {
		appErr := hyerrors.WrapUpdate(err, "failed to start update helper").
			WithContext("helper_path", helperPath).
			WithContext("launcher_path", exe).
			WithContext("update_file", tmp)
		hyerrors.Report(appErr)
		return appErr
	}

	_ = cmd.Process.Release()

	time.Sleep(500 * time.Millisecond)
	os.Exit(0)
	return nil
}

func findAppBundle(exePath string) string {
	// Walk up from the executable to find the .app bundle
	dir := filepath.Dir(exePath)
	for dir != "/" && dir != "." {
		if strings.HasSuffix(dir, ".app") {
			return dir
		}
		dir = filepath.Dir(dir)
	}
	return ""
}

func (a *App) checkUpdateSilently() {
	// Check for updates on all platforms
	asset, latestVersion, err := updater.CheckUpdate(a.ctx, config.LauncherVersion)
	if err != nil {
		logger.Warn("Silent update check failed", "error", err.Error())
		return
	}

	if asset == nil {
		logger.Debug("No update available", "current", config.LauncherVersion, "latest", latestVersion)
		return
	}

	logger.Info("Update available, emitting event", "current", config.LauncherVersion, "latest", latestVersion)
	wailsRuntime.EventsEmit(a.ctx, "update:available", asset)
}
