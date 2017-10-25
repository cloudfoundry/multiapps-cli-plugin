package clients

import (
	"net/http"

	baseclient "github.com/SAP/cf-mta-plugin/clients/baseclient"
	mtaclient "github.com/SAP/cf-mta-plugin/clients/mtaclient"
	restclient "github.com/SAP/cf-mta-plugin/clients/restclient"
)

// ClientFactory is a factory for creating XxxClientOperations instances
type ClientFactory interface {
	NewMtaClient(host, spaceID string, rt http.RoundTripper, jar http.CookieJar, tokenFactory baseclient.TokenFactory) mtaclient.MtaClientOperations
	NewRestClient(host, org, space string, rt http.RoundTripper, jar http.CookieJar, tokenfactory baseclient.TokenFactory) restclient.RestClientOperations
}

// DefaultClientFactory a default implementation of the ClientFactory
type DefaultClientFactory struct {
	mtaClient  mtaclient.MtaClientOperations
	restClient restclient.RestClientOperations
}

// NewDefaultClientFactory a default intialization method for the factory
func NewDefaultClientFactory() DefaultClientFactory {
	return DefaultClientFactory{mtaClient: nil, restClient: nil}
}

// NewMtaClient used for creating or returning cached value of the mta rest client
func (d DefaultClientFactory) NewMtaClient(host, spaceID string, rt http.RoundTripper, jar http.CookieJar, tokenFactory baseclient.TokenFactory) mtaclient.MtaClientOperations {
	if mtaClient == nil {
		mtaClient = mtaclient.NewMtaClient(host, spaceID, rt, jar, tokenFactory)
	}
	return mtaclient
}

// NewRestClient used for creating or returning cached value of the rest client
func (d DefaultClientFactory) NewRestClient(host, org, space string,
	rt http.RoundTripper, jar http.CookieJar, tokenfactory baseclient.TokenFactory) restclient.RestClientOperations {
	if restClient == nil {
		restClient = restclient.NewRetryableRestClient(host, org, space, rt, jar, tokenfactory)
	}
	return restclient
}
