package commands_test

//
// import (
// 	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/clients/models"
// 	restfake "github.com/cloudfoundry-incubator/multiapps-cli-plugin/clients/restclient/fakes"
// 	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/commands"
// 	cli_fakes "github.com/cloudfoundry-incubator/multiapps-cli-plugin/cli/fakes"
// 	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/testutil"
// 	plugin_fakes "github.com/cloudfoundry/cli/plugin/fakes"
// 	. "github.com/onsi/ginkgo"
// 	. "github.com/onsi/gomega"
// )
//
// var _ = Describe("UndeployCommand", func() {
// 	Describe("Execute", func() {
// 		const org = "test-org"
// 		const space = "test-space"
// 		const user = "test-user"
// 		const mtaID = "test"
//
// 		var name string
// 		var cliConnection *plugin_fakes.FakeCliConnection
// 		var restClient *restfake.FakeRestClientOperations
// 		var testClientFactory *commands.TestClientFactory
// 		var mta *models.Mta
// 		var command *commands.UndeployCommand
// 		var oc = testutil.NewUIOutputCapturer()
// 		var ex = testutil.NewUIExpector()
//
// 		var getOutputLines = func(processAborted bool, actionID, processID string) []string {
// 			lines := []string{}
// 			lines = append(lines,
// 				"Undeploying multi-target app "+mtaID+" in org "+org+" / space "+space+" as "+user+"...\n")
// 			if processAborted {
// 				lines = append(lines,
// 					"Aborting multi-target app operation with id test...\n",
// 					"OK\n")
// 			}
// 			lines = append(lines,
// 				"Monitoring process execution...\n",
// 				"Process finished.\n")
// 			return lines
// 		}
//
// 		var getProcessParameters = func(additional bool) map[string]string {
// 			params := map[string]string{
// 				"mtaId": "test",
// 			}
// 			if additional {
// 				params["noRestartSubscribedApps"] = "true"
// 				params["deleteServices"] = "true"
// 				params["deleteServiceBrokers"] = "true"
// 			}
// 			return params
// 		}
//
// 		var expectProcessParameters = func(expectedParameters map[string]string, processParameters models.ProcessParameters) {
// 			for _, processParameter := range processParameters.Parameters {
// 				if expectedParameters[*processParameter.ID] != "" {
// 					Expect(processParameter.Value).To(Equal(expectedParameters[*processParameter.ID]))
// 				}
// 			}
// 		}
//
// 		BeforeEach(func() {
// 			// ui.DisableTerminalOutput(true)
// 			name = command.GetPluginCommand().Name
// 			cliConnection = cli_fakes.NewFakeCliConnectionBuilder().
// 				CurrentOrg("test-org-guid", org, nil).
// 				CurrentSpace("test-space-guid", space, nil).
// 				Username(user, nil).
// 				AccessToken("bearer test-token", nil).
// 				APIEndpoint("https://api.test.ondemand.com", nil).Build()
// 			slmpClient = slmpfake.NewFakeSlmpClientBuilder().
// 				GetMetadata(&testutil.SlmpMetadataResult, nil).
// 				GetService("xs2-undeploy", testutil.GetService("xs2-undeploy", "Undeploy", []*models.Parameter{testutil.GetParameter("mtaId")}), nil).
// 				CreateServiceProcess("", nil, &testutil.ProcessResult, nil).Build()
// 			slppClient = slppfake.NewFakeSlppClientBuilder().
// 				GetMetadata(&testutil.SlppMetadataResult, nil).
// 				GetTasklistTask(&testutil.TaskResult, nil).Build()
// 			mtaModule := testutil.GetMtaModule("test-module", []string{}, []string{})
// 			mta = testutil.GetMta("test", "test-version", []*models.MtaModulesItems0{mtaModule}, []string{"test-mta-services"})
// 			restClient = restfake.NewFakeRestClientBuilder().
// 				GetMta("test", mta, nil).
// 				GetMta("test-non-existing-id", nil, newClientError(404, "404 Not Found", `MTA with id "test-non-existing-id" does not exist`)).
// 				GetOperations(nil, nil, testutil.GetOperations([]*models.Operation{testutil.GetOperation("test", "test-space-guid", "test", "undeploy", "SLP_TASK_STATE_ERROR", true)}), nil).Build()
// 			testClientFactory = commands.NewTestClientFactory(restClient)
// 			command = &commands.UndeployCommand{}
// 			testTokenFactory := commands.NewTestTokenFactory(cliConnection)
// 			command.InitializeAll(name, cliConnection, testutil.NewCustomTransport(200, nil), nil, testClientFactory, testTokenFactory)
// 		})
//
// 		// unknown flag - error
// 		Context("with an unknown flag", func() {
// 			It("should print incorrect usage, call cf help, and exit with a non-zero status", func() {
// 				output, status := oc.CaptureOutputAndStatus(func() int {
// 					return command.Execute([]string{"test-mta-id", "-unknown-flag"}).ToInt()
// 				})
// 				ex.ExpectFailure(status, output, "Incorrect usage. Unknown or wrong flag.")
// 				Expect(cliConnection.CliCommandArgsForCall(0)).To(Equal([]string{"help", name}))
// 			})
// 		})
//
// 		// wrong arguments - error
// 		Context("with wrong arguments", func() {
// 			It("should print incorrect usage, call cf help, and exit with a non-zero status", func() {
// 				output, status := oc.CaptureOutputAndStatus(func() int {
// 					return command.Execute([]string{"test-mta-id", "y", "z"}).ToInt()
// 				})
// 				ex.ExpectFailure(status, output, "Incorrect usage. Wrong arguments.")
// 				Expect(cliConnection.CliCommandArgsForCall(0)).To(Equal([]string{"help", name}))
// 			})
// 		})
//
// 		// no arguments - error
// 		Context("with no arguments", func() {
// 			It("should print incorrect usage, call cf help, and exit with a non-zero status", func() {
// 				output, status := oc.CaptureOutputAndStatus(func() int {
// 					return command.Execute([]string{}).ToInt()
// 				})
// 				ex.ExpectFailure(status, output, "Incorrect usage. Missing positional argument 'MTA_ID'.")
// 				Expect(cliConnection.CliCommandArgsForCall(0)).To(Equal([]string{"help", name}))
// 			})
// 		})
//
// 		// no MTA argument - error
// 		Context("with no mta id argument", func() {
// 			It("should print incorrect usage, call cf help, and exit with a non-zero status", func() {
// 				output, status := oc.CaptureOutputAndStatus(func() int {
// 					return command.Execute([]string{"-f"}).ToInt()
// 				})
// 				ex.ExpectFailure(status, output, "Incorrect usage. Missing positional argument 'MTA_ID'.")
// 				Expect(cliConnection.CliCommandArgsForCall(0)).To(Equal([]string{"help", name}))
// 			})
// 		})
//
// 		// non-existing MTA_ID - failure
// 		Context("with an incorrect mta id provided", func() {
// 			It("should display error and exit with non-zero status", func() {
// 				output, status := oc.CaptureOutputAndStatus(func() int {
// 					return command.Execute([]string{"test-non-existing-id", "-f"}).ToInt()
// 				})
// 				ex.ExpectFailureOnLine(status, output, "Multi-target app test-non-existing-id not found", 1)
// 			})
// 		})
//
// 		// existing MTA_ID and force option - failure
// 		Context("with an correct mta id provided and ongoing operation found and force option provided", func() {
// 			It("should try to abort the conflicting process", func() {
// 				output, status := oc.CaptureOutputAndStatus(func() int {
// 					return command.Execute([]string{"test", "-f"}).ToInt()
// 				})
// 				ex.ExpectSuccessWithOutput(status, output, getOutputLines(true, "abort", "test"))
// 				serviceID, process := slmpClient.CreateServiceProcessArgsForCall(0)
// 				Expect(serviceID).To(Equal("xs2-undeploy"))
// 				expectProcessParameters(getProcessParameters(false), process.Parameters)
// 			})
// 		})
//
// 		// existing MTA_ID and no ongoing operations - failure
// 		Context("with an correct mta id provided and no ongoing operations", func() {
// 			It("should proceed without trying to abort conflicting process", func() {
// 				testClientFactory.RestClient = restfake.NewFakeRestClientBuilder().
// 					GetOperations(nil, nil, testutil.GetOperations([]*models.Operation{}), nil).
// 					GetMta("test", mta, nil).Build()
// 				output, status := oc.CaptureOutputAndStatus(func() int {
// 					return command.Execute([]string{"test", "-f"}).ToInt()
// 				})
// 				ex.ExpectSuccessWithOutput(status, output, getOutputLines(false, "", ""))
// 				serviceID, process := slmpClient.CreateServiceProcessArgsForCall(0)
// 				Expect(serviceID).To(Equal("xs2-undeploy"))
// 				expectProcessParameters(getProcessParameters(false), process.Parameters)
// 			})
// 		})
//
// 		// existing MTA_ID and no ongoing operations - failure
// 		Context("with an correct mta id provided and no ongoing operations and more options provided", func() {
// 			It("should proceed without trying to abort conflicting process", func() {
// 				testClientFactory.RestClient = restfake.NewFakeRestClientBuilder().
// 					GetOperations(nil, nil, testutil.GetOperations([]*models.Operation{}), nil).
// 					GetMta("test", mta, nil).Build()
// 				output, status := oc.CaptureOutputAndStatus(func() int {
// 					return command.Execute([]string{"test", "-delete-services", "-no-restart-subscribed-apps", "-delete-service-brokers", "-do-not-fail-on-missing-permissions", "-f"}).ToInt()
// 				})
// 				ex.ExpectSuccessWithOutput(status, output, getOutputLines(false, "", ""))
// 				serviceID, process := slmpClient.CreateServiceProcessArgsForCall(0)
// 				Expect(serviceID).To(Equal("xs2-undeploy"))
// 				expectProcessParameters(getProcessParameters(true), process.Parameters)
// 			})
// 		})
// 	})
// })
