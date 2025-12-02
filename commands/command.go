package commands

import (
	"flag"

	"code.cloudfoundry.org/cli/v8/plugin"
	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/util"
)

// Command is an interface that should be implemented by all commands
type Command interface {
	GetPluginCommand() plugin.Command
	Initialize(name string, cliConnection plugin.CliConnection)
	Execute(args []string) ExecutionStatus

	executeInternal(positionalArgs []string, dsUrl string, flags *flag.FlagSet, cfTarget util.CloudFoundryTarget) ExecutionStatus
	defineCommandOptions(flags *flag.FlagSet)
}
