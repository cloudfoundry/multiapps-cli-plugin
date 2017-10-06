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
	var hash hash.Hash
	switch strings.ToUpper(algorithm) {
	case "MD5":
		hash = md5.New()
	case "SHA1":
		hash = sha1.New()
	case "SHA256":
		hash = sha256.New()
	case "SHA512":
		hash = sha512.New()
	default:
		return "", fmt.Errorf("Unsupported digest algorithm %q", algorithm)
	}

	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	_, err = io.Copy(hash, file)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(hash.Sum(nil)), nil
}
