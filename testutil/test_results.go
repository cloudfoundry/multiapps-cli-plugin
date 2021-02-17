package testutil

import (
	"io"
	"os"

	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/clients/models"
	"github.com/go-openapi/runtime"
)

var SimpleOperationResult = models.Operation{
	State:    "FINISHED",
	Messages: []*models.Message{&SimpleMessage},
}

var SimpleMessage = models.Message{
	ID:   0,
	Type: "INFO",
	Text: "Test message",
}

var GetMessage = func(id int64, message string) *models.Message {
	return &models.Message{
		ID:   id,
		Type: "INFO",
		Text: message,
	}
}

var OperationResult = models.Operation{
	State:       "FINISHED",
	ProcessID:   "1000",
	ProcessType: "DEPLOY",
	Messages:    []*models.Message{&SimpleMessage},
}

var SimpleMtaLog = models.Log{
	ID:          LogID,
	DisplayName: "Test log",
	Description: "Test log",
}

const ProcessID = "1000"
const LogID = "OPERATION.log"
const LogContent = "test-test-test"

//
type RuntimeResponse struct {
	code    int
	message string
}

func (r RuntimeResponse) Code() int {
	return r.code
}

func (r RuntimeResponse) Message() string {
	return r.message
}

func (r RuntimeResponse) GetHeader(header string) string {
	return ""
}

func (r RuntimeResponse) Body() io.ReadCloser {
	return nil
}

// generate a custom APIError for mocking test failures
func NewCustomError(customCode int, opName, customMessage string) *runtime.APIError {
	var customResponse = RuntimeResponse{
		code:    customCode,
		message: customMessage,
	}

	return runtime.NewAPIError(opName, customResponse, customCode)
}

//
var notFoundResponse = RuntimeResponse{
	code:    404,
	message: "Process with id 404 not found",
}

//
var ClientError = &runtime.APIError{
	OperationName: "Getting process",
	Response:      notFoundResponse,
	Code:          404,
}

//
//D41D8CD98F00B204E9800998ECF8427E -> MD5 hash for empty file
var SimpleFile = models.FileMetadata{
	ID:              "test.mtar",
	Digest:          "D41D8CD98F00B204E9800998ECF8427E",
	Name:            "test.mtar",
	DigestAlgorithm: "MD5",
	Space:           "test-space",
	Namespace:       "namespace",
}

//
func GetFile(file os.File, digest string, namespace string) *models.FileMetadata {
	stat, _ := os.Stat(file.Name())
	return &models.FileMetadata{
		ID:              stat.Name(),
		Space:           "test-space",
		Name:            stat.Name(),
		Namespace:       namespace,
		Digest:          digest,
		DigestAlgorithm: "MD5",
	}
}


func GetOperation(processID, spaceID string, mtaID string, namespace string, processType string, state string, acquiredLock bool) *models.Operation {
	return &models.Operation{
		ProcessID:    processID,
		ProcessType:  processType,
		StartedAt:    "2016-03-04T14:23:24.521Z[Etc/UTC]",
		SpaceID:      spaceID,
		User:         "admin",
		State:        models.State(state),
		AcquiredLock: acquiredLock,
		MtaID:        mtaID,
		Namespace:    namespace,
	}
}

//
func GetMta(id, version string, namespace string, modules []*models.Module, services []string) *models.Mta {
	return &models.Mta{
		Metadata: &models.Metadata{
			ID:        id,
			Version:   version,
			Namespace: namespace,
		},
		Modules:  modules,
		Services: services,
	}
}

//
func GetMtaModule(name string, services []string, providedDependencies []string) *models.Module {
	return &models.Module{
		ModuleName:            name,
		AppName:               name,
		Services:              services,
		ProvidedDendencyNames: providedDependencies,
	}
}

