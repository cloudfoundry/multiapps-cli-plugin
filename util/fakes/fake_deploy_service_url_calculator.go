package fakes

import "github.com/cloudfoundry-incubator/multiapps-cli-plugin/util"

type deployServiceURLFakeCalculatorImpl struct {
	deployServiceURL string
}

func NewDeployServiceURLFakeCalculator(deployServiceURL string) util.DeployServiceURLCalculator {
	return deployServiceURLFakeCalculatorImpl{deployServiceURL: deployServiceURL}
}

func (c deployServiceURLFakeCalculatorImpl) ComputeDeployServiceURL() (string, error) {
	return c.deployServiceURL, nil
}
