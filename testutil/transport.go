package testutil

import (
	"bytes"
	"io/ioutil"
	"net/http"

	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/clients/csrf"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"
)

type roundTripperFunc func(*http.Request) (*http.Response, error)

func (fn roundTripperFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return fn(req)
}

// NewCustomTransport creates a new custom transport to be used for testing
func NewCustomTransport(statusCode int) *csrf.Transport {
	csrfx := csrf.Csrf{Header: "", Token: ""}
	transport := roundTripperFunc(func(req *http.Request) (*http.Response, error) {
		var resp http.Response
		resp.StatusCode = statusCode
		resp.Header = make(http.Header)
		buf := bytes.NewBuffer(nil)
		resp.Body = ioutil.NopCloser(buf)
		return &resp, nil
	})
	return &csrf.Transport{OriginalTransport: transport, Csrf: &csrfx}
}

// NewCustomBearerToken creates a new bearer token to be used for testing
func NewCustomBearerToken(token string) runtime.ClientAuthInfoWriter {
	return runtime.ClientAuthInfoWriterFunc(func(r runtime.ClientRequest, _ strfmt.Registry) error {
		r.SetHeaderParam("Authorization", "Bearer "+token)
		return nil
	})
}
