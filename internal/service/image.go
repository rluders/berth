package service

import (
	"context"
	"fmt"

	imageTypes "github.com/docker/docker/api/types/image"
	dockerClient "github.com/docker/docker/client"
)

// ImageService defines the interface for image-related operations.
type ImageService interface {
	ImageList(ctx context.Context, options imageTypes.ListOptions) ([]imageTypes.Summary, error)
	ImageRemove(ctx context.Context, imageID string, options imageTypes.RemoveOptions) ([]imageTypes.DeleteResponse, error)
}

// dockerImageService is a concrete implementation of ImageService.
type dockerImageService struct {
	client dockerClient.APIClient
}

// NewImageService creates a new ImageService.
func NewImageService(client dockerClient.APIClient) ImageService {
	return &dockerImageService{client: client}
}

// ImageList lists all images.
func (s *dockerImageService) ImageList(ctx context.Context, options imageTypes.ListOptions) ([]imageTypes.Summary, error) {
	images, err := s.client.ImageList(ctx, options)
	if err != nil {
		return nil, fmt.Errorf("failed to list images: %w", err)
	}
	return images, nil
}

// ImageRemove removes an image.
func (s *dockerImageService) ImageRemove(ctx context.Context, imageID string, options imageTypes.RemoveOptions) ([]imageTypes.DeleteResponse, error) {
	resp, err := s.client.ImageRemove(ctx, imageID, options)
	if err != nil {
		return nil, fmt.Errorf("failed to remove image: %w", err)
	}
	return resp, nil
}
