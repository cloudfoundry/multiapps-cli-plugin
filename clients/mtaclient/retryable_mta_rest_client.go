package mtaclient

import (
	"net/http"
	"os"
	"time"

	"github.com/cloudfoundry/multiapps-cli-plugin/clients/baseclient"
	"github.com/cloudfoundry/multiapps-cli-plugin/clients/models"
)

type RetryableMtaRestClient struct {
	mtaClient       MtaClientOperations
	MaxRetriesCount int
	RetryInterval   time.Duration
}

func NewRetryableMtaRestClient(host, spaceID string, rt http.RoundTripper, tokenFactory baseclient.TokenFactory) RetryableMtaRestClient {
	mtaClient := NewMtaClient(host, spaceID, rt, tokenFactory)
	return RetryableMtaRestClient{mtaClient: mtaClient, MaxRetriesCount: 3, RetryInterval: time.Second * 3}
}

func (c RetryableMtaRestClient) ExecuteAction(operationID, actionID string) (ResponseHeader, error) {
	executeActionCb := func() (interface{}, error) {
		return c.mtaClient.ExecuteAction(operationID, actionID)
	}
	resp, err := baseclient.CallWithRetry(executeActionCb, c.MaxRetriesCount, c.RetryInterval)
	if err != nil {
		return ResponseHeader{}, err
	}
	return resp.(ResponseHeader), nil
}

func (c RetryableMtaRestClient) GetMta(mtaID string) (*models.Mta, error) {
	getMtaCb := func() (interface{}, error) {
		return c.mtaClient.GetMta(mtaID)
	}
	resp, err := baseclient.CallWithRetry(getMtaCb, c.MaxRetriesCount, c.RetryInterval)
	if err != nil {
		return nil, err
	}
	return resp.(*models.Mta), nil
}

func (c RetryableMtaRestClient) GetMtaFiles(namespace *string) ([]*models.FileMetadata, error) {
	getMtaFilesCb := func() (interface{}, error) {
		return c.mtaClient.GetMtaFiles(namespace)
	}
	resp, err := baseclient.CallWithRetry(getMtaFilesCb, c.MaxRetriesCount, c.RetryInterval)
	if err != nil {
		return nil, err
	}
	return resp.([]*models.FileMetadata), nil
}

func (c RetryableMtaRestClient) GetMtaOperation(operationID, embed string) (*models.Operation, error) {
	getMtaOperationCb := func() (interface{}, error) {
		return c.mtaClient.GetMtaOperation(operationID, embed)
	}
	resp, err := baseclient.CallWithRetry(getMtaOperationCb, c.MaxRetriesCount, c.RetryInterval)
	if err != nil {
		return nil, err
	}
	return resp.(*models.Operation), err
}

func (c RetryableMtaRestClient) GetMtaOperationLogs(operationID string) ([]*models.Log, error) {
	getMtaOperationLogsCb := func() (interface{}, error) {
		return c.mtaClient.GetMtaOperationLogs(operationID)
	}
	resp, err := baseclient.CallWithRetry(getMtaOperationLogsCb, c.MaxRetriesCount, c.RetryInterval)
	if err != nil {
		return nil, err
	}
	return resp.([]*models.Log), nil
}

func (c RetryableMtaRestClient) GetMtaOperations(mtaId *string, last *int64, status []string) ([]*models.Operation, error) {
	getMtaOperationsCb := func() (interface{}, error) {
		return c.mtaClient.GetMtaOperations(mtaId, last, status)
	}
	resp, err := baseclient.CallWithRetry(getMtaOperationsCb, c.MaxRetriesCount, c.RetryInterval)
	if err != nil {
		return nil, err
	}
	return resp.([]*models.Operation), nil
}

func (c RetryableMtaRestClient) GetMtas() ([]*models.Mta, error) {
	getMtasCb := func() (interface{}, error) {
		return c.mtaClient.GetMtas()
	}
	resp, err := baseclient.CallWithRetry(getMtasCb, c.MaxRetriesCount, c.RetryInterval)
	if err != nil {
		return nil, err
	}
	return resp.([]*models.Mta), nil
}

func (c RetryableMtaRestClient) GetOperationActions(operationID string) ([]string, error) {
	getOperationActionsCb := func() (interface{}, error) {
		return c.mtaClient.GetOperationActions(operationID)
	}
	resp, err := baseclient.CallWithRetry(getOperationActionsCb, c.MaxRetriesCount, c.RetryInterval)
	if err != nil {
		return nil, err
	}
	return resp.([]string), nil
}

func (c RetryableMtaRestClient) StartMtaOperation(operation models.Operation) (ResponseHeader, error) {
	startMtaOperationCb := func() (interface{}, error) {
		return c.mtaClient.StartMtaOperation(operation)
	}
	resp, err := baseclient.CallWithRetry(startMtaOperationCb, c.MaxRetriesCount, c.RetryInterval)
	if err != nil {
		return ResponseHeader{}, err
	}
	return resp.(ResponseHeader), nil
}

func (c RetryableMtaRestClient) UploadMtaFile(file os.File, fileSize int64, namespace *string) (*models.FileMetadata, error) {
	uploadMtaFileCb := func() (interface{}, error) {
		reopenedFile, err := os.Open(file.Name())
		if err != nil {
			return nil, err
		}

		return c.mtaClient.UploadMtaFile(*reopenedFile, fileSize, namespace)
	}
	resp, err := baseclient.CallWithRetry(uploadMtaFileCb, c.MaxRetriesCount, c.RetryInterval)
	if err != nil {
		return nil, err
	}
	return resp.(*models.FileMetadata), nil
}

func (c RetryableMtaRestClient) UploadMtaArchiveFromUrl(fileUrl string, namespace *string) (*models.FileMetadata, error) {
	uploadMtaArchiveFromUrlCb := func() (interface{}, error) {
		return c.mtaClient.UploadMtaArchiveFromUrl(fileUrl, namespace)
	}
	resp, err := baseclient.CallWithRetry(uploadMtaArchiveFromUrlCb, c.MaxRetriesCount, c.RetryInterval)
	if err != nil {
		return nil, err
	}
	return resp.(*models.FileMetadata), nil
}

func (c RetryableMtaRestClient) GetMtaOperationLogContent(operationID, logID string) (string, error) {
	getMtaOperationLogContentCb := func() (interface{}, error) {
		return c.mtaClient.GetMtaOperationLogContent(operationID, logID)
	}
	resp, err := baseclient.CallWithRetry(getMtaOperationLogContentCb, c.MaxRetriesCount, c.RetryInterval)
	if err != nil {
		return "", err
	}
	return resp.(string), nil
}
