package mtaclient_v2

import (
	models "github.com/cloudfoundry-incubator/multiapps-cli-plugin/clients/models"
	strfmt "github.com/go-openapi/strfmt"
)

// interface for all operations provided by the v2 client
type MtaV2ClientOperations interface {
	GetMtas(name, namespace *string, spaceGuid string) ([]*models.Mta, error)
	GetMtasForThisSpace(name, namespace *string) ([]*models.Mta, error)
}

// ResponseHeader response header
type ResponseHeader struct {
	Location strfmt.URI
}
