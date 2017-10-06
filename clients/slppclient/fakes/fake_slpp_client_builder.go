package fakes

import models "github.com/SAP/cf-mta-plugin/clients/models"

type logContentResult struct {
	string
	error
}

// FakeSlppClientBuilder is a builder of FakeSlppClientOperations instances
type FakeSlppClientBuilder struct {
	fakeSlppClient       FakeSlppClientOperations
	getLogContentResults map[string]logContentResult
	executeActionResults map[string]error
}

// NewFakeSlppClientBuilder creates a new builder
func NewFakeSlppClientBuilder() *FakeSlppClientBuilder {
	return &FakeSlppClientBuilder{}
}

// GetMetadata sets the metadata to return from the GetMetadata operation
func (b *FakeSlppClientBuilder) GetMetadata(metadata *models.Metadata, err error) *FakeSlppClientBuilder {
	b.fakeSlppClient.GetMetadataReturns(metadata, err)
	return b
}

// GetLogs sets the logs to return from the GetLogs operation
func (b *FakeSlppClientBuilder) GetLogs(logs models.Logs, err error) *FakeSlppClientBuilder {
	b.fakeSlppClient.GetLogsReturns(logs, err)
	return b
}

// GetLogContent sets the log content to return from the GetLogContent operation
func (b *FakeSlppClientBuilder) GetLogContent(logID string, logContent string, err error) *FakeSlppClientBuilder {
	if logID == "" {
		b.fakeSlppClient.GetLogContentReturns(logContent, err)
		return b
	}
	if b.getLogContentResults == nil {
		b.getLogContentResults = make(map[string]logContentResult)
	}
	b.getLogContentResults[logID] = logContentResult{logContent, err}
	if b.fakeSlppClient.GetLogContentStub == nil {
		b.fakeSlppClient.GetLogContentStub = func(arg0 string) (string, error) {
			result := b.getLogContentResults[arg0]
			return result.string, result.error
		}
	}
	return b
}

// GetTasklist sets the tasklist to return from the GetTasklist operation
func (b *FakeSlppClientBuilder) GetTasklist(tasklist models.Tasklist, err error) *FakeSlppClientBuilder {
	b.fakeSlppClient.GetTasklistReturns(tasklist, err)
	return b
}

// GetTasklistTask sets the task to return from the GetTasklistTask operation
func (b *FakeSlppClientBuilder) GetTasklistTask(task *models.Task, err error) *FakeSlppClientBuilder {
	b.fakeSlppClient.GetTasklistTaskReturns(task, err)
	return b
}

// GetServiceID sets the serviceID to return from the GetServiceID operation
func (b *FakeSlppClientBuilder) GetServiceID(serviceID string) *FakeSlppClientBuilder {
	b.fakeSlppClient.GetServiceIDReturns(serviceID)
	return b
}

// GetError sets the error to return from the GetError operation
func (b *FakeSlppClientBuilder) GetError(errorx *models.Error, err error) *FakeSlppClientBuilder {
	b.fakeSlppClient.GetErrorReturns(errorx, err)
	return b
}

// ExecuteAction sets the result to return from the ExecuteAction operation
func (b *FakeSlppClientBuilder) ExecuteAction(actionID string, err error) *FakeSlppClientBuilder {
	if actionID == "" {
		b.fakeSlppClient.ExecuteActionReturns(err)
		return b
	}
	if b.executeActionResults == nil {
		b.executeActionResults = make(map[string]error)
	}
	b.executeActionResults[actionID] = err
	if b.fakeSlppClient.ExecuteActionStub == nil {
		b.fakeSlppClient.ExecuteActionStub = func(arg0 string) error {
			return b.executeActionResults[arg0]
		}
	}
	return b
}

// GetActions sets the actions to return from the GetActions operation
func (b *FakeSlppClientBuilder) GetActions(actions models.Actions, err error) *FakeSlppClientBuilder {
	b.fakeSlppClient.GetActionsReturns(actions, err)
	return b
}

// Build builds a FakeSlppClientOperations instance
func (b *FakeSlppClientBuilder) Build() *FakeSlppClientOperations {
	return &b.fakeSlppClient
}
