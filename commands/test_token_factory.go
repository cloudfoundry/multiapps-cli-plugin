package commands

import (
	fakes "code.cloudfoundry.org/cli/v8/plugin/pluginfakes"
	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/testutil"
	"github.com/go-openapi/runtime"
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
	tokenString, _ := f.NewRawToken()
	return testutil.NewCustomBearerToken(tokenString), nil
}

func (f *TestTokenFactory) NewRawToken() (string, error) {
	return f.FakeCliConnection.AccessToken()
}
