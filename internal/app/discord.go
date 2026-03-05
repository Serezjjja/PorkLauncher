package app

import (
	"fmt"
	"time"

	"HyLauncher/internal/config"

	"github.com/hugolgst/rich-go/client"
)

func (a *App) discordRPC() {
	if !a.launcherCfg.DiscordRPC {
		return
	}

	now := time.Now()

	err := client.SetActivity(client.Activity{
		State:   "Игра HyTale",
		Details: "Залетай к нам на сервер!",
		Timestamps: &client.Timestamps{
			Start: &now,
		},
		Buttons: []*client.Button{
			{
				Label: "Дискорд",
				Url:   "https://discord.gg/RbreKRwsH7",
			},
			{
				Label: "Телеграм",
				Url:   "https://t.me/porkland",
			},
		},
	})

	_ = err
}

func (a *App) GetDiscordRPC() bool {
	return a.launcherCfg.DiscordRPC
}

func (a *App) SetDiscordRPC(enabled bool) error {
	// Update config
	if err := config.UpdateLauncher(func(cfg *config.LauncherConfig) error {
		cfg.DiscordRPC = enabled
		return nil
	}); err != nil {
		return fmt.Errorf("failed to save Discord RPC setting: %w", err)
	}

	a.launcherCfg.DiscordRPC = enabled

	if enabled {
		discordAppID := config.GetDiscordAppID()
		if discordAppID == "" {
			discordAppID = "1345687653965631540" // fallback for development
		}
		if err := client.Login(discordAppID); err != nil {
			return fmt.Errorf("failed to login to Discord: %w", err)
		}
		go a.discordRPC()
	} else {
		client.Logout()
	}

	return nil
}
