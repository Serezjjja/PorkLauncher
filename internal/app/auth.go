package app

import (
	"HyLauncher/internal/config"
	"HyLauncher/internal/service"
	"HyLauncher/pkg/hyerrors"
	"HyLauncher/pkg/logger"
	"HyLauncher/pkg/model"
)

// AuthResponse represents the response for authentication operations
type AuthResponse struct {
	Success  bool     `json:"success"`
	Username string   `json:"username,omitempty"`
	Roles    []string `json:"roles,omitempty"`
	Error    string   `json:"error,omitempty"`
}

// LoginRequest represents the login request from frontend
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// CurrentUserResponse represents the current user data response
type CurrentUserResponse struct {
	LoggedIn bool     `json:"loggedIn"`
	Username string   `json:"username,omitempty"`
	Roles    []string `json:"roles,omitempty"`
}

// Login authenticates the user with Azuriom and stores the token
func (a *App) Login(req LoginRequest) AuthResponse {
	azAuthSvc := service.NewAzuriomAuthService(a.ctx)

	authData, err := azAuthSvc.Login(req.Email, req.Password)
	if err != nil {
		errorMsg := err.Error()
		// Map error codes to user-friendly messages
		switch errorMsg {
		case "invalid_credentials":
			errorMsg = "Invalid email or password"
		case "account_blocked":
			errorMsg = "Account is blocked"
		case "2fa_required":
			errorMsg = "Two-factor authentication required"
		case "invalid_2fa":
			errorMsg = "Invalid 2FA code"
		case "session_expired":
			errorMsg = "Session expired"
		default:
			if len(errorMsg) > 100 {
				errorMsg = "Login failed"
			}
		}

		logger.Warn("Login failed", "error", errorMsg, "email", req.Email)
		return AuthResponse{
			Success: false,
			Error:   errorMsg,
		}
	}

	// Save the auth token to config
	if err := config.UpdateLauncher(func(cfg *config.LauncherConfig) error {
		cfg.AzuriomAuthToken = authData.Token
		// Also save the username as the nick
		cfg.Nick = authData.Username
		return nil
	}); err != nil {
		appErr := hyerrors.WrapConfig(err, "failed to save auth token")
		hyerrors.Report(appErr)
		return AuthResponse{
			Success: false,
			Error:   "Failed to save authentication data",
		}
	}

	// Update in-memory config
	a.launcherCfg.AzuriomAuthToken = authData.Token
	a.launcherCfg.Nick = authData.Username

	logger.Info("User logged in successfully", "username", authData.Username)

	return AuthResponse{
		Success:  true,
		Username: authData.Username,
		Roles:    authData.Roles,
	}
}

// Logout clears the authentication token both locally and on Azuriom
func (a *App) Logout() AuthResponse {
	token := a.launcherCfg.AzuriomAuthToken

	// Try to invalidate token on Azuriom server (best effort)
	if token != "" {
		azAuthSvc := service.NewAzuriomAuthService(a.ctx)
		if err := azAuthSvc.Logout(token); err != nil {
			logger.Warn("Failed to logout on server", "error", err.Error())
			// Continue with local logout even if server logout fails
		}
	}

	if err := config.UpdateLauncher(func(cfg *config.LauncherConfig) error {
		cfg.AzuriomAuthToken = ""
		cfg.Nick = ""
		return nil
	}); err != nil {
		appErr := hyerrors.WrapConfig(err, "failed to clear auth token")
		hyerrors.Report(appErr)
		return AuthResponse{
			Success: false,
			Error:   "Failed to logout",
		}
	}

	a.launcherCfg.AzuriomAuthToken = ""
	a.launcherCfg.Nick = ""

	logger.Info("User logged out")

	return AuthResponse{
		Success: true,
	}
}

// GetCurrentUser returns the current authenticated user data
func (a *App) GetCurrentUser() CurrentUserResponse {
	token := a.launcherCfg.AzuriomAuthToken
	if token == "" {
		return CurrentUserResponse{
			LoggedIn: false,
		}
	}

	azAuthSvc := service.NewAzuriomAuthService(a.ctx)
	user, err := azAuthSvc.GetUser(token)
	if err != nil {
		errorMsg := err.Error()
		// If session expired, clear the token
		if errorMsg == "session_expired" {
			_ = config.UpdateLauncher(func(cfg *config.LauncherConfig) error {
				cfg.AzuriomAuthToken = ""
				return nil
			})
			a.launcherCfg.AzuriomAuthToken = ""
		}

		logger.Warn("Failed to get current user", "error", errorMsg)
		return CurrentUserResponse{
			LoggedIn: false,
		}
	}

	return CurrentUserResponse{
		LoggedIn: true,
		Username: user.Username,
		Roles:    user.Roles,
	}
}

// CheckAuthOnStartup validates the stored token on app startup
func (a *App) CheckAuthOnStartup() CurrentUserResponse {
	return a.GetCurrentUser()
}

// HasPlayerRole checks if the current user has the "player" or "Участник" role
func (a *App) HasPlayerRole() bool {
	userResp := a.GetCurrentUser()
	if !userResp.LoggedIn {
		return false
	}

	allowedRoles := []string{"player", "Участник", "Member", "member"}
	for _, role := range userResp.Roles {
		for _, allowed := range allowedRoles {
			if role == allowed {
				return true
			}
		}
	}

	return false
}

// GetAuthToken returns the stored auth token (for internal use)
func (a *App) GetAuthToken() string {
	return a.launcherCfg.AzuriomAuthToken
}

// ValidatePlayerAccess checks if the user can launch the game
// Returns nil if access is granted, error otherwise
func (a *App) ValidatePlayerAccess() error {
	token := a.launcherCfg.AzuriomAuthToken
	if token == "" {
		return hyerrors.Validation("not authenticated").
			WithContext("reason", "no_auth_token")
	}

	azAuthSvc := service.NewAzuriomAuthService(a.ctx)
	user, err := azAuthSvc.GetUser(token)
	if err != nil {
		errorMsg := err.Error()
		if errorMsg == "session_expired" {
			// Clear the invalid token
			_ = config.UpdateLauncher(func(cfg *config.LauncherConfig) error {
				cfg.AzuriomAuthToken = ""
				return nil
			})
			a.launcherCfg.AzuriomAuthToken = ""

			return hyerrors.Validation("session expired, please login again").
				WithContext("reason", "session_expired")
		}

		return hyerrors.Wrap(err, hyerrors.CategoryInternal, "failed to validate user").
			WithContext("reason", "validation_failed")
	}

	// Check if user has an allowed role (player, Участник, Member)
	allowedRoles := []string{"player", "Участник", "Member", "member"}
	hasAllowedRole := false
	for _, role := range user.Roles {
		for _, allowed := range allowedRoles {
			if role == allowed {
				hasAllowedRole = true
				break
			}
		}
		if hasAllowedRole {
			break
		}
	}

	if !hasAllowedRole {
		return hyerrors.Validation("no player access").
			WithContext("reason", "missing_player_role").
			WithContext("username", user.Username).
			WithContext("roles", user.Roles)
	}

	return nil
}

// GetCustomUser fetches the user data from Azuriom API using the stored token
// This is used by the game service before launching
func (a *App) GetCustomUser() (*model.AzuriomUser, error) {
	token := a.launcherCfg.AzuriomAuthToken
	if token == "" {
		return nil, hyerrors.Validation("not authenticated")
	}

	azAuthSvc := service.NewAzuriomAuthService(a.ctx)
	return azAuthSvc.GetUser(token)
}
