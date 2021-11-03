package commands

import (
	"code.cloudfoundry.org/cli/cf/terminal"
	"code.cloudfoundry.org/cli/plugin"
	"flag"
	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/clients/baseclient"
	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/ui"
	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/util"
)

// MtasCommand is a command for listing all deployed MTAs
type MtasCommand struct {
	*BaseCommand
}

func NewMtasCommand() *MtasCommand {
	baseCmd := &BaseCommand{flagsParser: NewDefaultCommandFlagsParser(nil), flagsValidator: NewDefaultCommandFlagsValidator(nil)}
	mtasCmd := &MtasCommand{baseCmd}
	baseCmd.Command = mtasCmd
	return mtasCmd
}

// GetPluginCommand returns the plugin command details
func (c *MtasCommand) GetPluginCommand() plugin.Command {
	return plugin.Command{
		Name:     "mtas",
		HelpText: "List all multi-target apps",
		UsageDetails: plugin.Usage{
			Usage: "cf mtas [-u URL]",
			Options: map[string]string{
				deployServiceURLOpt: "Deploy service URL, by default 'deploy-service.<system-domain>'",
			},
		},
	}
}

func (c *MtasCommand) defineCommandOptions(flags *flag.FlagSet) {
	//no additional options to define
}

func (c *MtasCommand) executeInternal(positionalArgs []string, dsHost string, flags *flag.FlagSet, cfTarget util.CloudFoundryTarget) ExecutionStatus {
	// Print initial message
	ui.Say("Getting multi-target apps in org %s / space %s as %s...",
		terminal.EntityNameColor(cfTarget.Org.Name), terminal.EntityNameColor(cfTarget.Space.Name),
		terminal.EntityNameColor(cfTarget.Username))

	// Create new REST client
	mtaV2Client := c.NewMtaV2Client(dsHost, cfTarget)

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
