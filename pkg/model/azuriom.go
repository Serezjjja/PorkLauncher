package model

// AzuriomUser represents the user data returned from Azuriom custom auth API
type AzuriomUser struct {
	Username string   `json:"username"`
	Roles    []string `json:"roles"`
	Email    string   `json:"email,omitempty"`
}

// AzuriomLoginRequest represents the login request payload
type AzuriomLoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// AzuriomLoginResponse represents the login response from Azuriom API
type AzuriomLoginResponse struct {
	Success  bool     `json:"success"`
	Token    string   `json:"token"`
	Username string   `json:"username"`
	Roles    []string `json:"roles"`
	Message  string   `json:"message,omitempty"`
}

// AzuriomAuthData represents the stored authentication data
type AzuriomAuthData struct {
	Token    string   `json:"token"`
	Username string   `json:"username"`
	Roles    []string `json:"roles"`
}

// HasRole checks if the user has a specific role
func (u *AzuriomUser) HasRole(role string) bool {
	for _, r := range u.Roles {
		if r == role {
			return true
		}
	}
	return false
}

// HasAnyRole checks if the user has any of the specified roles
func (u *AzuriomUser) HasAnyRole(roles ...string) bool {
	for _, role := range roles {
		if u.HasRole(role) {
			return true
		}
	}
	return false
}
