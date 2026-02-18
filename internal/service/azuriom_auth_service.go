package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"HyLauncher/pkg/model"
)

const (
	// AZURIOM_BASE is the base URL for your Azuriom website
	AZURIOM_BASE = "https://porkland.net"
)

// AzuriomAuthService handles authentication with Azuriom's built-in Auth API
type AzuriomAuthService struct {
	ctx     context.Context
	baseURL string
	client  *http.Client
}

// NewAzuriomAuthService creates a new Azuriom authentication service
func NewAzuriomAuthService(ctx context.Context) *AzuriomAuthService {
	return &AzuriomAuthService{
		ctx:     ctx,
		baseURL: AZURIOM_BASE,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// AzuriomAuthResponse represents the response from Azuriom's authenticate endpoint
type AzuriomAuthResponse struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	UUID     string `json:"uuid"`
	Role     struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	} `json:"role"`
	AccessToken string `json:"access_token"`
	Status      string `json:"status"`
	Reason      string `json:"reason,omitempty"`
	Message     string `json:"message,omitempty"`
}

// Login authenticates the user with email and password using Azuriom's built-in API
// Returns the auth token and user data on success
func (s *AzuriomAuthService) Login(email, password string) (*model.AzuriomAuthData, error) {
	url := s.baseURL + "/api/auth/authenticate"

	reqBody := map[string]string{
		"email":    email,
		"password": password,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal login request: %w", err)
	}

	req, err := http.NewRequestWithContext(s.ctx, http.MethodPost, url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create login request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("login request failed: %w", err)
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read login response: %w", err)
	}

	var authResp AzuriomAuthResponse
	if err := json.Unmarshal(bodyBytes, &authResp); err != nil {
		return nil, fmt.Errorf("failed to decode login response: %w", err)
	}

	// Check for error status
	if authResp.Status == "error" {
		switch authResp.Reason {
		case "invalid_credentials":
			return nil, fmt.Errorf("invalid_credentials")
		case "user_banned":
			return nil, fmt.Errorf("account_blocked")
		case "2fa":
			return nil, fmt.Errorf("2fa_required")
		case "invalid_2fa":
			return nil, fmt.Errorf("invalid_2fa")
		default:
			return nil, fmt.Errorf("login failed: %s", authResp.Message)
		}
	}

	// Check for pending status (2FA required)
	if authResp.Status == "pending" {
		return nil, fmt.Errorf("2fa_required")
	}

	// Extract role name from role object
	roles := []string{}
	if authResp.Role.Name != "" {
		roles = append(roles, authResp.Role.Name)
	}

	return &model.AzuriomAuthData{
		Token:    authResp.AccessToken,
		Username: authResp.Username,
		Roles:    roles,
	}, nil
}

// GetUser fetches the current user data using the access token via verify endpoint
func (s *AzuriomAuthService) GetUser(token string) (*model.AzuriomUser, error) {
	url := s.baseURL + "/api/auth/verify"

	reqBody := map[string]string{
		"access_token": token,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal verify request: %w", err)
	}

	req, err := http.NewRequestWithContext(s.ctx, http.MethodPost, url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create verify request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("verify request failed: %w", err)
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read verify response: %w", err)
	}

	var authResp AzuriomAuthResponse
	if err := json.Unmarshal(bodyBytes, &authResp); err != nil {
		return nil, fmt.Errorf("failed to decode verify response: %w", err)
	}

	// Check for error status
	if authResp.Status == "error" {
		switch authResp.Reason {
		case "invalid_token":
			return nil, fmt.Errorf("session_expired")
		case "user_banned":
			return nil, fmt.Errorf("account_blocked")
		default:
			return nil, fmt.Errorf("verify failed: %s", authResp.Message)
		}
	}

	// Extract role name from role object
	roles := []string{}
	if authResp.Role.Name != "" {
		roles = append(roles, authResp.Role.Name)
	}

	return &model.AzuriomUser{
		Username: authResp.Username,
		Email:    authResp.Email,
		Roles:    roles,
	}, nil
}

// ValidateToken checks if the token is still valid and returns user data
func (s *AzuriomAuthService) ValidateToken(token string) (*model.AzuriomUser, error) {
	return s.GetUser(token)
}

// Logout invalidates the access token on Azuriom server
func (s *AzuriomAuthService) Logout(token string) error {
	url := s.baseURL + "/api/auth/logout"

	reqBody := map[string]string{
		"access_token": token,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("failed to marshal logout request: %w", err)
	}

	req, err := http.NewRequestWithContext(s.ctx, http.MethodPost, url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create logout request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := s.client.Do(req)
	if err != nil {
		return fmt.Errorf("logout request failed: %w", err)
	}
	defer resp.Body.Close()

	// Azuriom returns 200 on successful logout
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("logout failed with status: %s", resp.Status)
	}

	return nil
}

// HasRole checks if the user has a specific role
func HasRole(user *model.AzuriomUser, role string) bool {
	if user == nil {
		return false
	}
	return user.HasRole(role)
}

// HasAnyRole checks if the user has any of the specified roles
func HasAnyRole(user *model.AzuriomUser, roles ...string) bool {
	if user == nil {
		return false
	}
	return user.HasAnyRole(roles...)
}
