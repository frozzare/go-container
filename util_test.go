package container

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetFuncNameParameters(t *testing.T) {
	assert.Equal(t, []string{}, getFuncNameParameters(""))
	assert.Equal(t, []string{"string", "*user.User"}, getFuncNameParameters("func(string, *user.User) string"))
}

func TestGetFieldTypeName(t *testing.T) {
	assert.Equal(t, "", getFieldTypeName(""))
	assert.Equal(t, "*user.User", getFieldTypeName("<*user.User Value>"))
}
