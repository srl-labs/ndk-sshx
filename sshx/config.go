package sshx

import (
	"encoding/json"
)

const (
	SSHXBinPath  = "/opt/sshx/sshx-bin"
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
	a.configState = NewConfigState() // re-initialize config state
	if a.NDKAgent.Notifications.FullConfig != nil {
		err := json.Unmarshal(a.NDKAgent.Notifications.FullConfig, a.configState)
		if err != nil {
			a.logger.Error().Err(err).Msg("Failed to unmarshal config")
		}
	}
}

// processConfig processes the configuration received from the config notification stream
// and retrieves the uptime from the system.
func (a *App) processConfig() {
}
