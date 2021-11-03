package commands_test

import (
	"net/http"

	plugin_fakes "code.cloudfoundry.org/cli/plugin/pluginfakes"
	cli_fakes "github.com/cloudfoundry-incubator/multiapps-cli-plugin/cli/fakes"
	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/clients/baseclient"
	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/clients/models"
	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/clients/mtaclient"
	mtafake "github.com/cloudfoundry-incubator/multiapps-cli-plugin/clients/mtaclient/fakes"
	mtaV2fake "github.com/cloudfoundry-incubator/multiapps-cli-plugin/clients/mtaclient_v2/fakes"
	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/clients/restclient/fakes"
	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/commands"
	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/testutil"
	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/ui"
	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/util"
	util_fakes "github.com/cloudfoundry-incubator/multiapps-cli-plugin/util/fakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("BaseCommand", func() {
	const org = "test-org"
	const space = "test-space"
	const user = "test-user"

	var fakeCliConnection *plugin_fakes.FakeCliConnection
	var command *commands.BaseCommand
	var oc = testutil.NewUIOutputCapturer()
	var ex = testutil.NewUIExpector()
	var cfTarget util.CloudFoundryTarget

	fakeMtaClientBuilder := mtafake.NewFakeMtaClientBuilder()
	fakeMtaV2ClientBuilder := mtaV2fake.NewFakeMtaV2ClientBuilder()
	testTokenFactory := commands.NewTestTokenFactory(fakeCliConnection)

	BeforeEach(func() {
		ui.DisableTerminalOutput(true)
		command = &commands.BaseCommand{}
		fakeCliConnection = cli_fakes.NewFakeCliConnectionBuilder().
			CurrentOrg("test-org-guid", org, nil).
			CurrentSpace("test-space-guid", space, nil).
			Username(user, nil).
			AccessToken("bearer test-token", nil).Build()
	})

	Describe("CheckOngoingOperation", func() {
		var wasAborted bool
		var err error
		var mtaID string
		var namespace string
		var ongoingOperationToReturn *models.Operation

		var fakeRestClientBuilder *fakes.FakeRestClientBuilder
		var testClientFactory *commands.TestClientFactory

		BeforeEach(func() {
			mtaID = "mtaId"
			namespace = "namespace"
			ongoingOperationToReturn = testutil.GetOperation("test", "test-space-guid", mtaID, namespace, "deploy", "ERROR", true)

			fakeRestClientBuilder = fakes.NewFakeRestClientBuilder()
			testClientFactory = commands.NewTestClientFactory(fakeMtaClientBuilder.Build(), fakeMtaV2ClientBuilder.Build(), fakeRestClientBuilder.Build())

			testClientFactory.MtaClient = fakeMtaClientBuilder.
				GetMtaOperations(nil, nil, nil, []*models.Operation{ongoingOperationToReturn}, nil).
				GetOperationActions("test", []string{"abort", "retry"}, nil).
				Build()
			deployServiceURLCalculator := util_fakes.NewDeployServiceURLFakeCalculator("deploy-service.test.ondemand.com")

			command.InitializeAll("test", fakeCliConnection, testutil.NewCustomTransport(http.StatusOK), testClientFactory, testTokenFactory, deployServiceURLCalculator)
			cfTarget, _ = command.GetCFTarget()
		})
		Context("with valid ongoing operations", func() {
			It("should abort and exit with zero status", func() {
				output := oc.CaptureOutput(func() {
					wasAborted, err = command.CheckOngoingOperation(mtaID, namespace, "test-host", true, cfTarget)
				})
				Expect(err).ShouldNot(HaveOccurred())
				Expect(wasAborted).To(BeTrue())
				Expect(output).To(Equal([]string{"Executing action 'abort' on operation test...", "OK"}))
			})
		})
		Context("with one ongoing operation which does not have an MTA ID", func() {
			It("should exit with zero status", func() {
				nonConflictingOperation := testutil.GetOperation("111", "space-guid", "", "", "deploy", "ERROR", false)
				testClientFactory.MtaClient = fakeMtaClientBuilder.
					GetMtaOperations(nil, nil, nil, []*models.Operation{nonConflictingOperation}, nil).Build()
				wasAborted, err = command.CheckOngoingOperation(mtaID, namespace, "test-host", true, cfTarget)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(wasAborted).To(BeTrue())
			})
		})
		Context("with no ongoing operations", func() {
			It("should exit with zero status", func() {
				testClientFactory.MtaClient = fakeMtaClientBuilder.
					GetMtaOperations(nil, nil, nil, []*models.Operation{}, nil).Build()
				wasAborted, err = command.CheckOngoingOperation(mtaID, namespace, "test-host", true, cfTarget)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(wasAborted).To(BeTrue())
			})
		})
		Context("with valid ongoing operation but in a different namespace", func() {
			It("should exit with zero status", func() {
				wasAborted, err = command.CheckOngoingOperation(mtaID, "namespace2", "test-host", true, cfTarget)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(wasAborted).To(BeTrue())
			})
		})
		Context("with valid ongoing operation, namespace is empty", func() {
			It("should abort and exit with zero status", func() {
				conflictingOperationWithEmptyNamespace := testutil.GetOperation("test", "test-space-guid", mtaID, "", "deploy", "ERROR", true)
				testClientFactory.MtaClient = fakeMtaClientBuilder.
					GetMtaOperations(nil, nil, nil, []*models.Operation{conflictingOperationWithEmptyNamespace}, nil).Build()
				output := oc.CaptureOutput(func() {
					wasAborted, err = command.CheckOngoingOperation(mtaID, "", "test-host", true, cfTarget)
				})
				Expect(err).ShouldNot(HaveOccurred())
				Expect(wasAborted).To(BeTrue())
				Expect(output).To(Equal([]string{"Executing action 'abort' on operation test...", "OK"}))
			})
		})
		Context("with valid ongoing operations and no force option specified", func() {
			It("should exit with non-zero status", func() {
				output := oc.CaptureOutput(func() {
					wasAborted, err = command.CheckOngoingOperation(mtaID, namespace, "test-host", false, cfTarget)
				})
				Expect(err).ShouldNot(HaveOccurred())
				Expect(wasAborted).To(BeFalse())
				// The confirmation message is only printed to stdout, not to the output bucket,
				// that's why we only check for the cancellation message
				Expect(output).To(Equal([]string{"Deploy cancelled"}))
			})
		})
	})

	Describe("ExecuteAction", func() {
		var ongoingOperationToReturn *models.Operation
		fakeRestClientBuilder := fakes.NewFakeRestClientBuilder()
		testClientfactory := commands.NewTestClientFactory(fakeMtaClientBuilder.Build(), fakeMtaV2ClientBuilder.Build(), fakeRestClientBuilder.Build())
		BeforeEach(func() {
			ongoingOperationToReturn = testutil.GetOperation("test-process-id", "test-space-guid", "test", "namespace", "deploy", "ERROR", true)
			testClientfactory.MtaClient = fakeMtaClientBuilder.
				GetMtaOperations(nil, nil, nil, []*models.Operation{ongoingOperationToReturn}, nil).
				GetMtaOperation("test-process-id", "mesages", &testutil.SimpleOperationResult, nil).
				GetOperationActions("test-process-id", []string{"abort", "retry"}, nil).
				ExecuteAction("test-process-id", "abort", mtaclient.ResponseHeader{}, nil).
				ExecuteAction("test-process-id", "retry", mtaclient.ResponseHeader{Location: "operations/test-process-id?embed=messages"}, nil).Build()
			deployServiceURLCalculator := util_fakes.NewDeployServiceURLFakeCalculator("deploy-service.test.ondemand.com")
			command.InitializeAll("test", fakeCliConnection, testutil.NewCustomTransport(200), testClientfactory, testTokenFactory, deployServiceURLCalculator)
			cfTarget, _ = command.GetCFTarget()
		})
		Context("with valid process id and valid action id", func() {
			It("should abort and exit with zero status", func() {
				output, status := oc.CaptureOutputAndStatus(func() int {
					return command.ExecuteAction("test-process-id", "abort", 0, "test-host", cfTarget).ToInt()
				})
				ex.ExpectSuccessWithOutput(status, output, []string{"Executing action 'abort' on operation test-process-id...", "OK"})
			})
		})
		Context("with non-valid process id and valid action id", func() {
			It("should return error and exit with non-zero status", func() {
				output, status := oc.CaptureOutputAndStatus(func() int {
					return command.ExecuteAction("not-valid-process-id", "abort", 0, "test-host", cfTarget).ToInt()
				})
				ex.ExpectFailure(status, output, "Multi-target app operation with ID not-valid-process-id not found")
			})
		})

		Context("with valid process id and invalid action id", func() {
			It("should return error and exit with non-zero status", func() {
				output, status := oc.CaptureOutputAndStatus(func() int {
					return command.ExecuteAction("test-process-id", "not-existing-action", 0, "test-host", cfTarget).ToInt()
				})
				ex.ExpectFailure(status, output, "Invalid action not-existing-action")
			})
		})

		Context("with valid process id and valid action id", func() {
			It("should retry the process and exit with zero status", func() {
				output, status := oc.CaptureOutputAndStatus(func() int {
					return command.ExecuteAction("test-process-id", "retry", 0, "test-host", cfTarget).ToInt()
				})
				ex.ExpectSuccessWithOutput(status, output, []string{"Executing action 'retry' on operation test-process-id...", "OK",
					"Process finished.", "Use \"cf dmol -i test-process-id\" to download the logs of the process."})
			})
		})
	})
})

func newClientError(code int, status, description string) error {
	return &baseclient.ClientError{Code: code, Status: status, Description: description}
}
