package baseclient

import (
	"net/http"
	"net/http/httptest"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

// mockRoundTripper for testing
type mockRoundTripper struct {
	lastRequest *http.Request
}

func (m *mockRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	m.lastRequest = req
	return &http.Response{
		StatusCode: 200,
		Header:     make(http.Header),
		Body:       http.NoBody,
	}, nil
}

var _ = Describe("UserAgentTransport", func() {

	Describe("NewUserAgentTransport", func() {
		Context("when base transport is nil", func() {
			It("should use http.DefaultTransport as base", func() {
				transport := NewUserAgentTransport(nil)

				Expect(transport.Base).To(Equal(http.DefaultTransport))
			})

			It("should set User-Agent with correct prefix", func() {
				transport := NewUserAgentTransport(nil)

				Expect(transport.UserAgent).ToNot(BeEmpty())
				Expect(transport.UserAgent).To(HavePrefix("Multiapps-CF-plugin/"))
			})
		})

		Context("when custom base transport is provided", func() {
			It("should use the provided transport as base", func() {
				mockTransport := &mockRoundTripper{}
				transport := NewUserAgentTransport(mockTransport)

				Expect(transport.Base).To(Equal(mockTransport))
			})
		})
	})

	Describe("RoundTrip", func() {
		var mockTransport *mockRoundTripper
		var userAgentTransport *UserAgentTransport

		BeforeEach(func() {
			mockTransport = &mockRoundTripper{}
			userAgentTransport = NewUserAgentTransport(mockTransport)
		})

		Context("when making a request", func() {
			var req *http.Request

			BeforeEach(func() {
				req = httptest.NewRequest("GET", "http://example.com", nil)
				req.Header.Set("Existing-Header", "value")
			})

			It("should pass the request to base transport", func() {
				_, err := userAgentTransport.RoundTrip(req)

				Expect(err).ToNot(HaveOccurred())
				Expect(mockTransport.lastRequest).ToNot(BeNil())
			})

			It("should add User-Agent header to the request", func() {
				_, err := userAgentTransport.RoundTrip(req)

				Expect(err).ToNot(HaveOccurred())
				userAgent := mockTransport.lastRequest.Header.Get("User-Agent")
				Expect(userAgent).ToNot(BeEmpty())
				Expect(userAgent).To(HavePrefix("Multiapps-CF-plugin/"))
			})

			It("should preserve existing headers", func() {
				_, err := userAgentTransport.RoundTrip(req)

				Expect(err).ToNot(HaveOccurred())
				existingHeader := mockTransport.lastRequest.Header.Get("Existing-Header")
				Expect(existingHeader).To(Equal("value"))
			})

			It("should not modify the original request", func() {
				_, err := userAgentTransport.RoundTrip(req)

				Expect(err).ToNot(HaveOccurred())
				Expect(req.Header.Get("User-Agent")).To(BeEmpty())
			})
		})

		Context("when request has existing User-Agent header", func() {
			var req *http.Request

			BeforeEach(func() {
				req = httptest.NewRequest("GET", "http://example.com", nil)
				req.Header.Set("User-Agent", "existing-user-agent")
			})

			It("should override the existing User-Agent header", func() {
				_, err := userAgentTransport.RoundTrip(req)

				Expect(err).ToNot(HaveOccurred())
				userAgent := mockTransport.lastRequest.Header.Get("User-Agent")
				Expect(userAgent).ToNot(Equal("existing-user-agent"))
				Expect(userAgent).To(HavePrefix("Multiapps-CF-plugin/"))
			})

			It("should not modify the original request", func() {
				_, err := userAgentTransport.RoundTrip(req)

				Expect(err).ToNot(HaveOccurred())
				Expect(req.Header.Get("User-Agent")).To(Equal("existing-user-agent"))
			})
		})
	})
})
