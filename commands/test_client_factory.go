package commands

import (
	"net/http"

	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/clients/baseclient"
	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/clients/mtaclient"
	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/clients/restclient"
)

type TestClientFactory struct {
	RestClient restclient.RestClientOperations
	MtaClient  mtaclient.MtaClientOperations
}

func NewTestClientFactory(mtaClient mtaclient.MtaClientOperations,
	restClient restclient.RestClientOperations) *TestClientFactory {
	return &TestClientFactory{
		MtaClient:  mtaClient,
		RestClient: restClient,
	}
}

func (f *TestClientFactory) NewMtaClient(host, spaceID string, rt http.RoundTripper, jar http.CookieJar, tokenFactory baseclient.TokenFactory) mtaclient.MtaClientOperations {
	return f.MtaClient
}

func (f *TestClientFactory) NewRestClient(host string, rt http.RoundTripper, jar http.CookieJar, tokenFactory baseclient.TokenFactory) restclient.RestClientOperations {
	return f.RestClient
}

func (f *TestClientFactory) NewManagementMtaClient(host string, rt http.RoundTripper, jar http.CookieJar, tokenFactory baseclient.TokenFactory) mtaclient.MtaClientOperations {
	return f.MtaClient
}

func (f *TestClientFactory) NewManagementRestClient(host string, rt http.RoundTripper, jar http.CookieJar, tokenFactory baseclient.TokenFactory) restclient.RestClientOperations {
	return f.RestClient
}
