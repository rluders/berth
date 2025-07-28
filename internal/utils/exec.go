// Package utils provides utility functions for Berth.
package utils

import (
	"os/exec"
)

// RunCommand executes a shell command and returns its stdout and stderr.
// The `name` parameter is the command to run (e.g., "docker", "podman").
// The `args` parameter is a variadic slice of strings representing the arguments to the command.
// It returns the standard output, standard error, and an error if the command fails to execute.
func RunCommand(name string, args ...string) (string, string, error) {
	cmd := exec.Command(name, args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", string(out), err
	}
	return string(out), "", nil
}
