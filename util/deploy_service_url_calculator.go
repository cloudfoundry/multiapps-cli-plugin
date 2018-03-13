package util

import (
	"fmt"

	cfrestclient "github.com/SAP/cf-mta-plugin/clients/cfrestclient"
	"github.com/SAP/cf-mta-plugin/clients/models"
	"github.com/cloudfoundry/cli/plugin"
)

const deployServiceHost = "deploy-service"
const defaultDeployServiceHostHttpScheme = "https"
const defaultDeployServiceEndpoint = "/public/ping"

type DeployServiceURLCalculator interface {
	ComputeDeployServiceURL() (string, error)
}

type deployServiceURLCalculatorImpl struct {
	// TODO: remove the cliConnection dependency
	cliConnection      plugin.CliConnection
	cloudFoundryClient cfrestclient.CloudFoundryOperationsExtended
	httpGetExecutor    HttpSimpleGetExecutor
}

func NewDeployServiceURLCalculator(cliConnection plugin.CliConnection, cloudFoundryClient cfrestclient.CloudFoundryOperationsExtended) DeployServiceURLCalculator {
	return deployServiceURLCalculatorImpl{cliConnection: cliConnection, cloudFoundryClient: cloudFoundryClient, httpGetExecutor: NewSimpleGetExecutor()}
}

func NewDeployServiceURLCalculatorWithHttpExecutor(cliConnection plugin.CliConnection, cloudFoundryClient cfrestclient.CloudFoundryOperationsExtended, httpGetExecutor HttpSimpleGetExecutor) DeployServiceURLCalculator {
	return deployServiceURLCalculatorImpl{cliConnection: cliConnection, cloudFoundryClient: cloudFoundryClient, httpGetExecutor: httpGetExecutor}
}

func (c deployServiceURLCalculatorImpl) ComputeDeployServiceURL() (string, error) {
	result, err := c.cloudFoundryClient.GetSharedDomains()
	if err != nil {
		fmt.Printf("Maikooo err happened " + err.Error())
	}
	for _, domain := range result {
		fmt.Println(domain.Guid + " " + domain.Name + " " + domain.Url)
	}

	sharedDomain, err := c.findSharedDomain(result)
	if err != nil {
		return "", err
	}

	return deployServiceHost + "." + sharedDomain.Name, nil
}

func (c deployServiceURLCalculatorImpl) findSharedDomain(domains []models.SharedDomain) (models.SharedDomain, error) {
	for _, domain := range domains {
		if c.isCorrectDomain(domain.Name) {
			return domain, nil
		}
	}
	return models.SharedDomain{}, fmt.Errorf("Could not find any shared domains in space: %s //TODO:/./")
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
