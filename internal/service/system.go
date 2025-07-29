package service

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/api/types/system"
	"github.com/docker/docker/api/types/volume"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	dockerClient "github.com/docker/docker/client"
)

// SystemService defines the interface for system-related operations.
type SystemService interface {
	Info(ctx context.Context) (system.Info, error)
	DiskUsage(ctx context.Context, options types.DiskUsageOptions) (types.DiskUsage, error)
	ContainersPrune(ctx context.Context, pruneFilters filters.Args) (container.PruneReport, error)
	NetworksPrune(ctx context.Context, pruneFilters filters.Args) (network.PruneReport, error)
	ImagesPrune(ctx context.Context, pruneFilters filters.Args) (image.PruneReport, error)
	VolumesPrune(ctx context.Context, pruneFilters filters.Args) (volume.PruneReport, error)
}

// dockerSystemService is a concrete implementation of SystemService.
type dockerSystemService struct {
	client dockerClient.APIClient
}

// NewSystemService creates a new SystemService.
func NewSystemService(client dockerClient.APIClient) SystemService {
	return &dockerSystemService{client: client}
}

// Info returns information about the Docker system.
func (s *dockerSystemService) Info(ctx context.Context) (system.Info, error) {
	info, err := s.client.Info(ctx)
	if err != nil {
		return system.Info{}, fmt.Errorf("failed to get info: %w", err)
	}
	return info, nil
}

// DiskUsage returns disk usage statistics.
func (s *dockerSystemService) DiskUsage(ctx context.Context, options types.DiskUsageOptions) (types.DiskUsage, error) {
	diskUsage, err := s.client.DiskUsage(ctx, options)
	if err != nil {
		return types.DiskUsage{}, fmt.Errorf("failed to get disk usage: %w", err)
	}
	return diskUsage, nil
}

// ContainersPrune prunes unused containers.
func (s *dockerSystemService) ContainersPrune(ctx context.Context, pruneFilters filters.Args) (container.PruneReport, error) {
	report, err := s.client.ContainersPrune(ctx, pruneFilters)
	if err != nil {
		return container.PruneReport{}, fmt.Errorf("failed to prune containers: %w", err)
	}
	return report, nil
}

// NetworksPrune prunes unused networks.
func (s *dockerSystemService) NetworksPrune(ctx context.Context, pruneFilters filters.Args) (network.PruneReport, error) {
	report, err := s.client.NetworksPrune(ctx, pruneFilters)
	if err != nil {
		return network.PruneReport{}, fmt.Errorf("failed to prune networks: %w", err)
	}
	return report, nil
}

// ImagesPrune prunes unused images.
func (s *dockerSystemService) ImagesPrune(ctx context.Context, pruneFilters filters.Args) (image.PruneReport, error) {
	report, err := s.client.ImagesPrune(ctx, pruneFilters)
	if err != nil {
		return image.PruneReport{}, fmt.Errorf("failed to prune images: %w", err)
	}
	return report, nil
}

// VolumesPrune prunes unused volumes.
func (s *dockerSystemService) VolumesPrune(ctx context.Context, pruneFilters filters.Args) (volume.PruneReport, error) {
	report, err := s.client.VolumesPrune(ctx, pruneFilters)
	if err != nil {
		return volume.PruneReport{}, fmt.Errorf("failed to prune volumes: %w", err)
	}
	return report, nil
}
