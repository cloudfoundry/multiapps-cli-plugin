package commands

import (
	"github.com/SAP/cf-mta-plugin/clients/csrf"
	mtaclient "github.com/SAP/cf-mta-plugin/clients/mtaclient"
)

// Action interface representing actions to be excuted on processes
type Action interface {
	Execute(operationID, commandName string, mtaClient mtaclient.MtaClientOperations, sessionProvider csrf.SessionProvider) ExecutionStatus
}

// GetActionToExecute returns the action to execute specified with action id
func GetActionToExecute(actionID string) Action {
	switch actionID {
	case "abort":
		return &AbortAction{}
	case "retry":
		return &RetryAction{}
	case "monitor":
		return &MonitorAction{}
	case "resume":
		return &ResumeAction{}
	}

	return nil
}
