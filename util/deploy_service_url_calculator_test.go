package util_test

import (
	cfrestclient_fakes "github.com/SAP/cf-mta-plugin/clients/cfrestclient/fakes"
	"github.com/SAP/cf-mta-plugin/clients/models"
	"github.com/SAP/cf-mta-plugin/util"
	fakes "github.com/SAP/cf-mta-plugin/util/fakes"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("DeployServiceURLCalculator", func() {
	Describe("ComputeDeployServiceURL", func() {
		const domainName = "test.ondemand.com"
		const spaceName = "test-space"

		Context("when a space is targeted and there is one shared domain", func() {
			It("should return a URL constructed based on the first shared domain", func() {
				domains := []models.SharedDomain{
					models.SharedDomain{Name: "test.ondemand.com"},
				}
				fakeHttpExecutor := fakes.NewFakeHttpGetExecutor(200, nil)
				deployServiceURLCalculator := util.NewDeployServiceURLCalculatorWithHttpExecutor(cfrestclient_fakes.NewFakeCloudFoundryClient(domains, nil), fakeHttpExecutor)
				Expect(deployServiceURLCalculator.ComputeDeployServiceURL()).To(Equal("deploy-service.test.ondemand.com"))
			})
		})
		Context("when a space is targeted and there are two shared domains", func() {
			It("should return a URL constructed based on the first shared domain", func() {
				domains := []models.SharedDomain{
					models.SharedDomain{Name: "test1.ondemand.com"},
					models.SharedDomain{Name: "test2.ondemand.com"},
				}
				fakeHttpExecutor := fakes.NewFakeHttpGetExecutor(200, nil)
				deployServiceURLCalculator := util.NewDeployServiceURLCalculatorWithHttpExecutor(cfrestclient_fakes.NewFakeCloudFoundryClient(domains, nil), fakeHttpExecutor)
				Expect(deployServiceURLCalculator.ComputeDeployServiceURL()).To(Equal("deploy-service.test1.ondemand.com"))
			})
		})
		Context("when a space is targeted and there are no shared domains", func() {
			It("should return an error", func() {
				deployServiceURLCalculator := util.NewDeployServiceURLCalculator(cfrestclient_fakes.NewFakeCloudFoundryClient([]models.SharedDomain{}, nil))
				_, err := deployServiceURLCalculator.ComputeDeployServiceURL()
				Expect(err).Should(MatchError("Could not find any shared domains"))
			})
		})
	})
})
