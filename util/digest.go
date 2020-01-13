package util

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"hash"
	"io"
	"os"
	"strings"
)

// ComputeFileChecksum computes the checksum of the specified file based on the specified algorithm
func ComputeFileChecksum(filePath, algorithm string) (string, error) {
	algImpl, err := getAlgorithmImpl(algorithm)
	if err != nil {
		return "", err
	}

	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	_, err = io.Copy(algImpl, file)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(algImpl.Sum(nil)), nil
}

func getAlgorithmImpl(algorithm string) (hash.Hash, error) {
	switch strings.ToUpper(algorithm) {
	case "MD5":
		return md5.New(), nil
	case "SHA1":
		return sha1.New(), nil
	case "SHA256":
		return sha256.New(), nil
	case "SHA512":
		return sha512.New(), nil
	default:
		return nil, fmt.Errorf("Unsupported digest algorithm %q", algorithm)
	}
}
