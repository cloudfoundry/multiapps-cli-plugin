package commands_test

import (
	"fmt"

	cli_fakes "github.com/cloudfoundry-incubator/multiapps-cli-plugin/cli/fakes"
	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/clients/models"
	mtaV2fake "github.com/cloudfoundry-incubator/multiapps-cli-plugin/clients/mtaclient_v2/fakes"
	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/commands"
	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/configuration"
	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/testutil"
	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/ui"
	util_fakes "github.com/cloudfoundry-incubator/multiapps-cli-plugin/util/fakes"
	plugin_fakes "github.com/cloudfoundry/cli/plugin/fakes"
	plugin_models "github.com/cloudfoundry/cli/plugin/models"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("MtaCommand", func() {
	Describe("Execute", func() {
		const org = "test-org"
		const space = "test-space"
		const user = "test-user"

		var namespace = "namespace"
		var name string
		var cliConnection *plugin_fakes.FakeCliConnection
		var clientFactory *commands.TestClientFactory
		var command *commands.MtaCommand
		var oc = testutil.NewUIOutputCapturer()
		var ex = testutil.NewUIExpector()

		var getOutputLines = func(mtaID, version, namespace string, apps, services [][]string) []string {
			lines := []string{}
			lines = append(lines,
				fmt.Sprintf("Showing health and status for multi-target app %s in org %s / space %s as %s...\n", mtaID, org, space, user))
			lines = append(lines, "OK\n")
			lines = append(lines, fmt.Sprintf("Version: %s\n", version))
			lines = append(lines, fmt.Sprintf("Namespace: %s\n", namespace))
			lines = append(lines, "\nApps:\n")
			lines = append(lines, testutil.GetTableOutputLines(
				[]string{"name", "requested state", "instances", "memory", "disk", "urls"}, apps)...)
			if len(services) > 0 {
				lines = append(lines, "\nServices:\n")
				lines = append(lines, testutil.GetTableOutputLines(
					[]string{"name", "service", "plan", "bound apps", "last operation"}, services)...)
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
				AccessToken("bearer test-token", nil).
				GetApps([]plugin_models.GetAppsModel{getGetAppsModel("test-mta-module-1", "started", 1, 1, 512, 1024, "test-1", "bosh-lite.com")}, nil).
				GetServices([]plugin_models.GetServices_Model{getGetServicesModel("test-service-1", "test", "free", "create", "succeeded", []string{"test-mta-module-1"})}, nil).
				Build()
			mtaV2Client := mtaV2fake.NewFakeMtaV2ClientBuilder().
				GetMtas("any_mtaId", &namespace, "any_spaceGuid", nil, nil).Build()
			clientFactory = commands.NewTestClientFactory(nil, mtaV2Client, nil)
			command = &commands.MtaCommand{}
			testTokenFactory := commands.NewTestTokenFactory(cliConnection)
			deployServiceURLCalculator := util_fakes.NewDeployServiceURLFakeCalculator("deploy-service.test.ondemand.com")

			command.InitializeAll(name, cliConnection, testutil.NewCustomTransport(200), nil, clientFactory, testTokenFactory, deployServiceURLCalculator, configuration.NewSnapshot())
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

		// can't connect to backend - error
		Context("when can't connect to backend", func() {
			const host = "x"
			It("should print an error and exit with a non-zero status", func() {
				clientFactory.MtaV2Client = mtaV2fake.NewFakeMtaV2ClientBuilder().
					GetMtasForThisSpace("test", nil, nil, fmt.Errorf("Get https://%s/rest/test/test/mta: dial tcp: lookup %s: no such host", host, host)).Build()
				output, status := oc.CaptureOutputAndStatus(func() int {
					return command.Execute([]string{"test", "-u", host}).ToInt()
				})
				ex.ExpectFailureOnLine(status, output, "Could not get multi-target app test:", 2)
			})
		})

		// backend returns a "not found" response - error
		Context("with an error response returned by the backend", func() {
			It("should print an error and exit with a non-zero status", func() {
				clientFactory.MtaV2Client = mtaV2fake.NewFakeMtaV2ClientBuilder().
					GetMtasForThisSpace("test", nil, nil, newClientError(404, "404 Not Found", `MTA with id "test" does not exist`)).Build()
				output, status := oc.CaptureOutputAndStatus(func() int {
					return command.Execute([]string{"test"}).ToInt()
				})
				ex.ExpectFailureOnLine(status, output, "Multi-target app test not found", 1)
			})
		})

		// backend returns a non-empty response - success
		Context("with a non-empty response returned by the backend", func() {
			It("should print a information about the deployed MTA and exit with zero status", func() {
				clientFactory.MtaV2Client = mtaV2fake.NewFakeMtaV2ClientBuilder().
					GetMtasForThisSpace("test-mta-id", &namespace, []*models.Mta{testutil.GetMta("test-mta-id", "test-version", "namespace", []*models.Module{
						testutil.GetMtaModule("test-mta-module-1", []string{}, []string{})},
						[]string{"test-service-1"})}, nil).Build()
				output, status := oc.CaptureOutputAndStatus(func() int {
					return command.Execute([]string{"test-mta-id", "--namespace", namespace}).ToInt()
				})
				ex.ExpectSuccessWithOutput(status, output,
					getOutputLines("test-mta-id", "test-version", "namespace",
						[][]string{[]string{"test-mta-module-1", "started", "1/1", "512M", "1G", "test-1.bosh-lite.com"}},
						[][]string{[]string{"test-service-1", "test", "free", "test-mta-module-1", "create succeeded"}}))
			})
		})

		// backend returns a non-empty response - success
		Context("with a non-empty response returned by the backend but no namespace", func() {
			It("should print a information about the deployed MTA and exit with zero status", func() {
				clientFactory.MtaV2Client = mtaV2fake.NewFakeMtaV2ClientBuilder().
					GetMtasForThisSpace("test-mta-id", nil, []*models.Mta{testutil.GetMta("test-mta-id", "test-version", "", []*models.Module{
						testutil.GetMtaModule("test-mta-module-1", []string{}, []string{})},
						[]string{"test-service-1"})}, nil).Build()
				output, status := oc.CaptureOutputAndStatus(func() int {
					return command.Execute([]string{"test-mta-id"}).ToInt()
				})
				ex.ExpectSuccessWithOutput(status, output,
					getOutputLines("test-mta-id", "test-version", "",
						[][]string{[]string{"test-mta-module-1", "started", "1/1", "512M", "1G", "test-1.bosh-lite.com"}},
						[][]string{[]string{"test-service-1", "test", "free", "test-mta-module-1", "create succeeded"}}))
			})
		})

		// backend returns a non-empty response - success
		Context("with a non-empty response without services returned by the backend", func() {
			It("should print information about the deployed MTA and exit with zero status", func() {
				clientFactory.MtaV2Client = mtaV2fake.NewFakeMtaV2ClientBuilder().
					GetMtasForThisSpace("test-mta-id", &namespace, []*models.Mta{testutil.GetMta("test-mta-id", "test-version", "namespace", []*models.Module{
						testutil.GetMtaModule("test-mta-module-1", []string{}, []string{})},
						[]string{})}, nil).Build()
				output, status := oc.CaptureOutputAndStatus(func() int {
					return command.Execute([]string{"test-mta-id", "--namespace", namespace}).ToInt()
				})
				ex.ExpectSuccessWithOutput(status, output,
					getOutputLines("test-mta-id", "test-version", "namespace",
						[][]string{[]string{"test-mta-module-1", "started", "1/1", "512M", "1G", "test-1.bosh-lite.com"}},
						[][]string{}))
			})
		})

		// backend returns a non-empty response - success
		Context("with a non-empty response with unknown MTA version", func() {
			It("should print information about the deployed MTA and exit with zero status", func() {
				clientFactory.MtaV2Client = mtaV2fake.NewFakeMtaV2ClientBuilder().
					GetMtasForThisSpace("test-mta-id", nil, []*models.Mta{testutil.GetMta("test-mta-id", "0.0.0-unknown", "", []*models.Module{}, []string{})}, nil).Build()
				output, status := oc.CaptureOutputAndStatus(func() int {
					return command.Execute([]string{"test-mta-id"}).ToInt()
				})
				ex.ExpectSuccessWithOutput(status, output,
					getOutputLines("test-mta-id", "?", "", [][]string{}, [][]string{}))
			})
		})

		// backend returns a non-empty response - success
		Context("with a non-empty response without services and apps returned by the backend", func() {
			It("should print information about the deployed MTA and exit with zero status", func() {
				clientFactory.MtaV2Client = mtaV2fake.NewFakeMtaV2ClientBuilder().
					GetMtasForThisSpace("test-mta-id", &namespace, []*models.Mta{testutil.GetMta("test-mta-id", "test-version", "namespace", []*models.Module{}, []string{})}, nil).Build()
				output, status := oc.CaptureOutputAndStatus(func() int {
					return command.Execute([]string{"test-mta-id", "--namespace", namespace}).ToInt()
				})
				ex.ExpectSuccessWithOutput(status, output,
					getOutputLines("test-mta-id", "test-version", "namespace", [][]string{}, [][]string{}))
			})
		})
	})
})

func getGetAppsModel(name, state string, runningInstances, totalInstances int,
	memory, diskQuota int64, host, domain string) plugin_models.GetAppsModel {
	return plugin_models.GetAppsModel{
		Name:             name,
		State:            state,
		RunningInstances: runningInstances,
		TotalInstances:   totalInstances,
		Memory:           memory,
		DiskQuota:        diskQuota,
		Routes: []plugin_models.GetAppsRouteSummary{
			plugin_models.GetAppsRouteSummary{
				Host: host,
				Domain: plugin_models.GetAppsDomainFields{
					Name: domain,
				},
			},
		},
	}
}

func getGetServicesModel(name, offering, plan, opType, opState string, boundApplications []string) plugin_models.GetServices_Model {
	return plugin_models.GetServices_Model{
		Name: name,
		Service: plugin_models.GetServices_ServiceFields{
			Name: offering,
		},
		ServicePlan: plugin_models.GetServices_ServicePlan{
			Name: plan,
		},
		LastOperation: plugin_models.GetServices_LastOperation{
			Type:  opType,
			State: opState,
		},
		ApplicationNames: boundApplications,
	}
}
