package mtaclient

import (
	"os"

	models "github.com/cloudfoundry-incubator/multiapps-cli-plugin/clients/models"
	strfmt "github.com/go-openapi/strfmt"
)

// MtaClientOperations drun drun drun
type MtaClientOperations interface {
	ExecuteAction(operationID, actionID string) (ResponseHeader, error)
	GetMta(mtaID string) (*models.Mta, error)
	GetMtaFiles() ([]*models.FileMetadata, error)
	GetMtaOperation(operationID, embed string) (*models.Operation, error)
	GetMtaOperationLogs(operationID string) ([]*models.Log, error)
	GetMtaOperations(last *int64, status []string) ([]*models.Operation, error)
	GetMtas() ([]*models.Mta, error)
	GetOperationActions(operationID string) ([]string, error)
	StartMtaOperation(operation models.Operation) (ResponseHeader, error)
	UploadMtaFile(file os.File) (*models.FileMetadata, error)
	GetMtaOperationLogContent(operationID, logID string) (string, error)
}

// ResponseHeader response header
type ResponseHeader struct {
	Location strfmt.URI
}
