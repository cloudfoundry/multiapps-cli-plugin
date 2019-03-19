package configuration

import (
	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/ui"
	"os"
	"strconv"
)

const (
	ChunkSizeInMBEnv     = "CHUNK_SIZE_IN_MB"
	TargetURLEnv         = "DEPLOY_SERVICE_URL"
	DefaultChunkSizeInMB = uint64(45)
)

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

func GetTargetURL() string {
	targetURL := os.Getenv(TargetURLEnv)
	if targetURL != "" {
		ui.Say("Attention: You've specified a custom Deploy Service URL (%s) via the environment variable \"%s\". The application listening on that URL may be outdated, contain bugs or unreleased features or may even be modified by a potentially untrused person. Use at your own risk.\n", targetURL, TargetURLEnv)
	}
	return targetURL
}
