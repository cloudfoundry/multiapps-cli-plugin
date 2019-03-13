package restclient_test

import (
	"net/http"
	"net/http/cookiejar"
	"time"

	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/clients/baseclient"
	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/clients/restclient"
	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/testutil"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("RetryableRestClient", func() {
	expectedRetriesCount := 4

	Describe("PurgeConfiguration", func() {
		Context("when the backend returns not 204 No Content", func() {
			It("should return an error", func() {
				client := newRetryableRestClient(http.StatusInternalServerError, nil)
				err := client.PurgeConfiguration("org", "space")
				Expect(err).To(HaveOccurred())
			})
		})

		Context("when the backend returns 204 No Content", func() {
			It("should not return an error", func() {
				client := newRetryableRestClient(http.StatusNoContent, nil)
				err := client.PurgeConfiguration("org", "space")
				Expect(err).ToNot(HaveOccurred())
			})
		})
		Context("with an 500 server error returned by the backend", func() {
			It("should return the error, zero result and retry expectedRetriesCount times", func() {
				client := newRetryableRestClient(500, nil)
				retryableRestClient := client.(restclient.RetryableRestClient)
				err := retryableRestClient.PurgeConfiguration("org", "space")
				testutil.ExpectError(err)
				mockRestClient := retryableRestClient.RestClient.(*MockRestClient)
				Expect(mockRestClient.GetRetriesCount()).To(Equal(expectedRetriesCount))
			})
		})
		Context("with an 502 server error returned by the backend", func() {
			It("should return the error, zero result and retry expectedRetriesCount times", func() {
				client := newRetryableRestClient(502, nil)
				retryableRestClient := client.(restclient.RetryableRestClient)
				err := retryableRestClient.PurgeConfiguration("org", "space")
				testutil.ExpectError(err)
				mockRestClient := retryableRestClient.RestClient.(*MockRestClient)
				Expect(mockRestClient.GetRetriesCount()).To(Equal(expectedRetriesCount))
			})
		})
	})
})

func newRetryableRestClient(statusCode int, v interface{}) restclient.RestClientOperations {
	tokenFactory := baseclient.NewCustomTokenFactory("test-token")
	cookieJar, _ := cookiejar.New(nil)
	roundTripper := testutil.NewCustomTransport(statusCode, v)
	return NewMockRetryableRestClient("http://localhost:1000", "test-org", "test-space", roundTripper, cookieJar, tokenFactory)
}

// NewRetryableRestClient creates a new retryable REST client
func NewMockRetryableRestClient(host, org, space string, rt http.RoundTripper, jar http.CookieJar, tokenFactory baseclient.TokenFactory) restclient.RestClientOperations {
	restClient := NewMockRestClient(host, org, space, rt, jar, tokenFactory)
	return restclient.RetryableRestClient{RestClient: restClient, MaxRetriesCount: 3, RetryInterval: time.Microsecond * 1}
}

type MockRestClient struct {
	restClient   restclient.RestClientOperations
	retriesCount int
}

// NewMockRestClient creates a new Rest client
func NewMockRestClient(host, org, space string, rt http.RoundTripper, jar http.CookieJar, tokenFactory baseclient.TokenFactory) restclient.RestClientOperations {
	restClient := restclient.NewRestClient(host, org, space, rt, jar, tokenFactory)
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
