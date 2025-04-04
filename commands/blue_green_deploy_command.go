package commands

import (
	"flag"
	"os"
	"strconv"
	"time"

	"code.cloudfoundry.org/cli/plugin"
	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/util"
)

const noConfirmOpt = "no-confirm"

// BlueGreenDeployCommand is a command for blue green deployment of MTAs.
type BlueGreenDeployCommand struct {
	*DeployCommand
}

// NewBlueGreenDeployCommand creates a new BlueGreenDeployCommand.
func NewBlueGreenDeployCommand() *BlueGreenDeployCommand {
	baseCmd := &BaseCommand{flagsParser: deployCommandLineArgumentsParser{}, flagsValidator: deployCommandFlagsValidator{}}
	deployCmd := &DeployCommand{baseCmd, blueGreenDeployProcessParametersSetter(), &blueGreenDeployCommandProcessTypeProvider{}, os.Stdin, 30 * time.Second}
	bgDeployCmd := &BlueGreenDeployCommand{deployCmd}
	baseCmd.Command = bgDeployCmd
	return bgDeployCmd
}

// GetPluginCommand returns more information for the blue green deploy command.
func (c *BlueGreenDeployCommand) GetPluginCommand() plugin.Command {
	return plugin.Command{
		Name:     "bg-deploy",
		HelpText: "Deploy a multi-target app using blue-green deployment",
		UsageDetails: plugin.Usage{
			Usage: `Deploy a multi-target app using blue-green deployment
   cf bg-deploy MTA [-e EXT_DESCRIPTOR[,...]] [-t TIMEOUT] [--version-rule VERSION_RULE] [-u URL] [-f] [--retries RETRIES] [--no-start] [--namespace NAMESPACE] [--apply-namespace-app-names true/false] [--apply-namespace-service-names true/false] [--apply-namespace-app-routes true/false] [--apply-namespace-as-suffix true/false ] [--delete-services] [--delete-service-keys] [--delete-service-brokers] [--keep-files] [--no-restart-subscribed-apps] [--no-confirm] [--skip-idle-start] [--do-not-fail-on-missing-permissions] [--abort-on-error] [--apps-start-timeout TIMEOUT] [--apps-stage-timeout TIMEOUT] [--apps-upload-timeout TIMEOUT] [--apps-task-execution-timeout TIMEOUT]

   Perform action on an active deploy operation
   cf deploy -i OPERATION_ID -a ACTION [-u URL]`,
			Options: map[string]string{
				extDescriptorsOpt:                 "Extension descriptors",
				deployServiceURLOpt:               "Deploy service URL, by default 'deploy-service.<system-domain>'",
				versionRuleOpt:                    "Version rule (HIGHER, SAME_HIGHER, ALL)",
				operationIDOpt:                    "Active deploy operation ID",
				actionOpt:                         "Action to perform on active deploy operation (abort, retry, resume, monitor)",
				forceOpt:                          "Force deploy without confirmation for aborting conflicting processes",
				util.GetShortOption(noStartOpt):   "Do not start apps",
				util.GetShortOption(namespaceOpt): "(EXPERIMENTAL) Namespace for the MTA, applied on app names, app routes and service names",
				util.GetShortOption(applyNamespaceAppNamesOpt):                  "(EXPERIMENTAL) Apply namespace to application names: (true, false)",
				util.GetShortOption(applyNamespaceServiceNamesOpt):              "(EXPERIMENTAL) Apply namespace to service names: (true, false)",
				util.GetShortOption(applyNamespaceAppRoutesOpt):                 "(EXPERIMENTAL) Apply namespace to application routes: (true, false)",
				util.GetShortOption(applyNamespaceAsSuffix):                     "(EXPERIMENTAL) Apply namespace as a suffix rather than a prefix: (true, false)",
				util.GetShortOption(deleteServicesOpt):                          "Recreate changed services / delete discontinued services",
				util.GetShortOption(deleteServiceKeysOpt):                       "Delete existing service keys and apply the new ones",
				util.GetShortOption(deleteServiceBrokersOpt):                    "Delete discontinued service brokers",
				util.GetShortOption(keepFilesOpt):                               "Keep files used for deployment",
				util.GetShortOption(noRestartSubscribedAppsOpt):                 "Do not restart subscribed apps, updated during the deployment",
				util.GetShortOption(noConfirmOpt):                               "Do not require confirmation for deleting the previously deployed MTA apps",
				util.GetShortOption(noFailOnMissingPermissionsOpt):              "Do not fail on missing permissions for admin operations",
				util.GetShortOption(abortOnErrorOpt):                            "Auto-abort the process on any errors",
				util.GetShortOption(retriesOpt):                                 "Retry the operation N times in case a non-content error occurs (default 3)",
				util.GetShortOption(skipIdleStart):                              "Directly start the new MTA version as 'live', skipping the 'idle' phase of the resources. Do not require further confirmation or testing before deleting the old version",
				util.GetShortOption(stageTimeoutOpt):                            "Stage app timeout in seconds",
				util.GetShortOption(uploadTimeoutOpt):                           "Upload app timeout in seconds",
				util.GetShortOption(taskExecutionTimeoutOpt):                    "Task execution timeout in seconds",
				util.CombineFullAndShortParameters(startTimeoutOpt, timeoutOpt): "Start app timeout in seconds",
				util.GetShortOption(shouldBackupPreviousVersionOpt):             "(EXPERIMENTAL) Backup previous version of applications, use new cli command \"rollback-mta\" to rollback to the previous version",
			},
		},
	}
}

func (c *BlueGreenDeployCommand) defineCommandOptions(flags *flag.FlagSet) {
	c.DeployCommand.defineCommandOptions(flags)
	flags.Bool(noConfirmOpt, false, "")
}

// blueGreenDeployProcessParametersSetter returns a new ProcessParametersSetter.
func blueGreenDeployProcessParametersSetter() ProcessParametersSetter {
	return func(flags *flag.FlagSet, processBuilder *util.ProcessBuilder) {
		deployProcessParametersSetter()(flags, processBuilder)
		processBuilder.Parameter("keepOriginalAppNamesAfterDeploy", strconv.FormatBool(false))
		if GetBoolOpt(skipIdleStart, flags) {
			processBuilder.Parameter("noConfirm", strconv.FormatBool(true))
			processBuilder.Parameter("skipIdleStart", strconv.FormatBool(true))
			return
		}
		processBuilder.Parameter("noConfirm", strconv.FormatBool(GetBoolOpt(noConfirmOpt, flags)))
	}
}

type blueGreenDeployCommandProcessTypeProvider struct{}

func (bg blueGreenDeployCommandProcessTypeProvider) GetProcessType() string {
	return "BLUE_GREEN_DEPLOY"
}
