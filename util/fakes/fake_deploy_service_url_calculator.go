package fakes

import "github.com/cloudfoundry-incubator/multiapps-cli-plugin/util"

type deployServiceURLFakeCalculatorImpl struct {
	deployServiceURL string
}

func NewDeployServiceURLFakeCalculator(deployServiceURL string) util.DeployServiceURLCalculator {
	return deployServiceURLFakeCalculatorImpl{deployServiceURL: deployServiceURL}
}

func (c deployServiceURLFakeCalculatorImpl) ComputeDeployServiceURL(cmdOption string) (string, error) {
	if cmdOption != "" {
		return cmdOption, nil
	}
	return c.deployServiceURL, nil
}
