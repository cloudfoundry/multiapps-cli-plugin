package util

import (
	"fmt"

	cfrestclient "github.com/SAP/cf-mta-plugin/clients/cfrestclient"
	"github.com/SAP/cf-mta-plugin/clients/models"
)

const deployServiceHost = "deploy-service"
const defaultDeployServiceHostHttpScheme = "https"
const defaultDeployServiceEndpoint = "/public/ping"

type DeployServiceURLCalculator interface {
	ComputeDeployServiceURL() (string, error)
}

type deployServiceURLCalculatorImpl struct {
	cloudFoundryClient cfrestclient.CloudFoundryOperationsExtended
	httpGetExecutor    HttpSimpleGetExecutor
}

func NewDeployServiceURLCalculator(cloudFoundryClient cfrestclient.CloudFoundryOperationsExtended) DeployServiceURLCalculator {
	return deployServiceURLCalculatorImpl{cloudFoundryClient: cloudFoundryClient, httpGetExecutor: NewSimpleGetExecutor()}
}

func NewDeployServiceURLCalculatorWithHttpExecutor(cloudFoundryClient cfrestclient.CloudFoundryOperationsExtended, httpGetExecutor HttpSimpleGetExecutor) DeployServiceURLCalculator {
	return deployServiceURLCalculatorImpl{cloudFoundryClient: cloudFoundryClient, httpGetExecutor: httpGetExecutor}
}

func (c deployServiceURLCalculatorImpl) ComputeDeployServiceURL() (string, error) {
	result, err := c.cloudFoundryClient.GetSharedDomains()
	if err != nil {
		return "", err
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
	return models.SharedDomain{}, fmt.Errorf("Could not find any shared domains")
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
