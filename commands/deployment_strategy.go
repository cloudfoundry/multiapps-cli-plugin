package commands

import (
	"flag"
	"strconv"

	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/util"
)

type DeploymentStrategy interface {
	CreateProcessBuilder() *util.ProcessBuilder
}

type DeployCommandDeploymentStrategy struct{}

func (d *DeployCommandDeploymentStrategy) CreateProcessBuilder() *util.ProcessBuilder {
	processBuilder := util.NewProcessBuilder()
	processBuilder.ProcessType((deployCommandProcessTypeProvider{}).GetProcessType())
	return processBuilder
}

type BlueGreenDeployCommandDeploymentStrategy struct {
	noConfirm     bool
	skipIdleStart bool
}

func (b *BlueGreenDeployCommandDeploymentStrategy) CreateProcessBuilder() *util.ProcessBuilder {
	processBuilder := util.NewProcessBuilder()
	processBuilder.ProcessType((blueGreenDeployCommandProcessTypeProvider{}).GetProcessType())
	processBuilder.Parameter("noConfirm", strconv.FormatBool(b.noConfirm))
	processBuilder.Parameter("skipIdleStart", strconv.FormatBool(b.skipIdleStart))
	processBuilder.Parameter("keepOriginalAppNamesAfterDeploy", strconv.FormatBool(true))
	return processBuilder
}

func NewDeploymentStrategy(flags *flag.FlagSet, typeProvider ProcessTypeProvider) DeploymentStrategy {
	if typeProvider.GetProcessType() == (blueGreenDeployCommandProcessTypeProvider{}).GetProcessType() {
		return &BlueGreenDeployCommandDeploymentStrategy{GetBoolOpt(noConfirmOpt, flags), GetBoolOpt(skipIdleStart, flags)}
	}
	strategy := GetStringOpt(strategyOpt, flags)
	if strategy == "default" {
		return &DeployCommandDeploymentStrategy{}
	}
	if GetBoolOpt(skipIdleStart, flags) {
		return &BlueGreenDeployCommandDeploymentStrategy{true, true}
	}
	return &BlueGreenDeployCommandDeploymentStrategy{GetBoolOpt(skipTestingPhase, flags), false}
}

func AvailableStrategies() []string {
	return []string{"blue-green", "default"}
}
