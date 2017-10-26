package commands

import (
	"flag"
	"strconv"

	"github.com/SAP/cf-mta-plugin/util"
	"github.com/cloudfoundry/cli/plugin"
)

const noConfirmOpt = "no-confirm"

// BlueGreenDeployCommand is a command for blue green deployment of MTAs.
type BlueGreenDeployCommand struct {
	DeployCommand
}

// NewBlueGreenDeployCommand creates a new BlueGreenDeployCommand.
func NewBlueGreenDeployCommand() *BlueGreenDeployCommand {
	blueGreenDeployCommand := &BlueGreenDeployCommand{DeployCommand{BaseCommand{}, blueGreenDeployCommandFlagsDefiner(), blueGreenDeployProcessParametersSetter(), nil}}
	blueGreenDeployCommand.DeployCommand.processTypeProvider = blueGreenDeployCommand
	return blueGreenDeployCommand
}

// GetPluginCommand returns more information for the blue green deploy command.
func (c *BlueGreenDeployCommand) GetPluginCommand() plugin.Command {
	return plugin.Command{
		Name:     "bg-deploy",
		HelpText: "Deploy a multi-target app using blue-green deployment",
		UsageDetails: plugin.Usage{
			Usage: `Deploy a multi-target app using blue-green deployment
   cf bg-deploy MTA [-e EXT_DESCRIPTOR[,...]] [-t TIMEOUT] [--version-rule VERSION_RULE] [-u URL] [-f] [--no-start] [--use-namespaces] [--no-namespaces-for-services] [--delete-services] [--delete-service-keys] [--delete-service-brokers] [--keep-files] [--no-restart-subscribed-apps]  [--no-confirm] [--do-not-fail-on-missing-permissions]

   Perform action on an active deploy operation
   cf deploy -i OPERATION_ID -a ACTION [-u URL]`,
			Options: map[string]string{
				extDescriptorsOpt:                                  "Extension descriptors",
				deployServiceURLOpt:                                "Deploy service URL, by default 'deploy-service.<system-domain>'",
				timeoutOpt:                                         "Start timeout in seconds",
				versionRuleOpt:                                     "Version rule (HIGHER, SAME_HIGHER, ALL)",
				operationIDOpt:                                     "Active deploy operation id",
				actionOpt:                                          "Action to perform on active deploy operation (abort, retry, monitor)",
				forceOpt:                                           "Force deploy without confirmation for aborting conflicting processes",
				util.GetShortOption(noStartOpt):                    "Do not start apps",
				util.GetShortOption(useNamespacesOpt):              "Use namespaces in app and service names",
				util.GetShortOption(noNamespacesForServicesOpt):    "Do not use namespaces in service names",
				util.GetShortOption(deleteServicesOpt):             "Recreate changed services / delete discontinued services",
				util.GetShortOption(deleteServiceKeysOpt):          "Delete existing service keys and apply the new ones",
				util.GetShortOption(deleteServiceBrokersOpt):       "Delete discontinued service brokers",
				util.GetShortOption(keepFilesOpt):                  "Keep files used for deployment",
				util.GetShortOption(noRestartSubscribedAppsOpt):    "Do not restart subscribed apps, updated during the deployment",
				util.GetShortOption(noConfirmOpt):                  "Do not require confirmation for deleting the previously deployed MTA apps",
				util.GetShortOption(noFailOnMissingPermissionsOpt): "Do not fail on missing permissions for admin operations",
			},
		},
	}
}

// BlueGreenDeployCommandFlagsDefiner returns a new CommandFlagsDefiner.
func blueGreenDeployCommandFlagsDefiner() CommandFlagsDefiner {
	return func(flags *flag.FlagSet) map[string]interface{} {
		optionValues := deployCommandFlagsDefiner()(flags)
		optionValues[noConfirmOpt] = flags.Bool(noConfirmOpt, false, "")
		return optionValues
	}
}

// BlueGreenDeployProcessParametersSetter returns a new ProcessParametersSetter.
func blueGreenDeployProcessParametersSetter() ProcessParametersSetter {
	return func(optionValues map[string]interface{}, processBuilder *util.ProcessBuilder) {
		deployProcessParametersSetter()(optionValues, processBuilder)
		processBuilder.Parameter("noConfirm", strconv.FormatBool(GetBoolOpt(noConfirmOpt, optionValues)))
	}
}

func (bg BlueGreenDeployCommand) GetProcessType() string {
	return "bg-deploy"
}
