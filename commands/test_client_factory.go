package commands

import (
	"net/http"

	baseclient "github.com/SAP/cf-mta-plugin/clients/baseclient"
	restclient "github.com/SAP/cf-mta-plugin/clients/restclient"
	slmpclient "github.com/SAP/cf-mta-plugin/clients/slmpclient"
	slppclient "github.com/SAP/cf-mta-plugin/clients/slppclient"
)

type TestClientFactory struct {
	SlmpClient slmpclient.SlmpClientOperations
	SlppClient slppclient.SlppClientOperations
	RestClient restclient.RestClientOperations
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
