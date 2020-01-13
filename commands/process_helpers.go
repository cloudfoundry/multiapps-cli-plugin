package commands

import "github.com/cloudfoundry-incubator/multiapps-cli-plugin/util"

type ProcessTypeProvider interface {
	GetProcessType() string
}

// ProcessParametersSetter is a function that sets the startup parameters for
// an operation. It takes them from the list of parsed flags.
type ProcessParametersSetter func(options map[string]CommandOption, processBuilder *util.ProcessBuilder)
