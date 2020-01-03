package baseclient

import (
	"net/http"

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
func NewHTTPTransport(host, url, encodedUrl string, rt http.RoundTripper, jar http.CookieJar) *client.Runtime {
	// TODO: apply the changes made by Boyan here, as after the update of the dependencies the changes are not available
	transport := client.New(host, url, encodedUrl, schemes)
	transport.Consumers["text/html"] = runtime.TextConsumer()
	transport.Transport = rt
	transport.Jar = jar
	return transport
}
