package mtaclient

import (
	"os"

	models "github.com/SAP/cf-mta-plugin/clients/models"
	strfmt "github.com/go-openapi/strfmt"
)

// MtaClientOperations drun drun drun
type MtaClientOperations interface {
	ExecuteAction(operationID, actionID string) (ResponseHeader, error)
	GetMta(mtaID string) (*models.Mta, error)
	GetMtaFiles() ([]*File, error)
	GetMtaOperation(operationID, embed string) (*models.Operation, error)
	GetMtaOperationLogs(operationID string) ([]*models.Log, error)
	GetMtaOperations(last *string, status []string) ([]*models.Operation, error)
	GetMtas() ([]*models.Mta, error)
	GetOperationActions(operationID string) ([]string, error)
	StartMtaOperation(operation models.Operation) (ResponseHeader, error)
	UploadMtaFile(file os.File) (*models.File, error)
}

// ResponseHeader response header
type ResponseHeader struct {
	Location strfmt.URI
}
