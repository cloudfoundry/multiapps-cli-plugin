package commands

import (
	"crypto/tls"
	"flag"
	"fmt"
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
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"strings"
	"time"
)

const (
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

	optionParser OptionParser
	options      map[string]CommandOption
	flags        *flag.FlagSet
}

// Initialize initializes the command with the specified name and CLI connection
func (c *BaseCommand) Initialize(name string, cliConnection plugin.CliConnection) {
	log.Tracef("Initializing command '%s'\n", name)
	transport := newTransport()
	jar := newCookieJar()
	tokenFactory := NewDefaultTokenFactory(cliConnection)
	cloudFoundryClient := cfrestclient.NewCloudFoundryRestClient(getApiEndpoint(cliConnection), transport, jar, tokenFactory)
	c.createFlags()
	c.defineOptions()
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

func (c *BaseCommand) createFlags() {
	c.flags = flag.NewFlagSet(c.name, flag.ContinueOnError)
	c.flags.SetOutput(ioutil.Discard)
}

func deployServiceUrlOption() CommandOption {
	return CommandOption{new(string), "", "Deploy service URL, by default 'deploy-service.<system-domain>'", true}
}

func (c *BaseCommand) computeDeployServiceUrl() (string, error) {
	return c.deployServiceURLCalculator.ComputeDeployServiceURL(getStringOpt(deployServiceURLOpt, c.options))
}

// NewRestClient creates a new MTA deployer REST client
func (c *BaseCommand) NewRestClient(host string) restclient.RestClientOperations {
	return c.clientFactory.NewRestClient(host, c.transport, c.jar, c.tokenFactory)
}

// NewMtaClient creates a new MTA deployer REST client
func (c *BaseCommand) NewMtaClient(host string) (mtaclient.MtaClientOperations, error) {
	space, err := GetSpace(c.cliConnection)
	if err != nil {
		return nil, err
	}
	mtaClient := c.clientFactory.NewMtaClient(host, space.Guid, c.transport, c.jar, c.tokenFactory)
	return mtaClient, nil
}

func (c *BaseCommand) GetContext() (Context, error) {
	context, err := CreateContext(c.cliConnection)
	if err != nil {
		return Context{}, err
	}
	return context, nil
}

// ExecuteAction executes the action over the process specified with operationID
func (c *BaseCommand) ExecuteAction(operationID, actionID string, retries uint, host string) ExecutionStatus {
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
	action := GetActionToExecute(actionID, c.name, retries)
	if action == nil {
		ui.Failed("Invalid action %s", terminal.EntityNameColor(actionID))
		return Failure
	}

	// Executes the action specified with actionID
	return action.Execute(operationID, mtaClient)
}

// CheckOngoingOperation checks for ongoing operation for mta with the specified id and tries to abort it
func (c *BaseCommand) CheckOngoingOperation(mtaID string, force bool, host string) (bool, error) {
	mtaClient, err := c.NewMtaClient(host)
	if err != nil {
		return false, err
	}

	// Check if there is an ongoing operation for this MTA ID
	ongoingOperation, err := c.findOngoingOperation(mtaID, mtaClient)
	if err != nil {
		return false, err
	}
	if ongoingOperation == nil {
		return true, nil
	}
	if !c.shouldAbortConflictingOperation(mtaID, force) {
		ui.Warn("%s cancelled", strings.Title(ongoingOperation.ProcessType))
		return false, nil
	}
	action := GetNoRetriesActionToExecute("abort", c.name)
	status := action.Execute(ongoingOperation.ProcessID, mtaClient)
	if status == Failure {
		return false, nil
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
	space, err := GetSpace(c.cliConnection)
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

func (c *BaseCommand) defineOptions() {
	optionParser := c.optionParser
	if c.optionParser == nil {
		optionParser = DefaultOptionParser{AbstractOptionParser{}}
	}
	for name, option := range c.options {
		if parsed := optionParser.parseOption(name, option, c.flags); !parsed {
			optionParser.additionalParse(name, option, c.flags)
		}
	}
}

func getBoolOpt(name string, options map[string]CommandOption) bool {
	return *options[name].Value.(*bool)
}

func getStringOpt(name string, options map[string]CommandOption) string {
	return *options[name].Value.(*string)
}

func getUintOpt(name string, options map[string]CommandOption) uint {
	return *options[name].Value.(*uint)
}

func (c *BaseCommand) getOptionsForPluginCommand() map[string]string {
	options := make(map[string]string, len(c.options))
	for name, option := range c.options {
		options[formatOptionName(name, option.IsShortOpt)] = option.Usage
	}
	return options
}

func formatOptionName(name string, isShortOpt bool) string {
	if !isShortOpt {
		return util.GetShortOption(name)
	}
	return name
}

func newTransport() http.RoundTripper {
	csrfx := csrf.Csrf{Header: "", Token: "", IsInitialized: false, NonProtectedMethods: getNonProtectedMethods()}
	// TODO Make sure SSL verification is only skipped if the CLI is configured this way
	httpTransport := http.DefaultTransport.(*http.Transport)
	// Increase tls handshake timeout to cope with  of slow internet connection. 3 x default value =30s.
	httpTransport.TLSHandshakeTimeout = 30 * time.Second
	httpTransport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	return csrf.Transport{Transport: httpTransport, Csrf: &csrfx, Cookies: &csrf.Cookies{Cookies: []*http.Cookie{}}}
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
