package csrf

import "net/http"

type CsrfTokenUpdater interface {
	checkAndUpdateCsrfToken() error
	initializeToken(forceInitializing bool, url string) error
	isRetryNeeded(response *http.Response) (bool, error)
	updateCurrentCsrfToken(request *http.Request, t *Transport)
	isProtectionRequired(req *http.Request, t *Transport) bool
}
