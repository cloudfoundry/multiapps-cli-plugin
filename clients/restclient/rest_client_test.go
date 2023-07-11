package restclient_test

import (
	"github.com/cloudfoundry/multiapps-cli-plugin/clients/baseclient"
	"github.com/cloudfoundry/multiapps-cli-plugin/clients/restclient"
	"github.com/cloudfoundry/multiapps-cli-plugin/testutil"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"net/http"
)

var _ = Describe("RestClient", func() {

	Describe("PurgeConfiguration", func() {
		Context("when the backend returns not 204 No Content", func() {
			It("should return an error", func() {
				client := newRestClient(http.StatusInternalServerError)
				err := client.PurgeConfiguration("org", "space")
				Ω(err).Should(HaveOccurred())
			})
		})

		Context("when the backend returns 204 No Content", func() {
			It("should not return an error", func() {
				client := newRestClient(http.StatusNoContent)
				err := client.PurgeConfiguration("org", "space")
				Ω(err).ShouldNot(HaveOccurred())
			})
		})
	})
})

func newRestClient(statusCode int) restclient.RestClientOperations {
	tokenFactory := baseclient.NewCustomTokenFactory("test-token")
	roundTripper := testutil.NewCustomTransport(statusCode)
	return restclient.NewRestClient("http://localhost:1000", roundTripper, tokenFactory)
}
