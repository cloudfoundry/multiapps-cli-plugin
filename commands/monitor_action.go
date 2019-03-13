package commands

import (
	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/clients/mtaclient"
	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/ui"
	"github.com/cloudfoundry/cli/cf/terminal"
)

// MonitorAction monitors process execution
type MonitorAction struct {
	commandName string
}

// Execute executes monitor action on process with the specified id
func (a *MonitorAction) Execute(operationID string, mtaClient mtaclient.MtaClientOperations) ExecutionStatus {
	operation, err := getMonitoringOperation(operationID, mtaClient)
	if err != nil {
		ui.Failed("Could not monitor operation %s: %s", terminal.EntityNameColor(operationID), err)
		return Failure
	}

	return NewExecutionMonitor(a.commandName, operationID, "messages", operation.Messages, mtaClient).Monitor()
}
