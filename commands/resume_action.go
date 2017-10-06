package commands

import (
	"github.com/cloudfoundry/cli/cf/terminal"
	slppclient "github.com/SAP/cf-mta-plugin/clients/slppclient"
	"github.com/SAP/cf-mta-plugin/ui"
)

// ResumeAction retries the process with the specified id
type ResumeAction struct{}

// Execute executes resume action on process with the specified id
func (a *ResumeAction) Execute(processID, commandName string, slppClient slppclient.SlppClientOperations) ExecutionStatus {

	// Ensure session is not expired
	EnsureSlppSession(slppClient)

	progressMessages, err := getReportedProgressMessages(slppClient)
	if err != nil {
		ui.Failed("Could not get the already reported progress messages for process %s: %s", terminal.EntityNameColor(processID), err)
	}

	ui.Say("Resuming multi-target app operation with id %s...", terminal.EntityNameColor(processID))
	err = slppClient.ExecuteAction("resume")
	if err != nil {
		ui.Failed("Could not resume multi-target app operation with id %s: %s", terminal.EntityNameColor(processID), err)
		return Failure
	}
	ui.Ok()

	monitor := NewExecutionMonitor(processID, commandName, slppClient, progressMessages)
	return monitor.Monitor()
}
