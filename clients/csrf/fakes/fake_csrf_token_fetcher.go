package fakes

import (
	"net/http"

	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/clients/csrf/csrf_parameters"
)

const FakeCsrfTokenHeader = "fake-xcsrf-token-header"
const FakeCsrfTokenValue = "fake-xcsrf-token-value"

type FakeCsrfTokenFetcher struct {
}

func (c *FakeCsrfTokenFetcher) FetchCsrfToken(string, *http.Request) (*csrf_parameters.CsrfRequestHeader, error) {
	return &csrf_parameters.CsrfRequestHeader{FakeCsrfTokenHeader, FakeCsrfTokenValue}, nil
}

func NewFakeCsrfTokenFetcher() *FakeCsrfTokenFetcher {
	return &FakeCsrfTokenFetcher{}
}
