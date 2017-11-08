package util

import (
	"archive/zip"
	"errors"
	"io/ioutil"
	"fmt"

	"gopkg.in/yaml.v2"
)

const defaultDescriptorLocation string = "META-INF/mtad.yaml"
const defaultDescriptorName string = "mtad.yaml"
var defaultMaxDescriptorSize int64 = 1024 * 1024

type mtaDescriptor struct {
	SchemaVersion string `yaml:"_schema-version,omitempty"`
	ID            string `yaml:"ID,omitempty"`
	Version       string `yaml:"version,omitempty"`
}

// GetMtaIDFromArchive retrieves MTA ID from MTA archive
func GetMtaIDFromArchive(mtaArchveFilePath string) (string, error) {
	// Open the mta archive
	mtaArchiveReader, err := zip.OpenReader(mtaArchveFilePath)
	if err != nil {
		return "", err
	}
	defer mtaArchiveReader.Close()

	for _, file := range mtaArchiveReader.File {
		// Check for the mta descriptor
		if file.Name == defaultDescriptorLocation {

			descriptorBytes, err := readZipFile(file)
			if err != nil {
				return "", err
			}

			// Unmarshal the content of the temporary deployment descriptor into struct
			var descriptor mtaDescriptor
			err = yaml.Unmarshal(descriptorBytes, &descriptor)
			if err != nil {
				return "", err
			}

			// Return the MTA ID extracted from the deployment descriptor, if it is set
			if descriptor.ID != "" {
				return descriptor.ID, nil
			}
		}
	}
	return "", errors.New("Could not get MTA id from archive")
}

func readZipFile(file *zip.File) ([]byte, error) {
	err := validateFileSize(file, defaultMaxDescriptorSize)
	if err != nil {
		return nil, err
	}

	reader, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer reader.Close()
	return ioutil.ReadAll(reader)
}

func validateFileSize(file *zip.File, maxFileSize int64) error{
	fileInfo := file.FileInfo()

	if fileInfo.Size() > maxFileSize {
		return fmt.Errorf("The size %d of file %s exceeds max size limit %d", fileInfo.Size(), file.Name, defaultMaxDescriptorSize)
	}
	return nil
}

func setDefaultMaxDescriptorSize(maxFileSize int64){
	defaultMaxDescriptorSize = maxFileSize
}
