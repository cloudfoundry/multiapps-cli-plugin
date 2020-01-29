package commands_test

import (
	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/commands"
	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/commands/fakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"strconv"
)

var _ = Describe("Deployment Strategy", func() {
	const noConfirmOpt = "noConfirm"
	const keepOriginalNamesAfterDeploy = "keepOriginalAppNamesAfterDeploy"

	var deployProcessTypeProvider = &fakes.FakeDeployCommandProcessTypeProvider{}
	var bgDeployProcessTypeProvider = &fakes.FakeBlueGreenCommandProcessTypeProvider{}

	var createOptionsMap = func(noConfirm bool, strategy string) map[string]interface{} {
		options := make(map[string]interface{}, 3)
		skipTestingPhaseOption := true
		options["strategy"] = &strategy
		options["no-confirm"] = &noConfirm
		options["skip-testing-phase"] = &skipTestingPhaseOption
		return options
	}

	var testInputAndOperationProcessTypesMatch = func(provider commands.ProcessTypeProvider) {
		options := createOptionsMap(false, "default")
		processBuilder := commands.NewDeploymentStrategy(options, provider).CreateProcessBuilder()
		operation := processBuilder.Build()
		Expect(operation.ProcessType).To(Equal(provider.GetProcessType()))
	}

	Context("with a blue-green deploy command", func() {
		It("should build a blue-green deploy operation", func() {
			testInputAndOperationProcessTypesMatch(bgDeployProcessTypeProvider)
		})
	})

	Context("with a blue-green deploy command and --no-confirm flag", func() {
		It("should build a blue-green deploy operation with the noConfirm parameter set to true", func() {
			options := createOptionsMap(true, "default")

			processBuilder := commands.NewDeploymentStrategy(options, bgDeployProcessTypeProvider).CreateProcessBuilder()
			operation := processBuilder.Build()

			Expect(operation.Parameters[noConfirmOpt]).To(Equal(strconv.FormatBool(true)))
		})
	})

	Context("with a deploy command with strategy flag set to default", func() {
		It("should build a deploy operation", func() {
			testInputAndOperationProcessTypesMatch(deployProcessTypeProvider)
		})
	})

	Context("with a deploy command with strategy flag set to blue-green", func() {
		It("should build a blue-green deploy operation", func() {
			options := createOptionsMap(false, "blue-green")

			processBuilder := commands.NewDeploymentStrategy(options, deployProcessTypeProvider).CreateProcessBuilder()
			operation := processBuilder.Build()

			Expect(operation.ProcessType).To(Equal(bgDeployProcessTypeProvider.GetProcessType()))
			Expect(operation.Parameters[keepOriginalNamesAfterDeploy]).To(Equal(strconv.FormatBool(true)))
		})
	})

	Context("with a deploy command with strategy flag set to blue-green and --no-confirm flag present", func() {
		It("should build a blue-green deploy operation with the noConfirm parameter set to true", func() {
			options := createOptionsMap(true, "blue-green")

			processBuilder := commands.NewDeploymentStrategy(options, deployProcessTypeProvider).CreateProcessBuilder()
			operation := processBuilder.Build()

			Expect(operation.ProcessType).To(Equal(bgDeployProcessTypeProvider.GetProcessType()))
			Expect(operation.Parameters[noConfirmOpt]).To(Equal(strconv.FormatBool(true)))
			Expect(operation.Parameters[keepOriginalNamesAfterDeploy]).To(Equal(strconv.FormatBool(true)))
		})
	})
})
