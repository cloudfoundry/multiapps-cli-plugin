package util_test

import (
	. "github.com/cloudfoundry-incubator/multiapps-cli-plugin/util"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("CfCommandOptionsUtil", func() {
	Describe("DiscardIfEmpty", func() {
		Context("with regular string value", func() {
			It("should return the value as is", func() {
				value := "some value"
				Expect(DiscardIfEmpty(value)).To(Equal(&value))
			})
		})
		Context("with string value with only whitespace", func() {
			It("should return the value as is", func() {
				value := "   	"
				Expect(DiscardIfEmpty(value)).To(Equal(&value))
			})
		})
		Context("with empty string value", func() {
			It("should return nil", func() {
				Expect(TrimAndDiscardIfEmpty("")).To(BeNil())
			})
		})
	})
	Describe("TrimAndDiscardIfEmpty", func() {
		expected_value := "some value"

		Context("with regular string value", func() {
			It("should return the value as is", func() {
				Expect(TrimAndDiscardIfEmpty("some value")).To(Equal(&expected_value))
			})
		})
		Context("with string value with leading and trailing whitespace", func() {
			It("should trim the value", func() {
				Expect(TrimAndDiscardIfEmpty("   some value 	")).To(Equal(&expected_value))
			})
		})
		Context("with empty string value", func() {
			It("should return nil", func() {
				Expect(TrimAndDiscardIfEmpty("")).To(BeNil())
			})
		})
		Context("with only whitespace in value", func() {
			It("should return nil", func() {
				Expect(TrimAndDiscardIfEmpty("   	  	")).To(BeNil())
			})
		})
	})
})
