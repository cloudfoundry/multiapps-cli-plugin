package commands

import (
	"net/http"

	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/clients/baseclient"
	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/clients/mtaclient"
	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/clients/mtaclient_v2"
	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/clients/restclient"
)

type TestClientFactory struct {
	RestClient  restclient.RestClientOperations
	MtaClient   mtaclient.MtaClientOperations
	MtaV2Client mtaclient_v2.MtaV2ClientOperations
}

func NewTestClientFactory(mtaClient mtaclient.MtaClientOperations,
	mtaV2client mtaclient_v2.MtaV2ClientOperations,
	restClient restclient.RestClientOperations) *TestClientFactory {
	return &TestClientFactory{
		RestClient:  restClient,
		MtaClient:   mtaClient,
		MtaV2Client: mtaV2client,
	}
}

func (f *TestClientFactory) NewMtaClient(host, spaceID string, rt http.RoundTripper, tokenFactory baseclient.TokenFactory) mtaclient.MtaClientOperations {
	return f.MtaClient
}

func (f *TestClientFactory) NewRestClient(host string, rt http.RoundTripper, tokenFactory baseclient.TokenFactory) restclient.RestClientOperations {
	return f.RestClient
}

func (f *TestClientFactory) NewMtaV2Client(host, spaceGUID string, rt http.RoundTripper, tokenFactory baseclient.TokenFactory) mtaclient_v2.MtaV2ClientOperations {
	return f.MtaV2Client
}
