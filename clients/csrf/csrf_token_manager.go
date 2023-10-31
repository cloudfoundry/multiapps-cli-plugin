package csrf

import "net/http"

type CsrfTokenManager interface {
	updateToken(*http.Request) error
	invalidateToken()
}
