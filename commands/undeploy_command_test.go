package commands_test

import (
	"fmt"
	cliFakes "github.com/cloudfoundry-incubator/multiapps-cli-plugin/cli/fakes"
	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/clients/baseclient"
	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/clients/models"
	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/clients/mtaclient"
	mtaFake "github.com/cloudfoundry-incubator/multiapps-cli-plugin/clients/mtaclient/fakes"
	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/ui"

	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/commands"
	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/configuration"
	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/testutil"
	utilFakes "github.com/cloudfoundry-incubator/multiapps-cli-plugin/util/fakes"
	pluginFakes "github.com/cloudfoundry/cli/plugin/fakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("UndeployCommand", func() {
	Describe("Execute", func() {
		const org = "test-org"
		const space = "test-space"
		const spaceId = "test-space-guid"
		const user = "test-user"
		const mtaID = "test"
		const ongoingOperationId = "999"

		var name string
		var cliConnection *pluginFakes.FakeCliConnection
		var mtaClient mtaFake.FakeMtaClientOperations
		var testClientFactory *commands.TestClientFactory
		var command *commands.UndeployCommand
		var oc = testutil.NewUIOutputCapturer()
		var ex = testutil.NewUIExpector()

		var undeployOperation = models.Operation{
			State:       "FINISHED",
			ProcessID:   testutil.ProcessID,
			ProcessType: "UNDEPLOY",
			Messages:    []*models.Message{&testutil.SimpleMessage},
		}

		var ongoingOperation = models.Operation{
			AcquiredLock: true,
			Messages:     []*models.Message{&testutil.SimpleMessage},
			MtaID:        mtaID,
			ProcessID:    ongoingOperationId,
			ProcessType:  "DEPLOY",
			SpaceID:      spaceId,
			State:        "RUNNING",
			User:         user,
		}

		var ongoingOperations = []*models.Operation{&ongoingOperation}

		var getOutputLines = func(processID string, abortedProcessId string) []string {
			lines := []string{}
			lines = append(lines,
				"Undeploying multi-target app "+mtaID+" in org "+org+" / space "+space+" as "+user+"...\n")
			if abortedProcessId != "" {
				lines = append(lines,
					"Executing action 'abort' on operation "+abortedProcessId+"...\n",
					"OK\n")
			}
			lines = append(lines,
				"Operation id: "+processID+"\n",
				"Test message\n",
				"Process finished.\n",
				"Use \"cf dmol -i "+processID+"\" to download the logs of the process.\n")
			return lines
		}

		BeforeEach(func() {
			ui.DisableTerminalOutput(true)
			name = command.GetPluginCommand().Name
			cliConnection = cliFakes.NewFakeCliConnectionBuilder().
				CurrentOrg("test-org-guid", org, nil).
				CurrentSpace(spaceId, space, nil).
				Username(user, nil).
				AccessToken("bearer test-token", nil).Build()
			mtaClient = mtaFake.NewFakeMtaClientBuilder().
				GetMta(mtaID, nil, nil).
				GetMtaOperations(&[]string{mtaID}[0], nil, nil, nil, nil).
				StartMtaOperation(testutil.OperationResult, mtaclient.ResponseHeader{Location: "operations/1000?embed=messages"}, nil).
				GetMtaOperation(testutil.ProcessID, "messages", &undeployOperation, nil).Build()
			testClientFactory = commands.NewTestClientFactory(mtaClient, nil)
			command = commands.NewUndeployCommand()
			testTokenFactory := commands.NewTestTokenFactory(cliConnection)
			deployServiceURLCalculator := utilFakes.NewDeployServiceURLFakeCalculator("deploy-service.test.ondemand.com")
			command.InitializeAll(name, cliConnection, testutil.NewCustomTransport(200, nil), nil, testClientFactory, testTokenFactory, deployServiceURLCalculator, configuration.NewSnapshot())
		})

		// unknown flag - error
		Context("with an unknown flag", func() {
			It("should print incorrect usage, call cf help, and exit with a non-zero status", func() {
				output, status := oc.CaptureOutputAndStatus(func() int {
					return command.Execute([]string{"test-mta-id", "-unknown-flag"}).ToInt()
				})
				ex.ExpectFailure(status, output, "Incorrect usage. Unknown or wrong flag")
				Expect(cliConnection.CliCommandArgsForCall(0)).To(Equal([]string{"help", name}))
			})
		})

		// wrong arguments - error
		Context("with wrong arguments", func() {
			It("should print incorrect usage, call cf help, and exit with a non-zero status", func() {
				output, status := oc.CaptureOutputAndStatus(func() int {
					return command.Execute([]string{"test-mta-id", "y", "z"}).ToInt()
				})
				ex.ExpectFailure(status, output, "Incorrect usage. Wrong arguments")
				Expect(cliConnection.CliCommandArgsForCall(0)).To(Equal([]string{"help", name}))
			})
		})

		// no arguments - error
		Context("with no arguments", func() {
			It("should print incorrect usage, call cf help, and exit with a non-zero status", func() {
				output, status := oc.CaptureOutputAndStatus(func() int {
					return command.Execute([]string{}).ToInt()
				})
				ex.ExpectFailure(status, output, "Incorrect usage. Missing positional argument 'MTA_ID'")
				Expect(cliConnection.CliCommandArgsForCall(0)).To(Equal([]string{"help", name}))
			})
		})

		// no MTA argument - error
		Context("with no mta id argument", func() {
			It("should print incorrect usage, call cf help, and exit with a non-zero status", func() {
				output, status := oc.CaptureOutputAndStatus(func() int {
					return command.Execute([]string{"-f"}).ToInt()
				})
				ex.ExpectFailure(status, output, "Incorrect usage. Missing positional argument 'MTA_ID'")
				Expect(cliConnection.CliCommandArgsForCall(0)).To(Equal([]string{"help", name}))
			})
		})

		// non-existing MTA_ID - failure
		Context("with an incorrect mta id provided", func() {
			It("should display error and exit with non-zero status", func() {
				var clientError = baseclient.NewClientError(testutil.ClientError)
				testClientFactory.MtaClient = mtaFake.NewFakeMtaClientBuilder().
					GetMta("test-non-existing-id", nil, clientError).Build()
				output, status := oc.CaptureOutputAndStatus(func() int {
					return command.Execute([]string{"test-non-existing-id", "-f"}).ToInt()
				})
				ex.ExpectFailureOnLine(status, output, "Could not get multi-target app test-non-existing-id:", 1)
			})
		})

		// existing MTA_ID and ongoing operations and force option
		Context("with a correct mta id provided and ongoing operation found and force option provided", func() {
			It("should try to abort the conflicting process and fail it", func() {
				testClientFactory.MtaClient = mtaFake.NewFakeMtaClientBuilder().
					GetMta(mtaID, nil, nil).
					GetMtaOperations(&[]string{mtaID}[0], nil, nil, ongoingOperations, nil).
					GetOperationActions(ongoingOperationId, []string{"abort"}, nil).
					ExecuteAction(ongoingOperationId, "abort", mtaclient.ResponseHeader{Location: "operations/999?embed=messages"}, fmt.Errorf("test-error")).Build()
				output, status := oc.CaptureOutputAndStatus(func() int {
					return command.Execute([]string{mtaID, "-f"}).ToInt()
				})
				ex.ExpectFailureOnLine(status, output, "Could not execute action 'abort' on operation 999: test-error", 2)
			})

			It("should try to abort the conflicting process and success", func() {
				testClientFactory.MtaClient = mtaFake.NewFakeMtaClientBuilder().
					GetMta(mtaID, nil, nil).
					GetMtaOperations(&[]string{mtaID}[0], nil, nil, ongoingOperations, nil).
					GetOperationActions(ongoingOperationId, []string{"abort"}, nil).
					ExecuteAction(ongoingOperationId, "abort", mtaclient.ResponseHeader{Location: "operations/999?embed=messages"}, nil).
					StartMtaOperation(testutil.OperationResult, mtaclient.ResponseHeader{Location: "operations/1000?embed=messages"}, nil).
					GetMtaOperation(testutil.ProcessID, "messages", &undeployOperation, nil).Build()
				output, status := oc.CaptureOutputAndStatus(func() int {
					return command.Execute([]string{mtaID, "-f"}).ToInt()
				})
				ex.ExpectSuccessWithOutput(status, output, getOutputLines(testutil.ProcessID, ongoingOperationId))
			})
		})

		// existing MTA_ID and no ongoing operations - success
		Context("with a correct mta id provided and no ongoing operations", func() {
			It("should proceed without trying to abort conflicting process", func() {
				output, status := oc.CaptureOutputAndStatus(func() int {
					return command.Execute([]string{"test", "-f"}).ToInt()
				})
				ex.ExpectSuccessWithOutput(status, output, getOutputLines(testutil.ProcessID, ""))
			})

			It("should proceed without trying to abort conflicting process with more options", func() {
				output, status := oc.CaptureOutputAndStatus(func() int {
					return command.Execute([]string{"test", "-delete-services", "-no-restart-subscribed-apps", "-delete-service-brokers", "-do-not-fail-on-missing-permissions", "-f"}).ToInt()
				})
				ex.ExpectSuccessWithOutput(status, output, getOutputLines(testutil.ProcessID, ""))
			})
		})

		// unable to start operation - failure
		Context("with a correct mta id provided and failing start of operation", func() {
			It("should display error and exit with non-zero status", func() {
				testClientFactory.MtaClient = mtaFake.NewFakeMtaClientBuilder().
					GetMta(mtaID, nil, nil).
					GetMtaOperations(&[]string{mtaID}[0], nil, nil, nil, nil).
					StartMtaOperation(testutil.OperationResult, mtaclient.ResponseHeader{Location: "operations/1000?embed=messages"}, fmt.Errorf("test-error")).Build()

				output, status := oc.CaptureOutputAndStatus(func() int {
					return command.Execute([]string{testutil.ProcessID, "-f"}).ToInt()
				})
				ex.ExpectFailureOnLine(status, output, "Could not create undeploy process: test-error", 1)
			})
		})
	})
})
