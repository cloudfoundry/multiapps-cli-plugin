package slmpclient

import (
	"os"

	"github.com/SAP/cf-mta-plugin/clients/models"
)

// SlmpClientOperations is an interface having all SlmpClient operations
type SlmpClientOperations interface {
	GetMetadata() (*models.Metadata, error)
	GetServices() (models.Services, error)
	GetService(serviceID string) (*models.Service, error)
	GetServiceProcesses(serviceID string) (models.Processes, error)
	GetServiceFiles(serviceID string) (models.Files, error)
	CreateServiceFile(serviceID string, file os.File) (models.Files, error)
	CreateServiceProcess(serviceID string, process *models.Process) (*models.Process, error)
	GetServiceVersions(serviceID string) (models.Versions, error)
	GetProcess(processID string) (*models.Process, error)
}
