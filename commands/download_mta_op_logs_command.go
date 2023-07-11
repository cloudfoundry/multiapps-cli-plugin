package commands

import (
	"flag"
	"fmt"
	"github.com/cloudfoundry/multiapps-cli-plugin/clients/mtaclient"
	"github.com/cloudfoundry/multiapps-cli-plugin/util"
	"os"
	"path/filepath"
	"strings"

	"code.cloudfoundry.org/cli/cf/terminal"
	"code.cloudfoundry.org/cli/plugin"
	"github.com/cloudfoundry/multiapps-cli-plugin/clients/baseclient"
	"github.com/cloudfoundry/multiapps-cli-plugin/ui"
)

const (
	defaultDownloadDirPrefix string = "mta-op-"
	mtaOpt                   string = "mta"
	lastOpt                  string = "last"
	directoryOpt             string = "d"
)

// DownloadMtaOperationLogsCommand is a command for retrieving the logs of an MTA operation
type DownloadMtaOperationLogsCommand struct {
	*BaseCommand
}

func NewDmolCommand() *DownloadMtaOperationLogsCommand {
	baseCmd := &BaseCommand{flagsParser: NewDefaultCommandFlagsParser([]string{}), flagsValidator: dmolCommandFlagsValidator{}}
	dmolCmd := &DownloadMtaOperationLogsCommand{baseCmd}
	baseCmd.Command = dmolCmd
	return dmolCmd
}

// GetPluginCommand returns the plugin command details
func (c *DownloadMtaOperationLogsCommand) GetPluginCommand() plugin.Command {
	return plugin.Command{
		Name:     "download-mta-op-logs",
		Alias:    "dmol",
		HelpText: "Download logs of multi-target app operation",
		UsageDetails: plugin.Usage{
			Usage: `cf download-mta-op-logs -i OPERATION_ID [-d DIRECTORY] [-u URL]

   cf download-mta-op-logs --mta MTA [--last NUM] [-d DIRECTORY] [-u URL]`,
			Options: map[string]string{
				operationIDOpt:               "Operation ID",
				util.GetShortOption(mtaOpt):  "ID of the deployed MTA",
				util.GetShortOption(lastOpt): "Downloads last NUM operation logs. If not specified, logs for each process with the specified MTA_ID are downloaded",
				directoryOpt:                 "Root directory to download logs, by default the current working directory",
				deployServiceURLOpt:          "Deploy service URL, by default 'deploy-service.<system-domain>'",
			},
		},
	}
}

func (c *DownloadMtaOperationLogsCommand) defineCommandOptions(flags *flag.FlagSet) {
	flags.String(operationIDOpt, "", "")
	flags.String(directoryOpt, "", "")
	flags.String(mtaOpt, "", "")
	flags.Uint(lastOpt, 0, "")
}

func (c *DownloadMtaOperationLogsCommand) executeInternal(positionalArgs []string, dsHost string, flags *flag.FlagSet, cfTarget util.CloudFoundryTarget) ExecutionStatus {
	// Create new SLMP client
	mtaClient := c.NewMtaClient(dsHost, cfTarget)

	mtaId := GetStringOpt(mtaOpt, flags)
	last := GetUintOpt(lastOpt, flags)
	operationId := GetStringOpt(operationIDOpt, flags)
	downloadDirName := GetStringOpt(directoryOpt, flags)

	var operationIds []string

	if mtaId != "" {
		operations, err := mtaClient.GetMtaOperations(&mtaId, getOperationsCount(last), nil)
		if err != nil {
			ui.Failed("Could not get operations for MTA with ID %s: %s", mtaId, baseclient.NewClientError(err))
			return Failure
		}
		for _, op := range operations {
			operationIds = append(operationIds, op.ProcessID)
		}
	} else {
		operationIds = append(operationIds, operationId)
	}

	for _, opId := range operationIds {
		downloadPath := filepath.Join(downloadDirName, defaultDownloadDirPrefix+opId)
		err := downloadLogsForProcess(opId, downloadPath, mtaClient, cfTarget)
		if err != nil {
			ui.Failed(err.Error())
			return Failure
		}
	}
	return Success
}

func downloadLogsForProcess(operationId string, downloadPath string, mtaClient mtaclient.MtaClientOperations, cfTarget util.CloudFoundryTarget) error {
	// Print initial message
	ui.Say("Downloading logs of multi-target app operation with ID %s in org %s / space %s as %s...",
		terminal.EntityNameColor(operationId), terminal.EntityNameColor(cfTarget.Org.Name),
		terminal.EntityNameColor(cfTarget.Space.Name), terminal.EntityNameColor(cfTarget.Username))

	// Download all logs
	downloadedLogs := make(map[string]*string)
	logs, err := mtaClient.GetMtaOperationLogs(operationId)
	if err != nil {
		return fmt.Errorf("Could not get process logs: %s", baseclient.NewClientError(err))
	}
	for _, logx := range logs {
		content, err := mtaClient.GetMtaOperationLogContent(operationId, logx.ID)
		if err != nil {
			return fmt.Errorf("Could not get content of log %s: %s", terminal.EntityNameColor(logx.ID), baseclient.NewClientError(err))
		}
		downloadedLogs[logx.ID] = &content
	}
	ui.Ok()

	// Create the download directory
	downloadDir, err := createDownloadDirectory(downloadPath)
	if err != nil {
		return fmt.Errorf("Could not create download directory %s: %s", terminal.EntityNameColor(downloadPath), baseclient.NewClientError(err))
	}

	// Get all logs and save their contents to the download directory
	ui.Say("Saving logs to %s...", terminal.EntityNameColor(downloadDir))
	for logID, content := range downloadedLogs {
		err = saveLogContent(downloadDir, logID, content)
		if err != nil {
			return fmt.Errorf("Could not save log %s: %s", terminal.EntityNameColor(logID), baseclient.NewClientError(err))
		}
	}
	ui.Ok()
	return nil
}

func createDownloadDirectory(downloadDirName string) (string, error) {
	// Check if directory name ends with the os specific path separator
	if !strings.HasSuffix(downloadDirName, string(os.PathSeparator)) {
		//If there is no os specific path separator, put it at the end of the directory name
		downloadDirName = downloadDirName + string(os.PathSeparator)
	}

	// Check if the directory already exists
	if stat, _ := os.Stat(downloadDirName); stat != nil {
		return "", fmt.Errorf("File or directory already exists.")
	}

	// Create the directory
	err := os.MkdirAll(downloadDirName, 0755)
	if err != nil {
		return "", err
	}

	// Return the absolute path of the directory
	return filepath.Abs(filepath.Dir(downloadDirName))
}

func saveLogContent(downloadDir, logID string, content *string) error {
	ui.Say("  %s", logID)
	return os.WriteFile(filepath.Join(downloadDir, logID), []byte(*content), 0644)
}

type dmolCommandFlagsValidator struct{}

func (dmolCommandFlagsValidator) ValidateParsedFlags(flags *flag.FlagSet) error {
	if hasValue(flags, "i") && hasValue(flags, "mta") {
		return fmt.Errorf("Option -i and option --mta are incompatible")
	}
	return NewDefaultCommandFlagsValidator(map[string]bool{
		operationIDOpt: !hasValue(flags, mtaOpt),
		mtaOpt:         hasValue(flags, mtaOpt)}).ValidateParsedFlags(flags)
}

func hasValue(flags *flag.FlagSet, flagName string) bool {
	return GetStringOpt(flagName, flags) != ""
}
