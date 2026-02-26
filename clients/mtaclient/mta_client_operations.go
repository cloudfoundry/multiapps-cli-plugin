package mtaclient

import (
	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/clients/models"
	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/util"
	"github.com/go-openapi/strfmt"
	"net/http"
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
	UploadMtaFile(file util.NamedReadSeeker, fileSize int64, namespace *string) (*models.FileMetadata, error)
	StartUploadMtaArchiveFromUrl(fileUrl string, namespace *string) (http.Header, error)
	GetAsyncUploadJob(jobId string, namespace *string) (AsyncUploadJobResult, error)
	GetMtaOperationLogContent(operationID, logID string) (string, error)
}

// ResponseHeader response header
type ResponseHeader struct {
	Location strfmt.URI
}
