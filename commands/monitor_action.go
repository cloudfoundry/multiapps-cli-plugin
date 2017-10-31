package commands

import (
	"github.com/SAP/cf-mta-plugin/clients/csrf"
	mtaclient "github.com/SAP/cf-mta-plugin/clients/mtaclient"
	"github.com/SAP/cf-mta-plugin/ui"
	"github.com/cloudfoundry/cli/cf/terminal"
)

// MonitorAction monitors process execution
type MonitorAction struct{}

// Execute executes monitor action on process with the specified id
func (a *MonitorAction) Execute(operationID, commandName string, mtaClient mtaclient.MtaClientOperations, sessionProvider csrf.SessionProvider) ExecutionStatus {

	responseHeader, err := mtaClient.ExecuteAction(operationID, "monitor")
	if err != nil {
		ui.Failed("Could not monitor multi-target app operation with id %s: %s", terminal.EntityNameColor(operationID), err)
		return Failure
	}
	ui.Ok()

	return NewExecutionMonitor(commandName, responseHeader.Location.String(), mtaClient).Monitor()
}
