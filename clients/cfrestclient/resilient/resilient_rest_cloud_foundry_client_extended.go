package resilient

import (
	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/clients/cfrestclient"
	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/clients/models"
	"time"
)

type ResilientCloudFoundryRestClient struct {
	CloudFoundryRestClient cfrestclient.CloudFoundryOperationsExtended
	MaxRetriesCount        int
	RetryInterval          time.Duration
}

func NewResilientCloudFoundryClient(cloudFoundryRestClient cfrestclient.CloudFoundryOperationsExtended, maxRetriesCount int, retryIntervalInSeconds int) cfrestclient.CloudFoundryOperationsExtended {
	return &ResilientCloudFoundryRestClient{cloudFoundryRestClient, maxRetriesCount, time.Second * time.Duration(retryIntervalInSeconds)}
}

func (c ResilientCloudFoundryRestClient) GetSharedDomains() ([]models.SharedDomain, error) {
	sharedDomains, err := c.CloudFoundryRestClient.GetSharedDomains()
	for shouldRetry(c.MaxRetriesCount, err) {
		sharedDomains, err = c.CloudFoundryRestClient.GetSharedDomains()
		c.MaxRetriesCount--
		time.Sleep(c.RetryInterval)
	}
	return sharedDomains, err
}

func shouldRetry(retries int, err error) bool {
	return err != nil && retries > 0
}
