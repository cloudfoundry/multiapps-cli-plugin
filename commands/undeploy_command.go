package commands

import (
	"flag"
	"fmt"
	"sort"
	"strconv"
	"strings"

	"code.cloudfoundry.org/cli/cf/terminal"
	"code.cloudfoundry.org/cli/plugin"
	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/clients/baseclient"
	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/clients/models"
	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/ui"
	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/util"
)

// UndeployCommand is a command for undeploying MTAs
type UndeployCommand struct {
	*BaseCommand
	processTypeProvider ProcessTypeProvider
}

func NewUndeployCommand() *UndeployCommand {
	baseCmd := &BaseCommand{flagsParser: NewProcessActionExecutorCommandArgumentsParser([]string{"MTA_ID"}), flagsValidator: NewDefaultCommandFlagsValidator(nil)}
	undeployCmd := &UndeployCommand{baseCmd, &undeployCommandProcessTypeProvider{}}
	baseCmd.Command = undeployCmd
	return undeployCmd
}

// GetPluginCommand returns the plugin command details
func (c *UndeployCommand) GetPluginCommand() plugin.Command {
	return plugin.Command{
		Name:     "undeploy",
		HelpText: "Undeploy a multi-target app",
		UsageDetails: plugin.Usage{
			Usage: `Undeploy a multi-target app
   cf undeploy MTA_ID [-u URL] [-f] [--retries RETRIES] [--namespace NAMESPACE] [--delete-services] [--delete-service-keys] [--delete-service-brokers] [--no-restart-subscribed-apps] [--do-not-fail-on-missing-permissions] [--abort-on-error]

   Perform action on an active undeploy operation
   cf undeploy -i OPERATION_ID -a ACTION [-u URL]` + util.BaseEnvHelpText,
			Options: map[string]string{
				deployServiceURLOpt:                    "Deploy service URL, by default 'deploy-service.<system-domain>'",
				operationIDOpt:                         "Active undeploy operation ID",
				actionOpt:                              "Action to perform on the active undeploy operation (abort, retry, monitor)",
				forceOpt:                               "Force undeploy without confirmation",
				util.GetShortOption(deleteServicesOpt): "Delete services",
				util.GetShortOption(deleteServiceKeysOpt):          "Delete existing service keys",
				util.GetShortOption(deleteServiceBrokersOpt):       "Delete service brokers",
				util.GetShortOption(noRestartSubscribedAppsOpt):    "Do not restart subscribed apps, updated during the undeployment",
				util.GetShortOption(noFailOnMissingPermissionsOpt): "Do not fail on missing permissions for admin operations",
				util.GetShortOption(abortOnErrorOpt):               "Auto-abort the process on any errors",
				util.GetShortOption(retriesOpt):                    "Retry the operation N times in case a non-content error occurs (default 3)",
				util.GetShortOption(namespaceOpt):                  "Specify the (optional) namespace the target mta is in",
			},
		},
	}
}

func (c *UndeployCommand) defineCommandOptions(flags *flag.FlagSet) {
	flags.Bool(forceOpt, false, "")
	flags.String(operationIDOpt, "", "")
	flags.String(namespaceOpt, "", "")
	flags.String(actionOpt, "", "")
	flags.Bool(deleteServicesOpt, false, "")
	flags.Bool(deleteServiceKeysOpt, false, "")
	flags.Bool(noRestartSubscribedAppsOpt, false, "")
	flags.Bool(deleteServiceBrokersOpt, false, "")
	flags.Bool(noFailOnMissingPermissionsOpt, false, "")
	flags.Bool(abortOnErrorOpt, false, "")
	flags.Uint(retriesOpt, 3, "")
}

func (c *UndeployCommand) executeInternal(positionalArgs []string, dsHost string, flags *flag.FlagSet, cfTarget util.CloudFoundryTarget) ExecutionStatus {
	operationID := GetStringOpt(operationIDOpt, flags)
	actionID := GetStringOpt(actionOpt, flags)
	retries := GetUintOpt(retriesOpt, flags)

	if operationID != "" || actionID != "" {
		return c.ExecuteAction(operationID, actionID, retries, dsHost, cfTarget)
	}

	force := GetBoolOpt(forceOpt, flags)
	mtaID := positionalArgs[0]
	if !force && !ui.Confirm("Really undeploy multi-target app %s in org %s / space %s? (y/n)", terminal.EntityNameColor(mtaID),
		terminal.EntityNameColor(cfTarget.Org.Name),
		terminal.EntityNameColor(cfTarget.Space.Name)) {
		ui.Warn("Undeploy cancelled")
		return Failure
	}

	// Print initial message
	ui.Say("Undeploying multi-target app %s in org %s / space %s as %s...",
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
	processBuilder.Parameter("noRestartSubscribedApps", strconv.FormatBool(GetBoolOpt(noRestartSubscribedAppsOpt, flags)))
	processBuilder.Parameter("deleteServices", strconv.FormatBool(GetBoolOpt(deleteServicesOpt, flags)))
	processBuilder.Parameter("deleteServiceKeys", strconv.FormatBool(GetBoolOpt(deleteServiceKeysOpt, flags)))
	processBuilder.Parameter("deleteServiceBrokers", strconv.FormatBool(GetBoolOpt(deleteServiceBrokersOpt, flags)))
	processBuilder.Parameter("noFailOnMissingPermissions", strconv.FormatBool(GetBoolOpt(noFailOnMissingPermissionsOpt, flags)))
	processBuilder.Parameter("abortOnError", strconv.FormatBool(GetBoolOpt(abortOnErrorOpt, flags)))
	processBuilder.Parameter("namespace", namespace)
	operation := processBuilder.Build()

	// Create the new process
	responseHeader, err := mtaClient.StartMtaOperation(*operation)
	if err != nil {
		ui.Failed("Could not create undeploy process: %s", err)
		return Failure
	}

	executionMonitor := NewExecutionMonitorFromLocationHeader(c.name, responseHeader.Location.String(), retries, []*models.Message{}, mtaClient)
	return executionMonitor.Monitor()
}

type undeployCommandProcessTypeProvider struct{}

func (d undeployCommandProcessTypeProvider) GetProcessType() string {
	return "UNDEPLOY"
}

type ProcessActionExecutorCommandArgumentsParser struct {
	positionalArgNames []string
}

func NewProcessActionExecutorCommandArgumentsParser(positionalArgNames []string) ProcessActionExecutorCommandArgumentsParser {
	return ProcessActionExecutorCommandArgumentsParser{positionalArgNames: positionalArgNames}
}

func (p ProcessActionExecutorCommandArgumentsParser) ParseFlags(flags *flag.FlagSet, args []string) error {
	operationExecutorOptions := make(map[string]string)
	for _, arg := range args {
		optionFlag := flags.Lookup(strings.Replace(arg, "-", "", 1))
		if optionFlag != nil && (operationIDOpt == optionFlag.Name || actionOpt == optionFlag.Name) {
			operationExecutorOptions[optionFlag.Name] = arg
		}
	}

	if len(operationExecutorOptions) > 2 {
		return fmt.Errorf("Options %s and %s should be specified only once", operationIDOpt, actionOpt)
	}

	if len(operationExecutorOptions) == 1 {
		keys := append([]string{}, []string{operationIDOpt, actionOpt}...)
		sort.Strings(keys)
		return fmt.Errorf("All the %s options should be specified together", strings.Join(keys, " "))
	}

	return NewDefaultCommandFlagsParser(p.determinePositionalArguments(operationExecutorOptions)).ParseFlags(flags, args)
}

func (p ProcessActionExecutorCommandArgumentsParser) determinePositionalArguments(operationExecutorOptions map[string]string) []string {
	if len(operationExecutorOptions) == 2 {
		return []string{}
	}
	return p.positionalArgNames
}
