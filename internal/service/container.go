package service

import (
	"context"
	"fmt"
	"io"

	containerTypes "github.com/docker/docker/api/types/container"
	dockerClient "github.com/docker/docker/client"
)

// ContainerService defines the interface for container-related operations.
type ContainerService interface {
	ListContainers(ctx context.Context, options containerTypes.ListOptions) ([]containerTypes.Summary, error)
	StartContainer(ctx context.Context, containerID string, options containerTypes.StartOptions) error
	StopContainer(ctx context.Context, containerID string, options containerTypes.StopOptions) error
	RemoveContainer(ctx context.Context, containerID string, options containerTypes.RemoveOptions) error
	ContainerLogs(ctx context.Context, containerID string, options containerTypes.LogsOptions) (io.ReadCloser, error)
	ContainerInspect(ctx context.Context, containerID string) (containerTypes.InspectResponse, error)
}

// dockerContainerService is a concrete implementation of ContainerService.
type dockerContainerService struct {
	client dockerClient.APIClient
}

// NewContainerService creates a new ContainerService.
func NewContainerService(client dockerClient.APIClient) ContainerService {
	return &dockerContainerService{client: client}
}

// ListContainers lists all containers.
func (s *dockerContainerService) ListContainers(ctx context.Context, options containerTypes.ListOptions) ([]containerTypes.Summary, error) {
	containers, err := s.client.ContainerList(ctx, options)
	if err != nil {
		return nil, fmt.Errorf("failed to list containers: %w", err)
	}
	return containers, nil
}

// StartContainer starts a container.
func (s *dockerContainerService) StartContainer(ctx context.Context, containerID string, options containerTypes.StartOptions) error {
	return s.client.ContainerStart(ctx, containerID, options)
}

// StopContainer stops a container.
func (s *dockerContainerService) StopContainer(ctx context.Context, containerID string, options containerTypes.StopOptions) error {
	return s.client.ContainerStop(ctx, containerID, options)
}

// RemoveContainer removes a container.
func (s *dockerContainerService) RemoveContainer(ctx context.Context, containerID string, options containerTypes.RemoveOptions) error {
	return s.client.ContainerRemove(ctx, containerID, options)
}

// ContainerLogs retrieves container logs.
func (s *dockerContainerService) ContainerLogs(ctx context.Context, containerID string, options containerTypes.LogsOptions) (io.ReadCloser, error) {
	return s.client.ContainerLogs(ctx, containerID, options)
}

// ContainerInspect inspects a container.
func (s *dockerContainerService) ContainerInspect(ctx context.Context, containerID string) (containerTypes.InspectResponse, error) {
	inspect, err := s.client.ContainerInspect(ctx, containerID)
	if err != nil {
		return containerTypes.InspectResponse{}, fmt.Errorf("failed to inspect container %s: %w", containerID, err)
	}
	return inspect, nil
}
