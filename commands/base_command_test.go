package commands_test

//
import (
	"fmt"
	"net/http"

	"github.com/cloudfoundry/cli/cf/terminal"

	cli_fakes "github.com/cloudfoundry-incubator/multiapps-cli-plugin/cli/fakes"
	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/clients/baseclient"
	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/clients/models"
	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/clients/mtaclient"
	mtafake "github.com/cloudfoundry-incubator/multiapps-cli-plugin/clients/mtaclient/fakes"
	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/clients/restclient/fakes"
	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/commands"
	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/configuration"
	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/testutil"
	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/ui"
	util_fakes "github.com/cloudfoundry-incubator/multiapps-cli-plugin/util/fakes"
	plugin_fakes "github.com/cloudfoundry/cli/plugin/fakes"
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

	fakeMtaClientBuilder := mtafake.NewFakeMtaClientBuilder()
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

	Describe("GetOrg", func() {
		Context("with valid org returned by the CLI connection", func() {
			It("should not exit or report any errors", func() {
				command.Initialize("test", fakeCliConnection)
				o, err := command.GetOrg()
				Expect(err).ShouldNot(HaveOccurred())
				Expect(o.Name).To(Equal(org))
			})
		})
		Context("with no org returned by the CLI connection", func() {
			It("should print an error and exit with a non-zero status", func() {
				fakeCliConnection := cli_fakes.NewFakeCliConnectionBuilder().
					CurrentOrg("", "", nil).Build()
				command.Initialize("test", fakeCliConnection)
				_, err := command.GetOrg()
				Expect(err).To(MatchError(fmt.Errorf("No org and space targeted, use '%s' to target an org and a space", terminal.CommandColor("cf target -o ORG -s SPACE"))))
			})
		})
	})

	Describe("GetSpace", func() {
		Context("with valid space returned by the CLI connection", func() {
			It("should not exit or report any errors", func() {
				command.Initialize("test", fakeCliConnection)
				s, err := command.GetSpace()
				Expect(err).ShouldNot(HaveOccurred())
				Expect(s.Name).To(Equal(space))
			})
		})
		Context("with no space returned by the CLI connection", func() {
			It("should print an error and exit with a non-zero status", func() {
				fakeCliConnection := cli_fakes.NewFakeCliConnectionBuilder().
					CurrentSpace("", "", nil).Build()
				command.Initialize("test", fakeCliConnection)
				_, err := command.GetSpace()
				Expect(err).To(MatchError(fmt.Errorf("No space targeted, use '%s' to target a space", terminal.CommandColor("cf target -s"))))
			})
		})
	})

	Describe("GetUsername", func() {
		Context("with valid username returned by the CLI connection", func() {
			It("should not exit or report any errors", func() {
				command.Initialize("test", fakeCliConnection)
				Expect(command.GetUsername()).To(Equal(user))
			})
		})
		Context("with no space returned by the CLI connection", func() {
			It("should print an error and exit with a non-zero status", func() {
				fakeCliConnection := cli_fakes.NewFakeCliConnectionBuilder().
					Username("", nil).Build()
				command.Initialize("test", fakeCliConnection)
				_, err := command.GetUsername()
				Expect(err).To(MatchError(fmt.Errorf("Not logged in. Use '%s' to log in.", terminal.CommandColor("cf login"))))
			})
		})
	})

	Describe("CheckOngoingOperation", func() {
		var wasAborted bool
		var err error
		var mtaID string
		var ongoingOperationToReturn *models.Operation

		var fakeRestClientBuilder *fakes.FakeRestClientBuilder
		var testClientFactory *commands.TestClientFactory

		BeforeEach(func() {
			mtaID = "mtaId"
			ongoingOperationToReturn = testutil.GetOperation("test", "test-space-guid", mtaID, "deploy", "ERROR", true)

			fakeRestClientBuilder = fakes.NewFakeRestClientBuilder()
			testClientFactory = commands.NewTestClientFactory(fakeMtaClientBuilder.Build(), fakeRestClientBuilder.Build())

			testClientFactory.MtaClient = fakeMtaClientBuilder.
				GetMtaOperations(nil, nil, []*models.Operation{ongoingOperationToReturn}, nil).
				GetOperationActions("test", []string{"abort", "retry"}, nil).
				Build()
			deployServiceURLCalculator := util_fakes.NewDeployServiceURLFakeCalculator("deploy-service.test.ondemand.com")

			command.InitializeAll("test", fakeCliConnection, testutil.NewCustomTransport(http.StatusOK, nil), nil, testClientFactory, testTokenFactory, deployServiceURLCalculator, configuration.NewSnapshot())
		})
		Context("with valid ongoing operations", func() {
			It("should abort and exit with zero status", func() {
				output := oc.CaptureOutput(func() {
					wasAborted, err = command.CheckOngoingOperation(mtaID, "test-host", true)
				})
				Expect(err).ShouldNot(HaveOccurred())
				Expect(wasAborted).To(BeTrue())
				Expect(output).To(Equal([]string{"Executing action 'abort' on operation test...\n", "OK\n"}))
			})
		})
		Context("with one ongoing operation which does not have an MTA ID", func() {
			It("should exit with zero status", func() {
				nonConflictingOperation := testutil.GetOperation("111", "space-guid", "", "deploy", "ERROR", false)
				testClientFactory.MtaClient = fakeMtaClientBuilder.
					GetMtaOperations(nil, nil, []*models.Operation{nonConflictingOperation}, nil).Build()
				wasAborted, err = command.CheckOngoingOperation(mtaID, "test-host", true)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(wasAborted).To(BeTrue())
			})
		})
		Context("with no ongoing operations", func() {
			It("should exit with zero status", func() {
				testClientFactory.MtaClient = fakeMtaClientBuilder.
					GetMtaOperations(nil, nil, []*models.Operation{}, nil).Build()
				wasAborted, err = command.CheckOngoingOperation(mtaID, "test-host", true)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(wasAborted).To(BeTrue())
			})
		})
		Context("with valid ongoing operations and no force option specified", func() {
			It("should exit with non-zero status", func() {
				output := oc.CaptureOutput(func() {
					wasAborted, err = command.CheckOngoingOperation(mtaID, "test-host", false)
				})
				Expect(err).ShouldNot(HaveOccurred())
				Expect(wasAborted).To(BeFalse())
				Expect(output).To(Equal([]string{"",
					"There is an ongoing operation for multi-target app mtaId. Do you want to abort it? (y/n)> ",
					"Deploy cancelled\n"}))
			})
		})
	})

	Describe("ExecuteAction", func() {
		var ongoingOperationToReturn *models.Operation
		fakeRestClientBuilder := fakes.NewFakeRestClientBuilder()
		testClientfactory := commands.NewTestClientFactory(fakeMtaClientBuilder.Build(), fakeRestClientBuilder.Build())
		BeforeEach(func() {
			ongoingOperationToReturn = testutil.GetOperation("test-process-id", "test-space-guid", "test", "deploy", "ERROR", true)
			testClientfactory.MtaClient = fakeMtaClientBuilder.
				GetMtaOperations(nil, nil, []*models.Operation{ongoingOperationToReturn}, nil).
				GetMtaOperation("test-process-id", "mesages", &testutil.SimpleOperationResult, nil).
				GetOperationActions("test-process-id", []string{"abort", "retry"}, nil).
				ExecuteAction("test-process-id", "abort", mtaclient.ResponseHeader{}, nil).
				ExecuteAction("test-process-id", "retry", mtaclient.ResponseHeader{Location: "operations/test-process-id?embed=messages"}, nil).Build()
			deployServiceURLCalculator := util_fakes.NewDeployServiceURLFakeCalculator("deploy-service.test.ondemand.com")
			command.InitializeAll("test", fakeCliConnection, testutil.NewCustomTransport(200, nil), nil, testClientfactory, testTokenFactory, deployServiceURLCalculator, configuration.NewSnapshot())
		})
		Context("with valid process id and valid action id", func() {
			It("should abort and exit with zero status", func() {
				output, status := oc.CaptureOutputAndStatus(func() int {
					return command.ExecuteAction("test-process-id", "abort", 0, "test-host").ToInt()
				})
				ex.ExpectSuccessWithOutput(status, output, []string{"Executing action 'abort' on operation test-process-id...\n", "OK\n"})
			})
		})
		Context("with non-valid process id and valid action id", func() {
			It("should return error and exit with non-zero status", func() {
				output, status := oc.CaptureOutputAndStatus(func() int {
					return command.ExecuteAction("not-valid-process-id", "abort", 0, "test-host").ToInt()
				})
				ex.ExpectFailure(status, output, "Multi-target app operation with id not-valid-process-id not found")
			})
		})

		Context("with valid process id and invalid action id", func() {
			It("should return error and exit with non-zero status", func() {
				output, status := oc.CaptureOutputAndStatus(func() int {
					return command.ExecuteAction("test-process-id", "not-existing-action", 0, "test-host").ToInt()
				})
				ex.ExpectFailure(status, output, "Invalid action not-existing-action")
			})
		})

		Context("with valid process id and valid action id", func() {
			It("should retry the process and exit with zero status", func() {
				output, status := oc.CaptureOutputAndStatus(func() int {
					return command.ExecuteAction("test-process-id", "retry", 0, "test-host").ToInt()
				})
				ex.ExpectSuccessWithOutput(status, output, []string{"Executing action 'retry' on operation test-process-id...\n", "OK\n",
					"Process finished.\n", "Use \"cf dmol -i test-process-id\" to download the logs of the process.\n"})
			})
		})
	})
})

func newClientError(code int, status, description string) error {
	return &baseclient.ClientError{Code: code, Status: status, Description: description}
}
