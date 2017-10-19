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

//NewBaseClient creates a new client with the specified authentication info
func NewBaseClient(tokenFactory TokenFactory) *BaseClient {
	return &BaseClient{
		TokenFactory: tokenFactory,
	}
}

// NewHTTPTransport creates a new HTTP transport
func NewHTTPTransport(host, url, encodedUrl string, rt http.RoundTripper, jar http.CookieJar) *client.Runtime {
	transport := client.New(host, url, encodedUrl, schemes)
	transport.Consumers["text/html"] = runtime.TextConsumer()
	transport.DefaultMediaType = "application/json" //TODO support XML for the config registry
	transport.Transport = rt
	transport.Jar = jar
	return transport
}
