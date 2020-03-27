package util_test

import (
	"os"
	"path/filepath"

	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/util"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("ArchiveHandler", func() {
	Describe("GetMtaIDFromArchive", func() {
		var mtaArchiveFilePath, _ = filepath.Abs("../test_resources/commands/mtaArchive.mtar")
		Context("with valid mta archive", func() {
			It("should extract and return the id from deployment descriptor", func() {
				Expect(util.GetMtaIDFromArchive(mtaArchiveFilePath)).To(Equal("test"))
			})
		})
		Context("with valid mta archive and no deployment descriptor provided", func() {
			It("should return error", func() {
				mtaArchiveFilePath, _ = filepath.Abs("../test_resources/util/mtaArchiveNoDescriptor.mtar")
				_, err := util.GetMtaIDFromArchive(mtaArchiveFilePath)
				Expect(err).To(MatchError("Could not get MTA ID from archive"))
			})
		})

		Context("with invalid mta archive", func() {
			BeforeEach(func() {
				os.Create("test.mtar")
				mtaArchiveFilePath, _ = filepath.Abs("test.mtar")
			})
			It("should return error for not a valid zip archive", func() {
				_, err := util.GetMtaIDFromArchive(mtaArchiveFilePath)
				Expect(err).To(MatchError("zip: not a valid zip file"))
			})
			AfterEach(func() {
				os.Remove(mtaArchiveFilePath)
			})
		})
	})
})
