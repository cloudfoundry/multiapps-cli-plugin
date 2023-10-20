package configuration

import (
	"fmt"
	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/configuration/properties"
	"os"
)

const unknownError = "An unknown error occurred during the parsing of the environment variable \"%s\". Please report this! Value type: %T"

type Snapshot struct {
	backendURL               properties.ConfigurableProperty
	uploadChunkSizeInMB      properties.ConfigurableProperty
	uploadChunksSequentially properties.ConfigurableProperty
}

func NewSnapshot() Snapshot {
	return Snapshot{
		backendURL:               properties.BackendURL,
		uploadChunkSizeInMB:      properties.UploadChunkSizeInMB,
		uploadChunksSequentially: properties.UploadChunksSequentially,
	}
}

func (c Snapshot) GetBackendURL() string {
	return getStringProperty(c.backendURL)
}

func (c Snapshot) GetUploadChunkSizeInMB() uint64 {
	return getUint64Property(c.uploadChunkSizeInMB)
}

func (c Snapshot) GetUploadChunksSequentially() bool {
	return getBoolProperty(c.uploadChunksSequentially)
}

func getStringProperty(property properties.ConfigurableProperty) string {
	uncastedValue := getPropertyOrDefault(property)
	value, ok := uncastedValue.(string)
	if !ok {
		panic(fmt.Sprintf(unknownError, property.Name, uncastedValue))
	}
	return value
}

func getUint64Property(property properties.ConfigurableProperty) uint64 {
	uncastedValue := getPropertyOrDefault(property)
	value, ok := uncastedValue.(uint64)
	if !ok {
		panic(fmt.Sprintf(unknownError, property.Name, uncastedValue))
	}
	return value
}

func getBoolProperty(property properties.ConfigurableProperty) bool {
	uncastedValue := getPropertyOrDefault(property)
	value, ok := uncastedValue.(bool)
	if !ok {
		panic(fmt.Sprintf(unknownError, property.Name, uncastedValue))
	}
	return value
}

func getPropertyOrDefault(property properties.ConfigurableProperty) interface{} {
	value := getPropertyWithNameOrDefaultIfInvalid(property, property.Name)
	if value != nil {
		return value
	}
	return property.DefaultValue
}

func getPropertyWithNameOrDefaultIfInvalid(property properties.ConfigurableProperty, name string) interface{} {
	propertyValue, err := getPropertyWithName(name, property.Parser)
	if err != nil {
		propertyValue = os.Getenv(name)
		fmt.Printf(property.ParsingFailureMessage, propertyValue, name, property.DefaultValue)
		return property.DefaultValue
	}
	if propertyValue != nil {
		fmt.Printf(property.ParsingSuccessMessage, propertyValue, name)
		return propertyValue
	}
	return nil
}

func getPropertyWithName(name string, parser properties.Parser) (interface{}, error) {
	propertyValue := os.Getenv(name)
	if propertyValue != "" {
		return parser.Parse(propertyValue)
	}
	return nil, nil
}
