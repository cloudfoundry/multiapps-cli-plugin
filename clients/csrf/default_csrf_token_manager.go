package csrf

import (
	"net/http"
	"time"
)

const CsrfSessionTimeout = time.Minute + 30*time.Second

type DefaultCsrfTokenManager struct {
	csrfHelper       *CsrfTokenHelper
	csrfTokenFetcher CsrfTokenFetcher
}

func NewDefaultCsrfTokenManager(transport *Transport) *DefaultCsrfTokenManager {
	return &DefaultCsrfTokenManager{csrfHelper: transport.Csrf, csrfTokenFetcher: NewDefaultCsrfTokenFetcher(transport)}
}

func NewDefaultCsrfTokenManagerWithFetcher(csrf *CsrfTokenHelper, csrfTokenFetcher CsrfTokenFetcher) *DefaultCsrfTokenManager {
	return &DefaultCsrfTokenManager{csrfHelper: csrf, csrfTokenFetcher: csrfTokenFetcher}
}

func (c *DefaultCsrfTokenManager) updateToken(req *http.Request) error {
	if !c.csrfHelper.IsProtectionRequired(req) {
		return nil
	}

	c.csrfHelper.Mutex.Lock()
	defer c.csrfHelper.Mutex.Unlock()

	if c.csrfHelper.IsExpired(CsrfSessionTimeout) {
		csrfParams, err := c.csrfTokenFetcher.FetchCsrfToken(getCsrfTokenUrl(req.URL), req.Header.Get(AuthorizationHeader))
		if err != nil {
			return err
		}
		c.csrfHelper.Update(csrfParams)
	}

	c.csrfHelper.SetInRequest(req)
	return nil
}

func (c *DefaultCsrfTokenManager) invalidateToken() {
	c.csrfHelper.Mutex.Lock()
	c.csrfHelper.InvalidateToken()
	c.csrfHelper.Mutex.Unlock()
}
