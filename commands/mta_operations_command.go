package commands

import (
	"code.cloudfoundry.org/cli/cf/terminal"
	"code.cloudfoundry.org/cli/plugin"
	"flag"
	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/clients/baseclient"
	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/clients/models"
	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/clients/mtaclient"
	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/ui"
	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/util"
)

const allOpt = "all"

// MtaOperationsCommand is a command for listing all mta operations
type MtaOperationsCommand struct {
	*BaseCommand
}

func NewMtaOperationsCommand() *MtaOperationsCommand {
	baseCmd := &BaseCommand{flagsParser: NewDefaultCommandFlagsParser(nil), flagsValidator: NewDefaultCommandFlagsValidator(nil)}
	mtaOpsCmd := &MtaOperationsCommand{baseCmd}
	baseCmd.Command = mtaOpsCmd
	return mtaOpsCmd
}

// GetPluginCommand returns the plugin command details
func (c *MtaOperationsCommand) GetPluginCommand() plugin.Command {
	return plugin.Command{
		Name:     "mta-ops",
		HelpText: "List multi-target app operations",
		UsageDetails: plugin.Usage{
			Usage: "cf mta-ops [--mta MTA] [-u URL] [--last NUM] [--all]",
			Options: map[string]string{
				deployServiceURLOpt:          "Deploy service URL, by default 'deploy-service.<system-domain>'",
				util.GetShortOption(mtaOpt):  "ID of the deployed package",
				util.GetShortOption(lastOpt): "List last NUM operations",
				util.GetShortOption(allOpt):  "List all operations, not just the active ones",
			},
		},
	}
}

func (c *MtaOperationsCommand) defineCommandOptions(flags *flag.FlagSet) {
	flags.String(mtaOpt, "", "")
	flags.Uint(lastOpt, 0, "")
	flags.Bool(allOpt, false, "")
}

func (c *MtaOperationsCommand) executeInternal(positionalArgs []string, dsHost string, flags *flag.FlagSet, cfTarget util.CloudFoundryTarget) ExecutionStatus {
	mtaId := GetStringOpt(mtaOpt, flags)
	last := GetUintOpt(lastOpt, flags)
	all := GetBoolOpt(allOpt, flags)

	printInitialMessage(cfTarget, mtaId, all, last)

	// Create new REST client
	mtaClient := c.NewMtaClient(dsHost, cfTarget)

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
			mtaID := operation.MtaID
			if operation.MtaID == "" {
				mtaID = "N/A"
			}
			table.Add(operation.ProcessID, operation.ProcessType, mtaID, operation.Namespace, string(operation.State), operation.StartedAt, operation.User)
		}
		table.Print()
	} else {
		ui.Say("No multi-target app operations found")
	}
	return Success
}

func printInitialMessage(cfTarget util.CloudFoundryTarget, mtaId string, all bool, last uint) {
	var initialMessage string
	switch {
	case mtaId != "":
		initialMessage = "Getting multi-target app operations for %[1]s in org %[3]s / space %[4]s as %[5]s..."
	case all:
		initialMessage = "Getting all multi-target app operations in org %[3]s / space %[4]s as %[5]s..."
	case last == 1:
		initialMessage = "Getting last multi-target app operation in org %[3]s / space %[4]s as %[5]s..."
	case last != 0:
		initialMessage = "Getting last %[2]d multi-target app operations in org %[3]s / space %[4]s as %[5]s..."
	default:
		initialMessage = "Getting active multi-target app operations in org %[3]s / space %[4]s as %[5]s..."
	}
	ui.Say(initialMessage, terminal.EntityNameColor(mtaId), last, terminal.EntityNameColor(cfTarget.Org.Name),
		terminal.EntityNameColor(cfTarget.Space.Name), terminal.EntityNameColor(cfTarget.Username))
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
