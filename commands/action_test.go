package commands_test

import (
	"fmt"

	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/clients/models"
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
	var mtaClient *fakes.FakeMtaClientOperations

	var action commands.Action
	var oc = testutil.NewUIOutputCapturer()
	var ex = testutil.NewUIExpector()

	BeforeEach(func() {
		ui.DisableTerminalOutput(true)
	})

	Describe("AbortAction", func() {
		const actionID = "abort"
		Describe("ExecuteAction", func() {
			BeforeEach(func() {
				action = commands.GetNoRetriesActionToExecute(actionID, commandName)
				mtaClient = fakes.NewFakeMtaClientBuilder().
					GetOperationActions(operationID, []string{actionID}, nil).
					ExecuteAction(operationID, actionID, mtaclient.ResponseHeader{}, nil).
					Build()
			})
			Context("with no error returned from the backend", func() {
				It("should abort the process and exit with zero status", func() {
					output, status := oc.CaptureOutputAndStatus(func() int {
						return action.Execute(operationID, mtaClient).ToInt()
					})
					ex.ExpectSuccessWithOutput(status, output, []string{"Executing action \"abort\" on operation test-process-id...", "OK"})
				})
			})
			Context("with an error returned from the backend", func() {
				It("should return an error and exit with non-zero status", func() {
					mtaClient = fakes.NewFakeMtaClientBuilder().
						GetOperationActions(operationID, []string{actionID}, nil).
						ExecuteAction(operationID, "abort", mtaclient.ResponseHeader{}, fmt.Errorf("test-error")).
						Build()
					output, status := oc.CaptureOutputAndStatus(func() int {
						return action.Execute(operationID, mtaClient).ToInt()
					})
					ex.ExpectFailureOnLine(status, output, "Could not execute action \"abort\" on operation test-process-id: test-error", 2)
				})
			})
		})
	})

	Describe("RetryAction", func() {
		const actionID = "retry"
		Describe("ExecuteAction", func() {
			BeforeEach(func() {
				action = commands.GetNoRetriesActionToExecute(actionID, commandName)
				mtaClient = fakes.NewFakeMtaClientBuilder().
					GetOperationActions(operationID, []string{actionID}, nil).
					ExecuteAction(operationID, actionID, mtaclient.ResponseHeader{Location: "operations/" + operationID + "?embed=messages"}, nil).
					GetMtaOperation(operationID, "messages", &testutil.SimpleOperationResult, nil).
					Build()
			})
			Context("with no error returned from the backend", func() {
				It("should retry the process and exit with zero status", func() {
					output, status := oc.CaptureOutputAndStatus(func() int {
						return action.Execute(operationID, mtaClient).ToInt()
					})
					ex.ExpectSuccessWithOutput(status, output, []string{"Executing action \"retry\" on operation test-process-id...", "OK",
						"Process finished.", "Use \"cf dmol -i " + operationID + "\" to download the logs of the process."})
				})
			})
			Context("with an error returned from the backend", func() {
				It("should return an error and exit with non-zero status", func() {
					mtaClient = fakes.NewFakeMtaClientBuilder().
						GetOperationActions(operationID, []string{actionID}, nil).
						ExecuteAction(operationID, "retry", mtaclient.ResponseHeader{}, fmt.Errorf("test-error")).
						GetMtaOperation(operationID, "messages", &testutil.SimpleOperationResult, nil).
						Build()
					output, status := oc.CaptureOutputAndStatus(func() int {
						return action.Execute(operationID, mtaClient).ToInt()
					})
					ex.ExpectFailureOnLine(status, output, "Could not execute action \"retry\" on operation test-process-id: test-error", 2)
				})
			})
		})
	})

	Describe("MonitorAction", func() {
		const actionID = "monitor"
		Describe("ExecuteAction", func() {
			BeforeEach(func() {
				action = commands.GetNoRetriesActionToExecute(actionID, commandName)
			})
			Context("when the operation finishes successfully", func() {
				It("should monitor the operation successfully", func() {
					mtaClient = fakes.NewFakeMtaClientBuilder().
						GetMtaOperation(operationID, "messages", &testutil.SimpleOperationResult, nil).
						Build()
					output, status := oc.CaptureOutputAndStatus(func() int {
						return action.Execute(operationID, mtaClient).ToInt()
					})
					ex.ExpectSuccessWithOutput(status, output, []string{
						"Process finished.",
						"Use \"cf dmol -i " + operationID + "\" to download the logs of the process.",
					})
				})
			})
			Context("when the operation fails", func() {
				It("should fail with an error and show the available actions", func() {
					var errorMessage = &models.Message{
						ID:   0,
						Type: "ERROR",
						Text: "Could not create application 'foo'",
					}
					var operation = &models.Operation{
						State:    "ERROR",
						Messages: []*models.Message{errorMessage},
					}
					mtaClient = fakes.NewFakeMtaClientBuilder().
						GetMtaOperation(operationID, "messages", operation, nil).
						Build()
					_, status := oc.CaptureOutputAndStatus(func() int {
						return action.Execute(operationID, mtaClient).ToInt()
					})
					ex.ExpectNonZeroStatus(status)
				})
			})
		})
	})

	Describe("GetNoRetriesActionToExecute", func() {
		Context("with correct action id", func() {
			It("should return abort action to execute", func() {
				actionToExecute := commands.GetNoRetriesActionToExecute("abort", "deploy")
				Expect(actionToExecute).NotTo(BeNil())
			})
			It("should return retry action to execute", func() {
				actionToExecute := commands.GetNoRetriesActionToExecute("retry", "deploy")
				Expect(actionToExecute).NotTo(BeNil())
			})
			It("should return resume action to execute", func() {
				actionToExecute := commands.GetNoRetriesActionToExecute("resume", "deploy")
				Expect(actionToExecute).NotTo(BeNil())
			})
			It("should return monitor action to execute", func() {
				actionToExecute := commands.GetNoRetriesActionToExecute("monitor", "deploy")
				Expect(actionToExecute).NotTo(BeNil())
			})
		})
		Context("with incorrect action id", func() {
			It("should return nil", func() {
				actionToExecute := commands.GetNoRetriesActionToExecute("test", "deploy")
				Expect(actionToExecute).To(BeNil())
			})
		})
	})
})
