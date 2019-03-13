package csrf

import (
	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/clients/csrf/csrf_paramters"
	"net/http"
)

type CsrfTokenFetcher interface {
	FetchCsrfToken(url string, currentRequest *http.Request) (*csrf_paramters.CsrfRequestHeader, error)
}
