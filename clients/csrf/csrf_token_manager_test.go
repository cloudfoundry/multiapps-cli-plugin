package csrf

import (
	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/clients/csrf/fakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"net/http"
	"time"
)

const csrfTokenNotSet = ""

var _ = Describe("DefaultCsrfTokenUpdater", func() {
	Context("", func() {
		It("protection not needed", func() {
			transport, request := createTransport(), createRequest(http.MethodGet)
			Expect(transport.Csrf.IsProtectionRequired(request)).To(BeFalse())
		})
		It("protection not needed", func() {
			transport, request := createTransport(), createRequest(http.MethodOptions)
			Expect(transport.Csrf.IsProtectionRequired(request)).To(BeFalse())
		})
		It("protection not needed", func() {
			transport, request := createTransport(), createRequest(http.MethodHead)
			Expect(transport.Csrf.IsProtectionRequired(request)).To(BeFalse())
		})
		It("protection needed", func() {
			transport, request := createTransport(), createRequest(http.MethodPost)
			Expect(transport.Csrf.IsProtectionRequired(request)).To(BeTrue())
		})
		It("initialize new token", func() {
			transport := createTransport()
			request := createRequest(http.MethodPost)
			csrfTokenManager := NewDefaultCsrfTokenManagerWithFetcher(transport.Csrf, fakes.NewFakeCsrfTokenFetcher())
			err := csrfTokenManager.updateToken(request)
			Ω(err).ShouldNot(HaveOccurred())
			Expect(transport.Csrf.Header).To(Equal(fakes.FakeCsrfTokenHeader))
			Expect(transport.Csrf.Token).To(Equal(fakes.FakeCsrfTokenValue))
			Expect(transport.Csrf.LastUpdateTime).ToNot(Equal(time.Time{}))
		})
		It("should not update csrf tokens", func() {
			transport, request := createTransport(), createRequest(http.MethodGet)
			csrfTokenManager := NewDefaultCsrfTokenManagerWithFetcher(transport.Csrf, fakes.NewFakeCsrfTokenFetcher())
			err := csrfTokenManager.updateToken(request)
			Ω(err).ShouldNot(HaveOccurred())
			expectCsrfTokenIsProperlySet(request, csrfTokenNotSet, csrfTokenNotSet)
		})
		It("should not update csrf tokens", func() {
			transport, request := createTransport(), createRequest(http.MethodPost)
			csrfTokenManager := NewDefaultCsrfTokenManagerWithFetcher(transport.Csrf, fakes.NewFakeCsrfTokenFetcher())
			err := csrfTokenManager.updateToken(request)
			Ω(err).ShouldNot(HaveOccurred())
			expectCsrfTokenIsProperlySet(request, csrfTokenNotSet, csrfTokenNotSet)
		})
		It("should not update csrf tokens", func() {
			transport, request := createTransport(), createRequest(http.MethodGet)
			csrfTokenManager := NewDefaultCsrfTokenManagerWithFetcher(transport.Csrf, fakes.NewFakeCsrfTokenFetcher())
			err := csrfTokenManager.updateToken(request)
			Ω(err).ShouldNot(HaveOccurred())
			expectCsrfTokenIsProperlySet(request, csrfTokenNotSet, csrfTokenNotSet)
		})
		It("should update csrf tokens", func() {
			transport, request := createTransport(), createRequest(http.MethodPost)
			csrfTokenManager := NewDefaultCsrfTokenManagerWithFetcher(transport.Csrf, fakes.NewFakeCsrfTokenFetcher())
			err := csrfTokenManager.updateToken(request)
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
			It("should not get unauthorized", func() {
				transport, request := createTransport(), createRequest(http.MethodPost)
				csrfTokenManager := NewDefaultCsrfTokenManagerWithFetcher(transport.Csrf, fakes.NewFakeCsrfTokenFetcher())
				err := csrfTokenManager.updateToken(request)
				Ω(err).ShouldNot(HaveOccurred())
				Expect(request.Header.Get(fakes.FakeCsrfTokenHeader)).To(Equal(transport.Csrf.Token))

			})
		})
	})
})

func expectCsrfTokenIsProperlySet(request *http.Request, csrfTokenHeader, csrfTokenValue string) {
	Expect(request.Header.Get(csrfTokenHeader)).To(Equal(csrfTokenValue))
}

func createTransport() *Transport {
	return &Transport{Delegate: http.DefaultTransport.(*http.Transport),
		Csrf: &CsrfTokenHelper{NonProtectedMethods: getNonProtectedMethods()}}
}

func getNonProtectedMethods() map[string]struct{} {
	nonProtectedMethods := make(map[string]struct{}, 4)

	nonProtectedMethods[http.MethodGet] = struct{}{}
	nonProtectedMethods[http.MethodHead] = struct{}{}
	nonProtectedMethods[http.MethodTrace] = struct{}{}
	nonProtectedMethods[http.MethodOptions] = struct{}{}

	return nonProtectedMethods
}

func createValidCookies() []*http.Cookie {
	return []*http.Cookie{
		{
			Name:  "JSESSION",
			Value: "123",
		},
		{
			Name:  "__V_CAP__",
			Value: "321",
		},
	}
}

func createRequest(method string) *http.Request {
	r, _ := http.NewRequest(method, "http://localhost:1000", nil)
	return r
}
