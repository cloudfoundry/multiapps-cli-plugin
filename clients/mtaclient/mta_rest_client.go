package mtaclient

import (
	"context"
	"net/http"
	"os"

	baseclient "github.com/SAP/cf-mta-plugin/clients/baseclient"
	models "github.com/SAP/cf-mta-plugin/clients/models"
	operations "github.com/SAP/cf-mta-plugin/clients/mtaclient/operations"
	strfmt "github.com/go-openapi/strfmt"
)

const restBaseURL string = "api/v1/spaces/"

type MtaRestClient struct {
	baseclient.BaseClient
	client MtaClient
}

func NewMtaClient(host, spaceID string, rt http.RoundTripper, jar http.CookieJar, tokenFactory baseclient.TokenFactory) MtaClientOperations {
	restURL := restBaseURL + spaceID
	t := baseclient.NewHTTPTransport(host, restURL, restURL, rt, jar)
	httpMtaClient := New(t, strfmt.Default)
	return &MtaRestClient{baseclient.BaseClient{TokenFactory: tokenFactory}, client: httpMtaClient}
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
	return ResponseHeader{Location: resp.Payload}, nil
}

func (c MtaRestClient) GetMta(mtaID string) (*models.Mta, error) {
	params := &operations.GetMtaParams{
		Context: context.TODO(),
		MtaID:   mtaID,
	}

	result, err := executeRestOperation(c.TokenFactory, func(token string) (interface{}, error) {
		return c.client.Operations.GetMta(params, token)
	})

	return result.(*models.Mta), err
}

func (c MtaRestClient) GetMtaFiles() ([]*File, error) {
	params := &operations.GetMtaFilesParams{
		Context: context.TODO(),
	}
	token, err := c.TokenFactory.NewToken()
	if err != nil {
		return nil, baseclient.NewClientError(err)
	}
	resp, err := c.client.Operations.GetMtaFiles(params, token)
	if err != nil {
		return nil, baseclient.NewClientError(err)
	}
	return res.Payload, nil
}
func (c MtaRestClient) GetMtaOperation(operationID, embed string) (*models.Operation, error) {
	params := &operations.GetMtaOperationParms{
		Context:     context.TODO(),
		OperationID: operationID,
		Embed:       embed,
	}
	token, err := c.TokenFactory.NewToken()
	if err != nil {
		return nil, baseclient.NewClientError(err)
	}
	resp, err := c.client.Operations.GetMtaOperation(params, token)
	if err != nil {
		return nil, baseclient.NewClientError(err)
	}
	return res.Payload, nil
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
	return res.Payload, nil
}
func (c MtaRestClient) GetMtaOperations(last *string, status []string) ([]*models.Operation, error) {
	params := &operations.GetMtaOperationsParams{
		Context: context.TODO(),
		Last:    last,
		Status:  status,
	}
	token, err := c.TokenFactory.NewToken()
	if err != nil {
		return nil, baseclient.NewClientError(err)
	}
	resp, err := c.client.Operations.GetMtaOperations(params, token)
	if err != nil {
		return nil, baseclient.NewClientError(err)
	}
	return res.Payload, nil
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
	return res.Payload, nil
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
	return res.Payload, nil
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
	return ResponseHeader{Location: res.Location}, nil
}
func (c MtaRestClient) UploadMtaFile(file os.File) (*models.FileMetadata, error) {
	params := &operations.GetOperationActionsParams{
		Context: context.TODO(),
		File:    file,
	}

	result, err := executeRestOperation(c.TokenFactory, func(token string) (interface{}, error) {
		return c.client.Operations.UploadMtaFile(params, token)
	})

	return result.(*models.FileMetadata), err
}

func executeRestOperation(tokenProvider baseclient.TokenFactory, restOperation func(token string) (interface{}, error)) (interface{}, error) {
	token, err := tokenProvider.NewToken()
	if err != nil {
		return nil, baseclient.NewClientError(err)
	}
	resp, err := restOperation(token)
	if err != nil {
		return nil, baseclient.NewClientError(err)
	}
	return res.Payload, nil
}
