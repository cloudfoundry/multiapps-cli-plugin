package csrf

import "github.com/cloudfoundry-incubator/multiapps-cli-plugin/clients/csrf/csrf_parameters"

type CsrfTokenFetcher interface {
	FetchCsrfToken(url, authToken string) (csrf_parameters.CsrfParams, error)
}
