package clients

import (
	"net/http"

	baseclient "github.com/SAP/cf-mta-plugin/clients/baseclient"
	restclient "github.com/SAP/cf-mta-plugin/clients/restclient"
	slmpclient "github.com/SAP/cf-mta-plugin/clients/slmpclient"
	slppclient "github.com/SAP/cf-mta-plugin/clients/slppclient"
)

// ClientFactory is a factory for creating XxxClientOperations instances
type ClientFactory interface {
	NewSlmpClient(host, org, space string, rt http.RoundTripper, jar http.CookieJar, tokenfactory baseclient.TokenFactory) slmpclient.SlmpClientOperations
	NewSlppClient(host, org, space, serviceID, processID string, rt http.RoundTripper, jar http.CookieJar, tokenfactory baseclient.TokenFactory) slppclient.SlppClientOperations
	NewRestClient(host, org, space string, rt http.RoundTripper, jar http.CookieJar, tokenfactory baseclient.TokenFactory) restclient.RestClientOperations
}
