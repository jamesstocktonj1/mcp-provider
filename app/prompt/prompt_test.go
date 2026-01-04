package prompt

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewPromptHandler(t *testing.T) {
	handler, err := NewPromptHandler(nil, nil, nil)
	assert.NoError(t, err)
	assert.NotNil(t, handler)
	assert.NotNil(t, handler.wrpcHandler)
}

func TestHandlePutTargetLink(t *testing.T) {

}

func TestHandleDelTargetLink(t *testing.T) {

}
