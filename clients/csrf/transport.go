package csrf

import (
	"github.com/jinzhu/copier"
	"net/http"
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
	Transport http.RoundTripper
	Csrf      *Csrf
	Cookies   *Cookies
}

func (t Transport) RoundTrip(req *http.Request) (*http.Response, error) {
	req2 := http.Request{}
	copier.Copy(&req2, req)

	if t.Cookies != nil {
		UpdateCookiesIfNeeded(t.Cookies.Cookies, &req2)
	}

	csrfTokenManager := NewDefaultCsrfTokenUpdater(&t, &req2, NewDefaultCsrfTokenFetcher(&t))

	err := csrfTokenManager.updateCsrfToken()
	if err != nil {
		return nil, err
	}

	res, err := t.Transport.RoundTrip(&req2)
	if err != nil {
		return nil, err
	}
	isRetryNeeded, err := csrfTokenManager.isRetryNeeded(res)
	if err != nil {
		return nil, err
	}

	if isRetryNeeded {
		return nil, &ForbiddenError{}
	}

	return res, err
}
