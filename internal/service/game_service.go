package service

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"

	"HyLauncher/internal/env"
	"HyLauncher/internal/game"
	"HyLauncher/internal/java"
	"HyLauncher/internal/patch"
	"HyLauncher/internal/platform"
	"HyLauncher/internal/progress"
	"HyLauncher/internal/verify"
	"HyLauncher/pkg/fileutil"
	"HyLauncher/pkg/logger"
	"HyLauncher/pkg/model"
)

type GameService struct {
	ctx        context.Context
	reporter   *progress.Reporter
	authSvc    *AuthService
	authDomain string
	installMu  sync.Mutex
}

func NewGameService(ctx context.Context, reporter *progress.Reporter, svc *AuthService) *GameService {
	return &GameService{
		ctx:        ctx,
		reporter:   reporter,
		authSvc:    svc,
		authDomain: "sanasol.ws",
	}
}

func (s *GameService) EnsureGame(request model.InstanceModel) error {
	s.reporter.Report(progress.StageVerify, 0, "Verifying installation...")

	if err := java.EnsureJRE(s.ctx, request.Branch, s.reporter); err != nil {
		return fmt.Errorf("jre: %w", err)
	}

	if err := patch.EnsureButler(s.ctx, s.reporter); err != nil {
		return fmt.Errorf("butler: %w", err)
	}

	if err := game.CheckInstalled(s.ctx, request.Branch, request.BuildVersion); err != nil {
		return fmt.Errorf("game files: %w", err)
	}

	// Verify game files integrity before launching
	if err := s.verifyGameFiles(request); err != nil {
		logger.Warn("Game file verification failed", "error", err)
		// Don't fail launch on verification error, just log it
	}

	s.reporter.Report(progress.StageVerify, 100, "Ready")
	return nil
}

// verifyGameFiles performs integrity verification on game files
func (s *GameService) verifyGameFiles(request model.InstanceModel) error {
	gameDir := env.GetGameDir(request.Branch, request.BuildVersion)

	// Check if verification should be skipped
	if os.Getenv("HYTALE_SKIP_VERIFY") == "1" {
		logger.Info("Skipping game file verification (HYTALE_SKIP_VERIFY=1)")
		return nil
	}

	options := verify.Options{
		GameDir:       gameDir,
		Version:       request.BuildVersion,
		CreateBackups: true,
		ProgressCallback: func(current, total int64, fileName string) {
			// Update progress every 10%
			percent := float64(current) / float64(total) * 100
			if int(percent)%10 == 0 {
				s.reporter.Report(progress.StageVerify, percent, fmt.Sprintf("Verifying %s...", fileName))
			}
		},
	}

	report, err := verify.VerifyWithOptions(options)
	if err != nil {
		return fmt.Errorf("verification error: %w", err)
	}

	// Log verification results
	logger.Info("Game file verification complete",
		"version", report.Version,
		"total", report.Summary.TotalFiles,
		"passed", report.Summary.Passed,
		"failed", report.Summary.Failed,
		"warnings", report.Summary.Warnings,
		"status", report.OverallStatus)

	// Report any issues
	if report.OverallStatus == verify.StatusFailed {
		for _, file := range report.Files {
			if file.Status == verify.StatusFailed {
				logger.Error("File verification failed",
					"file", file.Path,
					"message", file.Message)
			}
		}
		return fmt.Errorf("file verification failed: %d files failed", report.Summary.Failed)
	}

	if report.OverallStatus == verify.StatusWarning {
		for _, file := range report.Files {
			if file.Status == verify.StatusWarning {
				logger.Warn("File verification warning",
					"file", file.Path,
					"message", file.Message)
			}
		}
	}

	return nil
}

func (s *GameService) EnsureInstalled(ctx context.Context, request model.InstanceModel, reporter *progress.Reporter) (string, error) {
	s.installMu.Lock()
	defer s.installMu.Unlock()

	if reporter != nil {
		reporter.Report(progress.StageVerify, 0, "Checking for updates...")
	}

	latest, err := patch.FindLatestVersion(request.Branch)
	if err != nil {
		return "", fmt.Errorf("fetch latest: %w", err)
	}

	switch request.BuildVersion {
	case "auto":
		return s.handleAutoVersion(ctx, request.Branch, latest, reporter)
	case "latest":
		return s.handleLatestVersion(ctx, request.Branch, latest, reporter)
	default:
		if err := s.EnsureGame(request); err == nil {
			return request.BuildVersion, nil
		}
		return "", fmt.Errorf("version %q not installed", request.BuildVersion)
	}
}

func (s *GameService) handleAutoVersion(ctx context.Context, branch string, latest int, reporter *progress.Reporter) (string, error) {
	autoDir := env.GetGameDir(branch, "auto")
	versionFile := filepath.Join(autoDir, ".version")

	currentVer := s.readVersionFile(versionFile)

	if currentVer == latest && game.CheckInstalled(ctx, branch, "auto") == nil {
		if reporter != nil {
			reporter.Report(progress.StageVerify, 100, "Up to date")
		}
		return "auto", nil
	}

	if reporter != nil {
		reporter.Report(progress.StageVerify, 50, fmt.Sprintf("Updating to %d...", latest))
	}

	if err := s.install(ctx, branch, "auto", latest, reporter); err != nil {
		return "", err
	}

	_ = os.WriteFile(versionFile, []byte(strconv.Itoa(latest)), 0644)
	return "auto", nil
}

func (s *GameService) handleLatestVersion(ctx context.Context, branch string, latest int, reporter *progress.Reporter) (string, error) {
	versionStr := strconv.Itoa(latest)

	if game.CheckInstalled(ctx, branch, versionStr) == nil {
		if reporter != nil {
			reporter.Report(progress.StageVerify, 100, "Up to date")
		}
		return versionStr, nil
	}

	if reporter != nil {
		reporter.Report(progress.StageVerify, 50, fmt.Sprintf("Installing %d...", latest))
	}

	if err := s.install(ctx, branch, versionStr, latest, reporter); err != nil {
		return "", err
	}

	return versionStr, nil
}

func (s *GameService) readVersionFile(path string) int {
	data, err := os.ReadFile(path)
	if err != nil {
		logger.Warn("Failed to read version file", "path", path, "error", err)
		return 0
	}
	content := strings.TrimSpace(string(data))
	logger.Info("Read version file", "path", path, "content", content)
	ver, err := strconv.Atoi(content)
	if err != nil {
		logger.Warn("Failed to parse version", "path", path, "content", content, "error", err)
		return 0
	}
	return ver
}

func (s *GameService) install(ctx context.Context, branch, version string, targetVer int, reporter *progress.Reporter) error {
	if err := java.EnsureJRE(ctx, branch, reporter); err != nil {
		return fmt.Errorf("jre: %w", err)
	}

	if err := patch.EnsureButler(ctx, reporter); err != nil {
		return fmt.Errorf("butler: %w", err)
	}

	currentVer := 0
	if version == "auto" {
		versionFile := filepath.Join(env.GetGameDir(branch, "auto"), ".version")
		currentVer = s.readVersionFile(versionFile)
		logger.Info("Read version file for install", "file", versionFile, "currentVer", currentVer, "targetVer", targetVer)
	}

	logger.Info("Starting patch download", "branch", branch, "currentVer", currentVer, "targetVer", targetVer, "versionDir", version)
	if err := patch.DownloadAndApplyPWR(ctx, branch, currentVer, targetVer, version, reporter); err != nil {
		logger.Error("Patch failed, attempting full reinstall", "error", err, "branch", branch, "version", version)

		if reporter != nil {
			reporter.Report(progress.StagePatch, 0, "Patch failed, cleaning for reinstall...")
		}

		// Clean the game directory for fresh install
		if cleanErr := s.cleanGameDirectory(branch, version); cleanErr != nil {
			logger.Error("Failed to clean game directory", "error", cleanErr)
			return fmt.Errorf("patch failed and cleanup failed: %v (original error: %w)", cleanErr, err)
		}

		// Retry with fresh install (currentVer = 0 forces full download)
		logger.Info("Retrying with fresh install", "branch", branch, "targetVer", targetVer, "versionDir", version)
		if reporter != nil {
			reporter.Report(progress.StagePatch, 0, "Downloading full game...")
		}

		if retryErr := patch.DownloadAndApplyPWR(ctx, branch, 0, targetVer, version, reporter); retryErr != nil {
			logger.Error("Full reinstall also failed", "error", retryErr)
			return fmt.Errorf("patch failed and full reinstall failed: %w (original error: %v)", retryErr, err)
		}

		logger.Info("Full reinstall successful")
	}

	if err := s.fixPermissions(branch, version); err != nil {
		return err
	}

	if err := s.applyAuthPatch(branch, version, reporter); err != nil {
		logger.Warn("Auth patch failed", "error", err)
	}

	if reporter != nil {
		reporter.Report(progress.StageComplete, 100, "Done")
	}
	return nil
}

// cleanGameDirectory removes all files in the game directory for fresh install
func (s *GameService) cleanGameDirectory(branch, version string) error {
	gameDir := env.GetGameDir(branch, version)
	logger.Info("Cleaning game directory", "dir", gameDir)

	// Check if directory exists
	if _, err := os.Stat(gameDir); os.IsNotExist(err) {
		// Directory doesn't exist, nothing to clean
		return nil
	}

	// Read all entries in the directory
	entries, err := os.ReadDir(gameDir)
	if err != nil {
		return fmt.Errorf("failed to read game directory: %w", err)
	}

	// Remove all entries
	for _, entry := range entries {
		path := filepath.Join(gameDir, entry.Name())
		if err := os.RemoveAll(path); err != nil {
			logger.Warn("Failed to remove file during cleanup", "path", path, "error", err)
			// Continue trying to remove other files
		}
	}

	// Also clean the version file for auto version
	if version == "auto" {
		versionFile := filepath.Join(gameDir, ".version")
		_ = os.Remove(versionFile)
	}

	logger.Info("Game directory cleaned successfully")
	return nil
}

func (s *GameService) fixPermissions(branch, version string) error {
	if runtime.GOOS != "darwin" {
		return nil
	}

	gameDir := env.GetGameDir(branch, version)
	clientExec := filepath.Join(gameDir, "Client", "Hytale.app", "Contents", "MacOS", "HytaleClient")

	if fileutil.FileExists(clientExec) {
		_ = os.Chmod(clientExec, 0755)
	}

	javaExec := filepath.Join(env.GetJREDir(), "bin", "java")
	if fileutil.FileExists(javaExec) {
		_ = os.Chmod(javaExec, 0755)
	}

	return nil
}

func (s *GameService) applyAuthPatch(branch, version string, reporter *progress.Reporter) error {
	if reporter != nil {
		reporter.Report(progress.StagePatch, 0, "Patching auth...")
	}

	req := model.InstanceModel{BuildVersion: version, Branch: branch}
	if err := patch.EnsureGamePatched(s.ctx, req, s.authDomain, reporter); err != nil {
		return err
	}

	if reporter != nil {
		reporter.Report(progress.StagePatch, 100, "Auth ready")
	}
	return nil
}

func (s *GameService) Launch(playerName string, request model.InstanceModel, serverIP ...string) error {
	session, err := s.authSvc.FetchGameSession(playerName)
	if err != nil {
		return err
	}

	s.reporter.Report(progress.StageLaunch, 0, "Launching...")

	gameDir := env.GetGameDir(request.Branch, request.BuildVersion)
	userDataDir := env.GetInstanceUserDataDir(request.InstanceID)

	if err := os.MkdirAll(userDataDir, 0755); err != nil {
		return fmt.Errorf("userdata: %w", err)
	}

	// Create ServerList.json with default server
	if err := s.createServerList(userDataDir); err != nil {
		// Log error but don't fail launch
		fmt.Printf("Failed to create ServerList.json: %v\n", err)
	}

	_ = patch.EnsureGamePatched(s.ctx, request, s.authDomain, nil)

	clientPath := env.GetGameClientPath(request.Branch, request.BuildVersion)
	if clientPath == "" {
		return fmt.Errorf("client not found")
	}

	javaBin, err := java.GetJavaExec(request.Branch)
	if err != nil {
		return fmt.Errorf("java: %w", err)
	}

	if runtime.GOOS == "darwin" {
		_ = os.Chmod(clientPath, 0755)
		_ = os.Chmod(javaBin, 0755)
	}

	args := []string{
		"--app-dir", mustAbs(gameDir),
		"--user-dir", mustAbs(userDataDir),
		"--java-exec", mustAbs(javaBin),
		"--auth-mode", "authenticated",
		"--uuid", session.UUID,
		"--name", session.Username,
		"--identity-token", session.IdentityToken,
		"--session-token", session.SessionToken,
	}

	if len(serverIP) > 0 && serverIP[0] != "" {
		args = append(args, "--server", serverIP[0])
	}

	cmd := exec.Command(clientPath, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	game.SetSDLVideoDriver(cmd)

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("start: %w", err)
	}

	if runtime.GOOS == "darwin" {
		_ = cmd.Process.Release()
		_ = platform.RemoveQuarantine(clientPath)
	}

	time.Sleep(500 * time.Millisecond)

	if cmd.Process != nil && runtime.GOOS != "windows" {
		if err := cmd.Process.Signal(syscall.Signal(0)); err != nil {
			_ = cmd.Wait()
			return fmt.Errorf("process exited")
		}
	}

	s.reporter.Report(progress.StageLaunch, 100, "Launched!")
	s.reporter.Reset()
	return nil
}

func mustAbs(path string) string {
	abs, _ := filepath.Abs(path)
	return abs
}

// ServerListEntry represents a saved server entry
type ServerListEntry struct {
	ID        string    `json:"Id"`
	Name      string    `json:"Name"`
	Address   string    `json:"Address"`
	DateSaved time.Time `json:"DateSaved"`
}

// ServerList represents the ServerList.json structure
type ServerList struct {
	SavedServers []ServerListEntry `json:"SavedServers"`
}

// createServerList creates the ServerList.json file in the UserData directory
func (s *GameService) createServerList(userDataDir string) error {
	now := time.Now()
	serverList := ServerList{
		SavedServers: []ServerListEntry{
			{
				ID:        "eb3933f0-1b63-4f8c-b98d-e7c4e7da4a8e",
				Name:      "P O R K L A N D",
				Address:   "play.porkland.net:5520",
				DateSaved: now,
			},
			{
				ID:        "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
				Name:      "P O R K L A N D",
				Address:   "play.porkland.net:5520",
				DateSaved: now,
			},
			{
				ID:        "b2c3d4e5-f6a7-8901-bcde-f12345678901",
				Name:      "P O R K L A N D",
				Address:   "play.porkland.net:5520",
				DateSaved: now,
			},
			{
				ID:        "c3d4e5f6-a7b8-9012-cdef-123456789012",
				Name:      "P O R K L A N D",
				Address:   "play.porkland.net:5520",
				DateSaved: now,
			},
			{
				ID:        "d4e5f6a7-b8c9-0123-defa-234567890123",
				Name:      "P O R K L A N D",
				Address:   "play.porkland.net:5520",
				DateSaved: now,
			},
			{
				ID:        "e5f6a7b8-c9d0-1234-efab-345678901234",
				Name:      "P O R K L A N D",
				Address:   "play.porkland.net:5520",
				DateSaved: now,
			},
			{
				ID:        "f6a7b8c9-d0e1-2345-fabc-456789012345",
				Name:      "P O R K L A N D",
				Address:   "play.porkland.net:5520",
				DateSaved: now,
			},
			{
				ID:        "a7b8c9d0-e1f2-3456-abcd-567890123456",
				Name:      "P O R K L A N D",
				Address:   "play.porkland.net:5520",
				DateSaved: now,
			},
			{
				ID:        "b8c9d0e1-f2a3-4567-bcde-678901234567",
				Name:      "P O R K L A N D",
				Address:   "play.porkland.net:5520",
				DateSaved: now,
			},
			{
				ID:        "c9d0e1f2-a3b4-5678-cdef-789012345678",
				Name:      "P O R K L A N D",
				Address:   "play.porkland.net:5520",
				DateSaved: now,
			},
		},
	}

	// Ensure UserData directory exists
	if err := os.MkdirAll(userDataDir, 0755); err != nil {
		return fmt.Errorf("failed to create UserData dir: %w", err)
	}

	// Marshal to JSON
	jsonData, err := json.MarshalIndent(serverList, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal ServerList: %w", err)
	}

	// Write to file
	serverListPath := filepath.Join(userDataDir, "ServerList.json")
	if err := os.WriteFile(serverListPath, jsonData, 0644); err != nil {
		return fmt.Errorf("failed to write ServerList.json: %w", err)
	}

	return nil
}
