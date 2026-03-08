package config

// Build-time configuration variables
// These are injected at build time using -ldflags
// Example: -ldflags="-X 'HyLauncher/internal/config.AzuriomBaseURL=https://example.com'"

var (
	// AzuriomBaseURL is the base URL for the Azuriom authentication server
	// Set this via GitHub Secrets or environment variable during build
	AzuriomBaseURL = ""

	// SessionServiceURL is the base URL for the game session service
	SessionServiceURL = ""

	// PatchDomain is the target domain for client patching (e.g., "porkln.fun")
	PatchDomain = ""

	// PatchAPIURL is the URL for the patch API endpoint
	// This is the patches-config endpoint that returns { patches_url: "..." }
	PatchAPIURL = ""

	// GamePatchesURL is the base URL for game patch files
	GamePatchesURL = ""

	// ServerLogoURL is the URL for the server logo image
	ServerLogoURL = ""

	// ServerBannerURL is the URL for the server banner image
	ServerBannerURL = ""

	// ServerIP is the game server IP address
	ServerIP = ""

	// ServerName is the display name of the game server
	ServerName = ""

	// JREManifestURL is the base URL for JRE manifest files
	JREManifestURL = ""

	// ButlerBaseURL is the base URL for Butler downloads
	ButlerBaseURL = ""

	// DiscordAppID is the Discord Rich Presence application ID
	DiscordAppID = ""

	// AuthDomain is the auth server domain (e.g., "auth.sanasol.ws")
	AuthDomain = ""
)

// Hytale-F2P API configuration
// Multi-source fallback chain for patches config
const (
	// Primary patches config source
	DefaultPatchesConfigPrimary = "https://auth.sanasol.ws/api/patches-config"
	// Backup patches config sources
	DefaultPatchesConfigBackup1 = "https://htdwnldsan.top/patches-config"
	DefaultPatchesConfigBackup2 = "https://dl1.htdwnldsan.top/patches-config"
	// Hardcoded fallback patches URL
	DefaultPatchesFallback = "https://dl.vboro.de/patches"
	// DNS TXT record for patches URL
	DefaultPatchesDNSTXT = "_patches.htdwnldsan.top"
	// Default auth domain
	DefaultAuthDomain = "auth.sanasol.ws"
)

// GetAzuriomBaseURL returns the Azuriom base URL
// Falls back to empty string if not set (will cause errors if used without being set)
func GetAzuriomBaseURL() string {
	return AzuriomBaseURL
}

// GetSessionServiceURL returns the session service URL
func GetSessionServiceURL() string {
	if SessionServiceURL == "" {
		return "https://localhost"
	}
	return SessionServiceURL
}

// GetPatchDomain returns the patch domain
func GetPatchDomain() string {
	if PatchDomain == "" {
		return "localhost"
	}
	return PatchDomain
}

// GetPatchAPIURL returns the patch API URL (patches-config endpoint)
// Returns the primary patches config URL by default
func GetPatchAPIURL() string {
	if PatchAPIURL != "" {
		return PatchAPIURL
	}
	return DefaultPatchesConfigPrimary
}

// GetPatchesConfigSources returns all patches config sources for fallback
func GetPatchesConfigSources() []string {
	if PatchAPIURL != "" {
		return []string{PatchAPIURL}
	}
	return []string{
		DefaultPatchesConfigPrimary,
		DefaultPatchesConfigBackup1,
		DefaultPatchesConfigBackup2,
	}
}

// GetPatchesFallbackURL returns the hardcoded fallback patches URL
func GetPatchesFallbackURL() string {
	return DefaultPatchesFallback
}

// GetAuthDomain returns the auth server domain
func GetAuthDomain() string {
	if AuthDomain != "" {
		return AuthDomain
	}
	return DefaultAuthDomain
}

// GetGamePatchesURL returns the game patches base URL
func GetGamePatchesURL() string {
	if GamePatchesURL == "" {
		return "https://localhost/patches"
	}
	return GamePatchesURL
}

// GetServerLogoURL returns the server logo URL
func GetServerLogoURL() string {
	return ServerLogoURL
}

// GetServerBannerURL returns the server banner URL
func GetServerBannerURL() string {
	return ServerBannerURL
}

// GetServerIP returns the game server IP
func GetServerIP() string {
	if ServerIP == "" {
		return "localhost:25565"
	}
	return ServerIP
}

// GetServerName returns the server display name
func GetServerName() string {
	if ServerName == "" {
		return "Game Server"
	}
	return ServerName
}

// GetJREManifestURL returns the JRE manifest base URL
func GetJREManifestURL() string {
	if JREManifestURL == "" {
		return "https://localhost/jre"
	}
	return JREManifestURL
}

// GetButlerBaseURL returns the Butler base URL
func GetButlerBaseURL() string {
	if ButlerBaseURL == "" {
		return "https://localhost/butler"
	}
	return ButlerBaseURL
}

// GetDiscordAppID returns the Discord App ID
func GetDiscordAppID() string {
	return DiscordAppID
}
