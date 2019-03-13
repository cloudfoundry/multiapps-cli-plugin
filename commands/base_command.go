package commands

import (
	"crypto/tls"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"os"
	"strings"
	"time"
	"unicode"

	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/clients"
	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/clients/baseclient"
	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/clients/cfrestclient"
	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/clients/csrf"
	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/clients/models"
	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/clients/mtaclient"
	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/clients/restclient"
	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/log"
	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/ui"
	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/util"
	"github.com/cloudfoundry/cli/cf/terminal"
	"github.com/cloudfoundry/cli/plugin"
	"github.com/cloudfoundry/cli/plugin/models"
)

const (
	// DeployServiceURLEnv is the deploy service URL environment variable
	DeployServiceURLEnv           = "DEPLOY_SERVICE_URL"
	deployServiceURLOpt           = "u"
	operationIDOpt                = "i"
	actionOpt                     = "a"
	forceOpt                      = "f"
	deleteServicesOpt             = "delete-services"
	deleteServiceBrokersOpt       = "delete-service-brokers"
	noRestartSubscribedAppsOpt    = "no-restart-subscribed-apps"
	noFailOnMissingPermissionsOpt = "do-not-fail-on-missing-permissions"
	abortOnErrorOpt               = "abort-on-error"
	deployServiceHost             = "deploy-service"
)

// BaseCommand represents a base command
type BaseCommand struct {
	name                       string
	cliConnection              plugin.CliConnection
	transport                  http.RoundTripper
	jar                        http.CookieJar
	clientFactory              clients.ClientFactory
	tokenFactory               baseclient.TokenFactory
	deployServiceURLCalculator util.DeployServiceURLCalculator
}

// Initialize initializes the command with the specified name and CLI connection
func (c *BaseCommand) Initialize(name string, cliConnection plugin.CliConnection) {
	log.Tracef("Initializing command '%s'\n", name)
	transport := newTransport()
	jar := newCookieJar()
	tokenFactory := NewDefaultTokenFactory(cliConnection)
	cloudFoundryClient := cfrestclient.NewCloudFoundryRestClient(getApiEndpoint(cliConnection), transport, jar, tokenFactory)
	c.InitializeAll(name, cliConnection, transport, jar, clients.NewDefaultClientFactory(), tokenFactory, util.NewDeployServiceURLCalculator(cloudFoundryClient))
}

// InitializeAll initializes the command with the specified name, CLI connection, transport and cookie jar.
func (c *BaseCommand) InitializeAll(name string, cliConnection plugin.CliConnection,
	transport http.RoundTripper, jar http.CookieJar, clientFactory clients.ClientFactory, tokenFactory baseclient.TokenFactory, deployServiceURLCalculator util.DeployServiceURLCalculator) {
	c.name = name
	c.cliConnection = cliConnection
	c.transport = transport
	c.jar = jar
	c.clientFactory = clientFactory
	c.tokenFactory = tokenFactory
	c.deployServiceURLCalculator = deployServiceURLCalculator
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
func (c *BaseCommand) CreateFlags(host *string) (*flag.FlagSet, error) {
	flags := flag.NewFlagSet(c.name, flag.ContinueOnError)
	deployServiceURL, err := c.GetDeployServiceURL()
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
func (c *BaseCommand) NewRestClient(host string) (restclient.RestClientOperations, error) {
	space, err := c.GetSpace()
	if err != nil {
		return nil, err
	}
	org, err := c.GetOrg()
	if err != nil {
		return nil, err
	}
	restClient := c.clientFactory.NewRestClient(host, org.Name, space.Name, c.transport, c.jar, c.tokenFactory)
	return restClient, nil
}

func (c *BaseCommand) NewManagementRestClient(host string) (restclient.RestClientOperations, error) {
	return c.clientFactory.NewManagementRestClient(host, c.transport, c.jar, c.tokenFactory), nil
}

// NewMtaClient creates a new MTA deployer REST client
func (c *BaseCommand) NewMtaClient(host string) (mtaclient.MtaClientOperations, error) {
	space, err := c.GetSpace()
	if err != nil {
		return nil, err
	}

	return c.clientFactory.NewMtaClient(host, space.Guid, c.transport, c.jar, c.tokenFactory), nil
}

func (c *BaseCommand) NewManagementMtaClient(host string) (mtaclient.MtaClientOperations, error) {
	return c.clientFactory.NewManagementMtaClient(host, c.transport, c.jar, c.tokenFactory), nil
}

// Context holding the username, Org and Space of the current used
type Context struct {
	Username string
	Org      string
	Space    string
}

// GetContext initializes and retrieves the Context
func (c *BaseCommand) GetContext() (Context, error) {
	username, err := c.GetUsername()
	if err != nil {
		return Context{}, err
	}
	org, err := c.GetOrg()
	if err != nil {
		return Context{}, err
	}
	space, err := c.GetSpace()
	if err != nil {
		return Context{}, err
	}
	return Context{Org: org.Name, Space: space.Name, Username: username}, nil
}

// GetOrg gets the current org name from the CLI connection
func (c *BaseCommand) GetOrg() (plugin_models.Organization, error) {
	org, err := c.cliConnection.GetCurrentOrg()
	if err != nil {
		return plugin_models.Organization{}, fmt.Errorf("Could not get current org: %s", err)
	}
	if org.Name == "" {
		return plugin_models.Organization{}, fmt.Errorf("No org and space targeted, use '%s' to target an org and a space", terminal.CommandColor("cf target -o ORG -s SPACE"))
	}
	return org, nil
}

// GetSpace gets the current space name from the CLI connection
func (c *BaseCommand) GetSpace() (plugin_models.Space, error) {
	space, err := c.cliConnection.GetCurrentSpace()
	if err != nil {
		return plugin_models.Space{}, fmt.Errorf("Could not get current space: %s", err)
	}

	if space.Name == "" || space.Guid == "" {
		return plugin_models.Space{}, fmt.Errorf("No space targeted, use '%s' to target a space", terminal.CommandColor("cf target -s"))
	}
	return space, nil
}

// GetUsername gets the username from the CLI connection
func (c *BaseCommand) GetUsername() (string, error) {
	username, err := c.cliConnection.Username()
	if err != nil {
		return "", fmt.Errorf("Could not get username: %s", err)
	}
	if username == "" {
		return "", fmt.Errorf("Not logged in. Use '%s' to log in.", terminal.CommandColor("cf login"))
	}
	return username, nil
}

// GetDeployServiceURL returns the deploy service URL
func (c *BaseCommand) GetDeployServiceURL() (string, error) {
	deployServiceURL := os.Getenv(DeployServiceURLEnv)
	if deployServiceURL == "" {
		return c.deployServiceURLCalculator.ComputeDeployServiceURL()
	}
	ui.Say(fmt.Sprintf("**Attention: You've specified a custom Deploy Service URL (%s) via the environment variable 'DEPLOY_SERVICE_URL'. The application listening on that URL may be outdated, contain bugs or unreleased features or may even be modified by a potentially untrused person. Use at your own risk.**\n", deployServiceURL))
	return deployServiceURL, nil
}

// ExecuteAction executes the action over the process specified with operationID
func (c *BaseCommand) ExecuteAction(operationID, actionID, host string) ExecutionStatus {
	// Create REST client
	mtaClient, err := c.NewMtaClient(host)
	if err != nil {
		ui.Failed(err.Error())
		return Failure
	}

	// find ongoing operation by the specified operationID
	ongoingOperation, err := c.findOngoingOperationByID(operationID, mtaClient)
	if err != nil {
		ui.Failed(err.Error())
		return Failure
	}

	if ongoingOperation == nil {
		ui.Failed("Multi-target app operation with id %s not found", terminal.EntityNameColor(operationID))
		return Failure
	}

	// Finds the action specified with the actionID
	action := GetActionToExecute(actionID, c.name)
	if action == nil {
		ui.Failed("Invalid action %s", terminal.EntityNameColor(actionID))
		return Failure
	}

	// Executes the action specified with actionID
	return action.Execute(operationID, mtaClient)
}

// CheckOngoingOperation checks for ongoing operation for mta with the specified id and tries to abort it
func (c *BaseCommand) CheckOngoingOperation(mtaID string, host string, force bool) (bool, error) {
	mtaClient, err := c.NewMtaClient(host)
	if err != nil {
		return false, err
	}

	// Check if there is an ongoing operation for this MTA ID
	ongoingOperation, err := c.findOngoingOperation(mtaID, mtaClient)
	if err != nil {
		return false, err
	}
	if ongoingOperation != nil {
		// Abort the conflict process if confirmed by the user
		if c.shouldAbortConflictingOperation(mtaID, force) {
			action := GetActionToExecute("abort", c.name)
			status := action.Execute(ongoingOperation.ProcessID, mtaClient)
			if status == Failure {
				return false, nil
			}
		} else {
			ui.Warn("%s cancelled", capitalizeFirst(string(ongoingOperation.ProcessType)))
			return false, nil
		}
	}

	return true, nil
}

func (c *BaseCommand) findOngoingOperationByID(processID string, mtaClient mtaclient.MtaClientOperations) (*models.Operation, error) {
	ongoingOperations, err := mtaClient.GetMtaOperations(nil, nil)
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
func (c *BaseCommand) findOngoingOperation(mtaID string, mtaClient mtaclient.MtaClientOperations) (*models.Operation, error) {
	activeStatesList := []string{"RUNNING", "ERROR", "ACTION_REQUIRED"}
	ongoingOperations, err := mtaClient.GetMtaOperations(nil, activeStatesList)
	if err != nil {
		return nil, fmt.Errorf("Could not get ongoing operations for multi-target app %s: %s", terminal.EntityNameColor(mtaID), err)
	}
	for _, ongoingOperation := range ongoingOperations {
		isConflicting, err := c.isConflicting(ongoingOperation, mtaID)
		if err != nil {
			return nil, err
		}
		if isConflicting {
			return ongoingOperation, nil
		}
	}
	return nil, nil
}

func (c *BaseCommand) isConflicting(operation *models.Operation, mtaID string) (bool, error) {
	space, err := c.GetSpace()
	if err != nil {
		return false, err
	}
	return operation.MtaID == mtaID && operation.SpaceID == space.Guid && operation.AcquiredLock, nil
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
