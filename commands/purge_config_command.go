package commands

import (
	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/clients/baseclient"
	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/log"
	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/ui"
	"github.com/cloudfoundry/cli/cf/terminal"
	"github.com/cloudfoundry/cli/plugin"
)

type PurgeConfigCommand struct {
	BaseCommand
}

func NewPurgeConfigCommand() *PurgeConfigCommand {
	return &PurgeConfigCommand{BaseCommand{options: getPurgeConfigCommandOptions()}}
}

func getPurgeConfigCommandOptions() map[string]CommandOption {
	return map[string]CommandOption{
		deployServiceURLOpt: deployServiceUrlOption(),
	}
}

func (c *PurgeConfigCommand) GetPluginCommand() plugin.Command {
	return plugin.Command{
		Name:     "purge-mta-config",
		HelpText: "Purge no longer valid configuration entries",
		UsageDetails: plugin.Usage{
			Usage:   "cf purge-mta-config [-u URL]",
			Options: c.getOptionsForPluginCommand(),
		},
	}
}

func (c *PurgeConfigCommand) Execute(args []string) ExecutionStatus {
	log.Tracef("Executing command %q with args %v\n", c.name, args)

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

	ui.Say("Purging configuration entries in org %s / space %s as %s",
		terminal.EntityNameColor(context.Org),
		terminal.EntityNameColor(context.Space),
		terminal.EntityNameColor(context.Username))

	rc := c.NewRestClient(host)
	// TODO: ensure session

	if err := rc.PurgeConfiguration(context.Org, context.Space); err != nil {
		ui.Failed("Could not purge configuration: %v\n", baseclient.NewClientError(err))
		return Failure
	}
	ui.Ok()
	return Success
}
