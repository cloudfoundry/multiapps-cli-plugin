package retrier

import (
	"time"

	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/log"
)

func Execute[T any](attempts int, callback func() (T, error), shouldRetry func(result T, err error) bool) (T, error) {
	for i := 0; i < attempts; i++ {
		result, err := callback()
		if shouldRetry(result, err) {
			logError[T](result, err)
			time.Sleep(3 * time.Second)
			continue
		}
		return result, err
	}
	return callback()
}

func logError[T any](result T, err error) {
	if err != nil {
		log.Tracef("retrying an operation that failed with: %v", err)
	}
	log.Tracef("result of the callback %v", result)
}
