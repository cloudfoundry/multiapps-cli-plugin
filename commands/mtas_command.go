package commands

import (
	"github.com/cloudfoundry/cli/cf/terminal"
	"github.com/cloudfoundry/cli/plugin"
	"github.com/SAP/cf-mta-plugin/log"
	"github.com/SAP/cf-mta-plugin/ui"
	"github.com/SAP/cf-mta-plugin/util"
)

// MtasCommand is a command for listing all deployed MTAs
type MtasCommand struct {
	BaseCommand
}

// GetPluginCommand returns the plugin command details
func (c *MtasCommand) GetPluginCommand() plugin.Command {
	return plugin.Command{
		Name:     "mtas",
		HelpText: "List all multi-target apps",
		UsageDetails: plugin.Usage{
			Usage: "cf mtas [-u URL]",
			Options: map[string]string{
				"u": "Deploy service URL, by default 'deploy-service.<system-domain>'",
			},
		},
	}
}

// Execute executes the command
func (c *MtasCommand) Execute(args []string) ExecutionStatus {
	log.Tracef("Executing command '"+c.name+"': args: '%v'\n", args)

	var host string

	// Parse command arguments and check for required options
	flags, err := c.CreateFlags(&host)
	if err != nil {
		ui.Failed(err.Error())
		return Failure
	}
	err = c.ParseFlags(args, nil, flags, nil)
	if err != nil {
		c.Usage(err.Error())
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
	restClient, err := c.NewRestClient(host)
	if err != nil {
		ui.Failed(err.Error())
		return Failure
	}

	// Get all deployed components
	components, err := restClient.GetComponents()
	if err != nil {
		ui.Failed("Could not get deployed components: %s", err)
		return Failure
	}
	ui.Ok()

	// Print all deployed MTAs
	mtas := components.Mtas.Mtas
	if len(mtas) > 0 {
		table := ui.Table([]string{"mta id", "version"})
		for _, mta := range mtas {
			table.Add(*mta.Metadata.ID, util.GetMtaVersionAsString(mta))
		}
		table.Print()
	} else {
		ui.Say("No multi-target apps found")
	}
	return Success
}
