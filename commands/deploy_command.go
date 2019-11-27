package commands

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/clients/baseclient"
	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/clients/models"
	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/configuration"
	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/log"
	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/ui"
	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/util"
	"github.com/cloudfoundry/cli/cf/terminal"
	"github.com/cloudfoundry/cli/plugin"
)

const (
	extDescriptorsOpt          = "e"
	timeoutOpt                 = "t"
	versionRuleOpt             = "version-rule"
	noStartOpt                 = "no-start"
	useNamespacesOpt           = "use-namespaces"
	noNamespacesForServicesOpt = "no-namespaces-for-services"
	deleteServiceKeysOpt       = "delete-service-keys"
	keepFilesOpt               = "keep-files"
	skipOwnershipValidationOpt = "skip-ownership-validation"
	moduleOpt                  = "m"
	resourceOpt                = "r"
	allModulesOpt              = "all-modules"
	allResourcesOpt            = "all-resources"
	verifyArchiveSignatureOpt  = "verify-archive-signature"
)

type listFlag struct {
	elements []string
}

func (variable listFlag) getElements() []string {
	return variable.elements
}

func (variable listFlag) getProcessList() string {
	return strings.Join(variable.elements, ",")
}

func (variable *listFlag) String() string {
	return fmt.Sprint(variable.elements)

}
func (variable *listFlag) Set(value string) error {
	variable.elements = append(variable.elements, value)
	return nil
}

var modulesList listFlag
var resourcesList listFlag
var reportedProgressMessages []string
var mtaElementsCalculator mtaElementsToAddCalculator

// DeployCommand is a command for deploying an MTA archive
type DeployCommand struct {
	BaseCommand
	commandFlagsDefiner     CommandFlagsDefiner
	processParametersSetter ProcessParametersSetter
	processTypeProvider     ProcessTypeProvider
}

// NewDeployCommand creates a new deploy command.
func NewDeployCommand() *DeployCommand {
	return &DeployCommand{BaseCommand{}, deployCommandFlagsDefiner(), deployProcessParametersSetter(), &deployCommandProcessTypeProvider{}}
}

// GetPluginCommand returns the plugin command details
func (c *DeployCommand) GetPluginCommand() plugin.Command {
	return plugin.Command{
		Name:     "deploy",
		HelpText: "Deploy a new multi-target app or sync changes to an existing one",
		UsageDetails: plugin.Usage{
			Usage: `Deploy a multi-target app archive
   cf deploy MTA [-e EXT_DESCRIPTOR[,...]] [-t TIMEOUT] [--version-rule VERSION_RULE] [-u URL] [-f] [--retries RETRIES] [--no-start] [--use-namespaces] [--no-namespaces-for-services] [--delete-services] [--delete-service-keys] [--delete-service-brokers] [--keep-files] [--no-restart-subscribed-apps] [--do-not-fail-on-missing-permissions] [--abort-on-error] [--skip-ownership-validation] [--verify-archive-signature]

   Perform action on an active deploy operation
   cf deploy -i OPERATION_ID -a ACTION [-u URL]`,
			Options: map[string]string{
				extDescriptorsOpt:                     "Extension descriptors",
				deployServiceURLOpt:                   "Deploy service URL, by default 'deploy-service.<system-domain>'",
				timeoutOpt:                            "Start timeout in seconds",
				versionRuleOpt:                        "Version rule (HIGHER, SAME_HIGHER, ALL)",
				operationIDOpt:                        "Active deploy operation id",
				actionOpt:                             "Action to perform on active deploy operation (abort, retry, monitor)",
				forceOpt:                              "Force deploy without confirmation for aborting conflicting processes",
				moduleOpt:                             "Deploy list of modules which are contained in the deployment descriptor, in the current location",
				resourceOpt:                           "Deploy list of resources which are contained in the deployment descriptor, in the current location",
				util.GetShortOption(noStartOpt):       "Do not start apps",
				util.GetShortOption(useNamespacesOpt): "Use namespaces in app and service names",
				util.GetShortOption(noNamespacesForServicesOpt):    "Do not use namespaces in service names",
				util.GetShortOption(deleteServicesOpt):             "Recreate changed services / delete discontinued services",
				util.GetShortOption(deleteServiceKeysOpt):          "Delete existing service keys and apply the new ones",
				util.GetShortOption(deleteServiceBrokersOpt):       "Delete discontinued service brokers",
				util.GetShortOption(keepFilesOpt):                  "Keep files used for deployment",
				util.GetShortOption(noRestartSubscribedAppsOpt):    "Do not restart subscribed apps, updated during the deployment",
				util.GetShortOption(noFailOnMissingPermissionsOpt): "Do not fail on missing permissions for admin operations",
				util.GetShortOption(abortOnErrorOpt):               "Auto-abort the process on any errors",
				util.GetShortOption(skipOwnershipValidationOpt):    "Skip the ownership validation that prevents the modification of entities managed by other multi-target apps",
				util.GetShortOption(allModulesOpt):                 "Deploy all modules which are contained in the deployment descriptor, in the current location",
				util.GetShortOption(allResourcesOpt):               "Deploy all resources which are contained in the deployment descriptor, in the current location",
				util.GetShortOption(verifyArchiveSignatureOpt):     "Verify the archive is correctly signed",
				util.GetShortOption(retriesOpt):                    "Retry the operation N times in case a non-content error occurs (default 3)",
			},
		},
	}
}

// ProcessParametersSetter is a function that sets the startup parameters for
// the deploy process. It takes them from the list of parsed flags.
type ProcessParametersSetter func(options map[string]interface{}, processBuilder *util.ProcessBuilder)

// DeployCommandFlagsDefiner returns a new CommandFlagsDefiner.
func deployCommandFlagsDefiner() CommandFlagsDefiner {
	return func(flags *flag.FlagSet) map[string]interface{} {
		optionValues := make(map[string]interface{})
		optionValues[extDescriptorsOpt] = flags.String(extDescriptorsOpt, "", "")
		optionValues[operationIDOpt] = flags.String(operationIDOpt, "", "")
		optionValues[actionOpt] = flags.String(actionOpt, "", "")
		optionValues[forceOpt] = flags.Bool(forceOpt, false, "")
		optionValues[timeoutOpt] = flags.String(timeoutOpt, "", "")
		optionValues[versionRuleOpt] = flags.String(versionRuleOpt, "", "")
		optionValues[deleteServicesOpt] = flags.Bool(deleteServicesOpt, false, "")
		optionValues[noStartOpt] = flags.Bool(noStartOpt, false, "")
		optionValues[useNamespacesOpt] = flags.Bool(useNamespacesOpt, false, "")
		optionValues[noNamespacesForServicesOpt] = flags.Bool(noNamespacesForServicesOpt, false, "")
		optionValues[deleteServiceKeysOpt] = flags.Bool(deleteServiceKeysOpt, false, "")
		optionValues[deleteServiceBrokersOpt] = flags.Bool(deleteServiceBrokersOpt, false, "")
		optionValues[keepFilesOpt] = flags.Bool(keepFilesOpt, false, "")
		optionValues[noRestartSubscribedAppsOpt] = flags.Bool(noRestartSubscribedAppsOpt, false, "")
		optionValues[noFailOnMissingPermissionsOpt] = flags.Bool(noFailOnMissingPermissionsOpt, false, "")
		optionValues[abortOnErrorOpt] = flags.Bool(abortOnErrorOpt, false, "")
		optionValues[skipOwnershipValidationOpt] = flags.Bool(skipOwnershipValidationOpt, false, "")
		optionValues[allModulesOpt] = flags.Bool(allModulesOpt, false, "")
		optionValues[allResourcesOpt] = flags.Bool(allResourcesOpt, false, "")
		optionValues[verifyArchiveSignatureOpt] = flags.Bool(verifyArchiveSignatureOpt, false, "")
		optionValues[retriesOpt] = flags.Uint(retriesOpt, 3, "")
		flags.Var(&modulesList, moduleOpt, "")
		flags.Var(&resourcesList, resourceOpt, "")
		return optionValues
	}
}

// DeployProcessParametersSetter returns a new ProcessParametersSetter.
func deployProcessParametersSetter() ProcessParametersSetter {
	return func(optionValues map[string]interface{}, processBuilder *util.ProcessBuilder) {
		processBuilder.Parameter("deleteServiceKeys", strconv.FormatBool(GetBoolOpt(deleteServiceKeysOpt, optionValues)))
		processBuilder.Parameter("deleteServices", strconv.FormatBool(GetBoolOpt(deleteServicesOpt, optionValues)))
		processBuilder.Parameter("noStart", strconv.FormatBool(GetBoolOpt(noStartOpt, optionValues)))
		processBuilder.Parameter("useNamespaces", strconv.FormatBool(GetBoolOpt(useNamespacesOpt, optionValues)))
		processBuilder.Parameter("useNamespacesForServices", strconv.FormatBool(!GetBoolOpt(noNamespacesForServicesOpt, optionValues)))
		processBuilder.Parameter("deleteServiceBrokers", strconv.FormatBool(GetBoolOpt(deleteServiceBrokersOpt, optionValues)))
		processBuilder.Parameter("startTimeout", GetStringOpt(timeoutOpt, optionValues))
		processBuilder.Parameter("versionRule", GetStringOpt(versionRuleOpt, optionValues))
		processBuilder.Parameter("keepFiles", strconv.FormatBool(GetBoolOpt(keepFilesOpt, optionValues)))
		processBuilder.Parameter("noRestartSubscribedApps", strconv.FormatBool(GetBoolOpt(noRestartSubscribedAppsOpt, optionValues)))
		processBuilder.Parameter("noFailOnMissingPermissions", strconv.FormatBool(GetBoolOpt(noFailOnMissingPermissionsOpt, optionValues)))
		processBuilder.Parameter("abortOnError", strconv.FormatBool(GetBoolOpt(abortOnErrorOpt, optionValues)))
		processBuilder.Parameter("skipOwnershipValidation", strconv.FormatBool(GetBoolOpt(skipOwnershipValidationOpt, optionValues)))
		processBuilder.Parameter("verifyArchiveSignature", strconv.FormatBool(GetBoolOpt(verifyArchiveSignatureOpt, optionValues)))
	}
}

func getMtaElementsList(mtaElements []string, optionValues map[string]interface{}) string {
	return strings.Join(mtaElements, ",")
}

// GetBoolOpt gets and dereferences the pointer identified by the specified name.
func GetBoolOpt(name string, optionValues map[string]interface{}) bool {
	return *optionValues[name].(*bool)
}

// GetStringOpt gets and dereferences the pointer identified by the specified name.
func GetStringOpt(name string, optionValues map[string]interface{}) string {
	return *optionValues[name].(*string)
}

// GetUintOpt gets and dereferences the pointer identified by the specified name.
func GetUintOpt(name string, optionValues map[string]interface{}) uint {
	return *optionValues[name].(*uint)
}

// Execute executes the command
func (c *DeployCommand) Execute(args []string) ExecutionStatus {
	log.Tracef("Executing command '"+c.name+"': args: '%v'\n", args)

	var host string

	// Parse command arguments and check for required options
	flags, err := c.CreateFlags(&host, args)
	if err != nil {
		ui.Failed(err.Error())
		return Failure
	}
	optionValues := c.commandFlagsDefiner(flags)
	parser := NewCommandFlagsParser(flags, newDeployCommandLineArgumentsParser(), NewDefaultCommandFlagsValidator(nil))
	err = parser.Parse(args)
	if err != nil {
		c.Usage(err.Error())
		return Failure
	}

	extDescriptors := GetStringOpt(extDescriptorsOpt, optionValues)
	operationID := GetStringOpt(operationIDOpt, optionValues)
	action := GetStringOpt(actionOpt, optionValues)
	force := GetBoolOpt(forceOpt, optionValues)
	retries := GetUintOpt(retriesOpt, optionValues)

	context, err := c.GetContext()
	if err != nil {
		ui.Failed(err.Error())
		return Failure
	}

	if operationID != "" || action != "" {
		return c.ExecuteAction(operationID, action, retries, host)
	}
	mtaElementsCalculator := mtaElementsToAddCalculator{shouldAddAllModules: false, shouldAddAllResources: false}
	mtaElementsCalculator.calculateElementsToDeploy(optionValues)

	mtaArchive, err := getMtaArchive(parser.Args(), mtaElementsCalculator)
	if err != nil {
		ui.Failed("Error retrieving MTA: %s", err.Error())
		return Failure
	}

	// Print initial message
	ui.Say("Deploying multi-target app archive %s in org %s / space %s as %s...\n",
		terminal.EntityNameColor(mtaArchive), terminal.EntityNameColor(context.Org),
		terminal.EntityNameColor(context.Space), terminal.EntityNameColor(context.Username))

	// Get the full path of the MTA archive
	mtaArchivePath, err := filepath.Abs(mtaArchive)
	if err != nil {
		ui.Failed("Could not get absolute path of file '%s'", mtaArchive)
		return Failure
	}

	// Get the full paths of the extension descriptors
	var extDescriptorPaths []string
	if extDescriptors != "" {
		extDescriptorFiles := strings.Split(extDescriptors, ",")
		for _, extDescriptorFile := range extDescriptorFiles {
			extDescriptorPath, err := filepath.Abs(extDescriptorFile)
			if err != nil {
				ui.Failed("Could not get absolute path of file '%s'", extDescriptorFile)
				return Failure
			}
			extDescriptorPaths = append(extDescriptorPaths, extDescriptorPath)
		}
	}

	streaming := configuration.IsStreamingFlagSet()

	if !streaming { //TODO extract to a function - check for and resolve conflict
		// Extract mta id from archive file
		mtaID, err := util.GetMtaIDFromArchive(mtaArchivePath)
		if os.IsNotExist(err) {
			ui.Failed("Could not find file %s", terminal.EntityNameColor(mtaArchivePath))
			return Failure
		} else if err != nil {
			ui.Failed("Could not get MTA id from deployment descriptor: %s", err)
			return Failure
		}

		// Check for an ongoing operation for this MTA ID and abort it
		wasAborted, err := c.CheckOngoingOperation(mtaID, host, force)
		if err != nil {
			ui.Failed("Could not get MTA operations: %s", baseclient.NewClientError(err))
			return Failure
		}
		if !wasAborted {
			return Failure
		}
	}
	// Check SLMP metadata
	// TODO: ensure session
	mtaClient, err := c.NewMtaClient(host)
	if err != nil {
		ui.Failed("Could not get space guid:", baseclient.NewClientError(err))
		return Failure
	}
	// Upload the MTA archive file
	mtaArchiveUploader := NewFileUploader([]string{mtaArchivePath}, mtaClient)
	uploadedMtaArchives, status := mtaArchiveUploader.UploadFiles()
	if status == Failure {
		return Failure
	}
	var uploadedArchivePartIds []string
	for _, uploadedMtaArchivePart := range uploadedMtaArchives {
		uploadedArchivePartIds = append(uploadedArchivePartIds, uploadedMtaArchivePart.ID)
	}

	// Upload the extension descriptor files
	var uploadedExtDescriptorIDs []string
	if len(extDescriptorPaths) != 0 {
		extDescriptorsUploader := NewFileUploader(extDescriptorPaths, mtaClient)
		uploadedExtDescriptors, status := extDescriptorsUploader.UploadFiles()
		if status == Failure {
			return Failure
		}
		for _, uploadedExtDescriptor := range uploadedExtDescriptors {
			uploadedExtDescriptorIDs = append(uploadedExtDescriptorIDs, uploadedExtDescriptor.ID)
		}
	}

	// Build the process instance
	processBuilder := util.NewProcessBuilder()
	processBuilder.ProcessType(c.processTypeProvider.GetProcessType())
	processBuilder.Parameter("appArchiveId", strings.Join(uploadedArchivePartIds, ","))
	processBuilder.Parameter("mtaExtDescriptorId", strings.Join(uploadedExtDescriptorIDs, ","))
	setModulesAndResourcesListParameters(modulesList, resourcesList, processBuilder, mtaElementsCalculator)
	c.processParametersSetter(optionValues, processBuilder)
	operation := processBuilder.Build()

	// Create the new process
	responseHeader, err := mtaClient.StartMtaOperation(*operation)
	if err != nil {
		ui.Failed("Could not create operation: %s", baseclient.NewClientError(err))
		return Failure
	}

	return NewExecutionMonitorFromLocationHeader(c.name, responseHeader.Location.String(), retries, []*models.Message{}, mtaClient).Monitor()
}

func setModulesAndResourcesListParameters(modulesList, resourcesList listFlag, processBuilder *util.ProcessBuilder, mtaElementsCalculator mtaElementsToAddCalculator) {

	if mtaElementsCalculator.shouldAddAllModules && mtaElementsCalculator.shouldAddAllResources {
		return
	}

	if mtaElementsCalculator.shouldAddAllModules && !mtaElementsCalculator.shouldAddAllResources {
		processBuilder.SetParameterWithoutCheck("resourcesForDeployment", resourcesList.getProcessList())
		return
	}

	if !mtaElementsCalculator.shouldAddAllModules && mtaElementsCalculator.shouldAddAllResources {
		processBuilder.SetParameterWithoutCheck("modulesForDeployment", modulesList.getProcessList())
		return
	}

	processBuilder.SetParameterWithoutCheck("resourcesForDeployment", resourcesList.getProcessList())
	processBuilder.SetParameterWithoutCheck("modulesForDeployment", modulesList.getProcessList())
}

func getMtaArchive(parsedArguments []string, mtaElementsCalculator mtaElementsToAddCalculator) (string, error) {
	if len(parsedArguments) == 0 {
		currentWorkingDirectory, err := os.Getwd()
		if err != nil {
			return "", fmt.Errorf("Could not get the current working directory: %s", err.Error())
		}
		return buildMtaArchiveFromDirectory(currentWorkingDirectory, mtaElementsCalculator)
	}

	mtaArgument := parsedArguments[0]
	fileInfo, err := os.Stat(mtaArgument)
	if err != nil && os.IsNotExist(err) {
		return "", fmt.Errorf("Could not find MTA %s", mtaArgument)
	}

	if !fileInfo.IsDir() {
		return mtaArgument, nil
	}

	return buildMtaArchiveFromDirectory(mtaArgument, mtaElementsCalculator)
}

func buildMtaArchiveFromDirectory(mtaDirectoryLocation string, mtaElementsCalculator mtaElementsToAddCalculator) (string, error) {
	modulesToAdd, err := getModulesToAdd(mtaDirectoryLocation, mtaElementsCalculator)
	if err != nil {
		return "", err
	}

	resourcesToAdd, err := getResourcesToAdd(mtaDirectoryLocation, mtaElementsCalculator)
	if err != nil {
		return "", err
	}

	return util.NewMtaArchiveBuilder(modulesToAdd, resourcesToAdd).Build(mtaDirectoryLocation)
}

func getModulesToAdd(mtaDirectoryLocation string, mtaElementsCalculator mtaElementsToAddCalculator) ([]string, error) {
	if mtaElementsCalculator.shouldAddAllModules {
		deploymentDescriptor, _, err := util.ParseDeploymentDescriptor(mtaDirectoryLocation)
		if err != nil {
			return []string{}, err
		}

		modulesToAdd := make([]string, 0)
		for _, module := range deploymentDescriptor.Modules {
			modulesToAdd = append(modulesToAdd, module.Name)
		}
		return modulesToAdd, nil
	}

	return modulesList.getElements(), nil
}

func getResourcesToAdd(mtaDirectoryLocation string, mtaElementsCalculator mtaElementsToAddCalculator) ([]string, error) {
	if mtaElementsCalculator.shouldAddAllResources {
		deploymentDescriptor, _, err := util.ParseDeploymentDescriptor(mtaDirectoryLocation)
		if err != nil {
			return []string{}, err
		}

		resourcesToAdd := make([]string, 0)
		for _, resource := range deploymentDescriptor.Resources {
			resourcesToAdd = append(resourcesToAdd, resource.Name)
		}
		return resourcesToAdd, nil
	}

	return resourcesList.getElements(), nil
}

type mtaElementsToAddCalculator struct {
	shouldAddAllModules   bool
	shouldAddAllResources bool
}

func (c *mtaElementsToAddCalculator) calculateElementsToDeploy(optionValues map[string]interface{}) {
	allModulesSpecified := GetBoolOpt(allModulesOpt, optionValues)
	allResourcesSpecified := GetBoolOpt(allResourcesOpt, optionValues)

	if !allResourcesSpecified && len(resourcesList.getElements()) == 0 && !allModulesSpecified && len(modulesList.getElements()) == 0 {
		// --all-resources ==false && no -r
		c.shouldAddAllResources = true
		c.shouldAddAllModules = true
		return
	}

	if allModulesSpecified {
		// --all-modules ==true , no matter if there is -m
		c.shouldAddAllModules = true
	}

	if allResourcesSpecified {
		// --all-modules ==true , no matter if there is -m
		c.shouldAddAllResources = true
	}
}

type deployCommandProcessTypeProvider struct{}

func (d deployCommandProcessTypeProvider) GetProcessType() string {
	return "DEPLOY"
}

type deployCommandLineArgumentsParser struct {
}

func newDeployCommandLineArgumentsParser() deployCommandLineArgumentsParser {
	return deployCommandLineArgumentsParser{}
}

func (p deployCommandLineArgumentsParser) ParseFlags(flags *flag.FlagSet, args []string) error {
	argument := findFirstNotFlagedArgument(flags, args)

	positionalArgumentsToValidate := determinePositionalArgumentsTovalidate(argument)

	return NewProcessActionExecutorCommandArgumentsParser(positionalArgumentsToValidate).ParseFlags(flags, args)
}

func findFirstNotFlagedArgument(flags *flag.FlagSet, args []string) string {
	if len(args) == 0 {
		return ""
	}
	optionFlag := flags.Lookup(strings.Replace(args[0], "-", "", 2))
	if optionFlag == nil {
		return args[0]
	}
	return ""
}

func determinePositionalArgumentsTovalidate(possitionalArgument string) []string {
	if possitionalArgument == "" {
		return []string{}
	}

	return []string{"MTA"}
}
