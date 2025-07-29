// Package engine provides functionality for detecting and interacting with container engines (Docker/Podman).
package engine

import (
	"bytes"
	"fmt"
	"os/exec"
)

// EngineType represents the type of container engine detected.
type EngineType string

const (
	// Docker represents the Docker container engine.
	Docker EngineType = "docker"
	// Podman represents the Podman container engine.
	Podman EngineType = "podman"
	// Unknown represents an unknown or undetected container engine.
	Unknown EngineType = "unknown"
)

var detectedEngine EngineType
var enginePath string

// init function is called automatically when the package is initialized.
// It attempts to detect whether Docker or Podman is available on the system
// and stores the detected engine type and its executable path.
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
// It returns the stdout, stderr, and any error encountered during execution.
func RunEngineCommand(args ...string) (string, string, error) {
	if detectedEngine == Unknown {
		return "", "", fmt.Errorf("no container engine detected")
	}
	cmd := exec.Command(enginePath, args...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	return stdout.String(), stderr.String(), err
}
