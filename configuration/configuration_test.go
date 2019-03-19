package configuration_test

import (
	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/configuration"
	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/ui"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"os"
	"strconv"
)

var _ = Describe("Configuration", func() {

	BeforeEach(func() {
		ui.DisableTerminalOutput(true)
	})

	Describe("GetTargetURL", func() {

		BeforeEach(func() {
			os.Unsetenv(configuration.TargetURLEnv)
		})

		Context("with a set environment variable", func() {
			It("should return its value", func() {
				targetURL := "http://my-multiapps-controller.domain.com"
				os.Setenv(configuration.TargetURLEnv, targetURL)
				Expect(configuration.GetTargetURL()).To(Equal(targetURL))
			})
		})
		Context("without a set environment variable", func() {
			It("should return an empty string", func() {
				Expect(configuration.GetTargetURL()).To(BeEmpty())
			})
		})

	})

	Describe("GetChunkSizeInMB", func() {

		BeforeEach(func() {
			os.Unsetenv(configuration.TargetURLEnv)
		})

		Context("with a set environment variable", func() {
			Context("containing a positive integer", func() {
				It("should return its value", func() {
					chunkSizeInMB := uint64(5)
					os.Setenv(configuration.ChunkSizeInMBEnv, strconv.Itoa(int(chunkSizeInMB)))
					Expect(configuration.GetChunkSizeInMB()).To(Equal(chunkSizeInMB))
				})
			})
			Context("containing zero", func() {
				It("should return the default value", func() {
					chunkSizeInMB := 0
					os.Setenv(configuration.ChunkSizeInMBEnv, strconv.Itoa(chunkSizeInMB))
					Expect(configuration.GetChunkSizeInMB()).To(Equal(configuration.DefaultChunkSizeInMB))
				})
			})
			Context("containing a string", func() {
				It("should return the default value", func() {
					chunkSizeInMB := "abc"
					os.Setenv(configuration.ChunkSizeInMBEnv, chunkSizeInMB)
					Expect(configuration.GetChunkSizeInMB()).To(Equal(configuration.DefaultChunkSizeInMB))
				})
			})
		})
		Context("without a set environment variable", func() {
			It("should return the default value", func() {
				Expect(configuration.GetChunkSizeInMB()).To(Equal(configuration.DefaultChunkSizeInMB))
			})
		})

	})

})
