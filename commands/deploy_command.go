package commands

import (
	"bufio"
	"encoding/base64"
	"errors"
	"flag"
	"fmt"
	"io/fs"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"code.cloudfoundry.org/cli/v8/cf/terminal"
	"code.cloudfoundry.org/cli/v8/plugin"
	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/clients/baseclient"
	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/clients/cfrestclient"
	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/clients/cfrestclient/resilient"
	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/clients/models"
	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/clients/mtaclient"
	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/commands/retrier"
	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/configuration"
	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/log"
	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/secure_parameters"
	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/ui"
	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/util"
	"gopkg.in/cheggaaa/pb.v1"
)

const (
	extDescriptorsOpt                = "e"
	timeoutOpt                       = "t"
	versionRuleOpt                   = "version-rule"
	noStartOpt                       = "no-start"
	deleteServiceKeysOpt             = "delete-service-keys"
	keepFilesOpt                     = "keep-files"
	skipOwnershipValidationOpt       = "skip-ownership-validation"
	moduleOpt                        = "m"
	resourceOpt                      = "r"
	allModulesOpt                    = "all-modules"
	allResourcesOpt                  = "all-resources"
	strategyOpt                      = "strategy"
	skipTestingPhase                 = "skip-testing-phase"
	skipIdleStart                    = "skip-idle-start"
	startTimeoutOpt                  = "apps-start-timeout"
	stageTimeoutOpt                  = "apps-stage-timeout"
	uploadTimeoutOpt                 = "apps-upload-timeout"
	taskExecutionTimeoutOpt          = "apps-task-execution-timeout"
	applyNamespaceAppNamesOpt        = "apply-namespace-app-names"
	applyNamespaceServiceNamesOpt    = "apply-namespace-service-names"
	applyNamespaceAppRoutesOpt       = "apply-namespace-app-routes"
	applyNamespaceAsSuffix           = "apply-namespace-as-suffix"
	maxNamespaceSize                 = 36
	shouldBackupPreviousVersionOpt   = "backup-previous-version"
	requireSecureParameters          = "require-secure-parameters"
	disposableUserProvidedServiceOpt = "disposable-user-provided-service"
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

	FileUrlReader      fs.File
	FileUrlReadTimeout time.Duration
	CfClient           cfrestclient.CloudFoundryOperationsExtended
}

// NewDeployCommand creates a new deploy command.
func NewDeployCommand() *DeployCommand {
	baseCmd := &BaseCommand{flagsParser: deployCommandLineArgumentsParser{}, flagsValidator: deployCommandFlagsValidator{}}
	deployCmd := &DeployCommand{baseCmd, deployProcessParametersSetter(), &deployCommandProcessTypeProvider{}, os.Stdin, 30 * time.Second, nil}
	baseCmd.Command = deployCmd
	return deployCmd
}

func (c *DeployCommand) Initialize(name string, cliConnection plugin.CliConnection) {
	c.BaseCommand.Initialize(name, cliConnection)
	delegate := cfrestclient.NewCloudFoundryRestClient(cliConnection)
	c.CfClient = resilient.NewResilientCloudFoundryClient(delegate, maxRetriesCount, retryIntervalInSeconds)
}

// GetPluginCommand returns the plugin command details
func (c *DeployCommand) GetPluginCommand() plugin.Command {
	return plugin.Command{
		Name:     "deploy",
		HelpText: "Deploy a new multi-target app or sync changes to an existing one",
		UsageDetails: plugin.Usage{
			Usage: `Deploy a multi-target app archive

   cf deploy MTA [-e EXT_DESCRIPTOR[,...]] [-t TIMEOUT] [--version-rule VERSION_RULE] [-u URL] [-f] [--retries RETRIES] [--no-start] [--namespace NAMESPACE] [--apply-namespace-app-names true/false] [--apply-namespace-service-names true/false] [--apply-namespace-app-routes true/false] [--apply-namespace-as-suffix true/false ] [--delete-services] [--delete-service-keys] [--delete-service-brokers] [--keep-files] [--no-restart-subscribed-apps] [--do-not-fail-on-missing-permissions] [--abort-on-error] [--strategy STRATEGY] [--skip-testing-phase] [--skip-idle-start] [--require-secure-parameters] [--disposable-user-provided-service] [--apps-start-timeout TIMEOUT] [--apps-stage-timeout TIMEOUT] [--apps-upload-timeout TIMEOUT] [--apps-task-execution-timeout TIMEOUT]

   Perform action on an active deploy operation
   cf deploy -i OPERATION_ID -a ACTION [-u URL]

   Deploy a multi-target app archive referenced by a remote URL
   <write MTA archive URL to STDOUT> | cf deploy [-e EXT_DESCRIPTOR[,...]] [-t TIMEOUT] [--version-rule VERSION_RULE] [-u MTA_CONTROLLER_URL] [--retries RETRIES] [--no-start] [--namespace NAMESPACE] [--apply-namespace-app-names true/false] [--apply-namespace-service-names true/false] [--apply-namespace-app-routes true/false] [--apply-namespace-as-suffix true/false ] [--delete-services] [--delete-service-keys] [--delete-service-brokers] [--keep-files] [--no-restart-subscribed-apps] [--do-not-fail-on-missing-permissions] [--abort-on-error] [--strategy STRATEGY] [--skip-testing-phase] [--skip-idle-start] [require-secure-parameters] [--disposable-user-provided-service] [--apps-start-timeout TIMEOUT] [--apps-stage-timeout TIMEOUT] [--apps-upload-timeout TIMEOUT] [--apps-task-execution-timeout TIMEOUT]` + util.UploadEnvHelpText,

			Options: map[string]string{
				extDescriptorsOpt:                 "Extension descriptors",
				deployServiceURLOpt:               "Deploy service URL, by default 'deploy-service.<system-domain>'",
				versionRuleOpt:                    "Version rule (HIGHER, SAME_HIGHER, ALL)",
				operationIDOpt:                    "Active deploy operation ID",
				actionOpt:                         "Action to perform on active deploy operation (abort, retry, resume, monitor)",
				forceOpt:                          "Force deploy without confirmation for aborting conflicting processes",
				moduleOpt:                         "Deploy list of modules which are contained in the deployment descriptor, in the current location",
				resourceOpt:                       "Deploy list of resources which are contained in the deployment descriptor, in the current location",
				util.GetShortOption(noStartOpt):   "Do not start apps",
				util.GetShortOption(namespaceOpt): "Namespace for the MTA, applied on app names, app routes and service names",
				util.GetShortOption(applyNamespaceAppNamesOpt):                  "Apply namespace to application names: (true, false)",
				util.GetShortOption(applyNamespaceServiceNamesOpt):              "Apply namespace to service names: (true, false)",
				util.GetShortOption(applyNamespaceAppRoutesOpt):                 "Apply namespace to application routes: (true, false)",
				util.GetShortOption(applyNamespaceAsSuffix):                     "Apply namespace as a suffix rather than a prefix: (true, false)",
				util.GetShortOption(deleteServicesOpt):                          "Recreate changed services / delete discontinued services",
				util.GetShortOption(deleteServiceKeysOpt):                       "Delete existing service keys and apply the new ones",
				util.GetShortOption(deleteServiceBrokersOpt):                    "Delete discontinued service brokers",
				util.GetShortOption(keepFilesOpt):                               "Keep files used for deployment",
				util.GetShortOption(noRestartSubscribedAppsOpt):                 "Do not restart subscribed apps, updated during the deployment",
				util.GetShortOption(noFailOnMissingPermissionsOpt):              "Do not fail on missing permissions for admin operations",
				util.GetShortOption(abortOnErrorOpt):                            "Auto-abort the process on any errors",
				util.GetShortOption(allModulesOpt):                              "Deploy all modules which are contained in the deployment descriptor, in the current location",
				util.GetShortOption(allResourcesOpt):                            "Deploy all resources which are contained in the deployment descriptor, in the current location",
				util.GetShortOption(retriesOpt):                                 "Retry the operation N times in case a non-content error occurs (default 3)",
				util.GetShortOption(strategyOpt):                                "Specify the deployment strategy when updating an mta (default, blue-green, incremental-blue-green)",
				util.GetShortOption(skipTestingPhase):                           "(STRATEGY: BLUE-GREEN, INCREMENTAL-BLUE-GREEN) Do not require confirmation for deleting the previously deployed MTA app",
				util.GetShortOption(skipIdleStart):                              "(STRATEGY: BLUE-GREEN, INCREMENTAL-BLUE-GREEN) Directly start the new MTA version as 'live', skipping the 'idle' phase of the resources. Do not require further confirmation or testing before deleting the old version",
				util.GetShortOption(stageTimeoutOpt):                            "Stage app timeout in seconds",
				util.GetShortOption(uploadTimeoutOpt):                           "Upload app timeout in seconds",
				util.GetShortOption(taskExecutionTimeoutOpt):                    "Task execution timeout in seconds",
				util.CombineFullAndShortParameters(startTimeoutOpt, timeoutOpt): "Start app timeout in seconds",
				util.GetShortOption(shouldBackupPreviousVersionOpt):             "(EXPERIMENTAL) (STRATEGY: BLUE-GREEN, INCREMENTAL-BLUE-GREEN) Backup previous version of applications, use new cli command \"rollback-mta\" to rollback to the previous version",
				util.GetShortOption(dependencyAwareStopOrderOpt):                "(EXPERIMENTAL) (STRATEGY: BLUE-GREEN, INCREMENTAL-BLUE-GREEN) Stop apps in a dependency-aware order during the resume phase of a blue-green deployment",
				util.GetShortOption(requireSecureParameters):                    "(EXPERIMENTAL) Pass secrets to the deploy service in a secure way",
				util.GetShortOption(disposableUserProvidedServiceOpt):           "Deploy when --require-secure-parameters flag is active for disposable UPS to be created and then deleted at the of the operation",
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
		processBuilder.Parameter("versionRule", GetStringOpt(versionRuleOpt, flags))
		processBuilder.Parameter("keepFiles", strconv.FormatBool(GetBoolOpt(keepFilesOpt, flags)))
		processBuilder.Parameter("noRestartSubscribedApps", strconv.FormatBool(GetBoolOpt(noRestartSubscribedAppsOpt, flags)))
		processBuilder.Parameter("noFailOnMissingPermissions", strconv.FormatBool(GetBoolOpt(noFailOnMissingPermissionsOpt, flags)))
		processBuilder.Parameter("abortOnError", strconv.FormatBool(GetBoolOpt(abortOnErrorOpt, flags)))
		processBuilder.Parameter("skipOwnershipValidation", strconv.FormatBool(GetBoolOpt(skipOwnershipValidationOpt, flags)))
		processBuilder.Parameter("appsStageTimeout", GetStringOpt(stageTimeoutOpt, flags))
		processBuilder.Parameter("appsUploadTimeout", GetStringOpt(uploadTimeoutOpt, flags))
		processBuilder.Parameter("appsTaskExecutionTimeout", GetStringOpt(taskExecutionTimeoutOpt, flags))
		processBuilder.Parameter("isSecurityEnabled", strconv.FormatBool(GetBoolOpt(requireSecureParameters, flags)))
		processBuilder.Parameter("isDisposableUserProvidedServiceEnabled", strconv.FormatBool(GetBoolOpt(disposableUserProvidedServiceOpt, flags)))

		var lastSetValue string = ""
		for i := 0; i < len(os.Args); i++ {
			arg := os.Args[i]
			if arg == "-t" {
				if i+1 < len(os.Args) {
					lastSetValue = os.Args[i+1]
					i++
				}
			} else if arg == "--apps-start-timeout" {
				if i+1 < len(os.Args) {
					lastSetValue = os.Args[i+1]
					i++
				}
			}
		}
		processBuilder.Parameter("startTimeout", lastSetValue)
	}
}

func (c *DeployCommand) defineCommandOptions(flags *flag.FlagSet) {
	flags.String(extDescriptorsOpt, "", "")
	flags.String(operationIDOpt, "", "")
	flags.String(actionOpt, "", "")
	flags.Bool(forceOpt, false, "")
	flags.String(versionRuleOpt, "", "")
	flags.Bool(deleteServicesOpt, false, "")
	flags.Bool(noStartOpt, false, "")
	flags.String(namespaceOpt, "", "")
	flags.String(applyNamespaceAppNamesOpt, "", "")
	flags.String(applyNamespaceServiceNamesOpt, "", "")
	flags.String(applyNamespaceAppRoutesOpt, "", "")
	flags.String(applyNamespaceAsSuffix, "", "")
	flags.Bool(deleteServiceKeysOpt, false, "")
	flags.Bool(deleteServiceBrokersOpt, false, "")
	flags.Bool(keepFilesOpt, false, "")
	flags.Bool(noRestartSubscribedAppsOpt, false, "")
	flags.Bool(noFailOnMissingPermissionsOpt, false, "")
	flags.Bool(abortOnErrorOpt, false, "")
	flags.Bool(skipOwnershipValidationOpt, false, "")
	flags.Bool(allModulesOpt, false, "")
	flags.Bool(allResourcesOpt, false, "")
	flags.Uint(retriesOpt, 3, "")
	flags.String(strategyOpt, "default", "")
	flags.Bool(skipTestingPhase, false, "")
	flags.Bool(skipIdleStart, false, "")
	flags.Var(&modulesList, moduleOpt, "")
	flags.Var(&resourcesList, resourceOpt, "")
	flags.String(timeoutOpt, "", "")
	flags.String(startTimeoutOpt, "", "")
	flags.String(stageTimeoutOpt, "", "")
	flags.String(uploadTimeoutOpt, "", "")
	flags.String(taskExecutionTimeoutOpt, "", "")
	flags.Bool(shouldBackupPreviousVersionOpt, false, "")
	flags.Bool(dependencyAwareStopOrderOpt, false, "")
	flags.Bool(requireSecureParameters, false, "")
	flags.Bool(disposableUserProvidedServiceOpt, false, "")
}

func (c *DeployCommand) executeInternal(positionalArgs []string, dsHost string, flags *flag.FlagSet, cfTarget util.CloudFoundryTarget) ExecutionStatus {
	operationID := GetStringOpt(operationIDOpt, flags)
	action := GetStringOpt(actionOpt, flags)
	retries := GetUintOpt(retriesOpt, flags)

	if operationID != "" || action != "" {
		return c.ExecuteAction(operationID, action, retries, dsHost, cfTarget)
	}

	mtaElementsCalculator := createMtaElementsCalculator(flags)

	rawMtaArchive, err := c.getMtaArchive(positionalArgs, mtaElementsCalculator)
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
	var disposableUserProvidedServiceName string

	// Check SLMP metadata
	// TODO: ensure session
	mtaClient := c.NewMtaClient(dsHost, cfTarget)

	namespace := strings.TrimSpace(GetStringOpt(namespaceOpt, flags))
	force := GetBoolOpt(forceOpt, flags)
	conf := configuration.NewSnapshot()
	uploadChunkSize := conf.GetUploadChunkSizeInMB()
	sequentialUpload := conf.GetUploadChunksSequentially()
	disableProgressBar := conf.GetDisableUploadProgressBar()
	fileUploader := NewFileUploader(mtaClient, namespace, uploadChunkSize, sequentialUpload, disableProgressBar)
	var yamlBytes []byte

	if isUrl {
		var fileId string
		var schemaVersion string

		asyncUploadJobResult := c.uploadFromUrl(mtaArchive, mtaClient, namespace, disableProgressBar)
		if asyncUploadJobResult.ExecutionStatus == Failure {
			return Failure
		}
		mtaId, fileId, schemaVersion = asyncUploadJobResult.MtaId, asyncUploadJobResult.FileId, asyncUploadJobResult.SchemaVersion
		// Check for an ongoing operation for this MTA ID and abort it
		wasAborted, err := c.CheckOngoingOperation(mtaId, namespace, dsHost, force, cfTarget)
		if err != nil {
			ui.Failed("Could not get MTA operations: %s", baseclient.NewClientError(err))
			return Failure
		}
		if !wasAborted {
			return Failure
		}

		uploadedArchivePartIds = append(uploadedArchivePartIds, fileId)

		if GetBoolOpt(requireSecureParameters, flags) {
			result := setUpSpecificsForDeploymentUsingSecrets(flags, c, mtaId, namespace, schemaVersion, &disposableUserProvidedServiceName, &yamlBytes)
			if result != Success {
				return Failure
			}
		}

		ui.Ok()
	} else {
		// Get the full path of the MTA archive
		mtaArchivePath, err := filepath.Abs(mtaArchive)
		if err != nil {
			ui.Failed("Could not get absolute path of file %q", mtaArchive)
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

		// Check for an ongoing operation for this MTA ID and abort it
		wasAborted, err := c.CheckOngoingOperation(mtaId, namespace, dsHost, force, cfTarget)
		if err != nil {
			ui.Failed("Could not get MTA operations: %s", baseclient.NewClientError(err))
			return Failure
		}
		if !wasAborted {
			return Failure
		}

		if GetBoolOpt(requireSecureParameters, flags) {
			result := setUpSpecificsForDeploymentUsingSecrets(flags, c, mtaId, namespace, descriptor.SchemaVersion, &disposableUserProvidedServiceName, &yamlBytes)
			if result != Success {
				return Failure
			}
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
				ui.Failed("Could not get absolute path of file %q", extDescriptorFile)
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

	if GetBoolOpt(requireSecureParameters, flags) {
		secureFileID, err := fileUploader.UploadBytes("__mta.secure.mtaext", yamlBytes)
		if err != nil {
			ui.Failed("Could not upload secure extension: %s", err)
			return Failure
		}
		uploadedExtDescriptorIDs = append(uploadedExtDescriptorIDs, secureFileID)
	}

	// Build the process instance
	processBuilder := NewDeploymentStrategy(flags, c.processTypeProvider).CreateProcessBuilder()
	processBuilder.Namespace(namespace)
	processBuilder.Parameter("applyNamespaceAppNames", GetStringOpt(applyNamespaceAppNamesOpt, flags))
	processBuilder.Parameter("applyNamespaceServiceNames", GetStringOpt(applyNamespaceServiceNamesOpt, flags))
	processBuilder.Parameter("applyNamespaceAppRoutes", GetStringOpt(applyNamespaceAppRoutesOpt, flags))
	processBuilder.Parameter("applyNamespaceAsSuffix", GetStringOpt(applyNamespaceAsSuffix, flags))

	processBuilder.Parameter("appArchiveId", strings.Join(uploadedArchivePartIds, ","))
	processBuilder.Parameter("mtaExtDescriptorId", strings.Join(uploadedExtDescriptorIDs, ","))
	processBuilder.Parameter("mtaId", mtaId)
	processBuilder.Parameter("disposableUserProvidedServiceName", disposableUserProvidedServiceName)
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
	return executionMonitor.Monitor()
}

func setUpSpecificsForDeploymentUsingSecrets(flags *flag.FlagSet, c *DeployCommand, mtaId, namespace, schemaVersion string, disposableUserProvidedServiceName *string, yamlBytes *[]byte) ExecutionStatus {
	// Collect special ENVs: __MTA___<name>, __MTA_JSON___<name>, __MTA_CERT___<name>
	parameters, err := secure_parameters.CollectFromEnv("__MTA")
	if err != nil {
		ui.Failed("Secure parameters error: %s", err)
		return Failure
	}

	if len(parameters) == 0 {
		ui.Failed("No secure parameters found in environment. Set variables like __MTA___<name>, __MTA_JSON___<name>, or __MTA_CERT___<name>.")
		return Failure
	}

	if GetBoolOpt(disposableUserProvidedServiceOpt, flags) {
		disposableUserProvidedServiceNameResult, err := secure_parameters.GetRandomisedUpsName(mtaId, namespace)
		if err != nil {
			ui.Failed("Failed to create disposable user-provided service name: %v", err)
			return Failure
		}

		isDisposableUpsCreated, _, err := secure_parameters.CreateDisposableUps(disposableUserProvidedServiceNameResult, c.cliConnection, c.CfClient)
		if err != nil {
			ui.Failed("Could not ensure disposable user-provided service %s: %v", disposableUserProvidedServiceName, err)
			return Failure
		}

		*disposableUserProvidedServiceName = disposableUserProvidedServiceNameResult
		if isDisposableUpsCreated {
			ui.Say("Created disposable user-provided service %s for secure parameters. Will be automatically deleted at the end of the operation!", terminal.EntityNameColor(disposableUserProvidedServiceNameResult))
		}
	} else {
		userProvidedServiceName := secure_parameters.GetUpsName(mtaId, namespace)

		isUpsCreated, _, err := secure_parameters.ValidateUpsExistsOrElseCreateIt(userProvidedServiceName, c.cliConnection, c.CfClient)
		if err != nil {
			ui.Failed("Could not ensure user-provided service %s: %v", userProvidedServiceName, err)
			return Failure
		}

		*disposableUserProvidedServiceName = ""
		if isUpsCreated {
			ui.Say("Created user-provided service %s for secure parameters.", terminal.EntityNameColor(userProvidedServiceName))
		} else {
			ui.Say("Using existing user-provided service %s for secure parameters.", terminal.EntityNameColor(userProvidedServiceName))
		}
	}

	yamlBytesResult, err := secure_parameters.BuildSecureExtension(parameters, mtaId, schemaVersion)
	if err != nil {
		ui.Failed("Could not build secure extension: %s", err)
		return Failure
	}
	if len(yamlBytesResult) == 0 {
		ui.Failed("Secure extension descriptor is empty: %s", err)
		return Failure
	}
	*yamlBytes = yamlBytesResult

	return Success
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

func (c *DeployCommand) uploadFromUrl(url string, mtaClient mtaclient.MtaClientOperations, namespace string,
	disableProgressBar bool) UploadFromUrlStatus {
	encodedFileUrl := base64.URLEncoding.EncodeToString([]byte(url))
	uploadStatus, _ := retrier.Execute[UploadFromUrlStatus](3, func() (UploadFromUrlStatus, error) {
		progressBar := c.tryFetchMtarSize(url, disableProgressBar)
		uploadFromUrlStatus := c.doUploadFromUrl(encodedFileUrl, mtaClient, namespace, progressBar)
		return uploadFromUrlStatus, nil
	}, func(result UploadFromUrlStatus, err error) bool {
		return shouldRetryUpload(result)
	})
	return uploadStatus
}

func (c *DeployCommand) doUploadFromUrl(encodedFileUrl string, mtaClient mtaclient.MtaClientOperations, namespace string, progressBar *pb.ProgressBar) UploadFromUrlStatus {
	responseHeaders, err := mtaClient.StartUploadMtaArchiveFromUrl(encodedFileUrl, &namespace)
	if err != nil {
		ui.Failed("Could not upload from url: %s", err)
		return UploadFromUrlStatus{
			FileId:          "",
			MtaId:           "",
			SchemaVersion:   "",
			ClientActions:   make([]string, 0),
			ExecutionStatus: Failure,
		}
	}

	var totalBytesProcessed int64 = 0
	if progressBar != nil {
		progressBar.Start()
		defer progressBar.Finish()
	}

	uploadJobUrl := responseHeaders.Get("Location")
	jobUrlParts := strings.Split(uploadJobUrl, "/")
	jobId := jobUrlParts[len(jobUrlParts)-1]

	timeout := time.NewTimer(time.Hour)
	defer timeout.Stop()
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	var file *models.FileMetadata
	var jobResult mtaclient.AsyncUploadJobResult
	for file == nil {
		jobResult, err = mtaClient.GetAsyncUploadJob(jobId, &namespace, responseHeaders.Get("x-cf-app-instance"))
		if err != nil {
			ui.Failed("Could not upload from url: %s", err)
			return UploadFromUrlStatus{
				FileId:          "",
				MtaId:           "",
				SchemaVersion:   "",
				ClientActions:   jobResult.ClientActions,
				ExecutionStatus: Failure,
			}
		}
		file = jobResult.File
		if len(jobResult.Error) != 0 {
			ui.Failed("Async upload job failed: %s", jobResult.Error)
			return UploadFromUrlStatus{
				FileId:          "",
				MtaId:           "",
				SchemaVersion:   "",
				ClientActions:   jobResult.ClientActions,
				ExecutionStatus: Failure,
			}
		}

		if progressBar != nil && jobResult.BytesProcessed != -1 {
			if jobResult.BytesProcessed < totalBytesProcessed {
				//retry happened in backend, rewind the progress bar
				progressBar.Add64(-totalBytesProcessed + jobResult.BytesProcessed)
			} else {
				progressBar.Add64(jobResult.BytesProcessed - totalBytesProcessed)
			}
			totalBytesProcessed = jobResult.BytesProcessed
		}

		if len(jobResult.MtaId) == 0 {
			select {
			case <-timeout.C:
				ui.Failed("Upload from URL timed out after 1 hour")
				return UploadFromUrlStatus{
					FileId:          "",
					MtaId:           "",
					SchemaVersion:   "",
					ClientActions:   make([]string, 0),
					ExecutionStatus: Failure,
				}
			case <-ticker.C:
			}
		}
	}
	if progressBar != nil && totalBytesProcessed < progressBar.Total {
		progressBar.Add64(progressBar.Total - totalBytesProcessed)
	}
	return UploadFromUrlStatus{
		FileId:          file.ID,
		MtaId:           jobResult.MtaId,
		SchemaVersion:   jobResult.SchemaVersion,
		ClientActions:   jobResult.ClientActions,
		ExecutionStatus: Success,
	}
}

func shouldRetryUpload(uploadFromUrlStatus UploadFromUrlStatus) bool {
	if uploadFromUrlStatus.ExecutionStatus == Success {
		return false
	}
	for _, clientAction := range uploadFromUrlStatus.ClientActions {
		if clientAction == "RETRY_UPLOAD" {
			ui.Warn("Upload request must be retried")
			return true
		}
	}
	return false
}

func (c *DeployCommand) uploadFiles(files []string, fileUploader *FileUploader) ([]string, ExecutionStatus) {
	var resultIds []string
	if len(files) == 0 {
		return resultIds, Success
	}

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
	if !mtaElementsCalculator.shouldAddAllModules {
		processBuilder.SetParameterWithoutCheck("modulesForDeployment", modulesList.getProcessList())
	}
	if !mtaElementsCalculator.shouldAddAllResources {
		processBuilder.SetParameterWithoutCheck("resourcesForDeployment", resourcesList.getProcessList())
	}
}

func (c *DeployCommand) getMtaArchive(parsedArguments []string, mtaElementsCalculator mtaElementsToAddCalculator) (interface{}, error) {
	if len(parsedArguments) == 0 {
		fileUrl := c.tryReadingFileUrl()
		if len(fileUrl) > 0 {
			return url.Parse(fileUrl)
		}

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

func (c *DeployCommand) tryReadingFileUrl() string {
	stat, err := c.FileUrlReader.Stat()
	if err != nil {
		return ""
	}

	if stat.Mode()&fs.ModeCharDevice == 0 {
		in := bufio.NewReader(c.FileUrlReader)
		input, _ := in.ReadString('\n')
		return strings.TrimSpace(input)
	}
	return ""
}

func (c *DeployCommand) tryFetchMtarSize(url string, disableProgressBar bool) *pb.ProgressBar {
	client := http.Client{Timeout: c.FileUrlReadTimeout}
	resp, err := client.Head(url)
	if err != nil {
		log.Tracef("could not call remote MTAR endpoint: %v\n", err)
		return nil
	}
	if resp.StatusCode/100 != 2 {
		log.Tracef("could not fetch remote MTAR size: %s\n", resp.Status)
		return nil
	}
	bar := pb.New64(resp.ContentLength).SetUnits(pb.U_BYTES)
	bar.ShowElapsedTime = true
	bar.ShowTimeLeft = false
	bar.NotPrint = disableProgressBar
	return bar
}

func buildMtaArchiveFromDirectory(mtaDirectoryLocation string, mtaElementsCalculator mtaElementsToAddCalculator) (string, error) {
	deploymentDescriptor, _, err := util.ParseDeploymentDescriptor(mtaDirectoryLocation)
	if err != nil {
		return "", err
	}

	modulesToAdd := mtaElementsCalculator.getModulesToAdd(deploymentDescriptor)
	resourcesToAdd := mtaElementsCalculator.getResourcesToAdd(deploymentDescriptor)

	return util.NewMtaArchiveBuilder(modulesToAdd, resourcesToAdd).Build(mtaDirectoryLocation)
}

type mtaElementsToAddCalculator struct {
	shouldAddAllModules   bool
	shouldAddAllResources bool
}

func createMtaElementsCalculator(flags *flag.FlagSet) mtaElementsToAddCalculator {
	return mtaElementsToAddCalculator{
		shouldAddAllModules:   GetBoolOpt(allModulesOpt, flags) || len(modulesList.getElements()) == 0,
		shouldAddAllResources: GetBoolOpt(allResourcesOpt, flags) || len(resourcesList.getElements()) == 0,
	}
}

func (c mtaElementsToAddCalculator) getModulesToAdd(deploymentDescriptor util.MtaDeploymentDescriptor) []string {
	if c.shouldAddAllModules {
		modulesToAdd := make([]string, 0)
		for _, module := range deploymentDescriptor.Modules {
			modulesToAdd = append(modulesToAdd, module.Name)
		}
		return modulesToAdd
	}

	return modulesList.getElements()
}

func (c mtaElementsToAddCalculator) getResourcesToAdd(deploymentDescriptor util.MtaDeploymentDescriptor) []string {
	if c.shouldAddAllResources {
		resourcesToAdd := make([]string, 0)
		for _, resource := range deploymentDescriptor.Resources {
			resourcesToAdd = append(resourcesToAdd, resource.Name)
		}
		return resourcesToAdd
	}

	return resourcesList.getElements()
}

type deployCommandProcessTypeProvider struct{}

func (d deployCommandProcessTypeProvider) GetProcessType() string {
	return "DEPLOY"
}

type deployCommandLineArgumentsParser struct{}

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
		switch f.Name {
		case strategyOpt:
			if f.Value.String() == "" {
				err = errors.New("strategy flag defined but no argument specified")
				return
			} else if !util.Contains(AvailableStrategies(), f.Value.String()) {
				err = fmt.Errorf("%s is not a valid deployment strategy, available strategies: %v", f.Value.String(), AvailableStrategies())
				return
			}
		case namespaceOpt:
			if len(f.Value.String()) > maxNamespaceSize {
				err = fmt.Errorf("Invalid value for namespace. The namespace cannot be more than %d symbols.", maxNamespaceSize)
				return
			}
		case timeoutOpt, startTimeoutOpt, stageTimeoutOpt, uploadTimeoutOpt, taskExecutionTimeoutOpt:
			if e := ValidateTimeoutOption(f.Name, flags, 259200); e != nil {
				err = e
				return
			}
		case applyNamespaceAppNamesOpt, applyNamespaceServiceNamesOpt, applyNamespaceAppRoutesOpt, applyNamespaceAsSuffix:
			if e := ValidateBooleanFlag(f.Name, flags); e != nil {
				err = e
				return
			}
		}
	})
	if err != nil {
		return err
	}
	return NewDefaultCommandFlagsValidator(nil).ValidateParsedFlags(flags)
}

func ValidateTimeoutOption(optionName string, flags *flag.FlagSet, maxAllowedValue int) error {
	optionValueStr := flags.Lookup(optionName).Value.String()

	optionValue, err := strconv.Atoi(optionValueStr)
	if err != nil || optionValue < 0 || optionValue > maxAllowedValue {
		return fmt.Errorf("Invalid value for %s: %s. Value must be in the range 0 to %d.", optionName, optionValueStr, maxAllowedValue)
	}
	return nil
}

func ValidateBooleanFlag(flagName string, flags *flag.FlagSet) error {
	flagValueStr := flags.Lookup(flagName).Value.String()

	if flagValueStr == "" {
		return fmt.Errorf("%s flag defined but no argument specified", flagName)
	}

	if flagValueStr != "true" && flagValueStr != "false" {
		return fmt.Errorf("Invalid value for %s: %s. Expected true or false.", flagName, flagValueStr)
	}

	return nil
}
