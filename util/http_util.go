package util

import (
	"net/http"
)

type HttpSimpleGetExecutor interface {
	ExecuteGetRequest(url string) (int, error)
}

type SimpleGetExecutor struct {
}

func (executor SimpleGetExecutor) ExecuteGetRequest(url string) (int, error) {
	resp, err := http.Get(url)
	if err != nil {
		return -1, err
	}
	defer resp.Body.Close()
	return resp.StatusCode, nil
}
