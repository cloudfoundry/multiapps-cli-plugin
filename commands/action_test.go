package commands_test

import (
	"fmt"

	slppclient "github.com/SAP/cf-mta-plugin/clients/slppclient"
	"github.com/SAP/cf-mta-plugin/clients/slppclient/fakes"
	"github.com/SAP/cf-mta-plugin/commands"
	"github.com/SAP/cf-mta-plugin/testutil"
	"github.com/SAP/cf-mta-plugin/ui"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Actions", func() {
	const processID = "test-process-id"
	const commandName = "deploy"
	const actionID = "abort"
	var slppClient slppclient.SlppClientOperations
	var action commands.Action
	var oc = testutil.NewUIOutputCapturer()
	var ex = testutil.NewUIExpector()

	BeforeEach(func() {
		ui.DisableTerminalOutput(true)
	})

	Describe("AbortAction", func() {
		Describe("ExecuteAction", func() {
			BeforeEach(func() {
				action = &commands.AbortAction{}
				slppClient = fakes.NewFakeSlppClientBuilder().
					GetMetadata(&testutil.SlppMetadataResult, nil).
					ExecuteAction("abort", nil).Build()
			})
			Context("with no error returned from backend", func() {
				It("should abort the process and exit with zero status", func() {
					output, status := oc.CaptureOutputAndStatus(func() int {
						return action.Execute(processID, commandName, slppClient).ToInt()
					})
					ex.ExpectSuccessWithOutput(status, output, []string{"Aborting multi-target app operation with id test-process-id...\n", "OK\n"})
				})
			})
			Context("with an error returned from backend", func() {
				It("should return error and exit with non-zero status", func() {
					slppClient = fakes.NewFakeSlppClientBuilder().
						GetMetadata(&testutil.SlppMetadataResult, nil).
						ExecuteAction("abort", fmt.Errorf("test-error")).Build()
					output, status := oc.CaptureOutputAndStatus(func() int {
						return action.Execute(processID, commandName, slppClient).ToInt()
					})
					ex.ExpectFailureOnLine(status, output, "Could not abort multi-target app operation with id test-process-id: test-error", 1)
				})
			})
		})
	})

	Describe("RetryAction", func() {
		Describe("ExecuteAction", func() {
			BeforeEach(func() {
				action = &commands.RetryAction{}
				slppClient = fakes.NewFakeSlppClientBuilder().
					GetMetadata(&testutil.SlppMetadataResult, nil).
					ExecuteAction("retry", nil).
					GetTasklistTask(&testutil.TaskResult, nil).
					Build()
			})
			Context("with no error returned from backend", func() {
				It("should retry the process and exit with zero status", func() {
					output, status := oc.CaptureOutputAndStatus(func() int {
						return action.Execute(processID, commandName, slppClient).ToInt()
					})
					ex.ExpectSuccessWithOutput(status, output, []string{"Retrying multi-target app operation with id test-process-id...\n", "OK\n",
						"Monitoring process execution...\n", "Process finished.\n"})
				})
			})
			Context("with an error returned from backend", func() {
				It("should return error and exit with non-zero status", func() {
					slppClient = fakes.NewFakeSlppClientBuilder().
						GetMetadata(&testutil.SlppMetadataResult, nil).
						ExecuteAction("retry", fmt.Errorf("test-error")).
						GetTasklistTask(&testutil.TaskResult, nil).
						Build()
					output, status := oc.CaptureOutputAndStatus(func() int {
						return action.Execute(processID, commandName, slppClient).ToInt()
					})
					ex.ExpectFailureOnLine(status, output, "Could not retry multi-target app operation with id test-process-id: test-error", 1)
				})
			})
		})
	})

	Describe("GetActionToExecute", func() {
		Context("with correct action id", func() {
			It("should return abort action to execute", func() {
				actionToExecute := commands.GetActionToExecute("abort")
				Expect(actionToExecute).NotTo(BeNil())
				Expect(actionToExecute).To(Equal(&commands.AbortAction{}))
			})
			It("should return retry action to execute", func() {
				actionToExecute := commands.GetActionToExecute("retry")
				Expect(actionToExecute).NotTo(BeNil())
				Expect(actionToExecute).To(Equal(&commands.RetryAction{}))
			})
		})
		Context("with incorrect action id", func() {
			It("should return nil", func() {
				actionToExecute := commands.GetActionToExecute("test")
				Expect(actionToExecute).To(BeNil())
			})
		})
	})
})
