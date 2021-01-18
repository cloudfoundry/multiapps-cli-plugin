package csrf

import (
	"net/http"
	"os"

	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/clients/csrf/csrf_parameters"
	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/log"
	"github.com/cloudfoundry/cli/plugin"
)

const CsrfTokenHeaderFetchValue = "Fetch"
const CsrfTokensApi = "/api/v1/csrf-token"
const ContentTypeHeader = "Content-Type"
const AuthorizationHeader = "Authorization"
const ApplicationJsonContentType = "application/json"
const CookieHeader = "Cookie"

type DefaultCsrfTokenFetcher struct {
	transport *Transport
}

func NewDefaultCsrfTokenFetcher(transport *Transport) *DefaultCsrfTokenFetcher {
	return &DefaultCsrfTokenFetcher{transport: transport}
}

func (c *DefaultCsrfTokenFetcher) FetchCsrfToken(url string, currentRequest *http.Request) (*csrf_parameters.CsrfRequestHeader, error) {

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

	response, err := c.transport.OriginalTransport.RoundTrip(fetchTokenRequest)
	if err != nil {
		return nil, err
	}
	// if there are set-cookie headers present in response - persist them in Transport
	if len(response.Cookies()) != 0 {
		log.Tracef("Set-Cookie headers present in response, updating current with '" + prettyPrintCookies(response.Cookies()) + "'\n")

		c.transport.Cookies.Cookies = response.Cookies()
	}

	log.Tracef("New CSRF Token fetched '" + response.Header.Get(XCsrfToken) + "'\n")
	return &csrf_parameters.CsrfRequestHeader{response.Header.Get(XCsrfHeader), response.Header.Get(XCsrfToken)}, nil
}

func getCsrfTokenUrl(req *http.Request) string {
	return string(req.URL.Scheme) + "://" + string(req.URL.Host) + CsrfTokensApi
}
