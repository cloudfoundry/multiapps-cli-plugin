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
	flags, err := c.CreateFlags(&host, args)
	if err != nil {
		ui.Failed(err.Error())
		return Failure
	}
	parser := NewCommandFlagsParser(flags, NewDefaultCommandFlagsParser(nil), NewDefaultCommandFlagsValidator(nil))
	err = parser.Parse(args)
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
	mtaV2Client, err := c.NewMtaV2Client(host)
	if err != nil {
		ui.Failed("Could not get space ID: %s", baseclient.NewClientError(err))
		return Failure
	}

	// Get all deployed components
	mtas, err := mtaV2Client.GetMtasForThisSpace(nil, nil)
	if err != nil {
		ui.Failed("Could not get deployed components: %s", baseclient.NewClientError(err))
		return Failure
	}
	ui.Ok()

	// Print all deployed MTAs
	if len(mtas) > 0 {
		table := ui.Table([]string{"mta id", "version", "namespace"})
		for _, mta := range mtas {
			table.Add(mta.Metadata.ID, util.GetMtaVersionAsString(mta), mta.Metadata.Namespace)
		}
		table.Print()
	} else {
		ui.Say("No multi-target apps found")
	}
	return Success
}
