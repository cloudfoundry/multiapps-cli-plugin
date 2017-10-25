package commands

import (
	mtaclient "github.com/SAP/cf-mta-plugin/clients/mtaclient"
	"github.com/SAP/cf-mta-plugin/ui"
	"github.com/cloudfoundry/cli/cf/terminal"
)

// AbortAction represents action abort
type AbortAction struct{}

// Execute aborts operation with the specified operation id
func (a *AbortAction) Execute(operationID, commandName string, mtaClient mtaclient.MtaClientOperations) ExecutionStatus {

	// TODO: Ensure the session is not expired
	// EnsureSlppSession(slppClient)

	ui.Say("Aborting multi-target app operation with id %s...", terminal.EntityNameColor(operationID))
	_, err := mtaClient.ExecuteAction(operationID, "abort")
	if err != nil {
		ui.Failed("Could not abort multi-target app operation with id %s: %s", terminal.EntityNameColor(operationID), err)
		return Failure
	}
	ui.Ok()
	return Success
}
