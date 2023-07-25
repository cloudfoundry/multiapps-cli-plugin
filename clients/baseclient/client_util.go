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
		retryErr, ok := err.(*RetryAfterError)
		if ok {
			ui.Warn("Retryable error occurred. Retrying after %s", retryErr.Duration)
			time.Sleep(retryErr.Duration)
			continue
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
	ae, ok := err.(*ClientError)
	if ok {
		httpCode := ae.Code
		httpCodeMajorDigit := httpCode / 100
		return httpCodeMajorDigit != 2
	}
	return true
}
