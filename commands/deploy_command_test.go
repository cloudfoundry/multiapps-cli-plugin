package commands_test

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	cli_fakes "github.com/cloudfoundry-incubator/multiapps-cli-plugin/cli/fakes"
	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/clients/models"
	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/clients/mtaclient"
	mtafake "github.com/cloudfoundry-incubator/multiapps-cli-plugin/clients/mtaclient/fakes"

	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/commands"
	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/testutil"
	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/ui"
	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/util"
	util_fakes "github.com/cloudfoundry-incubator/multiapps-cli-plugin/util/fakes"
	plugin_fakes "github.com/cloudfoundry/cli/plugin/fakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("DeployCommand", func() {
	Describe("Execute", func() {
		const org = "test-org"
		const space = "test-space"
		const namespace = "namespace"
		const user = "test-user"
		const testFilesLocation = "../test_resources/commands/"
		const testArchive = "mtaArchive.mtar"
		const mtaArchivePath = testFilesLocation + testArchive
		const extDescriptorPath = testFilesLocation + "extDescriptor.mtaext"

		var name string
		var cliConnection *plugin_fakes.FakeCliConnection
		// var fakeSession csrffake.FakeSessionProvider
		var mtaClient mtafake.FakeMtaClientOperations
		// var restClient *restfake.FakeRestClientOperations
		var testClientFactory *commands.TestClientFactory
		var command *commands.DeployCommand
		var oc = testutil.NewUIOutputCapturer()
		var ex = testutil.NewUIExpector()

		var fullMtaArchivePath, _ = filepath.Abs(mtaArchivePath)
		var fullExtDescriptorPath, _ = filepath.Abs(extDescriptorPath)
		var correctMtaUrl = "http://host/path/anatz.mtar?query=true"
		var incorrectMtaUrl = "http://alabala.com"

		var getLinesForAbortingProcess = func() []string {
			return []string{
				"Executing action 'abort' on operation test-process-id...\n",
				"OK\n",
			}
		}

		var getOutputLines = func(extDescriptor, processAborted, fromUrl bool) []string {
			var lines []string
			mtaNameToPrint := mtaArchivePath
			if fromUrl {
				mtaNameToPrint = "from url"
			}
			lines = append(lines,
				"Deploying multi-target app archive "+mtaNameToPrint+" in org "+org+" / space "+space+" as "+user+"...\n\n")
			if processAborted {
				lines = append(lines,
					"Executing action 'abort' on operation test-process-id...\n",
					"OK\n",
				)
			}
			if fromUrl {
				lines = append(lines, "OK\n")
			} else {
				lines = append(lines,
					"Uploading 1 files...\n",
					"  "+fullMtaArchivePath+"\n",
					"OK\n")
			}
			if extDescriptor {
				lines = append(lines,
					"Uploading 1 files...\n",
					"  "+fullExtDescriptorPath+"\n",
					"OK\n")
			}
			lines = append(lines, "Operation ID: 1000\n")
			lines = append(lines,
				"Test message\n",
				"Process finished.\n",
				"Use \"cf dmol -i 1000\" to download the logs of the process.\n")
			return lines
		}

		// var getProcessParameters = func(additional bool) map[string]string {
		// 	params := map[string]string{
		// 		"appArchiveId":   "mtaArchive.mtar",
		// 		"failOnCrashed":  "false",
		// 	}
		// 	if additional {
		// 		params["deleteServices"] = "true"
		// 		params["keepFiles"] = "true"
		// 		params["noStart"] = "true"
		// 	}
		// 	return params
		// }

		var getFile = func(path string) (*os.File, *models.FileMetadata) {
			file, _ := os.Open(path)
			digest, _ := util.ComputeFileChecksum(path, "MD5")
			f := testutil.GetFile(*file, strings.ToUpper(digest), namespace)
			return file, f
		}

		// var expectProcessParameters = func(expectedParameters map[string]string, processParameters map[string]interface{}) {
		// 	for processParam, processParamValue := range processParameters {
		// 		if expectedParameters[processParam] != "" {
		// 			Expect(processParamValue).To(Equal(expectedParameters[processParamValue.(string)]))
		// 		}
		// 	}
		// }

		BeforeEach(func() {
			ui.DisableTerminalOutput(true)
			name = command.GetPluginCommand().Name
			cliConnection = cli_fakes.NewFakeCliConnectionBuilder().
				CurrentOrg("test-org-guid", org, nil).
				CurrentSpace("test-space-guid", space, nil).
				Username(user, nil).
				AccessToken("bearer test-token", nil).Build()
			mtaArchiveFile, mtaArchive := getFile(mtaArchivePath)
			defer mtaArchiveFile.Close()
			extDescriptorFile, extDescriptor := getFile(extDescriptorPath)
			defer extDescriptorFile.Close()
			mtaClient = mtafake.NewFakeMtaClientBuilder().
				GetMtaFiles([]*models.FileMetadata{&testutil.SimpleFile}, nil).
				UploadMtaFile(*mtaArchiveFile, mtaArchive, nil).
				UploadMtaFile(*extDescriptorFile, extDescriptor, nil).
				UploadMtaArchiveFromUrl(correctMtaUrl, mtaArchive, nil).
				UploadMtaArchiveFromUrl(incorrectMtaUrl, nil, fmt.Errorf("connection refused")).
				StartMtaOperation(testutil.OperationResult, mtaclient.ResponseHeader{Location: "operations/1000?embed=messages"}, nil).
				GetMtaOperation("1000", "messages", &testutil.OperationResult, nil).
				GetMtaOperationLogContent("1000", testutil.LogID, testutil.LogContent, nil).
				GetMtaOperations(nil, nil, nil, []*models.Operation{&testutil.OperationResult}, nil).Build()
			testClientFactory = commands.NewTestClientFactory(mtaClient, nil, nil)
			command = commands.NewDeployCommand()
			testTokenFactory := commands.NewTestTokenFactory(cliConnection)
			deployServiceURLCalculator := util_fakes.NewDeployServiceURLFakeCalculator("deploy-service.test.ondemand.com")
			command.InitializeAll(name, cliConnection, testutil.NewCustomTransport(200), testClientFactory, testTokenFactory, deployServiceURLCalculator)
		})

		// unknown flag - error
		Context("with argument that is not a directory or MTA", func() {
			It("should print incorrect usage, call cf help, and exit with a non-zero status", func() {
				output, status := oc.CaptureOutputAndStatus(func() int {
					return command.Execute([]string{"x", "-l"}).ToInt()
				})
				ex.ExpectFailure(status, output, "Incorrect usage. Unknown or wrong flag")
				Expect(cliConnection.CliCommandArgsForCall(0)).To(Equal([]string{"help", name}))
			})
		})

		Context("with argument that is a directory or MTA and with unknown flag", func() {
			It("should print incorrect usage, call cf help, and exit with a non-zero status", func() {
				output, status := oc.CaptureOutputAndStatus(func() int {
					return command.Execute([]string{mtaArchivePath, "-l"}).ToInt()
				})
				ex.ExpectFailure(status, output, "Incorrect usage. Unknown or wrong flag")
				Expect(cliConnection.CliCommandArgsForCall(0)).To(Equal([]string{"help", name}))
			})
		})

		Context("with a correct URL argument", func() {
			It("should upload the MTAR from the correct URL and initiate a deploy", func() {
				output, status := oc.CaptureOutputAndStatus(func() int {
					return command.Execute([]string{correctMtaUrl}).ToInt()
				})
				ex.ExpectSuccessWithOutput(status, output, getOutputLines(false, false, true))
			})
		})

		Context("with an incorrect URL argument", func() {
			It("should fail with the error returned from the server", func() {
				output, status := oc.CaptureOutputAndStatus(func() int {
					return command.Execute([]string{incorrectMtaUrl}).ToInt()
				})
				ex.ExpectFailureOnLine(status, output, "Could not upload from url: connection refused", 1)
			})
		})

		// wrong arguments - error
		Context("with wrong arguments", func() {
			It("should print incorrect usage, call cf help, and exit with a non-zero status", func() {
				output, status := oc.CaptureOutputAndStatus(func() int {
					return command.Execute([]string{"x", "y", "z"}).ToInt()
				})
				ex.ExpectFailure(status, output, "Incorrect usage. Wrong arguments")
				Expect(cliConnection.CliCommandArgsForCall(0)).To(Equal([]string{"help", name}))
			})
		})

		// non-existing MTA archive - error
		Context("with a non-existing mta archive", func() {
			It("should print a file not found error and exit with a non-zero status", func() {
				const fileName = "non-existing-mtar.mtar"
				output, status := oc.CaptureOutputAndStatus(func() int {
					return command.Execute([]string{fileName}).ToInt()
				})
				ex.ExpectFailureOnLine(status, output, "Could not find MTA "+fileName, 0)
			})
		})

		// strategy flag set to "" - error
		Context("with strategy flag set to blank string", func() {
			It("should print incorrect usage, call cf help, and exit with a non-zero status", func() {
				output, status := oc.CaptureOutputAndStatus(func() int {
					return command.Execute([]string{mtaArchivePath, "--strategy"}).ToInt()
				})
				ex.ExpectFailure(status, output, "Incorrect usage. Unknown or wrong flag")
				Expect(cliConnection.CliCommandArgsForCall(0)).To(Equal([]string{"help", name}))
			})
		})

		// strategy flag set to invalid deployment strategy - error
		Context("with strategy flag set to an invalid deployment strategy", func() {
			It("should print the available strategies and exit with a non-zero status", func() {
				invalidStrategy := "asd"
				output, status := oc.CaptureOutputAndStatus(func() int {
					return command.Execute([]string{mtaArchivePath, "--strategy", invalidStrategy}).ToInt()
				})
				message := fmt.Sprintf("%s is not a valid deployment strategy, available strategies: %v", invalidStrategy, commands.AvailableStrategies())
				ex.ExpectFailureOnLine(status, output, message, 0)
			})
		})

		// TODO: can't connect to backend - error

		// TODO: backend returns an an error response - error

		// existing MTA archive - success
		Context("with an existing mta archive", func() {
			It("should upload 1 file and start the deployment process", func() {
				output, status := oc.CaptureOutputAndStatus(func() int {
					return command.Execute([]string{mtaArchivePath}).ToInt()
				})
				ex.ExpectSuccessWithOutput(status, output, getOutputLines(false, false, false))
				// operation := mtaClient.StartMtaOperationArgsForCall(1)
				// expectProcessParameters(getProcessParameters(false), operation.Parameters)
			})
		})

		// existing MTA archive and an extension descriptor - success
		Context("with an existing mta archive and an extension descriptor", func() {
			It("should upload 2 files and start the deployment process", func() {
				output, status := oc.CaptureOutputAndStatus(func() int {
					return command.Execute([]string{mtaArchivePath, "-e", extDescriptorPath}).ToInt()
				})
				ex.ExpectSuccessWithOutput(status, output, getOutputLines(true, false, false))
				// operation := mtaClient.StartMtaOperationArgsForCall(1)
				// expectProcessParameters(getProcessParameters(false), operation.Parameters)
			})
		})

		// existing MTA archive and additional options - success
		Context("with an existing mta archive and some options", func() {
			It("should upload 1 file and start the deployment process", func() {
				output, status := oc.CaptureOutputAndStatus(func() int {
					return command.Execute([]string{mtaArchivePath, "-f", "-delete-services", "-no-start", "-keep-files", "-do-not-fail-on-missing-permissions"}).ToInt()
				})
				ex.ExpectSuccessWithOutput(status, output, getOutputLines(false, false, false))
				// operation := mtaClient.StartMtaOperationArgsForCall(1)
				// expectProcessParameters(getProcessParameters(true), operation.Parameters)
			})
		})

		// non-existing ongoing operations - success
		// Context("with correct mta id from archive and no ongoing operations", func() {
		// 	It("should not try to abort confliction operations", func() {
		// 		testClientFactory.MtaClient = mtafake.NewFakeMtaClientBuilder().
		// 			GetMtaOperations(nil, nil, []*models.Operation{}, nil).
		// 			StartMtaOperation(models.Operation{}, mtaclient.ResponseHeader{Location: "operations/1000?embed=messages"}, nil).Build()
		// 		output, status := oc.CaptureOutputAndStatus(func() int {
		// 			return command.Execute([]string{mtaArchivePath}).ToInt()
		// 		})
		// 		fmt.Println(output)
		// 		ex.ExpectSuccessWithOutput(status, output, getOutputLines(false, false))
		// 		// operation := mtaClient.StartMtaOperationArgsForCall(1)
		// 		// expectProcessParameters(getProcessParameters(false), operation.Parameters)
		// 	})
		// })

		Context("with namespace starting with whitespace and ongoing operations", func() {
			It("should correctly detect a conflict with ongoing operation and result in error", func() {
				conflictingOperation := testutil.GetOperation("test", "test-space-guid", "test", "test", "deploy", "ERROR", true)
				testClientFactory.MtaClient = mtafake.NewFakeMtaClientBuilder().
					GetMtaOperations(nil, nil, nil, []*models.Operation{conflictingOperation}, nil).Build()
				output, status := oc.CaptureOutputAndStatus(func() int {
					return command.Execute([]string{mtaArchivePath, "--namespace", "    test   "}).ToInt()
				})
				ex.ExpectFailureOnLine(status, output, "There is an ongoing operation for multi-target app test", 1)
			})
		})

		// existing ongoing operations and force option not supplied - success
		Context("with correct mta id from archive, with ongoing operations provided and no force option", func() {
			It("should not try to abort confliction operations", func() {
				output, status := oc.CaptureOutputAndStatus(func() int {
					return command.Execute([]string{mtaArchivePath}).ToInt()
				})
				ex.ExpectSuccessWithOutput(status, output, getOutputLines(false, false, false))
				// operation := mtaClient.StartMtaOperationArgsForCall(1)
				// expectProcessParameters(getProcessParameters(false), operation.Parameters)
			})
		})

		// existing ongoing operations and force option supplied - success
		// Context("with correct mta id from archive, with ongoing operations provided and with force option", func() {
		// 	It("should try to abort confliction operations", func() {
		// 		testClientFactory.MtaClient = mtafake.NewFakeMtaClientBuilder().
		// 			GetMtaOperations(nil, nil, []*models.Operation{testutil.GetOperation("process-id", "test-space-guid", "test", "deploy", "ERROR", true)}, nil).Build()
		// 		output, status := oc.CaptureOutputAndStatus(func() int {
		// 			return command.Execute([]string{mtaArchivePath, "-f"}).ToInt()
		// 		})
		// 		ex.ExpectSuccessWithOutput(status, output, getOutputLines(false, true))
		// 		// operation := mtaClient.StartMtaOperationArgsForCall(1)
		// 		// expectProcessParameters(getProcessParameters(false), operation.Parameters)
		// 	})
		// })
		Context("with an error returned from getting ongoing operations", func() {
			It("should display error and exit witn non-zero status", func() {
				testClientFactory.MtaClient = mtafake.NewFakeMtaClientBuilder().
					GetMtaOperations(nil, nil, nil, []*models.Operation{}, fmt.Errorf("test-error-from backend")).Build()
				output, status := oc.CaptureOutputAndStatus(func() int {
					return command.Execute([]string{mtaArchivePath}).ToInt()
				})
				ex.ExpectFailureOnLine(status, output, "Could not get ongoing operation", 1)
			})
		})

		Context("with non-valid operation id and action id provided", func() {
			It("should return error and exit with non-zero status", func() {
				output, status := oc.CaptureOutputAndStatus(func() int {
					return command.Execute([]string{"-i", "test", "-a", "abort"}).ToInt()
				})
				ex.ExpectFailure(status, output, "Multi-target app operation with ID test not found")
			})
		})
		Context("with valid operation id and non-valid action id provided", func() {
			It("should return error and exit with non-zero status", func() {
				testClientFactory.MtaClient = mtafake.NewFakeMtaClientBuilder().
					GetMtaOperations(nil, nil, nil, []*models.Operation{
						testutil.GetOperation("test-process-id", space, "test-mta-id", namespace, "deploy", "ERROR", true),
					}, nil).Build()
				output, status := oc.CaptureOutputAndStatus(func() int {
					return command.Execute([]string{"-i", "test-process-id", "-a", "test"}).ToInt()
				})
				ex.ExpectFailure(status, output, "Invalid action test")
			})
		})
		Context("with valid operation id and no action id provided", func() {
			It("should return error and exit with non-zero status", func() {
				output, status := oc.CaptureOutputAndStatus(func() int {
					return command.Execute([]string{"-i", "test-process-id"}).ToInt()
				})
				ex.ExpectFailure(status, output, "All the a i options should be specified together")
			})
		})

		Context("with valid action id and no operation id provided", func() {
			It("should return error and exit with non-zero status", func() {
				output, status := oc.CaptureOutputAndStatus(func() int {
					return command.Execute([]string{"-a", "abort"}).ToInt()
				})
				ex.ExpectFailure(status, output, "All the a i options should be specified together")
			})
		})

		Context("with valid operation id and valid action id provided", func() {
			It("should execute action on the process specified with process id and exit with zero status", func() {
				testClientFactory.MtaClient = mtafake.NewFakeMtaClientBuilder().
					GetMtaOperations(nil, nil, nil, []*models.Operation{
						testutil.GetOperation("test-process-id", space, "test-mta-id", namespace, "deploy", "ERROR", true),
					}, nil).
					GetOperationActions("test", []string{"abort", "retry"}, nil).
					ExecuteAction("test-process-id", "test", mtaclient.ResponseHeader{Location: ""}, nil).Build()
				output, status := oc.CaptureOutputAndStatus(func() int {
					return command.Execute([]string{"-i", "test-process-id", "-a", "abort"}).ToInt()
				})
				ex.ExpectSuccessWithOutput(status, output, getLinesForAbortingProcess())
			})
		})
	})
})
