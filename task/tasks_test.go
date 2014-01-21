package task

import (
	"github.com/bmizerany/assert"
	"testing"
)

func TestCommandRunningString(t *testing.T) {
	cmdOut, err := ExecuteStringTask("echo hello\nworld")
	assert.Equal(t, nil, err)
	assert.Equal(t, "hello\nworld\n", cmdOut)
}

func TestCommandRunning(t *testing.T) {
	cmdOut, err := ExecuteTask("echo", "hello\nworld")
	assert.Equal(t, nil, err)
	assert.Equal(t, "hello\nworld\n", cmdOut)
}
