package secure_parameters

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strings"
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
	JSONContent   interface{}
}

func validateNoDuplicatesExist(name, prefix string, result map[string]ParameterValue) error {
	_, ok := result[name]
	if ok {
		return fmt.Errorf("secure parameter %q defined multiple ways (collision with %s)", name, prefix)
	}

	return nil
}

func getValue(parameter *ParameterValue) interface{} {
	switch parameter.Type {
	case typeJSON:
		return parameter.JSONContent

	default:
		return parameter.StringContent
	}
}

func CollectFromEnv(prefix string) (map[string]ParameterValue, error) {
	plainPrefix := prefix + "___"
	jsonPrefix := prefix + "_JSON___"
	certPrefix := prefix + "_CERT___"

	result := make(map[string]ParameterValue)

	for _, nameValuePair := range os.Environ() {
		equalsIndex := strings.IndexByte(nameValuePair, '=')
		if equalsIndex < 0 {
			continue
		}
		envName := nameValuePair[:equalsIndex]
		envValue := nameValuePair[equalsIndex+1:]

		var name string

		switch {
		case strings.HasPrefix(envName, jsonPrefix):
			name = strings.TrimPrefix(envName, jsonPrefix)

			err := addJSONValues(name, envValue, result)
			if err != nil {
				return nil, err
			}
		case strings.HasPrefix(envName, certPrefix):
			name = strings.TrimPrefix(envName, certPrefix)

			err := addCertificateValues(name, envValue, result)
			if err != nil {
				return nil, err
			}
		case strings.HasPrefix(envName, plainPrefix):
			name = strings.TrimPrefix(envName, plainPrefix)

			err := addPlainValues(name, envValue, result)
			if err != nil {
				return nil, err
			}
		default:
			continue
		}
	}
	return result, nil
}

func addJSONValues(name, raw string, result map[string]ParameterValue) error {
	if !nameRegex.MatchString(name) {
		return fmt.Errorf("invalid secure parameter name %q", name)
	}

	errDuplicated := validateNoDuplicatesExist(name, "__MTA_JSON", result)
	if errDuplicated != nil {
		return errDuplicated
	}
	var parsed interface{}

	errUnmarshal := json.Unmarshal([]byte(raw), &parsed)
	if errUnmarshal != nil {
		return fmt.Errorf("invalid JSON for %s: %w", name, errUnmarshal)
	}

	result[name] = ParameterValue{Type: typeJSON, JSONContent: parsed}
	return nil
}

func addCertificateValues(name, raw string, result map[string]ParameterValue) error {
	if !nameRegex.MatchString(name) {
		return fmt.Errorf("invalid secure parameter name %q", name)
	}

	err := validateNoDuplicatesExist(name, "__MTA_CERT", result)
	if err != nil {
		return err
	}

	decoded, err := base64.StdEncoding.DecodeString(raw)
	if err != nil {
		return fmt.Errorf("invalid base64 for %s: %w", name, err)
	}

	result[name] = ParameterValue{Type: typeMultiline, StringContent: string(decoded)}
	return nil
}

func addPlainValues(name, raw string, result map[string]ParameterValue) error {
	if !nameRegex.MatchString(name) {
		return fmt.Errorf("invalid secure parameter name %q", name)
	}

	err := validateNoDuplicatesExist(name, "__MTA", result)
	if err != nil {
		return err
	}

	result[name] = ParameterValue{Type: typeString, StringContent: raw}
	return nil
}
