package commands

import (
	slppclient "github.com/SAP/cf-mta-plugin/clients/slppclient"
	"github.com/SAP/cf-mta-plugin/ui"
	"github.com/cloudfoundry/cli/cf/terminal"
)

// AbortAction represents action abort
type AbortAction struct{}

// Execute aborts operation with the specified process id
func (a *AbortAction) Execute(processID, commandName string, slppClient slppclient.SlppClientOperations) ExecutionStatus {

	// Ensure the session is not expired
	EnsureSlppSession(slppClient)

	ui.Say("Aborting multi-target app operation with id %s...", terminal.EntityNameColor(processID))
	err := slppClient.ExecuteAction("abort")
	if err != nil {
		ui.Failed("Could not abort multi-target app operation with id %s: %s", terminal.EntityNameColor(processID), err)
		return Failure
	}
	ui.Ok()
	return Success
}
