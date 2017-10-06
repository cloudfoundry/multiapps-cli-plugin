package restclient_test

import (
	"net/http"
	"net/http/cookiejar"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/SAP/cf-mta-plugin/clients/baseclient"
	"github.com/SAP/cf-mta-plugin/clients/models"
	"github.com/SAP/cf-mta-plugin/clients/restclient"
	"github.com/SAP/cf-mta-plugin/testutil"
)

var _ = Describe("RetryableRestClient", func() {
	expectedRetriesCount := 4

	Describe("GetOperations", func() {
		Context("with valid ongoing operations returned by the backend", func() {
			It("should return the ongoing operations returned by the backend", func() {
				client := newRetryableRestClient(200, testutil.OperationsResult)
				result, err := client.GetOperations(nil, nil)
				testutil.ExpectNoErrorAndResult(err, result, testutil.OperationsResult)
			})
		})
		Context("with an error returned by the backend", func() {
			It("should return the error and a zero result", func() {
				client := newRetryableRestClient(404, nil)
				result, err := client.GetOperations(nil, nil)
				testutil.ExpectErrorAndZeroResult(err, result)
			})
		})
		Context("with an 500 server error returned by the backend", func() {
			It("should return the error, zero result and retry expectedRetriesCount times", func() {
				client := newRetryableRestClient(500, nil)
				retryableRestClient := client.(restclient.RetryableRestClient)
				result, err := retryableRestClient.GetOperations(nil, nil)
				testutil.ExpectErrorAndZeroResult(err, result)
				mockRestClient := retryableRestClient.RestClient.(*MockRestClient)
				Expect(mockRestClient.GetRetriesCount()).To(Equal(expectedRetriesCount))
			})
		})
		Context("with an 502 server error returned by the backend", func() {
			It("should return the error, zero result and retry expectedRetriesCount times", func() {
				client := newRetryableRestClient(502, nil)
				retryableRestClient := client.(restclient.RetryableRestClient)
				result, err := retryableRestClient.GetOperations(nil, nil)
				testutil.ExpectErrorAndZeroResult(err, result)
				mockRestClient := retryableRestClient.RestClient.(*MockRestClient)
				Expect(mockRestClient.GetRetriesCount()).To(Equal(expectedRetriesCount))
			})
		})
	})

	Describe("GetComponents", func() {
		Context("with valid deployed components returned by the backend", func() {
			It("should return the deployed components returned by the backend", func() {
				client := newRetryableRestClient(200, testutil.ComponentsResult)
				result, err := client.GetComponents()
				testutil.ExpectNoErrorAndResult(err, result, &testutil.ComponentsResult)
			})
		})
		Context("with an error returned by the backend", func() {
			It("should return the error and a zero result", func() {
				client := newRetryableRestClient(404, nil)
				result, err := client.GetComponents()
				testutil.ExpectErrorAndZeroResult(err, result)
			})
		})
		Context("with an 500 server error returned by the backend", func() {
			It("should return the error, zero result and retry expectedRetriesCount times", func() {
				client := newRetryableRestClient(500, nil)
				retryableRestClient := client.(restclient.RetryableRestClient)
				result, err := retryableRestClient.GetComponents()
				testutil.ExpectErrorAndZeroResult(err, result)
				mockRestClient := retryableRestClient.RestClient.(*MockRestClient)
				Expect(mockRestClient.GetRetriesCount()).To(Equal(expectedRetriesCount))
			})
		})
		Context("with an 502 server error returned by the backend", func() {
			It("should return the error, zero result and retry expectedRetriesCount times", func() {
				client := newRetryableRestClient(502, nil)
				retryableRestClient := client.(restclient.RetryableRestClient)
				result, err := retryableRestClient.GetComponents()
				testutil.ExpectErrorAndZeroResult(err, result)
				mockRestClient := retryableRestClient.RestClient.(*MockRestClient)
				Expect(mockRestClient.GetRetriesCount()).To(Equal(expectedRetriesCount))
			})
		})
	})

	Describe("GetMta", func() {
		Context("with a valid MTA returned by the backend", func() {
			It("should return the MTA returned by the backend", func() {
				client := newRetryableRestClient(200, testutil.MtaResult)
				result, err := client.GetMta(testutil.MtaID)
				testutil.ExpectNoErrorAndResult(err, result, &testutil.MtaResult)
			})
		})
		Context("with an error returned by the backend", func() {
			It("should return the error and a zero result", func() {
				client := newRetryableRestClient(404, nil)
				result, err := client.GetMta(testutil.MtaID)
				testutil.ExpectErrorAndZeroResult(err, result)
			})
		})
		Context("with an 500 server error returned by the backend", func() {
			It("should return the error, zero result and retry expectedRetriesCount times", func() {
				client := newRetryableRestClient(500, nil)
				retryableRestClient := client.(restclient.RetryableRestClient)
				result, err := retryableRestClient.GetMta(testutil.MtaID)
				testutil.ExpectErrorAndZeroResult(err, result)
				mockRestClient := retryableRestClient.RestClient.(*MockRestClient)
				Expect(mockRestClient.GetRetriesCount()).To(Equal(expectedRetriesCount))
			})
		})
		Context("with an 502 server error returned by the backend", func() {
			It("should return the error, zero result and retry expectedRetriesCount times", func() {
				client := newRetryableRestClient(502, nil)
				retryableRestClient := client.(restclient.RetryableRestClient)
				result, err := retryableRestClient.GetMta(testutil.MtaID)
				testutil.ExpectErrorAndZeroResult(err, result)
				mockRestClient := retryableRestClient.RestClient.(*MockRestClient)
				Expect(mockRestClient.GetRetriesCount()).To(Equal(expectedRetriesCount))
			})
		})
	})

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

// GetOperations retrieves all ongoing operations
func (c *MockRestClient) GetOperations(lastRequestedOperations *string, requestedStates []string) (models.Operations, error) {
	c.retriesCount++
	return c.restClient.GetOperations(lastRequestedOperations, requestedStates)
}

// GetComponents retrieves all deployed components (MTAs and standalone apps)
func (c *MockRestClient) GetComponents() (*models.Components, error) {
	c.retriesCount++
	return c.restClient.GetComponents()
}

// GetMta retrieves the deployed MTA with the specified MTA ID
func (c *MockRestClient) GetMta(mtaID string) (*models.Mta, error) {
	c.retriesCount++
	return c.restClient.GetMta(mtaID)
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
