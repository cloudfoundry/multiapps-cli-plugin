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

const processServicesOpt = "process-services"

// RevertDeployCommand is a command for revert deployment of MTAs.
type RevertDeployCommand struct {
	*BaseCommand
	processTypeProvider ProcessTypeProvider
}

// RevertDeployCommand creates a new RevertDeployCommand.
func NewRevertDeployCommand() *RevertDeployCommand {
	baseCmd := &BaseCommand{flagsParser: NewProcessActionExecutorCommandArgumentsParser([]string{"MTA_ID"}), flagsValidator: NewDefaultCommandFlagsValidator(nil)}
	revertDeployCmd := &RevertDeployCommand{baseCmd, &revertDeployCommandProcessTypeProvider{}}
	baseCmd.Command = revertDeployCmd
	return revertDeployCmd
}

// GetPluginCommand returns more information for the blue green deploy command.
func (c *RevertDeployCommand) GetPluginCommand() plugin.Command {
	return plugin.Command{
		Name:     "revert-deploy",
		HelpText: "Revert Deploy of a multi-target app",
		UsageDetails: plugin.Usage{
			Usage: `Revert Deploy of a multi-target app
   cf revert-deploy MTA_ID [-t TIMEOUT] [--version-rule VERSION_RULE] [-f] [--retries RETRIES] [--no-start] [--namespace NAMESPACE] [--no-restart-subscribed-apps] [--no-confirm] [--skip-idle-start] [--do-not-fail-on-missing-permissions] [--abort-on-error] [--apps-start-timeout TIMEOUT] [--apps-stage-timeout TIMEOUT] [--apps-upload-timeout TIMEOUT] [--apps-task-execution-timeout TIMEOUT]

   Perform action on an active deploy operation
   cf revert-deploy -i OPERATION_ID -a ACTION [-u URL]`,
			Options: map[string]string{
				deployServiceURLOpt:               "Deploy service URL, by default 'deploy-service.<system-domain>'",
				versionRuleOpt:                    "Version rule (HIGHER, SAME_HIGHER, ALL)",
				operationIDOpt:                    "Active deploy operation ID",
				actionOpt:                         "Action to perform on active deploy operation (abort, retry, resume, monitor)",
				forceOpt:                          "Force deploy without confirmation for aborting conflicting processes",
				util.GetShortOption(namespaceOpt): "(EXPERIMENTAL) Namespace for the MTA, applied on app names, app routes and service names",
				util.GetShortOption(noFailOnMissingPermissionsOpt): "Do not fail on missing permissions for admin operations",
				util.GetShortOption(abortOnErrorOpt):               "Auto-abort the process on any errors",
				util.GetShortOption(processServicesOpt):            "Enable processing of services during revert",
				util.GetShortOption(retriesOpt):                    "Retry the operation N times in case a non-content error occurs (default 3)",
			},
		},
	}
}

func (c *RevertDeployCommand) defineCommandOptions(flags *flag.FlagSet) {
	flags.Bool(forceOpt, false, "")
	flags.String(operationIDOpt, "", "")
	flags.String(namespaceOpt, "", "")
	flags.String(actionOpt, "", "")
	flags.Bool(noRestartSubscribedAppsOpt, false, "")
	flags.Bool(noFailOnMissingPermissionsOpt, false, "")
	flags.Bool(abortOnErrorOpt, false, "")
	flags.Bool(processServicesOpt, false, "")
	flags.Uint(retriesOpt, 3, "")
}

func (c *RevertDeployCommand) executeInternal(positionalArgs []string, dsHost string, flags *flag.FlagSet, cfTarget util.CloudFoundryTarget) ExecutionStatus {
	operationID := GetStringOpt(operationIDOpt, flags)
	actionID := GetStringOpt(actionOpt, flags)
	retries := GetUintOpt(retriesOpt, flags)

	if operationID != "" || actionID != "" {
		return c.ExecuteAction(operationID, actionID, retries, dsHost, cfTarget)
	}

	force := GetBoolOpt(forceOpt, flags)
	mtaID := positionalArgs[0]
	if !force && !ui.Confirm("Really revert multi-target app %s in org %s / space %s? (y/n)", terminal.EntityNameColor(mtaID),
		terminal.EntityNameColor(cfTarget.Org.Name),
		terminal.EntityNameColor(cfTarget.Space.Name)) {
		ui.Warn("Revert deploy cancelled")
		return Failure
	}

	// Print initial message
	ui.Say("Revert deploying multi-target app %s in org %s / space %s as %s...",
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
	processBuilder.Parameter("shouldProcessServices", strconv.FormatBool(GetBoolOpt(processServicesOpt, flags)))
	processBuilder.Parameter("namespace", namespace)
	operation := processBuilder.Build()

	// Create the new process
	responseHeader, err := mtaClient.StartMtaOperation(*operation)
	if err != nil {
		ui.Failed("Could not create revert deploy process: %s", err)
		return Failure
	}

	executionMonitor := NewExecutionMonitorFromLocationHeader(c.name, responseHeader.Location.String(), retries, []*models.Message{}, mtaClient)
	return executionMonitor.Monitor()
}

type revertDeployCommandProcessTypeProvider struct{}

func (revertDeploy revertDeployCommandProcessTypeProvider) GetProcessType() string {
	return "REVERT_DEPLOY"
}
