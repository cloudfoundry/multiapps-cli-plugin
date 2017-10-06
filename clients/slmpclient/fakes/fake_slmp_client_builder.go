package fakes

import (
	"os"
	"path/filepath"

	models "github.com/SAP/cf-mta-plugin/clients/models"
)

type serviceResult struct {
	*models.Service
	error
}

type processesResult struct {
	models.Processes
	error
}

type filesResult struct {
	models.Files
	error
}

type processResult struct {
	*models.Process
	error
}

type versionsResult struct {
	models.Versions
	error
}

type createServiceFileParams struct {
	serviceID string
	fileName  string
}

type createServiceProcessParams struct {
	serviceID string
	processID string
}

// FakeSlmpClientBuilder is a builder of FakeSlmpClientOperations instances
type FakeSlmpClientBuilder struct {
	fakeSlmpClient              FakeSlmpClientOperations
	getServiceResults           map[string]serviceResult
	getServiceProcessesResults  map[string]processesResult
	getServiceFilesResults      map[string]filesResult
	createServiceFileResults    map[createServiceFileParams]filesResult
	createServiceProcessResults map[createServiceProcessParams]processResult
	getServiceVersionsResults   map[string]versionsResult
	getProcessResults           map[string]processResult
}

// NewFakeSlmpClientBuilder creates a new builder
func NewFakeSlmpClientBuilder() *FakeSlmpClientBuilder {
	return &FakeSlmpClientBuilder{}
}

// GetMetadata sets the metadata to return from the GetMetadata operation
func (b *FakeSlmpClientBuilder) GetMetadata(metadata *models.Metadata, err error) *FakeSlmpClientBuilder {
	b.fakeSlmpClient.GetMetadataReturns(metadata, err)
	return b
}

// GetServices sets the services to return from the GetServices operations
func (b *FakeSlmpClientBuilder) GetServices(services models.Services, err error) *FakeSlmpClientBuilder {
	b.fakeSlmpClient.GetServicesReturns(services, err)
	return b
}

// GetProcess returns a process by its corresponding id

func (b *FakeSlmpClientBuilder) GetProcess(processID string, process *models.Process, err error) *FakeSlmpClientBuilder {
	if processID == "" {
		b.fakeSlmpClient.GetProcessReturns(process, err)
		return b
	}
	if b.getProcessResults == nil {
		b.getProcessResults = make(map[string]processResult)
	}
	b.getProcessResults[processID] = processResult{process, err}
	if b.fakeSlmpClient.GetProcessStub == nil {
		b.fakeSlmpClient.GetProcessStub = func(arg0 string) (*models.Process, error) {
			result := b.getProcessResults[arg0]
			return result.Process, result.error
		}
	}
	return b
}

// GetService sets the service to return from the GetService operation
func (b *FakeSlmpClientBuilder) GetService(serviceID string, service *models.Service, err error) *FakeSlmpClientBuilder {
	if serviceID == "" {
		b.fakeSlmpClient.GetServiceReturns(service, err)
		return b
	}
	if b.getServiceResults == nil {
		b.getServiceResults = make(map[string]serviceResult)
	}
	b.getServiceResults[serviceID] = serviceResult{service, err}
	if b.fakeSlmpClient.GetServiceStub == nil {
		b.fakeSlmpClient.GetServiceStub = func(arg0 string) (*models.Service, error) {
			result := b.getServiceResults[arg0]
			return result.Service, result.error
		}
	}
	return b
}

// GetServiceProcesses sets the services processes to return from the GetServiceProcesses operation
func (b *FakeSlmpClientBuilder) GetServiceProcesses(serviceID string, processes models.Processes, err error) *FakeSlmpClientBuilder {
	if serviceID == "" {
		b.fakeSlmpClient.GetServiceProcessesReturns(processes, err)
		return b
	}
	if b.getServiceProcessesResults == nil {
		b.getServiceProcessesResults = make(map[string]processesResult)
	}
	b.getServiceProcessesResults[serviceID] = processesResult{processes, err}
	if b.fakeSlmpClient.GetServiceProcessesStub == nil {
		b.fakeSlmpClient.GetServiceProcessesStub = func(arg0 string) (models.Processes, error) {
			result := b.getServiceProcessesResults[arg0]
			return result.Processes, result.error
		}
	}
	return b
}

// GetServiceFiles sets the service files to return from the GetServiceFiles operation
func (b *FakeSlmpClientBuilder) GetServiceFiles(serviceID string, files models.Files, err error) *FakeSlmpClientBuilder {
	if serviceID == "" {
		b.fakeSlmpClient.GetServiceFilesReturns(files, err)
		return b
	}
	if b.getServiceFilesResults == nil {
		b.getServiceFilesResults = make(map[string]filesResult)
	}
	b.getServiceFilesResults[serviceID] = filesResult{files, err}
	if b.fakeSlmpClient.GetServiceFilesStub == nil {
		b.fakeSlmpClient.GetServiceFilesStub = func(arg string) (models.Files, error) {
			result := b.getServiceFilesResults[arg]
			return result.Files, result.error
		}
	}
	return b
}

// CreateServiceFile sets the files to return from the CreateServiceFile operation
func (b *FakeSlmpClientBuilder) CreateServiceFile(serviceID string, file *os.File, files models.Files, err error) *FakeSlmpClientBuilder {
	if serviceID == "" && file == nil {
		b.fakeSlmpClient.CreateServiceFileReturns(files, err)
		return b
	}
	filePath, _ := filepath.Abs(file.Name())
	if b.createServiceFileResults == nil {
		b.createServiceFileResults = make(map[createServiceFileParams]filesResult)
	}
	b.createServiceFileResults[createServiceFileParams{serviceID, filePath}] = filesResult{files, err}
	if b.fakeSlmpClient.CreateServiceFileStub == nil {
		b.fakeSlmpClient.CreateServiceFileStub = func(arg0 string, arg1 os.File) (models.Files, error) {
			arg1Path, _ := filepath.Abs(arg1.Name())
			result := b.createServiceFileResults[createServiceFileParams{arg0, arg1Path}]
			return result.Files, result.error
		}
	}
	return b
}

// CreateServiceProcess sets the process to return from the CreateServiceProcess operation
func (b *FakeSlmpClientBuilder) CreateServiceProcess(serviceID string, process *models.Process, process2 *models.Process, err error) *FakeSlmpClientBuilder {
	if serviceID == "" && process == nil {
		b.fakeSlmpClient.CreateServiceProcessReturns(process2, err)
		return b
	}
	if b.createServiceProcessResults == nil {
		b.createServiceProcessResults = make(map[createServiceProcessParams]processResult)
	}
	b.createServiceProcessResults[createServiceProcessParams{serviceID, process.ID}] = processResult{process2, err}
	if b.fakeSlmpClient.CreateServiceProcessStub == nil {
		b.fakeSlmpClient.CreateServiceProcessStub = func(arg0 string, arg1 *models.Process) (*models.Process, error) {
			result := b.createServiceProcessResults[createServiceProcessParams{arg0, arg1.ID}]
			return result.Process, result.error
		}
	}
	return b
}

// GetServiceVersions sets the versions for the service with the specified ID
func (b *FakeSlmpClientBuilder) GetServiceVersions(serviceID string, versions models.Versions, err error) *FakeSlmpClientBuilder {
	if b.getServiceVersionsResults == nil {
		b.getServiceVersionsResults = make(map[string]versionsResult)
	}

	b.getServiceVersionsResults[serviceID] = versionsResult{versions, err}
	if b.fakeSlmpClient.GetServiceVersionsStub == nil {
		b.fakeSlmpClient.GetServiceVersionsStub = func(arg0 string) (models.Versions, error) {
			result := b.getServiceVersionsResults[serviceID]
			return result.Versions, result.error
		}

	}
	return b
}

// Build builds a FakeSlmpClientOperations instance
func (b *FakeSlmpClientBuilder) Build() *FakeSlmpClientOperations {
	return &b.fakeSlmpClient
}

func safeDeref(s *string) string {
	if s != nil {
		return *s
	}
	return ""
}
