package mtaclient

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/util"
	"github.com/go-openapi/runtime/client"
	"io"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"strconv"
	"strings"
	"time"

	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/clients/baseclient"
	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/clients/models"
	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/clients/mtaclient/operations"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"
)

const spacesURL string = "spaces/"
const restBaseURL string = "api/v1/"

type MtaRestClient struct {
	baseclient.BaseClient
	client *MtaClient

	dsHost    string
	spaceGuid string
}

type AsyncUploadJobResult struct {
	Status         string               `json:"status"`
	Error          string               `json:"error,omitempty"`
	File           *models.FileMetadata `json:"file,omitempty"`
	MtaId          string               `json:"mta_id,omitempty"`
	BytesProcessed int64                `json:"bytes_processed,omitempty"`
}

func NewMtaClient(host, spaceID string, rt http.RoundTripper, tokenFactory baseclient.TokenFactory) MtaClientOperations {
	restURL := restBaseURL + spacesURL + spaceID
	t := baseclient.NewHTTPTransport(host, restURL, rt)
	httpMtaClient := New(t, strfmt.Default)
	return &MtaRestClient{baseclient.BaseClient{TokenFactory: tokenFactory}, httpMtaClient, host, spaceID}
}

func (c MtaRestClient) ExecuteAction(operationID, actionID string) (ResponseHeader, error) {
	params := &operations.ExecuteOperationActionParams{
		OperationID: operationID,
		ActionID:    actionID,
		Context:     context.TODO(),
	}
	token, err := c.TokenFactory.NewToken()
	if err != nil {
		return ResponseHeader{}, baseclient.NewClientError(err)
	}
	resp, err := c.client.Operations.ExecuteOperationAction(params, token)
	if err != nil {
		return ResponseHeader{}, baseclient.NewClientError(err)
	}
	return ResponseHeader{Location: resp.Location}, nil
}

func (c MtaRestClient) GetMta(mtaID string) (*models.Mta, error) {
	params := &operations.GetMtaParams{
		Context: context.TODO(),
		MtaID:   mtaID,
	}

	result, err := executeRestOperation(c.TokenFactory, func(token runtime.ClientAuthInfoWriter) (interface{}, error) {
		return c.client.Operations.GetMta(params, token)
	})
	if err != nil {
		return nil, baseclient.NewClientError(err)
	}

	return result.(*operations.GetMtaOK).Payload, nil
}

func (c MtaRestClient) GetMtaFiles(namespace *string) ([]*models.FileMetadata, error) {
	params := &operations.GetMtaFilesParams{
		Context:   context.TODO(),
		Namespace: namespace,
	}
	token, err := c.TokenFactory.NewToken()
	if err != nil {
		return nil, baseclient.NewClientError(err)
	}
	resp, err := c.client.Operations.GetMtaFiles(params, token)
	if err != nil {
		return nil, baseclient.NewClientError(err)
	}
	return resp.Payload, nil
}

func (c MtaRestClient) GetMtaOperation(operationID, embed string) (*models.Operation, error) {
	params := &operations.GetMtaOperationParams{
		Context:     context.TODO(),
		OperationID: operationID,
		Embed:       &embed,
	}
	token, err := c.TokenFactory.NewToken()
	if err != nil {
		return nil, baseclient.NewClientError(err)
	}
	resp, err := c.client.Operations.GetMtaOperation(params, token)
	if err != nil {
		return nil, baseclient.NewClientError(err)
	}
	return resp.Payload, nil
}

func (c MtaRestClient) GetMtaOperationLogs(operationID string) ([]*models.Log, error) {
	params := &operations.GetMtaOperationLogsParams{
		Context:     context.TODO(),
		OperationID: operationID,
	}
	token, err := c.TokenFactory.NewToken()
	if err != nil {
		return nil, baseclient.NewClientError(err)
	}
	resp, err := c.client.Operations.GetMtaOperationLogs(params, token)
	if err != nil {
		return nil, baseclient.NewClientError(err)
	}
	return resp.Payload, nil
}

func (c MtaRestClient) GetMtaOperations(mtaId *string, last *int64, status []string) ([]*models.Operation, error) {
	params := &operations.GetMtaOperationsParams{
		Context: context.TODO(),
		MtaID:   mtaId,
		Last:    last,
		State:   status,
	}
	token, err := c.TokenFactory.NewToken()
	if err != nil {
		return nil, baseclient.NewClientError(err)
	}
	resp, err := c.client.Operations.GetMtaOperations(params, token)
	if err != nil {
		return nil, baseclient.NewClientError(err)
	}
	return resp.Payload, nil
}

func (c MtaRestClient) GetMtas() ([]*models.Mta, error) {
	token, err := c.TokenFactory.NewToken()
	if err != nil {
		return nil, baseclient.NewClientError(err)
	}
	resp, err := c.client.Operations.GetMtas(nil, token)
	if err != nil {
		return nil, baseclient.NewClientError(err)
	}
	return resp.Payload, nil
}

func (c MtaRestClient) GetOperationActions(operationID string) ([]string, error) {
	params := &operations.GetOperationActionsParams{
		Context:     context.TODO(),
		OperationID: operationID,
	}
	token, err := c.TokenFactory.NewToken()
	if err != nil {
		return nil, baseclient.NewClientError(err)
	}
	resp, err := c.client.Operations.GetOperationActions(params, token)
	if err != nil {
		return nil, baseclient.NewClientError(err)
	}
	return resp.Payload, nil
}

func (c MtaRestClient) StartMtaOperation(operation models.Operation) (ResponseHeader, error) {
	params := &operations.StartMtaOperationParams{
		Context:   context.TODO(),
		Operation: &operation,
	}
	token, err := c.TokenFactory.NewToken()
	if err != nil {
		return ResponseHeader{}, baseclient.NewClientError(err)
	}
	resp, err := c.client.Operations.StartMtaOperation(params, token)
	if err != nil {
		return ResponseHeader{}, baseclient.NewClientError(err)
	}
	return ResponseHeader{Location: resp.Location}, nil
}

func (c MtaRestClient) UploadMtaFile(file util.NamedReadSeeker, fileSize int64, namespace *string) (*models.FileMetadata, error) {
	requestUrl := "https://" + c.dsHost + "/" + restBaseURL + spacesURL + c.spaceGuid + "/files"
	if namespace != nil && len(*namespace) != 0 {
		requestUrl += "?namespace=" + *namespace
	}

	token, err := c.TokenFactory.NewRawToken()
	if err != nil {
		return nil, fmt.Errorf("could not get authentication token: %v", err)
	}

	contentLength, err := c.calculateRequestSize(file.Name(), fileSize)
	if err != nil {
		return nil, fmt.Errorf("could not calculate upload file request size: %v", err)
	}

	pipeReader, pipeWriter := io.Pipe()
	form := multipart.NewWriter(pipeWriter)

	errChan := make(chan error, 1)
	go func() {
		defer pipeWriter.Close()
		errChan <- c.writeFileToRequest(file, fileSize, form)
	}()

	ctx, done := context.WithTimeout(context.Background(), time.Hour)
	defer done()

	req, _ := http.NewRequestWithContext(ctx, http.MethodPost, requestUrl, pipeReader)
	req.Header.Set("Content-Type", form.FormDataContentType())
	req.Header.Set("Authorization", "Bearer "+token)
	req.ContentLength = contentLength

	cl := c.client.Transport.(*client.Runtime)
	httpClient := http.Client{Transport: cl.Transport, Jar: cl.Jar}

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("could not upload file %q: %v", file.Name(), err)
	}
	defer resp.Body.Close()

	pipeErr := <-errChan
	if pipeErr != nil {
		return nil, fmt.Errorf("could not upload file %q: %v", file.Name(), pipeErr)
	}
	if resp.StatusCode/100 != 2 {
		return nil, fmt.Errorf("could not upload file %q: %s", file.Name(), resp.Status)
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("could not read file upload response: %v", err)
	}

	fileEntry := &models.FileMetadata{}
	err = json.Unmarshal(bodyBytes, fileEntry)
	if err != nil {
		return nil, fmt.Errorf("could not deserialize file upload response: %v", err)
	}
	return fileEntry, nil
}

func (c MtaRestClient) calculateRequestSize(fileName string, fileSize int64) (int64, error) {
	var body bytes.Buffer
	form := multipart.NewWriter(&body)
	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition", fmt.Sprintf(`form-data; name="file"; filename="%s"`, escapeQuotes(fileName)))
	h.Set("Content-Type", "application/octet-stream")
	h.Set("Content-Length", strconv.FormatInt(fileSize, 10))
	_, err := form.CreatePart(h)
	if err != nil {
		return 0, err
	}
	err = form.Close()
	if err != nil {
		return 0, err
	}
	return int64(body.Len()) + fileSize, nil
}

func (c MtaRestClient) writeFileToRequest(file util.NamedReadSeeker, fileSize int64, form *multipart.Writer) error {
	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition", fmt.Sprintf(`form-data; name="file"; filename="%s"`, escapeQuotes(file.Name())))
	h.Set("Content-Type", "application/octet-stream")
	h.Set("Content-Length", strconv.FormatInt(fileSize, 10))
	fileWriter, err := form.CreatePart(h)
	if err != nil {
		return fmt.Errorf("could not create multipart file part: %v", err)
	}

	_, err = io.Copy(fileWriter, file)
	if err != nil {
		return fmt.Errorf("could not write file to HTTP request: %v", err)
	}

	err = form.Close()
	if err != nil {
		return fmt.Errorf("could not write end boundary to HTTP request: %v", err)
	}
	return nil
}

func escapeQuotes(s string) string {
	return strings.NewReplacer("\\", "\\\\", `"`, "\\\"").Replace(s)
}

func (c MtaRestClient) StartUploadMtaArchiveFromUrl(fileUrl string, namespace *string) (http.Header, error) {
	requestUrl := "https://" + c.dsHost + "/" + restBaseURL + spacesURL + c.spaceGuid + "/files/async"
	if namespace != nil && len(*namespace) != 0 {
		requestUrl += "?namespace=" + *namespace
	}

	body := struct {
		FileUrl string `json:"file_url"`
	}{fileUrl}

	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("could not serialize start async file upload request: %v", err)
	}

	req, err := http.NewRequest(http.MethodPost, requestUrl, bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, fmt.Errorf("could not create start async file upload request: %v", err)
	}

	token, err := c.TokenFactory.NewRawToken()
	if err != nil {
		return nil, fmt.Errorf("could not get authentication token: %v", err)
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	cl := c.client.Transport.(*client.Runtime)
	httpClient := http.Client{Transport: cl.Transport, Jar: cl.Jar, Timeout: 5 * time.Minute}
	httpClient.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("could not start async file upload: %v", err)
	}
	if resp.StatusCode == http.StatusTooManyRequests {
		return nil, c.handle429(resp.Header)
	}
	if resp.StatusCode/100 != 2 && resp.StatusCode/100 != 3 {
		return nil, fmt.Errorf("could not start async file upload: %s", resp.Status)
	}
	return resp.Header, nil
}

func (c MtaRestClient) handle429(headers http.Header) error {
	retryAfter := headers.Get("Retry-After")
	if len(retryAfter) == 0 {
		retryAfter = "3"
	}
	dur, err := time.ParseDuration(retryAfter + "s")
	if err != nil {
		return &baseclient.RetryAfterError{Duration: 3 * time.Second}
	}
	return &baseclient.RetryAfterError{Duration: dur}
}

func (c MtaRestClient) GetAsyncUploadJob(jobId string, namespace *string, appInstanceId string) (AsyncUploadJobResult, error) {
	requestUrl := "https://" + c.dsHost + "/" + restBaseURL + spacesURL + c.spaceGuid + "/files/jobs/" + jobId
	if namespace != nil && len(*namespace) != 0 {
		requestUrl += "?namespace=" + *namespace
	}

	req, err := http.NewRequest(http.MethodGet, requestUrl, nil)
	if err != nil {
		return AsyncUploadJobResult{}, fmt.Errorf("could not create get async file upload job request: %v", err)
	}

	token, err := c.TokenFactory.NewRawToken()
	if err != nil {
		return AsyncUploadJobResult{}, fmt.Errorf("could not get authentication token: %v", err)
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("x-cf-app-instance", appInstanceId)

	cl := c.client.Transport.(*client.Runtime)
	httpClient := http.Client{Transport: cl.Transport, Jar: cl.Jar, Timeout: 5 * time.Minute}

	resp, err := httpClient.Do(req)
	if err != nil {
		return AsyncUploadJobResult{}, fmt.Errorf("could not get async file upload job %s: %v", jobId, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode/100 != 2 {
		return AsyncUploadJobResult{}, fmt.Errorf("could not get async file upload job %s: %s", jobId, resp.Status)
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return AsyncUploadJobResult{}, fmt.Errorf("could not read async file upload job %s response: %v", jobId, err)
	}

	var responseBody AsyncUploadJobResult
	err = json.Unmarshal(bodyBytes, &responseBody)
	if err != nil {
		return AsyncUploadJobResult{}, fmt.Errorf("could not deserialize async file upload job %s response: %v", jobId, err)
	}
	return responseBody, nil
}

func (c MtaRestClient) GetMtaOperationLogContent(operationID, logID string) (string, error) {
	params := &operations.GetMtaOperationLogContentParams{
		Context:     context.TODO(),
		LogID:       logID,
		OperationID: operationID,
	}

	result, err := executeRestOperation(c.TokenFactory, func(token runtime.ClientAuthInfoWriter) (interface{}, error) {
		return c.client.Operations.GetMtaOperationLogContent(params, token)
	})

	if err != nil {
		return "", baseclient.NewClientError(err)
	}
	return result.(*operations.GetMtaOperationLogContentOK).Payload, nil
}

func executeRestOperation(tokenProvider baseclient.TokenFactory, restOperation func(token runtime.ClientAuthInfoWriter) (interface{}, error)) (interface{}, error) {
	token, err := tokenProvider.NewToken()
	if err != nil {
		return nil, baseclient.NewClientError(err)
	}
	return restOperation(token)
}
