package commands

import (
	"fmt"
	"strconv"
	"strings"

	baseclient "github.com/cloudfoundry-incubator/multiapps-cli-plugin/clients/baseclient"
	mtamodels "github.com/cloudfoundry-incubator/multiapps-cli-plugin/clients/models"
	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/log"
	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/ui"
	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/util"
	"github.com/cloudfoundry/cli/cf/formatters"
	"github.com/cloudfoundry/cli/cf/terminal"
	"github.com/cloudfoundry/cli/plugin"
	plugin_models "github.com/cloudfoundry/cli/plugin/models"
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
	flags, err := c.CreateFlags(&host, args)
	if err != nil {
		ui.Failed(err.Error())
		return Failure
	}

	parser := NewCommandFlagsParser(flags, NewDefaultCommandFlagsParser([]string{"MTA_ID"}), NewDefaultCommandFlagsValidator(map[string]bool{}))
	err = parser.Parse(args)
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
		ui.Failed("Could not get space ID: %s", baseclient.NewClientError(err))
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
	// GetApps() is more safe than GetApp(), because it retrieves all application statistics through a single call,
	// whereas GetApp(name) makes several calls, some of which may fail under specific conditions. For example,
	// the call to /v2/apps/<GUID>/instances may fail if the application is not yet staged. Weirdly enough, a call
	// to GetApps() is also faster than several calls to GetApp(name).
	apps, err := c.cliConnection.GetApps()
	if err != nil {
		ui.Failed("Could not get apps: %s", baseclient.NewClientError(err))
		return Failure
	}
	for _, app := range apps {
		if isMtaAssociatedApp(mta, app) {
			table.Add(app.Name, app.State, getInstances(app), size(app.Memory), size(app.DiskQuota), getRoutes(app))
		}
	}
	table.Print()

	if len(mta.Services) == 0 {
		return Success
	}
	ui.Say("\nServices:")
	table = ui.Table([]string{"name", "service", "plan", "bound apps", "last operation"})
	// Read the comment for GetApps(). The same applies for GetServices().
	services, err := c.cliConnection.GetServices()
	if err != nil {
		ui.Failed("Could not get services: %s", baseclient.NewClientError(err))
		return Failure
	}
	for _, service := range services {
		if isMtaAssociatedService(mta, service) {
			table.Add(service.Name, service.Service.Name, service.ServicePlan.Name, strings.Join(service.ApplicationNames, ", "), getLastOperation(service))
		}
	}
	table.Print()

	return Success
}

func size(n int64) string {
	return formatters.ByteSize(n * formatters.MEGABYTE)
}

func getInstances(app plugin_models.GetAppsModel) string {
	return strconv.Itoa(app.RunningInstances) + "/" + strconv.Itoa(app.TotalInstances)
}

func getRoutes(app plugin_models.GetAppsModel) string {
	var urls []string
	for _, route := range app.Routes {
		urls = append(urls, route.Host+"."+route.Domain.Name)
	}
	return strings.Join(urls, ", ")
}

func getLastOperation(service plugin_models.GetServices_Model) string {
	return service.LastOperation.Type + " " + service.LastOperation.State
}

func isMtaAssociatedApp(mta *mtamodels.Mta, app plugin_models.GetAppsModel) bool {
	for _, module := range mta.Modules {
		if module.AppName == app.Name {
			return true
		}
	}
	return false
}

func isMtaAssociatedService(mta *mtamodels.Mta, service plugin_models.GetServices_Model) bool {
	for _, serviceName := range mta.Services {
		if serviceName == service.Name {
			return true
		}
	}
	return false
}
