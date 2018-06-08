package util_test

import (
	. "github.com/cloudfoundry-incubator/multiapps-cli-plugin/util"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("CfCommandBuilder", func() {
	Describe("Build", func() {
		Context("with valid argument and options", func() {
			It("should return a valid command", func() {
				commandBuilder := NewCfCommandStringBuilder()
				commandBuilder.SetName("deploy")
				commandBuilder.AddArgument("jobscheduler.mtar")
				commandBuilder.AddBooleanOption("f")
				commandBuilder.AddOption("t", "100")
				commandBuilder.AddLongBooleanOption("no-start")
				commandBuilder.AddOption("p", "XSA")
				commandBuilder.AddLongOption("schema-version", "3")
				Expect(commandBuilder.Build()).To(Equal("cf deploy jobscheduler.mtar -f -t 100 --no-start -p XSA --schema-version 3"))
			})
		})
		Context("with just options and no arguments", func() {
			It("should return a valid command", func() {
				commandBuilder := NewCfCommandStringBuilder()
				commandBuilder.SetName("deploy")
				commandBuilder.AddOption("i", "12345")
				commandBuilder.AddOption("a", "abort")
				Expect(commandBuilder.Build()).To(Equal("cf deploy -i 12345 -a abort"))
			})
		})
		Context("with just arguments and no options", func() {
			It("should return a valid command", func() {
				commandBuilder := NewCfCommandStringBuilder()
				commandBuilder.SetName("deploy")
				commandBuilder.AddArgument("jobscheduler1.mtar")
				commandBuilder.AddArgument("jobscheduler2.mtar")
				Expect(commandBuilder.Build()).To(Equal("cf deploy jobscheduler1.mtar jobscheduler2.mtar"))
			})
		})
	})
})
