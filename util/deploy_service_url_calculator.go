package util

import (
	"fmt"

	"github.com/cloudfoundry/cli/plugin"
	"github.com/cloudfoundry/cli/plugin/models"
)

const deployServiceHost = "deploy-service"

type DeployServiceURLCalculator interface {
	ComputeDeployServiceURL() (string, error)
}

type deployServiceURLCalculatorImpl struct {
	cliConnection plugin.CliConnection
}

func NewDeployServiceURLCalculator(cliConnection plugin.CliConnection) DeployServiceURLCalculator {
	return deployServiceURLCalculatorImpl{cliConnection: cliConnection}
}

func (c deployServiceURLCalculatorImpl) ComputeDeployServiceURL() (string, error) {
	currentSpace, err := c.getCurrentSpace()
	if err != nil {
		return "", err
	}

	sharedDomain, err := findSharedDomain(currentSpace)
	if err != nil {
		return "", err
	}

	return deployServiceHost + "." + sharedDomain.Name, nil
}

func (c deployServiceURLCalculatorImpl) getCurrentSpace() (plugin_models.GetSpace_Model, error) {
	currentSpace, err := c.cliConnection.GetCurrentSpace()
	if err != nil {
		return plugin_models.GetSpace_Model{}, err
	}
	if currentSpace.Name == "" {
		return plugin_models.GetSpace_Model{}, fmt.Errorf("No space targeted, use 'cf target -s SPACE' to target a space.")
	}
	// The currentSpace object does not hold the shared domains for the space, so we must make one additional request to retrieve them:
	return c.cliConnection.GetSpace(currentSpace.Name)
}

func findSharedDomain(space plugin_models.GetSpace_Model) (plugin_models.GetSpace_Domains, error) {
	for _, domain := range space.Domains {
		if domain.Shared {
			return domain, nil
		}
	}
	return plugin_models.GetSpace_Domains{}, fmt.Errorf("Could not find any shared domains in space: %s", space.Name)
}
