package app

import (
	"context"
	"encoding/json"
)

const (
	defaultShell = "cli"
	bashShell    = "bash"
)

// ConfigState holds the application configuration and state.
type ConfigState struct {
	// AdminState is the name to use in the greeting.
	AdminState string `json:"admin-state,omitempty"`
	// Shell is the shell to be opened when sshx service is used.
	Shell string `json:"shell,omitempty"`
	URL   string `json:"url,omitempty"`
}

func NewConfigState() *ConfigState {
	return &ConfigState{
		Shell: defaultShell,
	}
}

// loadConfig loads configuration changes for greeter application.
func (a *App) loadConfig() {
	prevAdminState := a.configState.AdminState
	prevShell := a.configState.Shell

	a.configState = NewConfigState() // re-initialize config state
	if a.NDKAgent.Notifications.FullConfig != nil {
		err := json.Unmarshal(a.NDKAgent.Notifications.FullConfig, a.configState)
		if err != nil {
			a.logger.Error().Err(err).Msg("Failed to unmarshal config")
		}

		if a.configState.AdminState == "enable" && prevAdminState != a.configState.AdminState {
			a.logger.Info().
				Str("new admin-state", a.configState.AdminState).
				Str("prev admin-state", prevAdminState).
				Msg("Admin state changed")

			a.restartRequested = true
		}

		if prevShell != a.configState.Shell {
			a.logger.Info().
				Str("new shell", a.configState.Shell).
				Str("prev shell", prevShell).
				Msg("Shell changed")

			a.restartRequested = true
		}
	}
}

// processConfig processes the configuration received from the config notification stream
// and retrieves the uptime from the system.
func (a *App) processConfig(ctx context.Context) {
	a.logger.Info().Msg("Start processing config")

	if a.configState.AdminState == "enable" && a.restartRequested {
		a.startSSHX(ctx)
	}
}
