package commands

import (
	mtaclient "github.com/SAP/cf-mta-plugin/clients/mtaclient"
	"github.com/SAP/cf-mta-plugin/ui"
	"github.com/cloudfoundry/cli/cf/terminal"
)

// RetryAction retries the process with the specified id
type RetryAction struct{}

// Execute executes retry action on process with the specified id
func (a *RetryAction) Execute(operationID, commandName string, mtaClient mtaclient.MtaClientOperations) ExecutionStatus {

	// TODO: Ensure session is not expired
	// EnsureSlppSession(slppClient)

	ui.Say("Retrying multi-target app operation with id %s...", terminal.EntityNameColor(operationID))
	responseHeader, err = mtaClient.ExecuteAction(operationID, "retry")
	if err != nil {
		ui.Failed("Could not retry multi-target app operation with id %s: %s", terminal.EntityNameColor(operationID), err)
		return Failure
	}
	ui.Ok()

	monitor := NewExecutionMonitor(operationID, commandName, responseHeader.Location, mtaClient)
	return monitor.Monitor()
}
