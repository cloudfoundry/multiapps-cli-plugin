package fakes

import (
	"os"

	"github.com/cloudfoundry/multiapps-cli-plugin/clients/models"
	"github.com/cloudfoundry/multiapps-cli-plugin/clients/mtaclient"
)

type FakeMtaClientBuilder struct {
	FakeMtaClient FakeMtaClientOperations
}

func NewFakeMtaClientBuilder() *FakeMtaClientBuilder {
	return &FakeMtaClientBuilder{}
}

func (fb *FakeMtaClientBuilder) ExecuteAction(operationID, actionID string, resultHeader mtaclient.ResponseHeader, resultError error) *FakeMtaClientBuilder {
	fb.FakeMtaClient.ExecuteActionReturns(resultHeader, resultError)
	return fb
}
func (fb *FakeMtaClientBuilder) GetMta(mtaID string, resultMta *models.Mta, resultError error) *FakeMtaClientBuilder {
	fb.FakeMtaClient.GetMtaReturns(resultMta, resultError)
	return fb
}
func (fb *FakeMtaClientBuilder) GetMtaFiles(result []*models.FileMetadata, resultError error) *FakeMtaClientBuilder {
	fb.FakeMtaClient.GetMtaFilesReturns(result, resultError)
	return fb
}
func (fb *FakeMtaClientBuilder) GetMtaOperation(operaWtionID, embed string, result *models.Operation, resultErr error) *FakeMtaClientBuilder {
	fb.FakeMtaClient.GetMtaOperationReturns(result, resultErr)
	return fb
}
func (fb *FakeMtaClientBuilder) GetMtaOperationLogs(operationID string, result []*models.Log, resultErr error) *FakeMtaClientBuilder {
	fb.FakeMtaClient.GetMtaOperationLogsReturns(result, resultErr)
	return fb
}
func (fb *FakeMtaClientBuilder) GetMtaOperations(mtaId *string, last *int64, status []string, result []*models.Operation, resultErr error) *FakeMtaClientBuilder {
	fb.FakeMtaClient.GetMtaOperationsReturns(result, resultErr)
	return fb
}
func (fb *FakeMtaClientBuilder) GetMtas(result []*models.Mta, resultErr error) *FakeMtaClientBuilder {
	fb.FakeMtaClient.GetMtasReturns(result, resultErr)
	return fb
}
func (fb *FakeMtaClientBuilder) GetOperationActions(operationID string, result []string, resultErr error) *FakeMtaClientBuilder {
	fb.FakeMtaClient.GetOperationActionsReturns(result, resultErr)
	return fb
}
func (fb *FakeMtaClientBuilder) StartMtaOperation(operation models.Operation, result mtaclient.ResponseHeader, resultError error) *FakeMtaClientBuilder {
	fb.FakeMtaClient.StartMtaOperationReturns(result, resultError)
	return fb
}
func (fb *FakeMtaClientBuilder) UploadMtaFile(file os.File, result *models.FileMetadata, resultError error) *FakeMtaClientBuilder {
	fb.FakeMtaClient.UploadMtaFileReturns(result, resultError)
	return fb
}
func (fb *FakeMtaClientBuilder) UploadMtaArchiveFromUrl(fileUrl string, result *models.FileMetadata, resultError error) *FakeMtaClientBuilder {
	fb.FakeMtaClient.UploadMtaArchiveFromUrlReturnsOnCall(fileUrl, result, resultError)
	return fb
}
func (fb *FakeMtaClientBuilder) GetMtaOperationLogContent(operationID, logID string, result string, resultError error) *FakeMtaClientBuilder {
	fb.FakeMtaClient.GetMtaOperationLogContentReturns(result, resultError)
	return fb
}

func (fb *FakeMtaClientBuilder) Build() FakeMtaClientOperations {
	return fb.FakeMtaClient
}
