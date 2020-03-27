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

//ExecutionMonitor monitors execution of a process
type ExecutionMonitor struct {
	mtaClient          mtaclient.MtaClientOperations
	reportedMessages   map[int64]bool
	commandName        string
	monitoringLocation string
	operationID        string
	embed              string
	retries            uint
}

func NewExecutionMonitorFromLocationHeader(commandName, location string, retries uint, reportedOperationMessages []*models.Message, mtaClient mtaclient.MtaClientOperations) *ExecutionMonitor {
	operationID, embed := getMonitoringInformation(location)
	return &ExecutionMonitor{
		mtaClient:        mtaClient,
		reportedMessages: getAlreadyReportedOperationMessages(reportedOperationMessages),
		commandName:      commandName,
		operationID:      operationID,
		embed:            embed,
		retries:          retries,
	}
}

func getMonitoringInformation(monitoringLocation string) (string, string) {
	parsedURL, _ := url.Parse(monitoringLocation)
	path := parsedURL.Path
	parsedQuery, _ := url.ParseQuery(parsedURL.RawQuery)
	return strings.Split(path, "operations/")[1], parsedQuery["embed"][0]
}

//NewExecutionMonitor creates a new execution monitor
func NewExecutionMonitor(commandName, operationID, embed string, retries uint, reportedOperationMessages []*models.Message, mtaClient mtaclient.MtaClientOperations) *ExecutionMonitor {
	return &ExecutionMonitor{
		mtaClient:        mtaClient,
		reportedMessages: getAlreadyReportedOperationMessages(reportedOperationMessages),
		commandName:      commandName,
		operationID:      operationID,
		embed:            embed,
		retries:          retries,
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
	totalRetries := m.retries
	for {
		operation, err := m.mtaClient.GetMtaOperation(m.operationID, m.embed)
		if err != nil {
			ui.Failed("Could not get ongoing operation: %s", baseclient.NewClientError(err))
			return Failure
		}
		m.reportOperationMessages(operation)
		switch operation.State {
		case models.StateRUNNING:
			time.Sleep(3 * time.Second)
		case models.StateFINISHED:
			ui.Say("Process finished.")
			m.reportCommandForDownloadOfProcessLogs(m.operationID)
			return Success
		case models.StateABORTED:
			ui.Say("Process was aborted.")
			m.reportCommandForDownloadOfProcessLogs(m.operationID)
			return Failure
		case models.StateERROR:
			if canRetry(m.retries, operation) {
				ui.Say("Proceeding with automatic retry... (%d of %d attempts left)", m.retries, totalRetries)
				executeRetryAction(m)
				continue
			}
			messageInError := findErrorMessage(operation.Messages)
			if messageInError == nil {
				ui.Failed("There is no error message for operation with ID %s", m.operationID)
				return Failure
			}
			ui.Say("Process failed.")
			m.reportAvaiableActions(m.operationID)
			m.reportCommandForDownloadOfProcessLogs(m.operationID)
			return Failure
		case models.StateACTIONREQUIRED:
			intermediatePhase, flag := getIntermediatePhaseAndFlag(m.commandName)
			ui.Say("Process has entered %s phase. After testing your new deployment you can resume or abort the process.", intermediatePhase)
			m.reportAvaiableActions(m.operationID)
			ui.Say("Hint: Use the '%s' option of the %s command to skip this phase.", flag, m.commandName)
			return Success
		default:
			ui.Failed("Process is in illegal state %s.", terminal.EntityNameColor(string(operation.State)))
			return Failure
		}
	}
}

func getIntermediatePhaseAndFlag(commandName string) (string, string) {
	//for backwards compatibility until the bg-deploy deprecation period expires
	if commandName == "bg-deploy" {
		return "validation", "--no-confirm"
	}
	return "testing", "--skip-testing-phase"
}

func canRetry(retries uint, operation *models.Operation) bool {
	return retries > 0 && operation.ErrorType != models.ErrorTypeCONTENT
}

func executeRetryAction(executionMonitor *ExecutionMonitor) {
	retryAction := newAction("retry", VerbosityLevelSILENT)
	retryAction.Execute(executionMonitor.operationID, executionMonitor.mtaClient)
	executionMonitor.retries--
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
		ui.Say("%s", message.Text)
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
