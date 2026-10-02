package baseclient

import (
	"time"

	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/ui"
)

// CallWithRetry executes callback with retry
func CallWithRetry(callback func() (interface{}, error), maxRetriesCount int, retryInterval time.Duration) (interface{}, error) {
	for index := 0; index < maxRetriesCount; index++ {
		resp, err := callback()
		if !shouldRetry(err) {
			return resp, err
		}
		ui.Warn("Error occurred: %s. Retrying after: %s.", err.Error(), retryInterval)
		time.Sleep(retryInterval)
	}
	return callback()
}

func shouldRetry(err error) bool {
	if err == nil {
		return false
	}
	// A rate-limit (HTTP 429) response is not retried: the operation stops
	// immediately so the caller can surface the server's Retry-After hint.
	if _, ok := err.(*RetryAfterError); ok {
		return false
	}
	ae, ok := err.(*ClientError)
	if ok {
		httpCode := ae.Code
		httpCodeMajorDigit := httpCode / 100
		return httpCodeMajorDigit != 2
	}
	return true
}
