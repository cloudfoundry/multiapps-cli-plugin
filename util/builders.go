package util

import (
	"github.com/SAP/cf-mta-plugin/clients/models"
)

// ProcessBuilder a builder for models.Process instances
type ProcessBuilder struct {
	operation models.Operation
}

// NewProcessBuilder creates a new process builder
func NewProcessBuilder() *ProcessBuilder {
	return &ProcessBuilder{operation: models.Operation{Parameters: make(map[string]interface{})}}
}

func (pb *ProcessBuilder) ProcessType(processType string) *ProcessBuilder {
	pb.operation.ProcessType = processType
	return pb
}

// Parameter adds a parameter to the process if it is set
func (pb *ProcessBuilder) Parameter(parameterID string, value string) *ProcessBuilder {
	if value != "" {
		pb.operation.Parameters[parameterID] = value
	}
	return pb
}

// Build builds the process
func (pb *ProcessBuilder) Build() *models.Operation {
	return &pb.operation
}
