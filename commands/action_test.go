package commands_test

import (
	"fmt"

	csrf_fakes "github.com/cloudfoundry-incubator/multiapps-cli-plugin/clients/csrf/fakes"
	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/clients/mtaclient"
	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/clients/mtaclient/fakes"

	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/commands"
	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/testutil"
	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/ui"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Actions", func() {
	const operationID = "test-process-id"
	const commandName = "deploy"
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
		const actionID = "abort"
		Describe("ExecuteAction", func() {
			BeforeEach(func() {
				action = commands.GetActionToExecute(actionID, commandName)
				mtaClient = fakes.NewFakeMtaClientBuilder().
					ExecuteAction(operationID, actionID, mtaclient.ResponseHeader{}, nil).Build()
			})
			Context("with no error returned from backend", func() {
				It("should abort the process and exit with zero status", func() {
					output, status := oc.CaptureOutputAndStatus(func() int {
						return action.Execute(operationID, mtaClient, sessionProvider).ToInt()
					})
					ex.ExpectSuccessWithOutput(status, output, []string{"Executing action 'abort' on operation test-process-id...\n", "OK\n"})
				})
			})
			Context("with an error returned from backend", func() {
				It("should return error and exit with non-zero status", func() {
					mtaClient = fakes.NewFakeMtaClientBuilder().
						ExecuteAction(operationID, "abort", mtaclient.ResponseHeader{}, fmt.Errorf("test-error")).Build()
					output, status := oc.CaptureOutputAndStatus(func() int {
						return action.Execute(operationID, mtaClient, sessionProvider).ToInt()
					})
					ex.ExpectFailureOnLine(status, output, "Could not execute action 'abort' on operation test-process-id: test-error", 1)
				})
			})
		})
	})

	Describe("RetryAction", func() {
		const actionID = "retry"
		Describe("ExecuteAction", func() {
			BeforeEach(func() {
				action = commands.GetActionToExecute(actionID, commandName)
				mtaClient = fakes.NewFakeMtaClientBuilder().
					ExecuteAction(operationID, actionID, mtaclient.ResponseHeader{Location: "operations/" + operationID + "?embed=messages"}, nil).
					GetMtaOperation(operationID, "messages", &testutil.SimpleOperationResult, nil).
					Build()
			})
			Context("with no error returned from backend", func() {
				It("should retry the process and exit with zero status", func() {
					output, status := oc.CaptureOutputAndStatus(func() int {
						return action.Execute(operationID, mtaClient, sessionProvider).ToInt()
					})
					ex.ExpectSuccessWithOutput(status, output, []string{"Executing action 'retry' on operation test-process-id...\n", "OK\n",
						"Monitoring process " + operationID + "...\n", "Process finished.\n", "Use \"cf dmol -i " + operationID + "\" to download the logs of the process.\n"})
				})
			})
			Context("with an error returned from backend", func() {
				It("should return error and exit with non-zero status", func() {
					mtaClient = fakes.NewFakeMtaClientBuilder().
						ExecuteAction(operationID, "retry", mtaclient.ResponseHeader{}, fmt.Errorf("test-error")).
						GetMtaOperation(operationID, "messages", &testutil.SimpleOperationResult, nil).
						Build()
					output, status := oc.CaptureOutputAndStatus(func() int {
						return action.Execute(operationID, mtaClient, sessionProvider).ToInt()
					})
					ex.ExpectFailureOnLine(status, output, "Could not execute action 'retry' on operation test-process-id: test-error", 1)
				})
			})
		})
	})

	Describe("GetActionToExecute", func() {
		Context("with correct action id", func() {
			It("should return abort action to execute", func() {
				actionToExecute := commands.GetActionToExecute("abort", "deploy")
				Expect(actionToExecute).NotTo(BeNil())
			})
			It("should return retry action to execute", func() {
				actionToExecute := commands.GetActionToExecute("retry", "deploy")
				Expect(actionToExecute).NotTo(BeNil())
			})
			It("should return resume action to execute", func() {
				actionToExecute := commands.GetActionToExecute("resume", "deploy")
				Expect(actionToExecute).NotTo(BeNil())
			})
			It("should return monitor action to execute", func() {
				actionToExecute := commands.GetActionToExecute("monitor", "deploy")
				Expect(actionToExecute).NotTo(BeNil())
			})
		})
		Context("with incorrect action id", func() {
			It("should return nil", func() {
				actionToExecute := commands.GetActionToExecute("test", "deploy")
				Expect(actionToExecute).To(BeNil())
			})
		})
	})
})
