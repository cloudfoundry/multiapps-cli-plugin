package commands

import (
	mtaclient "github.com/SAP/cf-mta-plugin/clients/mtaclient"
	"github.com/SAP/cf-mta-plugin/ui"
	"github.com/cloudfoundry/cli/cf/terminal"
)

// MonitorAction monitors process execution
type MonitorAction struct{}

// Execute executes monitor action on process with the specified id
func (a *MonitorAction) Execute(operationID, commandName string, mtaClient mtaclient.MtaClientOperations) ExecutionStatus {
	// TODO: refactor the Execution monitor in order to use the new client
	// TODO: introduce a new way of building the reporting

	responseHeader, err = mtaClient.ExecuteAction(operationID, "monitor")
	if err != nil {
		ui.Failed("Could not monitor multi-target app operation with id %s: %s", terminal.EntityNameColor(operationID), err)
		return Failure
	}
	ui.Ok()

	monitor := NewExecutionMonitor(operationID, commandName, responseHeader.Location, mtaClient)
	return monitor.Monitor()
}
