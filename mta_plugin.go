package main

import (
	"fmt"
	"io/ioutil"
	defaultlog "log"
	"os"
	"strconv"
	"strings"

	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/commands"
	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/log"
	"github.com/cloudfoundry/cli/plugin"
)

// Version is the version of the CLI plugin. It is injected on linking time.
var Version string = "0.0.0"

// MtaPlugin represents a cf CLI plugin for executing operations on MTAs
type MtaPlugin struct{}

// Commands contains the commands supported by this plugin
var Commands = []commands.Command{
	commands.NewDeployCommand(),
	commands.NewBlueGreenDeployCommand(),
	&commands.MtasCommand{},
	&commands.DownloadMtaOperationLogsCommand{},
	commands.NewUndeployCommand(),
	&commands.MtaCommand{},
	&commands.MtaOperationsCommand{},
	&commands.PurgeConfigCommand{},
}

// Run runs this plugin
func (p *MtaPlugin) Run(cliConnection plugin.CliConnection, args []string) {
	disableStdOut()
	if args[0] == "CLI-MESSAGE-UNINSTALL" {
		return
	}
	command, err := findCommand(args[0])
	if err != nil {
		log.Fatalln(err)
	}
	command.Initialize(command.GetPluginCommand().Name, cliConnection)
	status := command.Execute(args[1:])
	if status == commands.Failure {
		os.Exit(1)
	}
}

// GetMetadata returns the metadata of this plugin
func (p *MtaPlugin) GetMetadata() plugin.PluginMetadata {
	metadata := plugin.PluginMetadata{
		Name:          "MtaPlugin",
		Version:       parseSemver(Version),
		MinCliVersion: plugin.VersionType{Major: 6, Minor: 7, Build: 0},
	}
	for _, command := range Commands {
		metadata.Commands = append(metadata.Commands, command.GetPluginCommand())
	}
	return metadata
}

func main() {
	plugin.Start(new(MtaPlugin))
}

func disableStdOut() {
	defaultlog.SetFlags(0)
	defaultlog.SetOutput(ioutil.Discard)
}

func findCommand(name string) (commands.Command, error) {
	for _, command := range Commands {
		pluginCommand := command.GetPluginCommand()
		if pluginCommand.Name == name || pluginCommand.Alias == name {
			return command, nil
		}
	}
	return nil, fmt.Errorf("Could not find command with name '%s'", name)
}

func parseSemver(version string) plugin.VersionType {
	mmb := strings.Split(version, ".")
	if len(mmb) != 3 {
		panic("invalid version: " + version)
	}
	major, _ := strconv.Atoi(mmb[0])
	minor, _ := strconv.Atoi(mmb[1])
	build, _ := strconv.Atoi(mmb[2])

	return plugin.VersionType{
		Major: major,
		Minor: minor,
		Build: build,
	}
}
