package slppclient

import (
	"github.com/SAP/cf-mta-plugin/clients/models"
)

// SlppClientOperations is an interface having all SlppClient operations
type SlppClientOperations interface {
	GetMetadata() (*models.Metadata, error)
	GetLogs() (models.Logs, error)
	GetLogContent(logID string) (string, error)
	GetTasklist() (models.Tasklist, error)
	GetTasklistTask(taskID string) (*models.Task, error)
	GetServiceID() string
	GetError() (*models.Error, error)
	ExecuteAction(actionID string) error
	GetActions() (models.Actions, error)
}
