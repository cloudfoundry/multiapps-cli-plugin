package commands_test

import (
	"fmt"

	. "github.com/onsi/ginkgo"
	"github.com/SAP/cf-mta-plugin/clients/models"
	slppclient "github.com/SAP/cf-mta-plugin/clients/slppclient"
	fakes "github.com/SAP/cf-mta-plugin/clients/slppclient/fakes"
	"github.com/SAP/cf-mta-plugin/commands"
	"github.com/SAP/cf-mta-plugin/testutil"
	"github.com/SAP/cf-mta-plugin/ui"
)

var _ = Describe("ExecutionMonitor", func() {
	Describe("Monitor", func() {
		var monitor *commands.ExecutionMonitor
		var client slppclient.SlppClientOperations
		processID := "1234"
		commandName := "deploy"
		var oc = testutil.NewUIOutputCapturer()
		var ex = testutil.NewUIExpector()

		fakeSlppClientBuilder := fakes.NewFakeSlppClientBuilder()

		var getOutputLines = func(processStatus models.SlpTaskState, errorMessage string, progressMessages []string) []string {
			var lines []string
			lines = append(lines, "Monitoring process execution...\n")
			if len(progressMessages) > 0 {
				lines = append(lines, progressMessages...)
			}
			switch processStatus {
			case models.SlpTaskStateSlpTaskStateFINISHED:
				lines = append(lines, "Process finished.\n")
			case models.SlpTaskStateSlpTaskStateABORTED:
				lines = append(lines, "Process was aborted.\n")
			case models.SlpTaskStateSlpTaskStateACTIONREQUIRED, models.SlpTaskStateSlpTaskStateDIALOG:
				lines = append(lines, "Process has entered validation phase. After testing your new deployment you can resume or abort the process.\n")
				lines = append(lines, fmt.Sprintf("Use \"cf %s -i %s -a resume\" to resume the process.\n", commandName, processID))
				lines = append(lines, fmt.Sprintf("Use \"cf %s -i %s -a abort\" to abort the process.\n", commandName, processID))
				lines = append(lines, "Hint: Use the '--no-confirm' option of the bg-deploy command to skip this phase.\n")
			default:
				lines = append(lines, fmt.Sprintf("Process is in illegal state %s.", processStatus))
			}

			if errorMessage != "" {
				lines = append(lines, fmt.Sprintf("Process failed: %s\n", errorMessage))
				lines = append(lines, fmt.Sprintf("Use \"cf %s -i %s -a retry\" to retry the process.\n", commandName, processID))
				lines = append(lines, fmt.Sprintf("Use \"cf %s -i %s -a abort\" to abort the process.\n", commandName, processID))
			}
			return lines
		}

		BeforeEach(func() {
			ui.DisableTerminalOutput(true)
		})

		Context("with no process task in the tasklist", func() {
			It("should return error", func() {
				tasklist := testutil.GetTaskList(testutil.GetTask("test-type", "test-state", []*models.ProgressMessage{}))
				client = fakeSlppClientBuilder.
					GetTasklistTask(tasklist.Tasks[0], nil).
					Build()
				monitor = commands.NewExecutionMonitor(processID, commandName, client, []*models.ProgressMessage{})
				output, status := oc.CaptureOutputAndStatus(func() int {
					return monitor.Monitor().ToInt()
				})
				ex.ExpectFailureOnLine(status, output, "The SLP task must be a Process task", 1)
			})
		})

		Context("with process task in state aborted and no progress messages in the tasklist", func() {
			It("should print info message and exit with zero status", func() {
				tasklist := testutil.GetTaskList(testutil.GetTask(models.SlpTaskTypeSlpTaskTypePROCESS, models.SlpTaskStateSlpTaskStateABORTED, []*models.ProgressMessage{}))
				client = fakeSlppClientBuilder.
					GetTasklistTask(tasklist.Tasks[0], nil).Build()
				monitor = commands.NewExecutionMonitor(processID, commandName, client, []*models.ProgressMessage{})
				output, status := oc.CaptureOutputAndStatus(func() int {
					return monitor.Monitor().ToInt()
				})
				ex.ExpectFailure(status, output, "Process was aborted.")
			})
		})
		Context("with process task in state error and no progress messages in the tasklist", func() {
			It("should return error and exit with non-zero status", func() {
				tasklist := testutil.GetTaskList(testutil.GetTask(models.SlpTaskTypeSlpTaskTypePROCESS, models.SlpTaskStateSlpTaskStateERROR, []*models.ProgressMessage{}))
				client = fakeSlppClientBuilder.
					GetTasklistTask(tasklist.Tasks[0], nil).
					GetError(testutil.ErrorResult, nil).GetActions(testutil.Actions, nil).Build()
				monitor = commands.NewExecutionMonitor(processID, commandName, client, []*models.ProgressMessage{})
				output, status := oc.CaptureOutputAndStatus(func() int {
					return monitor.Monitor().ToInt()
				})
				ex.ExpectFailureOnLine(status, output, "Use \"cf deploy -i 1234 -a abort\" to abort the process.\n", 2)
			})
		})
		Context("with process task in illegal state and no progress messages in the tasklist", func() {
			It("should return error and exit with non-zero status", func() {
				tasklist := testutil.GetTaskList(testutil.GetTask(models.SlpTaskTypeSlpTaskTypePROCESS, "slp.task.state.NON-EXISTING", []*models.ProgressMessage{}))
				client = fakeSlppClientBuilder.GetTasklistTask(tasklist.Tasks[0], nil).Build()
				monitor = commands.NewExecutionMonitor(processID, commandName, client, []*models.ProgressMessage{})
				output, status := oc.CaptureOutputAndStatus(func() int {
					return monitor.Monitor().ToInt()
				})
				ex.ExpectFailureOnLine(status, output, "Process is in illegal state slp.task.state.NON-EXISTING", 1)
			})
		})
		Context("with process task in state finished and no progress messages in the tasklist", func() {
			It("should print info message and exit with zero status", func() {
				const processStatus = models.SlpTaskStateSlpTaskStateFINISHED
				tasklist := testutil.GetTaskList(testutil.GetTask(models.SlpTaskTypeSlpTaskTypePROCESS, processStatus, []*models.ProgressMessage{}))
				client = fakeSlppClientBuilder.GetTasklistTask(tasklist.Tasks[0], nil).Build()
				monitor = commands.NewExecutionMonitor(processID, commandName, client, []*models.ProgressMessage{})
				output, status := oc.CaptureOutputAndStatus(func() int {
					return monitor.Monitor().ToInt()
				})
				ex.ExpectSuccessWithOutput(status, output, getOutputLines(processStatus, "", []string{}))
			})
		})

		Context("with process task in state finished and progress messages with non-repeating ids in the tasklist", func() {
			It("should print all progress messages and exit with zero status", func() {
				const processStatus = models.SlpTaskStateSlpTaskStateFINISHED
				tasklist := testutil.GetTaskList(testutil.GetTask(models.SlpTaskTypeSlpTaskTypePROCESS, processStatus, []*models.ProgressMessage{
					testutil.GetProgressMessage("1", "test-message-1"),
					testutil.GetProgressMessage("2", "test-message-2"),
					testutil.GetProgressMessage("3", "test-message-3"),
				}))
				client = fakeSlppClientBuilder.GetTasklistTask(tasklist.Tasks[0], nil).Build()
				monitor = commands.NewExecutionMonitor(processID, commandName, client, []*models.ProgressMessage{})
				output, status := oc.CaptureOutputAndStatus(func() int {
					return monitor.Monitor().ToInt()
				})
				ex.ExpectSuccessWithOutput(status, output, getOutputLines(processStatus, "", []string{"  test-message-1\n", "  test-message-2\n", "  test-message-3\n"}))
			})
		})
		Context("with process task in state finished and progress messages with repeating ids in the tasklist", func() {
			It("should print progress messages and exit with zero status", func() {
				const processStatus = models.SlpTaskStateSlpTaskStateFINISHED
				tasklist := testutil.GetTaskList(testutil.GetTask(models.SlpTaskTypeSlpTaskTypePROCESS, processStatus, []*models.ProgressMessage{
					testutil.GetProgressMessage("1", "test-message-1"),
					testutil.GetProgressMessage("1", "test-message-2"),
					testutil.GetProgressMessage("3", "test-message-3"),
					testutil.GetProgressMessage("4", "test-message-4"),
				}))
				client = fakeSlppClientBuilder.GetTasklistTask(tasklist.Tasks[0], nil).Build()
				monitor = commands.NewExecutionMonitor(processID, commandName, client, []*models.ProgressMessage{})
				output, status := oc.CaptureOutputAndStatus(func() int {
					return monitor.Monitor().ToInt()
				})
				ex.ExpectSuccessWithOutput(status, output, getOutputLines(processStatus, "", []string{"  test-message-1\n", "  test-message-3\n", "  test-message-4\n"}))
			})
		})
		Context("with process task in state finished and progress messages with non-repeating ids in the tasklist and already reported progress messages", func() {
			It("should print progress messages and exit with zero status", func() {
				const processStatus = models.SlpTaskStateSlpTaskStateFINISHED
				tasklist := testutil.GetTaskList(testutil.GetTask(models.SlpTaskTypeSlpTaskTypePROCESS, processStatus, []*models.ProgressMessage{
					testutil.GetProgressMessage("1", "test-message-1"),
					testutil.GetProgressMessage("2", "test-message-2"),
					testutil.GetProgressMessage("3", "test-message-3"),
					testutil.GetProgressMessage("4", "test-message-4"),
				}))
				client = fakeSlppClientBuilder.GetTasklistTask(tasklist.Tasks[0], nil).Build()
				monitor = commands.NewExecutionMonitor(processID, commandName, client, []*models.ProgressMessage{
					testutil.GetProgressMessage("1", "test-message-1"),
					testutil.GetProgressMessage("2", "test-message-2")})
				output, status := oc.CaptureOutputAndStatus(func() int {
					return monitor.Monitor().ToInt()
				})
				ex.ExpectSuccessWithOutput(status, output, getOutputLines(processStatus, "", []string{"  test-message-3\n", "  test-message-4\n"}))
			})
		})
		Context("with process task in state finished and already reported progress messages which are not part of the reported progressmessages", func() {
			It("should print progress messages and exit with zero status", func() {
				const processStatus = models.SlpTaskStateSlpTaskStateFINISHED
				client = fakeSlppClientBuilder.
					GetTasklist(testutil.GetTaskList(testutil.GetTask(models.SlpTaskTypeSlpTaskTypePROCESS, processStatus, []*models.ProgressMessage{
						testutil.GetProgressMessage("1", "test-message-1"),
						testutil.GetProgressMessage("2", "test-message-2"),
						testutil.GetProgressMessage("3", "test-message-3"),
						testutil.GetProgressMessage("4", "test-message-4"),
					})), nil).Build()
				monitor = commands.NewExecutionMonitor(processID, commandName, client, []*models.ProgressMessage{
					testutil.GetProgressMessage("5", "test-message-5"),
					testutil.GetProgressMessage("6", "test-message-6")})
				output, status := oc.CaptureOutputAndStatus(func() int {
					return monitor.Monitor().ToInt()
				})
				ex.ExpectSuccessWithOutput(status, output, getOutputLines(processStatus, "", []string{"  test-message-1\n", "  test-message-2\n", "  test-message-3\n", "  test-message-4\n"}))
			})
		})
		Context("with process task in state action required", func() {
			It("should print progress messages and exit with zero status", func() {
				const processStatus = models.SlpTaskStateSlpTaskStateACTIONREQUIRED
				tasklist := testutil.GetTaskList(testutil.GetTask(models.SlpTaskTypeSlpTaskTypePROCESS, processStatus, []*models.ProgressMessage{}))
				client = fakeSlppClientBuilder.
					GetTasklistTask(tasklist.Tasks[0], nil).
					GetActions(testutil.BlueGreenActions, nil).Build()
				monitor = commands.NewExecutionMonitor(processID, commandName, client, []*models.ProgressMessage{})
				output, status := oc.CaptureOutputAndStatus(func() int {
					return monitor.Monitor().ToInt()
				})
				ex.ExpectSuccessWithOutput(status, output, getOutputLines(processStatus, "", []string{}))
			})
		})
	})
})
