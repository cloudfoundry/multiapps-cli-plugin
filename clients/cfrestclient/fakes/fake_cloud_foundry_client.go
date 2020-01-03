package fakes

import (
	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/clients/cfrestclient"
	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/clients/models"
)

// TODO: use counterfeiter if the client becomes more sophisticated

type FakeCloudFoundryClient struct {
	domains []models.SharedDomain
	err     error
}

func NewFakeCloudFoundryClient(domains []models.SharedDomain, err error) cfrestclient.CloudFoundryOperationsExtended {
	return FakeCloudFoundryClient{domains: domains, err: err}
}

func (f FakeCloudFoundryClient) GetSharedDomains() ([]models.SharedDomain, error) {
	return f.domains, f.err
}
