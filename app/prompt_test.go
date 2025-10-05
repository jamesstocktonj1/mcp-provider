package app

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewPromptHandler(t *testing.T) {
	handler, err := NewPromptHandler()
	assert.NoError(t, err)
	assert.NotNil(t, handler)
	assert.NotNil(t, handler.wrpcHandler)
}
