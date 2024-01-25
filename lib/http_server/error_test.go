package http_server

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestError(t *testing.T) {
	err := NewError(500, "some error")
	assert.Equal(t, 500, err.Code)
	assert.Equal(t, "some error", err.Error())
}
