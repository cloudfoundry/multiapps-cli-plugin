package util

import (
	"fmt"

	"github.com/cloudfoundry/cli/plugin"
	"github.com/cloudfoundry/cli/plugin/models"
)

const deployServiceHost = "deploy-service"
const defaultDeployServiceHostHttpScheme = "https"
const defaultDeployServiceEndpoint = "/public/ping"

type DeployServiceURLCalculator interface {
	ComputeDeployServiceURL() (string, error)
}

type deployServiceURLCalculatorImpl struct {
	cliConnection   plugin.CliConnection
	httpGetExecutor HttpSimpleGetExecutor
}

func NewDeployServiceURLCalculator(cliConnection plugin.CliConnection) DeployServiceURLCalculator {
	return deployServiceURLCalculatorImpl{cliConnection: cliConnection, httpGetExecutor: NewSimpleGetExecutor()}
}

func NewDeployServiceURLCalculatorWithHttpExecutor(cliConnection plugin.CliConnection, httpGetExecutor HttpSimpleGetExecutor) DeployServiceURLCalculator {
	return deployServiceURLCalculatorImpl{cliConnection: cliConnection, httpGetExecutor: httpGetExecutor}
}

func (c deployServiceURLCalculatorImpl) ComputeDeployServiceURL() (string, error) {
	currentSpace, err := c.getCurrentSpace()
	if err != nil {
		return "", err
	}

	sharedDomain, err := c.findSharedDomain(currentSpace)
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

func (c deployServiceURLCalculatorImpl) findSharedDomain(space plugin_models.GetSpace_Model) (plugin_models.GetSpace_Domains, error) {
	for _, domain := range space.Domains {
		if domain.Shared {
			if c.isCorrectDomain(domain.Name) {
				return domain, nil
			}
		}
	}
	return plugin_models.GetSpace_Domains{}, fmt.Errorf("Could not find any shared domains in space: %s", space.Name)
}

func (c deployServiceURLCalculatorImpl) isCorrectDomain(domainName string) bool {
	statusCode, err := c.httpGetExecutor.ExecuteGetRequest(buildDeployServiceUrl(domainName))
	if err != nil {
		return false
	}

	return statusCode == 200
}

func buildDeployServiceUrl(domainName string) string {
	uriBuilder := NewUriBuilder()
	uri, err := uriBuilder.SetScheme(defaultDeployServiceHostHttpScheme).SetPath(defaultDeployServiceEndpoint).SetHost(deployServiceHost + "." + domainName).Build()
	if err != nil {
		return ""
	}
	return uri
}
