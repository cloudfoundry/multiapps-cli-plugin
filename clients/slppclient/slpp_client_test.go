package slppclient_test

import (
	"encoding/xml"
	"net/http/cookiejar"

	. "github.com/onsi/ginkgo"
	baseclient "github.com/SAP/cf-mta-plugin/clients/baseclient"
	models "github.com/SAP/cf-mta-plugin/clients/models"
	slppclient "github.com/SAP/cf-mta-plugin/clients/slppclient"
	"github.com/SAP/cf-mta-plugin/testutil"
)

var _ = Describe("SlppClient", func() {

	Describe("GetMetadata", func() {
		Context("with valid metadata returned by the backend", func() {
			It("should return the metadata returned by the backend", func() {
				client := newSlppClient(200, testutil.SlppMetadataResult)
				result, err := client.GetMetadata()
				testutil.ExpectNoErrorAndResult(err, result, &testutil.SlppMetadataResult)
			})
		})
		Context("with an error returned by the backend", func() {
			It("should return the error and a zero result", func() {
				client := newSlppClient(404, nil)
				result, err := client.GetMetadata()
				testutil.ExpectErrorAndZeroResult(err, result)
			})
		})
	})

	Describe("GetLogs", func() {
		Context("with valid process logs returned by the backend", func() {
			It("should return the process logs returned by the backend", func() {
				client := newSlppClient(200, testutil.LogsResult)
				result, err := client.GetLogs()
				testutil.ExpectNoErrorAndResult(err, result, testutil.LogsResult)
			})
		})
		Context("with an empty content returned by the backend", func() {
			It("should return empty process logs", func() {
				client := newSlppClient(200, nil)
				result, err := client.GetLogs()
				testutil.ExpectNoErrorAndResult(err, result,
					models.Logs{XMLName: xml.Name{Space: "", Local: ""}, Logs: nil})
			})
		})
		Context("with an error returned by the backend", func() {
			It("should return the error and a zero result", func() {
				client := newSlppClient(400, nil)
				result, err := client.GetLogs()
				testutil.ExpectErrorAndZeroResult(err, result)
			})
		})
	})

	Describe("GetLogContent", func() {
		Context("with valid process log content returned by the backend", func() {
			It("should return the process log content returned by the backend", func() {
				client := newSlppClient(200, testutil.LogContent)
				result, err := client.GetLogContent(testutil.LogID)
				testutil.ExpectNoErrorAndResult(err, result, testutil.LogContent)
			})
		})
		Context("with an empty content returned by the backend", func() {
			It("should return empty process log content", func() {
				client := newSlppClient(200, nil)
				result, err := client.GetLogContent(testutil.LogID)
				testutil.ExpectNoErrorAndResult(err, result, "")
			})
		})
		Context("with an error returned by the backend", func() {
			It("should return the error and a zero result", func() {
				client := newSlppClient(400, nil)
				result, err := client.GetLogContent(testutil.LogID)
				testutil.ExpectErrorAndZeroResult(err, result)
			})
		})
	})

	Describe("GetTasklist", func() {
		Context("with valid tasklist returned by the backend", func() {
			It("should return the tasklist returned by the backend", func() {
				client := newSlppClient(200, testutil.TasklistResult)
				result, err := client.GetTasklist()
				testutil.ExpectNoErrorAndResult(err, result, testutil.TasklistResult)
			})
		})
		Context("with an empty content returned by the backend", func() {
			It("should return empty tasklist", func() {
				client := newSlppClient(200, nil)
				result, err := client.GetTasklist()
				testutil.ExpectNoErrorAndResult(err, result, models.Tasklist{})
			})
		})
		Context("with an error returned by the backend", func() {
			It("should return the error and a zero result", func() {
				client := newSlppClient(400, nil)
				result, err := client.GetTasklist()
				testutil.ExpectErrorAndZeroResult(err, result)
			})
		})
	})

	Describe("GetTasklistTask", func() {
		Context("with valid task returned by the backend", func() {
			It("should return the task returned by the backend", func() {
				client := newSlppClient(200, &testutil.TaskResult)
				result, err := client.GetTasklistTask(testutil.ServiceID)
				testutil.ExpectNoErrorAndResult(err, result, &testutil.TaskResult)
			})
		})
		Context("with an empty content returned by the backend", func() {
			It("should return empty tasklist", func() {
				client := newSlppClient(200, nil)
				result, err := client.GetTasklistTask(testutil.ServiceID)
				testutil.ExpectNoErrorAndResult(err, result, &models.Task{})
			})
		})
		Context("with an error returned by the backend", func() {
			It("should return the error and a zero result", func() {
				client := newSlppClient(400, nil)
				result, err := client.GetTasklistTask(testutil.ServiceID)
				testutil.ExpectErrorAndZeroResult(err, result)
			})
		})
	})

	Describe("GetError", func() {
		Context("with valid client error returned by the backend", func() {
			It("should return the client error returned by the backend", func() {
				client := newSlppClient(200, testutil.ErrorResult)
				result, err := client.GetError()
				testutil.ExpectNoErrorAndResult(err, result, testutil.ErrorResult)
			})
		})
		Context("with an empty content returned by the backend", func() {
			It("should return empty client error", func() {
				client := newSlppClient(200, nil)
				result, err := client.GetError()
				testutil.ExpectNoErrorAndResult(err, result, &models.Error{})
			})
		})
		Context("with an error returned by the backend", func() {
			It("should return the error and a zero result", func() {
				client := newSlppClient(400, nil)
				result, err := client.GetError()
				testutil.ExpectErrorAndZeroResult(err, result)
			})
		})
	})

	Describe("ExecuteAction", func() {
		Context("with success returned by the backend", func() {
			It("should not return any errors", func() {
				client := newSlppClient(204, nil)
				err := client.ExecuteAction(testutil.ActionID)
				testutil.ExpectNoError(err)
			})
		})
		Context("with an error returned by the backend", func() {
			It("should return the error", func() {
				client := newSlppClient(400, nil)
				err := client.ExecuteAction(testutil.ActionID)
				testutil.ExpectError(err)
			})
		})
	})
})

func newSlppClient(statusCode int, v interface{}) slppclient.SlppClientOperations {
	tokenFactory := baseclient.NewCustomTokenFactory("test-token")
	cookieJar, _ := cookiejar.New(nil)
	roundTripper := testutil.NewCustomTransport(statusCode, v)
	return slppclient.NewSlppClient("http://localhost:1000", "test-org", "test-space", testutil.ServiceID, testutil.ProcessID,
		roundTripper, cookieJar, tokenFactory)
}
