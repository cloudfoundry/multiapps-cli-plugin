package util_test

import (
	"github.com/SAP/cf-mta-plugin/util"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("UriBuilderTest", func() {
	Describe("BuildUriTest", func() {
		const hostName = "test-host"
		const scheme = "test-scheme"
		const path = "test-path"

		Context("no scheme and host are provided", func() {
			It("should return an error", func() {
				uriBuilder := util.NewUriBuilder().SetPath(path)
				_, err := uriBuilder.Build()

				Expect(err).Should(MatchError("The host or scheme could not be empty"))
			})
		})
		Context("when scheme and path are provided", func() {
			It("should return an error", func() {
				uriBuilder := util.NewUriBuilder().SetScheme(scheme).SetPath(path)
				_, err := uriBuilder.Build()

				Expect(err).Should(MatchError("The host or scheme could not be empty"))
			})
		})
		Context("when scheme, host and path are provided", func() {
			It("should the built uri", func() {
				uriBuilder := util.NewUriBuilder().SetHost(hostName).SetScheme(scheme).SetPath(path)
				uri, err := uriBuilder.Build()

				Expect(uri).To(Equal("test-scheme://test-host/test-path"))
				Expect(err).To(BeNil())
			})
		})

	})
})
