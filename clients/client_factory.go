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
	NewManagementMtaClient(host string, rt http.RoundTripper, jar http.CookieJar, tokenFactory baseclient.TokenFactory) mtaclient.MtaClientOperations
}

// DefaultClientFactory a default implementation of the ClientFactory
type DefaultClientFactory struct {
	mtaClient           mtaclient.MtaClientOperations
	managementMtaClient mtaclient.MtaClientOperations
	restClient          restclient.RestClientOperations
}

// NewDefaultClientFactory a default intialization method for the factory
func NewDefaultClientFactory() *DefaultClientFactory {
	return &DefaultClientFactory{mtaClient: nil, restClient: nil, managementMtaClient: nil}
}

// NewMtaClient used for creating or returning cached value of the mta rest client
func (d *DefaultClientFactory) NewMtaClient(host, spaceID string, rt http.RoundTripper, jar http.CookieJar, tokenFactory baseclient.TokenFactory) mtaclient.MtaClientOperations {
	if d.mtaClient == nil {
		d.mtaClient = mtaclient.NewRetryableMtaRestClient(host, spaceID, rt, jar, tokenFactory)
	}
	return d.mtaClient
}

// NewManagementMtaClient used for creating or returning cached value of the mta rest client
func (d *DefaultClientFactory) NewManagementMtaClient(host string, rt http.RoundTripper, jar http.CookieJar, tokenFactory baseclient.TokenFactory) mtaclient.MtaClientOperations {
	if d.managementMtaClient == nil {
		d.managementMtaClient = mtaclient.NewRetryableManagementMtaRestClient(host, rt, jar, tokenFactory)
	}
	return d.managementMtaClient
}

// NewRestClient used for creating or returning cached value of the rest client
func (d *DefaultClientFactory) NewRestClient(host, org, space string,
	rt http.RoundTripper, jar http.CookieJar, tokenfactory baseclient.TokenFactory) restclient.RestClientOperations {
	if d.restClient == nil {
		d.restClient = restclient.NewRetryableRestClient(host, org, space, rt, jar, tokenfactory)
	}
	return d.restClient
}
