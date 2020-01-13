package commands

import (
	"errors"
	"flag"
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
	defaultDownloadDirPrefix = "mta-op-"
	directoryOpt             = "d"
)

// DownloadMtaOperationLogsCommand is a command for retrieving the logs of an MTA operation
type DownloadMtaOperationLogsCommand struct {
	BaseCommand
}

func NewDownloadMtaOperationLogsCommand() *DownloadMtaOperationLogsCommand {
	return &DownloadMtaOperationLogsCommand{BaseCommand{options: getDmolCommandOptions()}}
}

func getDmolCommandOptions() map[string]CommandOption {
	return map[string]CommandOption{
		deployServiceURLOpt: deployServiceUrlOption(),
		directoryOpt:        {new(string), "", "Directory to download logs, by default '" + defaultDownloadDirPrefix + "<OPERATION_ID>/'", true},
		operationIDOpt:      {new(string), "", "Operation id", true},
	}
}

// GetPluginCommand returns the plugin command details
func (c *DownloadMtaOperationLogsCommand) GetPluginCommand() plugin.Command {
	return plugin.Command{
		Name:     "download-mta-op-logs",
		Alias:    "dmol",
		HelpText: "Download logs of multi-target app operation",
		UsageDetails: plugin.Usage{
			Usage:   "cf download-mta-op-logs -i OPERATION_ID [-d DIRECTORY] [-u URL]",
			Options: c.getOptionsForPluginCommand(),
		},
	}
}

// Execute executes the command
func (c *DownloadMtaOperationLogsCommand) Execute(args []string) ExecutionStatus {
	log.Tracef("Executing command '" + c.name + "': args: '%v'\n", args)

	parser := NewCommandFlagsParserWithValidator(c.flags, NewDefaultCommandFlagsParser(0), &dmolCommandFlagsValidator{})
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
		ui.Failed(err.Error())
		return Failure
	}

	operationID := getStringOpt(operationIDOpt, c.options)

	// Print initial message
	ui.Say("Downloading logs of multi-target app operation with id %s in org %s / space %s as %s...",
		terminal.EntityNameColor(operationID), terminal.EntityNameColor(context.Org),
		terminal.EntityNameColor(context.Space), terminal.EntityNameColor(context.Username))

	// Create new SLMP client
	mtaClient, err := c.NewMtaClient(host)
	if err != nil {
		ui.Failed("Could not get space id: %s", baseclient.NewClientError(err))
		return Failure
	}

	// Download all logs
	logs, err := mtaClient.GetMtaOperationLogs(operationID)
	if err != nil {
		ui.Failed("Could not get process logs: %s", baseclient.NewClientError(err))
		return Failure
	}

	downloadedLogs := make(map[string]string)
	for _, logx := range logs {
		content, err := mtaClient.GetMtaOperationLogContent(operationID, logx.ID)
		if err != nil {
			ui.Failed("Could not get content of log %s: %s", terminal.EntityNameColor(logx.ID), baseclient.NewClientError(err))
			return Failure
		}
		downloadedLogs[logx.ID] = content
	}
	ui.Ok()

	downloadDirName := getStringOpt(directoryOpt, c.options)
	// Set the download directory if not specified
	if downloadDirName == "" {
		downloadDirName = defaultDownloadDirPrefix + operationID + "/"
	}

	// Create the download directory
	downloadDir, err := createDownloadDirectory(downloadDirName)
	if err != nil {
		ui.Failed("Could not create download directory %s: %s", terminal.EntityNameColor(downloadDirName), baseclient.NewClientError(err))
		return Failure
	}

	// Get all logs and save their contents to the download directory
	ui.Say("Saving logs to %s...", terminal.EntityNameColor(downloadDir))
	for logID, content := range downloadedLogs {
		err = saveLogContent(downloadDir, logID, content)
		if err != nil {
			ui.Failed("Could not save log %s: %s", terminal.EntityNameColor(logID), baseclient.NewClientError(err))
			return Failure
		}
	}
	ui.Ok()
	return Success
}

func createDownloadDirectory(downloadDirName string) (string, error) {
	// Check if directory name ends with the os specific path separator
	if !strings.HasSuffix(downloadDirName, string(os.PathSeparator)) {
		//If there is no os specific path separator, put it at the end of the directory name
		downloadDirName += string(os.PathSeparator)
	}

	// Check if the directory already exists
	if stat, _ := os.Stat(downloadDirName); stat != nil {
		return "", errors.New("File or directory already exists.")
	}

	// Create the directory
	err := os.MkdirAll(downloadDirName, 0755)
	if err != nil {
		return "", nil
	}

	// Return the absolute path of the directory
	return filepath.Abs(filepath.Dir(downloadDirName))
}

func saveLogContent(downloadDir, logID string, content string) error {
	ui.Say("  %s", logID)
	return ioutil.WriteFile(downloadDir + "/" + logID, []byte(content), 0644)
}

type dmolCommandFlagsValidator struct {
}

func (d *dmolCommandFlagsValidator) ValidateFlags(flags *flag.FlagSet, _ []string) error {
	var err error
	flags.VisitAll(func(f *flag.Flag) {
		if f.Name == operationIDOpt && f.Value.String() == "" {
			err = errors.New("Missing required option '" + operationIDOpt + "'")
		}
	})
	return err
}

func (*dmolCommandFlagsValidator) IsBeforeParsing() bool {
	return false
}
