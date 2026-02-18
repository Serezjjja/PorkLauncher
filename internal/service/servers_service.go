package service

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
	// Return your server - edit these values for your server
	servers := []ServerWithUrls{
		{
			Server: Server{
				ID:          1,
				Name:        "PorkLand",                                        // Change to your server name
				Description: "Супер пупер сервер HyTale",                       // Change description
				Logo:        "https://porkland.net/storage/img/logoresize.png", // URL to logo image or empty
				Banner:      "https://porkland.net/storage/img/bg2.jpg",        // URL to banner image or empty
				IP:          "play.porkland.net:5520",                          // Change to your server IP
			},
			LogoURL:   "https://porkland.net/storage/img/logoresize.png", // Full URL to logo or empty
			BannerURL: "https://porkland.net/storage/img/bg2.jpg",        // Full URL to banner or empty
		},
	}

	return servers, nil
}
