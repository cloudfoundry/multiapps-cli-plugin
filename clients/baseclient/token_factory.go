package baseclient

import "github.com/go-openapi/runtime"

// TokenFactory factory for generating new OAuth token
type TokenFactory interface {
	NewToken() (runtime.ClientAuthInfoWriter, error)
}
