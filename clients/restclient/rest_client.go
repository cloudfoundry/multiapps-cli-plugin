package restclient

import (
	"context"
	"net/http"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"

	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/clients/baseclient"
	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/clients/restclient/operations"
)

const restBaseURL string = "rest/"
const csrfRestBaseURL string = "api/v1/"

// RestClient represents a client for the MTA deployer REST protocol
type RestClient struct {
	baseclient.BaseClient
	Client *Rest
}

// NewRestClient creates a new Rest client
func NewRestClient(host string, rt http.RoundTripper, tokenFactory baseclient.TokenFactory) RestClientOperations {
	t := baseclient.NewHTTPTransport(host, restBaseURL, rt)
	client := New(t, strfmt.Default)
	return RestClient{baseclient.BaseClient{tokenFactory}, client}
}

func (c RestClient) PurgeConfiguration(org, space string) error {
	params := &operations.PurgeConfigurationParams{
		Org:     org,
		Space:   space,
		Context: context.TODO(),
	}
	_, err := executeRestOperation(c.TokenFactory, func(token runtime.ClientAuthInfoWriter) (interface{}, error) {
		return c.Client.Operations.PurgeConfiguration(params, token)
	})
	if err != nil {
		return baseclient.NewClientError(err)
	}
	return nil
}

func executeRestOperation(tokenProvider baseclient.TokenFactory, restOperation func(token runtime.ClientAuthInfoWriter) (interface{}, error)) (interface{}, error) {
	token, err := tokenProvider.NewToken()
	if err != nil {
		return nil, baseclient.NewClientError(err)
	}
	return restOperation(token)
}
