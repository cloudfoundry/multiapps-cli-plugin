package util_test

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/SAP/cf-mta-plugin/testutil"
	"github.com/SAP/cf-mta-plugin/util"

	. "github.com/onsi/ginkgo"
)

var _ = Describe("Digest", func() {

	Describe("ComputeFileChecksum", func() {
		const testFileName = "test-file.txt"

		var testFilePath string
		var testFile *os.File

		BeforeEach(func() {
			testFile, _ = os.Create(testFileName)
			testFilePath, _ = filepath.Abs(testFileName)
		})

		Context("with an unsupported digest alghoritm", func() {
			It("should return an error", func() {
				digest, err := util.ComputeFileChecksum(testFilePath, "unsupported-alghoritm-name")
				testutil.ExpectErrorAndZeroResult(err, digest)
			})
		})

		Context("with a supported digest alghoritm and an empty file", func() {
			It("should return the digest of the file", func() {
				digest, err := util.ComputeFileChecksum(testFilePath, "MD5")
				testutil.ExpectNoErrorAndResult(err, digest, "d41d8cd98f00b204e9800998ecf8427e")
			})
		})

		Context("with a supported digest alghoritm and a non-empty file", func() {
			It("should calculate the digest of the file and exit with zero status", func() {
				const testFileContent = "test file content"
				ioutil.WriteFile(testFile.Name(), []byte(testFileContent), 0644)
				digest, err := util.ComputeFileChecksum(testFilePath, "SHA1")
				testutil.ExpectNoErrorAndResult(err, digest, "9032bbc224ed8b39183cb93b9a7447727ce67f9d")
			})
		})

		AfterEach(func() {
			os.Remove(testFileName)
		})
	})
})
