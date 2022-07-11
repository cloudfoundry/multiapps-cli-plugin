package cfrestclient

import (
	models "github.com/cloudfoundry-incubator/multiapps-cli-plugin/clients/models"
)

type CloudFoundryOperationsExtended interface {
	GetSharedDomains() ([]models.Domain, error)
	GetApplications(mtaId, spaceGuid string) ([]models.CloudFoundryApplication, error)
	GetAppProcessStatistics(appGuid string) ([]models.ApplicationProcessStatistics, error)
	GetApplicationRoutes(appGuid string) ([]models.ApplicationRoute, error)
	GetServiceInstances(mtaId, spaceGuid string) ([]models.CloudFoundryServiceInstance, error)
	GetServiceBindings(serviceName string) ([]models.ServiceBinding, error)
}

type CloudFoundryUrlElements struct {
	Page           *string
	ResultsPerPage *string
}
