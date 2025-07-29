// Package controller provides the logic for interacting with container engines.
package controller

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/docker/docker/api/types/network"
	"github.com/rluders/berth/internal/engine"
	"github.com/rluders/berth/internal/service"
)

var networkService service.NetworkService

func init() {
	cli, err := engine.NewClient()
	if err != nil {
		panic(fmt.Errorf("failed to create Docker client: %w", err))
	}
	networkService = service.NewNetworkService(cli)
}

// Network represents a network's simplified information.
type Network struct {
	ID     string
	Name   string
	Driver string
	Scope  string
}

// ListNetworks lists all networks.
func ListNetworks() ([]Network, error) {
	networks, err := networkService.NetworkList(context.Background(), network.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list networks: %w", err)
	}

	var result []Network
	for _, n := range networks {
		result = append(result, Network{
			ID:     n.ID,
			Name:   n.Name,
			Driver: n.Driver,
			Scope:  n.Scope,
		})
	}

	return result, nil
}

// InspectNetwork inspects a network and returns its raw JSON output.
func InspectNetwork(idOrName string) (string, error) {
	network, err := networkService.NetworkInspect(context.Background(), idOrName, network.InspectOptions{})
	if err != nil {
		return "", fmt.Errorf("failed to inspect network %s: %w", idOrName, err)
	}

	jsonBytes, err := json.MarshalIndent(network, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal inspect data: %w", err)
	}

	return string(jsonBytes), nil
}
