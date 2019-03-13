package restclient

import (
	"net/http"
	"time"

	baseclient "github.com/cloudfoundry-incubator/multiapps-cli-plugin/clients/baseclient"
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

func NewRetryableManagementRestClient(host string, rt http.RoundTripper, jar http.CookieJar, tokenFactory baseclient.TokenFactory) RetryableRestClient {
	mtaManagementClient := NewManagementRestClient(host, rt, jar, tokenFactory)
	return RetryableRestClient{RestClient: mtaManagementClient, MaxRetriesCount: 3, RetryInterval: time.Second * 3}
}

// PurgeConfiguration purges a configuration
func (c RetryableRestClient) PurgeConfiguration(org, space string) error {
	purgeConfigurationCb := func() (interface{}, error) {
		return nil, c.RestClient.PurgeConfiguration(org, space)
	}
	_, err := baseclient.CallWithRetry(purgeConfigurationCb, c.MaxRetriesCount, c.RetryInterval)
	return err
}
