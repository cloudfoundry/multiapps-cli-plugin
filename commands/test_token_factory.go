package commands

import (
	"github.com/cloudfoundry/cli/plugin/fakes"
	"github.com/go-openapi/runtime"
	"github.com/SAP/cf-mta-plugin/testutil"
)

type TestTokenFactory struct {
	FakeCliConnection *fakes.FakeCliConnection
}

func NewTestTokenFactory(fakeCliConnection *fakes.FakeCliConnection) *TestTokenFactory {
	return &TestTokenFactory{
		FakeCliConnection: fakeCliConnection,
	}
}

func (f *TestTokenFactory) NewToken() (runtime.ClientAuthInfoWriter, error) {
	tokenString, _ := f.FakeCliConnection.AccessToken()
	return testutil.NewCustomBearerToken(tokenString), nil
}
