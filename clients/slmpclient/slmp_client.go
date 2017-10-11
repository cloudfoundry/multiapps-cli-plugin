package slmpclient

import (
	"context"
	"net/http"
	"os"

	"github.com/go-openapi/strfmt"

	baseclient "github.com/SAP/cf-mta-plugin/clients/baseclient"
	"github.com/SAP/cf-mta-plugin/clients/models"
	"github.com/SAP/cf-mta-plugin/clients/slmpclient/operations"
)

// SlmpClient represents a client for the SLMP protocol
type SlmpClient struct {
	baseclient.BaseClient
	Client *Slmp
}

// NewSlmpClient creates a new SLMP client
func NewSlmpClient(host, org, space string, rt http.RoundTripper, jar http.CookieJar, tokenFactory baseclient.TokenFactory) SlmpClientOperations {
	t := baseclient.NewHTTPTransport(host, getSlmpURL(org, space), getSlmpURL(baseclient.EncodeArg(org), baseclient.EncodeArg(org)), rt, jar)
	client := New(t, strfmt.Default)
	return SlmpClient{baseclient.BaseClient{TokenFactory: tokenFactory}, client}
}

// GetMetadata retrieves the SLMP metadata
func (c SlmpClient) GetMetadata() (*models.Metadata, error) {
	params := &operations.GetMetadataParams{
		Context: context.TODO(),
	}
	token, err := c.TokenFactory.NewToken()
	if err != nil {
		return nil, baseclient.NewClientError(err)
	}
	resp, err := c.Client.Operations.GetMetadata(params, token)
	if err != nil {
		return nil, baseclient.NewClientError(err)
	}
	return resp.Payload, nil
}

// GetServices retrieves all services
func (c SlmpClient) GetServices() (models.Services, error) {
	params := &operations.GetServicesParams{
		Context: context.TODO(),
	}
	token, err := c.TokenFactory.NewToken()
	if err != nil {
		return models.Services{}, baseclient.NewClientError(err)
	}
	resp, err := c.Client.Operations.GetServices(params, token)
	if err != nil {
		return models.Services{}, baseclient.NewClientError(err)
	}
	return resp.Payload, nil
}

// GetService retrieves the service with the specified service ID
func (c SlmpClient) GetService(serviceID string) (*models.Service, error) {
	params := &operations.GetServiceParams{
		ServiceID: serviceID,
		Context:   context.TODO(),
	}
	token, err := c.TokenFactory.NewToken()
	if err != nil {
		return nil, baseclient.NewClientError(err)
	}
	resp, err := c.Client.Operations.GetService(params, token)
	if err != nil {
		return nil, baseclient.NewClientError(err)
	}
	return resp.Payload, nil
}

// GetServiceProcesses retrieves all processes for the service with the specified service ID
func (c SlmpClient) GetServiceProcesses(serviceID string) (models.Processes, error) {
	params := &operations.GetServiceProcessesParams{
		ServiceID: serviceID,
		Context:   context.TODO(),
	}
	token, err := c.TokenFactory.NewToken()
	if err != nil {
		return models.Processes{}, baseclient.NewClientError(err)
	}
	resp, err := c.Client.Operations.GetServiceProcesses(params, token)
	if err != nil {
		return models.Processes{}, baseclient.NewClientError(err)
	}
	return resp.Payload, nil
}

func (c SlmpClient) GetProcess(processID string) (*models.Process, error) {
	params := &operations.GetProcessParams{
		ProcessID: processID,
		Context:   context.TODO(),
	}
	token, err := c.TokenFactory.NewToken()
	if err != nil {
		return nil, baseclient.NewClientError(err)
	}
	resp, err := c.Client.Operations.GetProcess(params, token)
	if err != nil {
		return nil, baseclient.NewClientError(err)
	}
	return resp.Payload, nil
}

// GetServiceFiles retrieves all files for the service with the specified service ID
func (c SlmpClient) GetServiceFiles(serviceID string) (models.Files, error) {
	params := &operations.GetServiceFilesParams{
		ServiceID: serviceID,
		Context:   context.TODO(),
	}
	token, err := c.TokenFactory.NewToken()
	if err != nil {
		return models.Files{}, baseclient.NewClientError(err)
	}
	resp, err := c.Client.Operations.GetServiceFiles(params, token)
	if err != nil {
		return models.Files{}, baseclient.NewClientError(err)
	}
	return resp.Payload, nil
}

// CreateServiceFile uploads a file for the service with the specified service ID
func (c SlmpClient) CreateServiceFile(serviceID string, file os.File) (models.Files, error) {
	params := &operations.CreateServiceFilesParams{
		ServiceID: serviceID,
		Files:     file,
		Context:   context.TODO(),
	}
	token, err := c.TokenFactory.NewToken()
	if err != nil {
		return models.Files{}, baseclient.NewClientError(err)
	}
	resp, err := c.Client.Operations.CreateServiceFiles(params, token)
	if err != nil {
		return models.Files{}, baseclient.NewClientError(err)
	}
	return resp.Payload, nil
}

// CreateServiceProcess creates a process for the service with the specified service ID
func (c SlmpClient) CreateServiceProcess(serviceID string, process *models.Process) (*models.Process, error) {
	params := &operations.CreateServiceProcessParams{
		ServiceID: serviceID,
		Process:   process,
		Context:   context.TODO(),
	}
	token, err := c.TokenFactory.NewToken()
	if err != nil {
		return nil, baseclient.NewClientError(err)
	}
	resp, err := c.Client.Operations.CreateServiceProcess(params, token)
	if err != nil {
		return nil, baseclient.NewClientError(err)
	}
	return resp.Payload, nil
}

// GetServiceVersions retrieves the versions for the service with the specified service ID
func (c SlmpClient) GetServiceVersions(serviceID string) (models.Versions, error) {
	params := &operations.GetServiceVersionsParams{
		ServiceID: serviceID,
		Context:   context.TODO(),
	}
	token, err := c.TokenFactory.NewToken()
	if err != nil {
		return models.Versions{}, baseclient.NewClientError(err)
	}
	resp, err := c.Client.Operations.GetServiceVersions(params, token)
	if err != nil {
		return models.Versions{}, baseclient.NewClientError(err)
	}
	return resp.Payload, nil
}

func getSlmpURL(org, space string) string {
	return "slprot/" + org + "/" + space + "/slp"
}
