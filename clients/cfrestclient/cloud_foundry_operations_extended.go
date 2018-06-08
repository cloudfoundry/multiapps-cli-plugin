package cfrestclient

import (
	models "github.com/cloudfoundry-incubator/multiapps-cli-plugin/clients/models"
)

type CloudFoundryOperationsExtended interface {
	GetSharedDomains() ([]models.SharedDomain, error)
}

type CloudFoundryUrlElements struct {
	Page           *string
	ResultsPerPage *string
	OrderDirection *string
}
