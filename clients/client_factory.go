package clients

import (
	"net/http"

	"github.com/cloudfoundry/multiapps-cli-plugin/clients/baseclient"
	"github.com/cloudfoundry/multiapps-cli-plugin/clients/mtaclient"
	"github.com/cloudfoundry/multiapps-cli-plugin/clients/mtaclient_v2"
	"github.com/cloudfoundry/multiapps-cli-plugin/clients/restclient"
)

// ClientFactory is a factory for creating XxxClientOperations instances
type ClientFactory interface {
	NewMtaClient(host, spaceID string, rt http.RoundTripper, tokenFactory baseclient.TokenFactory) mtaclient.MtaClientOperations
	NewRestClient(host string, rt http.RoundTripper, tokenFactory baseclient.TokenFactory) restclient.RestClientOperations
	NewMtaV2Client(host, spaceGUID string, rt http.RoundTripper, tokenFactory baseclient.TokenFactory) mtaclient_v2.MtaV2ClientOperations
}

// DefaultClientFactory a default implementation of the ClientFactory
type DefaultClientFactory struct {
	mtaClient   mtaclient.MtaClientOperations
	restClient  restclient.RestClientOperations
	mtaV2Client mtaclient_v2.MtaV2ClientOperations
}

// NewDefaultClientFactory a default initialization method for the factory
func NewDefaultClientFactory() *DefaultClientFactory {
	return &DefaultClientFactory{mtaClient: nil, restClient: nil}
}

// NewMtaClient used for creating or returning cached value of the mta rest client
func (d *DefaultClientFactory) NewMtaClient(host, spaceID string, rt http.RoundTripper, tokenFactory baseclient.TokenFactory) mtaclient.MtaClientOperations {
	if d.mtaClient == nil {
		d.mtaClient = mtaclient.NewRetryableMtaRestClient(host, spaceID, rt, tokenFactory)
	}
	return d.mtaClient
}

// NewMtaClient used for creating or returning cached value of the mta rest client
func (d *DefaultClientFactory) NewMtaV2Client(host, spaceGUID string, rt http.RoundTripper, tokenFactory baseclient.TokenFactory) mtaclient_v2.MtaV2ClientOperations {
	if d.mtaV2Client == nil {
		d.mtaV2Client = mtaclient_v2.NewRetryableMtaRestClient(host, spaceGUID, rt, tokenFactory)
	}
	return d.mtaV2Client
}

// NewRestClient used for creating or returning cached value of the rest client
func (d *DefaultClientFactory) NewRestClient(host string, rt http.RoundTripper, tokenFactory baseclient.TokenFactory) restclient.RestClientOperations {
	if d.restClient == nil {
		d.restClient = restclient.NewRetryableRestClient(host, rt, tokenFactory)
	}
	return d.restClient
}
