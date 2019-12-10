package mtaclient_v2

import (
	"net/http"
	"time"

	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/clients/baseclient"
	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/clients/models"
)

type RetryableMtaRestClient struct {
	mtaClient       MtaV2ClientOperations
	MaxRetriesCount int
	RetryInterval   time.Duration
}

func NewRetryableMtaRestClient(host string, spaceGUID string, rt http.RoundTripper, jar http.CookieJar, tokenFactory baseclient.TokenFactory) RetryableMtaRestClient {
	mtaClient := NewMtaClient(host, spaceGUID, rt, jar, tokenFactory)
	return RetryableMtaRestClient{mtaClient: mtaClient, MaxRetriesCount: 3, RetryInterval: time.Second * 3}
}

func (c RetryableMtaRestClient) GetMtas(name, namespace *string, spaceGuid string) ([]*models.Mta, error) {
	getMtasCb := func() (interface{}, error) {
		return c.mtaClient.GetMtas(name, namespace, spaceGuid)
	}
	resp, err := baseclient.CallWithRetry(getMtasCb, c.MaxRetriesCount, c.RetryInterval)
	if err != nil {
		return nil, err
	}
	return resp.([]*models.Mta), nil
}

func (c RetryableMtaRestClient) GetMtasForThisSpace(name, namespace *string) ([]*models.Mta, error) {
	getMtasCb := func() (interface{}, error) {
		return c.mtaClient.GetMtasForThisSpace(name, namespace)
	}
	resp, err := baseclient.CallWithRetry(getMtasCb, c.MaxRetriesCount, c.RetryInterval)
	if err != nil {
		return nil, err
	}
	return resp.([]*models.Mta), nil
}
