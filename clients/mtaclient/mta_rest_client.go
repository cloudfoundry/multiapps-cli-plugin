package mtaclient

import (
	"context"
	"net/http"
	"os"

	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/clients/baseclient"
	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/clients/models"
	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/clients/mtaclient/operations"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"
)

const spacesURL string = "spaces/"
const restBaseURL string = "api/v1/"

type MtaRestClient struct {
	baseclient.BaseClient
	client *MtaClient
}

func NewMtaClient(host, spaceID string, rt http.RoundTripper, tokenFactory baseclient.TokenFactory) MtaClientOperations {
	restURL := restBaseURL + spacesURL + spaceID
	t := baseclient.NewHTTPTransport(host, restURL, rt)
	httpMtaClient := New(t, strfmt.Default)
	return &MtaRestClient{baseclient.BaseClient{TokenFactory: tokenFactory}, httpMtaClient}
}

func (c MtaRestClient) ExecuteAction(operationID, actionID string) (ResponseHeader, error) {
	params := &operations.ExecuteOperationActionParams{
		OperationID: operationID,
		ActionID:    actionID,
		Context:     context.TODO(),
	}
	token, err := c.TokenFactory.NewToken()
	if err != nil {
		return ResponseHeader{}, baseclient.NewClientError(err)
	}
	resp, err := c.client.Operations.ExecuteOperationAction(params, token)
	if err != nil {
		return ResponseHeader{}, baseclient.NewClientError(err)
	}
	return ResponseHeader{Location: resp.Location}, nil
}

func (c MtaRestClient) GetMta(mtaID string) (*models.Mta, error) {
	params := &operations.GetMtaParams{
		Context: context.TODO(),
		MtaID:   mtaID,
	}

	result, err := executeRestOperation(c.TokenFactory, func(token runtime.ClientAuthInfoWriter) (interface{}, error) {
		return c.client.Operations.GetMta(params, token)
	})
	if err != nil {
		return nil, baseclient.NewClientError(err)
	}

	return result.(*operations.GetMtaOK).Payload, nil
}

func (c MtaRestClient) GetMtaFiles(namespace *string) ([]*models.FileMetadata, error) {
	params := &operations.GetMtaFilesParams{
		Context:   context.TODO(),
		Namespace: namespace,
	}
	token, err := c.TokenFactory.NewToken()
	if err != nil {
		return nil, baseclient.NewClientError(err)
	}
	resp, err := c.client.Operations.GetMtaFiles(params, token)
	if err != nil {
		return nil, baseclient.NewClientError(err)
	}
	return resp.Payload, nil
}

func (c MtaRestClient) GetMtaOperation(operationID, embed string) (*models.Operation, error) {
	params := &operations.GetMtaOperationParams{
		Context:     context.TODO(),
		OperationID: operationID,
		Embed:       &embed,
	}
	token, err := c.TokenFactory.NewToken()
	if err != nil {
		return nil, baseclient.NewClientError(err)
	}
	resp, err := c.client.Operations.GetMtaOperation(params, token)
	if err != nil {
		return nil, baseclient.NewClientError(err)
	}
	return resp.Payload, nil
}

func (c MtaRestClient) GetMtaOperationLogs(operationID string) ([]*models.Log, error) {
	params := &operations.GetMtaOperationLogsParams{
		Context:     context.TODO(),
		OperationID: operationID,
	}
	token, err := c.TokenFactory.NewToken()
	if err != nil {
		return nil, baseclient.NewClientError(err)
	}
	resp, err := c.client.Operations.GetMtaOperationLogs(params, token)
	if err != nil {
		return nil, baseclient.NewClientError(err)
	}
	return resp.Payload, nil
}

func (c MtaRestClient) GetMtaOperations(mtaId *string, last *int64, status []string) ([]*models.Operation, error) {
	params := &operations.GetMtaOperationsParams{
		Context: context.TODO(),
		MtaID:   mtaId,
		Last:    last,
		State:   status,
	}
	token, err := c.TokenFactory.NewToken()
	if err != nil {
		return nil, baseclient.NewClientError(err)
	}
	resp, err := c.client.Operations.GetMtaOperations(params, token)
	if err != nil {
		return nil, baseclient.NewClientError(err)
	}
	return resp.Payload, nil
}

func (c MtaRestClient) GetMtas() ([]*models.Mta, error) {
	token, err := c.TokenFactory.NewToken()
	if err != nil {
		return nil, baseclient.NewClientError(err)
	}
	resp, err := c.client.Operations.GetMtas(nil, token)
	if err != nil {
		return nil, baseclient.NewClientError(err)
	}
	return resp.Payload, nil
}

func (c MtaRestClient) GetOperationActions(operationID string) ([]string, error) {
	params := &operations.GetOperationActionsParams{
		Context:     context.TODO(),
		OperationID: operationID,
	}
	token, err := c.TokenFactory.NewToken()
	if err != nil {
		return nil, baseclient.NewClientError(err)
	}
	resp, err := c.client.Operations.GetOperationActions(params, token)
	if err != nil {
		return nil, baseclient.NewClientError(err)
	}
	return resp.Payload, nil
}

func (c MtaRestClient) StartMtaOperation(operation models.Operation) (ResponseHeader, error) {
	params := &operations.StartMtaOperationParams{
		Context:   context.TODO(),
		Operation: &operation,
	}
	token, err := c.TokenFactory.NewToken()
	if err != nil {
		return ResponseHeader{}, baseclient.NewClientError(err)
	}
	resp, err := c.client.Operations.StartMtaOperation(params, token)
	if err != nil {
		return ResponseHeader{}, baseclient.NewClientError(err)
	}
	return ResponseHeader{Location: resp.Location}, nil
}

func (c MtaRestClient) UploadMtaFile(file os.File, fileSize int64, namespace *string) (*models.FileMetadata, error) {
	params := &operations.UploadMtaFileParams{
		Context:   context.TODO(),
		File:      file,
		FileSize:  fileSize,
		Namespace: namespace,
	}

	result, err := executeRestOperation(c.TokenFactory, func(token runtime.ClientAuthInfoWriter) (interface{}, error) {
		return c.client.Operations.UploadMtaFile(params, token)
	})

	if err != nil {
		return nil, baseclient.NewClientError(err)
	}
	return result.(*operations.UploadMtaFileCreated).Payload, nil
}

func (c MtaRestClient) UploadMtaArchiveFromUrl(fileUrl string, namespace *string) (*models.FileMetadata, error) {
	params := &operations.UploadMtaArchiveFromUrlParams{
		Context:   context.TODO(),
		FileUrl:   fileUrl,
		Namespace: namespace,
	}

	result, err := executeRestOperation(c.TokenFactory, func(token runtime.ClientAuthInfoWriter) (interface{}, error) {
		return c.client.Operations.UploadMtaArchiveFromUrl(params, token)
	})

	if err != nil {
		return nil, baseclient.NewClientError(err)
	}
	return result.(*operations.UploadMtaArchiveFromUrlCreated).Payload, nil
}

func (c MtaRestClient) GetMtaOperationLogContent(operationID, logID string) (string, error) {
	params := &operations.GetMtaOperationLogContentParams{
		Context:     context.TODO(),
		LogID:       logID,
		OperationID: operationID,
	}

	result, err := executeRestOperation(c.TokenFactory, func(token runtime.ClientAuthInfoWriter) (interface{}, error) {
		return c.client.Operations.GetMtaOperationLogContent(params, token)
	})

	return result.(*operations.GetMtaOperationLogContentOK).Payload, baseclient.NewClientError(err)
}

func executeRestOperation(tokenProvider baseclient.TokenFactory, restOperation func(token runtime.ClientAuthInfoWriter) (interface{}, error)) (interface{}, error) {
	token, err := tokenProvider.NewToken()
	if err != nil {
		return nil, baseclient.NewClientError(err)
	}
	return restOperation(token)
}
