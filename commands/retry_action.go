package commands

import (
	"github.com/SAP/cf-mta-plugin/clients/baseclient"
	"github.com/SAP/cf-mta-plugin/clients/csrf"
	mtaclient "github.com/SAP/cf-mta-plugin/clients/mtaclient"
	"github.com/SAP/cf-mta-plugin/ui"
	"github.com/cloudfoundry/cli/cf/terminal"
)

// RetryAction retries the process with the specified id
type RetryAction struct{}

// Execute executes retry action on process with the specified id
func (a *RetryAction) Execute(operationID, commandName string, mtaClient mtaclient.MtaClientOperations, sessionProvider csrf.SessionProvider) ExecutionStatus {

	err := sessionProvider.GetSession()
	if err != nil {
		ui.Failed("Could not retrieve x-csrf-token for the current session: %s", baseclient.NewClientError(err))
		return Failure
	}

	ui.Say("Retrying multi-target app operation with id %s...", terminal.EntityNameColor(operationID))
	responseHeader, err := mtaClient.ExecuteAction(operationID, "retry")
	if err != nil {
		ui.Failed("Could not retry multi-target app operation with id %s: %s", operationID, baseclient.NewClientError(err))
		return Failure
	}
	ui.Ok()

	operation, err := getMonitoringOperation(operationID, mtaClient)
	if err != nil {
		ui.Failed("Could not monitor multi-target app operation with id %s: %s", terminal.EntityNameColor(operationID), baseclient.NewClientError(err))
		return Failure
	}

	monitor := NewExecutionMonitor(commandName, responseHeader.Location.String(), operation.Messages, mtaClient)
	return monitor.Monitor()
}
