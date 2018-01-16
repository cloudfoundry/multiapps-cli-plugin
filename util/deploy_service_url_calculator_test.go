package util

import (
	cli_fakes "github.com/SAP/cf-mta-plugin/cli/fakes"
	plugin_models "github.com/cloudfoundry/cli/plugin/models"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("DeployServiceURLCalculator", func() {
	Describe("ComputeDeployServiceURL", func() {
		const domainName = "test.ondemand.com"
		const spaceName = "test-space"

		Context("when a space is targeted and there is one shared domain", func() {
			It("should return a URL constructed based on the first shared domain", func() {
				cliConnection := cli_fakes.NewFakeCliConnectionBuilder().
					CurrentSpace("", spaceName, nil).
					GetSpace(spaceName, plugin_models.GetSpace_Model{
						GetSpaces_Model: plugin_models.GetSpaces_Model{
							Name: spaceName,
						},
						Domains: []plugin_models.GetSpace_Domains{
							plugin_models.GetSpace_Domains{Name: "custom.test.ondemand.com", Shared: false},
							plugin_models.GetSpace_Domains{Name: "test.ondemand.com", Shared: true},
						},
					}, nil).Build()
				deployServiceURLCalculator := NewDeployServiceURLCalculator(cliConnection)
				Expect(deployServiceURLCalculator.ComputeDeployServiceURL()).To(Equal("deploy-service.test.ondemand.com"))
			})
		})
		Context("when a space is targeted and there are two shared domains", func() {
			It("should return a URL constructed based on the first shared domain", func() {
				cliConnection := cli_fakes.NewFakeCliConnectionBuilder().
					CurrentSpace("", spaceName, nil).
					GetSpace(spaceName, plugin_models.GetSpace_Model{
						GetSpaces_Model: plugin_models.GetSpaces_Model{
							Name: spaceName,
						},
						Domains: []plugin_models.GetSpace_Domains{
							plugin_models.GetSpace_Domains{Name: "custom.test.ondemand.com", Shared: false},
							plugin_models.GetSpace_Domains{Name: "test1.ondemand.com", Shared: true},
							plugin_models.GetSpace_Domains{Name: "test2.ondemand.com", Shared: true},
						},
					}, nil).Build()
				deployServiceURLCalculator := NewDeployServiceURLCalculator(cliConnection)
				Expect(deployServiceURLCalculator.ComputeDeployServiceURL()).To(Equal("deploy-service.test1.ondemand.com"))
			})
		})
		Context("when a space is targeted and there are no shared domains", func() {
			It("should return an error", func() {
				cliConnection := cli_fakes.NewFakeCliConnectionBuilder().
					CurrentSpace("", spaceName, nil).
					GetSpace(spaceName, plugin_models.GetSpace_Model{
						GetSpaces_Model: plugin_models.GetSpaces_Model{
							Name: spaceName,
						},
						Domains: []plugin_models.GetSpace_Domains{
							plugin_models.GetSpace_Domains{Name: "custom.test.ondemand.com", Shared: false},
						},
					}, nil).Build()
				deployServiceURLCalculator := NewDeployServiceURLCalculator(cliConnection)
				_, err := deployServiceURLCalculator.ComputeDeployServiceURL()
				Expect(err).Should(MatchError("Could not find any shared domains in space: " + spaceName))
			})
		})
		Context("when no space is targeted", func() {
			It("should return an error", func() {
				cliConnection := cli_fakes.NewFakeCliConnectionBuilder().
					CurrentSpace("", "", nil).
					Build()
				deployServiceURLCalculator := NewDeployServiceURLCalculator(cliConnection)
				_, err := deployServiceURLCalculator.ComputeDeployServiceURL()
				Expect(err).Should(MatchError("No space targeted, use 'cf target -s SPACE' to target a space."))
			})
		})
	})
})
