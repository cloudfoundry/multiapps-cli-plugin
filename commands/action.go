package commands

import (
	"code.cloudfoundry.org/cli/cf/terminal"
	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/clients/baseclient"
	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/clients/models"
	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/clients/mtaclient"
	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/ui"
)

// Action interface representing actions to be excuted on processes
type Action interface {
	Execute(operationID string, mtaClient mtaclient.MtaClientOperations) ExecutionStatus
}

// GetActionToExecute returns the action to execute specified with action ID
func GetActionToExecute(actionID, commandName string, monitoringRetries uint) Action {
	switch actionID {
	case "abort":
		action := newAction(actionID, VerbosityLevelVERBOSE)
		return &action
	case "retry":
		action := newMonitoringAction(actionID, commandName, VerbosityLevelVERBOSE, monitoringRetries)
		return &action
	case "resume":
		action := newMonitoringAction(actionID, commandName, VerbosityLevelVERBOSE, monitoringRetries)
		return &action
	case "monitor":
		return &MonitorAction{
			commandName:       commandName,
			monitoringRetries: monitoringRetries,
		}
	}
	return nil
}

func GetNoRetriesActionToExecute(actionID, commandName string) Action {
	return GetActionToExecute(actionID, commandName, 0)
}

func newMonitoringAction(actionID, commandName string, verbosityLevel VerbosityLevel, monitoringRetries uint) monitoringAction {
	return monitoringAction{
		action:            newAction(actionID, verbosityLevel),
		commandName:       commandName,
		monitoringRetries: monitoringRetries,
	}
}

func newAction(actionID string, verbosityLevel VerbosityLevel) action {
	return action{
		actionID:       actionID,
		verbosityLevel: verbosityLevel,
	}
}

type VerbosityLevel int

const (
	VerbosityLevelVERBOSE VerbosityLevel = 0
	VerbosityLevelSILENT  VerbosityLevel = 1
)

type action struct {
	actionID       string
	verbosityLevel VerbosityLevel
}

func (a *action) Execute(operationID string, mtaClient mtaclient.MtaClientOperations) ExecutionStatus {
	return a.executeInSession(operationID, mtaClient)
}

func (a *action) executeInSession(operationID string, mtaClient mtaclient.MtaClientOperations) ExecutionStatus {
	if a.verbosityLevel == VerbosityLevelVERBOSE {
		ui.Say("Executing action %q on operation %s...", a.actionID, terminal.EntityNameColor(operationID))
	}
	_, err := mtaClient.ExecuteAction(operationID, a.actionID)
	if err != nil {
		ui.Failed("Could not execute action %q on operation %s: %s", a.actionID, terminal.EntityNameColor(operationID), err)
		return Failure
	}
	if a.verbosityLevel == VerbosityLevelVERBOSE {
		ui.Ok()
	}
	return Success
}

type monitoringAction struct {
	action
	commandName       string
	monitoringRetries uint
}

func (a *monitoringAction) Execute(operationID string, mtaClient mtaclient.MtaClientOperations) ExecutionStatus {
	// Get the messages of the operation before it's retried/resumed, so that the monitor knows they're from the previous execution and
	// should not show them again.
	operation, err := getMonitoringOperation(operationID, mtaClient)
	if err != nil {
		ui.Failed("Could not monitor multi-target app operation with ID %s: %s", terminal.EntityNameColor(operationID), baseclient.NewClientError(err))
		return Failure
	}

	status := a.executeInSession(operationID, mtaClient)
	if status == Failure {
		return status
	}

	return NewExecutionMonitor(a.commandName, operationID, "messages", a.monitoringRetries, operation.Messages, mtaClient).Monitor()
}

func getMonitoringOperation(operationID string, mtaClient mtaclient.MtaClientOperations) (*models.Operation, error) {
	return mtaClient.GetMtaOperation(operationID, "messages")
}
