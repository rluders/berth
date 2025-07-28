package utils

import (
	"os/exec"
)

// RunCommand executes a shell command and returns its stdout and stderr.
func RunCommand(name string, args ...string) (string, string, error) {
	cmd := exec.Command(name, args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", string(out), err
	}
	return string(out), "", nil
}
