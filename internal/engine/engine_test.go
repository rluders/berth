package engine

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDetectEngine_returnsKnownType(t *testing.T) {
	engine := DetectEngine()
	assert.Contains(t, []EngineType{Docker, Podman, Unknown}, engine,
		"detected engine must be Docker, Podman, or Unknown")
}

func TestGetEnginePath_emptyWhenUnknown(t *testing.T) {
	if DetectEngine() == Unknown {
		assert.Empty(t, GetEnginePath(), "engine path must be empty when engine is Unknown")
	}
}

func TestGetEnginePath_nonEmptyWhenDetected(t *testing.T) {
	if DetectEngine() != Unknown {
		assert.NotEmpty(t, GetEnginePath(), "engine path must be set when engine is detected")
	}
}

func TestGetEnginePath_consistentWithDetect(t *testing.T) {
	engine := DetectEngine()
	path := GetEnginePath()
	if engine == Unknown {
		assert.Empty(t, path)
	} else {
		assert.NotEmpty(t, path)
	}
}

func TestRunEngineCommand_failsWhenUnknown(t *testing.T) {
	if DetectEngine() != Unknown {
		t.Skip("engine is available — skipping unknown-engine error test")
	}
	_, _, err := RunEngineCommand("version")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no container engine detected")
}
