package engine

import (
	"fmt"
	"os/exec"

	"github.com/rluders/container-tui/internal/utils"
)

type EngineType string

const (
	Docker  EngineType = "docker"
	Podman  EngineType = "podman"
	Unknown EngineType = "unknown"
)

var detectedEngine EngineType
var enginePath string

func init() {
	if path, err := exec.LookPath(string(Docker)); err == nil {
		detectedEngine = Docker
		enginePath = path
		return
	}
	if path, err := exec.LookPath(string(Podman)); err == nil {
		detectedEngine = Podman
		enginePath = path
		return
	}
	detectedEngine = Unknown
	enginePath = ""
}

// DetectEngine returns the detected engine type.
func DetectEngine() EngineType {
	return detectedEngine
}

// GetEnginePath returns the absolute path to the detected engine binary.
func GetEnginePath() string {
	return enginePath
}

// RunEngineCommand executes a command using the detected container engine.
func RunEngineCommand(args ...string) (string, string, error) {
	if detectedEngine == Unknown {
		return "", "", fmt.Errorf("no container engine detected")
	}
	return utils.RunCommand(enginePath, args...)
}
