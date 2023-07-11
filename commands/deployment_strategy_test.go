package commands_test

import (
	"flag"
	"io"
	"strconv"

	"github.com/cloudfoundry/multiapps-cli-plugin/commands"
	"github.com/cloudfoundry/multiapps-cli-plugin/commands/fakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Deployment Strategy", func() {
	const noConfirmOpt = "noConfirm"
	const keepOriginalNamesAfterDeploy = "keepOriginalAppNamesAfterDeploy"
	const skipIdleStart = "skipIdleStart"

	var deployProcessTypeProvider = &fakes.FakeDeployCommandProcessTypeProvider{}
	var bgDeployProcessTypeProvider = &fakes.FakeBlueGreenCommandProcessTypeProvider{}

	var createFlags = func(noConfirm bool, skipIdleStart bool, strategy string) *flag.FlagSet {
		flags := flag.NewFlagSet("", flag.ContinueOnError)
		flags.SetOutput(io.Discard)

		flags.String("strategy", strategy, "")
		flags.Bool("no-confirm", noConfirm, "")
		flags.Bool("skip-testing-phase", true, "")
		flags.Bool("skip-idle-start", skipIdleStart, "")
		return flags
	}

	var testInputAndOperationProcessTypesMatch = func(provider commands.ProcessTypeProvider) {
		flags := createFlags(false, false, "default")
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
			flags := createFlags(true, false, "default")

			processBuilder := commands.NewDeploymentStrategy(flags, bgDeployProcessTypeProvider).CreateProcessBuilder()
			operation := processBuilder.Build()

			Expect(operation.Parameters[noConfirmOpt]).To(Equal(strconv.FormatBool(true)))
		})
	})

	Context("with a blue-green deploy command and --skip-idle-start flag", func() {
		It("should build a blue-green deploy operation with the skipIdleStart parameter set to true", func() {
			flags := createFlags(false, true, "default")

			processBuilder := commands.NewDeploymentStrategy(flags, bgDeployProcessTypeProvider).CreateProcessBuilder()
			operation := processBuilder.Build()

			Expect(operation.Parameters[skipIdleStart]).To(Equal(strconv.FormatBool(true)))
		})
	})

	Context("with a deploy command with strategy flag set to default", func() {
		It("should build a deploy operation", func() {
			testInputAndOperationProcessTypesMatch(deployProcessTypeProvider)
		})
	})

	Context("with a deploy command with strategy flag set to blue-green", func() {
		It("should build a blue-green deploy operation", func() {
			flags := createFlags(false, false, "blue-green")

			processBuilder := commands.NewDeploymentStrategy(flags, deployProcessTypeProvider).CreateProcessBuilder()
			operation := processBuilder.Build()

			Expect(operation.ProcessType).To(Equal(bgDeployProcessTypeProvider.GetProcessType()))
			Expect(operation.Parameters[keepOriginalNamesAfterDeploy]).To(Equal(strconv.FormatBool(true)))
		})
	})

	Context("with a deploy command with strategy flag set to blue-green and --no-confirm flag present", func() {
		It("should build a blue-green deploy operation with the noConfirm parameter set to true", func() {
			flags := createFlags(true, false, "blue-green")

			processBuilder := commands.NewDeploymentStrategy(flags, deployProcessTypeProvider).CreateProcessBuilder()
			operation := processBuilder.Build()

			Expect(operation.ProcessType).To(Equal(bgDeployProcessTypeProvider.GetProcessType()))
			Expect(operation.Parameters[noConfirmOpt]).To(Equal(strconv.FormatBool(true)))
			Expect(operation.Parameters[keepOriginalNamesAfterDeploy]).To(Equal(strconv.FormatBool(true)))
		})
	})

	Context("with a deploy command with strategy flag set to blue-green and skip-idl-start set to true", func() {
		It("should build a blue-green deploy operation", func() {
			flags := createFlags(false, true, "blue-green")

			processBuilder := commands.NewDeploymentStrategy(flags, deployProcessTypeProvider).CreateProcessBuilder()
			operation := processBuilder.Build()

			Expect(operation.ProcessType).To(Equal(bgDeployProcessTypeProvider.GetProcessType()))
			Expect(operation.Parameters[keepOriginalNamesAfterDeploy]).To(Equal(strconv.FormatBool(true)))
			Expect(operation.Parameters[skipIdleStart]).To(Equal(strconv.FormatBool(true)))
		})
	})
})
