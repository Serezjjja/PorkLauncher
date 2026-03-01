package service

import "HyLauncher/internal/config"

// Server represents a game server
type Server struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Logo        string `json:"logo"`
	Banner      string `json:"banner"`
	IP          string `json:"ip"`
}

// ServerWithUrls represents a server with full URLs
type ServerWithUrls struct {
	Server
	LogoURL   string `json:"logo_url"`
	BannerURL string `json:"banner_url"`
}

// ServersService handles servers - disabled (returns empty data)
type ServersService struct{}

// NewServersService creates a new servers service
func NewServersService() *ServersService {
	return &ServersService{}
}

// FetchServers returns your game server
func (s *ServersService) FetchServers() ([]ServerWithUrls, error) {
	// Return your server - configured via build-time variables
	servers := []ServerWithUrls{
		{
			Server: Server{
				ID:          1,
				Name:        config.GetServerName(),
				Description: "",
				Logo:        config.GetServerLogoURL(),
				Banner:      config.GetServerBannerURL(),
				IP:          config.GetServerIP(),
			},
			LogoURL:   config.GetServerLogoURL(),
			BannerURL: config.GetServerBannerURL(),
		},
	}

	return servers, nil
}
