package app

import (
	"encoding/json"
	"fmt"
)

// updateState updates the state of the application.
func (a *App) updateState() {
	// enum values should be provided in the form of "<ENUM-NAME>_VALUE"
	a.configState.Shell = fmt.Sprintf("SHELL_%s", a.configState.Shell)

	jsData, err := json.Marshal(a.configState)
	if err != nil {
		a.logger.Info().Msgf("failed to marshal json data: %v", err)
		return
	}

	err = a.NDKAgent.UpdateState(AppRoot, string(jsData))
	if err != nil {
		a.logger.Error().Msgf("failed to update state: %v", err)
	}
}
