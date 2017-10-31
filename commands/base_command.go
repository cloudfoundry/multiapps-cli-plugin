package commands

import (
	"crypto/tls"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"sort"
	"strings"
	"unicode"

	"github.com/SAP/cf-mta-plugin/clients"
	baseclient "github.com/SAP/cf-mta-plugin/clients/baseclient"
	"github.com/SAP/cf-mta-plugin/clients/csrf"
	"github.com/SAP/cf-mta-plugin/clients/models"
	mtaclient "github.com/SAP/cf-mta-plugin/clients/mtaclient"
	restclient "github.com/SAP/cf-mta-plugin/clients/restclient"
	"github.com/SAP/cf-mta-plugin/log"
	"github.com/SAP/cf-mta-plugin/ui"
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
)

// BaseCommand represents a base command
type BaseCommand struct {
	name          string
	cliConnection plugin.CliConnection
	transport     http.RoundTripper
	jar           http.CookieJar
	clientFactory clients.ClientFactory
	tokenFactory  baseclient.TokenFactory
}

// Initialize initializes the command with the specified name and CLI connection
func (c *BaseCommand) Initialize(name string, cliConnection plugin.CliConnection) {
	log.Tracef("Initializing command '%s'\n", name)
	c.InitializeAll(name, cliConnection, newTransport(), newCookieJar(), clients.NewDefaultClientFactory(), NewDefaultTokenFactory(cliConnection))
}

// InitializeAll initializes the command with the specified name, CLI connection, transport and cookie jar.
func (c *BaseCommand) InitializeAll(name string, cliConnection plugin.CliConnection,
	transport http.RoundTripper, jar http.CookieJar, clientFactory clients.ClientFactory, tokenFactory baseclient.TokenFactory) {
	c.name = name
	c.cliConnection = cliConnection
	c.transport = transport
	c.jar = jar
	c.clientFactory = clientFactory
	c.tokenFactory = tokenFactory
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
	flags.StringVar(host, "u", deployServiceURL, "")
	flags.SetOutput(ioutil.Discard)
	return flags, nil
}

// ParseFlags parses the flags and checks for wrong arguments and missing required flags
func (c *BaseCommand) ParseFlags(args []string, positionalArgNames []string, flags *flag.FlagSet,
	required map[string]bool) error {
	// Check for missing positional arguments
	positionalArgsCount := len(positionalArgNames)
	if len(args) < positionalArgsCount {
		return fmt.Errorf(fmt.Sprintf("Missing positional argument '%s'.", positionalArgNames[len(args)]))
	}
	for i := 0; i < positionalArgsCount; i++ {
		if flags.Lookup(strings.Replace(args[i], "-", "", 1)) != nil {
			return fmt.Errorf("Missing positional argument '%s'.", positionalArgNames[i])
		}
	}

	// Parse the arguments
	err := flags.Parse(args[positionalArgsCount:])
	if err != nil {
		return errors.New("Unknown or wrong flag.")
	}

	// Check for wrong arguments
	if flags.NArg() > 0 {
		return errors.New("Wrong arguments.")
	}

	var missingRequiredOptions []string
	// Check for missing required flags
	flags.VisitAll(func(f *flag.Flag) {
		log.Traceln(f.Name, f.Value)
		if required[f.Name] && f.Value.String() == "" {
			missingRequiredOptions = append(missingRequiredOptions, f.Name)
		}
	})
	if len(missingRequiredOptions) != 0 {
		return fmt.Errorf("Missing required options '%v'.", missingRequiredOptions)
	}
	return nil
}

// ContainsSpecificOptions checks if the argument list contains all the specific options
func ContainsSpecificOptions(flags *flag.FlagSet, args []string, specificOptions map[string]string) (bool, error) {
	var matchedOptions int
	for _, arg := range args {
		optionFlag := flags.Lookup(strings.Replace(arg, "-", "", 1))
		if optionFlag != nil && specificOptions[optionFlag.Name] == arg {
			matchedOptions++
		}
	}

	// TODO: Move this validation to the ParseFlags function.
	if matchedOptions > 0 && matchedOptions < len(specificOptions) {
		var keys []string
		for key := range specificOptions {
			keys = append(keys, key)
		}
		sort.Strings(keys)
		return false, fmt.Errorf("All the %s options should be specified together", strings.Join(keys, " "))
	}

	return matchedOptions == len(specificOptions), nil
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

// NewSessionProvider Returns a new SessionProvider - responponsible for giving a unique token each time
func (c *BaseCommand) NewSessionProvider(host string) (csrf.SessionProvider, error) {
	// TODO: introduce a factory for the different SessionProviders
	return c.NewManagementMtaClient(host)
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
		apiEndpoint, err := c.cliConnection.ApiEndpoint()
		if err != nil {
			return "", fmt.Errorf("Could not get API endpoint: %s", err)
		}
		if apiEndpoint == "" {
			return "", fmt.Errorf("No api endpoint set. Use '%s' to set an endpoint.", terminal.CommandColor("cf api"))
		}
		url, err := url.Parse(apiEndpoint)
		if err != nil {
			return "", fmt.Errorf("Could not parse API endpoint %s: %s", terminal.EntityNameColor(apiEndpoint), err)
		}
		if strings.HasPrefix(url.Host, "api.cf.") {
			deployServiceURL = "deploy-service.cfapps" + url.Host[6:]
		} else if strings.HasPrefix(url.Host, "api.") {
			deployServiceURL = "deploy-service" + url.Host[3:]
		}
	}
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
	action := GetActionToExecute(actionID)
	if action == nil {
		ui.Failed("Invalid action %s", terminal.EntityNameColor(actionID))
		return Failure
	}

	sessionProvider, _ := c.NewSessionProvider(host)

	// Executes the action specified with actionID
	return action.Execute(operationID, c.name, mtaClient, sessionProvider)
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
			action := GetActionToExecute("abort")
			sessionProvider, _ := c.NewSessionProvider(host)
			status := action.Execute(ongoingOperation.ProcessID, c.name, mtaClient, sessionProvider)
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
	ongoingOperations, err := mtaClient.GetMtaOperations(nil, nil)
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
	csrfx := csrf.Csrf{Header: "", Token: ""}
	// TODO Make sure SSL verification is only skipped if the CLI is configured this way
	httpTransport := http.DefaultTransport.(*http.Transport)
	httpTransport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	return csrf.Transport{Transport: httpTransport, Csrf: &csrfx}
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
