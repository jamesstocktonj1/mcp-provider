package app

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewServer(t *testing.T) {
	saveStdin := os.Stdin

	var err error
	os.Stdin, err = os.Open("../testdata/wasmcloud-provider.bin")
	assert.NoError(t, err)

	defer func() {
		os.Stdin = saveStdin
	}()

	s, err := NewServer()
	assert.NoError(t, err)
	assert.NotNil(t, s)

	time.Sleep(50 * time.Millisecond)
}
