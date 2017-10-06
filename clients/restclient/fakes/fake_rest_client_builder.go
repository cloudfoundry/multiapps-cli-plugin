package fakes

import (
	"strconv"

	"github.com/SAP/cf-mta-plugin/clients/models"
)

type mtaResult struct {
	*models.Mta
	error
}

// FakeRestClientBuilder is a builder of FakeRestClientOperations instances
type FakeRestClientBuilder struct {
	fakeRestClient FakeRestClientOperations
	mtaResults     map[string]mtaResult
}

// NewFakeRestClientBuilder creates a new builder
func NewFakeRestClientBuilder() *FakeRestClientBuilder {
	return &FakeRestClientBuilder{}
}

// GetOperations sets the operations to return from the GetOperations operation
func (b *FakeRestClientBuilder) GetOperations(lastRequestedOperations *string, requestedStates []string, operations models.Operations, err error) *FakeRestClientBuilder {
	ops := operations
	if lastRequestedOperations != nil {
		ops = getLastOperations(*lastRequestedOperations, operations)
	} else if requestedStates != nil && len(requestedStates) > 0 {
		ops = getOperationsByStates(requestedStates, operations)
	}
	b.fakeRestClient.GetOperationsReturns(ops, err)
	return b
}

func getOperationsByStates(requestedStates []string, operations models.Operations) models.Operations {
	var ops models.Operations
	requestedStatesMap := getRequestedStatesMap(requestedStates)
	for _, operation := range operations.Operations {
		if isRequested(&operation.State, requestedStatesMap) {
			ops.Operations = append(ops.Operations, operation)
		}
	}
	return ops
}

func getRequestedStatesMap(requestedStates []string) map[string]bool {
	var requestedStatesMap = make(map[string]bool)
	for _, requestedState := range requestedStates {
		requestedStatesMap[requestedState] = true
	}
	return requestedStatesMap
}

func isRequested(state *models.SlpTaskStateEnum, requestedStates map[string]bool) bool {
	return requestedStates[string(*state)]
}

func getLastOperations(lastRequestedOperationsCount string, operations models.Operations) models.Operations {
	operationsCount, _ := strconv.Atoi(lastRequestedOperationsCount)
	if operationsCount > len(operations.Operations) {
		return operations
	}
	operationsLength := len(operations.Operations)
	return models.Operations{Operations: operations.Operations[operationsLength-operationsCount:]}
}

// GetComponents sets the components to return from the GetComponents operations
func (b *FakeRestClientBuilder) GetComponents(components *models.Components, err error) *FakeRestClientBuilder {
	b.fakeRestClient.GetComponentsReturns(components, err)
	return b
}

// GetMta sets the MTA to return from the GetMta operation
func (b *FakeRestClientBuilder) GetMta(mtaID string, mta *models.Mta, err error) *FakeRestClientBuilder {
	if mtaID == "" {
		b.fakeRestClient.GetMtaReturns(mta, err)
		return b
	}
	if b.mtaResults == nil {
		b.mtaResults = make(map[string]mtaResult)
	}
	b.mtaResults[mtaID] = mtaResult{mta, err}
	if b.fakeRestClient.GetMtaStub == nil {
		b.fakeRestClient.GetMtaStub = func(arg0 string) (*models.Mta, error) {
			result := b.mtaResults[arg0]
			return result.Mta, result.error
		}
	}
	return b
}

// Build builds a FakeRestClientOperations instance
func (b *FakeRestClientBuilder) Build() *FakeRestClientOperations {
	return &b.fakeRestClient
}
