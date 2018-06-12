package commands

import (
	"strings"
	"time"

	"net/url"

	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/ui"
	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/util"
	"github.com/cloudfoundry/cli/cf/terminal"

	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/clients/baseclient"
	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/clients/models"
	mtaclient "github.com/cloudfoundry-incubator/multiapps-cli-plugin/clients/mtaclient"
)

const consoleOffset = "  "

//ExecutionMonitor monitors execution of a process
type ExecutionMonitor struct {
	mtaClient          mtaclient.MtaClientOperations
	reportedMessages   map[int64]bool
	commandName        string
	monitoringLocation string
	operationID        string
	embed              string
}

func NewExecutionMonitorFromLocationHeader(commandName, location string, reportedOperationMessages []*models.Message, mtaClient mtaclient.MtaClientOperations) *ExecutionMonitor {
	operationID, embed := getMonitoringInformation(location)
	return &ExecutionMonitor{
		mtaClient:        mtaClient,
		reportedMessages: getAlreadyReportedOperationMessages(reportedOperationMessages),
		commandName:      commandName,
		operationID:      operationID,
		embed:            embed,
	}
}

func getMonitoringInformation(monitoringLocation string) (string, string) {
	parsedURL, _ := url.Parse(monitoringLocation)
	path := parsedURL.Path
	parsedQuery, _ := url.ParseQuery(parsedURL.RawQuery)
	return strings.Split(path, "operations/")[1], parsedQuery["embed"][0]
}

//NewExecutionMonitor creates a new execution monitor
func NewExecutionMonitor(commandName, operationID, embed string, reportedOperationMessages []*models.Message, mtaClient mtaclient.MtaClientOperations) *ExecutionMonitor {
	return &ExecutionMonitor{
		mtaClient:        mtaClient,
		reportedMessages: getAlreadyReportedOperationMessages(reportedOperationMessages),
		commandName:      commandName,
		operationID:      operationID,
		embed:            embed,
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
	ui.Say("Monitoring process %s...", m.operationID)

	for {
		operation, err := m.mtaClient.GetMtaOperation(m.operationID, m.embed)
		if err != nil {
			ui.Failed("Could not get ongoing operation: %s", baseclient.NewClientError(err))
			return Failure
		}
		m.reportOperationMessages(operation)
		switch operation.State {
		case models.StateRUNNING:
			time.Sleep(2000)
		case models.StateFINISHED:
			ui.Say("Process finished.")
			m.reportCommandForDownloadOfProcessLogs(m.operationID)
			return Success
		case models.StateABORTED:
			ui.Say("Process was aborted.")
			m.reportCommandForDownloadOfProcessLogs(m.operationID)
			return Failure
		case models.StateERROR:
			messageInError := findErrorMessage(operation.Messages)
			if messageInError == nil {
				ui.Failed("There is not error message for operation with id %s", m.operationID)
				return Failure
			}
			ui.Say("Process failed: %s", messageInError.Text)
			m.reportAvaiableActions(m.operationID)
			m.reportCommandForDownloadOfProcessLogs(m.operationID)
			return Failure
		case models.StateACTIONREQUIRED:
			ui.Say("Process has entered validation phase. After testing your new deployment you can resume or abort the process.")
			m.reportAvaiableActions(m.operationID)
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
		ui.Say("%s%s", consoleOffset, message.Text)
	}
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
	ui.Say("Use \"%s\" to download the logs of the process.", commandBuilder.Build())
}

func (m *ExecutionMonitor) reportAvailableAction(action, operationID string) {
	commandBuilder := util.NewCfCommandStringBuilder()
	commandBuilder.SetName(m.commandName)
	commandBuilder.AddOption(operationIDOpt, operationID)
	commandBuilder.AddOption(actionOpt, action)
	ui.Say("Use \"%s\" to %s the process.", commandBuilder.Build(), action)
}
