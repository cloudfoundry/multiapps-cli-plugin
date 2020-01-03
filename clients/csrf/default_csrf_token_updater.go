package csrf

import (
	"net/http"
)

const XCsrfHeader = "X-Csrf-Header"
const XCsrfToken = "X-Csrf-Token"
const CsrfTokenHeaderRequiredValue = "Required"

type DefaultCsrfTokenUpdater struct {
	request          *http.Request
	transport        *Transport
	csrfTokenFetcher CsrfTokenFetcher
}

func NewDefaultCsrfTokenUpdater(transport *Transport, request *http.Request, csrfTokenFetcher CsrfTokenFetcher) *DefaultCsrfTokenUpdater {
	return &DefaultCsrfTokenUpdater{request: request, transport: transport, csrfTokenFetcher: csrfTokenFetcher}
}

func (c *DefaultCsrfTokenUpdater) updateCsrfToken() error {
	if !c.shouldInitialize() {
		return nil
	}
	err := c.initializeToken(false, getCsrfTokenUrl(c.request))
	if err != nil {
		return err
	}

	c.updateCurrentCsrfToken(c.request, c.transport)

	return nil
}

func (c *DefaultCsrfTokenUpdater) initializeToken(forceInitializing bool, url string) error {
	if forceInitializing || !c.transport.Csrf.IsInitialized {
		var err error
		csrfToken, err := c.csrfTokenFetcher.FetchCsrfToken(url, c.request)
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

func (c *DefaultCsrfTokenUpdater) isRetryNeeded(response *http.Response) (bool, error) {
	if !c.isProtectionRequired(c.request, c.transport) {
		return false, nil
	}
	if c.transport.Csrf.IsInitialized && (response.StatusCode == http.StatusForbidden) {
		csrfToken := response.Header.Get(XCsrfToken)

		if CsrfTokenHeaderRequiredValue == csrfToken {
			err := c.initializeToken(true, getCsrfTokenUrl(c.request))
			if err != nil {
				return false, err
			}

			return c.transport.Csrf.Token != "", nil
		}
	}

	return false, nil
}

func (c *DefaultCsrfTokenUpdater) updateCurrentCsrfToken(request *http.Request, t *Transport) {
	if c.transport.Csrf.Token != "" && c.transport.Csrf.Header != "" {
		request.Header.Set(XCsrfToken, t.Csrf.Token)
		request.Header.Set(XCsrfHeader, t.Csrf.Header)
	}
}

func (c *DefaultCsrfTokenUpdater) isProtectionRequired(req *http.Request, t *Transport) bool {
	return !t.Csrf.NonProtectedMethods[req.Method]
}

func (c *DefaultCsrfTokenUpdater) shouldInitialize() bool {
	return c.request != nil && c.isProtectionRequired(c.request, c.transport)
}
