package commands

import (
	"strconv"
	"time"

	"github.com/cloudfoundry/cli/cf/terminal"
	slppclient "github.com/SAP/cf-mta-plugin/clients/slppclient"
	"github.com/SAP/cf-mta-plugin/ui"
	"github.com/SAP/cf-mta-plugin/util"

	"github.com/SAP/cf-mta-plugin/clients/models"
)

const consoleOffset = "  "

//ExecutionMonitor monitors execution of a process
type ExecutionMonitor struct {
	slppClient               slppclient.SlppClientOperations
	reportedProgressMessages map[int]bool
	processID                string
	commandName              string
}

//NewExecutionMonitor creates a new execution monitor
func NewExecutionMonitor(processID, commandName string, slppClient slppclient.SlppClientOperations, reportedProgressMessages []*models.ProgressMessage) *ExecutionMonitor {
	return &ExecutionMonitor{
		slppClient:               slppClient,
		reportedProgressMessages: getProgressMessagesIds(reportedProgressMessages),
		processID:                processID,
		commandName:              commandName,
	}
}

func getProgressMessagesIds(reportedProgressMessages []*models.ProgressMessage) map[int]bool {
	result := make(map[int]bool)
	for _, progressMessage := range reportedProgressMessages {
		idNum, _ := strconv.Atoi(*progressMessage.ID)
		result[idNum] = true
	}
	return result
}

//Monitor monitors current state of the execution
func (m *ExecutionMonitor) Monitor() ExecutionStatus {
	ui.Say("Monitoring process execution...")

	for {
		processTask, err := m.slppClient.GetTasklistTask(m.slppClient.GetServiceID())
		if err != nil {
			ui.Failed("Could not get task: %s", err)
			return Failure
		}
		if processTask.Type != models.SlpTaskTypeSlpTaskTypePROCESS {
			ui.Failed("The SLP task must be a Process task")
			return Failure
		}
		m.reportProgressMessages(*processTask)
		switch processTask.Status {
		case models.SlpTaskStateSlpTaskStateRUNNING:
			time.Sleep(2000)
		case models.SlpTaskStateSlpTaskStateFINISHED:
			ui.Say("Process finished.")
			return Success
		case models.SlpTaskStateSlpTaskStateABORTED:
			ui.Say("Process was aborted.")
			return Failure
		case models.SlpTaskStateSlpTaskStateERROR:
			clientError, err := m.slppClient.GetError()
			if err != nil {
				ui.Failed("Could not get client error: %s", err)
				return Failure
			}
			ui.Say("Process failed: %s", clientError.Description)
			m.reportAvaiableActions()
			m.reportCommandForDownloadOfProcessLogs()
			return Failure
		case models.SlpTaskStateSlpTaskStateACTIONREQUIRED, models.SlpTaskStateSlpTaskStateDIALOG:
			ui.Say("Process has entered validation phase. After testing your new deployment you can resume or abort the process.")
			m.reportAvaiableActions()
			ui.Say("Hint: Use the '--no-confirm' option of the bg-deploy command to skip this phase.")
			return Success
		default:
			ui.Failed("Process is in illegal state %s.", terminal.EntityNameColor(string(processTask.Status)))
			return Failure
		}
	}
	ui.Failed("Error during process monitoring")
	return Failure
}

func (m *ExecutionMonitor) reportAvaiableActions() {
	actions, _ := m.slppClient.GetActions()
	actionsList := actions.Actions
	for i := 0; i < len(actionsList); i++ {
		m.reportAvailableAction(*actionsList[i])
	}
}

func (m *ExecutionMonitor) reportCommandForDownloadOfProcessLogs() {
	downloadProcessLogsCommand := DownloadMtaOperationLogsCommand{}
	commandBuilder := util.NewCfCommandStringBuilder()
	commandBuilder.SetName(downloadProcessLogsCommand.GetPluginCommand().Alias)
	commandBuilder.AddOption(operationIDOpt, m.processID)
	ui.Say("Use \"%s\" to download the logs of the process", commandBuilder.Build())
}

func (m *ExecutionMonitor) reportAvailableAction(action models.Action) {
	commandBuilder := util.NewCfCommandStringBuilder()
	commandBuilder.SetName(m.commandName)
	commandBuilder.AddOption(operationIDOpt, m.processID)
	commandBuilder.AddOption(actionOpt, *action.ID)
	ui.Say("Use \"%s\" to %s the process.", commandBuilder.Build(), *action.ID)
}

func (m *ExecutionMonitor) reportProgressMessages(task models.Task) {
	for _, progressMessage := range task.ProgressMessages.ProgressMessages {
		m.reportProgressMessage(progressMessage)
	}
}

func (m *ExecutionMonitor) reportProgressMessage(progressMessage *models.ProgressMessage) {
	idNum, _ := strconv.Atoi(*progressMessage.ID)
	if m.reportedProgressMessages[idNum] {
		return
	}

	m.reportedProgressMessages[idNum] = true
	ui.Say("%s%s", consoleOffset, *progressMessage.Message)
}

func contains(slice []string, element string) bool {
	for _, sliceElement := range slice {
		if sliceElement == element {
			return true
		}
	}
	return false
}
