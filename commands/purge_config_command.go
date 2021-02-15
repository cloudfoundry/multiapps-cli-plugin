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

	cfTarget, err := c.GetCFTarget()
	if err != nil {
		ui.Failed(err.Error())
		return Failure
	}

	ui.Say("Purging configuration entries in org %s / space %s as %s",
		terminal.EntityNameColor(cfTarget.Org.Name), terminal.EntityNameColor(cfTarget.Space.Name),
		terminal.EntityNameColor(cfTarget.Username))

	rc := c.NewRestClient(host)
	// TODO: ensure session

	if err := rc.PurgeConfiguration(cfTarget.Org.Name, cfTarget.Space.Name); err != nil {
		ui.Failed("Could not purge configuration: %v\n", baseclient.NewClientError(err))
		return Failure
	}
	ui.Ok()
	return Success
}
