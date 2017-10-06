package slmpclient

import (
	"context"
	"encoding/xml"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"

	"github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"
	"github.com/stretchr/testify/assert"
	"github.com/SAP/cf-mta-plugin/clients/csrf"
	"github.com/SAP/cf-mta-plugin/clients/models"
	"github.com/SAP/cf-mta-plugin/clients/slmpclient/operations"
	testutil "github.com/SAP/cf-mta-plugin/clients/testutil"
)

const xmlHeader = `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>`
const xmlns = "http://www.sap.com/lmsl/slp"

const metadataPayload = xmlHeader + `
<Metadata xmlns="http://www.sap.com/lmsl/slp">
  <slmpversion>1.2.0</slmpversion>
</Metadata>`

var metadataResult = models.Metadata{
	XMLName:     xml.Name{Space: xmlns, Local: "Metadata"},
	Slmpversion: "1.2.0",
}

const servicesPayload = xmlHeader + `
<services xmlns="http://www.sap.com/lmsl/slp">` +
	serviceElement + `
</services>`

const servicePayload = xmlHeader + serviceElement

const serviceElement = `
<Service xmlns="http://www.sap.com/lmsl/slp">
  <id>xs2-undeploy</id>
  <processes>services/xs2-undeploy/processes</processes>
	<parameters>` +
	parameterElement + `
	</parameters>
  <files>services/xs2-undeploy/files</files>
  <versions>services/xs2-undeploy/versions</versions>
  <slppversion>1.2.0</slppversion>
  <displayName>Undeploy</displayName>
  <description>Undeploy</description>
</Service>`

var servicesResult = models.Services{
	XMLName:  xml.Name{Space: xmlns, Local: "services"},
	Services: []*models.Service{&serviceResult},
}

var serviceResult = models.Service{
	XMLName:   xml.Name{Space: xmlns, Local: "Service"},
	ID:        uriptr(strfmt.URI("xs2-undeploy")),
	Processes: uriptr(strfmt.URI("services/xs2-undeploy/processes")),
	Parameters: models.ServiceParameters{
		Parameters: []*models.Parameter{&parameterResult},
	},
	Files:       strfmt.URI("services/xs2-undeploy/files"),
	Versions:    strfmt.URI("services/xs2-undeploy/versions"),
	Slppversion: strptr("1.2.0"),
	DisplayName: "Undeploy",
	Description: "Undeploy",
}

const processesPayload = xmlHeader + `
<processes xmlns="http://www.sap.com/lmsl/slp">` +
	processElement + `
</processes>`

const processPayload = xmlHeader + processElement

const processElement = `
<Process xmlns="http://www.sap.com/lmsl/slp">
  <id>1</id>
  <service>xs2-deploy</service>
  <status>slp.process.state.ACTIVE</status>
  <rootURL>runs/xs2-deploy/1</rootURL>
  <parameters>
    <Parameter>
      <id>appArchiveId</id>
      <type>slp.parameter.type.SCALAR</type>
      <required>true</required>
      <value>d204bf6a-5a56-4c91-952e-1e8dce81fca2</value>
    </Parameter>
  </parameters>
  <displayName>Deploy</displayName>
  <description>Deploy</description>
</Process>`

var processesResult = models.Processes{
	XMLName:   xml.Name{Space: xmlns, Local: "processes"},
	Processes: []*models.Process{&processResult},
}

var processResult = models.Process{
	XMLName: xml.Name{Space: xmlns, Local: "Process"},
	ID:      "1",
	Service: uriptr(strfmt.URI("xs2-deploy")),
	Status:  models.SlpProcessState("slp.process.state.ACTIVE"),
	RootURL: strfmt.URI("runs/xs2-deploy/1"),
	Parameters: models.ProcessParameters{
		Parameters: []*models.Parameter{
			&models.Parameter{
				XMLName:  xml.Name{Space: xmlns, Local: "Parameter"},
				ID:       strptr("appArchiveId"),
				Type:     models.SlpParameterType("slp.parameter.type.SCALAR"),
				Required: true,
				Value:    "d204bf6a-5a56-4c91-952e-1e8dce81fca2",
			},
		},
	},
	DisplayName: "Deploy",
	Description: "Deploy",
}

const versionsPayload = xmlHeader + `
<versions xmlns="http://www.sap.com/lmsl/slp">
  <ComponentVersion>
    <id>xs2-deploy_VERSIONS_1.0</id>
    <component>xs2-deploy</component>
    <version>1.0</version>
  </ComponentVersion>
</versions>`

var versionsResult = models.Versions{
	XMLName: xml.Name{Space: xmlns, Local: "versions"},
	ComponentVersions: []*models.ComponentVersion{
		&models.ComponentVersion{
			XMLName:   xml.Name{Space: xmlns, Local: "ComponentVersion"},
			ID:        strptr("xs2-deploy_VERSIONS_1.0"),
			Component: strptr("xs2-deploy"),
			Version:   strptr("1.0"),
		},
	},
}

const parametersPayload = xmlHeader + `
<parameters xmlns="http://www.sap.com/lmsl/slp">` +
	parameterElement + `
</parameters>`

const parameterPayload = xmlHeader + parameterElement

const parameterElement = `
<Parameter xmlns="http://www.sap.com/lmsl/slp">
	<id>mtaId</id>
	<type>slp.parameter.type.SCALAR</type>
	<required>true</required>
</Parameter>`

var parametersResult = models.Parameters{
	XMLName:    xml.Name{Space: xmlns, Local: "parameters"},
	Parameters: []*models.Parameter{&parameterResult},
}

var parameterResult = models.Parameter{
	XMLName:  xml.Name{Space: xmlns, Local: "Parameter"},
	ID:       strptr("mtaId"),
	Type:     models.SlpParameterType("slp.parameter.type.SCALAR"),
	Required: true,
}

const filesPayload = xmlHeader + `
<files xmlns="http://www.sap.com/lmsl/slp">` +
	fileElement + `
</files>`

const filePayload = xmlHeader + fileElement

const fileElement = `
<File xmlns="http://www.sap.com/lmsl/slp">
  <id>d204bf6a-5a56-4c91-952e-1e8dce81fca2</id>
  <filePath>xs2-deploy</filePath>
  <fileSize>5</fileSize>
  <fileName>test.txt</fileName>
  <digest>D8E8FCA2DC0F896FD7CB4CB0031BA249</digest>
  <digestAlgorithm>MD5</digestAlgorithm>
</File>`

var filesResult = models.Files{
	XMLName: xml.Name{Space: xmlns, Local: "files"},
	Files:   []*models.File{&fileResult},
}

// TODO Fix DateTime parsing by the client
// var dateTime, _ = strfmt.ParseDateTime("2016-03-04T00:00:00Z")
// var timestamp = models.SlpTimestamp(dateTime)

var fileResult = models.File{
	XMLName:  xml.Name{Space: xmlns, Local: "File"},
	ID:       strptr("d204bf6a-5a56-4c91-952e-1e8dce81fca2"),
	FilePath: "xs2-deploy",
	FileSize: 5,
	FileName: strptr("test.txt"),
	// ModificationTime: &timestamp,
	Digest:          "D8E8FCA2DC0F896FD7CB4CB0031BA249",
	DigestAlgorithm: "MD5",
}

func TestGetMetadata(t *testing.T) {
	server := testutil.NewGetXMLOKServer("/metadata", []byte(metadataPayload))
	defer server.Close()
	client := newClient(server, nil)
	res, err := client.Operations.GetMetadata(nil, nil)
	if assert.NoError(t, err) {
		testutil.CheckSuccess(t, metadataResult, *res.Payload)
	}
}

func TestGetMetadata_OAuth(t *testing.T) {
	server := testutil.NewGetXMLOAuthOKServer("/metadata", "dummy", []byte(metadataPayload))
	defer server.Close()
	slmp := newClient(server, nil)
	authInfo := client.BearerToken("dummy")
	res, err := slmp.Operations.GetMetadata(nil, authInfo)
	if assert.NoError(t, err) {
		testutil.CheckSuccess(t, metadataResult, *res.Payload)
	}
}

func TestGetMetadata_Error(t *testing.T) {
	server := testutil.NewStatusServer(http.StatusUnauthorized)
	defer server.Close()
	client := newClient(server, nil)
	res, err := client.Operations.GetMetadata(nil, nil)
	if assert.Error(t, err) {
		testutil.CheckError(t, err, "unknown error", http.StatusUnauthorized, res)
	}
}

func TestGetServices(t *testing.T) {
	server := testutil.NewGetXMLOKServer("/services", []byte(servicesPayload))
	defer server.Close()
	client := newClient(server, nil)
	res, err := client.Operations.GetServices(nil, nil)
	if assert.NoError(t, err) {
		testutil.CheckSuccess(t, servicesResult, res.Payload)
	}
}

func TestGetService(t *testing.T) {
	server := testutil.NewGetXMLOKServer("/services/xs2-undeploy", []byte(servicePayload))
	defer server.Close()
	client := newClient(server, nil)
	params := &operations.GetServiceParams{
		Context:   context.TODO(),
		ServiceID: "xs2-undeploy",
	}
	res, err := client.Operations.GetService(params, nil)
	if assert.NoError(t, err) {
		testutil.CheckSuccess(t, serviceResult, *res.Payload)
	}
}

func TestGetServiceProcesses(t *testing.T) {
	server := testutil.NewGetXMLOKServer("/services/xs2-deploy/processes", []byte(processesPayload))
	defer server.Close()
	client := newClient(server, nil)
	params := &operations.GetServiceProcessesParams{
		Context:   context.TODO(),
		ServiceID: "xs2-deploy",
	}
	res, err := client.Operations.GetServiceProcesses(params, nil)
	if assert.NoError(t, err) {
		testutil.CheckSuccess(t, processesResult, res.Payload)
	}
}

func TestCreateServiceProcess(t *testing.T) {
	process := models.Process{
		Service: uriptr(strfmt.URI("xs2-deploy")),
		Parameters: models.ProcessParameters{
			Parameters: []*models.Parameter{
				&models.Parameter{
					ID:    strptr("appArchiveId"),
					Type:  models.SlpParameterType("slp.parameter.type.SCALAR"),
					Value: "d204bf6a-5a56-4c91-952e-1e8dce81fca2",
				},
			},
		},
	}
	body, _ := xml.Marshal(process)
	server := testutil.NewPostXMLOKServer("/services/xs2-deploy/processes", body, []byte(processPayload))
	defer server.Close()
	client := newClient(server, nil)
	params := &operations.CreateServiceProcessParams{
		Context:   context.TODO(),
		ServiceID: "xs2-deploy",
		Process:   &process,
	}
	res, err := client.Operations.CreateServiceProcess(params, nil)
	if assert.NoError(t, err) {
		testutil.CheckSuccess(t, processResult, *res.Payload)
	}
}

func TestDeleteServiceProcess(t *testing.T) {
	server := testutil.NewDeleteNoContentServer("/services/xs2-deploy/processes/1")
	defer server.Close()
	client := newClient(server, nil)
	params := &operations.DeleteServiceProcessParams{
		Context:   context.TODO(),
		ServiceID: "xs2-deploy",
		ProcessID: "1",
	}
	_, err := client.Operations.DeleteServiceProcess(params, nil)
	assert.NoError(t, err)
}

func TestGetServiceParameters(t *testing.T) {
	server := testutil.NewGetXMLOKServer("/services/xs2-undeploy/parameters", []byte(parametersPayload))
	defer server.Close()
	client := newClient(server, nil)
	params := &operations.GetServiceParametersParams{
		Context:   context.TODO(),
		ServiceID: "xs2-undeploy",
	}
	res, err := client.Operations.GetServiceParameters(params, nil)
	if assert.NoError(t, err) {
		testutil.CheckSuccess(t, parametersResult, res.Payload)
	}
}

func TestGetServiceVersions(t *testing.T) {
	server := testutil.NewGetXMLOKServer("/services/xs2-deploy/versions", []byte(versionsPayload))
	defer server.Close()
	client := newClient(server, nil)
	params := &operations.GetServiceVersionsParams{
		Context:   context.TODO(),
		ServiceID: "xs2-deploy",
	}
	res, err := client.Operations.GetServiceVersions(params, nil)
	if assert.NoError(t, err) {
		testutil.CheckSuccess(t, versionsResult, res.Payload)
	}
}

func TestGetServiceFiles(t *testing.T) {
	server := testutil.NewGetXMLOKServer("/services/xs2-deploy/files", []byte(filesPayload))
	defer server.Close()
	client := newClient(server, nil)
	params := &operations.GetServiceFilesParams{
		Context:   context.TODO(),
		ServiceID: "xs2-deploy",
	}
	res, err := client.Operations.GetServiceFiles(params, nil)
	if assert.NoError(t, err) {
		testutil.CheckSuccess(t, filesResult, res.Payload)
	}
}

func TestCreateServiceFiles(t *testing.T) {
	fileName, data := "test.txt", []byte("This is a test")
	server := testutil.NewPostFileXMLOKServer("/services/xs2-deploy/files", "files", data, []byte(filesPayload))
	defer server.Close()
	err := ioutil.WriteFile(fileName, data, 0777)
	assert.NoError(t, err)
	defer os.Remove(fileName)
	file, err := os.Open(fileName)
	assert.NoError(t, err)
	defer file.Close()
	client := newClient(server, nil)
	params := &operations.CreateServiceFilesParams{
		Context:   context.TODO(),
		ServiceID: "xs2-deploy",
		Files:     *file,
	}
	res, err := client.Operations.CreateServiceFiles(params, nil)
	if assert.NoError(t, err) {
		testutil.CheckSuccess(t, filesResult, res.Payload)
	}
}

func TestDeleteServiceFiles(t *testing.T) {
	server := testutil.NewDeleteNoContentServer("/services/xs2-deploy/files")
	defer server.Close()
	client := newClient(server, nil)
	params := &operations.DeleteServiceFilesParams{
		Context:   context.TODO(),
		ServiceID: "xs2-deploy",
	}
	_, err := client.Operations.DeleteServiceFiles(params, nil)
	assert.NoError(t, err)
}

func TestDeleteServiceFiles_Csrf(t *testing.T) {
	server := testutil.NewGetXMLOKDeleteNoContentCsrfServer("/services/xs2-deploy/files", []byte(filesPayload))
	defer server.Close()
	csrfx := csrf.Csrf{Header: "", Token: ""}
	transport := csrf.Transport{Transport: http.DefaultTransport, Csrf: &csrfx}
	client := newClient(server, transport)
	params := &operations.GetServiceFilesParams{
		Context:   context.TODO(),
		ServiceID: "xs2-deploy",
	}
	res, err := client.Operations.GetServiceFiles(params, nil)
	if assert.NoError(t, err) {
		testutil.CheckSuccess(t, filesResult, res.Payload)
	}
	delparams := &operations.DeleteServiceFilesParams{
		Context:   context.TODO(),
		ServiceID: "xs2-deploy",
	}
	_, err = client.Operations.DeleteServiceFiles(delparams, nil)
	assert.NoError(t, err)
}

func TestDeleteServiceFile(t *testing.T) {
	server := testutil.NewDeleteNoContentServer("/services/xs2-deploy/files/d204bf6a-5a56-4c91-952e-1e8dce81fca2")
	defer server.Close()
	client := newClient(server, nil)
	params := &operations.DeleteServiceFileParams{
		Context:   context.TODO(),
		ServiceID: "xs2-deploy",
		FileID:    "d204bf6a-5a56-4c91-952e-1e8dce81fca2",
	}
	_, err := client.Operations.DeleteServiceFile(params, nil)
	assert.NoError(t, err)
}

func TestGetProcesses(t *testing.T) {
	server := testutil.NewGetXMLOKServer("/processes", []byte(processesPayload))
	defer server.Close()
	client := newClient(server, nil)
	res, err := client.Operations.GetProcesses(nil, nil)
	if assert.NoError(t, err) {
		testutil.CheckSuccess(t, processesResult, res.Payload)
	}
}

func TestGetProcess(t *testing.T) {
	server := testutil.NewGetXMLOKServer("/processes/1", []byte(processPayload))
	defer server.Close()
	client := newClient(server, nil)
	params := &operations.GetProcessParams{
		Context:   context.TODO(),
		ProcessID: "1",
	}
	res, err := client.Operations.GetProcess(params, nil)
	if assert.NoError(t, err) {
		testutil.CheckSuccess(t, processResult, *res.Payload)
	}
}

func TestDeleteProcess(t *testing.T) {
	server := testutil.NewDeleteNoContentServer("/processes/1")
	defer server.Close()
	client := newClient(server, nil)
	params := &operations.DeleteProcessParams{
		Context:   context.TODO(),
		ProcessID: "1",
	}
	_, err := client.Operations.DeleteProcess(params, nil)
	assert.NoError(t, err)
}

func TestGetProcessService(t *testing.T) {
	server := testutil.NewGetXMLOKServer("/processes/2/service", []byte(servicePayload))
	defer server.Close()
	client := newClient(server, nil)
	params := &operations.GetProcessServiceParams{
		Context:   context.TODO(),
		ProcessID: "2",
	}
	res, err := client.Operations.GetProcessService(params, nil)
	if assert.NoError(t, err) {
		testutil.CheckSuccess(t, serviceResult, *res.Payload)
	}
}

func newClient(server *httptest.Server, rt http.RoundTripper) *Slmp {
	hu, _ := url.Parse(server.URL)
	slmp := client.New(hu.Host, "/", []string{"http"})
	if rt != nil {
		slmp.Transport = rt
	}
	return New(slmp, strfmt.Default)
}

func strptr(s string) *string {
	return &s
}

func uriptr(u strfmt.URI) *strfmt.URI {
	return &u
}
