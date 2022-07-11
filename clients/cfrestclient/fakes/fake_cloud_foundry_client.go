package fakes

import (
	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/clients/cfrestclient"
	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/clients/models"
)

type FakeCloudFoundryClient struct {
	domains []models.Domain
	err     error

	Apps               []models.CloudFoundryApplication
	AppsErr            error
	AppProcessStats    []models.ApplicationProcessStatistics
	AppProcessStatsErr error
	AppRoutes          []models.ApplicationRoute
	AppRoutesErr       error
	Services           []models.CloudFoundryServiceInstance
	ServicesErr        error
	ServiceBindings    []models.ServiceBinding
	ServiceBindingsErr error
}

func NewFakeCloudFoundryClient(domains []models.Domain, err error) cfrestclient.CloudFoundryOperationsExtended {
	return FakeCloudFoundryClient{domains: domains, err: err}
}

func (f FakeCloudFoundryClient) GetSharedDomains() ([]models.Domain, error) {
	return f.domains, f.err
}

func (f FakeCloudFoundryClient) GetApplications(mtaId, spaceGuid string) ([]models.CloudFoundryApplication, error) {
	return f.Apps, f.AppsErr
}

func (f FakeCloudFoundryClient) GetAppProcessStatistics(appGuid string) ([]models.ApplicationProcessStatistics, error) {
	return f.AppProcessStats, f.AppProcessStatsErr
}

func (f FakeCloudFoundryClient) GetApplicationRoutes(appGuid string) ([]models.ApplicationRoute, error) {
	return f.AppRoutes, f.AppRoutesErr
}

func (f FakeCloudFoundryClient) GetServiceInstances(mtaId, spaceGuid string) ([]models.CloudFoundryServiceInstance, error) {
	return f.Services, f.ServicesErr
}

func (f FakeCloudFoundryClient) GetServiceBindings(serviceName string) ([]models.ServiceBinding, error) {
	return f.ServiceBindings, f.ServiceBindingsErr
}
