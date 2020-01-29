package commands

import (
	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/util"
	"strconv"
)

type DeploymentStrategy interface {
	CreateProcessBuilder() *util.ProcessBuilder
}

type DeployCommandDeploymentStrategy struct {}

func (d *DeployCommandDeploymentStrategy) CreateProcessBuilder() *util.ProcessBuilder {
	processBuilder := util.NewProcessBuilder()
	processBuilder.ProcessType((deployCommandProcessTypeProvider{}).GetProcessType())
	return processBuilder
}

type BlueGreenDeployCommandDeploymentStrategy struct {
	noConfirm    bool
}

func (b *BlueGreenDeployCommandDeploymentStrategy) CreateProcessBuilder() *util.ProcessBuilder {
	processBuilder := util.NewProcessBuilder()
	processBuilder.ProcessType((blueGreenDeployCommandProcessTypeProvider{}).GetProcessType())
	processBuilder.Parameter("noConfirm", strconv.FormatBool(b.noConfirm))
	processBuilder.Parameter("keepOriginalAppNamesAfterDeploy", strconv.FormatBool(true))
	return processBuilder
}

func NewDeploymentStrategy(options map[string]interface{}, typeProvider ProcessTypeProvider) DeploymentStrategy {
	if typeProvider.GetProcessType() == (blueGreenDeployCommandProcessTypeProvider{}).GetProcessType() {
		return &BlueGreenDeployCommandDeploymentStrategy{GetBoolOpt(noConfirmOpt, options)}
	}
	strategy := GetStringOpt(strategyOpt, options)
	if strategy == "default" {
		return &DeployCommandDeploymentStrategy{}
	}
	return &BlueGreenDeployCommandDeploymentStrategy{GetBoolOpt(skipTestingPhase, options)}
}

func AvailableStrategies() []string {
	return []string{"blue-green", "default"}
}
