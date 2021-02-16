package commands_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"

	cli_fakes "github.com/cloudfoundry-incubator/multiapps-cli-plugin/cli/fakes"
	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/clients/baseclient"
	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/clients/models"
	mtafake "github.com/cloudfoundry-incubator/multiapps-cli-plugin/clients/mtaclient/fakes"
	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/commands"
	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/testutil"
	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/ui"
	util_fakes "github.com/cloudfoundry-incubator/multiapps-cli-plugin/util/fakes"
	plugin_fakes "github.com/cloudfoundry/cli/plugin/fakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("DownloadMtaOperationLogsCommand", func() {
	Describe("Execute", func() {
		const org = "test-org"
		const space = "test-space"
		const user = "test-user"
		const mtaId = "test-mta-id"
		const spaceId = "test-space-guid"

		var name string
		var cliConnection *plugin_fakes.FakeCliConnection
		var mtaClient mtafake.FakeMtaClientOperations
		// var restClient *restfake.FakeRestClientOperations
		var clientFactory *commands.TestClientFactory
		var command *commands.DownloadMtaOperationLogsCommand
		var oc = testutil.NewUIOutputCapturer()
		var ex = testutil.NewUIExpector()
		var operations = []*models.Operation{
			newOperation(testutil.ProcessID, spaceId, user, mtaId),
			newOperation("1001", spaceId, user, mtaId),
			newOperation("1002", spaceId, user, mtaId),
		}

		var getOutputLines = func(dir string, processId string) []string {
			wd, _ := os.Getwd()
			return []string{
				fmt.Sprintf("Downloading logs of multi-target app operation with ID %s in org %s / space %s as %s...\n",
					processId, org, space, user),
				"OK\n",
				fmt.Sprintf("Saving logs to %s"+string(os.PathSeparator)+"%s...\n", wd, dir),
				fmt.Sprintf("  %s\n", testutil.LogID),
				"OK\n",
			}
		}

		var expectDirWithLog = func(dir string) {
			Expect(exists(dir)).To(Equal(true))
			Expect(exists(dir + "/" + testutil.LogID)).To(Equal(true))
			Expect(contentOf(dir + "/" + testutil.LogID)).To(Equal(testutil.LogContent))
		}

		BeforeEach(func() {
			ui.DisableTerminalOutput(true)
			name = command.GetPluginCommand().Name
			cliConnection = cli_fakes.NewFakeCliConnectionBuilder().
				CurrentOrg("test-org-guid", org, nil).
				CurrentSpace(spaceId, space, nil).
				Username(user, nil).
				AccessToken("bearer test-token", nil).Build()
			mtaClient = mtafake.NewFakeMtaClientBuilder().
				GetMtaOperationLogs(testutil.ProcessID, []*models.Log{&testutil.SimpleMtaLog}, nil).
				GetMtaOperations(&[]string{mtaId}[0], &[]int64{1}[0], nil, operations[:1], nil).
				GetMtaOperationLogContent(testutil.ProcessID, testutil.LogID, testutil.LogContent, nil).Build()
			clientFactory = commands.NewTestClientFactory(mtaClient, nil, nil)
			command = commands.NewDmolCommand()
			testTokenFactory := commands.NewTestTokenFactory(cliConnection)
			deployServiceURLCalculator := util_fakes.NewDeployServiceURLFakeCalculator("deploy-service.test.ondemand.com")
			command.InitializeAll(name, cliConnection, testutil.NewCustomTransport(200), clientFactory, testTokenFactory, deployServiceURLCalculator)
		})

		AfterEach(func() {
			os.RemoveAll("mta-op-" + testutil.ProcessID)
		})

		// unknown flag - error
		Context("with an unknown flag", func() {
			It("should print incorrect usage, call cf help, and exit with a non-zero status", func() {
				output, status := oc.CaptureOutputAndStatus(func() int {
					return command.Execute([]string{"-a"}).ToInt()
				})
				ex.ExpectFailure(status, output, "Incorrect usage. Unknown or wrong flag")
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

		// no arguments - error
		Context("with no arguments", func() {
			It("should print incorrect usage, call cf help, and exit with a non-zero status", func() {
				output, status := oc.CaptureOutputAndStatus(func() int {
					return command.Execute([]string{}).ToInt()
				})
				ex.ExpectFailure(status, output, "Incorrect usage. Missing required options '[i]'")
				Expect(cliConnection.CliCommandArgsForCall(0)).To(Equal([]string{"help", name}))
			})
		})

		// non-existing process id - error
		Context("with a non-existing process id", func() {
			It("should print an error and exit with a non-zero status", func() {
				os.Remove("mta-op-test")
				var clientError = baseclient.NewClientError(testutil.ClientError)
				clientFactory.MtaClient = mtafake.NewFakeMtaClientBuilder().
					GetMtaOperationLogs("test", []*models.Log{}, clientError).Build()
				output, status := oc.CaptureOutputAndStatus(func() int {
					return command.Execute([]string{"-i", "test"}).ToInt()
				})
				ex.ExpectFailureOnLine(status, output, "Could not get process logs: Process with id 404 not found (status 404): Process with id 404 not found", 1)
				Expect(exists("mta-op-test")).To(Equal(false))
			})
		})

		// non-existing mta id - error
		Context("with a non-existing mta id", func() {
			It("should print an error and exit with a non-zero status", func() {
				os.Remove("mta-op-test")
				var mtaId = "test-mta-id"
				var clientError = baseclient.NewClientError(testutil.ClientError)
				clientFactory.MtaClient = mtafake.NewFakeMtaClientBuilder().GetMtaOperations(&mtaId, nil, []string{}, []*models.Operation{}, clientError).Build()
				output, status := oc.CaptureOutputAndStatus(func() int {
					return command.Execute([]string{"--mta", mtaId}).ToInt()
				})
				ex.ExpectFailureOnLine(status, output, "Process with id 404 not found (status 404): Process with id 404 not found", 0)
				Expect(exists("mta-op-test")).To(Equal(false))
			})
		})

		// existing process id, backend returns an error response (GetLogs) - error
		Context("with an existing process id and an error response returned by the backend", func() {
			It("should print an error and exit with a non-zero status", func() {
				clientFactory.MtaClient = mtafake.NewFakeMtaClientBuilder().
					GetMtaOperationLogs(testutil.ProcessID, []*models.Log{}, fmt.Errorf("unknown error (status 500)")).Build()
				output, status := oc.CaptureOutputAndStatus(func() int {
					return command.Execute([]string{"-i", testutil.ProcessID}).ToInt()
				})
				ex.ExpectFailureOnLine(status, output, "Could not get process logs:", 1)
				Expect(exists("mta-op-" + testutil.ProcessID)).To(Equal(false))
			})
		})

		// existing process id, backend returns an error response (GetLogContent) - error
		Context("with an existing process id and an error response returned by the backend", func() {
			It("should print an error and exit with a non-zero status", func() {
				fakeMtaClient := mtafake.NewFakeMtaClientBuilder().
					GetMtaOperationLogs(testutil.ProcessID, []*models.Log{&testutil.SimpleMtaLog}, nil).
					GetMtaOperationLogContent("", "", "", fmt.Errorf("unknown error (status 500)")).Build()
				clientFactory.MtaClient = fakeMtaClient
				output, status := oc.CaptureOutputAndStatus(func() int {
					return command.Execute([]string{"-i", testutil.ProcessID}).ToInt()
				})
				ex.ExpectFailureOnLine(status, output, fmt.Sprintf("Could not get content of log %s:", testutil.LogID), 1)
				Expect(exists("mta-op-" + testutil.ProcessID)).To(Equal(false))
			})
		})

		// existing process id - success
		Context("with an existing process id", func() {
			const dir = "mta-op-" + testutil.ProcessID
			It("should download the logs for the current process and exit with zero status", func() {
				output, status := oc.CaptureOutputAndStatus(func() int {
					return command.Execute([]string{"-i", testutil.ProcessID}).ToInt()
				})
				ex.ExpectSuccessWithOutput(status, output, getOutputLines(dir, testutil.ProcessID))
				expectDirWithLog(dir)
			})
			AfterEach(func() {
				os.RemoveAll(dir)
			})
		})

		// existing mta id - success
		Context("with an existing mta id", func() {
			const dir = "mta-op-" + testutil.ProcessID
			It("should download the last logs for the specified mta and exit with zero status", func() {
				output, status := oc.CaptureOutputAndStatus(func() int {
					return command.Execute([]string{"--mta", mtaId, "--last", "1"}).ToInt()
				})
				ex.ExpectSuccessWithOutput(status, output, getOutputLines(dir, testutil.ProcessID))
				expectDirWithLog(dir)
			})
			AfterEach(func() {
				os.RemoveAll(dir)
			})
		})

		// existing process id and existing directory - error
		Context("with an existing process id and an existing directory", func() {
			const customDir string = "test"
			BeforeEach(func() {
				os.MkdirAll(customDir+"/mta-op-"+testutil.ProcessID, 0755)
			})
			It("should print an error and exit with a non-zero status", func() {
				output, status := oc.CaptureOutputAndStatus(func() int {
					return command.Execute([]string{"-i", testutil.ProcessID, "-d", customDir}).ToInt()
				})
				ex.ExpectFailureOnLine(status, output, fmt.Sprintf("Could not create download directory %s%s:", filepath.Join(customDir, "mta-op-"), testutil.ProcessID), 2)
			})
			AfterEach(func() {
				os.RemoveAll(customDir)
			})
		})

		if runtime.GOOS != "windows" {
			// existing process id and a directory that can't be written to - error
			Context("with an existing process id and a directory that can't be written to", func() {
				const customDir string = "test"
				BeforeEach(func() {
					os.Mkdir(customDir, 0000)
				})
				It("should print an error and exit with a non-zero status", func() {
					output, status := oc.CaptureOutputAndStatus(func() int {
						return command.Execute([]string{"-i", testutil.ProcessID, "-d", customDir + "/subdir"}).ToInt()
					})
					ex.ExpectFailureOnLine(status, output, fmt.Sprintf("Could not create download directory test/subdir/mta-op-%s:", testutil.ProcessID), 2)
				})
				AfterEach(func() {
					os.Chmod(customDir, 0755)
					os.RemoveAll(customDir)
				})
			})
		}

		// existing process id and non-existing directory - success
		Context("with an existing process id and a non-existing directory", func() {
			const customDir string = "test-non-existing"
			It("should create the directory, download the logs for the current process and exit with zero status", func() {
				output, status := oc.CaptureOutputAndStatus(func() int {
					return command.Execute([]string{"-i", testutil.ProcessID, "-d", customDir}).ToInt()
				})
				ex.ExpectSuccessWithOutput(status, output, getOutputLines(filepath.Join(customDir, "/mta-op-")+testutil.ProcessID, testutil.ProcessID))
				expectDirWithLog(customDir + "/mta-op-" + testutil.ProcessID)
			})
			AfterEach(func() {
				os.RemoveAll(customDir)
			})
		})

		// existing mta id and non-existing directory - success
		Context("with an existing mta id and a non-existing directory", func() {
			It("should create multiple directories, download the logs for the each process and exit with zero status", func() {
				clientFactory.MtaClient = mtafake.NewFakeMtaClientBuilder().
					GetMtaOperationLogs("1002", []*models.Log{&testutil.SimpleMtaLog}, nil).
					GetMtaOperationLogContent("1002", testutil.LogID, testutil.LogContent, nil).
					GetMtaOperations(&[]string{mtaId}[0], &[]int64{2}[0], nil, operations[1:], nil).Build()
				output, status := oc.CaptureOutputAndStatus(func() int {
					return command.Execute([]string{"--mta", mtaId, "--last", "2"}).ToInt()
				})
				expectedOutput := append(getOutputLines("mta-op-1001", "1001"), getOutputLines("mta-op-1002", "1002")...)
				ex.ExpectSuccessWithOutput(status, output, expectedOutput)
				expectDirWithLog("mta-op-1001")
				expectDirWithLog("mta-op-1002")
			})
			AfterEach(func() {
				os.RemoveAll("mta-op-1001")
				os.RemoveAll("mta-op-1002")
			})
		})

		// existing mta id and custom directory - success
		Context("with an existing mta id and a custom directory", func() {
			It("should create multiple directories in the custom one, download the logs for the each process and exit with zero status", func() {
				clientFactory.MtaClient = mtafake.NewFakeMtaClientBuilder().
					GetMtaOperationLogs("1002", []*models.Log{&testutil.SimpleMtaLog}, nil).
					GetMtaOperationLogContent("1002", testutil.LogID, testutil.LogContent, nil).
					GetMtaOperations(&[]string{mtaId}[0], &[]int64{2}[0], nil, operations[1:], nil).Build()
				output, status := oc.CaptureOutputAndStatus(func() int {
					return command.Execute([]string{"--mta", mtaId, "--last", "2", "-d", "custom-dir"}).ToInt()
				})
				expectedOutput := append(getOutputLines(filepath.Join("custom-dir", "mta-op-1001"), "1001"), getOutputLines(filepath.Join("custom-dir", "mta-op-1002"), "1002")...)
				ex.ExpectSuccessWithOutput(status, output, expectedOutput)
				expectDirWithLog("custom-dir/mta-op-1001")
				expectDirWithLog("custom-dir/mta-op-1002")
			})
			AfterEach(func() {
				os.RemoveAll("custom-dir/mta-op-1001")
				os.RemoveAll("custom-dir/mta-op-1002")
			})
		})

		// both mta id and process id - error
		Context("with mta id and process id", func() {
			It("should print incorrect usage - incompatible flags and exit with non-zero status", func() {
				output, status := oc.CaptureOutputAndStatus(func() int {
					return command.Execute([]string{"--mta", mtaId, "-i", testutil.ProcessID}).ToInt()
				})
				ex.ExpectFailure(status, output, "Incorrect usage. Option -i and option --mta are incompatible")
				Expect(cliConnection.CliCommandArgsForCall(0)).To(Equal([]string{"help", name}))
			})
		})
	})
})

func newOperation(operationId string, spaceId string, user string, mtaId string) *models.Operation {
	return &models.Operation{
		AcquiredLock: true,
		SpaceID:      spaceId,
		User:         user,
		ProcessID:    operationId,
		MtaID:        mtaId,
	}
}

func contentOf(fileName string) string {
	content, _ := ioutil.ReadFile(fileName)
	return string(content)
}

func exists(dirName string) bool {
	_, err := os.Stat(dirName)
	if err == nil {
		return true
	}
	return false
}
