package restclient_test

import (
	"net/http"
	"net/http/cookiejar"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	baseclient "github.com/SAP/cf-mta-plugin/clients/baseclient"
	restclient "github.com/SAP/cf-mta-plugin/clients/restclient"
	"github.com/SAP/cf-mta-plugin/testutil"
)

var _ = Describe("RestClient", func() {

	Describe("GetOperations", func() {
		Context("with valid ongoing operations returned by the backend", func() {
			It("should return the ongoing operations returned by the backend", func() {
				client := newRestClient(200, testutil.OperationsResult)
				result, err := client.GetOperations(nil, nil)
				testutil.ExpectNoErrorAndResult(err, result, testutil.OperationsResult)
			})
		})
		Context("with an error returned by the backend", func() {
			It("should return the error and a zero result", func() {
				client := newRestClient(404, nil)
				result, err := client.GetOperations(nil, nil)
				testutil.ExpectErrorAndZeroResult(err, result)
			})
		})
	})

	Describe("GetComponents", func() {
		Context("with valid deployed components returned by the backend", func() {
			It("should return the deployed components returned by the backend", func() {
				client := newRestClient(200, testutil.ComponentsResult)
				result, err := client.GetComponents()
				testutil.ExpectNoErrorAndResult(err, result, &testutil.ComponentsResult)
			})
		})
		Context("with an error returned by the backend", func() {
			It("should return the error and a zero result", func() {
				client := newRestClient(404, nil)
				result, err := client.GetComponents()
				testutil.ExpectErrorAndZeroResult(err, result)
			})
		})
	})

	Describe("GetMta", func() {
		Context("with a valid MTA returned by the backend", func() {
			It("should return the MTA returned by the backend", func() {
				client := newRestClient(200, testutil.MtaResult)
				result, err := client.GetMta(testutil.MtaID)
				testutil.ExpectNoErrorAndResult(err, result, &testutil.MtaResult)
			})
		})
		Context("with an error returned by the backend", func() {
			It("should return the error and a zero result", func() {
				client := newRestClient(404, nil)
				result, err := client.GetMta(testutil.MtaID)
				testutil.ExpectErrorAndZeroResult(err, result)
			})
		})
	})

	Describe("PurgeConfiguration", func() {
		Context("when the backend returns not 204 No Content", func() {
			It("should return an error", func() {
				client := newRestClient(http.StatusInternalServerError, nil)
				err := client.PurgeConfiguration("org", "space")
				Ω(err).Should(HaveOccurred())
			})
		})

		Context("when the backend returns 204 No Content", func() {
			It("should not return an error", func() {
				client := newRestClient(http.StatusNoContent, nil)
				err := client.PurgeConfiguration("org", "space")
				Ω(err).ShouldNot(HaveOccurred())
			})
		})
	})
})

func newRestClient(statusCode int, v interface{}) restclient.RestClientOperations {
	tokenFactory := baseclient.NewCustomTokenFactory("test-token")
	cookieJar, _ := cookiejar.New(nil)
	roundTripper := testutil.NewCustomTransport(statusCode, v)
	return restclient.NewRestClient("http://localhost:1000", "test-org", "test-space", roundTripper, cookieJar, tokenFactory)
}
