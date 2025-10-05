package app

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewServer(t *testing.T) {
	s, err := NewServer()
	assert.NoError(t, err)
	assert.NotNil(t, s)
}
