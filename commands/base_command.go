package commands

import (
	"crypto/tls"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"strings"
	"time"
	"unicode"

	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/clients"
	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/clients/baseclient"
	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/clients/cfrestclient"
	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/clients/cfrestclient/resilient"
	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/clients/csrf"
	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/clients/models"
	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/clients/mtaclient"
	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/clients/mtaclient_v2"
	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/clients/restclient"
	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/configuration"
	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/log"
	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/ui"
	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/util"
	"github.com/cloudfoundry/cli/cf/terminal"
	"github.com/cloudfoundry/cli/plugin"
)

const (
	// DeployServiceURLEnv is the deploy service URL environment variable
	deployServiceURLOpt           = "u"
	operationIDOpt                = "i"
	actionOpt                     = "a"
	forceOpt                      = "f"
	deleteServicesOpt             = "delete-services"
	deleteServiceBrokersOpt       = "delete-service-brokers"
	noRestartSubscribedAppsOpt    = "no-restart-subscribed-apps"
	noFailOnMissingPermissionsOpt = "do-not-fail-on-missing-permissions"
	abortOnErrorOpt               = "abort-on-error"
	retriesOpt                    = "retries"
	namespaceOpt                  = "namespace"
)

const maxRetriesCount = 3
const retryIntervalInSeconds = 10

// BaseCommand represents a base command
type BaseCommand struct {
	name                       string
	cliConnection              plugin.CliConnection
	transport                  http.RoundTripper
	jar                        http.CookieJar
	clientFactory              clients.ClientFactory
	tokenFactory               baseclient.TokenFactory
	deployServiceURLCalculator util.DeployServiceURLCalculator
	configurationSnapshot      configuration.Snapshot
}

// Initialize initializes the command with the specified name and CLI connection
func (c *BaseCommand) Initialize(name string, cliConnection plugin.CliConnection) {
	log.Tracef("Initializing command '%s'\n", name)
	transport := newTransport()
	jar := newCookieJar()
	tokenFactory := NewDefaultTokenFactory(cliConnection)
	cloudFoundryClient := cfrestclient.NewCloudFoundryRestClient(getApiEndpoint(cliConnection), transport, jar, tokenFactory)
	resilientCloudFoundryClient := resilient.NewResilientCloudFoundryClient(cloudFoundryClient, maxRetriesCount, retryIntervalInSeconds)
	c.InitializeAll(name, cliConnection, transport, jar, clients.NewDefaultClientFactory(), tokenFactory, util.NewDeployServiceURLCalculator(resilientCloudFoundryClient), configuration.NewSnapshot())
}

// InitializeAll initializes the command with the specified name, CLI connection, transport and cookie jar.
func (c *BaseCommand) InitializeAll(name string, cliConnection plugin.CliConnection,
	transport http.RoundTripper, jar http.CookieJar, clientFactory clients.ClientFactory, tokenFactory baseclient.TokenFactory, deployServiceURLCalculator util.DeployServiceURLCalculator, configurationSnapshot configuration.Snapshot) {
	c.name = name
	c.cliConnection = cliConnection
	c.transport = transport
	c.jar = jar
	c.clientFactory = clientFactory
	c.tokenFactory = tokenFactory
	c.deployServiceURLCalculator = deployServiceURLCalculator
	c.configurationSnapshot = configurationSnapshot
}

func getApiEndpoint(cliConnection plugin.CliConnection) string {
	api, err := cliConnection.ApiEndpoint()
	if err != nil {
		return ""
	}
	if strings.HasPrefix(api, "https://") {
		api = strings.Replace(api, "https://", "", -1)
	}
	return api
}

// Usage reports incorrect command usage
func (c *BaseCommand) Usage(message string) {
	ui.Say(terminal.FailureColor("FAILED"))
	ui.Say("Incorrect usage. %s\n", message)
	_, err := c.cliConnection.CliCommand("help", c.name)
	if err != nil {
		ui.Failed("Could not display help: %s", err)
	}
}

// CreateFlags creates a flag set to be used for parsing command arguments
func (c *BaseCommand) CreateFlags(host *string, args []string) (*flag.FlagSet, error) {
	flags := flag.NewFlagSet(c.name, flag.ContinueOnError)
	deployServiceURL, err := c.GetDeployServiceURL(args)
	if err != nil {
		return nil, err
	}

	flags.StringVar(host, deployServiceURLOpt, deployServiceURL, "")
	flags.SetOutput(ioutil.Discard)
	return flags, nil
}

func GetOptionValue(args []string, optionName string) string {
	for index, arg := range args {
		trimmedArg := strings.Trim(arg, "-")
		if optionName == trimmedArg && len(args) > index+1 {
			return args[index+1]
		}
	}
	return ""
}

// NewRestClient creates a new MTA deployer REST client
func (c *BaseCommand) NewRestClient(host string) restclient.RestClientOperations {
	return c.clientFactory.NewRestClient(host, c.transport, c.jar, c.tokenFactory)
}

// NewMtaClient creates a new MTA deployer REST client
func (c *BaseCommand) NewMtaClient(host string, cfTarget util.CloudFoundryTarget) mtaclient.MtaClientOperations {
	return c.clientFactory.NewMtaClient(host, cfTarget.Space.Guid, c.transport, c.jar, c.tokenFactory)
}

// NewMtaV2Client creates a new MTAV2 deployer REST client
func (c *BaseCommand) NewMtaV2Client(host string, cfTarget util.CloudFoundryTarget) mtaclient_v2.MtaV2ClientOperations {
	return c.clientFactory.NewMtaV2Client(host, cfTarget.Space.Guid, c.transport, c.jar, c.tokenFactory)
}

// GetCFTarget initializes and retrieves the CF Target with the current user
func (c *BaseCommand) GetCFTarget() (util.CloudFoundryTarget, error) {
	cfContext := util.NewCloudFoundryContext(c.cliConnection)
	username, err := cfContext.GetUsername()
	if err != nil {
		return util.CloudFoundryTarget{}, err
	}
	org, err := cfContext.GetOrg()
	if err != nil {
		return util.CloudFoundryTarget{}, err
	}
	space, err := cfContext.GetSpace()
	if err != nil {
		return util.CloudFoundryTarget{}, err
	}
	return util.NewCFTarget(org, space, username), nil
}

// GetDeployServiceURL returns the deploy service URL
func (c *BaseCommand) GetDeployServiceURL(args []string) (string, error) {
	customDeployServiceURL := c.GetCustomDeployServiceURL(args)
	if customDeployServiceURL != "" {
		return customDeployServiceURL, nil
	}

	return c.deployServiceURLCalculator.ComputeDeployServiceURL()
}

// GetCustomDeployServiceURL returns custom deploy service URL
func (c *BaseCommand) GetCustomDeployServiceURL(args []string) string {
	optionDeployServiceURL := GetOptionValue(args, deployServiceURLOpt)

	if optionDeployServiceURL != "" {
		ui.Say(fmt.Sprintf("**Attention: You've specified a custom Deploy Service URL (%s) via the command line option 'u'. The application listening on that URL may be outdated, contain bugs or unreleased features or may even be modified by a potentially untrused person. Use at your own risk.**\n", optionDeployServiceURL))
		return optionDeployServiceURL
	}
	return c.configurationSnapshot.GetBackendURL()
}

// ExecuteAction executes the action over the process specified with operationID
func (c *BaseCommand) ExecuteAction(operationID, actionID string, retries uint, host string, cfTarget util.CloudFoundryTarget) ExecutionStatus {
	mtaClient := c.NewMtaClient(host, cfTarget)

	// find ongoing operation by the specified operationID
	ongoingOperation, err := c.findOngoingOperationByID(operationID, mtaClient)
	if err != nil {
		ui.Failed(err.Error())
		return Failure
	}

	if ongoingOperation == nil {
		ui.Failed("Multi-target app operation with ID %s not found", terminal.EntityNameColor(operationID))
		return Failure
	}

	// Finds the action specified with the actionID
	action := GetActionToExecute(actionID, c.name, retries)
	if action == nil {
		ui.Failed("Invalid action %s", terminal.EntityNameColor(actionID))
		return Failure
	}

	// Executes the action specified with actionID
	return action.Execute(operationID, mtaClient)
}

// CheckOngoingOperation checks for ongoing operation for mta with the specified id and tries to abort it
func (c *BaseCommand) CheckOngoingOperation(mtaID string, namespace string, host string, force bool, cfTarget util.CloudFoundryTarget) (bool, error) {
	mtaClient := c.NewMtaClient(host, cfTarget)

	// Check if there is an ongoing operation for this MTA ID
	ongoingOperation, err := c.findOngoingOperation(mtaID, namespace, mtaClient, cfTarget)
	if err != nil {
		return false, err
	}
	if ongoingOperation != nil {
		// Abort the conflict process if confirmed by the user
		if c.shouldAbortConflictingOperation(mtaID, force) {
			action := GetNoRetriesActionToExecute("abort", c.name)
			status := action.Execute(ongoingOperation.ProcessID, mtaClient)
			if status == Failure {
				return false, nil
			}
		} else {
			ui.Warn("%s cancelled", capitalizeFirst(ongoingOperation.ProcessType))
			return false, nil
		}
	}

	return true, nil
}

func (c *BaseCommand) findOngoingOperationByID(processID string, mtaClient mtaclient.MtaClientOperations) (*models.Operation, error) {
	ongoingOperations, err := mtaClient.GetMtaOperations(nil, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("Could not get ongoing operation with id %s: %s", terminal.EntityNameColor(processID), err)
	}

	for _, ongoingOperation := range ongoingOperations {
		if ongoingOperation.ProcessID == processID {
			return ongoingOperation, nil
		}
	}
	return nil, nil
}

// FindOngoingOperation finds ongoing operation for mta with the specified id
func (c *BaseCommand) findOngoingOperation(mtaID string, namespace string, mtaClient mtaclient.MtaClientOperations, cfTarget util.CloudFoundryTarget) (*models.Operation, error) {
	activeStatesList := []string{"RUNNING", "ERROR", "ACTION_REQUIRED"}
	ongoingOperations, err := mtaClient.GetMtaOperations(&mtaID, nil, activeStatesList)
	if err != nil {
		return nil, fmt.Errorf("Could not get ongoing operations for multi-target app %s: %s", terminal.EntityNameColor(mtaID), err)
	}
	for _, ongoingOperation := range ongoingOperations {
		isConflicting := c.isConflicting(ongoingOperation, mtaID, namespace, cfTarget)
		if isConflicting {
			return ongoingOperation, nil
		}
	}
	return nil, nil
}

func (c *BaseCommand) isConflicting(operation *models.Operation, mtaID string, namespace string, cfTarget util.CloudFoundryTarget) bool {
	return operation.MtaID == mtaID &&
		operation.SpaceID == cfTarget.Space.Guid &&
		operation.Namespace == namespace &&
		operation.AcquiredLock
}

func (c *BaseCommand) shouldAbortConflictingOperation(mtaID string, force bool) bool {
	if force {
		return true
	}
	return ui.Confirm("There is an ongoing operation for multi-target app %s. Do you want to abort it? (y/n)",
		terminal.EntityNameColor(mtaID))
}

func newTransport() http.RoundTripper {
	csrfx := csrf.Csrf{Header: "", Token: "", IsInitialized: false, NonProtectedMethods: getNonProtectedMethods()}
	// TODO Make sure SSL verification is only skipped if the CLI is configured this way
	httpTransport := http.DefaultTransport.(*http.Transport)
	// Increase tls handshake timeout to cope with  of slow internet connection. 3 x default value =30s.
	httpTransport.TLSHandshakeTimeout = 30 * time.Second
	httpTransport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	return csrf.Transport{Transport: httpTransport, Csrf: &csrfx, Cookies: &csrf.Cookies{[]*http.Cookie{}}}
}

func getNonProtectedMethods() map[string]bool {
	nonProtectedMethods := make(map[string]bool)

	nonProtectedMethods[http.MethodGet] = true
	nonProtectedMethods[http.MethodHead] = true
	nonProtectedMethods[http.MethodOptions] = true

	return nonProtectedMethods
}

func newCookieJar() http.CookieJar {
	jar, err := cookiejar.New(nil)
	if err != nil {
		panic(fmt.Sprintf("Could not create cookie jar: %s", err))
	}
	return jar
}

func getTokenValue(tokenString string) string {
	// TODO(ivan): check whether there are >1 elements
	return strings.Fields(tokenString)[1]
}

func capitalizeFirst(s string) string {
	if s == "" {
		return s
	}
	a := []rune(s)
	a[0] = unicode.ToUpper(a[0])
	return string(a)
}
