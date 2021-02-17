package commands

import (
	"errors"
	"flag"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/clients/baseclient"
	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/clients/models"
	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/configuration"
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
	deleteServiceKeysOpt       = "delete-service-keys"
	keepFilesOpt               = "keep-files"
	skipOwnershipValidationOpt = "skip-ownership-validation"
	moduleOpt                  = "m"
	resourceOpt                = "r"
	allModulesOpt              = "all-modules"
	allResourcesOpt            = "all-resources"
	verifyArchiveSignatureOpt  = "verify-archive-signature"
	strategyOpt                = "strategy"
	skipTestingPhase           = "skip-testing-phase"
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

// DeployCommand is a command for deploying an MTA archive
type DeployCommand struct {
	*BaseCommand
	setProcessParameters ProcessParametersSetter
	processTypeProvider  ProcessTypeProvider
}

// NewDeployCommand creates a new deploy command.
func NewDeployCommand() *DeployCommand {
	baseCmd := &BaseCommand{flagsParser: deployCommandLineArgumentsParser{}, flagsValidator: deployCommandFlagsValidator{}}
	deployCmd := &DeployCommand{baseCmd, deployProcessParametersSetter(), &deployCommandProcessTypeProvider{}}
	baseCmd.Command = deployCmd
	return deployCmd
}

// GetPluginCommand returns the plugin command details
func (c *DeployCommand) GetPluginCommand() plugin.Command {
	return plugin.Command{
		Name:     "deploy",
		HelpText: "Deploy a new multi-target app or sync changes to an existing one",
		UsageDetails: plugin.Usage{
			Usage: `Deploy a multi-target app archive
   cf deploy MTA [-e EXT_DESCRIPTOR[,...]] [-t TIMEOUT] [--version-rule VERSION_RULE] [-u URL] [-f] [--retries RETRIES] [--no-start] [--namespace NAMESPACE] [--delete-services] [--delete-service-keys] [--delete-service-brokers] [--keep-files] [--no-restart-subscribed-apps] [--do-not-fail-on-missing-permissions] [--abort-on-error] [--verify-archive-signature] [--strategy STRATEGY] [--skip-testing-phase]

   Perform action on an active deploy operation
   cf deploy -i OPERATION_ID -a ACTION [-u URL]

   (EXPERIMENTAL) Deploy a multi-target app archive referenced by a remote URL
   cf deploy <MTA archive URL> [-e EXT_DESCRIPTOR[,...]] [-t TIMEOUT] [--version-rule VERSION_RULE] [-u MTA_CONTROLLER_URL] [--retries RETRIES] [--no-start] [--namespace NAMESPACE] [--delete-services] [--delete-service-keys] [--delete-service-brokers] [--keep-files] [--no-restart-subscribed-apps] [--do-not-fail-on-missing-permissions] [--abort-on-error] [--verify-archive-signature] [--strategy STRATEGY] [--skip-testing-phase]`,
			Options: map[string]string{
				extDescriptorsOpt:                      "Extension descriptors",
				deployServiceURLOpt:                    "Deploy service URL, by default 'deploy-service.<system-domain>'",
				timeoutOpt:                             "Start timeout in seconds",
				versionRuleOpt:                         "Version rule (HIGHER, SAME_HIGHER, ALL)",
				operationIDOpt:                         "Active deploy operation ID",
				actionOpt:                              "Action to perform on active deploy operation (abort, retry, resume, monitor)",
				forceOpt:                               "Force deploy without confirmation for aborting conflicting processes",
				moduleOpt:                              "Deploy list of modules which are contained in the deployment descriptor, in the current location",
				resourceOpt:                            "Deploy list of resources which are contained in the deployment descriptor, in the current location",
				util.GetShortOption(noStartOpt):        "Do not start apps",
				util.GetShortOption(namespaceOpt):      "(EXPERIMENTAL) Namespace for the mta, applied to app and service names as well",
				util.GetShortOption(deleteServicesOpt): "Recreate changed services / delete discontinued services",
				util.GetShortOption(deleteServiceKeysOpt):          "Delete existing service keys and apply the new ones",
				util.GetShortOption(deleteServiceBrokersOpt):       "Delete discontinued service brokers",
				util.GetShortOption(keepFilesOpt):                  "Keep files used for deployment",
				util.GetShortOption(noRestartSubscribedAppsOpt):    "Do not restart subscribed apps, updated during the deployment",
				util.GetShortOption(noFailOnMissingPermissionsOpt): "Do not fail on missing permissions for admin operations",
				util.GetShortOption(abortOnErrorOpt):               "Auto-abort the process on any errors",
				util.GetShortOption(allModulesOpt):                 "Deploy all modules which are contained in the deployment descriptor, in the current location",
				util.GetShortOption(allResourcesOpt):               "Deploy all resources which are contained in the deployment descriptor, in the current location",
				util.GetShortOption(verifyArchiveSignatureOpt):     "Verify the archive is correctly signed",
				util.GetShortOption(retriesOpt):                    "Retry the operation N times in case a non-content error occurs (default 3)",
				util.GetShortOption(strategyOpt):                   "Specify the deployment strategy when updating an mta (default, blue-green)",
				util.GetShortOption(skipTestingPhase):              "(STRATEGY: BLUE-GREEN) Do not require confirmation for deleting the previously deployed MTA apps",
			},
		},
	}
}

// ProcessParametersSetter is a function that sets the startup parameters for
// the deploy process. It takes them from the list of parsed flags.
type ProcessParametersSetter func(flags *flag.FlagSet, processBuilder *util.ProcessBuilder)

// deployProcessParametersSetter returns a new ProcessParametersSetter.
func deployProcessParametersSetter() ProcessParametersSetter {
	return func(flags *flag.FlagSet, processBuilder *util.ProcessBuilder) {
		processBuilder.Parameter("deleteServiceKeys", strconv.FormatBool(GetBoolOpt(deleteServiceKeysOpt, flags)))
		processBuilder.Parameter("deleteServices", strconv.FormatBool(GetBoolOpt(deleteServicesOpt, flags)))
		processBuilder.Parameter("noStart", strconv.FormatBool(GetBoolOpt(noStartOpt, flags)))
		processBuilder.Parameter("deleteServiceBrokers", strconv.FormatBool(GetBoolOpt(deleteServiceBrokersOpt, flags)))
		processBuilder.Parameter("startTimeout", GetStringOpt(timeoutOpt, flags))
		processBuilder.Parameter("versionRule", GetStringOpt(versionRuleOpt, flags))
		processBuilder.Parameter("keepFiles", strconv.FormatBool(GetBoolOpt(keepFilesOpt, flags)))
		processBuilder.Parameter("noRestartSubscribedApps", strconv.FormatBool(GetBoolOpt(noRestartSubscribedAppsOpt, flags)))
		processBuilder.Parameter("noFailOnMissingPermissions", strconv.FormatBool(GetBoolOpt(noFailOnMissingPermissionsOpt, flags)))
		processBuilder.Parameter("abortOnError", strconv.FormatBool(GetBoolOpt(abortOnErrorOpt, flags)))
		processBuilder.Parameter("skipOwnershipValidation", strconv.FormatBool(GetBoolOpt(skipOwnershipValidationOpt, flags)))
		processBuilder.Parameter("verifyArchiveSignature", strconv.FormatBool(GetBoolOpt(verifyArchiveSignatureOpt, flags)))
	}
}

func (c *DeployCommand) defineCommandOptions(flags *flag.FlagSet) {
	flags.String(extDescriptorsOpt, "", "")
	flags.String(operationIDOpt, "", "")
	flags.String(actionOpt, "", "")
	flags.Bool(forceOpt, false, "")
	flags.String(timeoutOpt, "", "")
	flags.String(versionRuleOpt, "", "")
	flags.Bool(deleteServicesOpt, false, "")
	flags.Bool(noStartOpt, false, "")
	flags.String(namespaceOpt, "", "")
	flags.Bool(deleteServiceKeysOpt, false, "")
	flags.Bool(deleteServiceBrokersOpt, false, "")
	flags.Bool(keepFilesOpt, false, "")
	flags.Bool(noRestartSubscribedAppsOpt, false, "")
	flags.Bool(noFailOnMissingPermissionsOpt, false, "")
	flags.Bool(abortOnErrorOpt, false, "")
	flags.Bool(skipOwnershipValidationOpt, false, "")
	flags.Bool(allModulesOpt, false, "")
	flags.Bool(allResourcesOpt, false, "")
	flags.Bool(verifyArchiveSignatureOpt, false, "")
	flags.Uint(retriesOpt, 3, "")
	flags.String(strategyOpt, "default", "")
	flags.Bool(skipTestingPhase, false, "")
	flags.Var(&modulesList, moduleOpt, "")
	flags.Var(&resourcesList, resourceOpt, "")
}

func (c *DeployCommand) executeInternal(positionalArgs []string, dsHost string, flags *flag.FlagSet, cfTarget util.CloudFoundryTarget) ExecutionStatus {
	operationID := GetStringOpt(operationIDOpt, flags)
	action := GetStringOpt(actionOpt, flags)
	retries := GetUintOpt(retriesOpt, flags)

	if operationID != "" || action != "" {
		return c.ExecuteAction(operationID, action, retries, dsHost, cfTarget)
	}

	mtaElementsCalculator := mtaElementsToAddCalculator{shouldAddAllModules: false, shouldAddAllResources: false}
	mtaElementsCalculator.calculateElementsToDeploy(flags)

	rawMtaArchive, err := getMtaArchive(positionalArgs, mtaElementsCalculator)
	if err != nil {
		ui.Failed("Error retrieving MTA: %s", err.Error())
		return Failure
	}

	isUrl, mtaArchive := parseMtaArchiveArgument(rawMtaArchive)

	mtaNameToPrint := terminal.EntityNameColor(mtaArchive)
	if isUrl {
		mtaNameToPrint = "from url"
	}

	// Print initial message
	ui.Say("Deploying multi-target app archive %s in org %s / space %s as %s...\n",
		mtaNameToPrint, terminal.EntityNameColor(cfTarget.Org.Name), terminal.EntityNameColor(cfTarget.Space.Name),
		terminal.EntityNameColor(cfTarget.Username))

	var uploadedArchivePartIds []string
	var uploadStatus ExecutionStatus
	var mtaId string

	// Check SLMP metadata
	// TODO: ensure session
	mtaClient := c.NewMtaClient(dsHost, cfTarget)

	namespace := strings.TrimSpace(GetStringOpt(namespaceOpt, flags))
	uploadChunkSizeInMB := configuration.NewSnapshot().GetUploadChunkSizeInMB()
	fileUploader := NewFileUploader(mtaClient, namespace, uploadChunkSizeInMB)

	if isUrl {
		uploadedArchive, err := mtaClient.UploadMtaArchiveFromUrl(mtaArchive, &namespace)
		if err != nil {
			ui.Failed("Could not upload from url: %s", baseclient.NewClientError(err))
			return Failure
		}
		uploadedArchivePartIds = append(uploadedArchivePartIds, uploadedArchive.ID)
		ui.Ok()
	} else {
		// Get the full path of the MTA archive
		mtaArchivePath, err := filepath.Abs(mtaArchive)
		if err != nil {
			ui.Failed("Could not get absolute path of file '%s'", mtaArchive)
			return Failure
		}

		// Extract mta id from archive file
		descriptor, err := util.GetMtaDescriptorFromArchive(mtaArchivePath)
		if os.IsNotExist(err) {
			ui.Failed("Could not find file %s", terminal.EntityNameColor(mtaArchivePath))
			return Failure
		} else if err != nil {
			ui.Failed("Could not get MTA ID from deployment descriptor: %s", err)
			return Failure
		}
		mtaId = descriptor.ID

		force := GetBoolOpt(forceOpt, flags)
		// Check for an ongoing operation for this MTA ID and abort it
		wasAborted, err := c.CheckOngoingOperation(descriptor.ID, namespace, dsHost, force, cfTarget)
		if err != nil {
			ui.Failed("Could not get MTA operations: %s", baseclient.NewClientError(err))
			return Failure
		}
		if !wasAborted {
			return Failure
		}

		// Upload the MTA archive file
		uploadedArchivePartIds, uploadStatus = c.uploadFiles([]string{mtaArchivePath}, fileUploader)
		if uploadStatus == Failure {
			return Failure
		}
	}

	extDescriptors := GetStringOpt(extDescriptorsOpt, flags)
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
	// Upload the extension descriptor files
	uploadedExtDescriptorIDs, uploadStatus := c.uploadFiles(extDescriptorPaths, fileUploader)
	if uploadStatus == Failure {
		return Failure
	}

	// Build the process instance
	processBuilder := NewDeploymentStrategy(flags, c.processTypeProvider).CreateProcessBuilder()
	processBuilder.Namespace(namespace)
	processBuilder.Parameter("appArchiveId", strings.Join(uploadedArchivePartIds, ","))
	processBuilder.Parameter("mtaExtDescriptorId", strings.Join(uploadedExtDescriptorIDs, ","))
	if !isUrl {
		processBuilder.Parameter("mtaId", mtaId)
	}
	setModulesAndResourcesListParameters(modulesList, resourcesList, processBuilder, mtaElementsCalculator)
	c.setProcessParameters(flags, processBuilder)

	operation := processBuilder.Build()

	// Create the new process
	responseHeader, err := mtaClient.StartMtaOperation(*operation)
	if err != nil {
		ui.Failed("Could not create operation: %s", baseclient.NewClientError(err))
		return Failure
	}
	executionMonitor := NewExecutionMonitorFromLocationHeader(c.name, responseHeader.Location.String(), retries, []*models.Message{}, mtaClient)
	ui.Say("Operation ID: %s", terminal.EntityNameColor(executionMonitor.operationID))
	return executionMonitor.Monitor()
}

func parseMtaArchiveArgument(rawMtaArchive interface{}) (bool, string) {
	switch castedMtaArchive := rawMtaArchive.(type) {
	case *url.URL:
		return true, castedMtaArchive.String()
	case string:
		return false, castedMtaArchive
	}
	return false, ""
}

func (c *DeployCommand) uploadFiles(files []string, fileUploader *FileUploader) ([]string, ExecutionStatus) {
	var resultIds []string

	uploadedFiles, status := fileUploader.UploadFiles(files)
	if status == Failure {
		return resultIds, Failure
	}

	for _, uploadedFilePart := range uploadedFiles {
		resultIds = append(resultIds, uploadedFilePart.ID)
	}
	return resultIds, Success
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

func getMtaArchive(parsedArguments []string, mtaElementsCalculator mtaElementsToAddCalculator) (interface{}, error) {
	if len(parsedArguments) == 0 {
		currentWorkingDirectory, err := os.Getwd()
		if err != nil {
			return "", fmt.Errorf("Could not get the current working directory: %s", err.Error())
		}
		return buildMtaArchiveFromDirectory(currentWorkingDirectory, mtaElementsCalculator)
	}

	mtaArgument := parsedArguments[0]

	if matched, _ := regexp.MatchString("^http[s]?://.+", mtaArgument); matched {
		return url.Parse(mtaArgument)
	}

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

func (c *mtaElementsToAddCalculator) calculateElementsToDeploy(flags *flag.FlagSet) {
	allModulesSpecified := GetBoolOpt(allModulesOpt, flags)
	allResourcesSpecified := GetBoolOpt(allResourcesOpt, flags)

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

type deployCommandLineArgumentsParser struct {}

func (p deployCommandLineArgumentsParser) ParseFlags(flags *flag.FlagSet, args []string) error {
	argument := p.findFirstNotFlaggedArgument(flags, args)
	positionalArgumentsToValidate := p.determinePositionalArgumentsToValidate(argument)
	return NewProcessActionExecutorCommandArgumentsParser(positionalArgumentsToValidate).ParseFlags(flags, args)
}

func (deployCommandLineArgumentsParser) findFirstNotFlaggedArgument(flags *flag.FlagSet, args []string) string {
	if len(args) == 0 || flags.Lookup(strings.Replace(args[0], "-", "", 2)) != nil {
		return ""
	}
	return args[0]
}

func (deployCommandLineArgumentsParser) determinePositionalArgumentsToValidate(positionalArgument string) []string {
	if positionalArgument == "" {
		return []string{}
	}
	return []string{"MTA"}
}

type deployCommandFlagsValidator struct{}

func (deployCommandFlagsValidator) ValidateParsedFlags(flags *flag.FlagSet) error {
	var err error
	flags.Visit(func(f *flag.Flag) {
		if f.Name == strategyOpt {
			if f.Value.String() == "" {
				err = errors.New("strategy flag defined but no argument specified")
			} else if !util.Contains(AvailableStrategies(), f.Value.String()) {
				err = fmt.Errorf("%s is not a valid deployment strategy, available strategies: %v", f.Value.String(), AvailableStrategies())
			}
		}
	})
	if err != nil {
		return err
	}
	return NewDefaultCommandFlagsValidator(nil).ValidateParsedFlags(flags)
}
