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
	restURL := restBaseUrl + spaceID
	t := baseclient.NewHTTPTransport(host, restURL, restURL, rt, jar)
	httpMtaClient := New(t, strfmt.Default)
	return &MtaRestClient{baseclient.BaseClient{TokenFactory: tokenFactory}, client: httpMtaClient}
}

func (c MtaRestClient) ExecuteAction(operationID, actionID string) (ResponseHeader, error) {
	params := &operations.ExecuteOperationActionParams{
		OperationID: operationID,
		ActionID: actionID,
		Context: context.TODO()
	}
	token, err := c.TokenFactory.NewToken()
	if err != nil {
		return ResponseHeader{}, baseclient.NewClientError(err)
	}
	// TODO: the token must be added to the accepting parameters. This should be done through the swagger. See the `auth` section :)
	resp, err := c.client.Operations.ExecuteOperationAction(params)
	if err != nil {
		return ResponseHeader{}, baseclient.NewClientError(err)
	}
	return ResponseHeader{Location: resp.Payload}, nil
}
func (client MtaRestClient) GetMta(mtaID string) (*models.Mta, error) {

}
func (client MtaRestClient) GetMtaFiles() ([]*File, error) {
}
func (client MtaRestClient) GetMtaOperation(operationID, embed string) (*models.Operation, error) {

}
func (client MtaRestClient) GetMtaOperationLogs(operationID string) ([]*models.Log, error) {
}
func (client MtaRestClient) GetMtaOperations() ([]*models.Operation, error) {
}
func (client MtaRestClient) GetMtas() ([]*models.Mta, error) {
}
func (client MtaRestClient) GetOperationActions(operationID string) ([]string, error) {
}
func (client MtaRestClient) StartMtaOperation(operation models.Operation) (ResponseHeader, error) {
}
func (client MtaRestClient) UploadMtaFile(file os.File) (*models.File, error) {
}
