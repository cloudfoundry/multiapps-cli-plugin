package mtaclient_v2

import (
	"context"
	"net/http"

	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/clients/baseclient"
	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/clients/models"
	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/clients/mtaclient_v2/operations"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"
)

const restBaseURL string = ""

type MtaV2RestClient struct {
	baseclient.BaseClient
	client    *HTTPMtaV2
	spaceGUID string
}

func NewMtaClient(host string, spaceGUID string, rt http.RoundTripper, jar http.CookieJar, tokenFactory baseclient.TokenFactory) MtaV2ClientOperations {
	t := baseclient.NewHTTPTransport(host, restBaseURL, restBaseURL, rt, jar)
	httpMtaV2Client := New(t, strfmt.Default)
	return &MtaV2RestClient{baseclient.BaseClient{TokenFactory: tokenFactory}, httpMtaV2Client, spaceGUID}
}

func (c MtaV2RestClient) GetMtas(name, namespace *string, spaceGuid string) ([]*models.Mta, error) {
	params := &operations.GetMtasV2Params{
		Context:   context.TODO(),
		Name:      name,
		Namespace: namespace,
		SpaceGUID: spaceGuid,
	}
	token, err := c.TokenFactory.NewToken()
	if err != nil {
		return nil, baseclient.NewClientError(err)
	}
	resp, err := c.client.Operations.GetMtasV2(params, token)
	if err != nil {
		return nil, baseclient.NewClientError(err)
	}
	return resp.Payload, nil
}

func (c MtaV2RestClient) GetMtasForThisSpace(name, namespace *string) ([]*models.Mta, error) {
	return c.GetMtas(name, namespace, c.spaceGUID)
}

func executeRestOperation(tokenProvider baseclient.TokenFactory, restOperation func(token runtime.ClientAuthInfoWriter) (interface{}, error)) (interface{}, error) {
	token, err := tokenProvider.NewToken()
	if err != nil {
		return nil, baseclient.NewClientError(err)
	}
	return restOperation(token)
}
