package controller

import (
	"fmt"
	"strings"

	"github.com/rluders/berth/internal/engine"
)

type Container struct {
	ID      string
	Image   string
	Command string
	Created string
	Status  string
	Ports   string
	Names   string
}

// ListContainers lists all running and stopped containers.
func ListContainers() ([]Container, error) {
	stdout, stderr, err := engine.RunEngineCommand("ps", "-a", "--format", "{{.ID}}\t{{.Image}}\t{{.Command}}\t{{.CreatedAt}}\t{{.Status}}\t{{.Ports}}\t{{.Names}}")
	if err != nil {
		return nil, fmt.Errorf("failed to list containers: %s, %w", stderr, err)
	}

	lines := strings.Split(strings.TrimSpace(stdout), "\n")
	var containers []Container
	for _, line := range lines {
		if line == "" {
			continue
		}
		fields := strings.Split(line, "\t")
		if len(fields) != 7 {
			// Log or handle malformed line
			continue
		}
		containers = append(containers, Container{
			ID:      fields[0],
			Image:   fields[1],
			Command: fields[2],
			Created: fields[3],
			Status:  fields[4],
			Ports:   fields[5],
			Names:   fields[6],
		})
	}
	return containers, nil
}

// StartContainer starts a container by its ID or name.
func StartContainer(idOrName string) error {
	_, stderr, err := engine.RunEngineCommand("start", idOrName)
	if err != nil {
		return fmt.Errorf("failed to start container %s: %s, %w", idOrName, stderr, err)
	}
	return nil
}

// StopContainer stops a container by its ID or name.
func StopContainer(idOrName string) error {
	_, stderr, err := engine.RunEngineCommand("stop", idOrName)
	if err != nil {
		return fmt.Errorf("failed to stop container %s: %s, %w", idOrName, stderr, err)
	}
	return nil
}

// RemoveContainer removes a container by its ID or name.
func RemoveContainer(idOrName string) error {
	_, stderr, err := engine.RunEngineCommand("rm", idOrName)
	if err != nil {
		return fmt.Errorf("failed to remove container %s: %s, %w", idOrName, stderr, err)
	}
	return nil
}

// GetContainerLogs retrieves the logs of a container.
func GetContainerLogs(idOrName string) (string, error) {
	stdout, stderr, err := engine.RunEngineCommand("logs", idOrName)
	if err != nil {
		return "", fmt.Errorf("failed to get logs for container %s: %s, %w", idOrName, stderr, err)
	}
	return stdout, nil
}

// InspectContainer inspects a container by its ID or name and returns its raw JSON output.
func InspectContainer(idOrName string) (string, error) {
	stdout, stderr, err := engine.RunEngineCommand("inspect", idOrName)
	if err != nil {
		return "", fmt.Errorf("failed to inspect container %s: %s, %w", idOrName, stderr, err)
	}
	return stdout, nil
}
