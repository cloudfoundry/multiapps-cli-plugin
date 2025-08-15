package baseclient

import (
	"net/http"

	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/util"
)

// UserAgentTransport wraps an existing RoundTripper and adds User-Agent header
type UserAgentTransport struct {
	Base      http.RoundTripper
	UserAgent string
}

// NewUserAgentTransport creates a new transport with User-Agent header support
func NewUserAgentTransport(base http.RoundTripper) *UserAgentTransport {
	if base == nil {
		base = http.DefaultTransport
	}

	return &UserAgentTransport{
		Base:      base,
		UserAgent: util.BuildUserAgent(),
	}
}

// RoundTrip implements the RoundTripper interface
func (uat *UserAgentTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	// Clone the request to avoid modifying the original
	reqCopy := req.Clone(req.Context())

	// Add or override the User-Agent header
	reqCopy.Header.Set("User-Agent", uat.UserAgent)

	// Execute the request with the base transport
	return uat.Base.RoundTrip(reqCopy)
}
