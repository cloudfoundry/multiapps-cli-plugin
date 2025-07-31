package commands

import (
	"flag"
	"fmt"
	"strconv"
	"strings"

	"code.cloudfoundry.org/cli/cf/terminal"
	"code.cloudfoundry.org/cli/plugin"
	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/clients/baseclient"
	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/clients/models"
	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/ui"
	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/util"
)

const processUserProvidedServicesOpt = "process-user-provided-services"

// RollbackMtaCommand is a command for rollback of deployed MTAs.
type RollbackMtaCommand struct {
	*BaseCommand
	processTypeProvider ProcessTypeProvider
}

func NewRollbackMtaCommand() *RollbackMtaCommand {
	baseCmd := &BaseCommand{flagsParser: NewProcessActionExecutorCommandArgumentsParser([]string{"MTA_ID"}), flagsValidator: NewDefaultCommandFlagsValidator(nil)}
	rollbackMtaCmd := &RollbackMtaCommand{baseCmd, &rollbackMtaCommandProcessTypeProvider{}}
	baseCmd.Command = rollbackMtaCmd
	return rollbackMtaCmd
}

// GetPluginCommand returns more information for the blue green deploy command.
func (c *RollbackMtaCommand) GetPluginCommand() plugin.Command {
	return plugin.Command{
		Name:     "rollback-mta",
		HelpText: "Rollback of a multi-target app works only if [--backup-previous-version] flag was used during blue-green deployment and backup applications exists in the space",
		UsageDetails: plugin.Usage{
			Usage: `Rollback of a multi-target app
   cf rollback-mta MTA_ID [-t TIMEOUT] [-f] [--retries RETRIES] [--namespace NAMESPACE] [--do-not-fail-on-missing-permissions] [--abort-on-error] [--apps-start-timeout TIMEOUT] [--apps-stage-timeout TIMEOUT] [--apps-upload-timeout TIMEOUT] [--apps-task-execution-timeout TIMEOUT]

   Perform action on an active deploy operation
   cf rollback-mta -i OPERATION_ID -a ACTION [-u URL]` + util.BaseEnvHelpText,
			Options: map[string]string{
				deployServiceURLOpt:                    "Deploy service URL, by default 'deploy-service.<system-domain>'",
				operationIDOpt:                         "Active deploy operation ID",
				actionOpt:                              "Action to perform on active deploy operation (abort, retry, resume, monitor)",
				forceOpt:                               "Force deploy without confirmation for aborting conflicting processes",
				util.GetShortOption(namespaceOpt):      "(EXPERIMENTAL) Namespace for the MTA, applied on app names, app routes and service names",
				util.GetShortOption(deleteServicesOpt): "Recreate changed services / delete discontinued services",
				util.GetShortOption(noFailOnMissingPermissionsOpt):              "Do not fail on missing permissions for admin operations",
				util.GetShortOption(abortOnErrorOpt):                            "Auto-abort the process on any errors",
				util.GetShortOption(processUserProvidedServicesOpt):             "Enable processing of user provided services during rollback",
				util.GetShortOption(retriesOpt):                                 "Retry the operation N times in case a non-content error occurs (default 3)",
				util.GetShortOption(stageTimeoutOpt):                            "Stage app timeout in seconds",
				util.GetShortOption(uploadTimeoutOpt):                           "Upload app timeout in seconds",
				util.GetShortOption(taskExecutionTimeoutOpt):                    "Task execution timeout in seconds",
				util.CombineFullAndShortParameters(startTimeoutOpt, timeoutOpt): "Start app timeout in seconds",
			},
		},
	}
}

func (c *RollbackMtaCommand) defineCommandOptions(flags *flag.FlagSet) {
	flags.Bool(forceOpt, false, "")
	flags.String(operationIDOpt, "", "")
	flags.String(namespaceOpt, "", "")
	flags.Bool(deleteServicesOpt, false, "")
	flags.String(actionOpt, "", "")
	flags.Bool(noFailOnMissingPermissionsOpt, false, "")
	flags.Bool(abortOnErrorOpt, false, "")
	flags.Bool(processUserProvidedServicesOpt, false, "")
	flags.Uint(retriesOpt, 3, "")
	flags.String(startTimeoutOpt, "", "")
	flags.String(stageTimeoutOpt, "", "")
	flags.String(uploadTimeoutOpt, "", "")
	flags.String(taskExecutionTimeoutOpt, "", "")
}

func (c *RollbackMtaCommand) executeInternal(positionalArgs []string, dsHost string, flags *flag.FlagSet, cfTarget util.CloudFoundryTarget) ExecutionStatus {
	operationID := GetStringOpt(operationIDOpt, flags)
	actionID := GetStringOpt(actionOpt, flags)
	retries := GetUintOpt(retriesOpt, flags)

	if operationID != "" || actionID != "" {
		return c.ExecuteAction(operationID, actionID, retries, dsHost, cfTarget)
	}

	force := GetBoolOpt(forceOpt, flags)
	mtaID := positionalArgs[0]
	if !force && !ui.Confirm("Really rollback multi-target app %s in org %s / space %s? (y/n)", terminal.EntityNameColor(mtaID),
		terminal.EntityNameColor(cfTarget.Org.Name),
		terminal.EntityNameColor(cfTarget.Space.Name)) {
		ui.Warn("Rollback mta cancelled")
		return Failure
	}

	// Print initial message
	ui.Say("Rollback multi-target app %s in org %s / space %s as %s...",
		terminal.EntityNameColor(mtaID), terminal.EntityNameColor(cfTarget.Org.Name),
		terminal.EntityNameColor(cfTarget.Space.Name), terminal.EntityNameColor(cfTarget.Username))

	// Create rest client
	mtaClient := c.NewMtaClient(dsHost, cfTarget)
	// Create new REST client for mtas V2 api
	mtaV2Client := c.NewMtaV2Client(dsHost, cfTarget)

	namespace := strings.TrimSpace(GetStringOpt(namespaceOpt, flags))

	// Check if a deployed MTA with the specified ID exists
	_, err := mtaV2Client.GetMtasForThisSpace(&mtaID, &namespace)
	if err != nil {
		ce, ok := err.(*baseclient.ClientError)
		if ok && ce.Code == 404 && strings.Contains(fmt.Sprint(ce.Description), mtaID) {
			if util.DiscardIfEmpty(namespace) != nil {
				ui.Failed("Multi-target app %s with namespace %s not found", terminal.EntityNameColor(mtaID), terminal.EntityNameColor(namespace))
			} else {
				ui.Failed("Multi-target app %s not found", terminal.EntityNameColor(mtaID))
			}
			return Failure
		}
		ui.Failed("Could not get multi-target app %s: %s", terminal.EntityNameColor(mtaID), baseclient.NewClientError(err))
		return Failure
	}

	// Check for an ongoing operation for this MTA ID and abort it
	wasAborted, err := c.CheckOngoingOperation(mtaID, namespace, dsHost, force, cfTarget)
	if err != nil {
		ui.Failed(err.Error())
		return Failure
	}
	if !wasAborted {
		return Failure
	}

	processBuilder := util.NewProcessBuilder()
	processBuilder.ProcessType(c.processTypeProvider.GetProcessType())
	processBuilder.Parameter("mtaId", mtaID)
	processBuilder.Parameter("noFailOnMissingPermissions", strconv.FormatBool(GetBoolOpt(noFailOnMissingPermissionsOpt, flags)))
	processBuilder.Parameter("abortOnError", strconv.FormatBool(GetBoolOpt(abortOnErrorOpt, flags)))
	processBuilder.Parameter("processUserProvidedServices", strconv.FormatBool(GetBoolOpt(processUserProvidedServicesOpt, flags)))
	processBuilder.Parameter("namespace", namespace)
	processBuilder.Parameter("deleteServices", strconv.FormatBool(GetBoolOpt(deleteServicesOpt, flags)))
	processBuilder.Parameter("appsStageTimeout", GetStringOpt(stageTimeoutOpt, flags))
	processBuilder.Parameter("appsUploadTimeout", GetStringOpt(uploadTimeoutOpt, flags))
	processBuilder.Parameter("appsTaskExecutionTimeout", GetStringOpt(taskExecutionTimeoutOpt, flags))
	operation := processBuilder.Build()

	// Create the new process
	responseHeader, err := mtaClient.StartMtaOperation(*operation)
	if err != nil {
		ui.Failed("Could not create rollback mta process: %s", err)
		return Failure
	}

	executionMonitor := NewExecutionMonitorFromLocationHeader(c.name, responseHeader.Location.String(), retries, []*models.Message{}, mtaClient)
	return executionMonitor.Monitor()
}

type rollbackMtaCommandProcessTypeProvider struct{}

func (rollbackMta rollbackMtaCommandProcessTypeProvider) GetProcessType() string {
	return "ROLLBACK_MTA"
}
