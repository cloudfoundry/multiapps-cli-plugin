package csrf

import (
	"net/http"
	"strconv"

	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/log"
	"github.com/jinzhu/copier"
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
	req2 := http.Request{}
	copier.Copy(&req2, req)

	UpdateCookiesIfNeeded(t.Cookies.Cookies, &req2)

	csrfTokenManager := NewDefaultCsrfTokenManager(&t, &req2)
	err := csrfTokenManager.updateToken()
	if err != nil {
		return nil, err
	}

	log.Tracef("Sending a request with CSRF '" + req2.Header.Get("X-Csrf-Header") + " : " + req2.Header.Get("X-Csrf-Token") + "'\n")
	log.Tracef("The sticky-session headers are: " + prettyPrintCookies(req2.Cookies()) + "\n")
	res, err := t.OriginalTransport.RoundTrip(&req2)
	if err != nil {
		return nil, err
	}
	tokenWasRefreshed, err := csrfTokenManager.refreshTokenIfNeeded(res)
	if err != nil {
		return nil, err
	}

	if tokenWasRefreshed {
		log.Tracef("Response code '" + strconv.Itoa(res.StatusCode) + "' from bad token. Must Retry.\n")
		return nil, &ForbiddenError{}
	}

	return res, err
}

func UpdateCookiesIfNeeded(cookies []*http.Cookie, request *http.Request) {
	if cookies != nil && len(cookies) > 0 {
		request.Header.Del(CookieHeader)
		for _, cookie := range cookies {
			request.AddCookie(cookie)
		}
	}
}

func prettyPrintCookies(cookies []*http.Cookie) string {
	result := ""
	for _, cookie := range cookies {
		result = result + cookie.String() + " "
	}
	return result
}
