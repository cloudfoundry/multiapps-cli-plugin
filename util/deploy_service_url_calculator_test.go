package util_test

import (
	cfrestclient_fakes "github.com/cloudfoundry-incubator/multiapps-cli-plugin/clients/cfrestclient/fakes"
	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/clients/models"
	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/util"
	fakes "github.com/cloudfoundry-incubator/multiapps-cli-plugin/util/fakes"

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
				fakeHttpExecutor := fakes.NewFakeHttpGetExecutor(map[string]int {
					"https://deploy-service.test.ondemand.com/public/ping": 200,
				})
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
				fakeHttpExecutor := fakes.NewFakeHttpGetExecutor(map[string]int {
					"https://deploy-service.test1.ondemand.com/public/ping": 200,
					"https://deploy-service.test2.ondemand.com/public/ping": 404,
				})
				deployServiceURLCalculator := util.NewDeployServiceURLCalculatorWithHttpExecutor(cfrestclient_fakes.NewFakeCloudFoundryClient(domains, nil), fakeHttpExecutor)
				Expect(deployServiceURLCalculator.ComputeDeployServiceURL()).To(Equal("deploy-service.test1.ondemand.com"))
			})
		})
		Context("when a space is targeted and there are two shared domains, but the Deploy Service does not respond on either of them", func() {
			It("should return an error", func() {
				domains := []models.SharedDomain{
					models.SharedDomain{Name: "test1.ondemand.com"},
					models.SharedDomain{Name: "test2.ondemand.com"},
				}
				fakeHttpExecutor := fakes.NewFakeHttpGetExecutor(map[string]int {
					"https://deploy-service.test1.ondemand.com/public/ping": 404,
					"https://deploy-service.test2.ondemand.com/public/ping": 404,
				})
				deployServiceURLCalculator := util.NewDeployServiceURLCalculatorWithHttpExecutor(cfrestclient_fakes.NewFakeCloudFoundryClient(domains, nil), fakeHttpExecutor)
				_, err := deployServiceURLCalculator.ComputeDeployServiceURL()
				Expect(err).Should(MatchError("The Deploy Service does not respond on any of the default URLs:\ndeploy-service.test1.ondemand.com\ndeploy-service.test2.ondemand.com\n\nYou can use the command line option -u or the DEPLOY_SERVICE_URL environment variable to specify a custom URL explicitly."))
			})
		})
		Context("when a space is targeted and there are two shared domains, but the Deploy Service is broken on one of them", func() {
			It("should return an error", func() {
				domains := []models.SharedDomain{
					models.SharedDomain{Name: "test1.ondemand.com"},
					models.SharedDomain{Name: "test2.ondemand.com"},
				}
				fakeHttpExecutor := fakes.NewFakeHttpGetExecutor(map[string]int {
					"https://deploy-service.test1.ondemand.com/public/ping": 404,
					"https://deploy-service.test2.ondemand.com/public/ping": 500,
				})
				deployServiceURLCalculator := util.NewDeployServiceURLCalculatorWithHttpExecutor(cfrestclient_fakes.NewFakeCloudFoundryClient(domains, nil), fakeHttpExecutor)
				_, err := deployServiceURLCalculator.ComputeDeployServiceURL()
				Expect(err).Should(MatchError("The Deploy Service does not respond on any of the default URLs:\ndeploy-service.test1.ondemand.com\ndeploy-service.test2.ondemand.com\n\nYou can use the command line option -u or the DEPLOY_SERVICE_URL environment variable to specify a custom URL explicitly."))
			})
		})
		Context("when a space is targeted and there are no shared domains", func() {
			It("should return an error", func() {
				deployServiceURLCalculator := util.NewDeployServiceURLCalculator(cfrestclient_fakes.NewFakeCloudFoundryClient([]models.SharedDomain{}, nil))
				_, err := deployServiceURLCalculator.ComputeDeployServiceURL()
				Expect(err).Should(MatchError("Could not compute the Deploy Service's URL as there are no shared domains on the landscape."))
			})
		})
	})
})
