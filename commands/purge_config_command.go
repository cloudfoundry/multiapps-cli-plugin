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

func (c *PurgeConfigCommand) Execute(args []string) ExecutionStatus {
	log.Tracef("Executing command %q with args %v\n", c.name, args)

	var host string
	flags, err := c.CreateFlags(&host)
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

	ui.Say("Purging configuration entries in org %s / space %s as %s",
		terminal.EntityNameColor(context.Org),
		terminal.EntityNameColor(context.Space),
		terminal.EntityNameColor(context.Username))

	rc, err := c.NewRestClient(host)
	if err != nil {
		c.reportError(baseclient.NewClientError(err))
		return Failure
	}
	// TODO: ensure session

	if err := rc.PurgeConfiguration(context.Org, context.Space); err != nil {
		c.reportError(baseclient.NewClientError(err))
		return Failure
	}
	ui.Ok()
	return Success
}

func (c *PurgeConfigCommand) reportError(err error) {
	ui.Failed("Could not purge configuration: %v\n", err)
}
