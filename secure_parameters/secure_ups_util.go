package secure_parameters

import (
	"crypto/rand"
	"fmt"
	"strings"

	"code.cloudfoundry.org/cli/v8/plugin"
	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/clients/cfrestclient"
)

func ValidateUpsExistsOrElseCreateIt(userProvidedServiceName string, cliConnection plugin.CliConnection, cfClient cfrestclient.CloudFoundryOperationsExtended) (upsCreatedByTheCli bool, encryptionKeyResult string, err error) {
	doesUpsExist, err := doesUpsExist(userProvidedServiceName, cliConnection, cfClient)
	if err != nil {
		return false, "", fmt.Errorf("Check if the UPS exists: %w", err)
	}

	if doesUpsExist {
		return false, "", nil
	}

	encryptionKey, err := getRandomEncryptionKey()
	if err != nil {
		return false, "", fmt.Errorf("Error while generating AES-256 encryption key: %w", err)
	}

	space, err := cliConnection.GetCurrentSpace()
	if err != nil {
		return false, "", fmt.Errorf("Failed to get the current space: %w", err)
	}

	if space.Guid == "" {
		return false, "", fmt.Errorf("Failed to get the current space Guid")
	}

	upsCredentials := map[string]string{
		"encryptionKey": encryptionKey,
	}

	_, err = cfClient.CreateUserProvidedServiceInstance(userProvidedServiceName, space.Guid, upsCredentials)
	if err != nil {
		return false, "", fmt.Errorf("Failed to create user-provided service %s: %w", userProvidedServiceName, err)
	}

	return true, encryptionKey, nil
}

func CreateDisposableUps(userProvidedServiceName string, cliConnection plugin.CliConnection, cfClient cfrestclient.CloudFoundryOperationsExtended) (upsCreatedByTheCli bool, encryptionKeyResult string, err error) {
	encryptionKey, err := getRandomEncryptionKey()
	if err != nil {
		return false, "", fmt.Errorf("Error while generating AES-256 encryption key: %w", err)
	}

	space, err := cliConnection.GetCurrentSpace()
	if err != nil {
		return false, "", fmt.Errorf("Failed to get the current space: %w", err)
	}

	if space.Guid == "" {
		return false, "", fmt.Errorf("Failed to get the current space Guid")
	}

	upsCredentials := map[string]string{
		"encryptionKey": encryptionKey,
	}

	_, err = cfClient.CreateUserProvidedServiceInstance(userProvidedServiceName, space.Guid, upsCredentials)
	if err != nil {
		return false, "", fmt.Errorf("Failed to create user-provided service %s: %w", userProvidedServiceName, err)
	}

	return true, encryptionKey, nil
}

func getRandomEncryptionKey() (string, error) {
	const alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789-_"

	encryptionKeyBytes := make([]byte, 32)
	if _, err := rand.Read(encryptionKeyBytes); err != nil {
		return "", err
	}

	for i := range encryptionKeyBytes {
		encryptionKeyBytes[i] = alphabet[int(encryptionKeyBytes[i]&63)]
	}

	return string(encryptionKeyBytes), nil
}

func doesUpsExist(userProvidedServiceName string, cliConnection plugin.CliConnection, cfClient cfrestclient.CloudFoundryOperationsExtended) (bool, error) {
	space, errSpace := cliConnection.GetCurrentSpace()
	if errSpace != nil {
		return false, fmt.Errorf("Cannot determine the current space")
	}
	spaceGuid := space.Guid

	_, errServiceInstance := cfClient.GetServiceInstanceByName(userProvidedServiceName, spaceGuid)
	if errServiceInstance != nil {
		if errServiceInstance.Error() == "service instance not found" {
			return false, nil
		}
		return false, fmt.Errorf("Error while checking if the UPS for secure encryption exists: %w", errServiceInstance)
	}

	return true, nil
}

func GetUpsName(mtaId, namespace string) string {
	if strings.TrimSpace(namespace) == "" {
		return "__mta-secure-" + mtaId
	}
	return "__mta-secure-" + mtaId + "-" + namespace
}

func GetRandomisedUpsName(mtaId, namespace string) (disposableUpsName string, err error) {
	randomisedPart, err := getRandomEncryptionKey()
	if err != nil {
		return "", err
	}
	resultSuffix := randomisedPart[:7]

	if strings.TrimSpace(namespace) == "" {
		return "__mta-secure-" + mtaId + "-" + resultSuffix, nil
	}
	return "__mta-secure-" + mtaId + "-" + namespace + "-" + resultSuffix, nil
}
