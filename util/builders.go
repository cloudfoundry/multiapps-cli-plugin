package util

import (
	"bytes"
	"fmt"
	"net/url"
	"strings"

	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/clients/models"
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

func (pb *ProcessBuilder) Namespace(namespace string) *ProcessBuilder {
	pb.operation.Namespace = namespace
	return pb
}

// Parameter adds a parameter to the process if it is set
func (pb *ProcessBuilder) Parameter(parameterID string, value string) *ProcessBuilder {
	if value != "" {
		pb.operation.Parameters[parameterID] = value
	}
	return pb
}

// SetParameterWithoutCheck sets the parameter without checking whether it is null
func (pb *ProcessBuilder) SetParameterWithoutCheck(parameterID string, value string) *ProcessBuilder {
	pb.operation.Parameters[parameterID] = value
	return pb
}

// Build builds the process
func (pb *ProcessBuilder) Build() *models.Operation {
	return &pb.operation
}

const hostAndSchemeSeparator = "://"

type UriBuilder struct {
	host   string
	scheme string
	path   string
}

func NewUriBuilder() *UriBuilder {
	return &UriBuilder{host: "", scheme: "", path: ""}
}

func (builder *UriBuilder) SetScheme(scheme string) *UriBuilder {
	builder.scheme = scheme
	return builder
}

func (builder *UriBuilder) SetHost(host string) *UriBuilder {
	builder.host = host
	return builder
}

func (builder *UriBuilder) SetPath(path string) *UriBuilder {
	builder.path = path
	return builder
}

func (builder *UriBuilder) Build() (string, error) {
	if builder.scheme == "" || builder.host == "" {
		return "", fmt.Errorf("The host or scheme could not be empty")
	}
	stringBuilder := bytes.Buffer{}
	stringBuilder.WriteString(builder.scheme)
	stringBuilder.WriteString(hostAndSchemeSeparator)
	stringBuilder.WriteString(builder.host)
	stringBuilder.WriteString(getPath(builder.path))
	builtUrl := stringBuilder.String()
	parsedUrl, err := url.Parse(builtUrl)
	if err != nil {
		return "", err
	}
	return parsedUrl.String(), nil
}

func getPath(path string) string {
	if strings.HasPrefix(path, "/") {
		return path
	}

	return "/" + path
}
