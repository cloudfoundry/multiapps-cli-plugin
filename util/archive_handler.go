package util

import (
	"archive/zip"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v2"
)

const defaultDescriptorLocation string = "META-INF/mtad.yaml"

type MtaDescriptor struct {
	SchemaVersion string `yaml:"_schema-version,omitempty"`
	ID            string `yaml:"ID,omitempty"`
	Version       string `yaml:"version,omitempty"`
	Namespace     string `yaml:"namespace,omitempty"`
}

// GetMtaDescriptorFromArchive retrieves MTA ID from MTA archive
func GetMtaDescriptorFromArchive(mtaArchiveFilePath string) (MtaDescriptor, error) {
	var descriptor MtaDescriptor
	// Open the mta archive
	mtaArchiveReader, err := zip.OpenReader(mtaArchiveFilePath)
	if err != nil {
		return MtaDescriptor{}, err
	}
	defer mtaArchiveReader.Close()

	for _, file := range mtaArchiveReader.File {
		// Check for the mta descriptor
		if file.Name == defaultDescriptorLocation {

			descriptorBytes, err := readZipFile(file)
			if err != nil {
				return MtaDescriptor{}, err
			}

			// Unmarshal the content of the temporary deployment descriptor into struct
			err = yaml.Unmarshal(descriptorBytes, &descriptor)
			if err != nil {
				return MtaDescriptor{}, err
			}

			// Return the MTA deployment descriptor, if it is valid
			if descriptor.ID != "" {
				return descriptor, nil
			}
		}
	}

	return MtaDescriptor{}, errors.New("Could not get a valid MTA descriptor from archive")
}

func CreateMtaArchive(source, target string) error {
	zipfile, err := os.Create(target)
	if err != nil {
		return err
	}
	defer zipfile.Close()

	archive := zip.NewWriter(zipfile)
	defer archive.Close()

	info, err := os.Stat(source)
	if err != nil {
		return nil
	}

	var baseDir string
	if info.IsDir() {
		baseDir = "."
	}

	return filepath.Walk(source, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.Name() == filepath.Base(source) {
			return nil
		}
		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}

		if baseDir != "" {
			pathWithoutSourceDirPrefix := strings.TrimPrefix(path, source)
			pathWithoutSourceDirPrefix = strings.TrimPrefix(pathWithoutSourceDirPrefix, string(os.PathSeparator))
			if pathWithoutSourceDirPrefix != "" {
				header.Name = filepath.Join(baseDir, pathWithoutSourceDirPrefix)
			}
		}

		if info.IsDir() {
			header.Name += string(os.PathSeparator)
		}
		header.Method = zip.Deflate

		header.Name = filepath.ToSlash(header.Name)

		writer, err := archive.CreateHeader(header)
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()
		_, err = io.Copy(writer, file)
		return err
	})
}

func readZipFile(file *zip.File) ([]byte, error) {
	reader, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer reader.Close()
	return ioutil.ReadAll(reader)
}

// ParseDeploymentDescriptor parses the deployment descriptor which is located in the provided direcotry
func ParseDeploymentDescriptor(deploymentDescriptorLocation string) (MtaDeploymentDescriptor, string, error) {
	if _, err := os.Stat(deploymentDescriptorLocation); os.IsNotExist(err) {
		return MtaDeploymentDescriptor{}, "", fmt.Errorf("Deployment descriptor location does not exist %s", deploymentDescriptorLocation)
	}
	deploymentDescriptor, err := findDeploymentDescriptor(deploymentDescriptorLocation)
	if err != nil {
		return MtaDeploymentDescriptor{}, "", err
	}
	deploymentDescriptorYaml, err := ioutil.ReadFile(deploymentDescriptor)
	if err != nil {
		return MtaDeploymentDescriptor{}, "", fmt.Errorf("Could not read deployment descriptor: %s", err.Error())
	}
	var descriptor MtaDeploymentDescriptor
	err = yaml.Unmarshal(deploymentDescriptorYaml, &descriptor)
	if err != nil {
		return MtaDeploymentDescriptor{}, "", fmt.Errorf("Could not unmarshal deployment descriptor from yaml: %s", err.Error())
	}

	return descriptor, deploymentDescriptor, nil
}

func findDeploymentDescriptor(deploymentDescriptorLocation string) (string, error) {
	deploymentDescriptorOccurances, err := filepath.Glob(filepath.Join(deploymentDescriptorLocation, deploymentDescriptorYamlName))
	if err != nil {
		return "", fmt.Errorf("Could not find deployment descriptor in location %s: %s", deploymentDescriptorLocation, err.Error())
	}

	if len(deploymentDescriptorOccurances) == 0 {
		return "", fmt.Errorf("No deployment descriptor with name %s was found in location %s", deploymentDescriptorYamlName, deploymentDescriptorLocation)
	}

	return deploymentDescriptorOccurances[0], nil
}
