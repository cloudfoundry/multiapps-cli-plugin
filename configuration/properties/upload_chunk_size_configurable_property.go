package properties

import (
	"errors"
	"strconv"
)

const DefaultUploadChunkSizeInMB = uint64(45)

var UploadChunkSizeInMB = ConfigurableProperty{
	Name: "MULTIAPPS_UPLOAD_CHUNK_SIZE",
	DeprecatedNames: []string{
		"CHUNK_SIZE_IN_MB",
	},
	Parser:                uploadChunkSizeParser{},
	ParsingSuccessMessage: "Attention: You've specified a custom chunk size (%d MB) via the environment variable \"%s\".\n",
	ParsingFailureMessage: "Attention: You've specified an INVALID custom chunk size (%s) via the environment variable \"%s\". Using default: %d\n",
	DefaultValue:          DefaultUploadChunkSizeInMB,
}

type uploadChunkSizeParser struct{}

func (p uploadChunkSizeParser) Parse(value string) (interface{}, error) {
	parsedValue, err := parseUint64(value)
	if err != nil {
		return nil, err
	}
	if parsedValue == 0 {
		return nil, errors.New("chunk size cannot be 0")
	}
	return parsedValue, nil
}

func parseUint64(value string) (uint64, error) {
	return strconv.ParseUint(value, 10, 64)
}
