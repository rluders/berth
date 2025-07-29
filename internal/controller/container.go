// Package controller provides the logic for interacting with container engines.
package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/docker/docker/api/types/container"
	"github.com/rluders/berth/internal/engine"
	"github.com/rluders/berth/internal/service"
)

var containerService service.ContainerService

func init() {
	cli, err := engine.NewClient()
	if err != nil {
		// Handle error, perhaps log it or panic if it's unrecoverable
		panic(fmt.Errorf("failed to create Docker client: %w", err))
	}
	containerService = service.NewContainerService(cli)
}

// Container represents a container's simplified information.
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
	containers, err := containerService.ListContainers(context.Background(), container.ListOptions{All: true})
	if err != nil {
		return nil, fmt.Errorf("failed to list containers: %w", err)
	}

	var result []Container
	for _, c := range containers {
		result = append(result, Container{
			ID:      c.ID[:12],
			Image:   c.Image,
			Command: c.Command,
			Created: fmt.Sprintf("%d", c.Created),
			Status:  c.Status,
			Ports:   fmt.Sprintf("%v", c.Ports),
			Names:   strings.Join(c.Names, ","),
		})
	}

	return result, nil
}

// StartContainer starts a container by its ID or name.
func StartContainer(idOrName string) error {
	return containerService.StartContainer(context.Background(), idOrName, container.StartOptions{})
}

// StopContainer stops a container by its ID or name.
func StopContainer(idOrName string) error {
	return containerService.StopContainer(context.Background(), idOrName, container.StopOptions{})
}

// RemoveContainer removes a container by its ID or name.
func RemoveContainer(idOrName string) error {
	return containerService.RemoveContainer(context.Background(), idOrName, container.RemoveOptions{})
}

// GetContainerLogs retrieves the logs of a container.
func GetContainerLogs(idOrName string) (string, error) {
	out, err := containerService.ContainerLogs(context.Background(), idOrName, container.LogsOptions{ShowStdout: true, ShowStderr: true})
	if err != nil {
		return "", fmt.Errorf("failed to get logs for container %s: %w", idOrName, err)
	}
	defer out.Close()

	buf := new(strings.Builder)
	_, err = io.Copy(buf, out)
	if err != nil {
		return "", fmt.Errorf("failed to read logs for container %s: %w", idOrName, err)
	}

	return buf.String(), nil
}

// InspectContainer inspects a container by its ID or name and returns its raw JSON output.
func InspectContainer(idOrName string) (string, error) {
	inspect, err := containerService.ContainerInspect(context.Background(), idOrName)
	if err != nil {
		return "", fmt.Errorf("failed to inspect container %s: %w", idOrName, err)
	}

	jsonBytes, err := json.MarshalIndent(inspect, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal inspect data: %w", err)
	}

	return string(jsonBytes), nil
}
