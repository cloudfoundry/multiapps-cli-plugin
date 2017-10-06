package restclient

import models "github.com/SAP/cf-mta-plugin/clients/models"

// RestClientOperations is an interface having all RestClient operations
type RestClientOperations interface {
	GetOperations(lastRequestedOperations *string, requestedStates []string) (models.Operations, error)
	GetComponents() (*models.Components, error)
	GetMta(mtaID string) (*models.Mta, error)
	PurgeConfiguration(org, space string) error
}
