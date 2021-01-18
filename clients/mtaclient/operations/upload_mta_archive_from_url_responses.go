// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"

	"github.com/go-openapi/strfmt"

	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/clients/baseclient"
	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/clients/models"
)

// UploadMtaArchiveFromUrlReader is a Reader for the UploadMtaArchiveFromUrl structure.
type UploadMtaArchiveFromUrlReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *UploadMtaArchiveFromUrlReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {

	case 201:
		result := NewUploadMtaArchiveFromUrlCreated()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil

	default:
		return nil, baseclient.BuildErrorResponse(response, consumer, o.formats)
	}
}

// NewUploadMtaArchiveFromUrlCreated creates a UploadMtaArchiveFromUrlCreated with default headers values
func NewUploadMtaArchiveFromUrlCreated() *UploadMtaArchiveFromUrlCreated {
	return &UploadMtaArchiveFromUrlCreated{}
}

/*UploadMtaArchiveFromUrlCreated handles this case with default header values.

Created
*/
type UploadMtaArchiveFromUrlCreated struct {
	Payload *models.FileMetadata
}

func (o *UploadMtaArchiveFromUrlCreated) Error() string {
	return fmt.Sprintf("[POST /files][%d] UploadMtaArchiveFromUrlCreated  %+v", 201, o.Payload)
}

func (o *UploadMtaArchiveFromUrlCreated) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.FileMetadata)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

