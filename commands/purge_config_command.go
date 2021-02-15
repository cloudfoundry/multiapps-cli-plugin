package commands

import (
	"flag"
	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/clients/baseclient"
	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/ui"
	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/util"
	"github.com/cloudfoundry/cli/cf/terminal"
	"github.com/cloudfoundry/cli/plugin"
)

type PurgeConfigCommand struct {
	*BaseCommand
}

func NewPurgeConfigCommand() *PurgeConfigCommand {
	baseCmd := &BaseCommand{flagsParser: NewDefaultCommandFlagsParser(nil), flagsValidator: NewDefaultCommandFlagsValidator(nil)}
	purgeConfigCmd := &PurgeConfigCommand{baseCmd}
	baseCmd.Command = purgeConfigCmd
	return purgeConfigCmd
}

func (c *PurgeConfigCommand) GetPluginCommand() plugin.Command {
	return plugin.Command{
		Name:     "purge-mta-config",
		HelpText: "Purge no longer valid configuration entries",
		UsageDetails: plugin.Usage{
			Usage: "cf purge-mta-config [-u URL]",
			Options: map[string]string{
				deployServiceURLOpt: "Deploy service URL, by default 'deploy-service.<system-domain>'",
			},
		},
	}
}

func (c *PurgeConfigCommand) defineCommandOptions(flags *flag.FlagSet) {
	//no additional options to define
}

func (c *PurgeConfigCommand) executeInternal(positionalArgs []string, dsHost string, flags *flag.FlagSet, cfTarget util.CloudFoundryTarget) ExecutionStatus {
	ui.Say("Purging configuration entries in org %s / space %s as %s",
		terminal.EntityNameColor(cfTarget.Org.Name), terminal.EntityNameColor(cfTarget.Space.Name),
		terminal.EntityNameColor(cfTarget.Username))

	rc := c.NewRestClient(dsHost)
	// TODO: ensure session

	if err := rc.PurgeConfiguration(cfTarget.Org.Name, cfTarget.Space.Name); err != nil {
		ui.Failed("Could not purge configuration: %v\n", baseclient.NewClientError(err))
		return Failure
	}
	ui.Ok()
	return Success
}
