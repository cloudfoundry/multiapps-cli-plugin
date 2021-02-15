package commands_test

import (
	"flag"
	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/commands"
	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/commands/fakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"io/ioutil"
	"strconv"
)

var _ = Describe("Deployment Strategy", func() {
	const noConfirmOpt = "noConfirm"
	const keepOriginalNamesAfterDeploy = "keepOriginalAppNamesAfterDeploy"

	var deployProcessTypeProvider = &fakes.FakeDeployCommandProcessTypeProvider{}
	var bgDeployProcessTypeProvider = &fakes.FakeBlueGreenCommandProcessTypeProvider{}

	var createFlags = func(noConfirm bool, strategy string) *flag.FlagSet {
		flags := flag.NewFlagSet("", flag.ContinueOnError)
		flags.SetOutput(ioutil.Discard)

		flags.String("strategy", strategy, "")
		flags.Bool("no-confirm", noConfirm, "")
		flags.Bool("skip-testing-phase", true, "")
		return flags
	}

	var testInputAndOperationProcessTypesMatch = func(provider commands.ProcessTypeProvider) {
		flags := createFlags(false, "default")
		processBuilder := commands.NewDeploymentStrategy(flags, provider).CreateProcessBuilder()
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
			flags := createFlags(true, "default")

			processBuilder := commands.NewDeploymentStrategy(flags, bgDeployProcessTypeProvider).CreateProcessBuilder()
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
			flags := createFlags(false, "blue-green")

			processBuilder := commands.NewDeploymentStrategy(flags, deployProcessTypeProvider).CreateProcessBuilder()
			operation := processBuilder.Build()

			Expect(operation.ProcessType).To(Equal(bgDeployProcessTypeProvider.GetProcessType()))
			Expect(operation.Parameters[keepOriginalNamesAfterDeploy]).To(Equal(strconv.FormatBool(true)))
		})
	})

	Context("with a deploy command with strategy flag set to blue-green and --no-confirm flag present", func() {
		It("should build a blue-green deploy operation with the noConfirm parameter set to true", func() {
			flags := createFlags(true, "blue-green")

			processBuilder := commands.NewDeploymentStrategy(flags, deployProcessTypeProvider).CreateProcessBuilder()
			operation := processBuilder.Build()

			Expect(operation.ProcessType).To(Equal(bgDeployProcessTypeProvider.GetProcessType()))
			Expect(operation.Parameters[noConfirmOpt]).To(Equal(strconv.FormatBool(true)))
			Expect(operation.Parameters[keepOriginalNamesAfterDeploy]).To(Equal(strconv.FormatBool(true)))
		})
	})
})
