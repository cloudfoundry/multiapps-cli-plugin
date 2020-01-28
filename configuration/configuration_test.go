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

	Describe("GetBackendURL", func() {

		BeforeEach(func() {
			os.Unsetenv(configuration.BackendURLConfigurableProperty.Name)
			for _, name := range configuration.BackendURLConfigurableProperty.DeprecatedNames {
				os.Unsetenv(name)
			}
		})

		Context("with a set environment variable", func() {
			It("should return its value", func() {
				backendURL := "http://my-multiapps-controller.domain.com"
				os.Setenv(configuration.BackendURLConfigurableProperty.Name, backendURL)
				Expect(configuration.GetBackendURL()).To(Equal(backendURL))
			})
		})
		Context("with a set environment variable (deprecated)", func() {
			It("should return its value", func() {
				if len(configuration.BackendURLConfigurableProperty.DeprecatedNames) > 0 {
					backendURL := "http://my-multiapps-controller.domain.com"
					os.Setenv(configuration.BackendURLConfigurableProperty.DeprecatedNames[0], backendURL)
					Expect(configuration.GetBackendURL()).To(Equal(backendURL))
				}
			})
		})
		Context("without a set environment variable", func() {
			It("should return an empty string", func() {
				Expect(configuration.GetBackendURL()).To(BeEmpty())
			})
		})

	})

	Describe("GetChunkSizeInMB", func() {

		BeforeEach(func() {
			os.Unsetenv(configuration.ChunkSizeInMBConfigurableProperty.Name)
			for _, name := range configuration.ChunkSizeInMBConfigurableProperty.DeprecatedNames {
				os.Unsetenv(name)
			}
		})

		Context("with a set environment variable", func() {
			Context("containing a positive integer", func() {
				It("should return its value", func() {
					chunkSizeInMB := uint64(5)
					os.Setenv(configuration.ChunkSizeInMBConfigurableProperty.Name, strconv.Itoa(int(chunkSizeInMB)))
					Expect(configuration.GetChunkSizeInMB()).To(Equal(chunkSizeInMB))
				})
			})
			Context("containing zero", func() {
				It("should return the default value", func() {
					chunkSizeInMB := 0
					os.Setenv(configuration.ChunkSizeInMBConfigurableProperty.Name, strconv.Itoa(chunkSizeInMB))
					Expect(configuration.GetChunkSizeInMB()).To(Equal(configuration.DefaultChunkSizeInMB))
				})
			})
			Context("containing a string", func() {
				It("should return the default value", func() {
					chunkSizeInMB := "abc"
					os.Setenv(configuration.ChunkSizeInMBConfigurableProperty.Name, chunkSizeInMB)
					Expect(configuration.GetChunkSizeInMB()).To(Equal(configuration.DefaultChunkSizeInMB))
				})
			})
		})
		Context("with a set environment variable (deprecated)", func() {
			It("should return its value", func() {
				if len(configuration.ChunkSizeInMBConfigurableProperty.DeprecatedNames) > 0 {
					chunkSizeInMB := uint64(5)
					os.Setenv(configuration.ChunkSizeInMBConfigurableProperty.DeprecatedNames[0], strconv.Itoa(int(chunkSizeInMB)))
					Expect(configuration.GetChunkSizeInMB()).To(Equal(chunkSizeInMB))
				}
			})
		})
		Context("without a set environment variable", func() {
			It("should return the default value", func() {
				Expect(configuration.GetChunkSizeInMB()).To(Equal(configuration.DefaultChunkSizeInMB))
			})
		})

	})

})
