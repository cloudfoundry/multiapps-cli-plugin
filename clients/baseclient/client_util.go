package baseclient

import (
	"net/url"
	"regexp"
	"strings"
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
	if isMatching(err) {
		return true
	}
	ae, ok := err.(*ClientError)
	if ok {
		httpCode := ae.Code
		httpCodeMajorDigit := httpCode / 100
		if httpCodeMajorDigit == 5 || httpCodeMajorDigit == 4 {
			return true
		}
	}
	return false
}

func isMatching(err error) bool {
	return strings.Contains(err.Error(), "retry is needed") || isErrorEOF(err)
}

func isErrorEOF(err error) bool {
	isMatching, _ := regexp.MatchString(" EOF$", err.Error())

	return isMatching
}

func EncodeArg(arg string) string {
	return strings.Replace(url.QueryEscape(arg), "+", "%20", -1)
}
