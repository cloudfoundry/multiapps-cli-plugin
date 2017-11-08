package commands

import (
	"fmt"
	"strconv"
	"strings"

	baseclient "github.com/SAP/cf-mta-plugin/clients/baseclient"
	"github.com/SAP/cf-mta-plugin/log"
	"github.com/SAP/cf-mta-plugin/ui"
	"github.com/SAP/cf-mta-plugin/util"
	"github.com/cloudfoundry/cli/cf/formatters"
	"github.com/cloudfoundry/cli/cf/terminal"
	"github.com/cloudfoundry/cli/plugin"
	"github.com/cloudfoundry/cli/plugin/models"
)

// MtaCommand is a command for listing a deployed MTA
type MtaCommand struct {
	BaseCommand
}

// GetPluginCommand returns the plugin command details
func (c *MtaCommand) GetPluginCommand() plugin.Command {
	return plugin.Command{
		Name:     "mta",
		HelpText: "Display health and status for a multi-target app",
		UsageDetails: plugin.Usage{
			Usage: "cf mta MTA_ID [-u URL]",
			Options: map[string]string{
				"u": "Deploy service URL, by default 'deploy-service.<system-domain>'",
			},
		},
	}
}

// Execute executes the command
func (c *MtaCommand) Execute(args []string) ExecutionStatus {
	log.Tracef("Executing command '"+c.name+"': args: '%v'\n", args)

	var host string

	// Parse command arguments and check for required options
	flags, err := c.CreateFlags(&host)
	if err != nil {
		ui.Failed(err.Error())
		return Failure
	}
	err = c.ParseFlags(args, []string{"MTA_ID"}, flags, nil)
	if err != nil {
		c.Usage(err.Error())
		return Failure
	}

	mtaID := args[0]

	context, err := c.GetContext()
	if err != nil {
		ui.Failed("Could not get org and space: %s", baseclient.NewClientError(err))
		return Failure
	}

	// Print initial message
	ui.Say("Showing health and status for multi-target app %s in org %s / space %s as %s...",
		terminal.EntityNameColor(mtaID), terminal.EntityNameColor(context.Org),
		terminal.EntityNameColor(context.Space), terminal.EntityNameColor(context.Username))

	// Create new REST client
	mtaClient, err := c.NewMtaClient(host)
	if err != nil {
		ui.Failed("Could not get space id: %s", baseclient.NewClientError(err))
		return Failure
	}

	// Get the MTA
	mta, err := mtaClient.GetMta(mtaID)
	if err != nil {
		ce, ok := err.(*baseclient.ClientError)
		if ok && ce.Code == 404 && strings.Contains(fmt.Sprint(ce.Description), mtaID) {
			ui.Failed("Multi-target app %s not found", terminal.EntityNameColor(mtaID))
			return Failure
		}
		ui.Failed("Could not get multi-target app %s: %s", terminal.EntityNameColor(mtaID), baseclient.NewClientError(err))
		return Failure

	}
	ui.Ok()

	// Display information about all apps and services
	ui.Say("Version: %s", util.GetMtaVersionAsString(mta))
	ui.Say("\nApps:")
	table := ui.Table([]string{"name", "requested state", "instances", "memory", "disk", "urls"})
	serviceApps := make(map[string][]string)
	for _, mtaModule := range mta.Modules {
		app, err := c.cliConnection.GetApp(mtaModule.AppName)
		if err != nil {
			ui.Failed("Could not get app %s: %s", terminal.EntityNameColor(mtaModule.AppName), baseclient.NewClientError(err))
			return Failure
		}
		table.Add(app.Name, app.State, getInstances(app), size(app.Memory), size(app.DiskQuota), getRoutes(app))
		for _, service := range app.Services {
			serviceApps[service.Name] = append(serviceApps[service.Name], app.Name)
		}
	}
	table.Print()
	if len(mta.Services) > 0 {
		ui.Say("\nServices:")
		table := ui.Table([]string{"name", "service", "plan", "bound apps", "last operation"})
		for _, serviceName := range mta.Services {
			service, err := c.cliConnection.GetService(serviceName)
			if err != nil {
				ui.Failed("Could not get service %s: %s", terminal.EntityNameColor(serviceName), baseclient.NewClientError(err))
				return Failure
			}
			table.Add(service.Name, service.ServiceOffering.Name, service.ServicePlan.Name,
				getBoundApps(service, serviceApps), getLastOperation(service))
		}
		table.Print()
	}
	return Success
}

func size(n int64) string {
	return formatters.ByteSize(n * formatters.MEGABYTE)
}

func getInstances(app plugin_models.GetAppModel) string {
	return strconv.Itoa(app.RunningInstances) + "/" + strconv.Itoa(app.InstanceCount)
}

func getRoutes(app plugin_models.GetAppModel) string {
	var urls []string
	for _, route := range app.Routes {
		urls = append(urls, route.Host+"."+route.Domain.Name)
	}
	return strings.Join(urls, ", ")
}

func getBoundApps(service plugin_models.GetService_Model, serviceApps map[string][]string) string {
	return strings.Join(serviceApps[service.Name], ", ")
}

func getLastOperation(service plugin_models.GetService_Model) string {
	return service.LastOperation.Type + " " + service.LastOperation.State
}
