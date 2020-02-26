package configuration

import (
	"fmt"
	"os"
)

const (
	unknownError         = "An unknown error occurred during the parsing of the environment variable \"%s\". Please report this! Value type: %T"
	DefaultChunkSizeInMB = uint64(45)
)

var BackendURLConfigurableProperty = configurableProperty{
	Name: "MULTIAPPS_CONTROLLER_URL",
	DeprecatedNames: []string{
		"DEPLOY_SERVICE_URL",
	},
	Parser:                noOpParser{},
	ParsingSuccessMessage: "Attention: You've specified a custom backend URL (%s) via the environment variable \"%s\". The application listening on that URL may be outdated, contain bugs or unreleased features or may even be modified by a potentially untrused person. Use at your own risk.\n",
	ParsingFailureMessage: "No validation implemented for custom backend URLs. If you're seeing this message then something has gone horribly wrong.\n",
	DefaultValue:          "",
}

var ChunkSizeInMBConfigurableProperty = configurableProperty{
	Name: "MULTIAPPS_UPLOAD_CHUNK_SIZE",
	DeprecatedNames: []string{
		"CHUNK_SIZE_IN_MB",
	},
	Parser:                chunkSizeInMBParser{},
	ParsingSuccessMessage: "Attention: You've specified a custom chunk size (%d MB) via the environment variable \"%s\".\n",
	ParsingFailureMessage: "Attention: You've specified an INVALID custom chunk size (%s) via the environment variable \"%s\". Using default: %d\n",
	DefaultValue:          DefaultChunkSizeInMB,
}

type Snapshot struct {
	backendURL    string
	chunkSizeInMB uint64
}

func NewSnapshot() Snapshot {
	return Snapshot{
		backendURL:    getBackendURLFromEnvironment(),
		chunkSizeInMB: getChunkSizeInMBFromEnvironment(),
	}
}

func (c Snapshot) GetBackendURL() string {
	return c.backendURL
}

func (c Snapshot) GetChunkSizeInMB() uint64 {
	return c.chunkSizeInMB
}

func getBackendURLFromEnvironment() string {
	return getStringProperty(BackendURLConfigurableProperty)
}

func getChunkSizeInMBFromEnvironment() uint64 {
	return getUint64Property(ChunkSizeInMBConfigurableProperty)
}

func getStringProperty(property configurableProperty) string {
	uncastedValue := getPropertyOrDefault(property)
	value, ok := uncastedValue.(string)
	if !ok {
		panic(fmt.Sprintf(unknownError, property.Name, uncastedValue))
	}
	return value
}

func getUint64Property(property configurableProperty) uint64 {
	uncastedValue := getPropertyOrDefault(property)
	value, ok := uncastedValue.(uint64)
	if !ok {
		panic(fmt.Sprintf(unknownError, property.Name, uncastedValue))
	}
	return value
}

func getPropertyOrDefault(property configurableProperty) interface{} {
	value := getPropertyWithNameOrDefaultIfInvalid(property, property.Name)
	if value != nil {
		return value
	}
	for _, deprecatedName := range property.DeprecatedNames {
		value := getPropertyWithNameOrDefaultIfInvalid(property, deprecatedName)
		if value != nil {
			fmt.Printf("Attention: You're using a deprecated environment variable \"%s\". Use \"%s\" instead.\n\n", deprecatedName, property.Name)
			return value
		}
	}
	return property.DefaultValue
}

func getPropertyWithNameOrDefaultIfInvalid(property configurableProperty, name string) interface{} {
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

func getPropertyWithName(name string, parser configurablePropertyParser) (interface{}, error) {
	propertyValue, isSet := os.LookupEnv(name)
	if isSet {
		return parser.Parse(propertyValue)
	}
	return nil, nil
}
