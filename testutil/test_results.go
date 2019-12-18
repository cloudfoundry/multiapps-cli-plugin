package testutil

import (
	"io"
	"os"

	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/clients/models"
	"github.com/go-openapi/runtime"
)

//
// import (
// 	"encoding/xml"
// 	"io"
// 	"os"
//
// 	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/clients/models"
// 	"github.com/go-openapi/runtime"
// 	"github.com/go-openapi/strfmt"
// )
//
// const xmlns string = "http://www.SAP.com/lmsl/slp"

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

// const Version = "1.2.0"
// const ServiceID = "xs2-undeploy"
const ProcessID = "1000"
const LogID = "MAIN_LOG"
const LogContent = "test-test-test"

// const ActionID = "slp.action.ABORT"
// const MtaID = "org.cloudfoundry.samples.music"
//
// var SlmpMetadataResult = models.Metadata{
// 	XMLName:     xml.Name{Space: xmlns, Local: "Metadata"},
// 	Slmpversion: Version,
// }
//
// var SlppMetadataResult = models.Metadata{
// 	XMLName:     xml.Name{Space: xmlns, Local: "Metadata"},
// 	Slppversion: Version,
// }
//
// var ServicesResult = models.Services{
// 	XMLName:  xml.Name{Space: xmlns, Local: "services"},
// 	Services: []*models.Service{&ServiceResult},
// }
//
// var ServiceResult = models.Service{
// 	XMLName:   xml.Name{Space: xmlns, Local: "Service"},
// 	ID:        uriptr(strfmt.URI(ServiceID)),
// 	Processes: uriptr(strfmt.URI("services/" + ServiceID + "/processes")),
// 	Parameters: models.ServiceParameters{
// 		Parameters: []*models.Parameter{&ParameterResult},
// 	},
// 	Files:       strfmt.URI("services/" + ServiceID + "/files"),
// 	Versions:    strfmt.URI("services/" + ServiceID + "/versions"),
// 	Slppversion: strptr(Version),
// 	DisplayName: "Undeploy",
// 	Description: "Undeploy",
// }
//
// var ParameterResult = models.Parameter{
// 	XMLName:  xml.Name{Space: xmlns, Local: "Parameter"},
// 	ID:       strptr("mtaId"),
// 	Type:     models.SlpParameterType("slp.parameter.type.SCALAR"),
// 	Required: true,
// }
//
// var ProcessesResult = models.Processes{
// 	XMLName:   xml.Name{Space: xmlns, Local: "processes"},
// 	Processes: []*models.Process{&ProcessResult},
// }
//
// var ProcessResult = models.Process{
// 	XMLName:     xml.Name{Space: xmlns, Local: "Process"},
// 	ID:          ProcessID,
// 	Service:     uriptr(strfmt.URI(ServiceID)),
// 	Status:      "slp.process.state.FINISHED",
// 	RootURL:     strfmt.URI("runs/" + ServiceID + "/" + ProcessID),
// 	DisplayName: "Undeploy",
// 	Description: "Undeploy",
// }
//
// var LogsResult = models.Logs{
// 	XMLName: xml.Name{Space: xmlns, Local: "logs"},
// 	Logs:    []*models.Log{&LogResult},
// }
//
// var LogResult = models.Log{
// 	XMLName: xml.Name{Space: xmlns, Local: "Log"},
// 	ID:      strptr(LogID),
// 	Content: uriptr(strfmt.URI("logs/" + LogID + "/content")),
// 	Format: &models.SlpLogFormat{
// 		SlpLogFormatEnum: "slp.log.format.TEXT",
// 	},
// }
//
// var OperationsResult = models.Operations{
// 	XMLName:    xml.Name{Local: "ongoing-operations"},
// 	Operations: []*models.Operation{&OperationResult},
// }
//
// var OperationResult = models.Operation{
// 	XMLName:      xml.Name{Local: "ongoing-operation"},
// 	ProcessID:    strptr(ProcessID),
// 	ProcessType:  models.ProcessType("deploy"),
// 	StartedAt:    strptr("2016-03-04T14:23:24.521Z[Etc/UTC]"),
// 	SpaceID:      strptr("5bea6497-6d70-4a31-9ad2-1ac64a520f8f"),
// 	User:         strptr("admin"),
// 	State:        models.SlpTaskStateEnum("SLP_TASK_STATE_ERROR"),
// 	AcquiredLock: nil,
// 	MtaID:        "test",
// }
//
// var ComponentsResult = models.Components{
// 	XMLName: xml.Name{Local: "components"},
// 	Mtas: models.ComponentsMtas{
// 		Mtas: []*models.Mta{&MtaResult},
// 	},
// 	StandaloneApps: models.ComponentsStandaloneApps{
// 		StandaloneApps: []string{"deploy-service"},
// 	},
// }
//
// var MtaResult = models.Mta{
// 	XMLName: xml.Name{Local: "mta"},
// 	Metadata: &models.MtaMetadata{
// 		ID:      strptr(MtaID),
// 		Version: strptr("1.0"),
// 	},
// 	Modules: models.MtaModules{
// 		Modules: []*models.MtaModulesItems0{
// 			&models.MtaModulesItems0{
// 				ModuleName: strptr("spring-music"),
// 				AppName:    strptr("spring-music"),
// 				Services: models.MtaModulesItems0Services{
// 					Services: []string{"postgresql"},
// 				},
// 				ProvidedDependencies: models.MtaModulesItems0ProvidedDependencies{
// 					ProvidedDependencies: []string{"spring-music"},
// 				},
// 			},
// 		},
// 	},
// 	Services: models.MtaServices{
// 		Services: []string{"postgresql"},
// 	},
// }
//
// var TasklistResult = models.Tasklist{
// 	XMLName: xml.Name{Space: xmlns, Local: "tasklist"},
// 	Tasks:   []*models.Task{&TaskResult},
// }
//
// var TaskResult = models.Task{
// 	XMLName:     xml.Name{Space: xmlns, Local: "Task"},
// 	ID:          strptr("startEvent"),
// 	Type:        models.SlpTaskType("slp.task.type.PROCESS"),
// 	Status:      models.SlpTaskState("slp.task.state.FINISHED"),
// 	Parent:      "roadmap_prepare",
// 	Progress:    100,
// 	RefreshRate: 10,
// 	DisplayName: "Start",
// 	Description: "Start",
// 	ProgressMessages: models.TaskProgressMessages{
// 		ProgressMessages: nil,
// 	},
// }
//
// var ErrorResult = &models.Error{
// 	XMLName:     xml.Name{Space: xmlns, Local: "Error"},
// 	ID:          strptr("test-id"),
// 	Code:        strptr("401"),
// 	DisplayName: strptr("Test display name"),
// 	Description: "Test description",
// }
//
// var FilesResult = models.Files{
// 	XMLName: xml.Name{Space: xmlns, Local: "files"},
// 	Files:   []*models.File{&FileResult},
// }
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
}

//
// var Actions = models.Actions{
// 	XMLName: xml.Name{Space: xmlns, Local: "actions"},
// 	Actions: []*models.Action{&RetryAction, &AbortAction},
// }
//
// var BlueGreenActions = models.Actions{
// 	XMLName: xml.Name{Space: xmlns, Local: "actions"},
// 	Actions: []*models.Action{&ResumeAction, &AbortAction},
// }
//
// var RetryAction = models.Action{
// 	XMLName: xml.Name{Space: xmlns, Local: "Action"},
// 	ID:      strptr("retry"),
// }
//
// var AbortAction = models.Action{
// 	XMLName: xml.Name{Space: xmlns, Local: "Action"},
// 	ID:      strptr("abort"),
// }
//
// var ResumeAction = models.Action{
// 	XMLName: xml.Name{Space: xmlns, Local: "Action"},
// 	ID:      strptr("resume"),
// }
//
// var ServiceVersion1_1 = models.Versions{
// 	XMLName:           xml.Name{Space: xmlns, Local: "Versions"},
// 	ComponentVersions: []*models.ComponentVersion{&models.ComponentVersion{Version: strptr("1.1")}},
// }
//
// func GetSlmpMetadata(version string) *models.Metadata {
// 	return &models.Metadata{
// 		XMLName:     xml.Name{Space: xmlns, Local: "Metadata"},
// 		Slmpversion: version,
// 	}
// }
//
// func GetSlppMetadata(version string) *models.Metadata {
// 	return &models.Metadata{
// 		XMLName:     xml.Name{Space: xmlns, Local: "Metadata"},
// 		Slppversion: version,
// 	}
// }
//
// func GetService(id, displayName string, parameters []*models.Parameter) *models.Service {
// 	return &models.Service{
// 		XMLName:   xml.Name{Space: xmlns, Local: "Service"},
// 		ID:        uriptr(strfmt.URI(id)),
// 		Processes: uriptr(strfmt.URI("services/" + id + "/processes")),
// 		Parameters: models.ServiceParameters{
// 			Parameters: parameters,
// 		},
// 		Files:       strfmt.URI("services/" + strfmt.URI(id) + "/files"),
// 		Versions:    strfmt.URI("services/" + strfmt.URI(id) + "/versions"),
// 		Slppversion: strptr(Version),
// 		DisplayName: displayName,
// 		Description: displayName,
// 	}
// }
//
// func GetParameter(id string) *models.Parameter {
// 	return &models.Parameter{
// 		XMLName:  xml.Name{Space: xmlns, Local: "Parameter"},
// 		ID:       strptr(id),
// 		Type:     models.SlpParameterType("slp.parameter.type.SCALAR"),
// 		Required: true,
// 	}
// }
//
// func GetFiles(files []*models.File) models.File {
// 	return models.File{
// 		XMLName: xml.Name{Space: xmlns, Local: "files"},
// 		Files:   files,
// 	}
// }

//
func GetFile(file os.File, digest string) *models.FileMetadata {
	stat, _ := os.Stat(file.Name())
	return &models.FileMetadata{
		ID:              stat.Name(),
		Space:           "test-space",
		Name:            stat.Name(),
		Digest:          digest,
		DigestAlgorithm: "MD5",
	}
}

//
// func GetComponents(mtas []*models.Mta, standaloneApps []string) *models.Components {
// 	return &models.Components{
// 		XMLName: xml.Name{Local: "components"},
// 		Mtas: models.ComponentsMtas{
// 			Mtas: mtas,
// 		},
// 		StandaloneApps: models.ComponentsStandaloneApps{
// 			StandaloneApps: standaloneApps,
// 		},
// 	}
// }
//

func GetOperation(processID, spaceID string, mtaID string, processType string, state string, acquiredLock bool) *models.Operation {
	return &models.Operation{
		ProcessID:    processID,
		ProcessType:  processType,
		StartedAt:    "2016-03-04T14:23:24.521Z[Etc/UTC]",
		SpaceID:      spaceID,
		User:         "admin",
		State:        models.State(state),
		AcquiredLock: acquiredLock,
		MtaID:        mtaID,
	}
}

//
func GetMta(id, version string, modules []*models.Module, services []string) *models.Mta {
	return &models.Mta{
		Metadata: &models.Metadata{
			ID:      id,
			Version: version,
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
		ProvidedDependencyNames: providedDependencies,
	}
}

//
// func GetTaskList(task models.Task) models.Tasklist {
// 	return models.Tasklist{
// 		XMLName: xml.Name{Space: xmlns, Local: "tasklist"},
// 		Tasks:   []*models.Task{&task},
// 	}
// }
//
// func GetTask(taskType models.SlpTaskType, state models.SlpTaskState, progressMessages []*models.ProgressMessage) models.Task {
// 	return models.Task{
// 		XMLName:     xml.Name{Space: xmlns, Local: "Task"},
// 		ID:          strptr("startEvent"),
// 		Type:        taskType,
// 		Status:      state,
// 		Parent:      "roadmap_prepare",
// 		Progress:    100,
// 		RefreshRate: 10,
// 		DisplayName: "Start",
// 		Description: "Start",
// 		ProgressMessages: models.TaskProgressMessages{
// 			ProgressMessages: progressMessages,
// 		},
// 	}
// }
//
// func GetProgressMessage(id, message string) *models.ProgressMessage {
// 	return &models.ProgressMessage{
// 		XMLName: xml.Name{Space: "ProgressMessage"},
// 		ID:      strptr(id),
// 		Message: strptr(message),
// 	}
// }
//
// func uriptr(u strfmt.URI) *strfmt.URI {
// 	return &u
// }
//
// func strptr(s string) *string {
// 	return &s
// }
//
// func boolptr(b bool) *bool {
// 	return &b
// }
