package commands

import (
	"github.com/cloudfoundry/cli/cf/terminal"
	"github.com/SAP/cf-mta-plugin/clients/models"
	slppclient "github.com/SAP/cf-mta-plugin/clients/slppclient"
	"github.com/SAP/cf-mta-plugin/ui"
)

// RetryAction retries the process with the specified id
type RetryAction struct{}

// Execute executes retry action on process with the specified id
func (a *RetryAction) Execute(processID, commandName string, slppClient slppclient.SlppClientOperations) ExecutionStatus {

	// Ensure session is not expired
	EnsureSlppSession(slppClient)

	progressMessages, err := getReportedProgressMessages(slppClient)
	if err != nil {
		ui.Failed("Could not get the already reported progress messages for process %s: %s", terminal.EntityNameColor(processID), err)
	}

	ui.Say("Retrying multi-target app operation with id %s...", terminal.EntityNameColor(processID))
	err = slppClient.ExecuteAction("retry")
	if err != nil {
		ui.Failed("Could not retry multi-target app operation with id %s: %s", terminal.EntityNameColor(processID), err)
		return Failure
	}
	ui.Ok()

	monitor := NewExecutionMonitor(processID, commandName, slppClient, progressMessages)
	return monitor.Monitor()
}

func getReportedProgressMessages(slppClient slppclient.SlppClientOperations) ([]*models.ProgressMessage, error) {
	processTaskID := slppClient.GetServiceID()
	processTask, err := slppClient.GetTasklistTask(processTaskID)
	if err != nil {
		return []*models.ProgressMessage{}, err
	}
	return processTask.ProgressMessages.ProgressMessages, nil
}
