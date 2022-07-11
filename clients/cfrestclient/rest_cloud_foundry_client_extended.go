package cfrestclient

import (
	"code.cloudfoundry.org/cli/plugin"
	"code.cloudfoundry.org/jsonry"
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/clients/baseclient"
	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/clients/cfrestclient/operations"
	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/clients/models"
	"github.com/go-openapi/strfmt"
)

const cfBaseUrl = "v3/"
const defaultDomainsPathPattern = "/domains"

type CloudFoundryRestClient struct {
	baseclient.BaseClient
	httpCloudFoundryClient *CloudFoundryClientExtended
	cliConn                plugin.CliConnection
}

func NewCloudFoundryRestClient(cliConn plugin.CliConnection, rt http.RoundTripper, tokenFactory baseclient.TokenFactory) CloudFoundryOperationsExtended {
	t := baseclient.NewHTTPTransport(getApiEndpoint(cliConn), cfBaseUrl, rt)
	httpCloudFoundryClient := New(t, strfmt.Default)
	return &CloudFoundryRestClient{baseclient.BaseClient{TokenFactory: tokenFactory}, httpCloudFoundryClient, cliConn}
}

func getApiEndpoint(cliConnection plugin.CliConnection) string {
	api, err := cliConnection.ApiEndpoint()
	if err != nil {
		return ""
	}
	if strings.HasPrefix(api, "https://") {
		api = strings.Replace(api, "https://", "", -1)
	}
	return api
}

func (c CloudFoundryRestClient) GetSharedDomains() ([]models.Domain, error) {
	var result []models.Domain
	var response *models.CloudFoundryResponse
	var err error
	var cloudFoundryUrlElements CloudFoundryUrlElements
	pathPattern := defaultDomainsPathPattern
	for pathPattern != "" {
		response, err = c.getSharedDomainsInternal(cloudFoundryUrlElements)
		if err != nil {
			return []models.Domain{}, models.HttpResponseError{Underlying: err}
		}
		result = append(result, toSharedDomains(response)...)
		cloudFoundryUrlElements, err = getPathQueryElements(response)
		if err != nil {
			return []models.Domain{}, models.HttpResponseError{Underlying: baseclient.NewClientError(err)}
		}
		pathPattern = response.Pagination.Next.Href
	}

	return result, nil
}

func (c CloudFoundryRestClient) GetApplications(mtaId, spaceGuid string) ([]models.CloudFoundryApplication, error) {
	token, err := c.cliConn.AccessToken()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve access token: %s", err)
	}
	apiEndpoint, _ := c.cliConn.ApiEndpoint()
	mtaIdHash := md5.Sum([]byte(mtaId))
	mtaIdHashStr := hex.EncodeToString(mtaIdHash[:])

	getAppsUrl := fmt.Sprintf("%s/%sapps?label_selector=mta_id=%s&space_guids=%s", apiEndpoint, cfBaseUrl, mtaIdHashStr, spaceGuid)
	return getPaginatedResources[models.CloudFoundryApplication](getAppsUrl, token)
}

func (c CloudFoundryRestClient) GetAppProcessStatistics(appGuid string) ([]models.ApplicationProcessStatistics, error) {
	token, err := c.cliConn.AccessToken()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve access token: %s", err)
	}
	apiEndpoint, _ := c.cliConn.ApiEndpoint()

	getAppProcessStatsUrl := fmt.Sprintf("%s/%sapps/%s/processes/web/stats", apiEndpoint, cfBaseUrl, appGuid)
	body, err := executeRequest(getAppProcessStatsUrl, token)
	if err != nil {
		return nil, err
	}
	processStats, err := parseBody[models.AppProcessStatisticsResponse](body)
	if err != nil {
		return nil, err
	}
	return processStats.Resources, nil
}

func (c CloudFoundryRestClient) GetApplicationRoutes(appGuid string) ([]models.ApplicationRoute, error) {
	token, err := c.cliConn.AccessToken()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve access token: %s", err)
	}
	apiEndpoint, _ := c.cliConn.ApiEndpoint()

	getAppRoutesUrl := fmt.Sprintf("%s/%sapps/%s/routes", apiEndpoint, cfBaseUrl, appGuid)
	return getPaginatedResources[models.ApplicationRoute](getAppRoutesUrl, token)
}

func (c CloudFoundryRestClient) GetServiceInstances(mtaId, spaceGuid string) ([]models.CloudFoundryServiceInstance, error) {
	token, err := c.cliConn.AccessToken()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve access token: %s", err)
	}
	apiEndpoint, _ := c.cliConn.ApiEndpoint()
	mtaIdHash := md5.Sum([]byte(mtaId))
	mtaIdHashStr := hex.EncodeToString(mtaIdHash[:])

	getServicesUrl := fmt.Sprintf("%s/%sservice_instances?fields[service_plan]=guid,name,relationships.service_offering&fields[service_plan.service_offering]=guid,name&space_guids=%s&label_selector=mta_id=%s",
		apiEndpoint, cfBaseUrl, spaceGuid, mtaIdHashStr)
	return getPaginatedResourcesWithIncluded(getServicesUrl, token, buildServiceInstance)
}

func (c CloudFoundryRestClient) GetServiceBindings(serviceName string) ([]models.ServiceBinding, error) {
	token, err := c.cliConn.AccessToken()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve access token: %s", err)
	}
	apiEndpoint, _ := c.cliConn.ApiEndpoint()

	getServiceBindingsUrl := fmt.Sprintf("%s/%sservice_credential_bindings?type=app&include=app&service_instance_names=%s", apiEndpoint, cfBaseUrl, serviceName)
	return getPaginatedResourcesWithIncluded(getServiceBindingsUrl, token, buildServiceBinding)
}

func (c CloudFoundryRestClient) getSharedDomainsInternal(cloudFoundryUrlElements CloudFoundryUrlElements) (*models.CloudFoundryResponse, error) {
	params := &operations.GetSharedDomainsParams{
		Page:           cloudFoundryUrlElements.Page,
		ResultsPerPage: cloudFoundryUrlElements.ResultsPerPage,
		Context:        context.TODO(),
	}
	token, err := c.TokenFactory.NewToken()
	if err != nil {
		return nil, baseclient.NewClientError(err)
	}

	//TODO filter by domains that have an owning organization
	// and maybe validate that organization is the current one targeted
	result, err := c.httpCloudFoundryClient.Operations.GetSharedDomains(params, token)
	if err != nil {
		return nil, baseclient.NewClientError(err)
	}
	return result.Payload, nil
}

func getPathQueryElements(response *models.CloudFoundryResponse) (CloudFoundryUrlElements, error) {
	nextUrl, err := url.Parse(response.Pagination.Next.Href)
	if err != nil {
		return CloudFoundryUrlElements{}, fmt.Errorf("could not parse pagintaion.next for getting shared domains: %s", response.Pagination.Next.Href)
	}
	nextUrlQuery := nextUrl.Query()
	page := nextUrlQuery.Get("page")
	resultsPerPage := nextUrlQuery.Get("per_page")

	return CloudFoundryUrlElements{Page: &page, ResultsPerPage: &resultsPerPage}, nil

}

func toSharedDomains(response *models.CloudFoundryResponse) []models.Domain {
	var result []models.Domain
	for _, cloudResource := range response.Resources {
		result = append(result, models.NewDomain(cloudResource.Name, cloudResource.GUID))
	}
	return result
}

func getPaginatedResources[T any](url, token string) ([]T, error) {
	var result []T
	for url != "" {
		body, err := executeRequest(url, token)
		if err != nil {
			return nil, err
		}
		response, err := parseBody[models.PaginatedResponse[T]](body)
		if err != nil {
			return nil, err
		}

		for _, entity := range response.Resources {
			result = append(result, entity)
		}
		url = response.Pagination.NextPage
	}
	return result, nil
}

func getPaginatedResourcesWithIncluded[T any, Auxiliary any](url, token string, auxiliaryContentHandler func(T, Auxiliary) T) ([]T, error) {
	var result []T
	for url != "" {
		body, err := executeRequest(url, token)
		if err != nil {
			return nil, err
		}
		response, err := parseBody[models.PaginatedResponseWithIncluded[T, Auxiliary]](body)
		if err != nil {
			return nil, err
		}

		for _, entity := range response.Resources {
			result = append(result, auxiliaryContentHandler(entity, response.Included))
		}
		url = response.Pagination.NextPage
	}
	return result, nil
}

func executeRequest(url, token string) ([]byte, error) {
	req, _ := http.NewRequest(http.MethodGet, url, nil)
	req.Header.Add("Authorization", token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode/100 != 2 {
		return nil, models.HttpResponseError{Underlying: fmt.Errorf("%s: %s", resp.Status, string(bytes))}
	}
	return bytes, nil
}

func parseBody[T any](body []byte) (T, error) {
	var result T
	err := jsonry.Unmarshal(body, &result)
	if err == nil {
		return result, nil
	}
	//jsonry doesn't work with raw objects like map, so try the base json decoder
	err = json.Unmarshal(body, &result)
	if err == nil {
		return result, nil
	}
	return result, fmt.Errorf("could not parse response: %s", err)
}

func buildServiceInstance(service models.CloudFoundryServiceInstance, auxiliaryContent models.ServiceInstanceAuxiliaryContent) models.CloudFoundryServiceInstance {
	servicePlan := findServicePlan(service.PlanGuid, auxiliaryContent.ServicePlans)
	service.Plan = servicePlan
	service.Offering = findServiceOffering(servicePlan, auxiliaryContent.ServiceOfferings)
	return service
}

func findServicePlan(planGuid string, plans []models.ServicePlan) models.ServicePlan {
	for _, plan := range plans {
		if plan.Guid == planGuid {
			return plan
		}
	}
	return models.ServicePlan{}
}

func findServiceOffering(plan models.ServicePlan, offerings []models.ServiceOffering) models.ServiceOffering {
	for _, offering := range offerings {
		if offering.Guid == plan.OfferingGuid {
			return offering
		}
	}
	return models.ServiceOffering{}
}

func buildServiceBinding(binding models.ServiceBinding, auxiliaryContent models.ServiceBindingAuxiliaryContent) models.ServiceBinding {
	binding.AppName = findApp(binding.AppGuid, auxiliaryContent.Apps).Name
	return binding
}

func findApp(appGuid string, apps []models.CloudFoundryApplication) models.CloudFoundryApplication {
	for _, app := range apps {
		if app.Guid == appGuid {
			return app
		}
	}
	return models.CloudFoundryApplication{}
}
