package fakes

import (
	"github.com/cloudfoundry/cli/plugin/fakes"
	plugin_models "github.com/cloudfoundry/cli/plugin/models"
)

// FakeCliConnectionBuilder is a builder of FakeCliConnection instances
type FakeCliConnectionBuilder struct {
	cliConn fakes.FakeCliConnection
}

// NewFakeCliConnectionBuilder creates a new builder
func NewFakeCliConnectionBuilder() *FakeCliConnectionBuilder {
	return &FakeCliConnectionBuilder{}
}

func (b *FakeCliConnectionBuilder) CurrentOrg(guid, name string, err error) *FakeCliConnectionBuilder {
	org := plugin_models.Organization{OrganizationFields: plugin_models.OrganizationFields{Guid: guid, Name: name}}
	b.cliConn.GetCurrentOrgReturns(org, err)
	return b
}

func (b *FakeCliConnectionBuilder) CurrentSpace(guid, name string, err error) *FakeCliConnectionBuilder {
	space := plugin_models.Space{SpaceFields: plugin_models.SpaceFields{Guid: guid, Name: name}}
	b.cliConn.GetCurrentSpaceReturns(space, err)
	return b
}

func (b *FakeCliConnectionBuilder) Username(username string, err error) *FakeCliConnectionBuilder {
	b.cliConn.UsernameReturns(username, err)
	return b
}

func (b *FakeCliConnectionBuilder) AccessToken(token string, err error) *FakeCliConnectionBuilder {
	b.cliConn.AccessTokenReturns(token, err)
	return b
}

func (b *FakeCliConnectionBuilder) APIEndpoint(apiURL string, err error) *FakeCliConnectionBuilder {
	b.cliConn.ApiEndpointReturns(apiURL, err)
	return b
}

func (b *FakeCliConnectionBuilder) GetApp(name string, app plugin_models.GetAppModel, err error) *FakeCliConnectionBuilder {
	b.cliConn.GetAppReturns(app, err) // TODO
	return b
}

func (b *FakeCliConnectionBuilder) GetApps(apps []plugin_models.GetAppsModel, err error) *FakeCliConnectionBuilder {
	b.cliConn.GetAppsReturns(apps, err)
	return b
}

func (b *FakeCliConnectionBuilder) GetService(name string, service plugin_models.GetService_Model, err error) *FakeCliConnectionBuilder {
	b.cliConn.GetServiceReturns(service, err) // TODO
	return b
}

func (b *FakeCliConnectionBuilder) GetServices(services []plugin_models.GetServices_Model, err error) *FakeCliConnectionBuilder {
	b.cliConn.GetServicesReturns(services, err)
	return b
}

func (b *FakeCliConnectionBuilder) GetSpace(name string, space plugin_models.GetSpace_Model, err error) *FakeCliConnectionBuilder {
	b.cliConn.GetSpaceReturns(space, err) // TODO
	return b
}

// Build builds a FakeCliConnection instance
func (b *FakeCliConnectionBuilder) Build() *fakes.FakeCliConnection {
	return &b.cliConn
}
