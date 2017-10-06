package slppclient_test

import (
	"encoding/xml"
	"net/http"
	"net/http/cookiejar"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	baseclient "github.com/SAP/cf-mta-plugin/clients/baseclient"
	"github.com/SAP/cf-mta-plugin/clients/models"
	slppclient "github.com/SAP/cf-mta-plugin/clients/slppclient"
	"github.com/SAP/cf-mta-plugin/testutil"
)

var _ = Describe("RetryableSlppClient", func() {
	expectedRetriesCount := 4

	Describe("GetMetadata", func() {
		Context("with valid metadata returned by the backend", func() {
			It("should return the metadata returned by the backend", func() {
				client := newRetryableSlppClient(200, testutil.SlppMetadataResult)
				result, err := client.GetMetadata()
				testutil.ExpectNoErrorAndResult(err, result, &testutil.SlppMetadataResult)
			})
		})
		Context("with an error returned by the backend", func() {
			It("should return the error and a zero result", func() {
				client := newRetryableSlppClient(404, nil)
				result, err := client.GetMetadata()
				testutil.ExpectErrorAndZeroResult(err, result)
			})
		})
		Context("with an 500 server error returned by the backend", func() {
			It("should return the error, zero result and retry expectedRetriesCount times", func() {
				client := newRetryableSlppClient(500, nil)
				retryableSlppClient := client.(slppclient.RetryableSlppClient)
				result, err := retryableSlppClient.GetMetadata()
				testutil.ExpectErrorAndZeroResult(err, result)
				mockSlppClient := retryableSlppClient.SlppClient.(*MockSlppClient)
				Expect(mockSlppClient.GetRetriesCount()).To(Equal(expectedRetriesCount))
			})
		})
		Context("with an 502 server error returned by the backend", func() {
			It("should return the error, zero result and retry expectedRetriesCount times", func() {
				client := newRetryableSlppClient(502, nil)
				retryableSlppClient := client.(slppclient.RetryableSlppClient)
				result, err := retryableSlppClient.GetMetadata()
				testutil.ExpectErrorAndZeroResult(err, result)
				mockSlppClient := retryableSlppClient.SlppClient.(*MockSlppClient)
				Expect(mockSlppClient.GetRetriesCount()).To(Equal(expectedRetriesCount))
			})
		})
	})

	Describe("GetLogs", func() {
		Context("with valid process logs returned by the backend", func() {
			It("should return the process logs returned by the backend", func() {
				client := newRetryableSlppClient(200, testutil.LogsResult)
				result, err := client.GetLogs()
				testutil.ExpectNoErrorAndResult(err, result, testutil.LogsResult)
			})
		})
		Context("with an empty content returned by the backend", func() {
			It("should return empty process logs", func() {
				client := newRetryableSlppClient(200, nil)
				result, err := client.GetLogs()
				testutil.ExpectNoErrorAndResult(err, result,
					models.Logs{XMLName: xml.Name{Space: "", Local: ""}, Logs: nil})
			})
		})
		Context("with an error returned by the backend", func() {
			It("should return the error and a zero result", func() {
				client := newRetryableSlppClient(400, nil)
				result, err := client.GetLogs()
				testutil.ExpectErrorAndZeroResult(err, result)
			})
		})
		Context("with an 500 server error returned by the backend", func() {
			It("should return the error, zero result and retry expectedRetriesCount times", func() {
				client := newRetryableSlppClient(500, nil)
				retryableSlppClient := client.(slppclient.RetryableSlppClient)
				result, err := retryableSlppClient.GetLogs()
				testutil.ExpectErrorAndZeroResult(err, result)
				mockSlppClient := retryableSlppClient.SlppClient.(*MockSlppClient)
				Expect(mockSlppClient.GetRetriesCount()).To(Equal(expectedRetriesCount))
			})
		})
		Context("with an 502 server error returned by the backend", func() {
			It("should return the error, zero result and retry expectedRetriesCount times", func() {
				client := newRetryableSlppClient(502, nil)
				retryableSlppClient := client.(slppclient.RetryableSlppClient)
				result, err := retryableSlppClient.GetLogs()
				testutil.ExpectErrorAndZeroResult(err, result)
				mockSlppClient := retryableSlppClient.SlppClient.(*MockSlppClient)
				Expect(mockSlppClient.GetRetriesCount()).To(Equal(expectedRetriesCount))
			})
		})
	})

	Describe("GetLogContent", func() {
		Context("with valid process log content returned by the backend", func() {
			It("should return the process log content returned by the backend", func() {
				client := newRetryableSlppClient(200, testutil.LogContent)
				result, err := client.GetLogContent(testutil.LogID)
				testutil.ExpectNoErrorAndResult(err, result, testutil.LogContent)
			})
		})
		Context("with an empty content returned by the backend", func() {
			It("should return empty process log content", func() {
				client := newRetryableSlppClient(200, nil)
				result, err := client.GetLogContent(testutil.LogID)
				testutil.ExpectNoErrorAndResult(err, result, "")
			})
		})
		Context("with an error returned by the backend", func() {
			It("should return the error and a zero result", func() {
				client := newRetryableSlppClient(400, nil)
				result, err := client.GetLogContent(testutil.LogID)
				testutil.ExpectErrorAndZeroResult(err, result)
			})
		})
		Context("with an 500 server error returned by the backend", func() {
			It("should return the error, zero result and retry expectedRetriesCount times", func() {
				client := newRetryableSlppClient(500, nil)
				retryableSlppClient := client.(slppclient.RetryableSlppClient)
				result, err := retryableSlppClient.GetLogContent(testutil.LogID)
				testutil.ExpectErrorAndZeroResult(err, result)
				mockSlppClient := retryableSlppClient.SlppClient.(*MockSlppClient)
				Expect(mockSlppClient.GetRetriesCount()).To(Equal(expectedRetriesCount))
			})
		})
		Context("with an 502 server error returned by the backend", func() {
			It("should return the error, zero result and retry expectedRetriesCount times", func() {
				client := newRetryableSlppClient(502, nil)
				retryableSlppClient := client.(slppclient.RetryableSlppClient)
				result, err := retryableSlppClient.GetLogContent(testutil.LogID)
				testutil.ExpectErrorAndZeroResult(err, result)
				mockSlppClient := retryableSlppClient.SlppClient.(*MockSlppClient)
				Expect(mockSlppClient.GetRetriesCount()).To(Equal(expectedRetriesCount))
			})
		})
	})

	Describe("GetTasklist", func() {
		Context("with valid tasklist returned by the backend", func() {
			It("should return the tasklist returned by the backend", func() {
				client := newRetryableSlppClient(200, testutil.TasklistResult)
				result, err := client.GetTasklist()
				testutil.ExpectNoErrorAndResult(err, result, testutil.TasklistResult)
			})
		})
		Context("with an empty content returned by the backend", func() {
			It("should return empty tasklist", func() {
				client := newRetryableSlppClient(200, nil)
				result, err := client.GetTasklist()
				testutil.ExpectNoErrorAndResult(err, result, models.Tasklist{})
			})
		})
		Context("with an error returned by the backend", func() {
			It("should return the error and a zero result", func() {
				client := newRetryableSlppClient(400, nil)
				result, err := client.GetTasklist()
				testutil.ExpectErrorAndZeroResult(err, result)
			})
		})
		Context("with an 500 server error returned by the backend", func() {
			It("should return the error, zero result and retry expectedRetriesCount times", func() {
				client := newRetryableSlppClient(500, nil)
				retryableSlppClient := client.(slppclient.RetryableSlppClient)
				result, err := retryableSlppClient.GetTasklist()
				testutil.ExpectErrorAndZeroResult(err, result)
				mockSlppClient := retryableSlppClient.SlppClient.(*MockSlppClient)
				Expect(mockSlppClient.GetRetriesCount()).To(Equal(expectedRetriesCount))
			})
		})
		Context("with an 502 server error returned by the backend", func() {
			It("should return the error, zero result and retry expectedRetriesCount times", func() {
				client := newRetryableSlppClient(502, nil)
				retryableSlppClient := client.(slppclient.RetryableSlppClient)
				result, err := retryableSlppClient.GetTasklist()
				testutil.ExpectErrorAndZeroResult(err, result)
				mockSlppClient := retryableSlppClient.SlppClient.(*MockSlppClient)
				Expect(mockSlppClient.GetRetriesCount()).To(Equal(expectedRetriesCount))
			})
		})
	})

	Describe("GetTasklistTask", func() {
		Context("with valid task returned by the backend", func() {
			It("should return the task returned by the backend", func() {
				client := newRetryableSlppClient(200, &testutil.TaskResult)
				result, err := client.GetTasklistTask(testutil.ServiceID)
				testutil.ExpectNoErrorAndResult(err, result, &testutil.TaskResult)
			})
		})
		Context("with an empty content returned by the backend", func() {
			It("should return empty task", func() {
				client := newRetryableSlppClient(200, nil)
				result, err := client.GetTasklistTask(testutil.ServiceID)
				testutil.ExpectNoErrorAndResult(err, result, &models.Task{})
			})
		})
		Context("with an error returned by the backend", func() {
			It("should return the error and a zero result", func() {
				client := newRetryableSlppClient(400, nil)
				result, err := client.GetTasklistTask(testutil.ServiceID)
				testutil.ExpectErrorAndZeroResult(err, result)
			})
		})
		Context("with an 500 server error returned by the backend", func() {
			It("should return the error, zero result and retry expectedRetriesCount times", func() {
				client := newRetryableSlppClient(500, nil)
				retryableSlppClient := client.(slppclient.RetryableSlppClient)
				result, err := retryableSlppClient.GetTasklistTask(testutil.ServiceID)
				testutil.ExpectErrorAndZeroResult(err, result)
				mockSlppClient := retryableSlppClient.SlppClient.(*MockSlppClient)
				Expect(mockSlppClient.GetRetriesCount()).To(Equal(expectedRetriesCount))
			})
		})
		Context("with an 502 server error returned by the backend", func() {
			It("should return the error, zero result and retry expectedRetriesCount times", func() {
				client := newRetryableSlppClient(502, nil)
				retryableSlppClient := client.(slppclient.RetryableSlppClient)
				result, err := retryableSlppClient.GetTasklistTask(testutil.ServiceID)
				testutil.ExpectErrorAndZeroResult(err, result)
				mockSlppClient := retryableSlppClient.SlppClient.(*MockSlppClient)
				Expect(mockSlppClient.GetRetriesCount()).To(Equal(expectedRetriesCount))
			})
		})
	})

	Describe("GetError", func() {
		Context("with valid client error returned by the backend", func() {
			It("should return the client error returned by the backend", func() {
				client := newRetryableSlppClient(200, testutil.ErrorResult)
				result, err := client.GetError()
				testutil.ExpectNoErrorAndResult(err, result, testutil.ErrorResult)
			})
		})
		Context("with an empty content returned by the backend", func() {
			It("should return empty client error", func() {
				client := newRetryableSlppClient(200, nil)
				result, err := client.GetError()
				testutil.ExpectNoErrorAndResult(err, result, &models.Error{})
			})
		})
		Context("with an error returned by the backend", func() {
			It("should return the error and a zero result", func() {
				client := newRetryableSlppClient(400, nil)
				result, err := client.GetError()
				testutil.ExpectErrorAndZeroResult(err, result)
			})
		})
		Context("with an 500 server error returned by the backend", func() {
			It("should return the error, zero result and retry expectedRetriesCount times", func() {
				client := newRetryableSlppClient(500, nil)
				retryableSlppClient := client.(slppclient.RetryableSlppClient)
				result, err := retryableSlppClient.GetError()
				testutil.ExpectErrorAndZeroResult(err, result)
				mockSlppClient := retryableSlppClient.SlppClient.(*MockSlppClient)
				Expect(mockSlppClient.GetRetriesCount()).To(Equal(expectedRetriesCount))
			})
		})
		Context("with an 502 server error returned by the backend", func() {
			It("should return the error, zero result and retry expectedRetriesCount times", func() {
				client := newRetryableSlppClient(502, nil)
				retryableSlppClient := client.(slppclient.RetryableSlppClient)
				result, err := retryableSlppClient.GetError()
				testutil.ExpectErrorAndZeroResult(err, result)
				mockSlppClient := retryableSlppClient.SlppClient.(*MockSlppClient)
				Expect(mockSlppClient.GetRetriesCount()).To(Equal(expectedRetriesCount))
			})
		})
	})

	Describe("ExecuteAction", func() {
		Context("with success returned by the backend", func() {
			It("should not return any errors", func() {
				client := newRetryableSlppClient(204, nil)
				err := client.ExecuteAction(testutil.ActionID)
				testutil.ExpectNoError(err)
			})
		})
		Context("with an error returned by the backend", func() {
			It("should return the error", func() {
				client := newRetryableSlppClient(400, nil)
				err := client.ExecuteAction(testutil.ActionID)
				testutil.ExpectError(err)
			})
		})
		Context("with an 500 server error returned by the backend", func() {
			It("should return the error, zero result and retry expectedRetriesCount times", func() {
				client := newRetryableSlppClient(500, nil)
				retryableSlppClient := client.(slppclient.RetryableSlppClient)
				err := retryableSlppClient.ExecuteAction(testutil.ActionID)
				testutil.ExpectError(err)
				mockSlppClient := retryableSlppClient.SlppClient.(*MockSlppClient)
				Expect(mockSlppClient.GetRetriesCount()).To(Equal(expectedRetriesCount))
			})
		})
		Context("with an 502 server error returned by the backend", func() {
			It("should return the error, zero result and retry expectedRetriesCount times", func() {
				client := newRetryableSlppClient(502, nil)
				retryableSlppClient := client.(slppclient.RetryableSlppClient)
				err := retryableSlppClient.ExecuteAction(testutil.ActionID)
				testutil.ExpectError(err)
				mockSlppClient := retryableSlppClient.SlppClient.(*MockSlppClient)
				Expect(mockSlppClient.GetRetriesCount()).To(Equal(expectedRetriesCount))
			})
		})
	})
})

func newRetryableSlppClient(statusCode int, v interface{}) slppclient.SlppClientOperations {
	tokenFactory := baseclient.NewCustomTokenFactory("test-token")
	cookieJar, _ := cookiejar.New(nil)
	roundTripper := testutil.NewCustomTransport(statusCode, v)
	return NewMockRetryableSlppClient("http://localhost:1000", "test-org", "test-space", testutil.ServiceID, testutil.ProcessID,
		roundTripper, cookieJar, tokenFactory)
}

// NewMockRetryableSlppClient creates a new retryable REST client
func NewMockRetryableSlppClient(host, org, space, serviceID, processID string, rt http.RoundTripper, jar http.CookieJar, tokenFactory baseclient.TokenFactory) slppclient.SlppClientOperations {
	slppClient := NewMockSlppClient(host, org, space, serviceID, processID, rt, jar, tokenFactory)
	return slppclient.RetryableSlppClient{SlppClient: slppClient, MaxRetriesCount: 3, RetryInterval: time.Microsecond * 1}
}

// MockSlppClient represents a client for the SLPP protocol
type MockSlppClient struct {
	slppClient   slppclient.SlppClientOperations
	retriesCount int
}

// NewMockSlppClient creates a new SLPP client
func NewMockSlppClient(host, org, space, serviceID, processID string, rt http.RoundTripper, jar http.CookieJar, tokenFactory baseclient.TokenFactory) slppclient.SlppClientOperations {
	slppClient := slppclient.NewSlppClient(host, org, space, serviceID, processID, rt, jar, tokenFactory)
	return &MockSlppClient{slppClient, 0}
}

// GetMetadata retrieves the SLPP metadata
func (c *MockSlppClient) GetMetadata() (*models.Metadata, error) {
	c.retriesCount++
	return c.slppClient.GetMetadata()
}

// GetLogs retrieves all process logs
func (c *MockSlppClient) GetLogs() (models.Logs, error) {
	c.retriesCount++
	return c.slppClient.GetLogs()
}

// GetLogContent retrieves the content of the specified logs
func (c *MockSlppClient) GetLogContent(logID string) (string, error) {
	c.retriesCount++
	return c.slppClient.GetLogContent(logID)
}

//GetTasklist retrieves the tasklist for the current process
func (c *MockSlppClient) GetTasklist() (models.Tasklist, error) {
	c.retriesCount++
	return c.slppClient.GetTasklist()
}

//GetTasklistTask retrieves concrete task by taskID
func (c *MockSlppClient) GetTasklistTask(taskID string) (*models.Task, error) {
	c.retriesCount++
	return c.slppClient.GetTasklistTask(taskID)
}

//GetServiceID returns serviceID
func (c *MockSlppClient) GetServiceID() string {
	c.retriesCount++
	return c.slppClient.GetServiceID()
}

//GetError retrieves the current client error
func (c *MockSlppClient) GetError() (*models.Error, error) {
	c.retriesCount++
	return c.slppClient.GetError()
}

//ExecuteAction executes an action specified with actionID
func (c *MockSlppClient) ExecuteAction(actionID string) error {
	c.retriesCount++
	return c.slppClient.ExecuteAction(actionID)
}

//GetActions retrieves the list of available actions for the current process
func (c *MockSlppClient) GetActions() (models.Actions, error) {
	c.retriesCount++
	return c.slppClient.GetActions()
}

// GetRetriesCount returns retries count
func (c *MockSlppClient) GetRetriesCount() int {
	return c.retriesCount
}
