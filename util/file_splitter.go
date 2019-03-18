package util

import (
	"io"
	"math"
	"os"
	"path/filepath"
	"strconv"

	"github.com/pborman/uuid"
)

func generateHash() string {
	return uuid.New()
}

// SplitFile ...
func SplitFile(filePath string, fileChunkSizeInMb uint64) ([]string, error) {
	if fileChunkSizeInMb == 0 {
		return []string{filePath}, nil
	}

	file, err := os.Open(filePath)

	if err != nil {
		return nil, err
	}

	defer file.Close()

	fileInfo, _ := file.Stat()

	var fileSize = uint64(fileInfo.Size())
	var fileChunkSize = toBytes(fileChunkSizeInMb)

	// calculate total number of parts the file will be chunked into
	totalPartsNum := uint64(math.Ceil(float64(fileSize) / float64(fileChunkSize)))

	baseFileName := filepath.Base(filePath)
	var fileParts []string
	if totalPartsNum <= 1 {
		return []string{filePath}, nil
	}
	partsTempDir := filepath.Join(os.TempDir(), generateHash())
	errCreatingTempDir := os.MkdirAll(partsTempDir, os.ModePerm)
	if errCreatingTempDir != nil {
		return nil, errCreatingTempDir
	}
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

func minUint64(a, b uint64) uint64 {
	if a < b {
		return a
	}
	return b
}

func toBytes(mb uint64) uint64 {
	return uint64(mb) * 1024 * 1024
}
