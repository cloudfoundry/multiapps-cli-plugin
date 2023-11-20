package csrf

import (
	"net/http"
	"sync"
	"time"

	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/clients/csrf/csrf_parameters"
)

const CookieHeader = "Cookie"

type CsrfTokenHelper struct {
	Header              string
	Token               string
	Cookies             []*http.Cookie
	NonProtectedMethods map[string]struct{}
	LastUpdateTime      time.Time
	Mutex               sync.Mutex
}

func (c *CsrfTokenHelper) IsProtectionRequired(req *http.Request) bool {
	_, present := c.NonProtectedMethods[req.Method]
	return !present
}

func (c *CsrfTokenHelper) IsExpired(timeout time.Duration) bool {
	return c.LastUpdateTime.IsZero() || c.LastUpdateTime.Add(timeout).Before(time.Now())
}

func (c *CsrfTokenHelper) Update(params csrf_parameters.CsrfParams) {
	c.Header = params.CsrfTokenHeader
	c.Token = params.CsrfTokenValue
	c.Cookies = params.Cookies
	c.LastUpdateTime = time.Now()
}

func (c *CsrfTokenHelper) SetInRequest(req *http.Request) {
	req.Header.Set(c.Header, c.Token)
	UpdateCookiesIfNeeded(c.Cookies, req)
}

func (c *CsrfTokenHelper) InvalidateToken() {
	c.LastUpdateTime = time.Time{}
}

func UpdateCookiesIfNeeded(cookies []*http.Cookie, request *http.Request) {
	if len(cookies) == 0 {
		return
	}
	request.Header.Del(CookieHeader)
	for _, cookie := range cookies {
		request.AddCookie(cookie)
	}
}
