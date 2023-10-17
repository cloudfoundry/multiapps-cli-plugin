package properties

import "strconv"

var UploadChunksInParallel = ConfigurableProperty{
	Name:                  "MULTIAPPS_UPLOAD_CHUNKS_IN_PARALLEL",
	Parser:                booleanParser{},
	ParsingSuccessMessage: "Attention: You've specified %v for the environment variable %s.\n",
	ParsingFailureMessage: "Invalid boolean value (%s) for environment variable %s. Using default value %v.\n",
	DefaultValue:          false,
}

type booleanParser struct{}

func (booleanParser) Parse(value string) (interface{}, error) {
	result, err := strconv.ParseBool(value)
	if err != nil {
		return false, err
	}
	return result, nil
}
