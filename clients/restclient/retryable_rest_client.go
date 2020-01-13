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
func NewRetryableRestClient(host string, rt http.RoundTripper, jar http.CookieJar, tokenFactory baseclient.TokenFactory) RestClientOperations {
	restClient := NewRestClient(host, rt, jar, tokenFactory)
	return RetryableRestClient{restClient, 3, time.Second * 3}
}

// PurgeConfiguration purges a configuration
func (c RetryableRestClient) PurgeConfiguration(org, space string) error {
	purgeConfigurationCb := func() (interface{}, error) {
		return nil, c.RestClient.PurgeConfiguration(org, space)
	}
	_, err := baseclient.CallWithRetry(purgeConfigurationCb, c.MaxRetriesCount, c.RetryInterval)
	return err
}
