package commands

import (
	"github.com/SAP/cf-mta-plugin/clients/baseclient"
	"github.com/SAP/cf-mta-plugin/clients/csrf"
	mtaclient "github.com/SAP/cf-mta-plugin/clients/mtaclient"
	"github.com/SAP/cf-mta-plugin/ui"
	"github.com/cloudfoundry/cli/cf/terminal"
)

// ResumeAction retries the process with the specified id
type ResumeAction struct{}

// Execute executes resume action on process with the specified id
func (a *ResumeAction) Execute(operationID, commandName string, mtaClient mtaclient.MtaClientOperations, sessionProvider csrf.SessionProvider) ExecutionStatus {

	// TODO: Ensure session is not expired
	err := sessionProvider.GetSession()
	if err != nil {
		ui.Failed("Could not get x-csrf-token for the current session: %s", err)
		return Failure
	}
	ui.Say("Resuming multi-target app operation with id %s...", terminal.EntityNameColor(operationID))
	responseHeader, err := mtaClient.ExecuteAction(operationID, "resume")
	if err != nil {
		ui.Failed("Could not resume multi-target app operation: %s", baseclient.NewClientError(err))
		return Failure
	}
	ui.Ok()

	operation, err := getMonitoringOperation(operationID, mtaClient)
	if err != nil {
		ui.Failed("Could not monitor multi-target app operation with id %s: %s", terminal.EntityNameColor(operationID), baseclient.NewClientError(err))
		return Failure
	}

	return NewExecutionMonitor(commandName, responseHeader.Location.String(), operation.Messages, mtaClient).Monitor()
}
