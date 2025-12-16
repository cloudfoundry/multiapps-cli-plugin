package secure_parameters

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"

	"gopkg.in/yaml.v3"
)

var nameRegex = regexp.MustCompile(`^[A-Za-z0-9._-]+$`)

type typeOfValue int

const (
	typeString typeOfValue = iota
	typeJSON
	typeMultiline
)

type ParameterValue struct {
	Type          typeOfValue
	StringContent string
	ObjectContent map[string]interface{}
}

func nameDuplicated(name, prefix string, result map[string]ParameterValue) error {
	_, ok := result[name]
	if ok {
		return fmt.Errorf("secure parameter %q defined multiple ways (collision with %s)", name, prefix)
	}

	return nil
}

func CollectFromEnv(prefix string) (map[string]ParameterValue, error) {
	plainValue := prefix + "___"
	jsonValue := prefix + "_JSON___"
	certificateValue := prefix + "_CERT___" //X509value beacuse the certiciates are of type X509 (should be renamed)

	result := map[string]ParameterValue{}

	for _, nameValuePair := range os.Environ() {
		equalsIndex := strings.IndexByte(nameValuePair, '=')
		if equalsIndex < 0 {
			continue
		}
		envName := nameValuePair[:equalsIndex]
		envValue := nameValuePair[equalsIndex+1:]

		var name string

		switch {
		case strings.HasPrefix(envName, jsonValue):
			name = strings.TrimPrefix(envName, jsonValue)

			if !nameRegex.MatchString(name) {
				return nil, fmt.Errorf("invalid secure parameter name %q", name)
			}

			err := nameDuplicated(name, "__MTA_JSON", result)
			if err != nil {
				return nil, err
			}

			var jsonObject map[string]interface{}

			err2 := json.Unmarshal([]byte(envValue), &jsonObject)
			if err2 != nil {
				return nil, fmt.Errorf("invalid JSON for %s: %w", name, err2)
			}
			result[name] = ParameterValue{Type: typeJSON, ObjectContent: jsonObject}

		case strings.HasPrefix(envName, certificateValue):
			name = strings.TrimPrefix(envName, certificateValue)

			if !nameRegex.MatchString(name) {
				return nil, fmt.Errorf("invalid secure parameter name %q", name)
			}

			err := nameDuplicated(name, "__MTA_CERT", result)
			if err != nil {
				return nil, err
			}

			decoded, err := base64.StdEncoding.DecodeString(envValue)
			if err != nil {
				return nil, fmt.Errorf("invalid base64 for %s: %w", name, err)
			}
			result[name] = ParameterValue{Type: typeMultiline, StringContent: string(decoded)}

		case strings.HasPrefix(envName, plainValue):
			name = strings.TrimPrefix(envName, plainValue)

			if !nameRegex.MatchString(name) {
				return nil, fmt.Errorf("invalid secure parameter name %q", name)
			}

			err := nameDuplicated(name, "__MTA", result)
			if err != nil {
				return nil, err
			}

			result[name] = ParameterValue{Type: typeString, StringContent: envValue}

		default:
			continue
		}
	}

	return result, nil
}

func BuildSecureExtension(parameters map[string]ParameterValue, mtaID string, schemaVersion string) ([]byte, error) {
	if len(parameters) == 0 {
		return nil, errors.New("no secure parameters collected")
	}

	if mtaID == "" {
		return nil, errors.New("mtaID is required for the extension descriptor's field 'extends'")
	}

	if schemaVersion == "" {
		schemaVersion = "3.3"
	}

	secureExtensionDescriptor := map[string]interface{}{
		"_schema-version": schemaVersion,
		"ID":              "__mta.secure",
		"extends":         mtaID,
		"parameters":      map[string]interface{}{},
	}

	parametersDescriptor := secureExtensionDescriptor["parameters"].(map[string]interface{})
	for name, currentParameterValue := range parameters {
		switch currentParameterValue.Type {
		case typeJSON:
			parametersDescriptor[name] = currentParameterValue.ObjectContent
		default:
			parametersDescriptor[name] = currentParameterValue.StringContent
		}
	}

	return yaml.Marshal(secureExtensionDescriptor)
}
