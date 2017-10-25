package commands

import (
	"net/http"

	baseclient "github.com/SAP/cf-mta-plugin/clients/baseclient"
	"github.com/SAP/cf-mta-plugin/clients/mtaclient"
	restclient "github.com/SAP/cf-mta-plugin/clients/restclient"
	"github.wdf.sap.corp/xs2ds/cf-cli-plugin/clients/slmpclient"
	"github.wdf.sap.corp/xs2ds/cf-cli-plugin/clients/slppclient"
)

type TestClientFactory struct {
	RestClient restclient.RestClientOperations
	MtaClient  mtaclient.MtaClientOperations
}

func NewTestClientFactory(slmpClient slmpclient.SlmpClientOperations, slppClient slppclient.SlppClientOperations,
	restClient restclient.RestClientOperations) *TestClientFactory {
	return &TestClientFactory{
		SlmpClient: slmpClient,
		SlppClient: slppClient,
		RestClient: restClient,
	}
}

func (f *TestClientFactory) NewSlmpClient(host, org, space string,
	rt http.RoundTripper, jar http.CookieJar, tokenFactory baseclient.TokenFactory) slmpclient.SlmpClientOperations {
	return f.SlmpClient
}

func (f *TestClientFactory) NewSlppClient(host, org, space, serviceID, processID string,
	rt http.RoundTripper, jar http.CookieJar, tokenFactory baseclient.TokenFactory) slppclient.SlppClientOperations {
	return f.SlppClient
}

func (f *TestClientFactory) NewRestClient(host, org, space string,
	rt http.RoundTripper, jar http.CookieJar, tokenFactory baseclient.TokenFactory) restclient.RestClientOperations {
	return f.RestClient
}
