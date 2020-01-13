package commands

import (
	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/clients/baseclient"
	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/clients/models"
	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/clients/mtaclient"
	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/log"
	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/ui"
	"github.com/cloudfoundry/cli/cf/terminal"
	"github.com/cloudfoundry/cli/plugin"
)

const (
	lastOpt = "last"
	allOpt  = "all"
)

// MtaOperationsCommand is a command for listing all mta operations
type MtaOperationsCommand struct {
	BaseCommand
}

func NewMtaOperationsCommand() *MtaOperationsCommand {
	return &MtaOperationsCommand{BaseCommand{options: getMtaOperationsCommand()}}
}

func getMtaOperationsCommand() map[string]CommandOption {
	return map[string]CommandOption{
		deployServiceURLOpt: deployServiceUrlOption(),
		lastOpt:             {new(uint), 0, "List last NUM operations", false},
		allOpt:              {new(bool), false, "List all operations, not just the active ones", false},
	}
}

// GetPluginCommand returns the plugin command details
func (c *MtaOperationsCommand) GetPluginCommand() plugin.Command {
	return plugin.Command{
		Name:     "mta-ops",
		HelpText: "List multi-target app operations",
		UsageDetails: plugin.Usage{
			Usage:   "cf mta-ops [-u URL] [--last NUM] [--all]",
			Options: c.getOptionsForPluginCommand(),
		},
	}
}

// Execute executes the command
func (c *MtaOperationsCommand) Execute(args []string) ExecutionStatus {
	log.Tracef("Executing command '" + c.name + "': args: '%v'\n", args)

	parser := NewCommandFlagsParser(c.flags, NewDefaultCommandFlagsParser(0))
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

	context, err := c.GetContext()
	if err != nil {
		ui.Failed("Could not get org and space: %s", baseclient.NewClientError(err))
		return Failure
	}

	last := getUintOpt(lastOpt, c.options)
	all := getBoolOpt(allOpt, c.options)

	printInitialMessage(context, all, last)

	// Create new REST client
	mtaClient, err := c.NewMtaClient(host)
	if err != nil {
		ui.Failed("Could not get space id: %s", baseclient.NewClientError(err))
		return Failure
	}

	// Get ongoing operations
	operationsToPrint, err := getOperationsToPrint(mtaClient, last, all)
	if err != nil {
		ui.Failed("Could not get multi-target app operations: %s", baseclient.NewClientError(err))
		return Failure
	}
	ui.Ok()

	if len(operationsToPrint) == 0 {
		ui.Say("No multi-target app operations found")
		return Success
	}

	table := ui.Table([]string{"id", "type", "mta id", "status", "started at", "started by"})
	for _, operation := range operationsToPrint {
		mtaID := operation.MtaID
		if operation.MtaID == "" {
			mtaID = "N/A"
		}
		table.Add(operation.ProcessID, operation.ProcessType, mtaID, string(operation.State), operation.StartedAt, operation.User)
	}
	table.Print()
	return Success
}

func printInitialMessage(context Context, all bool, last uint) {
	var initialMessage string
	switch {
	case all:
		initialMessage = "Getting all multi-target app operations in org %[2]s / space %[3]s as %[4]s..."
	case last == 1:
		initialMessage = "Getting last multi-target app operation in org %[2]s / space %[3]s as %[4]s..."
	case last != 0:
		initialMessage = "Getting last %[1]d multi-target app operations in org %[2]s / space %[3]s as %[4]s..."
	default:
		initialMessage = "Getting active multi-target app operations in org %[2]s / space %[3]s as %[4]s..."
	}
	ui.Say(initialMessage, last, terminal.EntityNameColor(context.Org), terminal.EntityNameColor(context.Space),
		terminal.EntityNameColor(context.Username))
}

var activeStates = []string{"RUNNING", "ERROR", "ACTION_REQUIRED"}

func getOperationsToPrint(mtaClient mtaclient.MtaClientOperations, last uint, all bool) ([]*models.Operation, error) {
	switch {
	case all:
		return mtaClient.GetMtaOperations(nil, nil)
	case last == 0:
		// Get operations in active state
		return mtaClient.GetMtaOperations(nil, activeStates)
	default:
		// Get last requested operations
		requestedOperationsCount := int64(last)
		return mtaClient.GetMtaOperations(&requestedOperationsCount, nil)
	}
}
