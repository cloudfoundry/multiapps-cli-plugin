package commands

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/clients/baseclient"
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
	strategyOpt                = "strategy"
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

func (variable listFlag) String() string {
	return fmt.Sprint(variable.elements)
}

func (variable listFlag) Set(value string) error {
	variable.elements = append(variable.elements, value)
	return nil
}

var modulesList listFlag
var resourcesList listFlag

// DeployCommand is a command for deploying an MTA archive
type DeployCommand struct {
	BaseCommand
	processParametersSetter ProcessParametersSetter
	processTypeProvider     ProcessTypeProvider
}

// NewDeployCommand creates a new deploy command.
func NewDeployCommand() *DeployCommand {
	return &DeployCommand{BaseCommand{optionParser: NewDeployCommandOptionParser(), options: getDeployCommandOptions()}, deployProcessParametersSetter(), deployCommandProcessTypeProvider{}}
}

func getDeployCommandOptions() map[string]CommandOption {
	return map[string]CommandOption{
		deployServiceURLOpt:           deployServiceUrlOption(),
		extDescriptorsOpt:             {new(string), "", "Extension descriptors", true},
		operationIDOpt:                {new(string), "", "Active deploy operation id", true},
		actionOpt:                     {new(string), "", "Action to perform on active deploy operation (abort, retry, monitor)", true},
		forceOpt:                      {new(bool), false, "Force deploy without confirmation for aborting conflicting processes", true},
		timeoutOpt:                    {new(string), "", "Start timeout in seconds", true},
		moduleOpt:                     {Value: &modulesList, Usage: "Deploy list of modules which are contained in the deployment descriptor, in the current location", IsShortOpt: true},
		resourceOpt:                   {Value: &resourcesList, Usage: "Deploy list of resources which are contained in the deployment descriptor, in the current location", IsShortOpt: true},
		versionRuleOpt:                {new(string), "", "Version rule (HIGHER, SAME_HIGHER, ALL)", false},
		deleteServicesOpt:             {new(bool), false, "Recreate changed services / delete discontinued services", false},
		deleteServiceKeysOpt:          {new(bool), false, "Delete existing service keys and apply the new ones", false},
		deleteServiceBrokersOpt:       {new(bool), false, "Delete discontinued service brokers", false},
		noStartOpt:                    {new(bool), false, "Do not start apps", false},
		useNamespacesOpt:              {new(bool), false, "Use namespaces in app and service names", false},
		noNamespacesForServicesOpt:    {new(bool), false, "Do not use namespaces in service names", false},
		keepFilesOpt:                  {new(bool), false, "Keep files used for deployment", false},
		noRestartSubscribedAppsOpt:    {new(bool), false, "Do not restart subscribed apps, updated during the deployment", false},
		noFailOnMissingPermissionsOpt: {new(bool), false, "Do not fail on missing permissions for admin operations", false},
		abortOnErrorOpt:               {new(bool), false, "Auto-abort the process on any errors", false},
		skipOwnershipValidationOpt:    {new(bool), false, "Skip the ownership validation that prevents the modification of entities managed by other multi-target apps", false},
		allModulesOpt:                 {new(bool), false, "Deploy all modules which are contained in the deployment descriptor, in the current location", false},
		allResourcesOpt:               {new(bool), false, "Deploy all resources which are contained in the deployment descriptor, in the current location", false},
		verifyArchiveSignatureOpt:     {new(bool), false, "Verify the archive is correctly signed", false},
		retriesOpt:                    {new(uint), 3, "Retry the operation N times in case a non-content error occurs (default 3)", false},
		strategyOpt:                   {new(string), "", "Specify the deployment strategy when updating an mta (default, blue-green)", false},
		noConfirmOpt:                  {new(bool), false, "Do not require confirmation for deleting the previously deployed MTA apps (only applicable when using blue-green deployment)", false},
	}
}

// GetPluginCommand returns the plugin command details
func (c *DeployCommand) GetPluginCommand() plugin.Command {
	return plugin.Command{
		Name:     "deploy",
		HelpText: "Deploy a new multi-target app or sync changes to an existing one",
		UsageDetails: plugin.Usage{
			Usage: `Deploy a multi-target app archive
   cf deploy MTA [-e EXT_DESCRIPTOR[,...]] [-t TIMEOUT] [--version-rule VERSION_RULE] [-u URL] [-f] [--retries RETRIES] [--no-start] [--use-namespaces] [--no-namespaces-for-services] [--delete-services] [--delete-service-keys] [--delete-service-brokers] [--keep-files] [--no-restart-subscribed-apps] [--do-not-fail-on-missing-permissions] [--abort-on-error] [--skip-ownership-validation] [--verify-archive-signature] [--strategy STRATEGY] [--no-confirm]

   Perform action on an active deploy operation
   cf deploy -i OPERATION_ID -a ACTION [-u URL]`,
			Options: c.getOptionsForPluginCommand(),
		},
	}
}

// DeployProcessParametersSetter returns a new ProcessParametersSetter.
func deployProcessParametersSetter() ProcessParametersSetter {
	return func(options map[string]CommandOption, processBuilder *util.ProcessBuilder) {
		processBuilder.Parameter("deleteServiceKeys", strconv.FormatBool(getBoolOpt(deleteServiceKeysOpt, options)))
		processBuilder.Parameter("deleteServices", strconv.FormatBool(getBoolOpt(deleteServicesOpt, options)))
		processBuilder.Parameter("noStart", strconv.FormatBool(getBoolOpt(noStartOpt, options)))
		processBuilder.Parameter("useNamespaces", strconv.FormatBool(getBoolOpt(useNamespacesOpt, options)))
		processBuilder.Parameter("useNamespacesForServices", strconv.FormatBool(!getBoolOpt(noNamespacesForServicesOpt, options)))
		processBuilder.Parameter("deleteServiceBrokers", strconv.FormatBool(getBoolOpt(deleteServiceBrokersOpt, options)))
		processBuilder.Parameter("startTimeout", getStringOpt(timeoutOpt, options))
		processBuilder.Parameter("versionRule", getStringOpt(versionRuleOpt, options))
		processBuilder.Parameter("keepFiles", strconv.FormatBool(getBoolOpt(keepFilesOpt, options)))
		processBuilder.Parameter("noRestartSubscribedApps", strconv.FormatBool(getBoolOpt(noRestartSubscribedAppsOpt, options)))
		processBuilder.Parameter("noFailOnMissingPermissions", strconv.FormatBool(getBoolOpt(noFailOnMissingPermissionsOpt, options)))
		processBuilder.Parameter("abortOnError", strconv.FormatBool(getBoolOpt(abortOnErrorOpt, options)))
		processBuilder.Parameter("skipOwnershipValidation", strconv.FormatBool(getBoolOpt(skipOwnershipValidationOpt, options)))
		processBuilder.Parameter("verifyArchiveSignature", strconv.FormatBool(getBoolOpt(verifyArchiveSignatureOpt, options)))
	}
}

type DeployCommandOptionParser struct {
	AbstractOptionParser
}

func NewDeployCommandOptionParser() DeployCommandOptionParser {
	return DeployCommandOptionParser{AbstractOptionParser{}}
}

func (DeployCommandOptionParser) additionalParse(name string, option CommandOption, flags *flag.FlagSet) {
	if val, ok := option.Value.(*listFlag); ok {
		flags.Var(val, name, "")
	}
}

// Execute executes the command
func (c *DeployCommand) Execute(args []string) ExecutionStatus {
	log.Tracef("Executing command '" + c.name + "': args: '%v'\n", args)

	mtaArgument, pos := getMtaArgumentAndPosition(c.flags, args)
	parser := NewCommandFlagsParserWithValidator(c.flags, NewProcessActionExecutorCommandArgumentsParser(pos), &deployCommandFlagsValidator{})
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

	operationID := getStringOpt(operationIDOpt, c.options)
	action := getStringOpt(actionOpt, c.options)
	retries := getUintOpt(retriesOpt, c.options)

	if operationID != "" || action != "" {
		return c.ExecuteAction(operationID, action, retries, host)
	}

	mtaElementsCalculator := mtaElementsToAddCalculator{shouldAddAllModules: false, shouldAddAllResources: false}
	mtaElementsCalculator.calculateElementsToDeploy(c.options)

	mtaArchive, err := getMtaArchive(mtaArgument, mtaElementsCalculator)
	if err != nil {
		ui.Failed("Error retrieving MTA: %s", err.Error())
		return Failure
	}

	context, err := c.GetContext()
	if err != nil {
		ui.Failed(err.Error())
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

	extDescriptors := getStringOpt(extDescriptorsOpt, c.options)

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

	// Extract mta id from archive file
	mtaID, err := util.GetMtaIDFromArchive(mtaArchivePath)
	if os.IsNotExist(err) {
		ui.Failed("Could not find file %s", terminal.EntityNameColor(mtaArchivePath))
		return Failure
	} else if err != nil {
		ui.Failed("Could not get MTA id from deployment descriptor: %s", err)
		return Failure
	}

	force := getBoolOpt(forceOpt, c.options)

	// Check for an ongoing operation for this MTA ID and abort it
	wasAborted, err := c.CheckOngoingOperation(mtaID, force, host)
	if err != nil {
		ui.Failed("Could not get MTA operations: %s", baseclient.NewClientError(err))
		return Failure
	}
	if !wasAborted {
		return Failure
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
	processBuilder := NewDeploymentStrategy(c.options, c.processTypeProvider).CreateProcessBuilder()
	c.processParametersSetter(c.options, processBuilder)
	processBuilder.Parameter("appArchiveId", strings.Join(uploadedArchivePartIds, ","))
	processBuilder.Parameter("mtaExtDescriptorId", strings.Join(uploadedExtDescriptorIDs, ","))
	processBuilder.Parameter("mtaId", mtaID)
	setModulesAndResourcesListParameters(modulesList, resourcesList, processBuilder, mtaElementsCalculator)

	operation := processBuilder.Build()

	// Create the new process
	responseHeader, err := mtaClient.StartMtaOperation(*operation)
	if err != nil {
		ui.Failed("Could not create operation: %s", baseclient.NewClientError(err))
		return Failure
	}

	return NewExecutionMonitorFromLocationHeader(c.name, responseHeader.Location.String(), retries, mtaClient).Monitor()
}

func getMtaArgumentAndPosition(flags *flag.FlagSet, args []string) (string, int) {
	if len(args) == 0 {
		return "", 0
	}

	if flags.Lookup(strings.Replace(args[0], "-", "", 2)) == nil {
		return args[0], 1
	}
	return "", 0
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

func getMtaArchive(mtaArgument string, mtaElementsCalculator mtaElementsToAddCalculator) (string, error) {
	if mtaArgument == "" {
		currentWorkingDirectory, err := os.Getwd()
		if err != nil {
			return "", fmt.Errorf("Could not get the current working directory: %s", err.Error())
		}
		return buildMtaArchiveFromDirectory(currentWorkingDirectory, mtaElementsCalculator)
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

func (c *mtaElementsToAddCalculator) calculateElementsToDeploy(options map[string]CommandOption) {
	allModulesSpecified := getBoolOpt(allModulesOpt, options)
	allResourcesSpecified := getBoolOpt(allResourcesOpt, options)

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

type deployCommandFlagsValidator struct {
}

func (d *deployCommandFlagsValidator) ValidateFlags(flags *flag.FlagSet, _ []string) error {
	var err error
	flags.Visit(func(f *flag.Flag) {
		if f.Name == strategyOpt {
			if f.Value.String() == "" {
				err = errors.New("strategy flag defined but no argument passed")
			} else if !isContainedIn(f.Value.String(), AvailableStrategies()) {
				err = fmt.Errorf("%s is not a valid deployment strategy\nAvailable strategies %v", f.Value.String(), AvailableStrategies())
			}
		}
	})
	return err
}

func (*deployCommandFlagsValidator) IsBeforeParsing() bool {
	return false
}

func isContainedIn(s string, arr []string) bool {
	for _, el := range arr {
		if el == s {
			return true
		}
	}
	return false
}
