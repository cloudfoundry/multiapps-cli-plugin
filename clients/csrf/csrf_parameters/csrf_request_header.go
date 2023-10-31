package csrf_parameters

import "net/http"

type CsrfParams struct {
	CsrfTokenHeader string
	CsrfTokenValue  string
	Cookies         []*http.Cookie
}
