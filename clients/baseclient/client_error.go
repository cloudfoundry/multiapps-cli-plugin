package baseclient

import (
	"bytes"
	"fmt"
	"time"

	"github.com/go-openapi/strfmt"

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
	ae, ok := err.(*runtime.APIError)
	if ok {
		response := ae.Response.(runtime.ClientResponse)
		return &ClientError{Code: ae.Code, Status: response.Message(), Description: response.Message()}
	}
	response, ok := err.(*ErrorResponse)
	if ok {
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

func (e *ErrorResponse) Error() string {
	return fmt.Sprintf("%s (status %d): %v ", e.Status, e.Code, e.Payload)
}

func (e *ErrorResponse) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {
	buf := new(bytes.Buffer)
	_, err := buf.ReadFrom(response.Body())
	if err != nil {
		return runtime.NewAPIError("unknown error", response, response.Code())
	}
	e.Payload = buf.String()

	return nil
}

type RetryAfterError struct {
	Duration time.Duration
}

func (e *RetryAfterError) Error() string {
	return "Retryable error: Retry-After " + e.Duration.String()
}
