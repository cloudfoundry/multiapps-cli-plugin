package fakes

import "github.com/cloudfoundry-incubator/multiapps-cli-plugin/clients/models"

type FakeCloudFoundryClient struct {
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

func (f FakeCloudFoundryClient) GetApplications(mtaId, namespace, spaceGuid string) ([]models.CloudFoundryApplication, error) {
	return f.Apps, f.AppsErr
}

func (f FakeCloudFoundryClient) GetAppProcessStatistics(appGuid string) ([]models.ApplicationProcessStatistics, error) {
	return f.AppProcessStats, f.AppProcessStatsErr
}

func (f FakeCloudFoundryClient) GetApplicationRoutes(appGuid string) ([]models.ApplicationRoute, error) {
	return f.AppRoutes, f.AppRoutesErr
}

func (f FakeCloudFoundryClient) GetServiceInstances(mtaId string, namespace string, spaceGuid string) ([]models.CloudFoundryServiceInstance, error) {
	return f.Services, f.ServicesErr
}

func (f FakeCloudFoundryClient) GetServiceBindings(serviceName string) ([]models.ServiceBinding, error) {
	return f.ServiceBindings, f.ServiceBindingsErr
}
