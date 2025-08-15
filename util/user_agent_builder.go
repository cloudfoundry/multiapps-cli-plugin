package util

import (
	"fmt"
	"os"
	"regexp"
	"runtime"
	"strings"
)

// Limit for MULTIAPPS_USER_AGENT_SUFFIX
const maxSuffixLength = 128

// pluginVersion stores the version set from the main package
var pluginVersion string = "0.0.0"

// SetPluginVersion sets the plugin version for use in User-Agent
func SetPluginVersion(version string) {
	pluginVersion = version
}

// GetPluginVersion returns the current plugin version
func GetPluginVersion() string {
	return pluginVersion
}

// BuildUserAgent creates a User-Agent string in the format:
// "Multiapps-CF-plugin/{version} ({operating system version}) {golang builder version} {custom_env_value}"
func BuildUserAgent() string {
	osInformation := getOperatingSystemInformation()
	goVersion := runtime.Version()
	customValue := getCustomEnvValue()

	userAgent := fmt.Sprintf("Multiapps-CF-plugin/%s (%s) %s", pluginVersion, osInformation, goVersion)

	if customValue != "" {
		userAgent = fmt.Sprintf("%s %s", userAgent, customValue)
	}

	return userAgent
}

// getOperatingSystemInformation returns OS name and architecture
func getOperatingSystemInformation() string {
	return fmt.Sprintf("%s %s", runtime.GOOS, runtime.GOARCH)
}

// getCustomEnvValue reads value from custom environment variable with validation
func getCustomEnvValue() string {
	value := os.Getenv("MULTIAPPS_USER_AGENT_SUFFIX")
	if value == "" {
		return ""
	}

	return sanitizeUserAgentSuffix(value)
}

// sanitizeUserAgentSuffix sanitizes the user agent suffix
func sanitizeUserAgentSuffix(value string) string {
	// Security constraints for HTTP User-Agent header:
	// 1. Max length to prevent server buffer overflow (Tomcat default: 8KB total headers)
	// 2. Only allow safe characters to prevent header injection
	// 3. Remove control characters and newlines

	// Remove control characters, CR, LF, and other dangerous characters
	value = strings.ReplaceAll(value, "\r", "")
	value = strings.ReplaceAll(value, "\n", "")
	value = strings.ReplaceAll(value, "\t", " ")

	// Only allow ASCII characters, spaces, hyphens, dots, underscores, and alphanumeric
	// This prevents header injection attacks
	invalidChars := regexp.MustCompile(`[^a-zA-Z0-9 .\-_]`)
	value = invalidChars.ReplaceAllString(value, "")

	// Remove sequences that could be interpreted as header separators
	value = strings.ReplaceAll(value, ":", "")
	value = strings.ReplaceAll(value, ";", "")

	// Trim whitespace and limit length
	value = strings.TrimSpace(value)
	if len(value) > maxSuffixLength {
		value = value[:maxSuffixLength]
		value = strings.TrimSpace(value) // Trim again in case we cut in the middle of whitespace
	}

	return value
}
