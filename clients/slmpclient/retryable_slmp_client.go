package slmpclient

import (
	"net/http"
	"os"
	"time"

	baseclient "github.com/SAP/cf-mta-plugin/clients/baseclient"
	"github.com/SAP/cf-mta-plugin/clients/models"
)

// RetryableSlmpClient represents a client for the SLMP protocol
type RetryableSlmpClient struct {
	SlmpClient      SlmpClientOperations
	MaxRetriesCount int
	RetryInterval   time.Duration
}

// NewRetryableSlmpClient creates a new retryable SLMP client
func NewRetryableSlmpClient(host, org, space string, rt http.RoundTripper, jar http.CookieJar, tokenFactory baseclient.TokenFactory) SlmpClientOperations {
	slmpClient := NewSlmpClient(host, org, space, rt, jar, tokenFactory)
	return RetryableSlmpClient{slmpClient, 3, time.Second * 3}
}

// GetMetadata retrieves the SLMP metadata
func (c RetryableSlmpClient) GetMetadata() (*models.Metadata, error) {
	getMetadataCb := func() (interface{}, error) {
		return c.SlmpClient.GetMetadata()
	}
	resp, err := baseclient.CallWithRetry(getMetadataCb, c.MaxRetriesCount, c.RetryInterval)
	return resp.(*models.Metadata), err
}

// GetServices retrieves all services
func (c RetryableSlmpClient) GetServices() (models.Services, error) {
	getServicesCb := func() (interface{}, error) {
		return c.SlmpClient.GetServices()
	}
	resp, err := baseclient.CallWithRetry(getServicesCb, c.MaxRetriesCount, c.RetryInterval)
	return resp.(models.Services), err
}

// GetService retrieves the service with the specified service ID
func (c RetryableSlmpClient) GetService(serviceID string) (*models.Service, error) {
	getServiceCb := func() (interface{}, error) {
		return c.SlmpClient.GetService(serviceID)
	}
	resp, err := baseclient.CallWithRetry(getServiceCb, c.MaxRetriesCount, c.RetryInterval)
	return resp.(*models.Service), err
}

// GetServiceProcesses retrieves all processes for the service with the specified service ID
func (c RetryableSlmpClient) GetServiceProcesses(serviceID string) (models.Processes, error) {
	getServiceProcessesCb := func() (interface{}, error) {
		return c.SlmpClient.GetServiceProcesses(serviceID)
	}
	resp, err := baseclient.CallWithRetry(getServiceProcessesCb, c.MaxRetriesCount, c.RetryInterval)
	return resp.(models.Processes), err
}

// GetProcess retrieves process with the specified process ID
func (c RetryableSlmpClient) GetProcess(processID string) (*models.Process, error) {
	getProcessCb := func() (interface{}, error) {
		return c.SlmpClient.GetProcess(processID)
	}
	resp, err := baseclient.CallWithRetry(getProcessCb, c.MaxRetriesCount, c.RetryInterval)
	return resp.(*models.Process), err
}

// GetServiceFiles retrieves all files for the service with the specified service ID
func (c RetryableSlmpClient) GetServiceFiles(serviceID string) (models.Files, error) {
	getServiceFilesCb := func() (interface{}, error) {
		return c.SlmpClient.GetServiceFiles(serviceID)
	}
	resp, err := baseclient.CallWithRetry(getServiceFilesCb, c.MaxRetriesCount, c.RetryInterval)
	return resp.(models.Files), err
}

// CreateServiceFile uploads a file for the service with the specified service ID
func (c RetryableSlmpClient) CreateServiceFile(serviceID string, file os.File) (models.Files, error) {
	createServiceFileCb := func() (interface{}, error) {
		return c.SlmpClient.CreateServiceFile(serviceID, file)
	}
	resp, err := baseclient.CallWithRetry(createServiceFileCb, c.MaxRetriesCount, c.RetryInterval)
	return resp.(models.Files), err
}

// CreateServiceProcess creates a process for the service with the specified service ID
func (c RetryableSlmpClient) CreateServiceProcess(serviceID string, process *models.Process) (*models.Process, error) {
	createServiceProcessCb := func() (interface{}, error) {
		return c.SlmpClient.CreateServiceProcess(serviceID, process)
	}
	resp, err := baseclient.CallWithRetry(createServiceProcessCb, c.MaxRetriesCount, c.RetryInterval)
	return resp.(*models.Process), err
}

// GetServiceVersions retrieves the versions for the service with the specified service ID
func (c RetryableSlmpClient) GetServiceVersions(serviceID string) (models.Versions, error) {
	getServiceVersionsCb := func() (interface{}, error) {
		return c.SlmpClient.GetServiceVersions(serviceID)
	}
	resp, err := baseclient.CallWithRetry(getServiceVersionsCb, c.MaxRetriesCount, c.RetryInterval)
	return resp.(models.Versions), err
}
