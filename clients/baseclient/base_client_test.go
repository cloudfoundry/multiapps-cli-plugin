package baseclient

import (
	"net/http"
	"net/http/httptest"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("BaseClient", func() {

	Describe("NewHTTPTransport", func() {
		Context("when creating transport with User-Agent functionality", func() {
			var server *httptest.Server
			var capturedHeaders http.Header

			BeforeEach(func() {
				server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					capturedHeaders = r.Header
					w.WriteHeader(http.StatusOK)
				}))
			})

			AfterEach(func() {
				if server != nil {
					server.Close()
				}
			})

			It("should include User-Agent header in requests", func() {
				transport := NewHTTPTransport(server.URL, "/", nil)

				req, err := http.NewRequest("GET", server.URL+"/test", nil)
				Expect(err).ToNot(HaveOccurred())

				_, err = transport.Transport.RoundTrip(req)
				Expect(err).ToNot(HaveOccurred())

				userAgent := capturedHeaders.Get("User-Agent")
				Expect(userAgent).ToNot(BeEmpty(), "Expected User-Agent header to be set")
				Expect(userAgent).To(HavePrefix("Multiapps-CF-plugin/"), "Expected User-Agent to start with 'Multiapps-CF-plugin/'")
			})
		})

		Context("when custom round tripper is provided", func() {
			It("should preserve the custom round tripper as base transport", func() {
				customTransport := &mockRoundTripper{}

				transport := NewHTTPTransport("example.com", "/", customTransport)

				userAgentTransport, ok := transport.Transport.(*UserAgentTransport)
				Expect(ok).To(BeTrue(), "Expected transport to be wrapped with UserAgentTransport")
				Expect(userAgentTransport.Base).To(Equal(customTransport), "Expected custom round tripper to be preserved as base transport")
			})
		})
	})
})
