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
)

var _ = Describe("ExecutionMonitor", func() {
	Describe("Monitor", func() {
		var monitor *commands.ExecutionMonitor
		var client mtaclient.MtaClientOperations
		processID := "1234"
		commandName := "deploy"
		var oc = testutil.NewUIOutputCapturer()
		var ex = testutil.NewUIExpector()

		fakeMtaClientBuilder := fakes.NewFakeMtaClientBuilder()

		var getOutputLines = func(processStatus models.State, errorMessage string, progressMessages []string) []string {
			var lines []string
			if len(progressMessages) > 0 {
				lines = append(lines, progressMessages...)
			}
			switch processStatus {
			case models.StateFINISHED:
				lines = append(lines, "Process finished.")
				lines = append(lines, fmt.Sprintf("Use \"cf dmol -i %s\" to download the logs of the process.", processID))
			case models.StateABORTED:
				lines = append(lines, "Process was aborted.")
				lines = append(lines, fmt.Sprintf("Use \"cf dmol -i %s\" to download the logs of the process.", processID))
			case models.StateACTIONREQUIRED:
				lines = append(lines, "Process has entered validation phase. After testing your new deployment you can resume or abort the process.")
				lines = append(lines, fmt.Sprintf("Use \"cf %s -i %s -a resume\" to resume the process.", commandName, processID))
				lines = append(lines, fmt.Sprintf("Use \"cf %s -i %s -a abort\" to abort the process.", commandName, processID))
				lines = append(lines, "Hint: Use the '--no-confirm' option of the bg-deploy command to skip this phase.")
			default:
				lines = append(lines, fmt.Sprintf("Process is in illegal state %s.", processStatus))
			}

			if errorMessage != "" {
				lines = append(lines, fmt.Sprintf("Process failed."))
				lines = append(lines, fmt.Sprintf("Use \"cf %s -i %s -a retry\" to retry the process.", commandName, processID))
				lines = append(lines, fmt.Sprintf("Use \"cf %s -i %s -a abort\" to abort the process.", commandName, processID))
				lines = append(lines, fmt.Sprintf("Use \"cf dmol -i %s\" to download the logs of the process.", processID))
			}
			return lines
		}

		BeforeEach(func() {
			ui.DisableTerminalOutput(true)
		})

		Context("with process task in state aborted", func() {
			It("should print info message and exit with non-zero status", func() {
				client = fakeMtaClientBuilder.
					GetMtaOperation(processID, "messages", &models.Operation{
						State:    "ABORTED",
						Messages: []*models.Message{},
					}, nil).Build()
				monitor = commands.NewExecutionMonitor(commandName, processID, "messages", 0, []*models.Message{}, client)
				output, status := oc.CaptureOutputAndStatus(func() int {
					return monitor.Monitor().ToInt()
				})
				ex.ExpectNonZeroStatus(status)
				ex.ExpectMessageOnLine(output, "Process was aborted.", 0)
			})
		})
		Context("with process task in state error and no progress messages in the tasklist, retrying 4 times", func() {
			It("should retry 4 times, return error and exit with non-zero status", func() {
				client = fakeMtaClientBuilder.
					GetMtaOperation(processID, "messages", &models.Operation{
						ProcessID: processID,
						State:     "ERROR",
						Messages: []*models.Message{
							&models.Message{
								Type: models.MessageTypeERROR,
								Text: "error message",
							},
						},
					}, nil).
					GetOperationActions(processID, []string{"abort", "retry"}, nil).Build()
				monitor = commands.NewExecutionMonitor(commandName, processID, "messages", 4, []*models.Message{}, client)
				output, status := oc.CaptureOutputAndStatus(func() int {
					return monitor.Monitor().ToInt()
				})
				ex.ExpectMessageOnLine(output, "Proceeding with automatic retry... (4 of 4 attempts left)", 1)
				ex.ExpectMessageOnLine(output, "Proceeding with automatic retry... (3 of 4 attempts left)", 2)
				ex.ExpectMessageOnLine(output, "Proceeding with automatic retry... (2 of 4 attempts left)", 3)
				ex.ExpectMessageOnLine(output, "Proceeding with automatic retry... (1 of 4 attempts left)", 4)
				ex.ExpectFailureOnLine(status, output, "Use \"cf deploy -i 1234 -a abort\" to abort the process.", 6)
			})
		})
		Context("with process task in illegal state and no progress messages in the tasklist", func() {
			It("should return error and exit with non-zero status", func() {
				client = fakeMtaClientBuilder.
					GetMtaOperation(processID, "messages", &models.Operation{
						State:    "UnknownState",
						Messages: []*models.Message{},
					}, nil).Build()
				monitor = commands.NewExecutionMonitor(commandName, processID, "messages", 0, []*models.Message{}, client)
				output, status := oc.CaptureOutputAndStatus(func() int {
					return monitor.Monitor().ToInt()
				})
				ex.ExpectFailure(status, output, "Process is in illegal state UnknownState")
			})
		})
		Context("with process task in state finished and no progress messages in the tasklist", func() {
			It("should print info message and exit with zero status", func() {
				const processStatus = models.StateFINISHED
				client = fakeMtaClientBuilder.
					GetMtaOperation(processID, "messages", &models.Operation{
						State:    "FINISHED",
						Messages: []*models.Message{},
					}, nil).Build()
				monitor = commands.NewExecutionMonitor(commandName, processID, "messages", 0, []*models.Message{}, client)
				output, status := oc.CaptureOutputAndStatus(func() int {
					return monitor.Monitor().ToInt()
				})
				ex.ExpectSuccessWithOutput(status, output, getOutputLines(processStatus, "", []string{}))
			})
		})

		Context("with process task in state finished and progress messages with non-repeating ids in the tasklist", func() {
			It("should print all progress messages and exit with zero status", func() {
				const processStatus = models.StateFINISHED
				client = fakeMtaClientBuilder.
					GetMtaOperation(processID, "messages", &models.Operation{
						State: "FINISHED",
						Messages: []*models.Message{
							testutil.GetMessage(1, "test-message-1"),
							testutil.GetMessage(2, "test-message-2"),
							testutil.GetMessage(31, "test-message-3"),
						},
					}, nil).Build()
				monitor = commands.NewExecutionMonitor(commandName, processID, "messages", 0, []*models.Message{}, client)
				output, status := oc.CaptureOutputAndStatus(func() int {
					return monitor.Monitor().ToInt()
				})
				ex.ExpectSuccessWithOutput(status, output, getOutputLines(processStatus, "", []string{"test-message-1", "test-message-2", "test-message-3"}))
			})
		})
		Context("with process task in state finished and progress messages with repeating ids in the tasklist", func() {
			It("should print progress messages and exit with zero status", func() {
				const processStatus = models.StateFINISHED
				client = fakeMtaClientBuilder.
					GetMtaOperation(processID, "messages", &models.Operation{
						State: "FINISHED",
						Messages: []*models.Message{
							testutil.GetMessage(1, "test-message-1"),
							testutil.GetMessage(1, "test-message-2"),
							testutil.GetMessage(3, "test-message-3"),
							testutil.GetMessage(4, "test-message-4"),
						},
					}, nil).Build()
				monitor = commands.NewExecutionMonitor(commandName, processID, "messages", 0, []*models.Message{}, client)
				output, status := oc.CaptureOutputAndStatus(func() int {
					return monitor.Monitor().ToInt()
				})
				ex.ExpectSuccessWithOutput(status, output, getOutputLines(processStatus, "", []string{"test-message-1", "test-message-3", "test-message-4"}))
			})
		})
		Context("with process task in state action required", func() {
			It("should print progress messages and exit with zero status", func() {
				const processStatus = models.StateFINISHED
				client = fakeMtaClientBuilder.
					GetMtaOperation(processID, "messages", &models.Operation{
						State:    "FINISHED",
						Messages: []*models.Message{},
					}, nil).
					GetOperationActions(processID, []string{"retry", "abort"}, nil).Build()
				monitor = commands.NewExecutionMonitorFromLocationHeader(commandName, "operations/"+processID+"?embed=messages", 0, []*models.Message{}, client)
				output, status := oc.CaptureOutputAndStatus(func() int {
					return monitor.Monitor().ToInt()
				})
				ex.ExpectSuccessWithOutput(status, output, getOutputLines(processStatus, "", []string{}))
			})
		})
	})
})
