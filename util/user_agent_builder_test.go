package util_test

import (
	"fmt"
	"os"
	"regexp"
	"runtime"
	"strings"

	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/util"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("UserAgentBuilder", func() {

	Describe("BuildUserAgent", func() {
		var originalVersion string
		var originalCfCliVersion string

		BeforeEach(func() {
			// Save original versions for restoration
			originalVersion = util.GetPluginVersion()
			originalCfCliVersion = util.GetCfCliVersion()
		})

		AfterEach(func() {
			// Restore original versions and clean up environment
			util.SetPluginVersion(originalVersion)
			util.SetCfCliVersion(originalCfCliVersion)
			os.Unsetenv("MULTIAPPS_USER_AGENT_SUFFIX")
		})

		Context("with default version", func() {
			BeforeEach(func() {
				util.SetPluginVersion("0.0.0")
			})

			It("should contain default version, OS, arch, and Go version", func() {
				userAgent := util.BuildUserAgent()

				Expect(userAgent).To(ContainSubstring("Multiapps-CF-plugin/0.0.0"))
				Expect(userAgent).To(ContainSubstring(runtime.GOOS))
				Expect(userAgent).To(ContainSubstring(runtime.GOARCH))
				Expect(userAgent).To(ContainSubstring(runtime.Version()))
				Expect(userAgent).To(HavePrefix("Multiapps-CF-plugin/"))
			})
		})

		Context("with custom version", func() {
			BeforeEach(func() {
				util.SetPluginVersion("3.6.0")
			})

			It("should contain the custom version", func() {
				userAgent := util.BuildUserAgent()

				Expect(userAgent).To(ContainSubstring("Multiapps-CF-plugin/3.6.0"))
				Expect(userAgent).To(ContainSubstring(runtime.GOOS))
				Expect(userAgent).To(ContainSubstring(runtime.GOARCH))
				Expect(userAgent).To(ContainSubstring(runtime.Version()))
			})
		})

		Context("with dev version", func() {
			BeforeEach(func() {
				util.SetPluginVersion("3.7.0-dev")
			})

			It("should contain the dev version", func() {
				userAgent := util.BuildUserAgent()

				Expect(userAgent).To(ContainSubstring("Multiapps-CF-plugin/3.7.0-dev"))
				Expect(userAgent).To(ContainSubstring(runtime.GOOS))
				Expect(userAgent).To(ContainSubstring(runtime.GOARCH))
				Expect(userAgent).To(ContainSubstring(runtime.Version()))
			})
		})

		Context("with custom CF CLI version", func() {
			BeforeEach(func() {
				util.SetPluginVersion("1.0.0")
				util.SetCfCliVersion("8.5.0")
			})

			It("should contain the CF CLI version in parentheses", func() {
				userAgent := util.BuildUserAgent()

				Expect(userAgent).To(ContainSubstring("(8.5.0)"))
				Expect(userAgent).To(ContainSubstring("Multiapps-CF-plugin/1.0.0"))
			})

			It("should have correct format with CF CLI version", func() {
				userAgent := util.BuildUserAgent()

				// Expected format: "Multiapps-CF-plugin/{version} ({os} {arch}) {go version} ({cf cli version})"
				expectedPattern := fmt.Sprintf("Multiapps-CF-plugin/1.0.0 \\(%s %s\\) %s \\(8.5.0\\)", runtime.GOOS, runtime.GOARCH, regexp.QuoteMeta(runtime.Version()))
				matched, _ := regexp.MatchString(expectedPattern, userAgent)
				Expect(matched).To(BeTrue(), fmt.Sprintf("User agent '%s' should match pattern '%s'", userAgent, expectedPattern))
			})
		})

		Context("with custom environment value", func() {
			BeforeEach(func() {
				util.SetPluginVersion("3.6.0")
				os.Setenv("MULTIAPPS_USER_AGENT_SUFFIX", "custom-suffix")
			})

			It("should include the custom suffix", func() {
				userAgent := util.BuildUserAgent()

				Expect(userAgent).To(ContainSubstring("Multiapps-CF-plugin/3.6.0"))
				Expect(userAgent).To(ContainSubstring(runtime.GOOS))
				Expect(userAgent).To(ContainSubstring(runtime.GOARCH))
				Expect(userAgent).To(ContainSubstring(runtime.Version()))
				Expect(userAgent).To(ContainSubstring("custom-suffix"))
			})
		})

		Context("with sanitized environment value", func() {
			BeforeEach(func() {
				util.SetPluginVersion("1.0.0")
			})

			It("should remove dangerous characters", func() {
				os.Setenv("MULTIAPPS_USER_AGENT_SUFFIX", "value:with;dangerous")
				userAgent := util.BuildUserAgent()

				Expect(userAgent).To(ContainSubstring("valuewithdangerous"))
				Expect(userAgent).ToNot(ContainSubstring(":"))
				Expect(userAgent).ToNot(ContainSubstring(";"))
			})

			It("should handle control characters", func() {
				os.Setenv("MULTIAPPS_USER_AGENT_SUFFIX", "value\r\nwith\tcontrol")
				userAgent := util.BuildUserAgent()

				Expect(userAgent).To(ContainSubstring("valuewith control"))
				Expect(userAgent).ToNot(ContainSubstring("\r"))
				Expect(userAgent).ToNot(ContainSubstring("\n"))
				Expect(userAgent).ToNot(ContainSubstring("\t"))
			})

			It("should remove unicode characters", func() {
				os.Setenv("MULTIAPPS_USER_AGENT_SUFFIX", "value†with•unicode")
				userAgent := util.BuildUserAgent()

				Expect(userAgent).To(ContainSubstring("valuewithunicode"))
				Expect(userAgent).ToNot(ContainSubstring("†"))
				Expect(userAgent).ToNot(ContainSubstring("•"))
			})

			It("should trim whitespace", func() {
				os.Setenv("MULTIAPPS_USER_AGENT_SUFFIX", "  trimmed  ")
				userAgent := util.BuildUserAgent()

				Expect(userAgent).To(ContainSubstring("trimmed"))
				Expect(userAgent).ToNot(ContainSubstring("  trimmed  "))
			})

			It("should handle empty value", func() {
				os.Setenv("MULTIAPPS_USER_AGENT_SUFFIX", "")
				userAgent := util.BuildUserAgent()

				// Should only contain the base user agent without suffix
				expectedBase := fmt.Sprintf("Multiapps-CF-plugin/1.0.0 (%s %s) %s (unknown-cf cli version)", runtime.GOOS, runtime.GOARCH, runtime.Version())
				Expect(userAgent).To(Equal(expectedBase))
			})

			It("should handle only whitespace", func() {
				os.Setenv("MULTIAPPS_USER_AGENT_SUFFIX", "   \t\r\n   ")
				userAgent := util.BuildUserAgent()

				// Should only contain the base user agent without suffix (whitespace gets trimmed to empty)
				expectedBase := fmt.Sprintf("Multiapps-CF-plugin/1.0.0 (%s %s) %s (unknown-cf cli version)", runtime.GOOS, runtime.GOARCH, runtime.Version())
				Expect(userAgent).To(Equal(expectedBase))
			})

			It("should truncate excessively long values", func() {
				longValue := strings.Repeat("a", 600) // 600 chars, should be truncated to 128
				os.Setenv("MULTIAPPS_USER_AGENT_SUFFIX", longValue)
				userAgent := util.BuildUserAgent()

				// Should contain truncated suffix
				expectedSuffix := strings.Repeat("a", 128)
				Expect(userAgent).To(ContainSubstring(expectedSuffix))
				// But not the full original value
				Expect(userAgent).ToNot(ContainSubstring(longValue))
			})
		})

		Context("with maximum length suffix", func() {
			BeforeEach(func() {
				util.SetPluginVersion("1.0.0")
				longSuffix := strings.Repeat("a", 128)
				os.Setenv("MULTIAPPS_USER_AGENT_SUFFIX", longSuffix)
			})

			It("should not exceed reasonable total length", func() {
				userAgent := util.BuildUserAgent()

				const maxTotalLength = 1024
				Expect(len(userAgent)).To(BeNumerically("<=", maxTotalLength))
			})

			It("should not contain dangerous characters", func() {
				userAgent := util.BuildUserAgent()

				Expect(userAgent).ToNot(ContainSubstring("\r"))
				Expect(userAgent).ToNot(ContainSubstring("\n"))
			})

			It("should not contain dangerous colon characters except in Go version", func() {
				userAgent := util.BuildUserAgent()

				if strings.Contains(userAgent, ":") {
					Expect(userAgent).To(ContainSubstring("go1."))
				}
			})
		})
	})

	Describe("SetPluginVersion", func() {
		var originalVersion string

		BeforeEach(func() {
			originalVersion = util.GetPluginVersion()
		})

		AfterEach(func() {
			util.SetPluginVersion(originalVersion)
		})

		It("should set the plugin version correctly", func() {
			testVersion := "1.2.3-test"
			util.SetPluginVersion(testVersion)

			Expect(util.GetPluginVersion()).To(Equal(testVersion))
		})

		It("should be used in BuildUserAgent", func() {
			testVersion := "1.2.3-test"
			util.SetPluginVersion(testVersion)

			userAgent := util.BuildUserAgent()
			expectedPrefix := fmt.Sprintf("Multiapps-CF-plugin/%s", testVersion)
			Expect(userAgent).To(HavePrefix(expectedPrefix))
		})
	})

	Describe("GetPluginVersion", func() {
		var originalVersion string

		BeforeEach(func() {
			originalVersion = util.GetPluginVersion()
		})

		AfterEach(func() {
			util.SetPluginVersion(originalVersion)
		})

		It("should return the current version", func() {
			testVersion := "2.4.6"
			util.SetPluginVersion(testVersion)

			Expect(util.GetPluginVersion()).To(Equal(testVersion))
		})

		It("should return default version initially", func() {
			util.SetPluginVersion("0.0.0")

			Expect(util.GetPluginVersion()).To(Equal("0.0.0"))
		})
	})

	Describe("SetCfCliVersion", func() {
		var originalVersion string

		BeforeEach(func() {
			originalVersion = util.GetCfCliVersion()
		})

		AfterEach(func() {
			util.SetCfCliVersion(originalVersion)
		})

		It("should set the CF CLI version correctly", func() {
			testVersion := "8.5.0"
			util.SetCfCliVersion(testVersion)

			Expect(util.GetCfCliVersion()).To(Equal(testVersion))
		})

		It("should be used in BuildUserAgent", func() {
			testVersion := "8.5.0"
			util.SetCfCliVersion(testVersion)
			util.SetPluginVersion("1.0.0")

			userAgent := util.BuildUserAgent()
			Expect(userAgent).To(ContainSubstring(fmt.Sprintf("(%s)", testVersion)))
		})
	})

	Describe("GetCfCliVersion", func() {
		var originalVersion string

		BeforeEach(func() {
			originalVersion = util.GetCfCliVersion()
		})

		AfterEach(func() {
			util.SetCfCliVersion(originalVersion)
		})

		It("should return the current CF CLI version", func() {
			testVersion := "8.6.0"
			util.SetCfCliVersion(testVersion)

			Expect(util.GetCfCliVersion()).To(Equal(testVersion))
		})

		It("should return default CF CLI version initially", func() {
			util.SetCfCliVersion(util.DefaultCliVersion)

			Expect(util.GetCfCliVersion()).To(Equal(util.DefaultCliVersion))
		})
	})

	Describe("Build-time user agent suffix", func() {
		var originalSuffix string

		BeforeEach(func() {
			originalSuffix = util.GetUserAgentSuffixOption()
			util.SetPluginVersion("1.0.0")
		})

		AfterEach(func() {
			util.SetUserAgentSuffixOption(originalSuffix)
			os.Unsetenv("MULTIAPPS_USER_AGENT_SUFFIX")
		})

		It("should use build-time suffix when environment variable is not set", func() {
			buildTimeSuffix := "build-time-suffix"
			util.SetUserAgentSuffixOption(buildTimeSuffix)
			os.Unsetenv("MULTIAPPS_USER_AGENT_SUFFIX")

			userAgent := util.BuildUserAgent()
			Expect(userAgent).To(ContainSubstring(buildTimeSuffix))
		})

		It("should prioritize environment variable over build-time suffix", func() {
			buildTimeSuffix := "build-time"
			envSuffix := "env-override"

			util.SetUserAgentSuffixOption(buildTimeSuffix)
			os.Setenv("MULTIAPPS_USER_AGENT_SUFFIX", envSuffix)

			userAgent := util.BuildUserAgent()
			Expect(userAgent).To(ContainSubstring(envSuffix))
			Expect(userAgent).ToNot(ContainSubstring(buildTimeSuffix))
		})

		It("should handle empty build-time suffix", func() {
			util.SetUserAgentSuffixOption("")
			os.Unsetenv("MULTIAPPS_USER_AGENT_SUFFIX")

			userAgent := util.BuildUserAgent()
			expectedBase := fmt.Sprintf("Multiapps-CF-plugin/1.0.0 (%s %s) %s (%s)", runtime.GOOS, runtime.GOARCH, runtime.Version(), util.GetCfCliVersion())
			Expect(userAgent).To(Equal(expectedBase))
		})
	})
})
