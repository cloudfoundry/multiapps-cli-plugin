package commands

import (
	"github.com/SAP/cf-mta-plugin/clients/baseclient"
	"github.com/SAP/cf-mta-plugin/clients/csrf"
	mtaclient "github.com/SAP/cf-mta-plugin/clients/mtaclient"
	"github.com/SAP/cf-mta-plugin/ui"
	"github.com/cloudfoundry/cli/cf/terminal"
)

// MonitorAction monitors process execution
type MonitorAction struct{}

// Execute executes monitor action on process with the specified id
func (a *MonitorAction) Execute(operationID, commandName string, mtaClient mtaclient.MtaClientOperations, sessionProvider csrf.SessionProvider) ExecutionStatus {

	err := sessionProvider.GetSession()
	if err != nil {
		ui.Failed("Could not retrieve x-csrf-token for the current session: %s", baseclient.NewClientError(err))
		return Failure
	}

	responseHeader, err := mtaClient.ExecuteAction(operationID, "monitor")
	if err != nil {
		ui.Failed("Could not monitor multi-target app operation with id %s: %s", terminal.EntityNameColor(operationID), err)
		return Failure
	}
	ui.Ok()

	operation, err := getMonitoringOperation(operationID, mtaClient)
	if err != nil {
		ui.Failed("Could not monitor multi-target app operation with id %s: %s", terminal.EntityNameColor(operationID), err)
		return Failure
	}

	return NewExecutionMonitor(commandName, responseHeader.Location.String(), operation.Messages, mtaClient).Monitor()
}
