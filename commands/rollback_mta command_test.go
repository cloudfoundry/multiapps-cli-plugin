package commands_test

import (
	"fmt"

	pluginFakes "code.cloudfoundry.org/cli/v8/plugin/pluginfakes"
	cliFakes "github.com/cloudfoundry-incubator/multiapps-cli-plugin/cli/fakes"
	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/clients/baseclient"
	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/clients/models"
	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/clients/mtaclient"
	mtaFake "github.com/cloudfoundry-incubator/multiapps-cli-plugin/clients/mtaclient/fakes"
	mtaV2Fake "github.com/cloudfoundry-incubator/multiapps-cli-plugin/clients/mtaclient_v2/fakes"
	mtaV2fake "github.com/cloudfoundry-incubator/multiapps-cli-plugin/clients/mtaclient_v2/fakes"
	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/commands"
	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/testutil"
	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/ui"
	utilFakes "github.com/cloudfoundry-incubator/multiapps-cli-plugin/util/fakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("RollbackMtaCommand", func() {
	const org = "test-org"
	const space = "test-space"
	const spaceId = "test-space-guid"
	const user = "test-user"
	const mtaID = "test"
	const ongoingOperationId = "999"

	var name string
	var cliConnection *pluginFakes.FakeCliConnection
	var mtaClient *mtaFake.FakeMtaClientOperations
	var testClientFactory *commands.TestClientFactory
	var command *commands.RollbackMtaCommand
	var oc = testutil.NewUIOutputCapturer()
	var ex = testutil.NewUIExpector()

	var rollbackMtaOperation = models.Operation{
		State:       "FINISHED",
		ProcessID:   testutil.ProcessID,
		ProcessType: "ROLLBACK_MTA",
		Messages:    []*models.Message{&testutil.SimpleMessage},
	}

	var ongoingOperation = models.Operation{
		AcquiredLock: true,
		Messages:     []*models.Message{&testutil.SimpleMessage},
		MtaID:        mtaID,
		ProcessID:    ongoingOperationId,
		ProcessType:  "ROLLBACK_MTA",
		SpaceID:      spaceId,
		State:        "RUNNING",
		User:         user,
	}

	var ongoingOperations = []*models.Operation{&ongoingOperation}

	var getOutputLines = func(processID string, abortedProcessId string) []string {
		lines := []string{}
		lines = append(lines,
			"Rollback multi-target app "+mtaID+" in org "+org+" / space "+space+" as "+user+"...")
		if abortedProcessId != "" {
			lines = append(lines,
				"Executing action \"abort\" on operation "+abortedProcessId+"...",
				"OK")
		}
		lines = append(lines,
			"Test message",
			"Process finished.",
			"Use \"cf dmol -i "+processID+"\" to download the logs of the process.")
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
			GetMtaOperation(testutil.ProcessID, "messages", &rollbackMtaOperation, nil).Build()
		mtaV2Client := mtaV2fake.NewFakeMtaV2ClientBuilder().
			GetMtasForThisSpace(mtaID, nil, nil, nil).Build()
		testClientFactory = commands.NewTestClientFactory(mtaClient, mtaV2Client, nil)
		command = commands.NewRollbackMtaCommand()
		testTokenFactory := commands.NewTestTokenFactory(cliConnection)
		deployServiceURLCalculator := utilFakes.NewDeployServiceURLFakeCalculator("deploy-service.test.ondemand.com")
		command.InitializeAll(name, cliConnection, testutil.NewCustomTransport(200), testClientFactory, testTokenFactory, deployServiceURLCalculator)
	})

	Describe("GetPluginCommand", func() {
		It("returns the correct plugin command", func() {
			pluginCmd := command.GetPluginCommand()
			Expect(pluginCmd.Name).To(Equal("rollback-mta"))
			Expect(pluginCmd.HelpText).To(Equal("(EXPERIMENTAL) Rollback of a multi-target app works only if [--backup-previous-version] flag was used during blue-green deployment and backup applications exists in the space"))
		})
	})

	Describe("Execute", func() {
		Context("when operation ID and action ID are provided", func() {
			It("executes the action", func() {
				testClientFactory.MtaClient = mtaFake.NewFakeMtaClientBuilder().
					GetMtaOperations(&[]string{mtaID}[0], nil, nil, ongoingOperations, nil).
					GetOperationActions(ongoingOperationId, []string{"abort"}, nil).
					ExecuteAction(ongoingOperationId, "abort", mtaclient.ResponseHeader{Location: "operations/999?embed=messages"}, nil).
					GetMtaOperation(testutil.ProcessID, "messages", &rollbackMtaOperation, nil).Build()
				testClientFactory.MtaV2Client = mtaV2Fake.NewFakeMtaV2ClientBuilder().
					GetMtasForThisSpace(mtaID, nil, nil, nil).Build()
				output, status := oc.CaptureOutputAndStatus(func() int {
					return command.Execute([]string{"-i", "999", "-a", "abort"}).ToInt()
				})
				ex.ExpectSuccessWithOutput(status, output, []string{"Executing action \"abort\" on operation 999...", "OK"})
			})

			It("executes the action and fails with not found operation", func() {
				testClientFactory.MtaClient = mtaFake.NewFakeMtaClientBuilder().
					GetMtaOperations(&[]string{mtaID}[0], nil, nil, []*models.Operation{}, nil).
					GetOperationActions(ongoingOperationId, []string{"abort"}, fmt.Errorf("not-found")).Build()
				testClientFactory.MtaV2Client = mtaV2Fake.NewFakeMtaV2ClientBuilder().
					GetMtasForThisSpace(mtaID, nil, nil, nil).Build()
				output, status := oc.CaptureOutputAndStatus(func() int {
					return command.Execute([]string{"-i", "999", "-a", "abort"}).ToInt()
				})
				ex.ExpectFailure(status, output, "Multi-target app operation with ID 999 not found")
			})
		})

		Context("with a correct mta id provided and ongoing operation found and force option provided", func() {
			It("should try to abort the conflicting process and fail it", func() {
				testClientFactory.MtaClient = mtaFake.NewFakeMtaClientBuilder().
					GetMtaOperations(&[]string{mtaID}[0], nil, nil, ongoingOperations, nil).
					GetOperationActions(ongoingOperationId, []string{"abort"}, nil).
					ExecuteAction(ongoingOperationId, "abort", mtaclient.ResponseHeader{Location: "operations/999?embed=messages"}, fmt.Errorf("test-error")).Build()
				testClientFactory.MtaV2Client = mtaV2Fake.NewFakeMtaV2ClientBuilder().
					GetMtasForThisSpace(mtaID, nil, nil, nil).Build()
				output, status := oc.CaptureOutputAndStatus(func() int {
					return command.Execute([]string{mtaID, "-f"}).ToInt()
				})
				ex.ExpectFailureOnLine(status, output, "Could not execute action \"abort\" on operation 999: test-error", 3)
			})

			It("should try to abort the conflicting process and success", func() {
				testClientFactory.MtaClient = mtaFake.NewFakeMtaClientBuilder().
					GetMtaOperations(&[]string{mtaID}[0], nil, nil, ongoingOperations, nil).
					GetOperationActions(ongoingOperationId, []string{"abort"}, nil).
					ExecuteAction(ongoingOperationId, "abort", mtaclient.ResponseHeader{Location: "operations/999?embed=messages"}, nil).
					StartMtaOperation(testutil.OperationResult, mtaclient.ResponseHeader{Location: "operations/1000?embed=messages"}, nil).
					GetMtaOperation(testutil.ProcessID, "messages", &rollbackMtaOperation, nil).Build()
				testClientFactory.MtaV2Client = mtaV2Fake.NewFakeMtaV2ClientBuilder().
					GetMtasForThisSpace(mtaID, nil, nil, nil).Build()
				output, status := oc.CaptureOutputAndStatus(func() int {
					return command.Execute([]string{mtaID, "-f"}).ToInt()
				})
				ex.ExpectSuccessWithOutput(status, output, getOutputLines(testutil.ProcessID, ongoingOperationId))
			})
		})

		Context("with error during during getting mtas for space", func() {
			It("should display error and exit with non-zero status", func() {
				var clientError = baseclient.NewClientError(testutil.ClientError)
				testClientFactory.MtaV2Client = mtaV2Fake.NewFakeMtaV2ClientBuilder().
					GetMtasForThisSpace("test-non-existing-id", nil, nil, clientError).Build()
				output, status := oc.CaptureOutputAndStatus(func() int {
					return command.Execute([]string{"test-non-existing-id", "-f"}).ToInt()
				})
				ex.ExpectFailureOnLine(status, output, "Could not get multi-target app test-non-existing-id:", 2)
			})
		})

		Context("with an incorrect mta id provided", func() {
			It("should display error and exit with non-zero status", func() {
				mta_id := "non-existing-mta"
				custom_error := testutil.NewCustomError(404, "mtas", "MTA with name \""+mta_id+"\" does not exist")
				var clientError = baseclient.NewClientError(custom_error)
				testClientFactory.MtaV2Client = mtaV2Fake.NewFakeMtaV2ClientBuilder().
					GetMtasForThisSpace(mta_id, nil, nil, clientError).Build()
				output, status := oc.CaptureOutputAndStatus(func() int {
					return command.Execute([]string{mta_id, "-f"}).ToInt()
				})
				ex.ExpectFailureOnLine(status, output, "Multi-target app "+mta_id+" not found", 2)
			})
		})

		Context("with an incorrect mta id and namespace provided", func() {
			It("should display error and exit with non-zero status", func() {
				mta_id := "non-existing-mta"
				namespace := "with-a-namespace"
				custom_error := testutil.NewCustomError(404, "mtas", "MTA with name \""+mta_id+"\" and namespace \""+namespace+"\" does not exist")
				var clientError = baseclient.NewClientError(custom_error)
				testClientFactory.MtaV2Client = mtaV2Fake.NewFakeMtaV2ClientBuilder().
					GetMtasForThisSpace(mta_id, &namespace, nil, clientError).Build()
				output, status := oc.CaptureOutputAndStatus(func() int {
					return command.Execute([]string{mta_id, "-f", "--namespace", namespace}).ToInt()
				})
				ex.ExpectFailureOnLine(status, output, "Multi-target app "+mta_id+" with namespace "+namespace+" not found", 2)
			})
		})

		Context("with a correct mta id provided and no ongoing operations", func() {
			It("should proceed without trying to abort conflicting process", func() {
				output, status := oc.CaptureOutputAndStatus(func() int {
					return command.Execute([]string{"test", "-f"}).ToInt()
				})
				ex.ExpectSuccessWithOutput(status, output, getOutputLines(testutil.ProcessID, ""))
			})

			It("should proceed without trying to abort conflicting process with more options", func() {
				output, status := oc.CaptureOutputAndStatus(func() int {
					return command.Execute([]string{"test", "-process-user-provided-services", "-do-not-fail-on-missing-permissions", "-f"}).ToInt()
				})
				ex.ExpectSuccessWithOutput(status, output, getOutputLines(testutil.ProcessID, ""))
			})
		})

		Context("with a correct mta id provided but fail to start operation", func() {
			It("should exit with non-zero status", func() {
				testClientFactory.MtaClient = mtaFake.NewFakeMtaClientBuilder().
					GetMtaOperations(&[]string{mtaID}[0], nil, nil, []*models.Operation{}, nil).
					StartMtaOperation(testutil.OperationResult, mtaclient.ResponseHeader{Location: "operations/1000?embed=messages"}, fmt.Errorf("server-error")).Build()
				testClientFactory.MtaV2Client = mtaV2Fake.NewFakeMtaV2ClientBuilder().
					GetMtasForThisSpace(mtaID, nil, nil, nil).Build()
				output, status := oc.CaptureOutputAndStatus(func() int {
					return command.Execute([]string{mtaID, "-f"}).ToInt()
				})
				ex.ExpectFailureOnLine(status, output, "Could not create rollback mta process: server-error", 2)
			})
		})

	})
})
