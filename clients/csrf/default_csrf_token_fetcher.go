package csrf

import (
	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/clients/csrf/csrf_paramters"
	"github.com/cloudfoundry/cli/plugin"
	"net/http"
	"os"
)

const CsrfTokenHeaderFetchValue = "Fetch"
const CsrfTokensApi = "/api/v1/csrf-token"
const ContentTypeHeader = "Content-Type"
const AuthorizationHeader = "Authorization"
const ApplicationJsonContentType = "application/json"
const CookieHeader = "CookieHeader"

type DefaultCsrfTokenFetcher struct {
	transport *Transport
}

func NewDefaultCsrfTokenFetcher(transport *Transport) *DefaultCsrfTokenFetcher {
	return &DefaultCsrfTokenFetcher{transport: transport}
}

func (c *DefaultCsrfTokenFetcher) FetchCsrfToken(url string, currentRequest *http.Request) (*csrf_paramters.CsrfRequestHeader, error) {
	fetchTokenRequest, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	fetchTokenRequest.Header.Set(XCsrfToken, CsrfTokenHeaderFetchValue)
	fetchTokenRequest.Header.Set(ContentTypeHeader, ApplicationJsonContentType)

	cliConnection := plugin.NewCliConnection(os.Args[1])
	token, err := cliConnection.AccessToken()
	if err != nil {
		return nil, err
	}
	fetchTokenRequest.Header.Set(AuthorizationHeader, token)
	UpdateCookiesIfNeeded(currentRequest.Cookies(), fetchTokenRequest)

	response, err := c.transport.Transport.RoundTrip(fetchTokenRequest)
	if err != nil {
		return nil, err
	}
	if len(response.Cookies()) != 0 {
		fetchTokenRequest.Header.Del(CookieHeader)
		UpdateCookiesIfNeeded(response.Cookies(), fetchTokenRequest)

		c.transport.Cookies.Cookies = fetchTokenRequest.Cookies()

		response, err = c.transport.Transport.RoundTrip(fetchTokenRequest)

		if err != nil {
			return nil, err
		}
	}

	return &csrf_paramters.CsrfRequestHeader{response.Header.Get(XCsrfHeader), response.Header.Get(XCsrfToken)}, nil
}

func getCsrfTokenUrl(req *http.Request) string {
	return req.URL.Scheme + "://" + req.URL.Host + CsrfTokensApi
}

func UpdateCookiesIfNeeded(cookies []*http.Cookie, request *http.Request) {
	if cookies != nil {
		request.Header.Del(CookieHeader)
		for _, cookie := range cookies {
			request.AddCookie(cookie)
		}
	}
}
