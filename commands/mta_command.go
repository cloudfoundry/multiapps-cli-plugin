package commands

import (
	"flag"
	"fmt"
	"strconv"
	"strings"

	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/clients/baseclient"
	mta_models "github.com/cloudfoundry-incubator/multiapps-cli-plugin/clients/models"
	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/ui"
	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/util"
	"github.com/cloudfoundry/cli/cf/formatters"
	"github.com/cloudfoundry/cli/cf/terminal"
	"github.com/cloudfoundry/cli/plugin"
	plugin_models "github.com/cloudfoundry/cli/plugin/models"
)

// MtaCommand is a command for listing a deployed MTA
type MtaCommand struct {
	*BaseCommand
}

func NewMtaCommand() *MtaCommand {
	baseCmd := &BaseCommand{flagsParser: NewDefaultCommandFlagsParser([]string{"MTA_ID"}), flagsValidator: NewDefaultCommandFlagsValidator(nil)}
	mtaCmd := &MtaCommand{baseCmd}
	baseCmd.Command = mtaCmd
	return mtaCmd
}

// GetPluginCommand returns the plugin command details
func (c *MtaCommand) GetPluginCommand() plugin.Command {
	return plugin.Command{
		Name:     "mta",
		HelpText: "Display health and status for a multi-target app",
		UsageDetails: plugin.Usage{
			Usage: "cf mta MTA_ID [--namespace NAMESPACE] [-u URL]",
			Options: map[string]string{
				util.GetShortOption(namespaceOpt): "(EXPERIMENTAL) namespace of the requested mta, empty by default",
				deployServiceURLOpt:               "Deploy service URL, by default 'deploy-service.<system-domain>'",
			},
		},
	}
}

func (c *MtaCommand) defineCommandOptions(flags *flag.FlagSet) {
	flags.String(namespaceOpt, "", "")
}

func (c *MtaCommand) executeInternal(positionalArgs []string, dsHost string, flags *flag.FlagSet, cfTarget util.CloudFoundryTarget) ExecutionStatus {
	mtaID := positionalArgs[0]
	// Print initial message
	ui.Say("Showing health and status for multi-target app %s in org %s / space %s as %s...",
		terminal.EntityNameColor(mtaID), terminal.EntityNameColor(cfTarget.Org.Name),
		terminal.EntityNameColor(cfTarget.Space.Name), terminal.EntityNameColor(cfTarget.Username))

	// Create new REST client
	mtaV2Client := c.NewMtaV2Client(dsHost, cfTarget)

	namespace := strings.TrimSpace(GetStringOpt(namespaceOpt, flags))
	// Get the MTA
	mtas, err := mtaV2Client.GetMtasForThisSpace(&mtaID, &namespace)
	if err != nil {
		ce, ok := err.(*baseclient.ClientError)
		if ok && ce.Code == 404 && strings.Contains(fmt.Sprint(ce.Description), mtaID) {
			ui.Failed("Multi-target app %s not found", terminal.EntityNameColor(mtaID))
			return Failure
		}
		ui.Failed("Could not get multi-target app %s: %s", terminal.EntityNameColor(mtaID), baseclient.NewClientError(err))
		return Failure

	}
	if len(mtas) > 1 {
		ui.Failed("Multiple multi-target apps exist for name %s, please enter namespace", terminal.EntityNameColor(mtaID))
		return Failure
	}
	mta := mtas[0]
	ui.Ok()

	// Display information about all apps and services
	ui.Say("Version: %s", util.GetMtaVersionAsString(mta))
	ui.Say("Namespace: %s", mta.Metadata.Namespace)
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

func isMtaAssociatedApp(mta *mta_models.Mta, app plugin_models.GetAppsModel) bool {
	for _, module := range mta.Modules {
		if module.AppName == app.Name {
			return true
		}
	}
	return false
}

func isMtaAssociatedService(mta *mta_models.Mta, service plugin_models.GetServices_Model) bool {
	for _, serviceName := range mta.Services {
		if serviceName == service.Name {
			return true
		}
	}
	return false
}
