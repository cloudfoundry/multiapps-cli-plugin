package csrf

import (
	"github.com/cloudfoundry/cli/plugin"
	"github.com/jinzhu/copier"
	"net/http"
	"os"
)

type Csrf struct {
	Header        string
	Token         string
	IsInitialized bool
	cookies       []*http.Cookie
}

type Transport struct {
	Transport           http.RoundTripper
	Csrf                *Csrf
	NonProtectedMethods map[string]bool
}

const xCsrfHeader = "X-Csrf-Header"
const xCsrfToken = "X-Csrf-Token"
const csrfTokenHeaderFetchValue = "Fetch"
const csrfTokenHeaderRequiredValue = "Required"
const cookieString = "Cookie"

func (t Transport) RoundTrip(req *http.Request) (*http.Response, error) {

	req2 := http.Request{}
	copier.Copy(&req2, req)

	if t.Csrf.cookies != nil {
		req2.Header.Del(cookieString)
		for _, cookie := range t.Csrf.cookies {
			req2.AddCookie(cookie)
		}
	}

	err := setCsrfToken(&req2, &t)

	if err != nil {
		return nil, err
	}

	res, err := t.Transport.RoundTrip(&req2)
	if err != nil {
		return nil, err
	}

	isRetryNeeded, err := isRetryNeeded(&req2, res, &t)

	if err != nil {
		return nil, err
	}

	if isRetryNeeded {
		return res, &ForbiddenError{}
	}

	return res, err
}

func setCsrfToken(request *http.Request, t *Transport) error {
	if request == nil || !isProtectionRequired(request, t) {
		return nil
	}

	err := initializeToken(false, getFetchNewTokenUrl(request), t, request)
	if err != nil {
		return err
	}

	updateCurrentCsrfTokens(request, t)

	return nil
}

func updateCurrentCsrfTokens(request *http.Request, t *Transport) {
	if t.Csrf.Token != "" {
		request.Header.Set(xCsrfToken, t.Csrf.Token)
	}
	if t.Csrf.Header != "" {
		request.Header.Set(xCsrfHeader, t.Csrf.Header)
	}
}

func initializeToken(force bool, url string, t *Transport, currentRequest *http.Request) error {
	if force || !t.Csrf.IsInitialized {
		var err error
		t.Csrf.Header, t.Csrf.Token, err = fetchNewCsrfToken(url, t, currentRequest)
		if err != nil {
			return err
		}
		t.Csrf.IsInitialized = true
	}

	return nil
}

func isRetryNeeded(request *http.Request, response *http.Response, t *Transport) (bool, error) {
	if !isProtectionRequired(request, t) {
		return false, nil
	}

	if t.Csrf.IsInitialized && (response.StatusCode == http.StatusForbidden) {
		csrfToken := response.Header.Get(xCsrfToken)

		if csrfTokenHeaderRequiredValue == csrfToken {
			err := initializeToken(true, getFetchNewTokenUrl(request), t, request)
			if err != nil {
				return false, err
			}

			return t.Csrf.Token != "", nil
		}
	}

	return false, nil
}

func isProtectionRequired(req *http.Request, t *Transport) bool {
	return !t.NonProtectedMethods[req.Method]
}

func fetchNewCsrfToken(url string, t *Transport, currentRequest *http.Request) (string, string, error) {
	fetchTokenRequest, _ := http.NewRequest(http.MethodGet, url, nil)
	fetchTokenRequest.Header.Set(xCsrfToken, csrfTokenHeaderFetchValue)

	fetchTokenRequest.Header.Set("Content-Type", "application/json")

	cliConnection := plugin.NewCliConnection(os.Args[1])

	token, err := cliConnection.AccessToken()
	if err != nil {
		return "", "", err
	}

	fetchTokenRequest.Header.Set("Authorization", token)

	for _, cookie := range currentRequest.Cookies() {
		fetchTokenRequest.AddCookie(cookie)
	}

	response, err := t.Transport.RoundTrip(fetchTokenRequest)
	if err != nil {
		return "", "", err
	}
	if len(response.Cookies()) != 0 {
		fetchTokenRequest.Header.Del(cookieString)
		for _, cookie := range response.Cookies() {
			fetchTokenRequest.AddCookie(cookie)
		}

		t.Csrf.cookies = fetchTokenRequest.Cookies()

		response, err = t.Transport.RoundTrip(fetchTokenRequest)
	}

	if err != nil {
		return "", "", err
	}

	return response.Header.Get(xCsrfHeader), response.Header.Get(xCsrfToken), nil
}

func getFetchNewTokenUrl(req *http.Request) string {
	return string(req.URL.Scheme) + "://" +
		string(req.URL.Host) + "/api/v1/csrf-token"
}
