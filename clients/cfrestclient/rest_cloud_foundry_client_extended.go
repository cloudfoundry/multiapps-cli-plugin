package cfrestclient

import (
	"context"
	"net/http"

	baseclient "github.com/SAP/cf-mta-plugin/clients/baseclient"
	"github.com/SAP/cf-mta-plugin/clients/cfrestclient/operations"
	models "github.com/SAP/cf-mta-plugin/clients/models"
	strfmt "github.com/go-openapi/strfmt"
)

const cfBaseUrl = "v2/"

type CloudFoundryRestClient struct {
	baseclient.BaseClient
	httpCloudFoundryClient *CloudFoundryClientExtended
}

func NewCloudFoundryRestClient(host string, rt http.RoundTripper, jar http.CookieJar, tokenFactory baseclient.TokenFactory) CloudFoundryOperationsExtended {
	t := baseclient.NewHTTPTransport(host, cfBaseUrl, cfBaseUrl, rt, jar)
	httpCloudFoundryClient := New(t, strfmt.Default)
	return &CloudFoundryRestClient{baseclient.BaseClient{TokenFactory: tokenFactory}, httpCloudFoundryClient}
}

func (c CloudFoundryRestClient) GetSharedDomains() ([]models.SharedDomain, error) {
	params := &operations.GetSharedDomainsParams{
		Context: context.TODO(),
	}
	token, err := c.TokenFactory.NewToken()
	if err != nil {
		return []models.SharedDomain{}, baseclient.NewClientError(err)
	}

	result, err := c.httpCloudFoundryClient.Operations.GetSharedDomains(params, token)
	if err != nil {
		return []models.SharedDomain{}, baseclient.NewClientError(err)
	}

	return toSharedDomains(result.Payload), nil
}

func toSharedDomains(response *models.CloudFoundryResponse) []models.SharedDomain {
	var result []models.SharedDomain
	for _, cloudResource := range response.Resources {
		result = append(result, models.NewSharedDomain(cloudResource.Entity.Name, cloudResource.Metadata.GUID, cloudResource.Metadata.URL))
	}
	return result
}
