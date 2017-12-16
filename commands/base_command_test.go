package commands_test

//
import (
	"fmt"
	"net/http"
	"os"

	"github.com/cloudfoundry/cli/cf/terminal"

	baseclient "github.com/SAP/cf-mta-plugin/clients/baseclient"
	"github.com/SAP/cf-mta-plugin/clients/models"
	"github.com/SAP/cf-mta-plugin/clients/mtaclient"
	mtafake "github.com/SAP/cf-mta-plugin/clients/mtaclient/fakes"
	fakes "github.com/SAP/cf-mta-plugin/clients/restclient/fakes"
	"github.com/SAP/cf-mta-plugin/commands"
	cmd_fakes "github.com/SAP/cf-mta-plugin/commands/fakes"
	"github.com/SAP/cf-mta-plugin/testutil"
	"github.com/SAP/cf-mta-plugin/ui"
	plugin_fakes "github.com/cloudfoundry/cli/plugin/fakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("BaseCommand", func() {
	const org = "test-org"
	const space = "test-space"
	const user = "test-user"
	const deployID = "DEPLOY"

	var fakeCliConnection *plugin_fakes.FakeCliConnection
	var command *commands.BaseCommand
	var oc = testutil.NewUIOutputCapturer()
	var ex = testutil.NewUIExpector()

	fakeMtaClientBuilder := mtafake.NewFakeMtaClientBuilder()
	testTokenFactory := commands.NewTestTokenFactory(fakeCliConnection)

	BeforeEach(func() {
		ui.DisableTerminalOutput(true)
		command = &commands.BaseCommand{}
		fakeCliConnection = cmd_fakes.NewFakeCliConnectionBuilder().
			CurrentOrg("test-org-guid", org, nil).
			CurrentSpace("test-space-guid", space, nil).
			Username(user, nil).
			AccessToken("bearer test-token", nil).
			APIEndpoint("https://api.test.ondemand.com", nil).Build()

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
				fakeCliConnection := cmd_fakes.NewFakeCliConnectionBuilder().
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
				fakeCliConnection := cmd_fakes.NewFakeCliConnectionBuilder().
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
				fakeCliConnection := cmd_fakes.NewFakeCliConnectionBuilder().
					Username("", nil).Build()
				command.Initialize("test", fakeCliConnection)
				_, err := command.GetUsername()
				Expect(err).To(MatchError(fmt.Errorf("Not logged in. Use '%s' to log in.", terminal.CommandColor("cf login"))))
			})
		})
	})

	Describe("GetDeployServiceHost", func() {
		Context("with an environment variable", func() {
			BeforeEach(func() {
				os.Setenv("DEPLOY_SERVICE_URL", "test")
			})
			It("should return the deploy service host set in the environment", func() {
				command.Initialize("test", fakeCliConnection)
				Expect(command.GetDeployServiceURL()).To(Equal("test"))
			})
			AfterEach(func() {
				os.Clearenv()
			})
		})
		Context("with valid API endpoint returned by the CLI connection", func() {
			It("should return the deploy service host constructed from the API endpoint", func() {
				command.Initialize("test", fakeCliConnection)
				Expect(command.GetDeployServiceURL()).To(Equal("deploy-service.test.ondemand.com"))
			})
		})
		Context("with no API endpoint returned by the CLI connection", func() {
			It("should print an error and exit with a non-zero status", func() {
				fakeCliConnection := cmd_fakes.NewFakeCliConnectionBuilder().
					APIEndpoint("", nil).Build()
				command.Initialize("test", fakeCliConnection)
				_, err := command.GetDeployServiceURL()
				Expect(err).To(MatchError(fmt.Errorf("No api endpoint set. Use '%s' to set an endpoint.", terminal.CommandColor("cf api"))))
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
				GetMtaOperations(nil, nil, []*models.Operation{ongoingOperationToReturn}, nil).Build()

			command.InitializeAll("test", fakeCliConnection, testutil.NewCustomTransport(http.StatusOK, nil), nil, testClientFactory, testTokenFactory)
		})
		Context("with valid ongoing operations", func() {
			It("should abort and exit with zero status", func() {
				output := oc.CaptureOutput(func() {
					wasAborted, err = command.CheckOngoingOperation(mtaID, "test-host", true)
				})
				Expect(err).ShouldNot(HaveOccurred())
				Expect(wasAborted).To(BeTrue())
				Expect(output).To(Equal([]string{"Aborting multi-target app operation with id test...\n", "OK\n"}))
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
				ExecuteAction("test-process-id", "abort", mtaclient.ResponseHeader{}, nil).
				ExecuteAction("test-process-id", "retry", mtaclient.ResponseHeader{Location: "operations/test-process-id?embed=messages"}, nil).Build()
			command.InitializeAll("test", fakeCliConnection, testutil.NewCustomTransport(200, nil), nil, testClientfactory, testTokenFactory)
		})
		Context("with valid process id and valid action id", func() {
			It("should abort and exit with zero status", func() {
				output, status := oc.CaptureOutputAndStatus(func() int {
					return command.ExecuteAction("test-process-id", "abort", "test-host").ToInt()
				})
				ex.ExpectSuccessWithOutput(status, output, []string{"Aborting multi-target app operation with id test-process-id...\n", "OK\n"})
			})
		})
		Context("with non-valid process id and valid action id", func() {
			It("should return error and exit with non-zero status", func() {
				output, status := oc.CaptureOutputAndStatus(func() int {
					return command.ExecuteAction("not-valid-process-id", "abort", "test-host").ToInt()
				})
				ex.ExpectFailureOnLine(status, output, "Multi-target app operation with id not-valid-process-id not found", 0)
			})
		})

		Context("with valid process id and invalid action id", func() {
			It("should return error and exit with non-zero status", func() {
				output, status := oc.CaptureOutputAndStatus(func() int {
					return command.ExecuteAction("test-process-id", "not-existing-action", "test-host").ToInt()
				})
				ex.ExpectFailureOnLine(status, output, "Invalid action not-existing-action", 0)
			})
		})

		Context("with valid process id and valid action id", func() {
			It("should retry the process and exit with zero status", func() {
				output, status := oc.CaptureOutputAndStatus(func() int {
					return command.ExecuteAction("test-process-id", "retry", "test-host").ToInt()
				})
				ex.ExpectSuccessWithOutput(status, output, []string{"Retrying multi-target app operation with id test-process-id...\n", "OK\n",
					"Monitoring process test-process-id...\n", "Process finished.\n", "Use \"cf dmol -i test-process-id\" to download the logs of the process.\n"})
			})
		})
	})
})

func newClientError(code int, status, description string) error {
	return &baseclient.ClientError{Code: code, Status: status, Description: description}
}
