package commands

import (
	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/util"
	"strconv"
)

type DeploymentStrategy interface {
	CreateProcessBuilder() *util.ProcessBuilder
}

type DeployCommandStrategy struct {
}

func (d *DeployCommandStrategy) CreateProcessBuilder() *util.ProcessBuilder {
	processBuilder := util.NewProcessBuilder()
	processBuilder.ProcessType(deployCommandProcessTypeProvider{}.GetProcessType())
	return processBuilder
}

type BlueGreenCommandStrategy struct {
	noConfirm bool
}

func (b *BlueGreenCommandStrategy) CreateProcessBuilder() *util.ProcessBuilder {
	processBuilder := util.NewProcessBuilder()
	processBuilder.ProcessType(blueGreenDeployCommandProcessTypeProvider{}.GetProcessType())
	processBuilder.Parameter("noConfirm", strconv.FormatBool(b.noConfirm))
	processBuilder.Parameter("keepOriginalAppNamesAfterDeploy", strconv.FormatBool(true))
	return processBuilder
}

func NewDeploymentStrategy(options map[string]CommandOption, typeProvider ProcessTypeProvider) DeploymentStrategy {
	if _, ok := typeProvider.(*blueGreenDeployCommandProcessTypeProvider); ok {
		return &BlueGreenCommandStrategy{getBoolOpt(noConfirmOpt, options)}
	}
	strategy := getStringOpt(strategyOpt, options)
	if strategy == "default" || strategy == "" {
		return &DeployCommandStrategy{}
	}
	return &BlueGreenCommandStrategy{getBoolOpt(noConfirmOpt, options)}
}

func AvailableStrategies() []string {
	return []string{"blue-green", "default"}
}
