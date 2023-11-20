package csrf

import (
	"net/http"
	"net/url"

	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/clients/csrf/csrf_parameters"
	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/log"
)

const CsrfTokenHeaderFetchValue = "Fetch"
const CsrfTokenHeaderRequiredValue = "Required"
const CsrfTokensApi = "/api/v1/csrf-token" //also available at /rest/csrf-token
const AuthorizationHeader = "Authorization"
const XCsrfHeader = "X-Csrf-Header"
const XCsrfToken = "X-Csrf-Token"

type DefaultCsrfTokenFetcher struct {
	transport *Transport
}

func NewDefaultCsrfTokenFetcher(transport *Transport) *DefaultCsrfTokenFetcher {
	return &DefaultCsrfTokenFetcher{transport: transport}
}

func (c *DefaultCsrfTokenFetcher) FetchCsrfToken(url, authToken string) (csrf_parameters.CsrfParams, error) {
	fetchTokenRequest, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return csrf_parameters.CsrfParams{}, err
	}
	fetchTokenRequest.Header.Set(XCsrfToken, CsrfTokenHeaderFetchValue)
	fetchTokenRequest.Header.Set(AuthorizationHeader, authToken)

	response, err := c.transport.Delegate.RoundTrip(fetchTokenRequest)
	if err != nil {
		return csrf_parameters.CsrfParams{}, err
	}

	// if there are set-cookie headers present in response - persist them in Transport
	cookies := response.Cookies()
	if len(cookies) != 0 {
		log.Tracef("Set-Cookie headers present in response, updating current with '" + prettyPrintCookies(cookies) + "'\n")
	}

	log.Tracef("New CSRF Token fetched '" + response.Header.Get(XCsrfToken) + "'\n")
	return csrf_parameters.CsrfParams{
		CsrfTokenHeader: response.Header.Get(XCsrfHeader),
		CsrfTokenValue:  response.Header.Get(XCsrfToken),
		Cookies:         cookies,
	}, nil
}

func getCsrfTokenUrl(url *url.URL) string {
	return url.Scheme + "://" + url.Host + CsrfTokensApi
}
