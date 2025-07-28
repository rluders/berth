package controller

import (
	"fmt"
	"strings"

	"github.com/rluders/container-tui/internal/engine"
)

type Network struct {
	ID      string
	Name    string
	Driver  string
	Scope   string
}

// ListNetworks lists all networks.
func ListNetworks() ([]Network, error) {
	stdout, stderr, err := engine.RunEngineCommand("network", "ls", "--format", "{{.ID}}\t{{.Name}}\t{{.Driver}}\t{{.Scope}}")
	if err != nil {
		return nil, fmt.Errorf("failed to list networks: %s, %w", stderr, err)
	}

	lines := strings.Split(strings.TrimSpace(stdout), "\n")
	var networks []Network
	for _, line := range lines {
		if line == "" {
			continue
		}
		fields := strings.Split(line, "\t")
		if len(fields) != 4 {
			// Log or handle malformed line
			continue
		}
		networks = append(networks, Network{
			ID:     fields[0],
			Name:   fields[1],
			Driver: fields[2],
			Scope:  fields[3],
		})
	}
	return networks, nil
}

// InspectNetwork inspects a network and returns its raw JSON output.
func InspectNetwork(idOrName string) (string, error) {
	stdout, stderr, err := engine.RunEngineCommand("network", "inspect", idOrName)
	if err != nil {
		return "", fmt.Errorf("failed to inspect network %s: %s, %w", idOrName, stderr, err)
	}
	return stdout, nil
}