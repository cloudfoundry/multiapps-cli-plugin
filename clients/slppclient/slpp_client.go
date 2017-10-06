package slppclient

import (
	"context"
	"net/http"

	"github.com/go-openapi/strfmt"

	baseclient "github.com/SAP/cf-mta-plugin/clients/baseclient"
	"github.com/SAP/cf-mta-plugin/clients/models"
	"github.com/SAP/cf-mta-plugin/clients/slppclient/operations"
)

// SlppClient represents a client for the SLPP protocol
type SlppClient struct {
	baseclient.BaseClient
	Client    *Slpp
	ServiceID string
}

// NewSlppClient creates a new SLPP client
func NewSlppClient(host, org, space, serviceID, processID string, rt http.RoundTripper, jar http.CookieJar, tokenFactory baseclient.TokenFactory) SlppClientOperations {
	t := baseclient.NewHTTPTransport(host, getSlppURL(org, space, serviceID, processID), rt, jar)
	client := New(t, strfmt.Default)
	return SlppClient{baseclient.BaseClient{TokenFactory: tokenFactory}, client, serviceID}

}

// GetMetadata retrieves the SLPP metadata
func (c SlppClient) GetMetadata() (*models.Metadata, error) {
	params := &operations.GetMetadataParams{
		Context: context.TODO(),
	}
	token, err := c.TokenFactory.NewToken()
	if err != nil {
		return nil, baseclient.NewClientError(err)
	}
	resp, err := c.Client.Operations.GetMetadata(params, token)
	if err != nil {
		return nil, baseclient.NewClientError(err)
	}
	return resp.Payload, nil
}

// GetLogs retrieves all process logs
func (c SlppClient) GetLogs() (models.Logs, error) {
	params := &operations.GetLogsParams{
		Context: context.TODO(),
	}
	token, err := c.TokenFactory.NewToken()
	if err != nil {
		return models.Logs{}, baseclient.NewClientError(err)
	}
	resp, err := c.Client.Operations.GetLogs(params, token)
	if err != nil {
		return models.Logs{}, baseclient.NewClientError(err)
	}
	return resp.Payload, nil
}

// GetLogContent retrieves the content of the specified logs
func (c SlppClient) GetLogContent(logID string) (string, error) {
	params := &operations.GetLogContentParams{
		LogID:   logID,
		Context: context.TODO(),
	}
	token, err := c.TokenFactory.NewToken()
	if err != nil {
		return "", baseclient.NewClientError(err)
	}
	resp, err := c.Client.Operations.GetLogContent(params, token)
	if err != nil {
		return "", baseclient.NewClientError(err)
	}
	return resp.Payload, nil
}

//GetTasklist retrieves the tasklist for the current process
func (c SlppClient) GetTasklist() (models.Tasklist, error) {
	params := &operations.GetTasklistParams{
		Context: context.TODO(),
	}
	token, err := c.TokenFactory.NewToken()
	if err != nil {
		return models.Tasklist{}, baseclient.NewClientError(err)
	}
	resp, err := c.Client.Operations.GetTasklist(params, token)
	if err != nil {
		return models.Tasklist{}, baseclient.NewClientError(err)
	}
	return resp.Payload, nil
}

//GetTasklistTask retrieves concrete task by taskID
func (c SlppClient) GetTasklistTask(taskID string) (*models.Task, error) {
	params := &operations.GetTasklistTaskParams{
		TaskID:  taskID,
		Context: context.TODO(),
	}
	token, err := c.TokenFactory.NewToken()
	if err != nil {
		return nil, baseclient.NewClientError(err)
	}
	resp, err := c.Client.Operations.GetTasklistTask(params, token)
	if err != nil {
		return nil, baseclient.NewClientError(err)
	}
	return resp.Payload, nil
}

//GetServiceID returns serviceID
func (c SlppClient) GetServiceID() string {
	return c.ServiceID
}

//GetError retrieves the current client error
func (c SlppClient) GetError() (*models.Error, error) {
	params := &operations.GetErrorParams{
		Context: context.TODO(),
	}
	token, err := c.TokenFactory.NewToken()
	if err != nil {
		return nil, baseclient.NewClientError(err)
	}
	resp, err := c.Client.Operations.GetError(params, token)
	if err != nil {
		return nil, baseclient.NewClientError(err)
	}
	return resp.Payload, nil
}

//ExecuteAction executes an action specified with actionID
func (c SlppClient) ExecuteAction(actionID string) error {
	params := &operations.ExecuteActionParams{
		ActionID: actionID,
		Context:  context.TODO(),
	}
	token, err := c.TokenFactory.NewToken()
	if err != nil {
		return baseclient.NewClientError(err)
	}
	_, err = c.Client.Operations.ExecuteAction(params, token)
	if err != nil {
		return baseclient.NewClientError(err)
	}
	return nil
}

//GetActions retrieves the list of available actions for the current process
func (c SlppClient) GetActions() (models.Actions, error) {
	params := &operations.GetActionsParams{
		Context: context.TODO(),
	}
	token, err := c.TokenFactory.NewToken()
	if err != nil {
		return models.Actions{}, baseclient.NewClientError(err)
	}
	resp, err := c.Client.Operations.GetActions(params, token)
	if err != nil {
		return models.Actions{}, baseclient.NewClientError(err)
	}
	return resp.Payload, nil
}

func getSlppURL(org, space, serviceID, processID string) string {
	return "slprot/" + org + "/" + space + "/slp/runs/" + serviceID + "/" + processID
}
