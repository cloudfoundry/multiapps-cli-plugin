package slmpclient_test

import (
	"net/http/cookiejar"

	. "github.com/onsi/ginkgo"
	baseclient "github.com/SAP/cf-mta-plugin/clients/baseclient"
	slmpclient "github.com/SAP/cf-mta-plugin/clients/slmpclient"
	"github.com/SAP/cf-mta-plugin/testutil"
)

var _ = Describe("SlmpClient", func() {

	Describe("GetMetadata", func() {
		Context("with valid metadata returned by the backend", func() {
			It("should return the metadata returned by the backend", func() {
				client := newSlmpClient(200, testutil.SlmpMetadataResult)
				result, err := client.GetMetadata()
				testutil.ExpectNoErrorAndResult(err, result, &testutil.SlmpMetadataResult)
			})
		})
		Context("with an error returned by the backend", func() {
			It("should return the error and a zero result", func() {
				client := newSlmpClient(404, nil)
				result, err := client.GetMetadata()
				testutil.ExpectErrorAndZeroResult(err, result)
			})
		})
	})

	Describe("GetServices", func() {
		Context("with valid services returned by the backend", func() {
			It("should return all services returned by the backend", func() {
				client := newSlmpClient(200, testutil.ServicesResult)
				result, err := client.GetServices()
				testutil.ExpectNoErrorAndResult(err, result, testutil.ServicesResult)
			})
		})
		Context("with an error returned by the backend", func() {
			It("should return the error and a zero result", func() {
				client := newSlmpClient(404, nil)
				result, err := client.GetServices()
				testutil.ExpectErrorAndZeroResult(err, result)
			})
		})
	})

	Describe("GetService", func() {
		Context("with a valid service returned by the backend", func() {
			It("should return the specific service returned by the backend", func() {
				client := newSlmpClient(200, testutil.ServiceResult)
				result, err := client.GetService(testutil.ServiceID)
				testutil.ExpectNoErrorAndResult(err, result, &testutil.ServiceResult)
			})
		})
		Context("with an error returned by the backend", func() {
			It("should return the error and a zero result", func() {
				client := newSlmpClient(404, nil)
				result, err := client.GetService(testutil.ServiceID)
				testutil.ExpectErrorAndZeroResult(err, result)
			})
		})
	})

	Describe("GetServiceProcesses", func() {
		Context("with valid service processes returned by the backend", func() {
			It("should return all service processes returned by the backend", func() {
				client := newSlmpClient(200, testutil.ProcessesResult)
				result, err := client.GetServiceProcesses(testutil.ServiceID)
				testutil.ExpectNoErrorAndResult(err, result, testutil.ProcessesResult)
			})
		})
		Context("with an error returned by the backend", func() {
			It("should return the error and a zero result", func() {
				client := newSlmpClient(404, nil)
				result, err := client.GetServiceProcesses(testutil.ServiceID)
				testutil.ExpectErrorAndZeroResult(err, result)
			})
		})
	})

	Describe("GetServiceFiles", func() {
		Context("with valid service files returned by the backend", func() {
			It("should return all service files returned by the backend", func() {
				client := newSlmpClient(200, testutil.FilesResult)
				result, err := client.GetServiceFiles("xs2-deploy")
				testutil.ExpectNoErrorAndResult(err, result, testutil.FilesResult)
			})
		})
		Context("with an error returned by the backend", func() {
			It("should return the error and a zero result", func() {
				client := newSlmpClient(404, nil)
				result, err := client.GetServiceFiles("xs2-deploy")
				testutil.ExpectErrorAndZeroResult(err, result)
			})
		})
	})

	Describe("CreateServiceProcess", func() {
		Context("with a valid newly created process retuned by the backend", func() {
			It("should return the process returned by the backend", func() {
				client := newSlmpClient(200, testutil.ProcessResult)
				result, err := client.CreateServiceProcess(testutil.ServiceID, &testutil.ProcessResult)
				testutil.ExpectNoErrorAndResult(err, result, &testutil.ProcessResult)
			})
		})
		Context("with an error returned by the backend", func() {
			It("should return the error and a zero result", func() {
				client := newSlmpClient(404, nil)
				result, err := client.CreateServiceProcess(testutil.ServiceID, &testutil.ProcessResult)
				testutil.ExpectErrorAndZeroResult(err, result)
			})
		})
	})
})

func newSlmpClient(statusCode int, v interface{}) slmpclient.SlmpClientOperations {
	tokenFactory := baseclient.NewCustomTokenFactory("test-token")
	cookieJar, _ := cookiejar.New(nil)
	roundTripper := testutil.NewCustomTransport(statusCode, v)
	return slmpclient.NewSlmpClient("http://localhost:1000", "test-org", "test-space", roundTripper, cookieJar, tokenFactory)
}
