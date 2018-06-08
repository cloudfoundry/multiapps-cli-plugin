package util

import (
	"fmt"
	"strings"

	cfrestclient "github.com/cloudfoundry-incubator/multiapps-cli-plugin/clients/cfrestclient"
	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/clients/models"
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
	sharedDomains, err := c.cloudFoundryClient.GetSharedDomains()
	if err != nil {
		return "", err
	}

	deployServiceURL, err := c.computeDeployServiceURL(sharedDomains)
	if err != nil {
		return "", err
	}

	return deployServiceURL, nil
}

func (c deployServiceURLCalculatorImpl) computeDeployServiceURL(domains []models.SharedDomain) (string, error) {
	if len(domains) == 0 {
		return "", fmt.Errorf("Could not compute the Deploy Service's URL as there are no shared domains on the landscape.")
	}
	possibleDeployServiceURLs := buildPossibleDeployServiceURLs(domains)
	for _, possibleDeployServiceURL := range possibleDeployServiceURLs {
		if c.isCorrectURL(possibleDeployServiceURL) {
			return possibleDeployServiceURL, nil
		}
	}
	return "", fmt.Errorf("The Deploy Service does not respond on any of the default URLs:\n" + strings.Join(possibleDeployServiceURLs, "\n") + "\n\nYou can use the command line option -u or the DEPLOY_SERVICE_URL environment variable to specify a custom URL explicitly.")
}

func buildPossibleDeployServiceURLs(domains []models.SharedDomain) ([]string) {
	var possibleDeployServiceURLs []string
	for _, domain := range domains {
		possibleDeployServiceURLs = append(possibleDeployServiceURLs, deployServiceHost + "." + domain.Name)
	}
	return possibleDeployServiceURLs
}

func (c deployServiceURLCalculatorImpl) isCorrectURL(deployServiceURL string) bool {
	uriBuilder := NewUriBuilder()
	uri, err := uriBuilder.SetScheme(defaultDeployServiceHostHttpScheme).SetPath(defaultDeployServiceEndpoint).SetHost(deployServiceURL).Build()
	statusCode, err := c.httpGetExecutor.ExecuteGetRequest(uri)
	if err != nil {
		return false
	}
	return statusCode == 200
}


