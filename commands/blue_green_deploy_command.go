package commands

import (
	"strconv"

	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/util"
	"github.com/cloudfoundry/cli/plugin"
)

const noConfirmOpt = "no-confirm"

// BlueGreenDeployCommand is a command for blue green deployment of MTAs.
type BlueGreenDeployCommand struct {
	DeployCommand
}

// NewBlueGreenDeployCommand creates a new BlueGreenDeployCommand.
func NewBlueGreenDeployCommand() *BlueGreenDeployCommand {
	return &BlueGreenDeployCommand{DeployCommand{BaseCommand{options: getBlueGreenDeployCommandOptions()}, blueGreenDeployProcessParametersSetter(), blueGreenDeployCommandProcessTypeProvider{}}}
}

func getBlueGreenDeployCommandOptions() map[string]CommandOption {
	deployOptions := getDeployCommandOptions()
	delete(deployOptions, strategyOpt)
	return deployOptions
}

// GetPluginCommand returns more information for the blue green deploy command.
func (c *BlueGreenDeployCommand) GetPluginCommand() plugin.Command {
	return plugin.Command{
		Name:     "bg-deploy",
		HelpText: "Deploy a multi-target app using blue-green deployment",
		UsageDetails: plugin.Usage{
			Usage: `Deploy a multi-target app using blue-green deployment
   cf bg-deploy MTA [-e EXT_DESCRIPTOR[,...]] [-t TIMEOUT] [--version-rule VERSION_RULE] [-u URL] [-f]  [--retries RETRIES] [--no-start] [--use-namespaces] [--no-namespaces-for-services] [--delete-services] [--delete-service-keys] [--delete-service-brokers] [--keep-files] [--no-restart-subscribed-apps]  [--no-confirm] [--do-not-fail-on-missing-permissions] [--abort-on-error] [--skip-ownership-validation] [--verify-archive-signature]

   Perform action on an active deploy operation
   cf deploy -i OPERATION_ID -a ACTION [-u URL]`,
			Options: c.getOptionsForPluginCommand(),
		},
	}
}

// BlueGreenDeployProcessParametersSetter returns a new ProcessParametersSetter.
func blueGreenDeployProcessParametersSetter() ProcessParametersSetter {
	return func(options map[string]CommandOption, processBuilder *util.ProcessBuilder) {
		deployProcessParametersSetter()(options, processBuilder)
		processBuilder.Parameter("noConfirm", strconv.FormatBool(getBoolOpt(noConfirmOpt, options)))
		processBuilder.Parameter("keepOriginalAppNamesAfterDeploy", strconv.FormatBool(false))
	}
}

type blueGreenDeployCommandProcessTypeProvider struct{}

func (bg blueGreenDeployCommandProcessTypeProvider) GetProcessType() string {
	return "BLUE_GREEN_DEPLOY"
}
