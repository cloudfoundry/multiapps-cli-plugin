package commands

import (
	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/clients/baseclient"
	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/log"
	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/ui"
	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/util"
	"github.com/cloudfoundry/cli/cf/terminal"
	"github.com/cloudfoundry/cli/plugin"
)

// MtasCommand is a command for listing all deployed MTAs
type MtasCommand struct {
	BaseCommand
}

func NewMtasCommand() *MtasCommand {
	return &MtasCommand{BaseCommand{options: getMtasCommandOptions()}}
}

func getMtasCommandOptions() map[string]CommandOption {
	return map[string]CommandOption{
		deployServiceURLOpt: deployServiceUrlOption(),
	}
}

// GetPluginCommand returns the plugin command details
func (c *MtasCommand) GetPluginCommand() plugin.Command {
	return plugin.Command{
		Name:     "mtas",
		HelpText: "List all multi-target apps",
		UsageDetails: plugin.Usage{
			Usage:   "cf mtas [-u URL]",
			Options: c.getOptionsForPluginCommand(),
		},
	}
}

// Execute executes the command
func (c *MtasCommand) Execute(args []string) ExecutionStatus {
	log.Tracef("Executing command '" + c.name + "': args: '%v'\n", args)

	parser := NewCommandFlagsParser(c.flags, NewDefaultCommandFlagsParser(0))
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

	// Print initial message
	ui.Say("Getting multi-target apps in org %s / space %s as %s...",
		terminal.EntityNameColor(context.Org), terminal.EntityNameColor(context.Space), terminal.EntityNameColor(context.Username))

	// Create new REST client
	mtaClient, err := c.NewMtaClient(host)
	if err != nil {
		ui.Failed("Could not get space id: %s", baseclient.NewClientError(err))
		return Failure
	}

	// Get all deployed components
	mtas, err := mtaClient.GetMtas()
	if err != nil {
		ui.Failed("Could not get deployed components: %s", baseclient.NewClientError(err))
		return Failure
	}
	ui.Ok()

	// Print all deployed MTAs
	if len(mtas) > 0 {
		table := ui.Table([]string{"mta id", "version"})
		for _, mta := range mtas {
			table.Add(mta.Metadata.ID, util.GetMtaVersionAsString(mta))
		}
		table.Print()
	} else {
		ui.Say("No multi-target apps found")
	}
	return Success
}
