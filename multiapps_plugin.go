package main

import (
	"fmt"
	"io"
	defaultlog "log"
	"os"
	"strconv"
	"strings"

	"code.cloudfoundry.org/cli/plugin"
	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/commands"
	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/log"
	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/util"
)

// Version is the version of the CLI plugin. It is injected on linking time.
var Version string = "0.0.0"

// MultiappsUserAgentSuffixOption is the default user agent suffix option. It is injected on linking time.
var MultiappsUserAgentSuffixOption string = ""

// MultiappsPlugin represents a cf CLI plugin for executing operations on MTAs
type MultiappsPlugin struct{}

// Commands contains the commands supported by this plugin
var Commands = []commands.Command{
	commands.NewDeployCommand(),
	commands.NewBlueGreenDeployCommand(),
	commands.NewMtasCommand(),
	commands.NewDmolCommand(),
	commands.NewUndeployCommand(),
	commands.NewMtaCommand(),
	commands.NewMtaOperationsCommand(),
	commands.NewPurgeConfigCommand(),
	commands.NewRollbackMtaCommand(),
}

// Run runs this plugin
func (p *MultiappsPlugin) Run(cliConnection plugin.CliConnection, args []string) {
	disableStdOut()
	if args[0] == "CLI-MESSAGE-UNINSTALL" {
		return
	}
	command, err := findCommand(args[0])
	if err != nil {
		log.Fatalln(err)
	}
	versionOutput, err := cliConnection.CliCommandWithoutTerminalOutput("version")
	if err != nil {
		log.Traceln(err)
		versionOutput = []string{util.DefaultCliVersion}
	}
	util.SetCfCliVersion(strings.Join(versionOutput, " "))
	util.SetPluginVersion(Version)
	util.SetUserAgentSuffixOption(MultiappsUserAgentSuffixOption)
	command.Initialize(command.GetPluginCommand().Name, cliConnection)
	status := command.Execute(args[1:])
	if status == commands.Failure {
		os.Exit(1)
	}
}

// GetMetadata returns the metadata of this plugin
func (p *MultiappsPlugin) GetMetadata() plugin.PluginMetadata {
	metadata := plugin.PluginMetadata{
		Name:          "multiapps",
		Version:       parseSemver(Version),
		MinCliVersion: plugin.VersionType{Major: 6, Minor: 7, Build: 0},
	}
	for _, command := range Commands {
		metadata.Commands = append(metadata.Commands, command.GetPluginCommand())
	}
	return metadata
}

func main() {
	plugin.Start(new(MultiappsPlugin))
}

func disableStdOut() {
	defaultlog.SetFlags(0)
	defaultlog.SetOutput(io.Discard)
}

func findCommand(name string) (commands.Command, error) {
	for _, command := range Commands {
		pluginCommand := command.GetPluginCommand()
		if pluginCommand.Name == name || pluginCommand.Alias == name {
			return command, nil
		}
	}
	return nil, fmt.Errorf("Could not find command with name %q", name)
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
