package commands

import (
	"code.cloudfoundry.org/cli/cf/terminal"
	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/clients/mtaclient"
	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/ui"
)

// MonitorAction monitors process execution
type MonitorAction struct {
	commandName       string
	monitoringRetries uint
}

// Execute executes monitor action on process with the specified id
func (a *MonitorAction) Execute(operationID string, mtaClient mtaclient.MtaClientOperations) ExecutionStatus {
	operation, err := getMonitoringOperation(operationID, mtaClient)
	if err != nil {
		ui.Failed("Could not monitor operation %s: %s", terminal.EntityNameColor(operationID), err)
		return Failure
	}

	return NewExecutionMonitor(a.commandName, operationID, "messages", a.monitoringRetries, operation.Messages, mtaClient).Monitor()
}
