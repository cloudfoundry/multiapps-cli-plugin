package baseclient

import (
	"bytes"
	"fmt"

	strfmt "github.com/go-openapi/strfmt"

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
	if err == nil {
		return nil
	}
	if ae, ok := err.(*runtime.APIError); ok {
		response := ae.Response.(runtime.ClientResponse)
		return &ClientError{Code: ae.Code, Status: response.Message(), Description: response.Message()}
	}
	if response, ok := err.(*ErrorResponse); ok {
		return &ClientError{Code: response.Code, Status: response.Status, Description: response.Payload}
	}
	return err
}

func BuildErrorResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {
	result := &ErrorResponse{
		Code:   response.Code(),
		Status: response.Message(), // this isn't the body!
	}
	if err := result.readResponse(response, consumer, formats); err != nil {
		return err
	}
	return result
}

// ErrorResponse handles error cases
type ErrorResponse struct {
	Code    int
	Status  string
	Payload string
}

func (ce *ErrorResponse) Error() string {
	return fmt.Sprintf("%s (status %d): %v ", ce.Status, ce.Code, ce.Payload)
}

func (o *ErrorResponse) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {
	buf := new(bytes.Buffer)
	_, err := buf.ReadFrom(response.Body())
	if err != nil {
		return runtime.NewAPIError("unknown error", response, response.Code())
	}
	o.Payload = buf.String()

	return nil
}
