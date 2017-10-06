package commands

import (
	"fmt"

	"github.com/SAP/cf-mta-plugin/clients/models"
)

// ServiceID is an 'enum' representing the service IDs supported by the deploy
// service backend.
type ServiceID int

const (
	DeployServiceID ServiceID = iota
	UndeployServiceID
	BgDeployServiceID
)

var processTypes = []models.ProcessType{"deploy", "undeploy", "blue-green-deploy"}

var serviceIDs = []string{"xs2-deploy", "xs2-undeploy", "xs2-bg-deploy"}

// ToServiceID returns the service ID corresponding to the given process type.
func ToServiceID(processType models.ProcessType) (ServiceID, error) {
	for i, value := range processTypes {
		if value == processType {
			return ServiceID(i), nil
		}
	}
	return DeployServiceID, fmt.Errorf("Unknown process type: %s", processType)
}

// ProcessType returns the process type corresponding to the given service ID.
func (serviceID ServiceID) ProcessType() models.ProcessType {
	return processTypes[serviceID]
}

func (serviceID ServiceID) String() string {
	return serviceIDs[serviceID]
}
