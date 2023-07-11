package util

import (
	"fmt"
	"io"
	"math"
	"os"
	"path/filepath"
	"strconv"

	"github.com/cloudfoundry/multiapps-cli-plugin/configuration/properties"
	"github.com/pborman/uuid"
)

const MaxFileChunkCount = 50

func generateHash() string {
	return uuid.New()
}

// SplitFile ...
func SplitFile(filePath string, fileChunkSizeInMB uint64) ([]string, error) {
	if fileChunkSizeInMB == 0 {
		return []string{filePath}, nil
	}

	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	fileInfo, _ := file.Stat()

	var fileSize = uint64(fileInfo.Size())
	var fileChunkSize = toBytes(fileChunkSizeInMB)

	// calculate total number of parts the file will be chunked into
	totalPartsNum := uint64(math.Ceil(float64(fileSize) / float64(fileChunkSize)))
	if totalPartsNum <= 1 {
		return []string{filePath}, nil
	}

	partsTempDir := filepath.Join(os.TempDir(), generateHash())
	errCreatingTempDir := os.MkdirAll(partsTempDir, os.ModePerm)
	if errCreatingTempDir != nil {
		return nil, errCreatingTempDir
	}

	baseFileName := filepath.Base(filePath)
	var fileParts []string

	for i := uint64(0); i < totalPartsNum; i++ {
		filePartName := baseFileName + ".part." + strconv.FormatUint(i, 10)
		tempFile := filepath.Join(partsTempDir, filePartName)
		filePart, err := os.Create(tempFile)
		if err != nil {
			return nil, err
		}
		defer filePart.Close()

		partSize := int64(minUint64(fileChunkSize, fileSize-i*fileChunkSize))
		_, err = io.CopyN(filePart, file, partSize)
		if err != nil {
			return nil, err
		}

		fileParts = append(fileParts, filePart.Name())
	}
	return fileParts, nil
}

// ValidateChunkSize validate the chunk size
func ValidateChunkSize(filePath string, fileChunkSizeInMB uint64) error {
	if fileChunkSizeInMB == 0 {
		return nil
	}

	if fileChunkSizeInMB == properties.DefaultUploadChunkSizeInMB {
		return nil
	}

	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return err
	}

	var fileSize = uint64(fileInfo.Size())
	var fileSizeInMb = toMegabytes(fileSize)

	var minFileChunkSizeInMb = uint64(math.Ceil(float64(fileSizeInMb) / float64(MaxFileChunkCount)))

	if fileChunkSizeInMB < minFileChunkSizeInMb {
		return fmt.Errorf("The specified chunk size (%d MB) is below the minimum chunk size (%d MB) for an archive with a size of %d MBs", fileChunkSizeInMB, minFileChunkSizeInMb, fileSizeInMb)
	}
	return nil
}

func minUint64(a, b uint64) uint64 {
	if a < b {
		return a
	}
	return b
}

func toBytes(mb uint64) uint64 {
	return mb * 1024 * 1024
}

func toMegabytes(bytes uint64) uint64 {
	return bytes / 1024 / 1024
}
