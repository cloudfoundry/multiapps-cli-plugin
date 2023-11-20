package csrf

import (
	"errors"
	"net/http"
	"strings"

	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/log"
)

type Transport struct {
	Delegate http.RoundTripper
	Csrf     *CsrfTokenHelper
}

func (t *Transport) RoundTrip(req *http.Request) (*http.Response, error) {
	reqCopy := req.Clone(req.Context())

	csrfTokenManager := NewDefaultCsrfTokenManager(t)
	err := csrfTokenManager.updateToken(reqCopy)
	if err != nil {
		return nil, err
	}

	log.Tracef("Sending a request with CSRF %s\n", reqCopy.Header.Get("X-Csrf-Token"))
	log.Tracef("Cookies used are: %s\n", prettyPrintCookies(reqCopy.Cookies()))

	resp, err := t.Delegate.RoundTrip(reqCopy)
	if err != nil {
		return nil, err
	}
	if isCsrfError(resp) {
		csrfTokenManager.invalidateToken()
		return nil, errors.New("invalid CSRF token")
	}

	return resp, nil
}

func prettyPrintCookies(cookies []*http.Cookie) string {
	var result strings.Builder
	for _, cookie := range cookies {
		result.WriteString(cookie.String())
		result.WriteRune(' ')
	}
	return result.String()
}

func isCsrfError(resp *http.Response) bool {
	return resp.StatusCode == http.StatusForbidden && resp.Header.Get(XCsrfToken) == CsrfTokenHeaderRequiredValue

}
