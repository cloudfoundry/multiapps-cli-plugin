package mtaclient

import (
	"os"

	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/clients/models"
	"github.com/go-openapi/strfmt"
)

// MtaClientOperations drun drun drun
type MtaClientOperations interface {
	ExecuteAction(operationID, actionID string) (ResponseHeader, error)
	GetMta(mtaID string) (*models.Mta, error)
	GetMtaFiles(namespace *string) ([]*models.FileMetadata, error)
	GetMtaOperation(operationID, embed string) (*models.Operation, error)
	GetMtaOperationLogs(operationID string) ([]*models.Log, error)
	GetMtaOperations(mtaId *string, last *int64, status []string) ([]*models.Operation, error)
	GetMtas() ([]*models.Mta, error)
	GetOperationActions(operationID string) ([]string, error)
	StartMtaOperation(operation models.Operation) (ResponseHeader, error)
	UploadMtaFile(file os.File, fileSize int64, namespace *string) (*models.FileMetadata, error)
	UploadMtaArchiveFromUrl(fileUrl string, namespace *string) (*models.FileMetadata, error)
	GetMtaOperationLogContent(operationID, logID string) (string, error)
}

// ResponseHeader response header
type ResponseHeader struct {
	Location strfmt.URI
}
