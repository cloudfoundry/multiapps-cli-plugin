package restclient

import (
	"context"
	"net/http"

	"github.com/go-openapi/strfmt"

	"github.com/SAP/cf-mta-plugin/clients/baseclient"
	"github.com/SAP/cf-mta-plugin/clients/models"
	operations "github.com/SAP/cf-mta-plugin/clients/restclient/operations"
)

// RestClient represents a client for the MTA deployer REST protocol
type RestClient struct {
	baseclient.BaseClient
	Client *Rest
}

// NewRestClient creates a new Rest client
func NewRestClient(host, org, space string, rt http.RoundTripper, jar http.CookieJar, tokenFactory baseclient.TokenFactory) RestClientOperations {
	t := baseclient.NewHTTPTransport(host, getRestURL(org, space), rt, jar)
	client := New(t, strfmt.Default)
	return RestClient{baseclient.BaseClient{tokenFactory}, client}
}

// GetOperations retrieves all ongoing operations
func (c RestClient) GetOperations(lastRequestedOperations *string, requestedStates []string) (models.Operations, error) {
	params := &operations.GetOperationsParams{
		Last:    lastRequestedOperations,
		Status:  requestedStates,
		Context: context.TODO(),
	}
	token, err := c.TokenFactory.NewToken()
	if err != nil {
		return models.Operations{}, baseclient.NewClientError(err)
	}
	resp, err := c.Client.Operations.GetOperations(params, token)
	if err != nil {
		return models.Operations{}, baseclient.NewClientError(err)
	}
	return resp.Payload, nil
}

// GetComponents retrieves all deployed components (MTAs and standalone apps)
func (c RestClient) GetComponents() (*models.Components, error) {
	params := &operations.GetComponentsParams{
		Context: context.TODO(),
	}
	token, err := c.TokenFactory.NewToken()
	if err != nil {
		return nil, baseclient.NewClientError(err)
	}
	resp, err := c.Client.Operations.GetComponents(params, token)
	if err != nil {
		return nil, baseclient.NewClientError(err)
	}
	return resp.Payload, nil
}

// GetMta retrieves the deployed MTA with the specified MTA ID
func (c RestClient) GetMta(mtaID string) (*models.Mta, error) {
	params := &operations.GetMtaParams{
		MtaID:   mtaID,
		Context: context.TODO(),
	}
	token, err := c.TokenFactory.NewToken()
	if err != nil {
		return nil, baseclient.NewClientError(err)
	}
	resp, err := c.Client.Operations.GetMta(params, token)
	if err != nil {
		return nil, baseclient.NewClientError(err)
	}
	return resp.Payload, nil
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
