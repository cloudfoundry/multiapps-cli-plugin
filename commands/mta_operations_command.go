package commands

import (
	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/clients/baseclient"
	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/clients/models"
	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/clients/mtaclient"
	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/log"
	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/ui"
	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/util"
	"github.com/cloudfoundry/cli/cf/terminal"
	"github.com/cloudfoundry/cli/plugin"
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
			Usage: "cf mta-ops [--mta MTA] [-u URL] [--last NUM] [--all]",
			Options: map[string]string{
				"u": "Deploy service URL, by default 'deploy-service.<system-domain>'",
				util.GetShortOption("mta"):  "ID of the deployed package",
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
	var mtaId string
	var last uint
	var all bool

	// Parse command arguments and check for required options
	flags, err := c.CreateFlags(&host, args)
	if err != nil {
		ui.Failed(err.Error())
		return Failure
	}
	flags.StringVar(&mtaId, "mta", "", "")
	flags.UintVar(&last, "last", 0, "")
	flags.BoolVar(&all, "all", false, "")
	parser := NewCommandFlagsParser(flags, NewDefaultCommandFlagsParser(nil), NewDefaultCommandFlagsValidator(nil))
	err = parser.Parse(args)
	if err != nil {
		c.Usage(err.Error())
		return Failure
	}

	context, err := c.GetContext()
	if err != nil {
		ui.Failed("Could not get org and space: %s", baseclient.NewClientError(err))
		return Failure
	}

	printInitialMessage(context, mtaId, all, last)

	// Create new REST client
	mtaClient, err := c.NewMtaClient(host)
	if err != nil {
		ui.Failed("Could not get space ID: %s", baseclient.NewClientError(err))
		return Failure
	}

	// Get ongoing operations
	operationsToPrint, err := getOperationsToPrint(mtaClient, mtaId, last, all)
	if err != nil {
		ui.Failed("Could not get multi-target app operations: %s", baseclient.NewClientError(err))
		return Failure
	}
	ui.Ok()

	if len(operationsToPrint) > 0 {
		table := ui.Table([]string{"id", "type", "mta id", "namespace", "status", "started at", "started by"})
		for _, operation := range operationsToPrint {
			var mtaid string = operation.MtaID
			if operation.MtaID == "" {
				mtaid = "N/A"
			}
			table.Add(operation.ProcessID, string(operation.ProcessType), mtaid, operation.Namespace, string(operation.State), operation.StartedAt, operation.User)
		}
		table.Print()
	} else {
		ui.Say("No multi-target app operations found")
	}
	return Success
}

func printInitialMessage(context Context, mtaId string, all bool, last uint) {
	var initialMessage string
	if mtaId != "" {
		initialMessage = "Getting multi-target app operations for %[1]s in org %[3]s / space %[4]s as %[5]s..."
	} else if all {
		initialMessage = "Getting all multi-target app operations in org %[3]s / space %[4]s as %[5]s..."
	} else if last == 1 {
		initialMessage = "Getting last multi-target app operation in org %[3]s / space %[4]s as %[5]s..."
	} else if last != 0 {
		initialMessage = "Getting last %[2]d multi-target app operations in org %[3]s / space %[4]s as %[5]s..."
	} else {
		initialMessage = "Getting active multi-target app operations in org %[3]s / space %[4]s as %[5]s..."
	}
	ui.Say(initialMessage, terminal.EntityNameColor(mtaId), last, terminal.EntityNameColor(context.Org),
		terminal.EntityNameColor(context.Space), terminal.EntityNameColor(context.Username))
}

func getOperationsToPrint(mtaClient mtaclient.MtaClientOperations, mtaId string, last uint, all bool) ([]*models.Operation, error) {
	var ops []*models.Operation
	var err error
	if all {
		// Get all operations
		ops, err = mtaClient.GetMtaOperations(&mtaId, nil, nil)
	} else {
		ops, err = mtaClient.GetMtaOperations(&mtaId, getOperationsCount(last), getActiveStatesList(last))
	}
	if err != nil {
		return []*models.Operation{}, err
	}
	return ops, nil
}

func getOperationsCount(last uint) *int64 {
	if last == 0 {
		return nil
	}
	requestedOps := int64(last)
	return &requestedOps
}

func getActiveStatesList(last uint) []string {
	if last == 0 {
		return activeStatesList
	}
	return nil
}

var activeStatesList = []string{"RUNNING", "ERROR", "ACTION_REQUIRED"}
