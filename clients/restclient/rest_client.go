package restclient

import (
	"context"
	"net/http"

	"github.com/go-openapi/strfmt"

	"github.com/SAP/cf-mta-plugin/clients/baseclient"
	operations "github.com/SAP/cf-mta-plugin/clients/restclient/operations"
)

// RestClient represents a client for the MTA deployer REST protocol
type RestClient struct {
	baseclient.BaseClient
	Client *Rest
}

// NewRestClient creates a new Rest client
func NewRestClient(host, org, space string, rt http.RoundTripper, jar http.CookieJar, tokenFactory baseclient.TokenFactory) RestClientOperations {
	t := baseclient.NewHTTPTransport(host, getRestURL(org, space), getRestURL(baseclient.EncodeArg(org), baseclient.EncodeArg(space)), rt, jar)
	client := New(t, strfmt.Default)
	return RestClient{baseclient.BaseClient{tokenFactory}, client}
}

func getRestURL(org, space string) string {
	return "rest/" + org + "/" + space
}

func (c RestClient) PurgeConfiguration(org, space string) error {
	params := &operations.PurgeConfigurationParams{
		Org:     org,
		Space:   space,
		Context: context.TODO(),
	}
	token, err := c.TokenFactory.NewToken()
	if err != nil {
		return baseclient.NewClientError(err)
	}
	_, err = c.Client.Operations.PurgeConfiguration(params, token)
	if err != nil {
		return baseclient.NewClientError(err)
	}
	return nil
}
