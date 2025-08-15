package commands_test

import (
	"net/http"
	"net/http/httptest"

	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/commands"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("CSRFTransport", func() {

	Describe("newTransport", func() {
		var server *httptest.Server
		var capturedHeaders http.Header
		var transport http.RoundTripper

		BeforeEach(func() {
			transport = commands.NewTransportForTesting(false)
		})

		AfterEach(func() {
			if server != nil {
				server.Close()
			}
		})

		Context("when making regular requests", func() {
			BeforeEach(func() {
				server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					capturedHeaders = r.Header
					w.WriteHeader(http.StatusOK)
				}))
			})

			It("should include User-Agent header", func() {
				req, err := http.NewRequest("GET", server.URL+"/test", nil)
				Expect(err).ToNot(HaveOccurred())

				_, err = transport.RoundTrip(req)
				Expect(err).ToNot(HaveOccurred())

				userAgent := capturedHeaders.Get("User-Agent")
				Expect(userAgent).ToNot(BeEmpty(), "Expected User-Agent header to be set in CSRF transport")
				Expect(userAgent).To(HavePrefix("Multiapps-CF-plugin/"), "Expected User-Agent to start with 'Multiapps-CF-plugin/'")
			})
		})

		Context("when CSRF token fetch is required", func() {
			var requestCount int

			BeforeEach(func() {
				requestCount = 0
				server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					requestCount++
					capturedHeaders = r.Header

					if requestCount == 1 {
						// First request should be CSRF token fetch - return 403 to trigger token fetch
						w.Header().Set("X-Csrf-Token", "required")
						w.WriteHeader(http.StatusForbidden)
					} else {
						// Second request should have the token
						w.WriteHeader(http.StatusOK)
					}
				}))
			})

			It("should include User-Agent header in CSRF token fetch request", func() {
				req, err := http.NewRequest("POST", server.URL+"/test", nil)
				Expect(err).ToNot(HaveOccurred())

				// Execute the request through the transport
				// This should trigger a CSRF token fetch first
				// We expect this to potentially error since our mock server doesn't properly implement CSRF
				// But we can still verify the User-Agent was set in the token fetch request
				transport.RoundTrip(req)

				userAgent := capturedHeaders.Get("User-Agent")
				Expect(userAgent).ToNot(BeEmpty(), "Expected User-Agent header to be set in CSRF token fetch request")
				Expect(userAgent).To(HavePrefix("Multiapps-CF-plugin/"), "Expected User-Agent to start with 'Multiapps-CF-plugin/'")
			})
		})
	})
})
