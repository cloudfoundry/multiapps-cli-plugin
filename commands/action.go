package commands

import (
	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/clients/baseclient"
	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/clients/csrf"
	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/clients/models"
	mtaclient "github.com/cloudfoundry-incubator/multiapps-cli-plugin/clients/mtaclient"
	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/ui"
	"github.com/cloudfoundry/cli/cf/terminal"
)

// Action interface representing actions to be excuted on processes
type Action interface {
	Execute(operationID string, mtaClient mtaclient.MtaClientOperations, sessionProvider csrf.SessionProvider) ExecutionStatus
}

// GetActionToExecute returns the action to execute specified with action id
func GetActionToExecute(actionID, commandName string) Action {
	switch actionID {
	case "abort":
		action := newAction(actionID)
		return &action
	case "retry":
		action := newMonitoringAction(actionID, commandName)
		return &action
	case "resume":
		action := newMonitoringAction(actionID, commandName)
		return &action
	case "monitor":
		return &MonitorAction{}
	}
	return nil
}

func newMonitoringAction(actionID, commandName string) monitoringAction {
	return monitoringAction{
		action:      newAction(actionID),
		commandName: commandName,
	}
}

func newAction(actionID string) action {
	return action{
		actionID: actionID,
	}
}

type action struct {
	actionID string
}

func (a *action) Execute(operationID string, mtaClient mtaclient.MtaClientOperations, sessionProvider csrf.SessionProvider) ExecutionStatus {
	err := sessionProvider.GetSession()
	if err != nil {
		ui.Failed("Could not retrieve CSRF token for the current session: %s", baseclient.NewClientError(err))
		return Failure
	}
	return a.executeInSession(operationID, mtaClient)
}

func (a *action) executeInSession(operationID string, mtaClient mtaclient.MtaClientOperations) ExecutionStatus {
	ui.Say("Executing action '%s' on operation %s...", a.actionID, terminal.EntityNameColor(operationID))
	_, err := mtaClient.ExecuteAction(operationID, a.actionID)
	if err != nil {
		ui.Failed("Could not execute action '%s' on operation %s: %s", a.actionID, terminal.EntityNameColor(operationID), err)
		return Failure
	}
	ui.Ok()
	return Success
}

type monitoringAction struct {
	action
	commandName string
}

func (a *monitoringAction) Execute(operationID string, mtaClient mtaclient.MtaClientOperations, sessionProvider csrf.SessionProvider) ExecutionStatus {
	err := sessionProvider.GetSession()
	if err != nil {
		ui.Failed("Could not retrieve CSRF token for the current session: %s", baseclient.NewClientError(err))
		return Failure
	}

	// Get the messages of the operation before it's retried/resumed, so that the monitor knows they're from the previous execution and
	// should not show them again.
	operation, err := getMonitoringOperation(operationID, mtaClient)
	if err != nil {
		ui.Failed("Could not monitor multi-target app operation with id %s: %s", terminal.EntityNameColor(operationID), baseclient.NewClientError(err))
		return Failure
	}

	status := a.executeInSession(operationID, mtaClient)
	if status == Failure {
		return status
	}

	return NewExecutionMonitor(a.commandName, operationID, "messages", operation.Messages, mtaClient).Monitor()
}

func getMonitoringOperation(operationID string, mtaClient mtaclient.MtaClientOperations) (*models.Operation, error) {
	return mtaClient.GetMtaOperation(operationID, "messages")
}
