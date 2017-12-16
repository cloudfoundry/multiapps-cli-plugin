package commands_test

import (
	"fmt"

	csrf_fakes "github.com/SAP/cf-mta-plugin/clients/csrf/fakes"
	"github.com/SAP/cf-mta-plugin/clients/mtaclient"
	"github.com/SAP/cf-mta-plugin/clients/mtaclient/fakes"

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
	var mtaClient fakes.FakeMtaClientOperations

	var sessionProvider csrf_fakes.FakeSessionProvider
	var action commands.Action
	var oc = testutil.NewUIOutputCapturer()
	var ex = testutil.NewUIExpector()

	BeforeEach(func() {
		sessionProvider = csrf_fakes.NewFakeSessionProviderBuilder().GetSession(nil).Build()
		ui.DisableTerminalOutput(true)
	})

	Describe("AbortAction", func() {
		Describe("ExecuteAction", func() {
			BeforeEach(func() {
				action = &commands.AbortAction{}
				mtaClient = fakes.NewFakeMtaClientBuilder().
					ExecuteAction(processID, "abort", mtaclient.ResponseHeader{}, nil).Build()
			})
			Context("with no error returned from backend", func() {
				It("should abort the process and exit with zero status", func() {
					output, status := oc.CaptureOutputAndStatus(func() int {
						return action.Execute(processID, commandName, mtaClient, sessionProvider).ToInt()
					})
					ex.ExpectSuccessWithOutput(status, output, []string{"Aborting multi-target app operation with id test-process-id...\n", "OK\n"})
				})
			})
			Context("with an error returned from backend", func() {
				It("should return error and exit with non-zero status", func() {
					mtaClient = fakes.NewFakeMtaClientBuilder().
						ExecuteAction(processID, "abort", mtaclient.ResponseHeader{}, fmt.Errorf("test-error")).Build()
					output, status := oc.CaptureOutputAndStatus(func() int {
						return action.Execute(processID, commandName, mtaClient, sessionProvider).ToInt()
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
				mtaClient = fakes.NewFakeMtaClientBuilder().
					ExecuteAction(processID, "retry", mtaclient.ResponseHeader{Location: "operations/" + processID + "?embed=messages"}, nil).
					GetMtaOperation(processID, "messages", &testutil.SimpleOperationResult, nil).
					Build()
			})
			Context("with no error returned from backend", func() {
				It("should retry the process and exit with zero status", func() {
					output, status := oc.CaptureOutputAndStatus(func() int {
						return action.Execute(processID, commandName, mtaClient, sessionProvider).ToInt()
					})
					ex.ExpectSuccessWithOutput(status, output, []string{"Retrying multi-target app operation with id test-process-id...\n", "OK\n",
						"Monitoring process " + processID + "...\n", "Process finished.\n", "Use \"cf dmol -i " + processID + "\" to download the logs of the process.\n"})
				})
			})
			Context("with an error returned from backend", func() {
				It("should return error and exit with non-zero status", func() {
					mtaClient = fakes.NewFakeMtaClientBuilder().
						ExecuteAction(processID, "retry", mtaclient.ResponseHeader{}, fmt.Errorf("test-error")).
						GetMtaOperation(processID, "messages", &testutil.SimpleOperationResult, nil).
						Build()
					output, status := oc.CaptureOutputAndStatus(func() int {
						return action.Execute(processID, commandName, mtaClient, sessionProvider).ToInt()
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
