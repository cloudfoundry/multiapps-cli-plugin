package cfrestclient

import (
	models "github.com/cloudfoundry-incubator/multiapps-cli-plugin/clients/models"
)

type CloudFoundryOperationsExtended interface {
	GetApplications(mtaId, mtaNamespace, spaceGuid string) ([]models.CloudFoundryApplication, error)
	GetAppProcessStatistics(appGuid string) ([]models.ApplicationProcessStatistics, error)
	GetApplicationRoutes(appGuid string) ([]models.ApplicationRoute, error)
	GetServiceInstances(mtaId, mtaNamespace, spaceGuid string) ([]models.CloudFoundryServiceInstance, error)
	GetServiceBindings(serviceName string) ([]models.ServiceBinding, error)
	GetServiceInstanceByName(serviceName, spaceGuid string) (models.CloudFoundryServiceInstance, error)
}
