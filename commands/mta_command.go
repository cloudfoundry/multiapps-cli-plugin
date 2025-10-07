package commands

import (
	"flag"
	"fmt"
	"strconv"
	"strings"

	"code.cloudfoundry.org/cli/cf/formatters"
	"code.cloudfoundry.org/cli/cf/terminal"
	"code.cloudfoundry.org/cli/plugin"
	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/clients/baseclient"
	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/clients/cfrestclient"
	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/clients/cfrestclient/resilient"
	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/clients/models"
	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/ui"
	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/util"
)

// MtaCommand is a command for listing a deployed MTA
type MtaCommand struct {
	*BaseCommand

	CfClient cfrestclient.CloudFoundryOperationsExtended
}

func NewMtaCommand() *MtaCommand {
	baseCmd := &BaseCommand{flagsParser: NewDefaultCommandFlagsParser([]string{"MTA_ID"}), flagsValidator: NewDefaultCommandFlagsValidator(nil)}
	mtaCmd := &MtaCommand{BaseCommand: baseCmd}
	baseCmd.Command = mtaCmd
	return mtaCmd
}

func (c *MtaCommand) Initialize(name string, cliConnection plugin.CliConnection) {
	c.BaseCommand.Initialize(name, cliConnection)
	delegate := cfrestclient.NewCloudFoundryRestClient(cliConnection)
	c.CfClient = resilient.NewResilientCloudFoundryClient(delegate, maxRetriesCount, retryIntervalInSeconds)
}

// GetPluginCommand returns the plugin command details
func (c *MtaCommand) GetPluginCommand() plugin.Command {
	return plugin.Command{
		Name:     "mta",
		HelpText: "Display health and status for a multi-target app",
		UsageDetails: plugin.Usage{
			Usage: "cf mta MTA_ID [--namespace NAMESPACE] [-u URL]" + util.BaseEnvHelpText,
			Options: map[string]string{
				util.GetShortOption(namespaceOpt): "namespace of the requested mta, empty by default",
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

	apps, err := c.CfClient.GetApplications(mta.Metadata.ID, mta.Metadata.Namespace, cfTarget.Space.Guid)
	if err != nil {
		ui.Failed("Could not get apps: %s", err)
		return Failure
	}

	table := ui.Table([]string{"name", "requested state", "instances", "memory", "disk", "urls"})

	for _, app := range apps {
		processes, err := c.CfClient.GetAppProcessStatistics(app.Guid)
		if err != nil {
			ui.Failed("Could not get app %q process statistics: %s", app.Name, err)
			return Failure
		}

		routes, err := c.CfClient.GetApplicationRoutes(app.Guid)
		if err != nil {
			ui.Failed("Could not get app %q routes: %s", app.Name, err)
			return Failure
		}

		memory := int64(0)
		disk := int64(0)
		if len(processes) > 0 {
			memory = processes[0].Memory
			disk = processes[0].Disk
		}
		table.Add(app.Name, app.State, getInstances(processes), size(memory), size(disk), formatRoutes(routes))
	}
	table.Print()

	ui.Say("\nServices:")
	if len(mta.Services) == 0 {
		return Success
	}

	services, err := c.CfClient.GetServiceInstances(mta.Metadata.ID, mta.Metadata.Namespace, cfTarget.Space.Guid)
	if err != nil {
		ui.Failed("Could not get services: %s", err)
		return Failure
	}

	table = ui.Table([]string{"name", "service", "plan", "bound apps", "last operation"})

	for _, service := range services {
		serviceBindings, err := c.CfClient.GetServiceBindings(service.Name)
		if err != nil {
			ui.Failed("Could not get service bindings: %s", err)
			return Failure
		}

		table.Add(service.Name, service.Offering.Name, service.Plan.Name, formatAppNames(serviceBindings), getLastOperation(service))
	}
	table.Print()

	return Success
}

func size(n int64) string {
	return formatters.ByteSize(n)
}

func getInstances(processes []models.ApplicationProcessStatistics) string {
	var runningProcesses []models.ApplicationProcessStatistics
	for _, process := range processes {
		if process.State == "RUNNING" {
			runningProcesses = append(runningProcesses, process)
		}
	}
	return strconv.Itoa(len(runningProcesses)) + "/" + strconv.Itoa(len(processes))
}

func formatRoutes(routes []models.ApplicationRoute) string {
	var urls []string
	for _, route := range routes {
		urls = append(urls, route.Url)
	}
	return strings.Join(urls, ", ")
}

func formatAppNames(bindings []models.ServiceBinding) string {
	var appNames []string
	for _, binding := range bindings {
		appNames = append(appNames, binding.AppName)
	}
	return strings.Join(appNames, ", ")
}

func getLastOperation(service models.CloudFoundryServiceInstance) string {
	return service.LastOperation.Type + " " + service.LastOperation.State
}
