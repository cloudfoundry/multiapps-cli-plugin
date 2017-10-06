package slmpclient_test

import (
	"net/http"
	"net/http/cookiejar"
	"os"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	baseclient "github.com/SAP/cf-mta-plugin/clients/baseclient"
	"github.com/SAP/cf-mta-plugin/clients/models"
	slmpclient "github.com/SAP/cf-mta-plugin/clients/slmpclient"
	"github.com/SAP/cf-mta-plugin/testutil"
)

var _ = Describe("RetryableSlmpClient", func() {
	expectedRetriesCount := 4

	Describe("GetMetadata", func() {
		Context("with valid metadata returned by the backend", func() {
			It("should return the metadata returned by the backend", func() {
				client := newRetryableSlmpClient(200, testutil.SlmpMetadataResult)
				result, err := client.GetMetadata()
				testutil.ExpectNoErrorAndResult(err, result, &testutil.SlmpMetadataResult)
			})
		})
		Context("with an error returned by the backend", func() {
			It("should return the error and a zero result", func() {
				client := newRetryableSlmpClient(404, nil)
				result, err := client.GetMetadata()
				testutil.ExpectErrorAndZeroResult(err, result)
			})
		})
		Context("with an 500 server error returned by the backend", func() {
			It("should return the error, zero result and retry expectedRetriesCount times", func() {
				client := newRetryableSlmpClient(500, nil)
				retryableSlmpClient := client.(slmpclient.RetryableSlmpClient)
				result, err := retryableSlmpClient.GetMetadata()
				testutil.ExpectErrorAndZeroResult(err, result)
				mockSlmpClient := retryableSlmpClient.SlmpClient.(*MockSlmpClient)
				Expect(mockSlmpClient.GetRetriesCount()).To(Equal(expectedRetriesCount))
			})
		})
		Context("with an 502 server error returned by the backend", func() {
			It("should return the error, zero result and retry expectedRetriesCount times", func() {
				client := newRetryableSlmpClient(502, nil)
				retryableSlmpClient := client.(slmpclient.RetryableSlmpClient)
				result, err := retryableSlmpClient.GetMetadata()
				testutil.ExpectErrorAndZeroResult(err, result)
				mockSlmpClient := retryableSlmpClient.SlmpClient.(*MockSlmpClient)
				Expect(mockSlmpClient.GetRetriesCount()).To(Equal(expectedRetriesCount))
			})
		})
	})

	Describe("GetServices", func() {
		Context("with valid services returned by the backend", func() {
			It("should return all services returned by the backend", func() {
				client := newRetryableSlmpClient(200, testutil.ServicesResult)
				result, err := client.GetServices()
				testutil.ExpectNoErrorAndResult(err, result, testutil.ServicesResult)
			})
		})
		Context("with an error returned by the backend", func() {
			It("should return the error and a zero result", func() {
				client := newRetryableSlmpClient(404, nil)
				result, err := client.GetServices()
				testutil.ExpectErrorAndZeroResult(err, result)
			})
		})
		Context("with an 500 server error returned by the backend", func() {
			It("should return the error, zero result and retry expectedRetriesCount times", func() {
				client := newRetryableSlmpClient(500, nil)
				retryableSlmpClient := client.(slmpclient.RetryableSlmpClient)
				result, err := retryableSlmpClient.GetServices()
				testutil.ExpectErrorAndZeroResult(err, result)
				mockSlmpClient := retryableSlmpClient.SlmpClient.(*MockSlmpClient)
				Expect(mockSlmpClient.GetRetriesCount()).To(Equal(expectedRetriesCount))
			})
		})
		Context("with an 502 server error returned by the backend", func() {
			It("should return the error, zero result and retry expectedRetriesCount times", func() {
				client := newRetryableSlmpClient(502, nil)
				retryableSlmpClient := client.(slmpclient.RetryableSlmpClient)
				result, err := retryableSlmpClient.GetServices()
				testutil.ExpectErrorAndZeroResult(err, result)
				mockSlmpClient := retryableSlmpClient.SlmpClient.(*MockSlmpClient)
				Expect(mockSlmpClient.GetRetriesCount()).To(Equal(expectedRetriesCount))
			})
		})
	})

	Describe("GetService", func() {
		Context("with a valid service returned by the backend", func() {
			It("should return the specific service returned by the backend", func() {
				client := newRetryableSlmpClient(200, testutil.ServiceResult)
				result, err := client.GetService(testutil.ServiceID)
				testutil.ExpectNoErrorAndResult(err, result, &testutil.ServiceResult)
			})
		})
		Context("with an error returned by the backend", func() {
			It("should return the error and a zero result", func() {
				client := newRetryableSlmpClient(404, nil)
				result, err := client.GetService(testutil.ServiceID)
				testutil.ExpectErrorAndZeroResult(err, result)
			})
		})
		Context("with an 500 server error returned by the backend", func() {
			It("should return the error, zero result and retry expectedRetriesCount times", func() {
				client := newRetryableSlmpClient(500, nil)
				retryableSlmpClient := client.(slmpclient.RetryableSlmpClient)
				result, err := retryableSlmpClient.GetService(testutil.ServiceID)
				testutil.ExpectErrorAndZeroResult(err, result)
				mockSlmpClient := retryableSlmpClient.SlmpClient.(*MockSlmpClient)
				Expect(mockSlmpClient.GetRetriesCount()).To(Equal(expectedRetriesCount))
			})
		})
		Context("with an 502 server error returned by the backend", func() {
			It("should return the error, zero result and retry expectedRetriesCount times", func() {
				client := newRetryableSlmpClient(502, nil)
				retryableSlmpClient := client.(slmpclient.RetryableSlmpClient)
				result, err := retryableSlmpClient.GetService(testutil.ServiceID)
				testutil.ExpectErrorAndZeroResult(err, result)
				mockSlmpClient := retryableSlmpClient.SlmpClient.(*MockSlmpClient)
				Expect(mockSlmpClient.GetRetriesCount()).To(Equal(expectedRetriesCount))
			})
		})
	})

	Describe("GetServiceProcesses", func() {
		Context("with valid service processes returned by the backend", func() {
			It("should return all service processes returned by the backend", func() {
				client := newRetryableSlmpClient(200, testutil.ProcessesResult)
				result, err := client.GetServiceProcesses(testutil.ServiceID)
				testutil.ExpectNoErrorAndResult(err, result, testutil.ProcessesResult)
			})
		})
		Context("with an error returned by the backend", func() {
			It("should return the error and a zero result", func() {
				client := newRetryableSlmpClient(404, nil)
				result, err := client.GetServiceProcesses(testutil.ServiceID)
				testutil.ExpectErrorAndZeroResult(err, result)
			})
		})
		Context("with an 500 server error returned by the backend", func() {
			It("should return the error, zero result and retry expectedRetriesCount times", func() {
				client := newRetryableSlmpClient(500, nil)
				retryableSlmpClient := client.(slmpclient.RetryableSlmpClient)
				result, err := retryableSlmpClient.GetServiceProcesses(testutil.ServiceID)
				testutil.ExpectErrorAndZeroResult(err, result)
				mockSlmpClient := retryableSlmpClient.SlmpClient.(*MockSlmpClient)
				Expect(mockSlmpClient.GetRetriesCount()).To(Equal(expectedRetriesCount))
			})
		})
		Context("with an 502 server error returned by the backend", func() {
			It("should return the error, zero result and retry expectedRetriesCount times", func() {
				client := newRetryableSlmpClient(502, nil)
				retryableSlmpClient := client.(slmpclient.RetryableSlmpClient)
				result, err := retryableSlmpClient.GetServiceProcesses(testutil.ServiceID)
				testutil.ExpectErrorAndZeroResult(err, result)
				mockSlmpClient := retryableSlmpClient.SlmpClient.(*MockSlmpClient)
				Expect(mockSlmpClient.GetRetriesCount()).To(Equal(expectedRetriesCount))
			})
		})
	})

	Describe("GetServiceFiles", func() {
		Context("with valid service files returned by the backend", func() {
			It("should return all service files returned by the backend", func() {
				client := newRetryableSlmpClient(200, testutil.FilesResult)
				result, err := client.GetServiceFiles("xs2-deploy")
				testutil.ExpectNoErrorAndResult(err, result, testutil.FilesResult)
			})
		})
		Context("with an error returned by the backend", func() {
			It("should return the error and a zero result", func() {
				client := newRetryableSlmpClient(404, nil)
				result, err := client.GetServiceFiles("xs2-deploy")
				testutil.ExpectErrorAndZeroResult(err, result)
			})
		})
		Context("with an 500 server error returned by the backend", func() {
			It("should return the error, zero result and retry expectedRetriesCount times", func() {
				client := newRetryableSlmpClient(500, nil)
				retryableSlmpClient := client.(slmpclient.RetryableSlmpClient)
				result, err := retryableSlmpClient.GetServiceFiles("xs2-deploy")
				testutil.ExpectErrorAndZeroResult(err, result)
				mockSlmpClient := retryableSlmpClient.SlmpClient.(*MockSlmpClient)
				Expect(mockSlmpClient.GetRetriesCount()).To(Equal(expectedRetriesCount))
			})
		})
		Context("with an 502 server error returned by the backend", func() {
			It("should return the error, zero result and retry expectedRetriesCount times", func() {
				client := newRetryableSlmpClient(502, nil)
				retryableSlmpClient := client.(slmpclient.RetryableSlmpClient)
				result, err := retryableSlmpClient.GetServiceFiles("xs2-deploy")
				testutil.ExpectErrorAndZeroResult(err, result)
				mockSlmpClient := retryableSlmpClient.SlmpClient.(*MockSlmpClient)
				Expect(mockSlmpClient.GetRetriesCount()).To(Equal(expectedRetriesCount))
			})
		})
	})

	Describe("CreateServiceProcess", func() {
		Context("with a valid newly created process retuned by the backend", func() {
			It("should return the process returned by the backend", func() {
				client := newRetryableSlmpClient(200, testutil.ProcessResult)
				result, err := client.CreateServiceProcess(testutil.ServiceID, &testutil.ProcessResult)
				testutil.ExpectNoErrorAndResult(err, result, &testutil.ProcessResult)
			})
		})
		Context("with an error returned by the backend", func() {
			It("should return the error and a zero result", func() {
				client := newRetryableSlmpClient(404, nil)
				result, err := client.CreateServiceProcess(testutil.ServiceID, &testutil.ProcessResult)
				testutil.ExpectErrorAndZeroResult(err, result)
			})
		})
		Context("with an 500 server error returned by the backend", func() {
			It("should return the error, zero result and retry expectedRetriesCount times", func() {
				client := newRetryableSlmpClient(500, nil)
				retryableSlmpClient := client.(slmpclient.RetryableSlmpClient)
				result, err := retryableSlmpClient.CreateServiceProcess(testutil.ServiceID, &testutil.ProcessResult)
				testutil.ExpectErrorAndZeroResult(err, result)
				mockSlmpClient := retryableSlmpClient.SlmpClient.(*MockSlmpClient)
				Expect(mockSlmpClient.GetRetriesCount()).To(Equal(expectedRetriesCount))
			})
		})
		Context("with an 502 server error returned by the backend", func() {
			It("should return the error, zero result and retry expectedRetriesCount times", func() {
				client := newRetryableSlmpClient(502, nil)
				retryableSlmpClient := client.(slmpclient.RetryableSlmpClient)
				result, err := retryableSlmpClient.CreateServiceProcess(testutil.ServiceID, &testutil.ProcessResult)
				testutil.ExpectErrorAndZeroResult(err, result)
				mockSlmpClient := retryableSlmpClient.SlmpClient.(*MockSlmpClient)
				Expect(mockSlmpClient.GetRetriesCount()).To(Equal(expectedRetriesCount))
			})
		})
	})
})

func newRetryableSlmpClient(statusCode int, v interface{}) slmpclient.SlmpClientOperations {
	tokenFactory := baseclient.NewCustomTokenFactory("test-token")
	cookieJar, _ := cookiejar.New(nil)
	roundTripper := testutil.NewCustomTransport(statusCode, v)
	return NewMockRetryableSlmpClient("http://localhost:1000", "test-org", "test-space", roundTripper, cookieJar, tokenFactory)
}

// NewMockRetryableSlmpClient creates a new retryable REST client
func NewMockRetryableSlmpClient(host, org, space string, rt http.RoundTripper, jar http.CookieJar, tokenFactory baseclient.TokenFactory) slmpclient.SlmpClientOperations {
	slmpClient := NewMockSlmpClient(host, org, space, rt, jar, tokenFactory)
	return slmpclient.RetryableSlmpClient{SlmpClient: slmpClient, MaxRetriesCount: 3, RetryInterval: time.Microsecond * 1}
}

type MockSlmpClient struct {
	slmpClient   slmpclient.SlmpClientOperations
	retriesCount int
}

// NewMockSlmpClient creates a new Rest client
func NewMockSlmpClient(host, org, space string, rt http.RoundTripper, jar http.CookieJar, tokenFactory baseclient.TokenFactory) slmpclient.SlmpClientOperations {
	slmpClient := slmpclient.NewSlmpClient(host, org, space, rt, jar, tokenFactory)
	return &MockSlmpClient{slmpClient, 0}
}

// GetMetadata retrieves the SLMP metadata
func (c *MockSlmpClient) GetMetadata() (*models.Metadata, error) {
	c.retriesCount++
	return c.slmpClient.GetMetadata()
}

// GetServices retrieves all services
func (c *MockSlmpClient) GetServices() (models.Services, error) {
	c.retriesCount++
	return c.slmpClient.GetServices()
}

// GetService retrieves the service with the specified service ID
func (c *MockSlmpClient) GetService(serviceID string) (*models.Service, error) {
	c.retriesCount++
	return c.slmpClient.GetService(serviceID)
}

// GetServiceProcesses retrieves all processes for the service with the specified service ID
func (c *MockSlmpClient) GetServiceProcesses(serviceID string) (models.Processes, error) {
	c.retriesCount++
	return c.slmpClient.GetServiceProcesses(serviceID)
}

func (c *MockSlmpClient) GetProcess(processID string) (*models.Process, error) {
	c.retriesCount++
	return c.slmpClient.GetProcess(processID)
}

// GetServiceFiles retrieves all files for the service with the specified service ID
func (c *MockSlmpClient) GetServiceFiles(serviceID string) (models.Files, error) {
	c.retriesCount++
	return c.slmpClient.GetServiceFiles(serviceID)
}

// CreateServiceFile uploads a file for the service with the specified service ID
func (c *MockSlmpClient) CreateServiceFile(serviceID string, file os.File) (models.Files, error) {
	c.retriesCount++
	return c.slmpClient.CreateServiceFile(serviceID, file)
}

// CreateServiceProcess creates a process for the service with the specified service ID
func (c *MockSlmpClient) CreateServiceProcess(serviceID string, process *models.Process) (*models.Process, error) {
	c.retriesCount++
	return c.slmpClient.CreateServiceProcess(serviceID, process)
}

// GetServiceVersions retrieves the versions for the service with the specified service ID
func (c *MockSlmpClient) GetServiceVersions(serviceID string) (models.Versions, error) {
	c.retriesCount++
	return c.slmpClient.GetServiceVersions(serviceID)
}

// GetRetriesCount returns retries count
func (c *MockSlmpClient) GetRetriesCount() int {
	return c.retriesCount
}
