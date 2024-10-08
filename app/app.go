package app

import (
	"context"

	"github.com/rs/zerolog"
	"github.com/srl-labs/bond"
)

const (
	AppName = "sshx"
	AppRoot = "/" + AppName
)

// App defines the application.
type App struct {
	Name string
	// configState holds the application configuration and state.
	configState *ConfigState
	logger      *zerolog.Logger
	NDKAgent    *bond.Agent

	restartRequested bool // flag to indicate that an sshx restart is requested
	sshxPid          int  // pid of a running sshx process
}

// New creates a new App instance and connects to the NDK socket.
func New(logger *zerolog.Logger, agent *bond.Agent) *App {
	return &App{
		Name:     AppName,
		NDKAgent: agent,

		configState: NewConfigState(),
		logger:      logger,
	}
}

// Start starts the application.
func (a *App) Start(ctx context.Context) {
	for {
		select {
		case <-a.NDKAgent.Notifications.FullConfigReceived:
			a.logger.Info().Msg("Received full config")

			a.loadConfig()

			a.processConfig(ctx)

			a.updateState()

		case <-ctx.Done():
			return
		}
	}
}
