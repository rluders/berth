package service

import (
	"context"
	"fmt"

	volumeTypes "github.com/docker/docker/api/types/volume"
	dockerClient "github.com/docker/docker/client"
)

// VolumeService defines the interface for volume-related operations.
type VolumeService interface {
	VolumeList(ctx context.Context, options volumeTypes.ListOptions) (volumeTypes.ListResponse, error)
	VolumeRemove(ctx context.Context, volumeID string, force bool) error
}

// dockerVolumeService is a concrete implementation of VolumeService.
type dockerVolumeService struct {
	client dockerClient.APIClient
}

// NewVolumeService creates a new VolumeService.
func NewVolumeService(client dockerClient.APIClient) VolumeService {
	return &dockerVolumeService{client: client}
}

// VolumeList lists all volumes.
func (s *dockerVolumeService) VolumeList(ctx context.Context, options volumeTypes.ListOptions) (volumeTypes.ListResponse, error) {
	volumes, err := s.client.VolumeList(ctx, options)
	if err != nil {
		return volumeTypes.ListResponse{}, fmt.Errorf("failed to list volumes: %w", err)
	}
	return volumes, nil
}

// VolumeRemove removes a volume.
func (s *dockerVolumeService) VolumeRemove(ctx context.Context, volumeID string, force bool) error {
	return s.client.VolumeRemove(ctx, volumeID, force)
}
