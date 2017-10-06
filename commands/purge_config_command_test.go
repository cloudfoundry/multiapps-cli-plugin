package commands_test

import (
	"fmt"

	plugin_fakes "github.com/cloudfoundry/cli/plugin/fakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/SAP/cf-mta-plugin/clients/models"
	restfake "github.com/SAP/cf-mta-plugin/clients/restclient/fakes"
	"github.com/SAP/cf-mta-plugin/commands"
	cmd_fakes "github.com/SAP/cf-mta-plugin/commands/fakes"
	"github.com/SAP/cf-mta-plugin/testutil"
	"github.com/SAP/cf-mta-plugin/ui"
)

var _ = Describe("PurgeConfigCommand", func() {

	Describe("Execute", func() {
		const (
			org   = "test-org"
			space = "test-space"
			user  = "test-user"
		)

		var name string
		var cliConnection *plugin_fakes.FakeCliConnection
		var clientFactory *commands.TestClientFactory
		var command *commands.PurgeConfigCommand
		var testTokenFactory *commands.TestTokenFactory

		var oc testutil.OutputCapturer
		var ex testutil.Expector

		BeforeEach(func() {
			ui.DisableTerminalOutput(true)
			name = "purge-mta-config"
			cliConnection = cmd_fakes.NewFakeCliConnectionBuilder().
				CurrentOrg("test-org-guid", org, nil).
				CurrentSpace("test-space-guid", space, nil).
				Username(user, nil).
				AccessToken("bearer test-token", nil).
				APIEndpoint("https://api.test.ondemand.com", nil).Build()

			testTokenFactory = commands.NewTestTokenFactory(cliConnection)
			clientFactory = commands.NewTestClientFactory(nil, nil, nil)

			command = &commands.PurgeConfigCommand{}
			command.InitializeAll(name, cliConnection, testutil.NewCustomTransport(200, nil), nil, clientFactory, testTokenFactory)

			oc = testutil.NewUIOutputCapturer()
			ex = testutil.NewUIExpector()
		})

		Context("with an unknown flag", func() {
			It("should print incorrect usage, call cf help, and exit with a non-zero status", func() {
				output, status := oc.CaptureOutputAndStatus(func() int {
					return command.Execute([]string{"-a"}).ToInt()
				})
				ex.ExpectFailure(status, output, "Incorrect usage. Unknown or wrong flag.")
				Expect(cliConnection.CliCommandArgsForCall(0)).To(Equal([]string{"help", name}))
			})
		})

		Context("with wrong arguments", func() {
			It("should print incorrect usage, call cf help, and exit with a non-zero status", func() {
				output, status := oc.CaptureOutputAndStatus(func() int {
					return command.Execute([]string{"x", "y", "z"}).ToInt()
				})
				ex.ExpectFailure(status, output, "Incorrect usage. Wrong arguments.")
				Expect(cliConnection.CliCommandArgsForCall(0)).To(Equal([]string{"help", name}))
			})
		})

		Context("with an error response returned by the backend", func() {
			It("should print an error and exit with a non-zero status", func() {
				clientFactory.RestClient = restfake.NewFakeRestClientBuilder().
					GetComponents(nil, fmt.Errorf("unknown error (status 404)")).Build()
				output, status := oc.CaptureOutputAndStatus(func() int {
					return command.Execute([]string{}).ToInt()
				})
				ex.ExpectFailureOnLine(status, output, "Could not purge configuration", 1)
			})
		})

		Context("with an success response returned by the backend", func() {
			It("should print a message and exit with zero status", func() {
				clientFactory.RestClient = restfake.NewFakeRestClientBuilder().
					GetComponents(testutil.GetComponents([]*models.Mta{}, []string{}), nil).Build()
				output, status := oc.CaptureOutputAndStatus(func() int {
					return command.Execute([]string{}).ToInt()
				})
				expectedOutput := []string{
					"Purging configuration entries in org test-org / space test-space as test-user\n",
					"OK\n",
				}
				ex.ExpectSuccessWithOutput(status, output, expectedOutput)
			})
		})

	})
})
