package app

import (
	"HyLauncher/internal/patch"
	"HyLauncher/pkg/hyerrors"
)

type VersionsResponse struct {
	Versions []int  `json:"versions"`
	Error    string `json:"error,omitempty"`
}

type LaunchResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error,omitempty"`
}

func (a *App) DownloadAndLaunch(playerName string) LaunchResponse {
	return a.downloadAndLaunchInternal(playerName, "")
}

func (a *App) DownloadAndLaunchWithServer(playerName string, serverIP string) LaunchResponse {
	return a.downloadAndLaunchInternal(playerName, serverIP)
}

func (a *App) downloadAndLaunchInternal(playerName string, serverIP string) LaunchResponse {
	// Check authentication and player role before launching
	if err := a.ValidatePlayerAccess(); err != nil {
		hyerrors.Report(hyerrors.Validation("player access check failed").WithContext("error", err.Error()))
		return LaunchResponse{Success: false, Error: err.Error()}
	}

	// Get the authenticated username from Azuriom instead of using the provided playerName
	userResp := a.GetCurrentUser()
	if !userResp.LoggedIn {
		return LaunchResponse{Success: false, Error: "not authenticated"}
	}

	// Use the username from Azuriom as the player name
	authPlayerName := userResp.Username

	if err := a.validatePlayerName(authPlayerName); err != nil {
		hyerrors.Report(hyerrors.Validation("provided invalid username"))
		return LaunchResponse{Success: false, Error: err.Error()}
	}

	_ = a.SyncInstanceState()

	installedVersion, err := a.gameSvc.EnsureInstalled(a.ctx, a.instance, a.progress)
	if err != nil {
		appErr := hyerrors.WrapGame(err, "failed to install game").
			WithContext("branch", a.instance.Branch).
			WithContext("requestedVersion", a.instance.BuildVersion)
		hyerrors.Report(appErr)
		return LaunchResponse{Success: false, Error: appErr.Error()}
	}

	if installedVersion != a.instance.BuildVersion {
		a.instance.BuildVersion = installedVersion
		if err := a.UpdateInstanceVersion(installedVersion); err != nil {
			_ = err
		}
	}

	if err := a.gameSvc.Launch(authPlayerName, a.instance, serverIP); err != nil {
		appErr := hyerrors.GameCritical("failed to launch game").
			WithDetails(err.Error()).
			WithContext("player", authPlayerName).
			WithContext("branch", a.instance.Branch).
			WithContext("version", a.instance.BuildVersion)
		hyerrors.Report(appErr)
		return LaunchResponse{Success: false, Error: appErr.Error()}
	}

	return LaunchResponse{Success: true}
}

func (a *App) GetReleaseVersions() VersionsResponse {
	release, err := patch.ListAllVersions("release")
	if err != nil {
		return VersionsResponse{Error: err.Error()}
	}
	return VersionsResponse{Versions: release}
}

func (a *App) GetPreReleaseVersions() VersionsResponse {
	release, err := patch.ListAllVersions("pre-release")
	if err != nil {
		return VersionsResponse{Error: err.Error()}
	}
	return VersionsResponse{Versions: release}
}
