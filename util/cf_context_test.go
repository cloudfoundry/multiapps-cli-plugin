package util

import (
	"fmt"

	"code.cloudfoundry.org/cli/cf/terminal"

	plugin_fakes "code.cloudfoundry.org/cli/plugin/pluginfakes"
	cli_fakes "github.com/cloudfoundry-incubator/multiapps-cli-plugin/cli/fakes"
	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/ui"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("CloudFoundryContext", func() {
	const (
		org   = "test-org"
		space = "test-space"
		user  = "test-user"
	)

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
				cfContext := NewCloudFoundryContext(fakeCliConnection)
				org, err := cfContext.GetOrg()
				Expect(err).ShouldNot(HaveOccurred())
				Expect(org).To(Equal(org))
			})
		})
		Context("with no org returned by the CLI connection", func() {
			It("should print an error and exit with a non-zero status", func() {
				fakeCliConnection := cli_fakes.NewFakeCliConnectionBuilder().
					CurrentOrg("", "", nil).Build()
				cfContext := NewCloudFoundryContext(fakeCliConnection)
				_, err := cfContext.GetOrg()
				Expect(err).To(MatchError(fmt.Errorf("No org and space targeted, use %q to target an org and a space", terminal.CommandColor("cf target -o ORG -s SPACE"))))
			})
		})
	})

	Describe("GetSpace", func() {
		Context("with valid space returned by the CLI connection", func() {
			It("should not exit or report any errors", func() {
				cfContext := NewCloudFoundryContext(fakeCliConnection)
				space, err := cfContext.GetSpace()
				Expect(err).ShouldNot(HaveOccurred())
				Expect(space).To(Equal(space))
			})
		})
		Context("with no space returned by the CLI connection", func() {
			It("should print an error and exit with a non-zero status", func() {
				fakeCliConnection := cli_fakes.NewFakeCliConnectionBuilder().
					CurrentSpace("", "", nil).Build()
				cfContext := NewCloudFoundryContext(fakeCliConnection)
				_, err := cfContext.GetSpace()
				Expect(err).To(MatchError(fmt.Errorf("No space targeted, use %q to target a space", terminal.CommandColor("cf target -s SPACE"))))
			})
		})
	})

	Describe("GetUsername", func() {
		Context("with valid username returned by the CLI connection", func() {
			It("should not exit or report any errors", func() {
				cfContext := NewCloudFoundryContext(fakeCliConnection)
				user, err := cfContext.GetUsername()
				Expect(err).ShouldNot(HaveOccurred())
				Expect(user).To(Equal(user))
			})
		})
		Context("with no space returned by the CLI connection", func() {
			It("should print an error and exit with a non-zero status", func() {
				fakeCliConnection := cli_fakes.NewFakeCliConnectionBuilder().
					Username("", nil).Build()
				cfContext := NewCloudFoundryContext(fakeCliConnection)
				_, err := cfContext.GetUsername()
				Expect(err).To(MatchError(fmt.Errorf("Not logged in. Use %q to log in.", terminal.CommandColor("cf login"))))
			})
		})
	})
})
