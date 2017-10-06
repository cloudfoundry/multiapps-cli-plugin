package commands

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/cloudfoundry/cli/cf/terminal"
	"github.com/cloudfoundry/cli/plugin"
	baseclient "github.com/SAP/cf-mta-plugin/clients/baseclient"
	"github.com/SAP/cf-mta-plugin/clients/models"
	"github.com/SAP/cf-mta-plugin/log"
	"github.com/SAP/cf-mta-plugin/ui"
	"github.com/SAP/cf-mta-plugin/util"
)

//UndeployCommand is a command for undeploying MTAs
type UndeployCommand struct {
	BaseCommand
}

// GetPluginCommand returns the plugin command details
func (c *UndeployCommand) GetPluginCommand() plugin.Command {
	return plugin.Command{
		Name:     "undeploy",
		HelpText: "Undeploy a multi-target app",
		UsageDetails: plugin.Usage{
			Usage: `Undeploy a multi-target app
   cf undeploy MTA_ID [-u URL] [-f] [--delete-services] [--delete-service-brokers] [--no-restart-subscribed-apps] [--do-not-fail-on-missing-permissions]

   Perform action on an active undeploy operation
   cf undeploy -i OPERATION_ID -a ACTION [-u URL]`,
			Options: map[string]string{
				deployServiceURLOpt: "Deploy service URL, by default 'deploy-service.<system-domain>'",
				operationIDOpt:      "Active undeploy operation id",
				actionOpt:           "Action to perform on the active undeploy operation (abort, retry, monitor)",
				forceOpt:            "Force undeploy without confirmation",
				util.GetShortOption(deleteServicesOpt):             "Delete services",
				util.GetShortOption(deleteServiceBrokersOpt):       "Delete service brokers",
				util.GetShortOption(noRestartSubscribedAppsOpt):    "Do not restart subscribed apps, updated during the undeployment",
				util.GetShortOption(noFailOnMissingPermissionsOpt): "Do not fail on missing permissions for admin operations",
			},
		},
	}
}

// ServiceID returns the service ID of the processes started by UndeployCommand
func (c *UndeployCommand) ServiceID() ServiceID {
	return UndeployServiceID
}

// Execute executes the command
func (c *UndeployCommand) Execute(args []string) ExecutionStatus {
	log.Tracef("Executing command '"+c.name+"': args: '%v'\n", args)

	var serviceID = c.ServiceID()

	var host string
	var operationID string
	var actionID string
	var force bool
	var deleteServices bool
	var noRestartSubscribedApps bool
	var deleteServiceBrokers bool
	var noFailOnMissingPermissions bool
	flags, err := c.CreateFlags(&host)
	if err != nil {
		ui.Failed(err.Error())
		return Failure
	}
	flags.BoolVar(&force, forceOpt, false, "")
	flags.StringVar(&operationID, operationIDOpt, "", "")
	flags.StringVar(&actionID, actionOpt, "", "")
	flags.BoolVar(&deleteServices, deleteServicesOpt, false, "")
	flags.BoolVar(&noRestartSubscribedApps, noRestartSubscribedAppsOpt, false, "")
	flags.BoolVar(&deleteServiceBrokers, deleteServiceBrokersOpt, false, "")
	flags.BoolVar(&noFailOnMissingPermissions, noFailOnMissingPermissionsOpt, false, "")
	shouldExecuteActionOnExistingProcess, _ := ContainsSpecificOptions(flags, args, map[string]string{"i": "-i", "a": "-a"})
	var positionalArgNames []string
	if !shouldExecuteActionOnExistingProcess {
		positionalArgNames = []string{"MTA_ID"}
	}
	err = c.ParseFlags(args, positionalArgNames, flags, nil)
	if err != nil {
		c.Usage(err.Error())
		return Failure
	}

	context, err := c.GetContext()
	if err != nil {
		ui.Failed(err.Error())
		return Failure
	}

	if operationID != "" || actionID != "" {
		return c.ExecuteAction(operationID, actionID, host, serviceID)
	}

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
	restClient, err := c.NewRestClient(host)
	if err != nil {
		ui.Failed(err.Error())
		return Failure
	}

	// Check if a deployed MTA with the specified ID exists
	_, err = restClient.GetMta(mtaID)
	if err != nil {
		ce, ok := err.(*baseclient.ClientError)
		// TODO(ivan): This is crap. Expecting the error message returned by a
		// remote component to decide what to print ot the user based on a
		// simple Contains method is not sane.
		if ok && ce.Code == 404 && strings.Contains(fmt.Sprint(ce.Description), mtaID) {
			ui.Failed("Multi-target app %s not found", terminal.EntityNameColor(mtaID))
			return Failure
		}
		ui.Failed("Could not get multi-target app %s: %s", terminal.EntityNameColor(mtaID), err)
		return Failure

	}

	// Check for an ongoing operation for this MTA ID and abort it
	wasAborted, err := c.CheckOngoingOperation(mtaID, host, force)
	if err != nil {
		ui.Failed(err.Error())
		return Failure
	}
	if !wasAborted {
		return Failure
	}

	// Create SLMP client
	slmpClient, err := c.NewSlmpClient(host)
	if err != nil {
		ui.Failed(err.Error())
		return Failure
	}

	// Check SLMP metadata
	err = CheckSlmpMetadata(slmpClient)
	if err != nil {
		ui.Failed(err.Error())
		return Failure
	}

	ui.Say("Starting undeployment process...")

	processBuilder := util.NewProcessBuilder()
	processBuilder.ServiceID(serviceID.String())
	processBuilder.Parameter("mtaId", mtaID)
	processBuilder.Parameter("noRestartSubscribedApps", strconv.FormatBool(noRestartSubscribedApps))
	processBuilder.Parameter("deleteServices", strconv.FormatBool(deleteServices))
	processBuilder.Parameter("deleteServiceBrokers", strconv.FormatBool(deleteServiceBrokers))
	processBuilder.Parameter("noFailOnMissingPermissions", strconv.FormatBool(noFailOnMissingPermissions))
	process := processBuilder.Build()

	// Create the new process
	createdProcess, err := slmpClient.CreateServiceProcess(serviceID.String(), process)
	if err != nil {
		ui.Failed("Could not create process for service %s: %s", terminal.EntityNameColor(serviceID.String()), err)
		return Failure
	}
	ui.Ok()

	processID := createdProcess.ID
	// Create SLPP client
	slppClient, err := c.NewSlppClient(host, serviceID.String(), processID)
	if err != nil {
		ui.Failed(err.Error())
		return Failure
	}

	// Check SLPP metadata
	err = CheckSlppMetadata(slppClient)
	if err != nil {
		ui.Failed(err.Error())
		return Failure
	}

	// Monitor process execution
	monitor := NewExecutionMonitor(processID, c.name, slppClient, []*models.ProgressMessage{})
	return monitor.Monitor()
}
