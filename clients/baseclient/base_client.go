package baseclient

import (
	"net/http"
	"net/http/cookiejar"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/client"
)

var schemes = []string{"http", "https"}

// BaseClient represents a base SLP client
type BaseClient struct {
	TokenFactory TokenFactory
}

// GetTokenFactory returns the authentication info
func (c *BaseClient) GetTokenFactory() TokenFactory {
	return c.TokenFactory
}

// SetTokenFactory sets the authentication info
func (c *BaseClient) SetTokenFactory(tokenFactory TokenFactory) {
	c.TokenFactory = tokenFactory
}

// NewHTTPTransport creates a new HTTP transport
func NewHTTPTransport(host, url string, rt http.RoundTripper) *client.Runtime {
	// TODO: apply the changes made by Boyan here, as after the update of the dependencies the changes are not available
	transport := client.New(host, url, schemes)
	transport.Consumers["text/html"] = runtime.TextConsumer()

	// Wrap the RoundTripper with User-Agent support
	userAgentTransport := NewUserAgentTransport(rt)
	transport.Transport = userAgentTransport

	jar, _ := cookiejar.New(nil)
	transport.Jar = jar
	return transport
}
