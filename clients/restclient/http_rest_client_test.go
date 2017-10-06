package restclient

import (
	"context"
	"encoding/xml"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"
	"github.com/stretchr/testify/assert"
	"github.com/SAP/cf-mta-plugin/clients/models"
	"github.com/SAP/cf-mta-plugin/clients/restclient/operations"
	"github.com/SAP/cf-mta-plugin/clients/testutil"
)

const xmlHeader = `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>`

const operationsPayload = xmlHeader + `
<ongoing-operations>` +
	operationElement + `
</ongoing-operations>`

const activeOperationsPayload = xmlHeader + `
<ongoing-operations>` +
	operationElement +
	activeElement + `
</ongoing-operations>`

const operationPayload = xmlHeader + operationElement

const operationElement = `
<ongoing-operation>
  <process-id>1</process-id>
  <process-type>deploy</process-type>
  <started-at>2016-03-04T14:23:24.521Z[Etc/UTC]</started-at>
  <space-id>5bea6497-6d70-4a31-9ad2-1ac64a520f8f</space-id>
  <user>admin</user>
  <state>SLP_TASK_STATE_ERROR</state>
  <acquired-lock>false</acquired-lock>
</ongoing-operation>`

const activeElement = `
<ongoing-operation>
  <process-id>2</process-id>
  <process-type>deploy</process-type>
  <started-at>2016-03-04T14:23:24.521Z[Etc/UTC]</started-at>
  <space-id>5bea6497-6d70-4a31-9ad2-1ac64a520f8f</space-id>
  <user>admin</user>
  <state>SLP_TASK_STATE_RUNNING</state>
  <acquired-lock>false</acquired-lock>
</ongoing-operation>`

var operationsResult = models.Operations{
	XMLName:    xml.Name{Local: "ongoing-operations"},
	Operations: []*models.Operation{&operationResult},
}

var mixedOperationsResult = models.Operations{
	XMLName:    xml.Name{Local: "ongoing-operations"},
	Operations: []*models.Operation{&operationResult, &activeOperationResult},
}

var operationResult = models.Operation{
	XMLName:      xml.Name{Local: "ongoing-operation"},
	ProcessID:    strptr("1"),
	ProcessType:  models.ProcessType("deploy"),
	StartedAt:    strptr("2016-03-04T14:23:24.521Z[Etc/UTC]"),
	SpaceID:      strptr("5bea6497-6d70-4a31-9ad2-1ac64a520f8f"),
	User:         strptr("admin"),
	State:        models.SlpTaskStateEnum("SLP_TASK_STATE_ERROR"),
	AcquiredLock: boolptr(false),
}

var activeOperationResult = models.Operation{
	XMLName:      xml.Name{Local: "ongoing-operation"},
	ProcessID:    strptr("2"),
	ProcessType:  models.ProcessType("deploy"),
	StartedAt:    strptr("2016-03-04T14:23:24.521Z[Etc/UTC]"),
	SpaceID:      strptr("5bea6497-6d70-4a31-9ad2-1ac64a520f8f"),
	User:         strptr("admin"),
	State:        models.SlpTaskStateEnum("SLP_TASK_STATE_RUNNING"),
	AcquiredLock: boolptr(false),
}

const componentsPayload = xmlHeader + `
<components>
  <mtas>` +
	mtaElement + `
  </mtas>
  <standaloneApps>
    <standaloneApp>deploy-service</standaloneApp>
  </standaloneApps>
</components>`

const mtaPayload = xmlHeader + mtaElement

const mtaElement = `
<mta>
  <metadata>
    <id>org.cloudfoundry.samples.music</id>
    <version>1.0</version>
  </metadata>
  <modules>
    <module>
      <moduleName>spring-music</moduleName>
      <appName>spring-music</appName>
      <services>
        <service>postgresql</service>
      </services>
      <providedDependencies>
        <providedDependency>spring-music</providedDependency>
      </providedDependencies>
    </module>
  </modules>
  <services>
    <service>postgresql</service>
  </services>
</mta>`

var componentsResult = models.Components{
	XMLName: xml.Name{Local: "components"},
	Mtas: models.ComponentsMtas{
		Mtas: []*models.Mta{&mtaResult},
	},
	StandaloneApps: models.ComponentsStandaloneApps{
		StandaloneApps: []string{"deploy-service"},
	},
}

var mtaResult = models.Mta{
	XMLName: xml.Name{Local: "mta"},
	Metadata: &models.MtaMetadata{
		ID:      strptr("org.cloudfoundry.samples.music"),
		Version: strptr("1.0"),
	},
	Modules: models.MtaModules{
		Modules: []*models.MtaModulesItems0{
			&models.MtaModulesItems0{
				ModuleName: strptr("spring-music"),
				AppName:    strptr("spring-music"),
				Services: models.MtaModulesItems0Services{
					Services: []string{"postgresql"},
				},
				ProvidedDependencies: models.MtaModulesItems0ProvidedDependencies{
					ProvidedDependencies: []string{"spring-music"},
				},
			},
		},
	},
	Services: models.MtaServices{
		Services: []string{"postgresql"},
	},
}

func TestGetOperations(t *testing.T) {
	server := testutil.NewGetXMLOKServer("/operations", []byte(operationsPayload))
	defer server.Close()
	client := newClient(server, nil)
	params := &operations.GetOperationsParams{
		Context: context.TODO(),
	}
	res, err := client.Operations.GetOperations(params, nil)
	if assert.NoError(t, err) {
		testutil.CheckSuccess(t, operationsResult, res.Payload)
	}
}

func TestGetOperationsWithLastOperationCount(t *testing.T) {
	server := testutil.NewGetXMLOKServer("/operations", []byte(operationsPayload))
	defer server.Close()
	client := newClient(server, nil)
	params := &operations.GetOperationsParams{
		Context: context.TODO(),
		Last:    strptr("1"),
	}
	res, err := client.Operations.GetOperations(params, nil)
	if assert.NoError(t, err) {
		testutil.CheckSuccess(t, operationsResult, res.Payload)
	}
}

func TestGetOperationsWithMixedOperationsRequested(t *testing.T) {
	server := testutil.NewGetXMLOKServer("/operations", []byte(activeOperationsPayload))
	defer server.Close()
	client := newClient(server, nil)
	params := &operations.GetOperationsParams{
		Context: context.TODO(),
		Status:  []string{"SLP_TASK_STATE_RUNNING", "SLP_TASK_STATE_ERROR"},
	}
	res, err := client.Operations.GetOperations(params, nil)
	if assert.NoError(t, err) {
		testutil.CheckSuccess(t, mixedOperationsResult, res.Payload)
	}
}

func TestGetOperation(t *testing.T) {
	server := testutil.NewGetXMLOKServer("/operations/1", []byte(operationPayload))
	defer server.Close()
	client := newClient(server, nil)
	params := &operations.GetOperationParams{
		Context:   context.TODO(),
		ProcessID: "1",
	}
	res, err := client.Operations.GetOperation(params, nil)
	if assert.NoError(t, err) {
		testutil.CheckSuccess(t, operationResult, *res.Payload)
	}
}

func TestGetComponents(t *testing.T) {
	server := testutil.NewGetXMLOKServer("/components", []byte(componentsPayload))
	defer server.Close()
	client := newClient(server, nil)
	params := &operations.GetComponentsParams{
		Context: context.TODO(),
	}
	res, err := client.Operations.GetComponents(params, nil)
	if assert.NoError(t, err) {
		testutil.CheckSuccess(t, &componentsResult, res.Payload)
	}
}

func TestGetMta(t *testing.T) {
	server := testutil.NewGetXMLOKServer("/components/org.cloudfoundry.samples.music", []byte(mtaPayload))
	defer server.Close()
	client := newClient(server, nil)
	params := &operations.GetMtaParams{
		Context: context.TODO(),
		MtaID:   "org.cloudfoundry.samples.music",
	}
	res, err := client.Operations.GetMta(params, nil)
	if assert.NoError(t, err) {
		testutil.CheckSuccess(t, &mtaResult, res.Payload)
	}
}

func newClient(server *httptest.Server, rt http.RoundTripper) *Rest {
	hu, _ := url.Parse(server.URL)
	runtime := client.New(hu.Host, "/", []string{"http"})
	if rt != nil {
		runtime.Transport = rt
	}
	return New(runtime, strfmt.Default)
}

func strptr(s string) *string {
	return &s
}

func boolptr(b bool) *bool {
	return &b
}
