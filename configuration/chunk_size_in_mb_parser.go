package configuration

import (
	"errors"
	"strconv"
)

type chunkSizeInMBParser struct {
}

func (p chunkSizeInMBParser) Parse(value string) (interface{}, error) {
	chunkSizeInMB, err := strconv.ParseUint(value, 10, 64)
	if err != nil {
		return nil, err
	}
	if chunkSizeInMB == 0 {
		return nil, errors.New("chunk size cannot be 0")
	}
	return chunkSizeInMB, nil
}
