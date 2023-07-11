package fakes

import (
	"github.com/cloudfoundry/multiapps-cli-plugin/clients/models"
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

// Build builds a FakeRestClientOperations instance
func (b *FakeRestClientBuilder) Build() *FakeRestClientOperations {
	return &b.fakeRestClient
}
