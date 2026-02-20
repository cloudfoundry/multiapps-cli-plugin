package secure_parameters

import (
	"encoding/base64"
	"strings"
	"testing"

	"gopkg.in/yaml.v3"
)

func setEnv(t *testing.T, nameOfEnv, valueOfEnv string) {
	t.Helper()
	t.Setenv(nameOfEnv, valueOfEnv)
}

func TestCollectFromEnv(t *testing.T) {
	setEnv(t, "__MTA___fakePassword", "secretValue")
	setEnv(t, "__MTA_JSON___fakeJson", `{"a":1,"b":"secretValueJson"}`)
	testCertificate := "-----BEGIN CERTIFICATE-----\nMIBgNVBAYTAPXwBc63heW9WrP3qnDEm+UZE4V0Au7OWnOeiobq\n-----END CERTIFICATE-----\n"
	setEnv(t, "__MTA_CERT___fakeCertificate", base64.StdEncoding.EncodeToString([]byte(testCertificate)))
	setEnv(t, "unrelatedEnvFirst", "exampleValueFirst")
	setEnv(t, "unrelatedEnvSecond", "exampleValueSecond")

	resultToTest, err := CollectFromEnv("__MTA")
	if err != nil {
		t.Fatalf("Collecting environment variables has failed: %s", err.Error())
	}

	testNormalVariable(t, &resultToTest)

	testJsonVariable(t, &resultToTest)

	testCertificateVariable(t, &resultToTest, testCertificate)

	if _, exists := resultToTest["other"]; exists {
		t.Fatalf("Unexpected value and environment variable")
	}
}

func testNormalVariable(t *testing.T, resultToTest *map[string]ParameterValue) {
	parameterValue, ok := (*resultToTest)["fakePassword"]
	if !ok {
		t.Fatalf("Missing key 'fakePassword' in map")
	}

	if parameterValue.Type != typeString || parameterValue.StringContent != "secretValue" {
		t.Fatalf("The value of 'fakePassword' key is not correct")
	}
}

func testJsonVariable(t *testing.T, resultToTest *map[string]ParameterValue) {
	jsonValue, ok := (*resultToTest)["fakeJson"]
	if !ok {
		t.Fatalf("Missing key 'fakeJson' in map")
	}

	if jsonValue.Type != typeJSON {
		t.Fatalf("The value of 'fakeJson' key is not correct")
	}

	castedValue, ok := jsonValue.JSONContent.(map[string]interface{})
	if !ok {
		t.Fatal("fakeJson is not an Object")
	}

	if firstJsonValue, ok := castedValue["a"].(float64); !ok || firstJsonValue != 1 {
		t.Fatalf("The first value of the json is not what it should be: %v", castedValue["a"])
	}

	if castedValue["b"] != "secretValueJson" {
		t.Fatalf("The second value of the json is not what it should be: %v", castedValue["b"])
	}
}

func testCertificateVariable(t *testing.T, resultToTest *map[string]ParameterValue, testCertificate string) {
	certificateValue, ok := (*resultToTest)["fakeCertificate"]

	if !ok {
		t.Fatalf("The value of the certificate is not present")
	}

	if certificateValue.Type != typeMultiline || certificateValue.StringContent != testCertificate {
		t.Fatalf("The value of the certificate is not what it should be: %v", certificateValue)
	}
}

func TestCollectFromEnvWhenWrongName(t *testing.T) {
	setEnv(t, "__MTA___fake spaced parameter", "x")

	_, err := CollectFromEnv("__MTA")
	if err == nil || err.Error() != `invalid secure parameter name "fake spaced parameter"` {
		t.Fatalf("Expected invalid name error: %v", err)
	}
}

func TestCollectFromEnvWhenInvalidJson(t *testing.T) {
	setEnv(t, "__MTA_JSON___fakeJson", `{wrongFormat - fake}`)

	_, err := CollectFromEnv("__MTA")
	if err == nil || !strings.Contains(err.Error(), "invalid JSON for fakeJson") {
		t.Fatalf("Expected invalid JSON error: %v", err)
	}
}

func TestCollectFromEnvWhenDuplicateNames(t *testing.T) {
	setEnv(t, "__MTA_JSON___duplicate", `{"fakeValueName":"value"}`)
	setEnv(t, "__MTA___duplicate", "randomValue")

	_, err := CollectFromEnv("__MTA")
	if err == nil || !strings.Contains(err.Error(), `secure parameter "duplicate" defined multiple ways`) {
		t.Fatalf("Expected duplication error: %v", err)
	}
}

func TestCollectFromEnvWhenInvalidCertificate(t *testing.T) {
	setEnv(t, "__MTA_CERT___fakeCertificate", "%**@&@#!#&notBase64*@&$)@!")

	_, err := CollectFromEnv("__MTA")
	if err == nil || !strings.Contains(err.Error(), "invalid base64 for fakeCertificate") {
		t.Fatalf("Expected invalid base64 error: %v", err)
	}
}

func TestCollectFromEnvWhenDifferentPrefix(t *testing.T) {
	setEnv(t, "__MTA_JSON___myJson", `{"apple":"green"}`)

	result, err := CollectFromEnv("__OTHER")

	if err != nil {
		t.Fatalf("Error while trying to collect environment variables with a different prefix: %s", err.Error())
	}

	if len(result) > 0 {
		t.Fatalf("There should be zero environment variables collected, but there are: %d", len(result))
	}
}

func TestBuildSecureExtension(t *testing.T) {
	parameters := map[string]ParameterValue{
		"password":        {Type: typeString, StringContent: "secretValue"},
		"fakeJson":        {Type: typeJSON, JSONContent: map[string]interface{}{"secretParameterFirst": "secretValueOne", "secretParameterSecond": "secretValueTwo"}},
		"fakeCertificate": {Type: typeMultiline, StringContent: "-----BEGIN CERTIFICATE-----\nMIBgNVBAYTAPXwBc63heW9WrP3qnDEm+UZE4V0Au7OWnOeiobq\n-----END CERTIFICATE-----\n"},
	}

	yamlResult, err := BuildSecureExtension(parameters, "test-mta", "")

	if err != nil {
		t.Fatalf("Error while building the secure extension descriptor: %s", err.Error())
	}

	var unmarshaledBack map[string]interface{}

	err2 := yaml.Unmarshal(yamlResult, &unmarshaledBack)

	if err2 != nil {
		t.Fatalf("Error while unmarshaling extension descriptor: %s", err.Error())
	}

	if unmarshaledBack["_schema-version"] != "3.3" {
		t.Fatalf("Schema version is not what it should be: %v", unmarshaledBack["_schema-version"])
	}

	if unmarshaledBack["ID"] != "__mta.secure" {
		t.Fatalf("ID of the secure extension descriptor is not what it should be: %v", unmarshaledBack["ID"])
	}

	if unmarshaledBack["extends"] != "test-mta" {
		t.Fatalf("Extends of secure extension descriptor is not what it should be: %v", unmarshaledBack["extends"])
	}

	parametersUnmarshaled, ok := unmarshaledBack["parameters"].(map[string]interface{})

	if !ok {
		t.Fatalf("Parameters is not a map, but rather: %T", unmarshaledBack["parameters"])
	}

	if parametersUnmarshaled["password"] != "secretValue" {
		t.Fatalf("Value of password is incorrect: %v", parametersUnmarshaled["password"])
	}

	fakeJson, ok := parametersUnmarshaled["fakeJson"].(map[string]interface{})

	if !ok {
		t.Fatalf("fakeJson is not an object but: %T", parametersUnmarshaled["fakeJson"])
	}

	if fakeJson["secretParameterFirst"] != "secretValueOne" {
		t.Fatalf("fakeJson.secretParameterFirst is not what it should be: %v", fakeJson["secretParameterFirst"])
	}

	if fakeJson["secretParameterSecond"] != "secretValueTwo" {
		t.Fatalf("fakeJson.secretParameterSecond is not what it should be: %v", fakeJson["secretParameterSecond"])
	}

	if parametersUnmarshaled["fakeCertificate"] != "-----BEGIN CERTIFICATE-----\nMIBgNVBAYTAPXwBc63heW9WrP3qnDEm+UZE4V0Au7OWnOeiobq\n-----END CERTIFICATE-----\n" {
		t.Fatalf("fakeCertificate is not what it should be: %v", parametersUnmarshaled["fakeCertificate"])
	}
}

func TestBuildSecureExtensionWhenExplicitSchema(t *testing.T) {
	parameters := map[string]ParameterValue{
		"password": {Type: typeString, StringContent: "secretValue"},
	}

	yamlResult, err := BuildSecureExtension(parameters, "test-mta", "3.1")

	if err != nil {
		t.Fatalf("Error while building the secure extension descriotor: %s", err.Error())
	}

	var unmarshaledBack map[string]interface{}

	err2 := yaml.Unmarshal(yamlResult, &unmarshaledBack)

	if err2 != nil {
		t.Fatalf("Error while unmarshaling extension descriptor: %s", err.Error())
	}

	if unmarshaledBack["_schema-version"] != "3.1" {
		t.Fatalf("Schema version  must be 3.1, but it is: %v", unmarshaledBack["_schema-version"])
	}
}

func TestBuildSecureExtensionWhenNoParameters(t *testing.T) {
	_, err := BuildSecureExtension(map[string]ParameterValue{}, "test-mta", "")

	if err == nil || err.Error() != "no secure parameters collected" {
		t.Fatalf("Expected no parameters error, but rather got: %v", err)
	}
}

func TestBuildSecureExtensionWhenNoMtaId(t *testing.T) {
	parameters := map[string]ParameterValue{
		"password": {Type: typeString, StringContent: "secretValue"},
	}

	_, err := BuildSecureExtension(parameters, "", "")
	if err == nil || err.Error() != "mtaID is required for the extension descriptor's field 'extends'" {
		t.Fatalf("Expected missing mta id error, but rather got: %v", err)
	}
}
