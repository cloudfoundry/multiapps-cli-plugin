package restclient

import (
	"net/http"
	"time"

	baseclient "github.com/SAP/cf-mta-plugin/clients/baseclient"
	"github.com/SAP/cf-mta-plugin/clients/models"
)

// RetryableRestClient represents retryable REST client for the MTA deployer REST protocol
type RetryableRestClient struct {
	RestClient      RestClientOperations
	MaxRetriesCount int
	RetryInterval   time.Duration
}

// NewRetryableRestClient creates a new retryable REST client
func NewRetryableRestClient(host, org, space string, rt http.RoundTripper, jar http.CookieJar, tokenFactory baseclient.TokenFactory) RestClientOperations {
	restClient := NewRestClient(host, org, space, rt, jar, tokenFactory)
	return RetryableRestClient{restClient, 3, time.Second * 3}
}

// GetOperations retrieves all ongoing operations
func (c RetryableRestClient) GetOperations(lastRequestedOperations *string, requestedStates []string) (models.Operations, error) {
	getOperationsCb := func() (interface{}, error) {
		return c.RestClient.GetOperations(lastRequestedOperations, requestedStates)
	}
	resp, err := baseclient.CallWithRetry(getOperationsCb, c.MaxRetriesCount, c.RetryInterval)
	return resp.(models.Operations), err
}

// GetComponents retrieves all deployed components (MTAs and standalone apps)
func (c RetryableRestClient) GetComponents() (*models.Components, error) {
	getComponentsCb := func() (interface{}, error) {
		return c.RestClient.GetComponents()
	}
	resp, err := baseclient.CallWithRetry(getComponentsCb, c.MaxRetriesCount, c.RetryInterval)
	return resp.(*models.Components), err
}

// GetMta retrieves the deployed MTA with the specified MTA ID
func (c RetryableRestClient) GetMta(mtaID string) (*models.Mta, error) {
	GetMtaCb := func() (interface{}, error) {
		return c.RestClient.GetMta(mtaID)
	}
	resp, err := baseclient.CallWithRetry(GetMtaCb, c.MaxRetriesCount, c.RetryInterval)
	return resp.(*models.Mta), err
}

// PurgeConfiguration purges a configuration
func (c RetryableRestClient) PurgeConfiguration(org, space string) error {
	purgeConfigurationCb := func() (interface{}, error) {
		return nil, c.RestClient.PurgeConfiguration(org, space)
	}
	_, err := baseclient.CallWithRetry(purgeConfigurationCb, c.MaxRetriesCount, c.RetryInterval)
	return err
}
