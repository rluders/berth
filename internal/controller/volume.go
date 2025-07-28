package controller

import (
	"fmt"
	"strings"

	"github.com/rluders/container-tui/internal/engine"
)

type Volume struct {
	Name       string
	Driver     string
	Scope      string
	Mountpoint string
}

// ListVolumes lists all volumes.
func ListVolumes() ([]Volume, error) {
	stdout, stderr, err := engine.RunEngineCommand("volume", "ls", "--format", "{{.Name}}\t{{.Driver}}\t{{.Scope}}\t{{.Mountpoint}}")
	if err != nil {
		return nil, fmt.Errorf("failed to list volumes: %s, %w", stderr, err)
	}

	lines := strings.Split(strings.TrimSpace(stdout), "\n")
	var volumes []Volume
	for _, line := range lines {
		if line == "" {
			continue
		}
		fields := strings.Split(line, "\t")
		if len(fields) != 4 {
			// Log or handle malformed line
			continue
		}
		volumes = append(volumes, Volume{
			Name:       fields[0],
			Driver:     fields[1],
			Scope:      fields[2],
			Mountpoint: fields[3],
		})
	}
	return volumes, nil
}

// RemoveVolume removes a volume by its name.
func RemoveVolume(name string) error {
	_, stderr, err := engine.RunEngineCommand("volume", "rm", name)
	if err != nil {
		return fmt.Errorf("failed to remove volume %s: %s, %w", name, stderr, err)
	}
	return nil
}
