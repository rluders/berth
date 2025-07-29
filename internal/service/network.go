package service

import (
	"context"
	"fmt"

	networkTypes "github.com/docker/docker/api/types/network"
	dockerClient "github.com/docker/docker/client"
)

// NetworkService defines the interface for network-related operations.
type NetworkService interface {
	NetworkList(ctx context.Context, options networkTypes.ListOptions) ([]networkTypes.Summary, error)
	NetworkInspect(ctx context.Context, networkID string, options networkTypes.InspectOptions) (networkTypes.Inspect, error)
}

// dockerNetworkService is a concrete implementation of NetworkService.
type dockerNetworkService struct {
	client dockerClient.APIClient
}

// NewNetworkService creates a new NetworkService.
func NewNetworkService(client dockerClient.APIClient) NetworkService {
	return &dockerNetworkService{client: client}
}

// NetworkList lists all networks.
func (s *dockerNetworkService) NetworkList(ctx context.Context, options networkTypes.ListOptions) ([]networkTypes.Summary, error) {
	networks, err := s.client.NetworkList(ctx, options)
	if err != nil {
		return nil, fmt.Errorf("failed to list networks: %w", err)
	}
	return networks, nil
}

// NetworkInspect inspects a network.
func (s *dockerNetworkService) NetworkInspect(ctx context.Context, networkID string, options networkTypes.InspectOptions) (networkTypes.Inspect, error) {
	network, err := s.client.NetworkInspect(ctx, networkID, options)
	if err != nil {
		return networkTypes.Inspect{}, fmt.Errorf("failed to inspect network %s: %w", networkID, err)
	}
	return network, nil
}
