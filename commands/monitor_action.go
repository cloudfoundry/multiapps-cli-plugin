package commands

import (
	"github.com/cloudfoundry/cli/cf/terminal"
	slppclient "github.com/SAP/cf-mta-plugin/clients/slppclient"
	"github.com/SAP/cf-mta-plugin/ui"
)

// MonitorAction monitors process execution
type MonitorAction struct{}

// Execute executes monitor action on process with the specified id
func (a *MonitorAction) Execute(processID, commandName string, slppClient slppclient.SlppClientOperations) ExecutionStatus {
	progressMessages, err := getReportedProgressMessages(slppClient)
	if err != nil {
		ui.Failed("Could not get the already reported progress messages for process %s: %s", terminal.EntityNameColor(processID), err)
	}

	monitor := NewExecutionMonitor(processID, commandName, slppClient, progressMessages)
	return monitor.Monitor()
}
