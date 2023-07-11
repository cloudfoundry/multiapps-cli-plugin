package util

import (
	"code.cloudfoundry.org/cli/plugin"
	"fmt"
	"github.com/cloudfoundry/multiapps-cli-plugin/configuration"
	"github.com/cloudfoundry/multiapps-cli-plugin/ui"
	"strings"
)

const deployServiceHost = "deploy-service"

type DeployServiceURLCalculator interface {
	ComputeDeployServiceURL(cmdOption string) (string, error)
}

type deployServiceURLCalculatorImpl struct {
	cliConn plugin.CliConnection
}

func NewDeployServiceURLCalculator(cliConn plugin.CliConnection) DeployServiceURLCalculator {
	return deployServiceURLCalculatorImpl{cliConn: cliConn}
}

func (c deployServiceURLCalculatorImpl) ComputeDeployServiceURL(cmdOption string) (string, error) {
	if cmdOption != "" {
		ui.Say(fmt.Sprintf("**Attention: You've specified a custom Deploy Service URL (%s) via the command line option 'u'. "+
			"The application listening on that URL may be outdated, contain bugs or unreleased features or may even be modified by a potentially untrused person. Use at your own risk.**\n", cmdOption))
		return cmdOption, nil
	}

	urlFromEnv := configuration.NewSnapshot().GetBackendURL()
	if urlFromEnv != "" {
		return urlFromEnv, nil
	}

	cfApi, _ := c.cliConn.ApiEndpoint()
	domainSeparatorIndex := strings.IndexByte(cfApi, '.')
	domain := cfApi[domainSeparatorIndex+1:]

	return deployServiceHost + "." + domain, nil
}
