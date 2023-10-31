package fakes

import "github.com/cloudfoundry-incubator/multiapps-cli-plugin/clients/csrf/csrf_parameters"

const FakeCsrfTokenHeader = "fake-xcsrf-token-header"
const FakeCsrfTokenValue = "fake-xcsrf-token-value"

type FakeCsrfTokenFetcher struct{}

func (c *FakeCsrfTokenFetcher) FetchCsrfToken(string, string) (csrf_parameters.CsrfParams, error) {
	return csrf_parameters.CsrfParams{CsrfTokenHeader: FakeCsrfTokenHeader, CsrfTokenValue: FakeCsrfTokenValue}, nil
}

func NewFakeCsrfTokenFetcher() *FakeCsrfTokenFetcher {
	return &FakeCsrfTokenFetcher{}
}
