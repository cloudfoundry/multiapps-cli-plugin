package fakes

import (
	"os"
	"sync"

	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/clients/models"
	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/clients/mtaclient"
)

type FakeMtaClientOperations struct {
	ExecuteActionStub        func(operationID, actionID string) (mtaclient.ResponseHeader, error)
	executeActionMutex       sync.RWMutex
	executeActionArgsForCall []struct {
		operationID string
		actionID    string
	}
	executeActionReturns struct {
		result1 mtaclient.ResponseHeader
		result2 error
	}
	executeActionReturnsOnCall map[int]struct {
		result1 mtaclient.ResponseHeader
		result2 error
	}
	GetMtaStub        func(mtaID string) (*models.Mta, error)
	getMtaMutex       sync.RWMutex
	getMtaArgsForCall []struct {
		mtaID string
	}
	getMtaReturns struct {
		result1 *models.Mta
		result2 error
	}
	getMtaReturnsOnCall map[int]struct {
		result1 *models.Mta
		result2 error
	}
	GetMtaFilesStub        func() ([]*models.FileMetadata, error)
	getMtaFilesMutex       sync.RWMutex
	getMtaFilesArgsForCall []struct{}
	getMtaFilesReturns     struct {
		result1 []*models.FileMetadata
		result2 error
	}
	getMtaFilesReturnsOnCall map[int]struct {
		result1 []*models.FileMetadata
		result2 error
	}
	GetMtaOperationStub        func(operationID, embed string) (*models.Operation, error)
	getMtaOperationMutex       sync.RWMutex
	getMtaOperationArgsForCall []struct {
		operationID string
		embed       string
	}
	getMtaOperationReturns struct {
		result1 *models.Operation
		result2 error
	}
	getMtaOperationReturnsOnCall map[int]struct {
		result1 *models.Operation
		result2 error
	}
	GetMtaOperationLogsStub        func(operationID string) ([]*models.Log, error)
	getMtaOperationLogsMutex       sync.RWMutex
	getMtaOperationLogsArgsForCall []struct {
		operationID string
	}
	getMtaOperationLogsReturns struct {
		result1 []*models.Log
		result2 error
	}
	getMtaOperationLogsReturnsOnCall map[int]struct {
		result1 []*models.Log
		result2 error
	}
	GetMtaOperationsStub        func(last *int64, status []string) ([]*models.Operation, error)
	getMtaOperationsMutex       sync.RWMutex
	getMtaOperationsArgsForCall []struct {
		last   *int64
		status []string
	}
	getMtaOperationsReturns struct {
		result1 []*models.Operation
		result2 error
	}
	getMtaOperationsReturnsOnCall map[int]struct {
		result1 []*models.Operation
		result2 error
	}
	GetMtasStub        func() ([]*models.Mta, error)
	getMtasMutex       sync.RWMutex
	getMtasArgsForCall []struct{}
	getMtasReturns     struct {
		result1 []*models.Mta
		result2 error
	}
	getMtasReturnsOnCall map[int]struct {
		result1 []*models.Mta
		result2 error
	}
	GetOperationActionsStub        func(operationID string) ([]string, error)
	getOperationActionsMutex       sync.RWMutex
	getOperationActionsArgsForCall []struct {
		operationID string
	}
	getOperationActionsReturns struct {
		result1 []string
		result2 error
	}
	getOperationActionsReturnsOnCall map[int]struct {
		result1 []string
		result2 error
	}
	StartMtaOperationStub        func(operation models.Operation) (mtaclient.ResponseHeader, error)
	startMtaOperationMutex       sync.RWMutex
	startMtaOperationArgsForCall []struct {
		operation models.Operation
	}
	startMtaOperationReturns struct {
		result1 mtaclient.ResponseHeader
		result2 error
	}
	startMtaOperationReturnsOnCall map[int]struct {
		result1 mtaclient.ResponseHeader
		result2 error
	}
	UploadMtaFileStub        func(file os.File) (*models.FileMetadata, error)
	uploadMtaFileMutex       sync.RWMutex
	uploadMtaFileArgsForCall []struct {
		file os.File
	}
	uploadMtaFileReturns struct {
		result1 *models.FileMetadata
		result2 error
	}
	uploadMtaFileReturnsOnCall map[int]struct {
		result1 *models.FileMetadata
		result2 error
	}
	GetMtaOperationLogContentStub        func(operationID, logID string) (string, error)
	getMtaOperationLogContentMutex       sync.RWMutex
	getMtaOperationLogContentArgsForCall []struct {
		operationID string
		logID       string
	}
	getMtaOperationLogContentReturns struct {
		result1 string
		result2 error
	}
	getMtaOperationLogContentReturnsOnCall map[int]struct {
		result1 string
		result2 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake FakeMtaClientOperations) ExecuteAction(operationID string, actionID string) (mtaclient.ResponseHeader, error) {
	fake.executeActionMutex.Lock()
	ret, specificReturn := fake.executeActionReturnsOnCall[len(fake.executeActionArgsForCall)]
	fake.executeActionArgsForCall = append(fake.executeActionArgsForCall, struct {
		operationID string
		actionID    string
	}{operationID, actionID})
	fake.recordInvocation("ExecuteAction", []interface{}{operationID, actionID})
	fake.executeActionMutex.Unlock()
	if fake.ExecuteActionStub != nil {
		return fake.ExecuteActionStub(operationID, actionID)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fake.executeActionReturns.result1, fake.executeActionReturns.result2
}

func (fake *FakeMtaClientOperations) ExecuteActionCallCount() int {
	fake.executeActionMutex.RLock()
	defer fake.executeActionMutex.RUnlock()
	return len(fake.executeActionArgsForCall)
}

func (fake *FakeMtaClientOperations) ExecuteActionArgsForCall(i int) (string, string) {
	fake.executeActionMutex.RLock()
	defer fake.executeActionMutex.RUnlock()
	return fake.executeActionArgsForCall[i].operationID, fake.executeActionArgsForCall[i].actionID
}

func (fake *FakeMtaClientOperations) ExecuteActionReturns(result1 mtaclient.ResponseHeader, result2 error) {
	fake.ExecuteActionStub = nil
	fake.executeActionReturns = struct {
		result1 mtaclient.ResponseHeader
		result2 error
	}{result1, result2}
}

func (fake *FakeMtaClientOperations) ExecuteActionReturnsOnCall(i int, result1 mtaclient.ResponseHeader, result2 error) {
	fake.ExecuteActionStub = nil
	if fake.executeActionReturnsOnCall == nil {
		fake.executeActionReturnsOnCall = make(map[int]struct {
			result1 mtaclient.ResponseHeader
			result2 error
		})
	}
	fake.executeActionReturnsOnCall[i] = struct {
		result1 mtaclient.ResponseHeader
		result2 error
	}{result1, result2}
}

func (fake FakeMtaClientOperations) GetMta(mtaID string) (*models.Mta, error) {
	fake.getMtaMutex.Lock()
	ret, specificReturn := fake.getMtaReturnsOnCall[len(fake.getMtaArgsForCall)]
	fake.getMtaArgsForCall = append(fake.getMtaArgsForCall, struct {
		mtaID string
	}{mtaID})
	fake.recordInvocation("GetMta", []interface{}{mtaID})
	fake.getMtaMutex.Unlock()
	if fake.GetMtaStub != nil {
		return fake.GetMtaStub(mtaID)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fake.getMtaReturns.result1, fake.getMtaReturns.result2
}

func (fake *FakeMtaClientOperations) GetMtaCallCount() int {
	fake.getMtaMutex.RLock()
	defer fake.getMtaMutex.RUnlock()
	return len(fake.getMtaArgsForCall)
}

func (fake *FakeMtaClientOperations) GetMtaArgsForCall(i int) string {
	fake.getMtaMutex.RLock()
	defer fake.getMtaMutex.RUnlock()
	return fake.getMtaArgsForCall[i].mtaID
}

func (fake *FakeMtaClientOperations) GetMtaReturns(result1 *models.Mta, result2 error) {
	fake.GetMtaStub = nil
	fake.getMtaReturns = struct {
		result1 *models.Mta
		result2 error
	}{result1, result2}
}

func (fake *FakeMtaClientOperations) GetMtaReturnsOnCall(i int, result1 *models.Mta, result2 error) {
	fake.GetMtaStub = nil
	if fake.getMtaReturnsOnCall == nil {
		fake.getMtaReturnsOnCall = make(map[int]struct {
			result1 *models.Mta
			result2 error
		})
	}
	fake.getMtaReturnsOnCall[i] = struct {
		result1 *models.Mta
		result2 error
	}{result1, result2}
}

func (fake FakeMtaClientOperations) GetMtaFiles() ([]*models.FileMetadata, error) {
	fake.getMtaFilesMutex.Lock()
	ret, specificReturn := fake.getMtaFilesReturnsOnCall[len(fake.getMtaFilesArgsForCall)]
	fake.getMtaFilesArgsForCall = append(fake.getMtaFilesArgsForCall, struct{}{})
	fake.recordInvocation("GetMtaFiles", []interface{}{})
	fake.getMtaFilesMutex.Unlock()
	if fake.GetMtaFilesStub != nil {
		return fake.GetMtaFilesStub()
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fake.getMtaFilesReturns.result1, fake.getMtaFilesReturns.result2
}

func (fake *FakeMtaClientOperations) GetMtaFilesCallCount() int {
	fake.getMtaFilesMutex.RLock()
	defer fake.getMtaFilesMutex.RUnlock()
	return len(fake.getMtaFilesArgsForCall)
}

func (fake *FakeMtaClientOperations) GetMtaFilesReturns(result1 []*models.FileMetadata, result2 error) {
	fake.GetMtaFilesStub = nil
	fake.getMtaFilesReturns = struct {
		result1 []*models.FileMetadata
		result2 error
	}{result1, result2}
}

func (fake *FakeMtaClientOperations) GetMtaFilesReturnsOnCall(i int, result1 []*models.FileMetadata, result2 error) {
	fake.GetMtaFilesStub = nil
	if fake.getMtaFilesReturnsOnCall == nil {
		fake.getMtaFilesReturnsOnCall = make(map[int]struct {
			result1 []*models.FileMetadata
			result2 error
		})
	}
	fake.getMtaFilesReturnsOnCall[i] = struct {
		result1 []*models.FileMetadata
		result2 error
	}{result1, result2}
}

func (fake FakeMtaClientOperations) GetMtaOperation(operationID string, embed string) (*models.Operation, error) {
	fake.getMtaOperationMutex.Lock()
	ret, specificReturn := fake.getMtaOperationReturnsOnCall[len(fake.getMtaOperationArgsForCall)]
	fake.getMtaOperationArgsForCall = append(fake.getMtaOperationArgsForCall, struct {
		operationID string
		embed       string
	}{operationID, embed})
	fake.recordInvocation("GetMtaOperation", []interface{}{operationID, embed})
	fake.getMtaOperationMutex.Unlock()
	if fake.GetMtaOperationStub != nil {
		return fake.GetMtaOperationStub(operationID, embed)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fake.getMtaOperationReturns.result1, fake.getMtaOperationReturns.result2
}

func (fake *FakeMtaClientOperations) GetMtaOperationCallCount() int {
	fake.getMtaOperationMutex.RLock()
	defer fake.getMtaOperationMutex.RUnlock()
	return len(fake.getMtaOperationArgsForCall)
}

func (fake *FakeMtaClientOperations) GetMtaOperationArgsForCall(i int) (string, string) {
	fake.getMtaOperationMutex.RLock()
	defer fake.getMtaOperationMutex.RUnlock()
	return fake.getMtaOperationArgsForCall[i].operationID, fake.getMtaOperationArgsForCall[i].embed
}

func (fake *FakeMtaClientOperations) GetMtaOperationReturns(result1 *models.Operation, result2 error) {
	fake.GetMtaOperationStub = nil
	fake.getMtaOperationReturns = struct {
		result1 *models.Operation
		result2 error
	}{result1, result2}
}

func (fake *FakeMtaClientOperations) GetMtaOperationReturnsOnCall(i int, result1 *models.Operation, result2 error) {
	fake.GetMtaOperationStub = nil
	if fake.getMtaOperationReturnsOnCall == nil {
		fake.getMtaOperationReturnsOnCall = make(map[int]struct {
			result1 *models.Operation
			result2 error
		})
	}
	fake.getMtaOperationReturnsOnCall[i] = struct {
		result1 *models.Operation
		result2 error
	}{result1, result2}
}

func (fake FakeMtaClientOperations) GetMtaOperationLogs(operationID string) ([]*models.Log, error) {
	fake.getMtaOperationLogsMutex.Lock()
	ret, specificReturn := fake.getMtaOperationLogsReturnsOnCall[len(fake.getMtaOperationLogsArgsForCall)]
	fake.getMtaOperationLogsArgsForCall = append(fake.getMtaOperationLogsArgsForCall, struct {
		operationID string
	}{operationID})
	fake.recordInvocation("GetMtaOperationLogs", []interface{}{operationID})
	fake.getMtaOperationLogsMutex.Unlock()
	if fake.GetMtaOperationLogsStub != nil {
		return fake.GetMtaOperationLogsStub(operationID)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fake.getMtaOperationLogsReturns.result1, fake.getMtaOperationLogsReturns.result2
}

func (fake *FakeMtaClientOperations) GetMtaOperationLogsCallCount() int {
	fake.getMtaOperationLogsMutex.RLock()
	defer fake.getMtaOperationLogsMutex.RUnlock()
	return len(fake.getMtaOperationLogsArgsForCall)
}

func (fake *FakeMtaClientOperations) GetMtaOperationLogsArgsForCall(i int) string {
	fake.getMtaOperationLogsMutex.RLock()
	defer fake.getMtaOperationLogsMutex.RUnlock()
	return fake.getMtaOperationLogsArgsForCall[i].operationID
}

func (fake *FakeMtaClientOperations) GetMtaOperationLogsReturns(result1 []*models.Log, result2 error) {
	fake.GetMtaOperationLogsStub = nil
	fake.getMtaOperationLogsReturns = struct {
		result1 []*models.Log
		result2 error
	}{result1, result2}
}

func (fake *FakeMtaClientOperations) GetMtaOperationLogsReturnsOnCall(i int, result1 []*models.Log, result2 error) {
	fake.GetMtaOperationLogsStub = nil
	if fake.getMtaOperationLogsReturnsOnCall == nil {
		fake.getMtaOperationLogsReturnsOnCall = make(map[int]struct {
			result1 []*models.Log
			result2 error
		})
	}
	fake.getMtaOperationLogsReturnsOnCall[i] = struct {
		result1 []*models.Log
		result2 error
	}{result1, result2}
}

func (fake FakeMtaClientOperations) GetMtaOperations(last *int64, status []string) ([]*models.Operation, error) {
	var statusCopy []string
	if status != nil {
		statusCopy = make([]string, len(status))
		copy(statusCopy, status)
	}
	fake.getMtaOperationsMutex.Lock()
	ret, specificReturn := fake.getMtaOperationsReturnsOnCall[len(fake.getMtaOperationsArgsForCall)]
	fake.getMtaOperationsArgsForCall = append(fake.getMtaOperationsArgsForCall, struct {
		last   *int64
		status []string
	}{last, statusCopy})
	fake.recordInvocation("GetMtaOperations", []interface{}{last, statusCopy})
	fake.getMtaOperationsMutex.Unlock()
	if fake.GetMtaOperationsStub != nil {
		return fake.GetMtaOperationsStub(last, status)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fake.getMtaOperationsReturns.result1, fake.getMtaOperationsReturns.result2
}

func (fake *FakeMtaClientOperations) GetMtaOperationsCallCount() int {
	fake.getMtaOperationsMutex.RLock()
	defer fake.getMtaOperationsMutex.RUnlock()
	return len(fake.getMtaOperationsArgsForCall)
}

func (fake *FakeMtaClientOperations) GetMtaOperationsArgsForCall(i int) (*int64, []string) {
	fake.getMtaOperationsMutex.RLock()
	defer fake.getMtaOperationsMutex.RUnlock()
	return fake.getMtaOperationsArgsForCall[i].last, fake.getMtaOperationsArgsForCall[i].status
}

func (fake *FakeMtaClientOperations) GetMtaOperationsReturns(result1 []*models.Operation, result2 error) {
	fake.GetMtaOperationsStub = nil
	fake.getMtaOperationsReturns = struct {
		result1 []*models.Operation
		result2 error
	}{result1, result2}
}

func (fake *FakeMtaClientOperations) GetMtaOperationsReturnsOnCall(i int, result1 []*models.Operation, result2 error) {
	fake.GetMtaOperationsStub = nil
	if fake.getMtaOperationsReturnsOnCall == nil {
		fake.getMtaOperationsReturnsOnCall = make(map[int]struct {
			result1 []*models.Operation
			result2 error
		})
	}
	fake.getMtaOperationsReturnsOnCall[i] = struct {
		result1 []*models.Operation
		result2 error
	}{result1, result2}
}

func (fake FakeMtaClientOperations) GetMtas() ([]*models.Mta, error) {
	fake.getMtasMutex.Lock()
	ret, specificReturn := fake.getMtasReturnsOnCall[len(fake.getMtasArgsForCall)]
	fake.getMtasArgsForCall = append(fake.getMtasArgsForCall, struct{}{})
	fake.recordInvocation("GetMtas", []interface{}{})
	fake.getMtasMutex.Unlock()
	if fake.GetMtasStub != nil {
		return fake.GetMtasStub()
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fake.getMtasReturns.result1, fake.getMtasReturns.result2
}

func (fake *FakeMtaClientOperations) GetMtasCallCount() int {
	fake.getMtasMutex.RLock()
	defer fake.getMtasMutex.RUnlock()
	return len(fake.getMtasArgsForCall)
}

func (fake *FakeMtaClientOperations) GetMtasReturns(result1 []*models.Mta, result2 error) {
	fake.GetMtasStub = nil
	fake.getMtasReturns = struct {
		result1 []*models.Mta
		result2 error
	}{result1, result2}
}

func (fake *FakeMtaClientOperations) GetMtasReturnsOnCall(i int, result1 []*models.Mta, result2 error) {
	fake.GetMtasStub = nil
	if fake.getMtasReturnsOnCall == nil {
		fake.getMtasReturnsOnCall = make(map[int]struct {
			result1 []*models.Mta
			result2 error
		})
	}
	fake.getMtasReturnsOnCall[i] = struct {
		result1 []*models.Mta
		result2 error
	}{result1, result2}
}

func (fake FakeMtaClientOperations) GetOperationActions(operationID string) ([]string, error) {
	fake.getOperationActionsMutex.Lock()
	ret, specificReturn := fake.getOperationActionsReturnsOnCall[len(fake.getOperationActionsArgsForCall)]
	fake.getOperationActionsArgsForCall = append(fake.getOperationActionsArgsForCall, struct {
		operationID string
	}{operationID})
	fake.recordInvocation("GetOperationActions", []interface{}{operationID})
	fake.getOperationActionsMutex.Unlock()
	if fake.GetOperationActionsStub != nil {
		return fake.GetOperationActionsStub(operationID)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fake.getOperationActionsReturns.result1, fake.getOperationActionsReturns.result2
}

func (fake *FakeMtaClientOperations) GetOperationActionsCallCount() int {
	fake.getOperationActionsMutex.RLock()
	defer fake.getOperationActionsMutex.RUnlock()
	return len(fake.getOperationActionsArgsForCall)
}

func (fake *FakeMtaClientOperations) GetOperationActionsArgsForCall(i int) string {
	fake.getOperationActionsMutex.RLock()
	defer fake.getOperationActionsMutex.RUnlock()
	return fake.getOperationActionsArgsForCall[i].operationID
}

func (fake *FakeMtaClientOperations) GetOperationActionsReturns(result1 []string, result2 error) {
	fake.GetOperationActionsStub = nil
	fake.getOperationActionsReturns = struct {
		result1 []string
		result2 error
	}{result1, result2}
}

func (fake *FakeMtaClientOperations) GetOperationActionsReturnsOnCall(i int, result1 []string, result2 error) {
	fake.GetOperationActionsStub = nil
	if fake.getOperationActionsReturnsOnCall == nil {
		fake.getOperationActionsReturnsOnCall = make(map[int]struct {
			result1 []string
			result2 error
		})
	}
	fake.getOperationActionsReturnsOnCall[i] = struct {
		result1 []string
		result2 error
	}{result1, result2}
}

func (fake FakeMtaClientOperations) StartMtaOperation(operation models.Operation) (mtaclient.ResponseHeader, error) {
	fake.startMtaOperationMutex.Lock()
	ret, specificReturn := fake.startMtaOperationReturnsOnCall[len(fake.startMtaOperationArgsForCall)]
	fake.startMtaOperationArgsForCall = append(fake.startMtaOperationArgsForCall, struct {
		operation models.Operation
	}{operation})
	fake.recordInvocation("StartMtaOperation", []interface{}{operation})
	fake.startMtaOperationMutex.Unlock()
	if fake.StartMtaOperationStub != nil {
		return fake.StartMtaOperationStub(operation)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fake.startMtaOperationReturns.result1, fake.startMtaOperationReturns.result2
}

func (fake *FakeMtaClientOperations) StartMtaOperationCallCount() int {
	fake.startMtaOperationMutex.RLock()
	defer fake.startMtaOperationMutex.RUnlock()
	return len(fake.startMtaOperationArgsForCall)
}

func (fake *FakeMtaClientOperations) StartMtaOperationArgsForCall(i int) models.Operation {
	fake.startMtaOperationMutex.RLock()
	defer fake.startMtaOperationMutex.RUnlock()
	return fake.startMtaOperationArgsForCall[i].operation
}

func (fake *FakeMtaClientOperations) StartMtaOperationReturns(result1 mtaclient.ResponseHeader, result2 error) {
	fake.StartMtaOperationStub = nil
	fake.startMtaOperationReturns = struct {
		result1 mtaclient.ResponseHeader
		result2 error
	}{result1, result2}
}

func (fake *FakeMtaClientOperations) StartMtaOperationReturnsOnCall(i int, result1 mtaclient.ResponseHeader, result2 error) {
	fake.StartMtaOperationStub = nil
	if fake.startMtaOperationReturnsOnCall == nil {
		fake.startMtaOperationReturnsOnCall = make(map[int]struct {
			result1 mtaclient.ResponseHeader
			result2 error
		})
	}
	fake.startMtaOperationReturnsOnCall[i] = struct {
		result1 mtaclient.ResponseHeader
		result2 error
	}{result1, result2}
}

func (fake FakeMtaClientOperations) UploadMtaFile(file os.File) (*models.FileMetadata, error) {
	fake.uploadMtaFileMutex.Lock()
	ret, specificReturn := fake.uploadMtaFileReturnsOnCall[len(fake.uploadMtaFileArgsForCall)]
	fake.uploadMtaFileArgsForCall = append(fake.uploadMtaFileArgsForCall, struct {
		file os.File
	}{file})
	fake.recordInvocation("UploadMtaFile", []interface{}{file})
	fake.uploadMtaFileMutex.Unlock()
	if fake.UploadMtaFileStub != nil {
		return fake.UploadMtaFileStub(file)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fake.uploadMtaFileReturns.result1, fake.uploadMtaFileReturns.result2
}

func (fake *FakeMtaClientOperations) UploadMtaFileCallCount() int {
	fake.uploadMtaFileMutex.RLock()
	defer fake.uploadMtaFileMutex.RUnlock()
	return len(fake.uploadMtaFileArgsForCall)
}

func (fake *FakeMtaClientOperations) UploadMtaFileArgsForCall(i int) os.File {
	fake.uploadMtaFileMutex.RLock()
	defer fake.uploadMtaFileMutex.RUnlock()
	return fake.uploadMtaFileArgsForCall[i].file
}

func (fake *FakeMtaClientOperations) UploadMtaFileReturns(result1 *models.FileMetadata, result2 error) {
	fake.UploadMtaFileStub = nil
	fake.uploadMtaFileReturns = struct {
		result1 *models.FileMetadata
		result2 error
	}{result1, result2}
}

func (fake *FakeMtaClientOperations) UploadMtaFileReturnsOnCall(i int, result1 *models.FileMetadata, result2 error) {
	fake.UploadMtaFileStub = nil
	if fake.uploadMtaFileReturnsOnCall == nil {
		fake.uploadMtaFileReturnsOnCall = make(map[int]struct {
			result1 *models.FileMetadata
			result2 error
		})
	}
	fake.uploadMtaFileReturnsOnCall[i] = struct {
		result1 *models.FileMetadata
		result2 error
	}{result1, result2}
}

func (fake FakeMtaClientOperations) GetMtaOperationLogContent(operationID string, logID string) (string, error) {
	fake.getMtaOperationLogContentMutex.Lock()
	ret, specificReturn := fake.getMtaOperationLogContentReturnsOnCall[len(fake.getMtaOperationLogContentArgsForCall)]
	fake.getMtaOperationLogContentArgsForCall = append(fake.getMtaOperationLogContentArgsForCall, struct {
		operationID string
		logID       string
	}{operationID, logID})
	fake.recordInvocation("GetMtaOperationLogContent", []interface{}{operationID, logID})
	fake.getMtaOperationLogContentMutex.Unlock()
	if fake.GetMtaOperationLogContentStub != nil {
		return fake.GetMtaOperationLogContentStub(operationID, logID)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fake.getMtaOperationLogContentReturns.result1, fake.getMtaOperationLogContentReturns.result2
}

func (fake *FakeMtaClientOperations) GetMtaOperationLogContentCallCount() int {
	fake.getMtaOperationLogContentMutex.RLock()
	defer fake.getMtaOperationLogContentMutex.RUnlock()
	return len(fake.getMtaOperationLogContentArgsForCall)
}

func (fake *FakeMtaClientOperations) GetMtaOperationLogContentArgsForCall(i int) (string, string) {
	fake.getMtaOperationLogContentMutex.RLock()
	defer fake.getMtaOperationLogContentMutex.RUnlock()
	return fake.getMtaOperationLogContentArgsForCall[i].operationID, fake.getMtaOperationLogContentArgsForCall[i].logID
}

func (fake *FakeMtaClientOperations) GetMtaOperationLogContentReturns(result1 string, result2 error) {
	fake.GetMtaOperationLogContentStub = nil
	fake.getMtaOperationLogContentReturns = struct {
		result1 string
		result2 error
	}{result1, result2}
}

func (fake *FakeMtaClientOperations) GetMtaOperationLogContentReturnsOnCall(i int, result1 string, result2 error) {
	fake.GetMtaOperationLogContentStub = nil
	if fake.getMtaOperationLogContentReturnsOnCall == nil {
		fake.getMtaOperationLogContentReturnsOnCall = make(map[int]struct {
			result1 string
			result2 error
		})
	}
	fake.getMtaOperationLogContentReturnsOnCall[i] = struct {
		result1 string
		result2 error
	}{result1, result2}
}

func (fake *FakeMtaClientOperations) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.executeActionMutex.RLock()
	defer fake.executeActionMutex.RUnlock()
	fake.getMtaMutex.RLock()
	defer fake.getMtaMutex.RUnlock()
	fake.getMtaFilesMutex.RLock()
	defer fake.getMtaFilesMutex.RUnlock()
	fake.getMtaOperationMutex.RLock()
	defer fake.getMtaOperationMutex.RUnlock()
	fake.getMtaOperationLogsMutex.RLock()
	defer fake.getMtaOperationLogsMutex.RUnlock()
	fake.getMtaOperationsMutex.RLock()
	defer fake.getMtaOperationsMutex.RUnlock()
	fake.getMtasMutex.RLock()
	defer fake.getMtasMutex.RUnlock()
	fake.getOperationActionsMutex.RLock()
	defer fake.getOperationActionsMutex.RUnlock()
	fake.startMtaOperationMutex.RLock()
	defer fake.startMtaOperationMutex.RUnlock()
	fake.uploadMtaFileMutex.RLock()
	defer fake.uploadMtaFileMutex.RUnlock()
	fake.getMtaOperationLogContentMutex.RLock()
	defer fake.getMtaOperationLogContentMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakeMtaClientOperations) recordInvocation(key string, args []interface{}) {
	fake.invocationsMutex.Lock()
	defer fake.invocationsMutex.Unlock()
	if fake.invocations == nil {
		fake.invocations = map[string][][]interface{}{}
	}
	if fake.invocations[key] == nil {
		fake.invocations[key] = [][]interface{}{}
	}
	fake.invocations[key] = append(fake.invocations[key], args)
}

var _ mtaclient.MtaClientOperations = new(FakeMtaClientOperations)
