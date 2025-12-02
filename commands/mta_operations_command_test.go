package commands_test

import (
	"fmt"

	plugin_fakes "code.cloudfoundry.org/cli/v8/plugin/pluginfakes"
	cli_fakes "github.com/cloudfoundry-incubator/multiapps-cli-plugin/cli/fakes"
	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/clients/models"
	mtafake "github.com/cloudfoundry-incubator/multiapps-cli-plugin/clients/mtaclient/fakes"
	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/commands"
	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/testutil"
	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/ui"
	util_fakes "github.com/cloudfoundry-incubator/multiapps-cli-plugin/util/fakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("MtaOperationsCommand", func() {
	Describe("Execute", func() {
		const org = "test-org"
		const space = "test-space"
		const user = "test-user"

		var name string
		var cliConnection *plugin_fakes.FakeCliConnection
		var clientFactory *commands.TestClientFactory
		var command *commands.MtaOperationsCommand
		var oc = testutil.NewUIOutputCapturer()
		var ex = testutil.NewUIExpector()

		var getOutputLines = func(operationsDetails [][]string) []string {
			lines := []string{}
			if len(operationsDetails) > 0 {
				lines = append(lines, testutil.GetTableOutputLines([]string{"id", "type", "mta id", "namespace", "status", "started at", "started by"}, operationsDetails)...)
			} else {
				lines = append(lines, "No multi-target app operations found")
			}

			return lines
		}

		BeforeEach(func() {
			ui.DisableTerminalOutput(true)
			name = command.GetPluginCommand().Name
			cliConnection = cli_fakes.NewFakeCliConnectionBuilder().
				CurrentOrg("test-org-guid", org, nil).
				CurrentSpace("test-space-guid", space, nil).
				Username(user, nil).
				AccessToken("bearer test-token", nil).Build()
			mtaClient := mtafake.NewFakeMtaClientBuilder().
				GetMta("test", nil, nil).Build()
			clientFactory = commands.NewTestClientFactory(mtaClient, nil, nil)
			command = commands.NewMtaOperationsCommand()
			testTokenFactory := commands.NewTestTokenFactory(cliConnection)
			deployServiceURLCalculator := util_fakes.NewDeployServiceURLFakeCalculator("deploy-service.test.ondemand.com")
			command.InitializeAll(name, cliConnection, testutil.NewCustomTransport(200), clientFactory, testTokenFactory, deployServiceURLCalculator)
		})

		// with an unknown flag - error
		Context("with an unknown flag", func() {
			It("should print incorrect usage, call cf help, and exit with a non-zero status", func() {
				output, status := oc.CaptureOutputAndStatus(func() int {
					return command.Execute([]string{"-a"}).ToInt()
				})
				ex.ExpectFailure(status, output, "Incorrect usage. Unknown or wrong flags: -a")
				Expect(cliConnection.CliCommandArgsForCall(0)).To(Equal([]string{"help", name}))
			})
		})

		// with an unknown flag and one valid flag
		Context("with an unknown flag and one valid flag", func() {
			It("should print incorrect usage, call cf help, and exit with a non-zero status", func() {
				output, status := oc.CaptureOutputAndStatus(func() int {
					return command.Execute([]string{"-a", "--last"}).ToInt()
				})
				ex.ExpectFailure(status, output, "Incorrect usage. Unknown or wrong flags: -a")
				Expect(cliConnection.CliCommandArgsForCall(0)).To(Equal([]string{"help", name}))
			})
		})

		// wrong arguments
		Context("with wrong arguments", func() {
			It("should print incorrect usage, call cf help, and exit with a non-zero status", func() {
				output, status := oc.CaptureOutputAndStatus(func() int {
					return command.Execute([]string{"x", "y", "z"}).ToInt()
				})
				ex.ExpectFailure(status, output, "Incorrect usage. Wrong arguments")
				Expect(cliConnection.CliCommandArgsForCall(0)).To(Equal([]string{"help", name}))
			})
		})

		// can't connect to backend - error
		Context("when can't connect to backend", func() {
			const host = "x"
			It("should print an error and exit with a non-zero status", func() {
				clientFactory.MtaClient = mtafake.NewFakeMtaClientBuilder().
					GetMtaOperations(nil, nil, nil, []*models.Operation{}, fmt.Errorf("Get https://%s/rest/test/test/mta: dial tcp: lookup %s: no such host", host, host)).Build()
				output, status := oc.CaptureOutputAndStatus(func() int {
					return command.Execute([]string{"-u", host}).ToInt()
				})
				ex.ExpectFailureOnLine(status, output, "Could not get multi-target app operations:", 2)
			})
		})

		// backend returns an an error response - error
		Context("with an error response returned by the backend", func() {
			It("should print an error and exit with a non-zero status", func() {
				clientFactory.MtaClient = mtafake.NewFakeMtaClientBuilder().
					GetMtaOperations(nil, nil, nil, []*models.Operation{}, fmt.Errorf("unknown error (status 404)")).Build()
				output, status := oc.CaptureOutputAndStatus(func() int {
					return command.Execute([]string{}).ToInt()
				})
				ex.ExpectFailureOnLine(status, output, "Could not get multi-target app operations:", 2)
			})
		})

		Context("with empty response returned by the backend", func() {
			It("should print info and return with zero status", func() {
				clientFactory.MtaClient = mtafake.NewFakeMtaClientBuilder().
					GetMtaOperations(nil, nil, nil, []*models.Operation{}, nil).Build()
				output, status := oc.CaptureOutputAndStatus(func() int {
					return command.Execute([]string{}).ToInt()
				})
				expectedOutput := getOutputLines([][]string{})
				expectedOutput = append([]string{
					"Getting active multi-target app operations in org test-org / space test-space as test-user...",
					"OK",
				}, expectedOutput...)
				ex.ExpectSuccessWithOutput(status, output, expectedOutput)
			})
		})
		Context("with non-empty response returned by the backend", func() {
			It("should print info and return with zero status", func() {
				clientFactory.MtaClient = mtafake.NewFakeMtaClientBuilder().
					GetMtaOperations(nil, nil, nil, []*models.Operation{
						testutil.GetOperation("111", "test-space", "test", "namespace", "deploy", "ERROR", false)}, nil).Build()
				output, status := oc.CaptureOutputAndStatus(func() int {
					return command.Execute([]string{}).ToInt()
				})
				expectedOutput := getOutputLines([][]string{
					{"111", "deploy", "test", "namespace", "ERROR", "2016-03-04T14:23:24.521Z[Etc/UTC]", "admin"},
				})
				expectedOutput = append([]string{
					"Getting active multi-target app operations in org test-org / space test-space as test-user...",
					"OK",
				}, expectedOutput...)
				ex.ExpectSuccessWithOutput(status, output, expectedOutput)
			})
		})
		Context("with non-empty response returned by the backend containing a nil value for MTA ID", func() {
			It("should print info and return with zero status", func() {
				clientFactory.MtaClient = mtafake.NewFakeMtaClientBuilder().
					GetMtaOperations(nil, nil, nil, []*models.Operation{
						testutil.GetOperation("111", "test-space", "", "", "deploy", "ERROR", false)}, nil).Build()
				output, status := oc.CaptureOutputAndStatus(func() int {
					return command.Execute([]string{}).ToInt()
				})
				expectedOutput := getOutputLines([][]string{
					{"111", "deploy", "N/A", "", "ERROR", "2016-03-04T14:23:24.521Z[Etc/UTC]", "admin"},
				})
				expectedOutput = append([]string{
					"Getting active multi-target app operations in org test-org / space test-space as test-user...",
					"OK",
				}, expectedOutput...)
				ex.ExpectSuccessWithOutput(status, output, expectedOutput)
			})
		})
		Context("with more than 1 operations returned by the backend", func() {
			It("should print info and return with zero status", func() {
				clientFactory.MtaClient = mtafake.NewFakeMtaClientBuilder().
					GetMtaOperations(nil, nil, nil, []*models.Operation{
						testutil.GetOperation("test-1", "test-space", "test-mta-1", "namespace", "deploy", "ERROR", true),
						testutil.GetOperation("test-2", "test-space", "test-mta-2", "namespace", "deploy", "ERROR", false),
					}, nil).Build()
				output, status := oc.CaptureOutputAndStatus(func() int {
					return command.Execute([]string{}).ToInt()
				})
				expectedOutput := getOutputLines([][]string{
					{"test-1", "deploy", "test-mta-1", "namespace", "ERROR", "2016-03-04T14:23:24.521Z[Etc/UTC]", "admin"},
					{"test-2", "deploy", "test-mta-2", "namespace", "ERROR", "2016-03-04T14:23:24.521Z[Etc/UTC]", "admin"},
				})
				expectedOutput = append([]string{
					"Getting active multi-target app operations in org test-org / space test-space as test-user...",
					"OK",
				}, expectedOutput...)
				ex.ExpectSuccessWithOutput(status, output, expectedOutput)
			})
		})
		Context("with more than 1 operations returned by the backend and last option provided", func() {
			It("should print the info of the last 2 operations and return with zero status", func() {
				clientFactory.MtaClient = mtafake.NewFakeMtaClientBuilder().
					GetMtaOperations(nil, &[]int64{2}[0], nil, []*models.Operation{
						testutil.GetOperation("test-2", "test-space", "test-mta-2", "namespace", "deploy", "ERROR", false),
						testutil.GetOperation("test-3", "test-space", "test-mta-3", "namespace", "deploy", "ERROR", false),
					}, nil).Build()
				output, status := oc.CaptureOutputAndStatus(func() int {
					return command.Execute([]string{"-last", "2"}).ToInt()
				})
				expectedOutput := getOutputLines([][]string{
					{"test-2", "deploy", "test-mta-2", "namespace", "ERROR", "2016-03-04T14:23:24.521Z[Etc/UTC]", "admin"},
					{"test-3", "deploy", "test-mta-3", "namespace", "ERROR", "2016-03-04T14:23:24.521Z[Etc/UTC]", "admin"},
				})
				expectedOutput = append([]string{
					"Getting last 2 multi-target app operations in org test-org / space test-space as test-user...",
					"OK",
				}, expectedOutput...)
				ex.ExpectSuccessWithOutput(status, output, expectedOutput)
			})
		})
		Context("with more than 1 operations returned by the backend and last option provided", func() {
			It("should print the info for all of the operations and return with zero status", func() {
				clientFactory.MtaClient = mtafake.NewFakeMtaClientBuilder().
					GetMtaOperations(nil, &[]int64{10}[0], nil, []*models.Operation{
						testutil.GetOperation("test-1", "test-space", "test-mta-1", "namespace", "deploy", "ERROR", true),
						testutil.GetOperation("test-2", "test-space", "test-mta-2", "namespace", "deploy", "ERROR", false),
						testutil.GetOperation("test-3", "test-space", "test-mta-3", "namespace", "deploy", "ERROR", false),
					}, nil).Build()
				output, status := oc.CaptureOutputAndStatus(func() int {
					return command.Execute([]string{"-last", "10"}).ToInt()
				})
				expectedOutput := getOutputLines([][]string{
					[]string{"test-1", "deploy", "test-mta-1", "namespace", "ERROR", "2016-03-04T14:23:24.521Z[Etc/UTC]", "admin"},
					[]string{"test-2", "deploy", "test-mta-2", "namespace", "ERROR", "2016-03-04T14:23:24.521Z[Etc/UTC]", "admin"},
					[]string{"test-3", "deploy", "test-mta-3", "namespace", "ERROR", "2016-03-04T14:23:24.521Z[Etc/UTC]", "admin"},
				})
				expectedOutput = append([]string{
					"Getting last 10 multi-target app operations in org test-org / space test-space as test-user...",
					"OK",
				}, expectedOutput...)
				ex.ExpectSuccessWithOutput(status, output, expectedOutput)
			})
		})
		Context("with more than 1 operations returned by the backend and last option provided", func() {
			It("should print the info for the last operation and return with zero status", func() {
				clientFactory.MtaClient = mtafake.NewFakeMtaClientBuilder().
					GetMtaOperations(nil, &[]int64{10}[0], nil, []*models.Operation{
						testutil.GetOperation("test-1", "test-space", "test-mta-1", "namespace", "deploy", "ERROR", true),
						testutil.GetOperation("test-2", "test-space", "test-mta-2", "namespace", "deploy", "ERROR", false),
						testutil.GetOperation("test-3", "test-space", "test-mta-3", "namespace", "deploy", "ERROR", false),
					}, nil).Build()
				output, status := oc.CaptureOutputAndStatus(func() int {
					return command.Execute([]string{"-last", "1"}).ToInt()
				})
				expectedOutput := getOutputLines([][]string{
					{"test-1", "deploy", "test-mta-1", "namespace", "ERROR", "2016-03-04T14:23:24.521Z[Etc/UTC]", "admin"},
					{"test-2", "deploy", "test-mta-2", "namespace", "ERROR", "2016-03-04T14:23:24.521Z[Etc/UTC]", "admin"},
					{"test-3", "deploy", "test-mta-3", "namespace", "ERROR", "2016-03-04T14:23:24.521Z[Etc/UTC]", "admin"},
				})
				expectedOutput = append([]string{
					"Getting last multi-target app operation in org test-org / space test-space as test-user...",
					"OK",
				}, expectedOutput...)
				ex.ExpectSuccessWithOutput(status, output, expectedOutput)
			})
		})
		Context("with empty response returned by the backend and last option provided", func() {
			It("should print the info for all of the operations and return with zero status", func() {
				clientFactory.MtaClient = mtafake.NewFakeMtaClientBuilder().
					GetMtaOperations(nil, &[]int64{10}[0], nil, []*models.Operation{}, nil).Build()
				output, status := oc.CaptureOutputAndStatus(func() int {
					return command.Execute([]string{"-last", "10"}).ToInt()
				})
				expectedOutput := getOutputLines([][]string{})
				expectedOutput = append([]string{
					"Getting last 10 multi-target app operations in org test-org / space test-space as test-user...",
					"OK",
				}, expectedOutput...)
				ex.ExpectSuccessWithOutput(status, output, expectedOutput)
			})
		})
		Context("with more than 1 operations returned by the backend and no options provided", func() {
			It("should print the info for operations in active state and return with zero status", func() {
				clientFactory.MtaClient = mtafake.NewFakeMtaClientBuilder().
					GetMtaOperations(nil, &[]int64{10}[0], nil, []*models.Operation{
						testutil.GetOperation("test-1", "test-space", "test-mta-1", "namespace", "deploy", "ERROR", true),
						testutil.GetOperation("test-2", "test-space", "test-mta-2", "namespace", "deploy", "RUNNING", false),
						testutil.GetOperation("test-3", "test-space", "test-mta-3", "namespace", "deploy", "ERROR", false),
					}, nil).Build()
				output, status := oc.CaptureOutputAndStatus(func() int {
					return command.Execute([]string{}).ToInt()
				})
				expectedOutput := getOutputLines([][]string{
					{"test-1", "deploy", "test-mta-1", "namespace", "ERROR", "2016-03-04T14:23:24.521Z[Etc/UTC]", "admin"},
					{"test-2", "deploy", "test-mta-2", "namespace", "RUNNING", "2016-03-04T14:23:24.521Z[Etc/UTC]", "admin"},
					{"test-3", "deploy", "test-mta-3", "namespace", "ERROR", "2016-03-04T14:23:24.521Z[Etc/UTC]", "admin"},
				})
				expectedOutput = append([]string{
					"Getting active multi-target app operations in org test-org / space test-space as test-user...",
					"OK",
				}, expectedOutput...)
				ex.ExpectSuccessWithOutput(status, output, expectedOutput)
			})
		})
		Context("with more than 1 operations returned by the backend and no options provided", func() {
			It("should print the info for operations in active state, not include operations in finished state and return with zero status", func() {
				clientFactory.MtaClient = mtafake.NewFakeMtaClientBuilder().
					GetMtaOperations(nil, nil, []string{"SLP_TASK_STATE_ERROR", "SLP_TASK_STATE_RUNNING"}, []*models.Operation{
						testutil.GetOperation("test-1", "test-space", "test-mta-1", "namespace", "deploy", "ERROR", true),
						testutil.GetOperation("test-2", "test-space", "test-mta-2", "namespace", "deploy", "RUNNING", false),
					}, nil).Build()
				output, status := oc.CaptureOutputAndStatus(func() int {
					return command.Execute([]string{}).ToInt()
				})
				expectedOutput := getOutputLines([][]string{
					{"test-1", "deploy", "test-mta-1", "namespace", "ERROR", "2016-03-04T14:23:24.521Z[Etc/UTC]", "admin"},
					{"test-2", "deploy", "test-mta-2", "namespace", "RUNNING", "2016-03-04T14:23:24.521Z[Etc/UTC]", "admin"},
				})
				expectedOutput = append([]string{
					"Getting active multi-target app operations in org test-org / space test-space as test-user...",
					"OK",
				}, expectedOutput...)
				ex.ExpectSuccessWithOutput(status, output, expectedOutput)
			})
		})
		Context("with more than 1 operations returned by the backend and all option provided", func() {
			It("should print the info for operations in active and finished state and return with zero status", func() {
				clientFactory.MtaClient = mtafake.NewFakeMtaClientBuilder().
					GetMtaOperations(nil, nil, nil, []*models.Operation{
						testutil.GetOperation("test-1", "test-space", "test-mta-1", "namespace", "deploy", "ERROR", true),
						testutil.GetOperation("test-2", "test-space", "test-mta-2", "namespace", "deploy", "RUNNING", false),
						testutil.GetOperation("test-3", "test-space", "test-mta-3", "namespace", "deploy", "FINISHED", false),
					}, nil).Build()
				output, status := oc.CaptureOutputAndStatus(func() int {
					return command.Execute([]string{"-all"}).ToInt()
				})
				expectedOutput := getOutputLines([][]string{
					{"test-1", "deploy", "test-mta-1", "namespace", "ERROR", "2016-03-04T14:23:24.521Z[Etc/UTC]", "admin"},
					{"test-2", "deploy", "test-mta-2", "namespace", "RUNNING", "2016-03-04T14:23:24.521Z[Etc/UTC]", "admin"},
					{"test-3", "deploy", "test-mta-3", "namespace", "FINISHED", "2016-03-04T14:23:24.521Z[Etc/UTC]", "admin"},
				})
				expectedOutput = append([]string{
					"Getting all multi-target app operations in org test-org / space test-space as test-user...",
					"OK",
				}, expectedOutput...)
				ex.ExpectSuccessWithOutput(status, output, expectedOutput)
			})
		})
		Context("with mta id provided", func() {
			It("should print the info for the operation and return with zero status", func() {
				clientFactory.MtaClient = mtafake.NewFakeMtaClientBuilder().
					GetMtaOperations(&[]string{"test-mta-id"}[0], &[]int64{1}[0], nil, []*models.Operation{
						testutil.GetOperation("test-1", "test-space", "test-mta-id", "namespace", "deploy", "ERROR", true),
					}, nil).Build()
				output, status := oc.CaptureOutputAndStatus(func() int {
					return command.Execute([]string{"--mta", "test-mta-id"}).ToInt()
				})
				expectedOutput := getOutputLines([][]string{
					{"test-1", "deploy", "test-mta-id", "namespace", "ERROR", "2016-03-04T14:23:24.521Z[Etc/UTC]", "admin"},
				})
				expectedOutput = append([]string{
					"Getting multi-target app operations for test-mta-id in org test-org / space test-space as test-user...",
					"OK",
				}, expectedOutput...)
				ex.ExpectSuccessWithOutput(status, output, expectedOutput)
			})
		})
	})
})
