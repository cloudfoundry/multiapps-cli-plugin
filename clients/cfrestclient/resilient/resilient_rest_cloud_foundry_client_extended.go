package resilient

import (
	"time"

	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/clients/cfrestclient"
	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/clients/models"
)

type ResilientCloudFoundryRestClient struct {
	CloudFoundryRestClient cfrestclient.CloudFoundryOperationsExtended
	MaxRetriesCount        int
	RetryInterval          time.Duration
}

func NewResilientCloudFoundryClient(cloudFoundryRestClient cfrestclient.CloudFoundryOperationsExtended, maxRetriesCount int, retryIntervalInSeconds int) cfrestclient.CloudFoundryOperationsExtended {
	return &ResilientCloudFoundryRestClient{cloudFoundryRestClient, maxRetriesCount, time.Second * time.Duration(retryIntervalInSeconds)}
}

func (c ResilientCloudFoundryRestClient) GetApplications(mtaId, mtaNamespace, spaceGuid string) ([]models.CloudFoundryApplication, error) {
	return retryOnError(func() ([]models.CloudFoundryApplication, error) {
		return c.CloudFoundryRestClient.GetApplications(mtaId, mtaNamespace, spaceGuid)
	}, c.MaxRetriesCount, c.RetryInterval)
}

func (c ResilientCloudFoundryRestClient) GetAppProcessStatistics(appGuid string) ([]models.ApplicationProcessStatistics, error) {
	return retryOnError(func() ([]models.ApplicationProcessStatistics, error) {
		return c.CloudFoundryRestClient.GetAppProcessStatistics(appGuid)
	}, c.MaxRetriesCount, c.RetryInterval)
}

func (c ResilientCloudFoundryRestClient) GetApplicationRoutes(appGuid string) ([]models.ApplicationRoute, error) {
	return retryOnError(func() ([]models.ApplicationRoute, error) {
		return c.CloudFoundryRestClient.GetApplicationRoutes(appGuid)
	}, c.MaxRetriesCount, c.RetryInterval)
}

func (c ResilientCloudFoundryRestClient) GetServiceInstances(mtaId, mtaNamespace, spaceGuid string) ([]models.CloudFoundryServiceInstance, error) {
	return retryOnError(func() ([]models.CloudFoundryServiceInstance, error) {
		return c.CloudFoundryRestClient.GetServiceInstances(mtaId, mtaNamespace, spaceGuid)
	}, c.MaxRetriesCount, c.RetryInterval)
}

func (c ResilientCloudFoundryRestClient) GetServiceBindings(serviceName string) ([]models.ServiceBinding, error) {
	return retryOnError(func() ([]models.ServiceBinding, error) {
		return c.CloudFoundryRestClient.GetServiceBindings(serviceName)
	}, c.MaxRetriesCount, c.RetryInterval)
}

func (c ResilientCloudFoundryRestClient) GetServiceInstanceByName(serviceName, spaceGuid string) (models.CloudFoundryServiceInstance, error) {
	return retryOnError(func() (models.CloudFoundryServiceInstance, error) {
		return c.CloudFoundryRestClient.GetServiceInstanceByName(serviceName, spaceGuid)
	}, c.MaxRetriesCount, c.RetryInterval)
}

func retryOnError[T any](operation func() (T, error), retries int, retryInterval time.Duration) (T, error) {
	result, err := operation()
	for shouldRetry(retries, err) {
		time.Sleep(retryInterval)
		retries--
		result, err = operation()
	}
	return result, err
}

func shouldRetry(retries int, err error) bool {
	_, isResponseError := err.(models.HttpResponseError)
	return isResponseError && retries > 0
}
