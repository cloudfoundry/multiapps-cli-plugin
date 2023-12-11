package cfrestclient

import (
	"crypto/md5"
	"crypto/tls"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"code.cloudfoundry.org/cli/plugin"
	"code.cloudfoundry.org/jsonry"
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
	body, err := executeRequest(getAppProcessStatsUrl, token, c.isSslDisabled)
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

func getPaginatedResources[T any](url, token string, isSslDisabled bool) ([]T, error) {
	var result []T
	for url != "" {
		body, err := executeRequest(url, token, isSslDisabled)
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
		body, err := executeRequest(url, token, isSslDisabled)
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

func executeRequest(url, token string, isSslDisabled bool) ([]byte, error) {
	req, _ := http.NewRequest(http.MethodGet, url, nil)
	req.Header.Add("Authorization", token)
	httpTransport := http.DefaultTransport.(*http.Transport).Clone()
	httpTransport.TLSClientConfig = &tls.Config{InsecureSkipVerify: isSslDisabled}
	client := http.DefaultClient
	client.Transport = httpTransport
	resp, err := client.Do(req)
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
