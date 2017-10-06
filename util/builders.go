package util

import (
	"github.com/go-openapi/strfmt"
	"github.com/SAP/cf-mta-plugin/clients/models"
)

// ProcessBuilder a builder for models.Process instances
type ProcessBuilder struct {
	process models.Process
}

// NewProcessBuilder creates a new process builder
func NewProcessBuilder() *ProcessBuilder {
	return &ProcessBuilder{}
}

func (pb *ProcessBuilder) ServiceID(serviceID string) *ProcessBuilder {
	uri := strfmt.URI(serviceID)
	pb.process.Service = &uri
	return pb
}

// ProcessParameter adds a parameter to the process if it is set
func (pb *ProcessBuilder) Parameter(parameterID string, value string) *ProcessBuilder {
	if value != "" {
		scalarType := models.SlpParameterType("slp.parameter.type.SCALAR")
		processParam := models.Parameter{ID: &parameterID, Value: value, Type: scalarType}
		pb.process.Parameters.Parameters = append(pb.process.Parameters.Parameters, &processParam)
	}
	return pb
}

// Build builds the process
func (pb *ProcessBuilder) Build() *models.Process {
	return &pb.process
}
