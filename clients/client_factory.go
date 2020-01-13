package clients

import (
	"net/http"

	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/clients/baseclient"
	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/clients/mtaclient"
	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/clients/restclient"
)

// ClientFactory is a factory for creating XxxClientOperations instances
type ClientFactory interface {
	NewMtaClient(host, spaceID string, rt http.RoundTripper, jar http.CookieJar, tokenFactory baseclient.TokenFactory) mtaclient.MtaClientOperations
	NewRestClient(host string, rt http.RoundTripper, jar http.CookieJar, tokenFactory baseclient.TokenFactory) restclient.RestClientOperations
}

// DefaultClientFactory a default implementation of the ClientFactory
type DefaultClientFactory struct {
	mtaClient  mtaclient.MtaClientOperations
	restClient restclient.RestClientOperations
}

// NewDefaultClientFactory a default initialization method for the factory
func NewDefaultClientFactory() *DefaultClientFactory {
	return &DefaultClientFactory{}
}

// NewMtaClient used for creating or returning cached value of the mta rest client
func (d *DefaultClientFactory) NewMtaClient(host, spaceID string, rt http.RoundTripper, jar http.CookieJar,
	tokenFactory baseclient.TokenFactory) mtaclient.MtaClientOperations {
	if d.mtaClient == nil {
		d.mtaClient = mtaclient.NewRetryableMtaRestClient(host, spaceID, rt, jar, tokenFactory)
	}
	return d.mtaClient
}

// NewRestClient used for creating or returning cached value of the rest client
func (d *DefaultClientFactory) NewRestClient(host string, rt http.RoundTripper, jar http.CookieJar,
	tokenFactory baseclient.TokenFactory) restclient.RestClientOperations {
	if d.restClient == nil {
		d.restClient = restclient.NewRetryableRestClient(host, rt, jar, tokenFactory)
	}
	return d.restClient
}
