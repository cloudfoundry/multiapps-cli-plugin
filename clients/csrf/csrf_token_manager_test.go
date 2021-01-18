package csrf

import (
	"net/http"
	"net/url"

	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/clients/csrf/fakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

const testUrl = "http://localhost:1000"

const csrfTokenNotSet = ""

var _ = Describe("DefaultCsrfTokenUpdater", func() {
	Context("", func() {
		It("protection not needed", func() {
			transport, request := createTransport(), createRequest(http.MethodGet)
			csrfTokenManager := NewDefaultCsrfTokenManager(transport, request)
			Expect(csrfTokenManager.isProtectionRequired(request, transport)).To(BeFalse())
		})
		It("protection not needed", func() {
			transport, request := createTransport(), createRequest(http.MethodOptions)
			csrfTokenManager := NewDefaultCsrfTokenManager(transport, request)
			Expect(csrfTokenManager.isProtectionRequired(request, transport)).To(BeFalse())
		})
		It("protection not needed", func() {
			transport, request := createTransport(), createRequest(http.MethodHead)
			csrfTokenManager := NewDefaultCsrfTokenManager(transport, request)
			Expect(csrfTokenManager.isProtectionRequired(request, transport)).To(BeFalse())
		})
		It("protection needed", func() {
			transport, request := createTransport(), createRequest(http.MethodPost)
			csrfTokenManager := NewDefaultCsrfTokenManager(transport, request)
			Expect(csrfTokenManager.isProtectionRequired(request, transport)).To(BeTrue())
		})
		It("retry is not needed", func() {
			transport, request := createTransport(), createRequest(http.MethodPost)
			csrfTokenManager := NewDefaultCsrfTokenManager(transport, request)
			Expect(csrfTokenManager.refreshTokenIfNeeded(createResponse(http.StatusOK, ""))).To(BeFalse())
		})
		It("retry is not needed", func() {
			transport, request := createTransport(), createRequest(http.MethodPost)
			csrfTokenManager := NewDefaultCsrfTokenManager(transport, request)
			Expect(csrfTokenManager.refreshTokenIfNeeded(createResponse(http.StatusForbidden, CsrfTokenHeaderRequiredValue))).To(BeFalse())
		})
		It("retry is needed", func() {
			transport := createTransport()
			transport.Csrf.IsInitialized = true
			request := createRequest(http.MethodPost)
			csrfTokenManager := NewDefaultCsrfTokenManagerWithFetcher(transport, request, fakes.NewFakeCsrfTokenFetcher())
			isRetryNeeded, err := csrfTokenManager.refreshTokenIfNeeded(createResponse(http.StatusForbidden, CsrfTokenHeaderRequiredValue))
			Ω(err).ShouldNot(HaveOccurred())
			Expect(isRetryNeeded).To(BeTrue())
		})
		It("initialize new token", func() {
			transport := createTransport()
			transport.Csrf.IsInitialized = true
			request := createRequest(http.MethodPost)
			csrfTokenManager := NewDefaultCsrfTokenManagerWithFetcher(transport, request, fakes.NewFakeCsrfTokenFetcher())
			err := csrfTokenManager.initializeToken(true)
			Ω(err).ShouldNot(HaveOccurred())
			Expect(transport.Csrf.Header).To(Equal(fakes.FakeCsrfTokenHeader))
			Expect(transport.Csrf.Token).To(Equal(fakes.FakeCsrfTokenValue))
			Expect(transport.Csrf.IsInitialized).To(BeTrue())
		})
		It("update current csrf tokens", func() {
			transport := createTransport()
			request := createRequest(http.MethodGet)
			csrfTokenManager := NewDefaultCsrfTokenManagerWithFetcher(transport, request, fakes.NewFakeCsrfTokenFetcher())
			err := csrfTokenManager.initializeToken(true)
			Ω(err).ShouldNot(HaveOccurred())
			csrfTokenManager.updateTokenInRequest()
			expectCsrfTokenIsProperlySet(request, fakes.FakeCsrfTokenHeader, fakes.FakeCsrfTokenValue)
		})
		It("should not update csrf tokens", func() {
			transport, request := createTransport(), createRequest(http.MethodGet)
			csrfTokenManager := NewDefaultCsrfTokenManagerWithFetcher(transport, request, fakes.NewFakeCsrfTokenFetcher())
			err := csrfTokenManager.updateToken()
			Ω(err).ShouldNot(HaveOccurred())
			expectCsrfTokenIsProperlySet(request, csrfTokenNotSet, csrfTokenNotSet)
		})
		It("should not update csrf tokens", func() {
			transport, request := createTransport(), createRequest(http.MethodPost)
			transport.Csrf.IsInitialized = true
			csrfTokenManager := NewDefaultCsrfTokenManagerWithFetcher(transport, request, fakes.NewFakeCsrfTokenFetcher())
			err := csrfTokenManager.updateToken()
			Ω(err).ShouldNot(HaveOccurred())
			expectCsrfTokenIsProperlySet(request, csrfTokenNotSet, csrfTokenNotSet)
		})
		It("should not update csrf tokens", func() {
			transport, request := createTransport(), createRequest(http.MethodGet)
			csrfTokenManager := NewDefaultCsrfTokenManagerWithFetcher(transport, request, fakes.NewFakeCsrfTokenFetcher())
			err := csrfTokenManager.updateToken()
			Ω(err).ShouldNot(HaveOccurred())
			expectCsrfTokenIsProperlySet(request, csrfTokenNotSet, csrfTokenNotSet)
		})
		It("should update csrf tokens", func() {
			transport, request := createTransport(), createRequest(http.MethodPost)
			csrfTokenManager := NewDefaultCsrfTokenManagerWithFetcher(transport, request, fakes.NewFakeCsrfTokenFetcher())
			err := csrfTokenManager.updateToken()
			Ω(err).ShouldNot(HaveOccurred())
			expectCsrfTokenIsProperlySet(request, fakes.FakeCsrfTokenHeader, fakes.FakeCsrfTokenValue)
		})
		Context("set cookies in the request, valid cookies", func() {
			It("should be equal", func() {
				request := createRequest(http.MethodGet)
				cookies := createValidCookies()
				UpdateCookiesIfNeeded(cookies, request)
				Expect(cookies).To(Equal(request.Cookies()))
			})
		})
	})
})

func expectCsrfTokenIsProperlySet(request *http.Request, csrfTokenHeader, csrfTokenValue string) {
	Expect(request.Header.Get(XCsrfHeader)).To(Equal(csrfTokenHeader))
	Expect(request.Header.Get(XCsrfToken)).To(Equal(csrfTokenValue))
}

func createResponse(httpStatusCode int, csrfToken string) *http.Response {
	response := &http.Response{}
	response.Header = make(http.Header)
	response.StatusCode = httpStatusCode
	response.Header.Set(XCsrfToken, csrfToken)

	return response
}

func createTransport() *Transport {
	return &Transport{http.DefaultTransport.(*http.Transport),
		&Csrf{"", "", false, getNonProtectedMethods()}, &Cookies{[]*http.Cookie{}}}
}

func getNonProtectedMethods() map[string]bool {
	nonProtectedMethods := make(map[string]bool)

	nonProtectedMethods[http.MethodGet] = true
	nonProtectedMethods[http.MethodHead] = true
	nonProtectedMethods[http.MethodOptions] = true

	return nonProtectedMethods
}

func createValidCookies() []*http.Cookie {
	var cookies []*http.Cookie
	cookie1 := &http.Cookie{}
	cookie1.Name = "JSESSION"
	cookie1.Value = "123"
	cookie2 := &http.Cookie{}
	cookie2.Name = "__V_CAP__"
	cookie2.Value = "321"
	cookies = append(cookies, cookie1)
	cookies = append(cookies, cookie2)

	return cookies
}

func createRequest(method string) *http.Request {
	request := &http.Request{}
	requestUrl := &url.URL{}
	requestUrl.Scheme = "http"
	requestUrl.Host = "localhost:1000"
	request.URL = requestUrl
	request.Header = make(http.Header)
	request.Method = method

	return request
}
