package csrf

import (
	"net/http"

	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/clients/csrf/csrf_parameters"
)

type CsrfTokenFetcher interface {
	FetchCsrfToken(url string, currentRequest *http.Request) (*csrf_parameters.CsrfRequestHeader, error)
}
