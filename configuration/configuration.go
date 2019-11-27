package configuration

import (
	"os"
	"strconv"

	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/ui"
)

const (
	// ChunkSizeInMBEnv Defines the chunk size of MTAR in MB
	ChunkSizeInMBEnv = "CHUNK_SIZE_IN_MB"
	// TargetURLEnv Defines the URL of the deploy service
	TargetURLEnv = "DEPLOY_SERVICE_URL"
	// I accept the risk of conflicts env variable
	IAcceptTheRiskOfEnv = "I_ACCEPT_THE_RISK_OF"
	// I want to stream the MTA env variable
	IWantToStreamTheMTAFromEnv = "I_WANT_TO_STREAM_THE_MTA_FROM"
	// DefaultChunkSizeInMB ...
	DefaultChunkSizeInMB = uint64(45)
)

// GetChunkSizeInMB Retrieves the MTAR chunk size from environment or uses the default one
func GetChunkSizeInMB() uint64 {
	chunkSizeInMb, isSet := os.LookupEnv(ChunkSizeInMBEnv)
	if isSet {
		parsedChunkSizeInMb, err := strconv.ParseUint(chunkSizeInMb, 10, 64)
		if err == nil && parsedChunkSizeInMb != 0 {
			ui.Say("Attention: You've specified a custom chunk size (%d MB) via the environment variable \"%s\".", parsedChunkSizeInMb, ChunkSizeInMBEnv)
			return parsedChunkSizeInMb
		}
		ui.Warn("Attention: You've specified an INVALID custom chunk size (%s) via the environment variable \"%s\". Using default: %d", chunkSizeInMb, ChunkSizeInMBEnv, DefaultChunkSizeInMB)
	}
	return DefaultChunkSizeInMB
}

// GetTargetURL Retrieves the URL of the deploy service if set in the environment
func GetTargetURL() string {
	targetURL := os.Getenv(TargetURLEnv)
	if targetURL != "" {
		ui.Say("Attention: You've specified a custom Deploy Service URL (%s) via the environment variable \"%s\". The application listening on that URL may be outdated, contain bugs or unreleased features or may even be modified by a potentially untrused person. Use at your own risk.\n", targetURL, TargetURLEnv)
	}
	return targetURL
}

// IsStreamingFlagSet gives answer wheter the user wants to stream the mta or not
func IsStreamingFlagSet() bool {
	iAcceptTheRiskOfConflict := os.Getenv(IAcceptTheRiskOfEnv)
	iWantToStreamTheMTA := os.Getenv(IWantToStreamTheMTAFromEnv)
	if iAcceptTheRiskOfConflict == "CONFLICT" && iWantToStreamTheMTA == "FILE" {
		ui.Say("Attention: You've requested streaming the mta to the deploy-service. Conflict checks and partial deploy are disabled")
		return true
	}
	return false
}
