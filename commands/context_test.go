package commands_test

import (
	"fmt"
	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/commands"

	"github.com/cloudfoundry/cli/cf/terminal"

	cli_fakes "github.com/cloudfoundry-incubator/multiapps-cli-plugin/cli/fakes"
	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/ui"
	plugin_fakes "github.com/cloudfoundry/cli/plugin/fakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Context", func() {
	const org = "test-org"
	const space = "test-space"
	const user = "test-user"

	var fakeCliConnection *plugin_fakes.FakeCliConnection

	BeforeEach(func() {
		ui.DisableTerminalOutput(true)
		fakeCliConnection = cli_fakes.NewFakeCliConnectionBuilder().
			CurrentOrg("test-org-guid", org, nil).
			CurrentSpace("test-space-guid", space, nil).
			Username(user, nil).
			AccessToken("bearer test-token", nil).Build()
	})

	Describe("GetOrg", func() {
		Context("with valid org returned by the CLI connection", func() {
			It("should not exit or report any errors", func() {
				o, err := commands.GetOrg(fakeCliConnection)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(o.Name).To(Equal(org))
			})
		})
		Context("with no org returned by the CLI connection", func() {
			It("should print an error and exit with a non-zero status", func() {
				fakeCliConnection := cli_fakes.NewFakeCliConnectionBuilder().
					CurrentOrg("", "", nil).Build()
				_, err := commands.GetOrg(fakeCliConnection)
				Expect(err).To(MatchError(fmt.Errorf("No org and space targeted, use '%s' to target an org and a space", terminal.CommandColor("cf target -o ORG -s SPACE"))))
			})
		})
	})

	Describe("GetSpace", func() {
		Context("with valid space returned by the CLI connection", func() {
			It("should not exit or report any errors", func() {
				s, err := commands.GetSpace(fakeCliConnection)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(s.Name).To(Equal(space))
			})
		})
		Context("with no space returned by the CLI connection", func() {
			It("should print an error and exit with a non-zero status", func() {
				fakeCliConnection := cli_fakes.NewFakeCliConnectionBuilder().
					CurrentSpace("", "", nil).Build()
				_, err := commands.GetSpace(fakeCliConnection)
				Expect(err).To(MatchError(fmt.Errorf("No space targeted, use '%s' to target a space", terminal.CommandColor("cf target -s"))))
			})
		})
	})

	Describe("GetUsername", func() {
		Context("with valid username returned by the CLI connection", func() {
			It("should not exit or report any errors", func() {
				Expect(commands.GetUsername(fakeCliConnection)).To(Equal(user))
			})
		})
		Context("with no space returned by the CLI connection", func() {
			It("should print an error and exit with a non-zero status", func() {
				fakeCliConnection := cli_fakes.NewFakeCliConnectionBuilder().
					Username("", nil).Build()
				_, err := commands.GetUsername(fakeCliConnection)
				Expect(err).To(MatchError(fmt.Errorf("Not logged in. Use '%s' to log in.", terminal.CommandColor("cf login"))))
			})
		})
	})
})
