package secure_parameters

import (
	"errors"

	"gopkg.in/yaml.v3"
)

func BuildSecureExtension(parameters map[string]ParameterValue, mtaID string, schemaVersion string) ([]byte, error) {
	if len(parameters) == 0 {
		return nil, errors.New("no secure parameters collected")
	}

	if mtaID == "" {
		return nil, errors.New("mtaID is required for the secure extension descriptor's field 'extends'")
	}

	if schemaVersion == "" {
		return nil, errors.New("schemaVersion is required for the secure extension descriptor")
	}

	secureExtensionDescriptor := map[string]interface{}{
		"_schema-version": schemaVersion,
		"ID":              "__mta.secure",
		"extends":         mtaID,
		"parameters":      map[string]interface{}{},
	}

	parametersDescriptor := secureExtensionDescriptor["parameters"].(map[string]interface{})
	for name, currentParameterValue := range parameters {
		parametersDescriptor[name] = getValue(&currentParameterValue)
	}

	return yaml.Marshal(secureExtensionDescriptor)
}
