package commands

import (
	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/clients/baseclient"
	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/clients/csrf"
	mtaclient "github.com/cloudfoundry-incubator/multiapps-cli-plugin/clients/mtaclient"
	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/ui"
	"github.com/cloudfoundry/cli/cf/terminal"
)

// AbortAction represents action abort
type AbortAction struct{}

// Execute aborts operation with the specified operation id
func (a *AbortAction) Execute(operationID, commandName string, mtaClient mtaclient.MtaClientOperations, sessionProvider csrf.SessionProvider) ExecutionStatus {

	err := sessionProvider.GetSession()
	if err != nil {
		ui.Failed("Could not retrieve x-csrf-token for the current session: %s", baseclient.NewClientError(err))
		return Failure
	}
	ui.Say("Aborting multi-target app operation with id %s...", terminal.EntityNameColor(operationID))
	_, err = mtaClient.ExecuteAction(operationID, "abort")
	if err != nil {
		ui.Failed("Could not abort multi-target app operation with id %s: %s", terminal.EntityNameColor(operationID), err)
		return Failure
	}
	ui.Ok()
	return Success
}
