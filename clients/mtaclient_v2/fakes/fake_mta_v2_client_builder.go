package fakes

import (
	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/clients/models"
)

type FakeMtaV2ClientBuilder struct {
	FakeMtaV2Client FakeMtaV2ClientOperations
}

func NewFakeMtaV2ClientBuilder() *FakeMtaV2ClientBuilder {
	return &FakeMtaV2ClientBuilder{}
}

func (fb *FakeMtaV2ClientBuilder) GetMtas(mtaID string, namespace *string, spaceGuid string, resultMta []*models.Mta, resultError error) *FakeMtaV2ClientBuilder {
	fb.FakeMtaV2Client.GetMtasReturns(resultMta, resultError)
	return fb
}

func (fb *FakeMtaV2ClientBuilder) GetMtasForThisSpace(mtaID string, namespace *string, resultMta []*models.Mta, resultError error) *FakeMtaV2ClientBuilder {
	fb.FakeMtaV2Client.GetMtasForThisSpaceReturns(resultMta, resultError)
	return fb
}

func (fb *FakeMtaV2ClientBuilder) Build() FakeMtaV2ClientOperations {
	return fb.FakeMtaV2Client
}
