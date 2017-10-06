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
func SplitFile(filePath string) ([]string, error) {
	file, err := os.Open(filePath)

	if err != nil {
		return nil, err
	}

	defer file.Close()

	fileInfo, _ := file.Stat()

	var fileSize = fileInfo.Size()

	const fileChunk = 45 * (1 << 20) // 45 MB

	// calculate total number of parts the file will be chunked into

	totalPartsNum := uint64(math.Ceil(float64(fileSize) / float64(fileChunk)))

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

		partSize := int64(math.Min(fileChunk, float64(fileSize-int64(i*fileChunk))))
		_, err = io.CopyN(filePart, file, partSize)
		if err != nil {
			return nil, err
		}

		fileParts = append(fileParts, filePart.Name())
	}
	return fileParts, nil
}
