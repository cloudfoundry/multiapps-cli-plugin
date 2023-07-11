package util_test

import (
	"os"
	"path/filepath"

	"github.com/cloudfoundry/multiapps-cli-plugin/testutil"
	"github.com/cloudfoundry/multiapps-cli-plugin/util"

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

		Context("with an unsupported digest algorithm", func() {
			It("should return an error", func() {
				digest, err := util.ComputeFileChecksum(testFilePath, "unsupported-algorithm-name")
				testutil.ExpectErrorAndZeroResult(err, digest)
			})
		})

		Context("with a supported digest algorithm and an empty file", func() {
			It("should return the digest of the file", func() {
				digest, err := util.ComputeFileChecksum(testFilePath, "MD5")
				testutil.ExpectNoErrorAndResult(err, digest, "d41d8cd98f00b204e9800998ecf8427e")
			})
		})

		Context("with a supported digest algorithm and a non-empty file", func() {
			It("should calculate the digest of the file and exit with zero status", func() {
				const testFileContent = "test file content"
				os.WriteFile(testFile.Name(), []byte(testFileContent), 0644)
				digest, err := util.ComputeFileChecksum(testFilePath, "SHA1")
				testutil.ExpectNoErrorAndResult(err, digest, "9032bbc224ed8b39183cb93b9a7447727ce67f9d")
			})
		})

		AfterEach(func() {
			testFile.Close()
			os.Remove(testFileName)
		})
	})
})
