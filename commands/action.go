package commands

import slppclient "github.com/SAP/cf-mta-plugin/clients/slppclient"

// Action interface representing actions to be excuted on processes
type Action interface {
	Execute(processID, commandName string, slppClient slppclient.SlppClientOperations) ExecutionStatus
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
