package commands

import (
	"strings"
	"time"

	"net/url"

	"github.com/SAP/cf-mta-plugin/ui"
	"github.com/SAP/cf-mta-plugin/util"
	"github.com/cloudfoundry/cli/cf/terminal"

	"github.com/SAP/cf-mta-plugin/clients/models"
	mtaclient "github.com/SAP/cf-mta-plugin/clients/mtaclient"
)

const consoleOffset = "  "

//ExecutionMonitor monitors execution of a process
type ExecutionMonitor struct {
	mtaClient          mtaclient.MtaClientOperations
	reportedMessages   map[int64]bool
	operationID        string
	commandName        string
	monitoringLocation string
}

//NewExecutionMonitor creates a new execution monitor
func NewExecutionMonitor(commandName, monitoringLocation string, reportedOperationMessages []*models.Message, mtaClient mtaclient.MtaClientOperations) *ExecutionMonitor {
	return &ExecutionMonitor{
		mtaClient:          mtaClient,
		reportedMessages:   getAlreadyReportedOperationMessages(reportedOperationMessages),
		commandName:        commandName,
		monitoringLocation: monitoringLocation,
	}
}

func getAlreadyReportedOperationMessages(reportedOperationMessages []*models.Message) map[int64]bool {
	result := make(map[int64]bool)
	for _, message := range reportedOperationMessages {
		result[message.ID] = true
	}
	return result
}

func (m *ExecutionMonitor) Monitor() ExecutionStatus {
	ui.Say("Monitoring process execution...")

	for {
		operation, err := getOperation(m.monitoringLocation, m.mtaClient)
		if err != nil {
			ui.Failed("Could not get ongoing operation: %s", err)
			return Failure
		}
		m.reportOperationMessages(operation)
		switch operation.State {
		case models.StateRUNNING:
			time.Sleep(2000)
		case models.StateFINISHED:
			ui.Say("Process finished.")
			return Success
		case models.StateABORTED:
			ui.Say("Process was aborted.")
			return Failure
		case models.StateERROR:
			messageInError := findErrorMessage(operation.Messages)
			if messageInError == nil {
				ui.Failed("There is not error message for operation with id %s", operation.ProcessID)
				return Failure
			}
			ui.Say("Process failed: %s", messageInError.Message)
			m.reportAvaiableActions(operation.ProcessID)
			m.reportCommandForDownloadOfProcessLogs(operation.ProcessID)
			return Failure
		case models.StateACTIONREQUIRED:
			ui.Say("Process has entered validation phase. After testing your new deployment you can resume or abort the process.")
			m.reportAvaiableActions(operation.ProcessID)
			ui.Say("Hint: Use the '--no-confirm' option of the bg-deploy command to skip this phase.")
			return Success
		default:
			ui.Failed("Process is in illegal state %s.", terminal.EntityNameColor(string(operation.State)))
			return Failure
		}
	}
}

func findErrorMessage(messages models.OperationMessages) *models.Message {
	for _, message := range messages {
		if message.Type == models.MessageTypeERROR {
			return message
		}
	}
	return nil
}

func (m *ExecutionMonitor) reportOperationMessages(operation *models.Operation) {
	for _, message := range operation.Messages {
		if m.reportedMessages[message.ID] {
			continue
		}
		m.reportedMessages[message.ID] = true
		ui.Say("%s%s", consoleOffset, message.Message)
	}
}

func getMonitoringInformation(monitoringLocation string) (string, string, error) {
	parsedUrl, _ := url.Parse(monitoringLocation)
	path := parsedUrl.Path
	parsedQuery, _ := url.ParseQuery(parsedUrl.RawQuery)
	return strings.Split(path, "operations/")[1], parsedQuery["embed"][0], nil
}

func getOperation(monitoringLocation string, mtaClient mtaclient.MtaClientOperations) (*models.Operation, error) {
	operationID, embedMessages, _ := getMonitoringInformation(monitoringLocation)
	return mtaClient.GetMtaOperation(operationID, embedMessages)
}

func (m *ExecutionMonitor) reportAvaiableActions(operationID string) {
	actions, _ := m.mtaClient.GetOperationActions(operationID)
	for _, action := range actions {
		m.reportAvailableAction(action, operationID)
	}
}

func (m *ExecutionMonitor) reportCommandForDownloadOfProcessLogs(operationID string) {
	downloadProcessLogsCommand := DownloadMtaOperationLogsCommand{}
	commandBuilder := util.NewCfCommandStringBuilder()
	commandBuilder.SetName(downloadProcessLogsCommand.GetPluginCommand().Alias)
	commandBuilder.AddOption(operationIDOpt, operationID)
	ui.Say("Use \"%s\" to download the logs of the process", commandBuilder.Build())
}

func (m *ExecutionMonitor) reportAvailableAction(action, operationID string) {
	commandBuilder := util.NewCfCommandStringBuilder()
	commandBuilder.SetName(m.commandName)
	commandBuilder.AddOption(operationIDOpt, operationID)
	commandBuilder.AddOption(actionOpt, action)
	ui.Say("Use \"%s\" to %s the process.", commandBuilder.Build(), action)
}
