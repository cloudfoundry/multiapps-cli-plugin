package util_test

import (
	"os"
	"path/filepath"

	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/util"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("ArchiveHandler", func() {
	Describe("GetMtaDescriptorFromArchive", func() {
		var mtaArchiveFilePath, _ = filepath.Abs("../test_resources/commands/mtaArchive.mtar")
		Context("with valid mta archive", func() {
			It("should extract and return the id from deployment descriptor", func() {
				descriptor, err := util.GetMtaDescriptorFromArchive(mtaArchiveFilePath)
				Expect(err).To(BeZero())
				Expect(descriptor.ID).To(Equal("test"))
			})
		})
		Context("with valid mta archive and no deployment descriptor provided", func() {
			It("should return error", func() {
				mtaArchiveFilePath, _ = filepath.Abs("../test_resources/util/mtaArchiveNoDescriptor.mtar")
				_, err := util.GetMtaDescriptorFromArchive(mtaArchiveFilePath)
				Expect(err).To(MatchError("Could not get a valid MTA descriptor from archive"))
			})
		})

		Context("with invalid mta archive", func() {
			const testMtarName = "test.mtar"
			var testFile *os.File
			BeforeEach(func() {
				testFile, _ = os.Create(testMtarName)
				mtaArchiveFilePath, _ = filepath.Abs(testMtarName)
			})
			It("should return error for not a valid zip archive", func() {
				_, err := util.GetMtaDescriptorFromArchive(mtaArchiveFilePath)
				Expect(err).To(MatchError("zip: not a valid zip file"))
			})
			AfterEach(func() {
				testFile.Close()
				os.Remove(testMtarName)
			})
		})
	})
})
