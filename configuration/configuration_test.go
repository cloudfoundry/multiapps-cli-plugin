package configuration_test

import (
	"github.com/cloudfoundry/multiapps-cli-plugin/configuration"
	"github.com/cloudfoundry/multiapps-cli-plugin/configuration/properties"
	"github.com/cloudfoundry/multiapps-cli-plugin/ui"
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
			os.Unsetenv(properties.BackendURL.Name)
		})

		Context("with a set environment variable", func() {
			It("should return its value", func() {
				backendURL := "http://my-multiapps-controller.domain.com"
				os.Setenv(properties.BackendURL.Name, backendURL)
				configurationSnapshot := configuration.NewSnapshot()
				Expect(configurationSnapshot.GetBackendURL()).To(Equal(backendURL))
			})
		})
		Context("without a set environment variable", func() {
			It("should return an empty string", func() {
				configurationSnapshot := configuration.NewSnapshot()
				Expect(configurationSnapshot.GetBackendURL()).To(BeEmpty())
			})
		})

	})

	Describe("GetChunkSizeInMB", func() {

		BeforeEach(func() {
			os.Unsetenv(properties.UploadChunkSizeInMB.Name)
		})

		Context("with a set environment variable", func() {
			Context("containing a positive integer", func() {
				It("should return its value", func() {
					uploadChunkSizeInMB := uint64(5)
					os.Setenv(properties.UploadChunkSizeInMB.Name, strconv.Itoa(int(uploadChunkSizeInMB)))
					configurationSnapshot := configuration.NewSnapshot()
					Expect(configurationSnapshot.GetUploadChunkSizeInMB()).To(Equal(uploadChunkSizeInMB))
				})
			})
			Context("containing zero", func() {
				It("should return the default value", func() {
					uploadChunkSizeInMB := 0
					os.Setenv(properties.UploadChunkSizeInMB.Name, strconv.Itoa(uploadChunkSizeInMB))
					configurationSnapshot := configuration.NewSnapshot()
					Expect(configurationSnapshot.GetUploadChunkSizeInMB()).To(Equal(properties.DefaultUploadChunkSizeInMB))
				})
			})
			Context("containing a string", func() {
				It("should return the default value", func() {
					uploadChunkSizeInMB := "abc"
					os.Setenv(properties.UploadChunkSizeInMB.Name, uploadChunkSizeInMB)
					configurationSnapshot := configuration.NewSnapshot()
					Expect(configurationSnapshot.GetUploadChunkSizeInMB()).To(Equal(properties.DefaultUploadChunkSizeInMB))
				})
			})
		})
		Context("without a set environment variable", func() {
			It("should return the default value", func() {
				configurationSnapshot := configuration.NewSnapshot()
				Expect(configurationSnapshot.GetUploadChunkSizeInMB()).To(Equal(properties.DefaultUploadChunkSizeInMB))
			})
		})

	})

})
