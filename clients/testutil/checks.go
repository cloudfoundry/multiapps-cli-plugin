package testutil

import (
	"testing"

	"github.com/go-openapi/runtime"
	"github.com/stretchr/testify/assert"
)

func CheckSuccess(t *testing.T, expected interface{}, actual interface{}) {
	assert.IsType(t, expected, actual)
	assert.EqualValues(t, expected, actual)
}

func CheckError(t *testing.T, err error, opName string, status int, res interface{}) {
	assert.IsType(t, (*runtime.APIError)(nil), err)
	assert.Equal(t, opName, err.(*runtime.APIError).OperationName)
	assert.Equal(t, status, err.(*runtime.APIError).Code)
	assert.Nil(t, res)
}
