package cfrestclient

import (
	"bytes"
	"crypto/md5"
	"crypto/tls"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"code.cloudfoundry.org/cli/v8/plugin"
	"code.cloudfoundry.org/jsonry"
	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/clients/baseclient"
	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/clients/models"
	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/log"
)

const cfBaseUrl = "v3/"

type CloudFoundryRestClient struct {
	cliConn       plugin.CliConnection
	isSslDisabled bool
}

func NewCloudFoundryRestClient(cliConn plugin.CliConnection) CloudFoundryOperationsExtended {
	isSslDisabled, err := cliConn.IsSSLDisabled()
	if err != nil {
		log.Tracef("Error while determining skip-ssl-validation: %v", err)
		isSslDisabled = false
	}
	return &CloudFoundryRestClient{cliConn, isSslDisabled}
}

func (c CloudFoundryRestClient) GetApplications(mtaId, mtaNamespace, spaceGuid string) ([]models.CloudFoundryApplication, error) {
	token, err := c.cliConn.AccessToken()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve access token: %s", err)
	}
	apiEndpoint, _ := c.cliConn.ApiEndpoint()
	mtaIdHash := md5.Sum([]byte(mtaId))
	mtaIdHashStr := hex.EncodeToString(mtaIdHash[:])

	getAppsUrl := fmt.Sprintf("%s/%sapps?space_guids=%s&label_selector=mta_id=%s", apiEndpoint, cfBaseUrl, spaceGuid, mtaIdHashStr)
	if mtaNamespace != "" {
		namespaceHash := md5.Sum([]byte(mtaNamespace))
		namespaceHashStr := hex.EncodeToString(namespaceHash[:])
		getAppsUrl = fmt.Sprintf("%s,mta_namespace=%s", getAppsUrl, namespaceHashStr)
	} else {
		getAppsUrl = fmt.Sprintf("%s,!mta_namespace", getAppsUrl)
	}
	return getPaginatedResources[models.CloudFoundryApplication](getAppsUrl, token, c.isSslDisabled)
}

func (c CloudFoundryRestClient) GetAppProcessStatistics(appGuid string) ([]models.ApplicationProcessStatistics, error) {
	token, err := c.cliConn.AccessToken()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve access token: %s", err)
	}
	apiEndpoint, _ := c.cliConn.ApiEndpoint()

	getAppProcessStatsUrl := fmt.Sprintf("%s/%sapps/%s/processes/web/stats", apiEndpoint, cfBaseUrl, appGuid)
	body, err := executeRequest(http.MethodGet, getAppProcessStatsUrl, token, c.isSslDisabled, nil)
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
	return getPaginatedResources[models.ApplicationRoute](getAppRoutesUrl, token, c.isSslDisabled)
}

func (c CloudFoundryRestClient) GetServiceInstances(mtaId, mtaNamespace, spaceGuid string) ([]models.CloudFoundryServiceInstance, error) {
	token, err := c.cliConn.AccessToken()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve access token: %s", err)
	}
	apiEndpoint, _ := c.cliConn.ApiEndpoint()
	mtaIdHash := md5.Sum([]byte(mtaId))
	mtaIdHashStr := hex.EncodeToString(mtaIdHash[:])

	getServicesUrl := fmt.Sprintf("%s/%sservice_instances?fields[service_plan]=guid,name,relationships.service_offering&fields[service_plan.service_offering]=guid,name&space_guids=%s&label_selector=mta_id=%s",
		apiEndpoint, cfBaseUrl, spaceGuid, mtaIdHashStr)
	if mtaNamespace != "" {
		namespaceHash := md5.Sum([]byte(mtaNamespace))
		namespaceHashStr := hex.EncodeToString(namespaceHash[:])
		getServicesUrl = fmt.Sprintf("%s,mta_namespace=%s", getServicesUrl, namespaceHashStr)
	} else {
		getServicesUrl = fmt.Sprintf("%s,!mta_namespace", getServicesUrl)
	}
	return getPaginatedResourcesWithIncluded(getServicesUrl, token, c.isSslDisabled, buildServiceInstance)
}

func (c CloudFoundryRestClient) GetServiceBindings(serviceName string) ([]models.ServiceBinding, error) {
	token, err := c.cliConn.AccessToken()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve access token: %s", err)
	}
	apiEndpoint, _ := c.cliConn.ApiEndpoint()

	getServiceBindingsUrl := fmt.Sprintf("%s/%sservice_credential_bindings?type=app&include=app&service_instance_names=%s", apiEndpoint, cfBaseUrl, serviceName)
	return getPaginatedResourcesWithIncluded(getServiceBindingsUrl, token, c.isSslDisabled, buildServiceBinding)
}

func (c CloudFoundryRestClient) GetServiceInstanceByName(serviceName, spaceGuid string) (models.CloudFoundryServiceInstance, error) {
	token, err := c.cliConn.AccessToken()
	if err != nil {
		return models.CloudFoundryServiceInstance{}, fmt.Errorf("failed to retrieve access token: %s", err)
	}
	apiEndpoint, _ := c.cliConn.ApiEndpoint()

	getServicesUrl := fmt.Sprintf("%s/%sservice_instances?names=%s&space_guids=%s",
		apiEndpoint, cfBaseUrl, serviceName, spaceGuid)
	services, err := getPaginatedResourcesWithIncluded(getServicesUrl, token, c.isSslDisabled, buildServiceInstance)
	if err != nil {
		return models.CloudFoundryServiceInstance{}, err
	}
	if len(services) == 0 {
		return models.CloudFoundryServiceInstance{}, fmt.Errorf("service instance not found")
	}

	resultService := services[0]
	return resultService, nil
}

func (c CloudFoundryRestClient) CreateUserProvidedServiceInstance(serviceName string, spaceGuid string, credentials map[string]string) (models.CloudFoundryServiceInstance, error) {
	token, err := c.cliConn.AccessToken()
	if err != nil {
		return models.CloudFoundryServiceInstance{}, fmt.Errorf("failed to retrieve access token: %s", err)
	}

	apiEndpoint, _ := c.cliConn.ApiEndpoint()

	createServiceURL := fmt.Sprintf("%s/%sservice_instances", apiEndpoint, cfBaseUrl)

	payload := map[string]any{
		"type": "user-provided",
		"name": serviceName,
		"relationships": map[string]any{
			"space": map[string]any{
				"data": map[string]any{
					"guid": spaceGuid,
				},
			},
		},
	}

	if credentials != nil {
		payload["credentials"] = credentials
	}

	jsonBody, err := json.Marshal(payload)
	if err != nil {
		return models.CloudFoundryServiceInstance{}, fmt.Errorf("failed to marshal create UPS request: %w", err)
	}

	body, err := executeRequest(http.MethodPost, createServiceURL, token, c.isSslDisabled, jsonBody)
	if err != nil {
		return models.CloudFoundryServiceInstance{}, err
	}

	response, err := parseBody[models.CloudFoundryUserProvidedServiceInstance](body)
	if err != nil {
		return models.CloudFoundryServiceInstance{}, err
	}

	return models.CloudFoundryServiceInstance{
		Guid:          response.Guid,
		Name:          response.Name,
		Type:          response.Type,
		LastOperation: response.LastOperation,
		SpaceGuid:     response.SpaceGuid,
		Metadata:      response.Metadata,
	}, nil
}

func getPaginatedResources[T any](url, token string, isSslDisabled bool) ([]T, error) {
	var result []T
	for url != "" {
		body, err := executeRequest(http.MethodGet, url, token, isSslDisabled, nil)
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

func getPaginatedResourcesWithIncluded[T any, Auxiliary any](url, token string, isSslDisabled bool, auxiliaryContentHandler func(T, Auxiliary) T) ([]T, error) {
	var result []T
	for url != "" {
		body, err := executeRequest(http.MethodGet, url, token, isSslDisabled, nil)
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

func executeRequest(methodType, url, token string, isSslDisabled bool, body []byte) ([]byte, error) {
	var reader io.Reader

	if body != nil {
		reader = bytes.NewReader(body)
	}
	request, err := http.NewRequest(methodType, url, reader)
	if err != nil {
		return nil, err
	}

	request.Header.Add("Authorization", token)
	request.Header.Set("Accept", "application/json")
	if body != nil {
		request.Header.Set("Content-Type", "application/json")
	}

	// Create transport with TLS configuration
	httpTransport := http.DefaultTransport.(*http.Transport).Clone()
	httpTransport.TLSClientConfig = &tls.Config{InsecureSkipVerify: isSslDisabled}

	// Wrap with User-Agent transport
	userAgentTransport := baseclient.NewUserAgentTransport(httpTransport)

	client := &http.Client{Transport: userAgentTransport}
	resp, err := client.Do(request)
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
