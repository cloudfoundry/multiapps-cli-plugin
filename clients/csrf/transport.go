package csrf

import (
	"net/http"
	"strings"

	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/log"
)

type Csrf struct {
	Header              string
	Token               string
	IsInitialized       bool
	NonProtectedMethods map[string]bool
}

type Cookies struct {
	Cookies []*http.Cookie
}

type Transport struct {
	OriginalTransport http.RoundTripper
	Csrf              *Csrf
	Cookies           *Cookies
}

func (t Transport) RoundTrip(req *http.Request) (*http.Response, error) {
	reqCopy := req.Clone(req.Context())

	csrfTokenManager := NewDefaultCsrfTokenManager(&t, reqCopy)
	err := csrfTokenManager.updateToken()
	if err != nil {
		return nil, err
	}

	log.Tracef("Sending a request with CSRF '%s : %s'\n", reqCopy.Header.Get("X-Csrf-Header"), reqCopy.Header.Get("X-Csrf-Token"))
	log.Tracef("Cookies used are: %s\n", prettyPrintCookies(reqCopy.Cookies()))
	res, err := t.OriginalTransport.RoundTrip(reqCopy)
	if err != nil {
		return nil, err
	}
	tokenWasRefreshed, err := csrfTokenManager.refreshTokenIfNeeded(res)
	if err != nil {
		return nil, err
	}

	if tokenWasRefreshed {
		log.Tracef("Response code '%d' from bad token. Must Retry.\n", res.StatusCode)
		return nil, &ForbiddenError{}
	}
	return res, err
}

func prettyPrintCookies(cookies []*http.Cookie) string {
	var result strings.Builder
	for _, cookie := range cookies {
		result.WriteString(cookie.String())
		result.WriteRune(' ')
	}
	return result.String()
}
