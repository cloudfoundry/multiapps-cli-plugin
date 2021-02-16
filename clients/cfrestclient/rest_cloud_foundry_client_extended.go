package cfrestclient

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	baseclient "github.com/cloudfoundry-incubator/multiapps-cli-plugin/clients/baseclient"
	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/clients/cfrestclient/operations"
	models "github.com/cloudfoundry-incubator/multiapps-cli-plugin/clients/models"
	strfmt "github.com/go-openapi/strfmt"
)

const cfBaseUrl = "v2/"
const defaultSharedDomainsPathPattern = "/shared_domains"

type CloudFoundryRestClient struct {
	baseclient.BaseClient
	httpCloudFoundryClient *CloudFoundryClientExtended
}

func NewCloudFoundryRestClient(host string, rt http.RoundTripper, tokenFactory baseclient.TokenFactory) CloudFoundryOperationsExtended {
	t := baseclient.NewHTTPTransport(host, cfBaseUrl, cfBaseUrl, rt)
	httpCloudFoundryClient := New(t, strfmt.Default)
	return &CloudFoundryRestClient{baseclient.BaseClient{TokenFactory: tokenFactory}, httpCloudFoundryClient}
}

func (c CloudFoundryRestClient) GetSharedDomains() ([]models.SharedDomain, error) {
	result := []models.SharedDomain{}
	var response *models.CloudFoundryResponse
	var err error
	var cloudFoundryUrlElememnts CloudFoundryUrlElements
	pathPattern := defaultSharedDomainsPathPattern
	for pathPattern != "" {
		response, err = c.getSharedDomainsInternal(cloudFoundryUrlElememnts)
		if err != nil {
			return []models.SharedDomain{}, err
		}
		result = append(result, toSharedDomains(response)...)
		cloudFoundryUrlElememnts, err = getPathQueryElements(response)
		if err != nil {
			return []models.SharedDomain{}, baseclient.NewClientError(err)
		}
		pathPattern = response.NextURL
	}

	return result, err
}

func (c CloudFoundryRestClient) getSharedDomainsInternal(cloudFoundryUrlElements CloudFoundryUrlElements) (*models.CloudFoundryResponse, error) {
	params := &operations.GetSharedDomainsParams{
		Page:           cloudFoundryUrlElements.Page,
		ResultsPerPage: cloudFoundryUrlElements.ResultsPerPage,
		OrderDirection: cloudFoundryUrlElements.OrderDirection,
		Context:        context.TODO(),
	}
	token, err := c.TokenFactory.NewToken()
	if err != nil {
		return nil, baseclient.NewClientError(err)
	}

	result, err := c.httpCloudFoundryClient.Operations.GetSharedDomains(params, token)
	if err != nil {
		return nil, baseclient.NewClientError(err)
	}
	return result.Payload, nil
}

func getPathQueryElements(response *models.CloudFoundryResponse) (CloudFoundryUrlElements, error) {
	nextUrl, err := url.Parse(response.NextURL)
	if err != nil {
		return CloudFoundryUrlElements{}, fmt.Errorf("Could not parse next_url for getting shared domains: %s", response.NextURL)
	}
	nextUrlQuery := nextUrl.Query()
	page := nextUrlQuery.Get("page")
	resultsPerPage := nextUrlQuery.Get("results-per-page")
	orderDirection := nextUrlQuery.Get("order-direction")

	return CloudFoundryUrlElements{Page: &page, ResultsPerPage: &resultsPerPage, OrderDirection: &orderDirection}, nil

}

func toSharedDomains(response *models.CloudFoundryResponse) []models.SharedDomain {
	var result []models.SharedDomain
	for _, cloudResource := range response.Resources {
		result = append(result, models.NewSharedDomain(cloudResource.Entity.Name, cloudResource.Metadata.GUID, cloudResource.Metadata.URL))
	}
	return result
}
