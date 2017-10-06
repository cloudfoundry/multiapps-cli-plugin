package baseclient

import (
	"fmt"

	"github.com/go-openapi/runtime"
)

type ClientError struct {
	Code        int
	Status      string
	Description interface{}
}

func (ce *ClientError) Error() string {
	return fmt.Sprintf("%s (status %d): %v ", ce.Status, ce.Code, ce.Description)
}

func NewClientError(err error) error {
	ae, ok := err.(*runtime.APIError)
	if ok {
		resp := ae.Response.(runtime.ClientResponse)
		return &ClientError{Code: ae.Code, Status: resp.Message(), Description: resp.Message()}
	}
	return err
}
