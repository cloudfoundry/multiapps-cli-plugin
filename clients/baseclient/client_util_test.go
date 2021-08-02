package baseclient

import (
	"fmt"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("ClientUtil", func() {
	Describe("ShouldRetry", func() {
		Context("when zero value pointer error(uninitialized value - nil) is passed as argument", func() {
			It("Retry operation call shouldn't be made", func() {
				result := shouldRetry(nil)
				Expect(result).To(Equal(false))
			})
		})

		Context("when the backend returns status code with first digit 1 (1xx)", func() {
			It("Retry operation call is expected to be made", func() {
				err := ClientError{102, "Processing", "The request has been accepted but it's not yet complete"}
				result := shouldRetry(&err)
				Expect(result).To(Equal(true))
			})
		})

		Context("when the backend returns status code with first digit 2 (2xx)", func() {
			It("Retry operation call shouldn't be made", func() {
				err := ClientError{200, "OK", "The request has succeeded"}
				result := shouldRetry(&err)
				Expect(result).To(Equal(false))
			})
		})

		Context("when the backend returns status code with first digit 3 (3xx)", func() {
			It("Retry operation call is expected to be made", func() {
				err := ClientError{301, "Moved Permanently", "URI of requested resource has been changed"}
				result := shouldRetry(&err)
				Expect(result).To(Equal(true))
			})
		})

		Context("when the backend returns status code with first digit 4 (4xx)", func() {
			It("Retry operation call is expected to be made", func() {
				err := ClientError{404, "Not Found", "The server cannot find requested resource"}
				result := shouldRetry(&err)
				Expect(result).To(Equal(true))
			})
		})

		Context("when the backend returns status code with first digit 5 (5xx)", func() {
			It("Retry operation call is expected to be made", func() {
				err := ClientError{500, "Internal Server Error", "The server got an invalid response"}
				result := shouldRetry(&err)
				Expect(result).To(Equal(true))
			})
		})

		Context("when the passed error is not of type 'ClientError'", func() {
			It("Retry operation call is expected to be made", func() {
				err := MockError{999, "Not ClientError"}
				result := shouldRetry(&err)
				Expect(result).To(Equal(true))
			})
		})
	})
})

var _ = Describe("ClientUtil", func() {
	Describe("CallWithRetry", func() {
		Context("when retry operation call is Not needed", func() {
			It("Just retrun required response", func() {
				getInterface := func() (interface{}, error) {
					toReturn := testStruct{}
					return toReturn, nil
				}
				result, err := CallWithRetry(getInterface, 4, (time.Duration(0) * time.Second))
				Expect(err).NotTo(HaveOccurred())
				Expect(result).To(Equal(result))
			})
		})
		Context("when retry operation call IS needed", func() {
			It("Callback with retry", func() {
				mockCallback := func() (interface{}, error) {
					toReturn := testStruct{}
					err := ClientError{404, "Not Found", "The server cannot find requesred resource"}
					return toReturn, &err
				}
				result, err := CallWithRetry(mockCallback, 4, (time.Duration(0) * time.Second))
				Expect(err).To(HaveOccurred())
				Expect(result).To(Equal(result))
			})
		})
	})
})

type MockError struct {
	Code   int
	Status string
}

func (ce *MockError) Error() string {
	return fmt.Sprintf("%s (status %d)", ce.Status, ce.Code)
}

type testStruct struct {
}
