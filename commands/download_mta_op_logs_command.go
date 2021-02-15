package commands

import (
	"flag"
	"fmt"
	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/clients/mtaclient"
	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/util"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/clients/baseclient"
	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/log"
	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/ui"
	"github.com/cloudfoundry/cli/cf/terminal"
	"github.com/cloudfoundry/cli/plugin"
)

const (
	defaultDownloadDirPrefix string = "mta-op-"
)

// DownloadMtaOperationLogsCommand is a command for retrieving the logs of an MTA operation
type DownloadMtaOperationLogsCommand struct {
	BaseCommand
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
				"i": "Operation ID",
				util.GetShortOption("mta"):  "ID of the deployed MTA",
				util.GetShortOption("last"): "Downloads last NUM operation logs. If not specified, logs for each process with the specified MTA_ID are downloaded",
				"d": "Root directory to download logs, by default the current working directory",
				"u": "Deploy service URL, by default 'deploy-service.<system-domain>'",
			},
		},
	}
}

// Execute executes the command
func (c *DownloadMtaOperationLogsCommand) Execute(args []string) ExecutionStatus {
	log.Tracef("Executing command '"+c.name+"': args: '%v'\n", args)

	var host string
	var operationId string
	var mtaId string
	var last uint
	var downloadDirName string

	// Parse command arguments and check for required options
	flags, err := c.CreateFlags(&host, args)
	if err != nil {
		ui.Failed(err.Error())
		return Failure
	}
	flags.StringVar(&operationId, "i", "", "")
	flags.StringVar(&downloadDirName, "d", "", "")
	flags.StringVar(&mtaId, "mta", "", "")
	flags.UintVar(&last, "last", 0, "")
	parser := NewCommandFlagsParser(flags, NewDefaultCommandFlagsParser([]string{}), dmolCommandFlagsValidator{})
	err = parser.Parse(args)
	if err != nil {
		c.Usage(err.Error())
		return Failure
	}

	cfTarget, err := c.GetCFTarget()
	if err != nil {
		ui.Failed(err.Error())
		return Failure
	}

	// Create new SLMP client
	mtaClient := c.NewMtaClient(host, cfTarget)

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
		err = downloadLogsForProcess(opId, downloadPath, mtaClient, cfTarget)
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
	return ioutil.WriteFile(filepath.Join(downloadDir, logID), []byte(*content), 0644)
}

type dmolCommandFlagsValidator struct{}

func (dmolCommandFlagsValidator) ValidateParsedFlags(flags *flag.FlagSet) error {
	if hasValue(flags, "i") && hasValue(flags, "mta") {
		return fmt.Errorf("Option -i and option --mta are incompatible")
	}
	return NewDefaultCommandFlagsValidator(map[string]bool{
		"i":   !hasValue(flags, "mta"),
		"mta": hasValue(flags, "mta")}).ValidateParsedFlags(flags)
}

func hasValue(flags *flag.FlagSet, flagName string) bool {
	return flags.Lookup(flagName).Value.String() != ""
}
