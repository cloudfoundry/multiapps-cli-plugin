package baseclient

import (
	"time"
)

// CallWithRetry executes callback with retry
func CallWithRetry(callback func() (interface{}, error), maxRetriesCount int, retryInterval time.Duration) (interface{}, error) {
	for index := 0; index < maxRetriesCount; index++ {
		resp, err := callback()
		if !shouldRetry(err) {
			return resp, err
		}
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
