package testutil

import (
	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/clients/baseclient"
	"github.com/go-openapi/runtime"
)

// NewCustomTokenFactory represents mock of the TokenFactory
func NewCustomTokenFactory(token string) baseclient.TokenFactory {
	return &customTokenFactory{tokenString: token}
}

type customTokenFactory struct {
	tokenString string
}

func (c *customTokenFactory) NewToken() (runtime.ClientAuthInfoWriter, error) {
	return NewCustomBearerToken(c.tokenString), nil
}

func (c *customTokenFactory) NewRawToken() (string, error) {
	return c.tokenString, nil
}
