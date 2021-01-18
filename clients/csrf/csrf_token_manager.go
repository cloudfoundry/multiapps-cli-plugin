package csrf

import "net/http"

type CsrfTokenManager interface {
	initializeToken(forceInitializing bool) error
	updateToken()
	refreshTokenIfNeeded(response *http.Response) (bool, error)
	updateCsrfTokenInRequest()
}
