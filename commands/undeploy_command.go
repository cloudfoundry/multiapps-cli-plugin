package commands

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/clients/baseclient"
	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/log"
	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/ui"
	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/util"
	"github.com/cloudfoundry/cli/cf/terminal"
	"github.com/cloudfoundry/cli/plugin"
)

//UndeployCommand is a command for undeploying MTAs
type UndeployCommand struct {
	BaseCommand
	processParametersSetter ProcessParametersSetter
	processTypeProvider     ProcessTypeProvider
}

func NewUndeployCommand() *UndeployCommand {
	return &UndeployCommand{BaseCommand{options: getUndeployCommandOptions()}, undeployProcessParametersSetter(), &undeployCommandProcessTypeProvider{}}
}

func getUndeployCommandOptions() map[string]CommandOption {
	return map[string]CommandOption{
		deployServiceURLOpt:           deployServiceUrlOption(),
		operationIDOpt:                {new(string), "", "Active undeploy operation id", true},
		actionOpt:                     {new(string), "", "Action to perform on active undeploy operation (abort, retry, monitor)", true},
		forceOpt:                      {new(bool), false, "Force undeploy without confirmation for aborting conflicting processes", true},
		deleteServicesOpt:             {new(bool), false, "Delete services", false},
		deleteServiceBrokersOpt:       {new(bool), false, "Delete service brokers", false},
		noRestartSubscribedAppsOpt:    {new(bool), false, "Do not restart subscribed apps, updated during the undeployment", false},
		noFailOnMissingPermissionsOpt: {new(bool), false, "Do not fail on missing permissions for admin operations", false},
		abortOnErrorOpt:               {new(bool), false, "Auto-abort the process on any errors", false},
		retriesOpt:                    {new(uint), 3, "Retry the operation N times in case a non-content error occurs (default 3)", false},
	}
}

// GetPluginCommand returns the plugin command details
func (c *UndeployCommand) GetPluginCommand() plugin.Command {
	return plugin.Command{
		Name:     "undeploy",
		HelpText: "Undeploy a multi-target app",
		UsageDetails: plugin.Usage{
			Usage: `Undeploy a multi-target app
   cf undeploy MTA_ID [-u URL] [-f] [--retries RETRIES] [--delete-services] [--delete-service-brokers] [--no-restart-subscribed-apps] [--do-not-fail-on-missing-permissions] [--abort-on-error]

   Perform action on an active undeploy operation
   cf undeploy -i OPERATION_ID -a ACTION [-u URL]`,
			Options: c.getOptionsForPluginCommand(),
		},
	}
}

func undeployProcessParametersSetter() ProcessParametersSetter {
	return func(options map[string]CommandOption, processBuilder *util.ProcessBuilder) {
		processBuilder.Parameter("noRestartSubscribedApps", strconv.FormatBool(getBoolOpt(noRestartSubscribedAppsOpt, options)))
		processBuilder.Parameter("deleteServices", strconv.FormatBool(getBoolOpt(deleteServicesOpt, options)))
		processBuilder.Parameter("deleteServiceBrokers", strconv.FormatBool(getBoolOpt(deleteServiceBrokersOpt, options)))
		processBuilder.Parameter("noFailOnMissingPermissions", strconv.FormatBool(getBoolOpt(noFailOnMissingPermissionsOpt, options)))
		processBuilder.Parameter("abortOnError", strconv.FormatBool(getBoolOpt(abortOnErrorOpt, options)))
	}
}

// Execute executes the command
func (c *UndeployCommand) Execute(args []string) ExecutionStatus {
	log.Tracef("Executing command '" + c.name + "': args: '%v'\n", args)

	parser := NewCommandFlagsParserWithValidator(c.flags, NewProcessActionExecutorCommandArgumentsParser(1), NewPositionalArgumentsFlagsValidator([]string{"MTA_ID"}))
	err := parser.Parse(args)
	if err != nil {
		c.Usage(err.Error())
		return Failure
	}

	host, err := c.computeDeployServiceUrl()
	if err != nil {
		ui.Failed("Could not compute deploy service URL: %s", err.Error())
		return Failure
	}

	context, err := c.GetContext()
	if err != nil {
		ui.Failed(err.Error())
		return Failure
	}

	operationID := getStringOpt(operationIDOpt, c.options)
	actionID := getStringOpt(actionOpt, c.options)
	retries := getUintOpt(retriesOpt, c.options)

	if operationID != "" || actionID != "" {
		return c.ExecuteAction(operationID, actionID, retries, host)
	}

	force := getBoolOpt(forceOpt, c.options)
	mtaID := args[0]

	if !force && !ui.Confirm("Really undeploy multi-target app %s? (y/n)", terminal.EntityNameColor(mtaID)) {
		ui.Warn("Undeploy cancelled")
		return Failure
	}

	// Print initial message
	ui.Say("Undeploying multi-target app %s in org %s / space %s as %s...",
		terminal.EntityNameColor(mtaID), terminal.EntityNameColor(context.Org),
		terminal.EntityNameColor(context.Space), terminal.EntityNameColor(context.Username))

	// Create rest client
	mtaClient, err := c.NewMtaClient(host)
	if err != nil {
		ui.Failed(err.Error())
		return Failure
	}

	// Check if a deployed MTA with the specified ID exists
	_, err = mtaClient.GetMta(mtaID)
	if err != nil {
		ce, ok := err.(*baseclient.ClientError)
		if ok && ce.Code == 404 && strings.Contains(fmt.Sprint(ce.Description), mtaID) {
			ui.Failed("Multi-target app %s not found", terminal.EntityNameColor(mtaID))
			return Failure
		}
		ui.Failed("Could not get multi-target app %s: %s", terminal.EntityNameColor(mtaID), baseclient.NewClientError(err))
		return Failure
	}

	// Check for an ongoing operation for this MTA ID and abort it
	wasAborted, err := c.CheckOngoingOperation(mtaID, force, host)
	if err != nil {
		ui.Failed(err.Error())
		return Failure
	}
	if !wasAborted {
		return Failure
	}

	processBuilder := util.NewProcessBuilder()
	processBuilder.ProcessType(c.processTypeProvider.GetProcessType())
	c.processParametersSetter(c.options, processBuilder)
	processBuilder.Parameter("mtaId", mtaID)

	operation := processBuilder.Build()

	// Create the new process
	responseHeader, err := mtaClient.StartMtaOperation(*operation)
	if err != nil {
		ui.Failed("Could not create undeploy process: %s", err)
		return Failure
	}

	// Monitor process execution
	return NewExecutionMonitorFromLocationHeader(c.name, responseHeader.Location.String(), retries, mtaClient).Monitor()
}

type undeployCommandProcessTypeProvider struct{}

func (d undeployCommandProcessTypeProvider) GetProcessType() string {
	return "UNDEPLOY"
}
