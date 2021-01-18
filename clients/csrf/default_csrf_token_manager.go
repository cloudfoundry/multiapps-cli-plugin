package csrf

import (
	"net/http"
)

const XCsrfHeader = "X-Csrf-Header"
const XCsrfToken = "X-Csrf-Token"
const CsrfTokenHeaderRequiredValue = "Required"

type DefaultCsrfTokenManager struct {
	request          *http.Request
	transport        *Transport
	csrfTokenFetcher CsrfTokenFetcher
}

func NewDefaultCsrfTokenManager(transport *Transport, request *http.Request) *DefaultCsrfTokenManager {
	return &DefaultCsrfTokenManager{request: request, transport: transport, csrfTokenFetcher: NewDefaultCsrfTokenFetcher(transport)}
}

func NewDefaultCsrfTokenManagerWithFetcher(transport *Transport, request *http.Request, csrfTokenFetcher CsrfTokenFetcher) *DefaultCsrfTokenManager {
	return &DefaultCsrfTokenManager{request: request, transport: transport, csrfTokenFetcher: csrfTokenFetcher}
}

func (c *DefaultCsrfTokenManager) updateToken() error {
	if !c.shouldInitialize() {
		return nil
	}
	err := c.initializeToken(false)
	if err != nil {
		return err
	}

	c.updateTokenInRequest()

	return nil
}

func (c *DefaultCsrfTokenManager) initializeToken(forceInitializing bool) error {
	if forceInitializing || !c.transport.Csrf.IsInitialized {
		var err error
		csrfToken, err := c.csrfTokenFetcher.FetchCsrfToken(getCsrfTokenUrl(c.request), c.request)
		if csrfToken == nil {
			return nil
		}
		c.transport.Csrf.Header, c.transport.Csrf.Token = csrfToken.CsrfTokenHeader, csrfToken.CsrfTokenValue
		if err != nil {
			return err
		}
		c.transport.Csrf.IsInitialized = true
	}

	return nil
}

func (c *DefaultCsrfTokenManager) refreshTokenIfNeeded(response *http.Response) (bool, error) {
	if !c.isProtectionRequired(c.request, c.transport) {
		return false, nil
	}
	if c.transport.Csrf.IsInitialized && (response.StatusCode == http.StatusForbidden) {
		csrfToken := response.Header.Get(XCsrfToken)

		if CsrfTokenHeaderRequiredValue == csrfToken {
			err := c.initializeToken(true)
			if err != nil {
				return false, err
			}
			// the token was refreshed successfully so the client should retry the request
			return c.transport.Csrf.Token != "", nil
		}
	}

	return false, nil
}

func (c *DefaultCsrfTokenManager) updateTokenInRequest() {
	if c.transport.Csrf.Token != "" && c.transport.Csrf.Header != "" {
		c.request.Header.Set(XCsrfToken, c.transport.Csrf.Token)
		c.request.Header.Set(XCsrfHeader, c.transport.Csrf.Header)
	}
}

func (c *DefaultCsrfTokenManager) isProtectionRequired(req *http.Request, t *Transport) bool {
	return !t.Csrf.NonProtectedMethods[req.Method]
}

func (c *DefaultCsrfTokenManager) shouldInitialize() bool {
	return c.request != nil && c.isProtectionRequired(c.request, c.transport)
}
