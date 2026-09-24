package mtaclient_test

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"time"

	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/clients/baseclient"
	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/clients/csrf"
	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/clients/models"
	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/clients/mtaclient"
	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/testutil"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("MtaRestClient", func() {
	Describe("StartMtaOperation", func() {
		Context("when the backend rate-limits the request with HTTP 429", func() {
			It("should return a RetryAfterError carrying the wait time so the caller can stop and report it", func() {
				client := newMtaClient(http.StatusTooManyRequests)

				_, err := client.StartMtaOperation(models.Operation{})

				Expect(err).Should(HaveOccurred())
				var retryErr *baseclient.RetryAfterError
				Expect(errors.As(err, &retryErr)).Should(BeTrue())
			})

			It("should carry the server's Retry-After header value", func() {
				client := newMtaClientWithHeaders(http.StatusTooManyRequests, http.Header{"Retry-After": []string{"7"}})

				_, err := client.StartMtaOperation(models.Operation{})

				Expect(err).Should(HaveOccurred())
				var retryErr *baseclient.RetryAfterError
				Expect(errors.As(err, &retryErr)).Should(BeTrue())
				Expect(retryErr.Duration).Should(Equal(7 * time.Second))
			})

			It("should fall back to the default when Retry-After is absent", func() {
				client := newMtaClientWithHeaders(http.StatusTooManyRequests, http.Header{})

				_, err := client.StartMtaOperation(models.Operation{})

				Expect(err).Should(HaveOccurred())
				var retryErr *baseclient.RetryAfterError
				Expect(errors.As(err, &retryErr)).Should(BeTrue())
				Expect(retryErr.Duration).Should(Equal(3 * time.Second))
			})
		})

		Context("when the backend returns a non-429 error", func() {
			It("should return a ClientError, not a RetryAfterError", func() {
				client := newMtaClient(http.StatusInternalServerError)

				_, err := client.StartMtaOperation(models.Operation{})

				Expect(err).Should(HaveOccurred())
				var retryErr *baseclient.RetryAfterError
				Expect(errors.As(err, &retryErr)).Should(BeFalse())
				var clientErr *baseclient.ClientError
				Expect(errors.As(err, &clientErr)).Should(BeTrue())
			})
		})
	})
})

func newMtaClient(statusCode int) mtaclient.MtaClientOperations {
	tokenFactory := testutil.NewCustomTokenFactory("test-token")
	roundTripper := testutil.NewCustomTransport(statusCode)
	return mtaclient.NewMtaClient("http://localhost:1000", "test-space-guid", roundTripper, tokenFactory)
}

func newMtaClientWithHeaders(statusCode int, headers http.Header) mtaclient.MtaClientOperations {
	tokenFactory := testutil.NewCustomTokenFactory("test-token")
	roundTripper := newHeaderTransport(statusCode, headers)
	return mtaclient.NewMtaClient("http://localhost:1000", "test-space-guid", roundTripper, tokenFactory)
}

type headerRoundTripperFunc func(*http.Request) (*http.Response, error)

func (fn headerRoundTripperFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return fn(req)
}

func newHeaderTransport(statusCode int, headers http.Header) *csrf.Transport {
	transport := headerRoundTripperFunc(func(req *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: statusCode,
			Header:     headers,
			Body:       io.NopCloser(bytes.NewBuffer(nil)),
		}, nil
	})
	userAgentTransport := baseclient.NewUserAgentTransport(transport)
	return &csrf.Transport{Delegate: userAgentTransport, Csrf: &csrf.CsrfTokenHelper{}}
}
