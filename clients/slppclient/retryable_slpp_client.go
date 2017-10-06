package slppclient

import (
	"net/http"
	"time"

	baseclient "github.com/SAP/cf-mta-plugin/clients/baseclient"
	"github.com/SAP/cf-mta-plugin/clients/models"
)

// RetryableSlppClient represents a client for the SLPP protocol
type RetryableSlppClient struct {
	SlppClient      SlppClientOperations
	MaxRetriesCount int
	RetryInterval   time.Duration
}

// NewRetryableSlppClient creates a new retryable SLPP client
func NewRetryableSlppClient(host, org, space, serviceID, processID string, rt http.RoundTripper, jar http.CookieJar, tokenFactory baseclient.TokenFactory) SlppClientOperations {
	slppClient := NewSlppClient(host, org, space, serviceID, processID, rt, jar, tokenFactory)
	return RetryableSlppClient{slppClient, 3, time.Second * 3}
}

// GetMetadata retrieves the SLPP metadata
func (c RetryableSlppClient) GetMetadata() (*models.Metadata, error) {
	getMetadataCb := func() (interface{}, error) {
		return c.SlppClient.GetMetadata()
	}
	resp, err := baseclient.CallWithRetry(getMetadataCb, c.MaxRetriesCount, c.RetryInterval)
	return resp.(*models.Metadata), err
}

// GetLogs retrieves all process logs
func (c RetryableSlppClient) GetLogs() (models.Logs, error) {
	getLogsCb := func() (interface{}, error) {
		return c.SlppClient.GetLogs()
	}
	resp, err := baseclient.CallWithRetry(getLogsCb, c.MaxRetriesCount, c.RetryInterval)
	return resp.(models.Logs), err
}

// GetLogContent retrieves the content of the specified logs
func (c RetryableSlppClient) GetLogContent(logID string) (string, error) {
	getLogContentCb := func() (interface{}, error) {
		return c.SlppClient.GetLogContent(logID)
	}
	resp, err := baseclient.CallWithRetry(getLogContentCb, c.MaxRetriesCount, c.RetryInterval)
	return resp.(string), err
}

//GetTasklist retrieves the tasklist for the current process
func (c RetryableSlppClient) GetTasklist() (models.Tasklist, error) {
	getTasklistCb := func() (interface{}, error) {
		return c.SlppClient.GetTasklist()
	}
	resp, err := baseclient.CallWithRetry(getTasklistCb, c.MaxRetriesCount, c.RetryInterval)
	return resp.(models.Tasklist), err
}

//GetTasklistTask retrieves concrete task by taskID
func (c RetryableSlppClient) GetTasklistTask(taskID string) (*models.Task, error) {
	getTasklistTaskCb := func() (interface{}, error) {
		return c.SlppClient.GetTasklistTask(taskID)
	}
	resp, err := baseclient.CallWithRetry(getTasklistTaskCb, c.MaxRetriesCount, c.RetryInterval)
	return resp.(*models.Task), err
}

//GetServiceID returns serviceID
func (c RetryableSlppClient) GetServiceID() string {
	return c.SlppClient.GetServiceID()
}

//GetError retrieves the current client error
func (c RetryableSlppClient) GetError() (*models.Error, error) {
	getErrorCb := func() (interface{}, error) {
		return c.SlppClient.GetError()
	}
	resp, err := baseclient.CallWithRetry(getErrorCb, c.MaxRetriesCount, c.RetryInterval)
	return resp.(*models.Error), err
}

//ExecuteAction executes an action specified with actionID
func (c RetryableSlppClient) ExecuteAction(actionID string) error {
	executeActionCb := func() (interface{}, error) {
		return nil, c.SlppClient.ExecuteAction(actionID)
	}
	_, err := baseclient.CallWithRetry(executeActionCb, c.MaxRetriesCount, c.RetryInterval)
	return err
}

//GetActions retrieves the list of available actions for the current process
func (c RetryableSlppClient) GetActions() (models.Actions, error) {
	getActionsCb := func() (interface{}, error) {
		return c.SlppClient.GetActions()
	}
	resp, err := baseclient.CallWithRetry(getActionsCb, c.MaxRetriesCount, c.RetryInterval)
	return resp.(models.Actions), err
}
