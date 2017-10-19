package mtaclient

import (
	"net/http"
	"os"

	baseclient "github.com/SAP/cf-mta-plugin/clients/baseclient"
	models "github.com/SAP/cf-mta-plugin/clients/models"
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

func ExecuteAction(operationID, actionID string) (ResponseHeader, error) {

}
func GetMta(mtaID string) (*models.Mta, error) {

}
func GetMtaFiles() ([]*File, error) {
}
func GetMtaOperation(operationID, embed string) (*models.Operation, error) {

}
func GetMtaOperationLogs(operationID string) ([]*models.Log, error) {
}
func GetMtaOperations() ([]*models.Operation, error) {
}
func GetMtas() ([]*models.Mta, error) {
}
func GetOperationActions(operationID string) ([]string, error) {
}
func StartMtaOperation(operation models.Operation) (ResponseHeader, error) {
}
func UploadMtaFile(file os.File) (*models.File, error) {
}
