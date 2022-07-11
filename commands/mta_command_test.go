package commands_test

import (
	"fmt"

	plugin_fakes "code.cloudfoundry.org/cli/plugin/pluginfakes"
	cli_fakes "github.com/cloudfoundry-incubator/multiapps-cli-plugin/cli/fakes"
	cf_client_fakes "github.com/cloudfoundry-incubator/multiapps-cli-plugin/clients/cfrestclient/fakes"
	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/clients/models"
	mtaV2fake "github.com/cloudfoundry-incubator/multiapps-cli-plugin/clients/mtaclient_v2/fakes"
	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/commands"
	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/testutil"
	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/ui"
	util_fakes "github.com/cloudfoundry-incubator/multiapps-cli-plugin/util/fakes"
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
		var cfClient cf_client_fakes.FakeCloudFoundryClient
		var command *commands.MtaCommand
		var oc = testutil.NewUIOutputCapturer()
		var ex = testutil.NewUIExpector()

		var getOutputLines = func(mtaID, version, namespace string, apps, services [][]string) []string {
			var lines []string
			lines = append(lines,
				fmt.Sprintf("Showing health and status for multi-target app %s in org %s / space %s as %s...", mtaID, org, space, user))
			lines = append(lines, "OK")
			lines = append(lines, fmt.Sprintf("Version: %s", version))
			lines = append(lines, fmt.Sprintf("Namespace: %s", namespace))
			lines = append(lines, "")
			lines = append(lines, "Apps:")
			lines = append(lines, testutil.GetTableOutputLines(
				[]string{"name", "requested state", "instances", "memory", "disk", "urls"}, apps)...)
			lines = append(lines, "")
			lines = append(lines, "Services:")
			if len(services) > 0 {
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
				APIEndpoint("https://example.com", nil).
				Build()
			cfClient = cf_client_fakes.FakeCloudFoundryClient{
				Apps:            getApps("test-mta-module-1", "started"),
				AppProcessStats: getProcessStats(1, 1, 512*1024*1024, 1024*1024*1024),
				AppRoutes:       getAppRoutes("test-1", "bosh-lite.com"),
				Services:        getServices("test-service-1", "test", "free", "create", "succeeded"),
				ServiceBindings: getServiceBindings([]string{"test-mta-module-1"}),
			}
			mtaV2Client := mtaV2fake.NewFakeMtaV2ClientBuilder().
				GetMtas("any_mtaId", &namespace, "any_spaceGuid", nil, nil).Build()
			clientFactory = commands.NewTestClientFactory(nil, mtaV2Client, nil)
			command = commands.NewMtaCommand()
			testTokenFactory := commands.NewTestTokenFactory(cliConnection)
			deployServiceURLCalculator := util_fakes.NewDeployServiceURLFakeCalculator("deploy-service.test.ondemand.com")

			command.InitializeAll(name, cliConnection, testutil.NewCustomTransport(200), clientFactory, testTokenFactory, deployServiceURLCalculator)
			command.CfClient = cfClient
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
				ex.ExpectFailureOnLine(status, output, "Multi-target app test not found", 2)
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
						[][]string{{"test-mta-module-1", "started", "1/1", "512M", "1G", "test-1.bosh-lite.com"}},
						[][]string{{"test-service-1", "test", "free", "test-mta-module-1", "create succeeded"}}))
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
						[][]string{{"test-mta-module-1", "started", "1/1", "512M", "1G", "test-1.bosh-lite.com"}},
						[][]string{{"test-service-1", "test", "free", "test-mta-module-1", "create succeeded"}}))
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

func getApps(name, state string) []models.CloudFoundryApplication {
	return []models.CloudFoundryApplication{
		{
			Name:  name,
			Guid:  "app-guid",
			State: state,
		},
	}
}

func getProcessStats(runningInstances, totalInstances int64, memory, diskQuota int64) []models.ApplicationProcessStatistics {
	var processes []models.ApplicationProcessStatistics
	for i := int64(0); i < runningInstances; i++ {
		processMemory := memory / runningInstances
		if i == 0 {
			processMemory += memory % runningInstances
		}
		disk := diskQuota / runningInstances
		if i == 0 {
			disk += diskQuota % runningInstances
		}
		processes = append(processes, models.ApplicationProcessStatistics{
			State:  "RUNNING",
			Memory: processMemory,
			Disk:   disk,
		})
	}
	for i := runningInstances; i < totalInstances; i++ {
		processes = append(processes, models.ApplicationProcessStatistics{
			State:  "STOPPED",
			Memory: 0,
			Disk:   0,
		})
	}
	return processes
}

func getAppRoutes(host, domain string) []models.ApplicationRoute {
	return []models.ApplicationRoute{
		{
			Host: host,
			Url:  host + "." + domain,
		},
	}
}

func getServices(name, offering, plan, opType, opState string) []models.CloudFoundryServiceInstance {
	return []models.CloudFoundryServiceInstance{
		{
			Guid: "service-guid",
			Name: name,
			Type: "managed",
			LastOperation: models.LastOperation{
				Type:  opType,
				State: opState,
			},
			PlanGuid:  "plan-guid",
			SpaceGuid: "space-guid",
			Plan: models.ServicePlan{
				Guid:         "plan-guid",
				Name:         plan,
				OfferingGuid: "offering-guid",
			},
			Offering: models.ServiceOffering{
				Guid: "offering-guid",
				Name: offering,
			},
		},
	}
}

func getServiceBindings(boundApplications []string) []models.ServiceBinding {
	var bindings []models.ServiceBinding
	for _, appName := range boundApplications {
		bindings = append(bindings, models.ServiceBinding{
			Guid:    "binding-guid",
			AppGuid: "app-guid",
			AppName: appName,
		})
	}
	return bindings
}
