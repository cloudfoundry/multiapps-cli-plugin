package util_test

import (
	"os"
	"path/filepath"

	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/testutil"
	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/util"

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

		Context("with SHA-384 algorithm name (as returned by deploy service for FIPS+NCS)", func() {
			It("should compute the SHA-384 digest of the file", func() {
				const testFileContent = "test file content"
				os.WriteFile(testFile.Name(), []byte(testFileContent), 0644)
				digest, err := util.ComputeFileChecksum(testFilePath, "SHA-384")
				testutil.ExpectNoErrorAndResult(err, digest, "4b87814537ab46771af4f37f259d6321f8c36b9e0d9ed1eda10005031cf7de752b9937a9cfe02744a3fc9785c3106317")
			})
		})

		Context("with legacy hyphenated SHA-256 algorithm name", func() {
			It("should compute the SHA-256 digest of the file", func() {
				const testFileContent = "test file content"
				os.WriteFile(testFile.Name(), []byte(testFileContent), 0644)
				digest, err := util.ComputeFileChecksum(testFilePath, "SHA-256")
				testutil.ExpectNoErrorAndResult(err, digest, "60f5237ed4049f0382661ef009d2bc42e48c3ceb3edb6600f7024e7ab3b838f3")
			})
		})

		Context("with legacy hyphenated SHA-1 algorithm name", func() {
			It("should compute the SHA-1 digest of the file", func() {
				const testFileContent = "test file content"
				os.WriteFile(testFile.Name(), []byte(testFileContent), 0644)
				digest, err := util.ComputeFileChecksum(testFilePath, "SHA-1")
				testutil.ExpectNoErrorAndResult(err, digest, "9032bbc224ed8b39183cb93b9a7447727ce67f9d")
			})
		})

		Context("with legacy hyphenated SHA-512 algorithm name", func() {
			It("should compute the SHA-512 digest of the file", func() {
				const testFileContent = "test file content"
				os.WriteFile(testFile.Name(), []byte(testFileContent), 0644)
				digest, err := util.ComputeFileChecksum(testFilePath, "SHA-512")
				testutil.ExpectNoErrorAndResult(err, digest, "6543084c6981d2d3e44d295d20adee4e5fa3fc1bb2b7d3ac80ca456b0141dfd75657ba40ea5ae398824b3a871c8b7762c8bdbcbc252fa7f358d33d90ee455d86")
			})
		})

		Context("with SHA3-256 algorithm", func() {
			It("should compute the SHA3-256 digest of the file", func() {
				const testFileContent = "test file content"
				os.WriteFile(testFile.Name(), []byte(testFileContent), 0644)
				digest, err := util.ComputeFileChecksum(testFilePath, "SHA3-256")
				testutil.ExpectNoErrorAndResult(err, digest, "0de0da51c6d2c645c7822739d886b15894c0baa51e668415b683b3d028a12758")
			})
		})

		Context("with SHA3-512 algorithm", func() {
			It("should compute the SHA3-512 digest of the file", func() {
				const testFileContent = "test file content"
				os.WriteFile(testFile.Name(), []byte(testFileContent), 0644)
				digest, err := util.ComputeFileChecksum(testFilePath, "SHA3-512")
				testutil.ExpectNoErrorAndResult(err, digest, "fa9fc2b12fcb7ca5e1758dfb907d139d9f803676a8659de0b0527e77adf4b332631a22306bad174e46ce47a4e5b6e4fd2ccdc58cefa51b233c0e5ce6651330f8")
			})
		})

		AfterEach(func() {
			testFile.Close()
			os.Remove(testFileName)
		})
	})
})
