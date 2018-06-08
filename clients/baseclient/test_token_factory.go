package baseclient

import (
	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/testutil"
	"github.com/go-openapi/runtime"
)

// NewCustomTokenFactory represents mock of the TokenFactory
func NewCustomTokenFactory(token string) TokenFactory {
	return &customTokenfactory{tokenString: token}
}

type customTokenfactory struct {
	tokenString string
}

func (c *customTokenfactory) NewToken() (runtime.ClientAuthInfoWriter, error) {
	return testutil.NewCustomBearerToken(c.tokenString), nil
}
