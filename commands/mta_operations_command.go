package commands

import (
	"strconv"
	"strings"

	"github.com/cloudfoundry/cli/cf/terminal"
	"github.com/cloudfoundry/cli/plugin"
	"github.com/SAP/cf-mta-plugin/clients/models"
	restclient "github.com/SAP/cf-mta-plugin/clients/restclient"
	"github.com/SAP/cf-mta-plugin/log"
	"github.com/SAP/cf-mta-plugin/ui"
	"github.com/SAP/cf-mta-plugin/util"
)

// MtaOperationsCommand is a command for listing all mta operations
type MtaOperationsCommand struct {
	BaseCommand
}

// GetPluginCommand returns the plugin command details
func (c *MtaOperationsCommand) GetPluginCommand() plugin.Command {
	return plugin.Command{
		Name:     "mta-ops",
		HelpText: "List multi-target app operations",
		UsageDetails: plugin.Usage{
			Usage: "cf mta-ops [-u URL] [--last NUM] [--all]",
			Options: map[string]string{
				"u": "Deploy service URL, by default 'deploy-service.<system-domain>'",
				util.GetShortOption("last"): "List last NUM operations",
				util.GetShortOption("all"):  "List all operations, not just the active ones",
			},
		},
	}
}

// Execute executes the command
func (c *MtaOperationsCommand) Execute(args []string) ExecutionStatus {
	log.Tracef("Executing command '"+c.name+"': args: '%v'\n", args)

	var host string
	var last uint
	var all bool

	// Parse command arguments and check for required options
	flags, err := c.CreateFlags(&host)
	if err != nil {
		ui.Failed(err.Error())
		return Failure
	}
	flags.UintVar(&last, "last", 0, "")
	flags.BoolVar(&all, "all", false, "")
	err = c.ParseFlags(args, nil, flags, nil)
	if err != nil {
		c.Usage(err.Error())
		return Failure
	}

	context, err := c.GetContext()
	if err != nil {
		ui.Failed(err.Error())
		return Failure
	}

	printInitialMessage(context, all, last)

	// Create new REST client
	restClient, err := c.NewRestClient(host)
	if err != nil {
		ui.Failed(err.Error())
		return Failure
	}

	// Get ongoing operations
	operationsToPrint, err := getOperationsToPrint(restClient, last, all)
	if err != nil {
		ui.Failed("Could not get multi-target app operations: %s", err)
		return Failure
	}
	ui.Ok()

	if len(operationsToPrint) > 0 {
		table := ui.Table([]string{"id", "type", "mta id", "status", "started at", "started by"})
		for _, operation := range operationsToPrint {
			var mtaid string = operation.MtaID
			if operation.MtaID == "" {
				mtaid = "N/A"
			}
			table.Add(*operation.ProcessID, string(operation.ProcessType), mtaid, getStatus(operation.State), *operation.StartedAt, *operation.User)
		}
		table.Print()
	} else {
		ui.Say("No multi-target app operations found")
	}
	return Success
}

func printInitialMessage(context Context, all bool, last uint) {
	var initialMessage string
	if all {
		initialMessage = "Getting all multi-target app operations in org %[2]s / space %[3]s as %[4]s..."
	} else if last == 1 {
		initialMessage = "Getting last multi-target app operation in org %[2]s / space %[3]s as %[4]s..."
	} else if last != 0 {
		initialMessage = "Getting last %[1]d multi-target app operations in org %[2]s / space %[3]s as %[4]s..."
	} else {
		initialMessage = "Getting active multi-target app operations in org %[2]s / space %[3]s as %[4]s..."
	}
	ui.Say(initialMessage, last, terminal.EntityNameColor(context.Org), terminal.EntityNameColor(context.Space),
		terminal.EntityNameColor(context.Username))
}

func getOperationsToPrint(restClient restclient.RestClientOperations, last uint, all bool) ([]*models.Operation, error) {
	var ops models.Operations
	var err error
	if all {
		// Get all operations
		ops, err = restClient.GetOperations(nil, nil)
	} else {
		if last == 0 {
			// Get operations in active state
			ops, err = restClient.GetOperations(nil, activeStatesList)
		} else {
			// Get last requested operations
			requestedOperationsCount := strconv.Itoa(int(last))
			ops, err = restClient.GetOperations(&requestedOperationsCount, nil)
		}
	}
	if err != nil {
		return []*models.Operation{}, err
	}
	return ops.Operations, nil
}

var activeStatesList = []string{"SLP_TASK_STATE_RUNNING", "SLP_TASK_STATE_ERROR", "SLP_TASK_STATE_DIALOG", "SLP_TASK_STATE_ACTION_REQUIRED", "SLP_TASK_STATE_BREAKPOINT"}

func getStatus(state models.SlpTaskStateEnum) string {
	if strings.HasPrefix(string(state), "SLP_TASK_STATE_") {
		return string(state[len("SLP_TASK_STATE_"):])
	}
	return string(state)
}
