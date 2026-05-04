package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRunCommand_success(t *testing.T) {
	stdout, stderr, err := RunCommand("echo", "hello")
	require.NoError(t, err)
	assert.Contains(t, stdout, "hello")
	assert.Empty(t, stderr)
}

func TestRunCommand_invalidBinary(t *testing.T) {
	_, _, err := RunCommand("berth-nonexistent-binary-xyz")
	assert.Error(t, err)
}

func TestRunCommand_nonZeroExit(t *testing.T) {
	_, stderr, err := RunCommand("false")
	assert.Error(t, err)
	_ = stderr
}
