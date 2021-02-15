package commands

import (
	"flag"
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

func NewDeploymentStrategy(flags *flag.FlagSet, typeProvider ProcessTypeProvider) DeploymentStrategy {
	if typeProvider.GetProcessType() == (blueGreenDeployCommandProcessTypeProvider{}).GetProcessType() {
		return &BlueGreenDeployCommandDeploymentStrategy{GetBoolOpt(noConfirmOpt, flags)}
	}
	strategy := GetStringOpt(strategyOpt, flags)
	if strategy == "default" {
		return &DeployCommandDeploymentStrategy{}
	}
	return &BlueGreenDeployCommandDeploymentStrategy{GetBoolOpt(skipTestingPhase, flags)}
}

func AvailableStrategies() []string {
	return []string{"blue-green", "default"}
}
