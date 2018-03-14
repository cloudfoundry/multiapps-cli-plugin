package cfrestclient

import (
	models "github.com/SAP/cf-mta-plugin/clients/models"
)

type CloudFoundryOperationsExtended interface {
	GetSharedDomains() ([]models.SharedDomain, error)
}

type CloudFoundryUrlElements struct {
	Page           *string
	ResultsPerPage *string
	OrderDirection *string
}
