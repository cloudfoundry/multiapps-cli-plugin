package restclient_test

import (
	"net/http"
	"time"

	"github.com/cloudfoundry/multiapps-cli-plugin/clients/baseclient"
	"github.com/cloudfoundry/multiapps-cli-plugin/clients/restclient"
	"github.com/cloudfoundry/multiapps-cli-plugin/testutil"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("RetryableRestClient", func() {
	expectedRetriesCount := 4

	Describe("PurgeConfiguration", func() {
		Context("when the backend returns not 204 No Content", func() {
			It("should return an error", func() {
				client := newRetryableRestClient(http.StatusInternalServerError)
				err := client.PurgeConfiguration("org", "space")
				Expect(err).To(HaveOccurred())
			})
		})

		Context("when the backend returns 204 No Content", func() {
			It("should not return an error", func() {
				client := newRetryableRestClient(http.StatusNoContent)
				err := client.PurgeConfiguration("org", "space")
				Expect(err).ToNot(HaveOccurred())
			})
		})
		Context("with an 500 server error returned by the backend", func() {
			It("should return the error, zero result and retry expectedRetriesCount times", func() {
				client := newRetryableRestClient(500)
				retryableRestClient := client.(restclient.RetryableRestClient)
				err := retryableRestClient.PurgeConfiguration("org", "space")
				testutil.ExpectError(err)
				mockRestClient := retryableRestClient.RestClient.(*MockRestClient)
				Expect(mockRestClient.GetRetriesCount()).To(Equal(expectedRetriesCount))
			})
		})
		Context("with an 502 server error returned by the backend", func() {
			It("should return the error, zero result and retry expectedRetriesCount times", func() {
				client := newRetryableRestClient(502)
				retryableRestClient := client.(restclient.RetryableRestClient)
				err := retryableRestClient.PurgeConfiguration("org", "space")
				testutil.ExpectError(err)
				mockRestClient := retryableRestClient.RestClient.(*MockRestClient)
				Expect(mockRestClient.GetRetriesCount()).To(Equal(expectedRetriesCount))
			})
		})
	})
})

func newRetryableRestClient(statusCode int) restclient.RestClientOperations {
	tokenFactory := baseclient.NewCustomTokenFactory("test-token")
	roundTripper := testutil.NewCustomTransport(statusCode)
	return NewMockRetryableRestClient("http://localhost:1000", roundTripper, tokenFactory)
}

// NewRetryableRestClient creates a new retryable REST client
func NewMockRetryableRestClient(host string, rt http.RoundTripper, tokenFactory baseclient.TokenFactory) restclient.RestClientOperations {
	restClient := NewMockRestClient(host, rt, tokenFactory)
	return restclient.RetryableRestClient{RestClient: restClient, MaxRetriesCount: 3, RetryInterval: time.Microsecond * 1}
}

type MockRestClient struct {
	restClient   restclient.RestClientOperations
	retriesCount int
}

// NewMockRestClient creates a new Rest client
func NewMockRestClient(host string, rt http.RoundTripper, tokenFactory baseclient.TokenFactory) restclient.RestClientOperations {
	restClient := restclient.NewRestClient(host, rt, tokenFactory)
	return &MockRestClient{restClient, 0}
}

// PurgeConfiguration purges configuration
func (c *MockRestClient) PurgeConfiguration(org, space string) error {
	c.retriesCount++
	return c.restClient.PurgeConfiguration(org, space)
}

// GetRetriesCount returns retries count
func (c *MockRestClient) GetRetriesCount() int {
	return c.retriesCount
}
