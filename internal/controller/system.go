// Package controller provides the logic for interacting with container engines.
package controller

import (
	"context"
	"fmt"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/rluders/berth/internal/engine"
	"github.com/rluders/berth/internal/service"
)

var systemService service.SystemService

func init() {
	cli, err := engine.NewClient()
	if err != nil {
		panic(fmt.Errorf("failed to create Docker client: %w", err))
	}
	systemService = service.NewSystemService(cli)
}

// SystemInfo holds system-wide statistics about containers, images, volumes, and networks.
type SystemInfo struct {
	Containers int // Total number of containers.
	Running    int // Number of running containers.
	Paused     int // Number of paused containers.
	Stopped    int // Number of stopped containers.
	Images     int // Total number of images.
	Volumes    int // Total number of volumes.
	Networks   int // Total number of networks.
	DiskUsage  string
}

// GetSystemInfo retrieves system-wide information about containers, images, and volumes.
func GetSystemInfo() (SystemInfo, error) {
	info, err := systemService.Info(context.Background())
	if err != nil {
		return SystemInfo{}, fmt.Errorf("failed to get info: %w", err)
	}

	diskUsage, err := systemService.DiskUsage(context.Background(), types.DiskUsageOptions{})
	if err != nil {
		return SystemInfo{}, fmt.Errorf("failed to get disk usage: %w", err)
	}

	return SystemInfo{
		Containers: info.Containers,
		Running:    info.ContainersRunning,
		Paused:     info.ContainersPaused,
		Stopped:    info.ContainersStopped,
		Images:     info.Images,
		Volumes:    len(diskUsage.Volumes),
		Networks:   len(info.DriverStatus),
		DiskUsage:  fmt.Sprintf("%d", diskUsage.LayersSize),
	}, nil
}

// BasicCleanup removes stopped containers, unused networks, and unused images.
func BasicCleanup() (string, error) {
	var output strings.Builder

	_, err := systemService.ContainersPrune(context.Background(), filters.Args{})
	if err != nil {
		output.WriteString(fmt.Sprintf("Failed to prune containers: %s\n", err))
	} else {
		output.WriteString("Containers pruned successfully\n")
	}

	_, err = systemService.NetworksPrune(context.Background(), filters.Args{})
	if err != nil {
		output.WriteString(fmt.Sprintf("Failed to prune networks: %s\n", err))
	} else {
		output.WriteString("Networks pruned successfully\n")
	}

	_, err = systemService.ImagesPrune(context.Background(), filters.Args{})
	if err != nil {
		output.WriteString(fmt.Sprintf("Failed to prune images: %s\n", err))
	} else {
		output.WriteString("Images pruned successfully\n")
	}

	return output.String(), nil
}

// AdvancedCleanup removes dangling volumes and dangling images.
func AdvancedCleanup() (string, error) {
	var output strings.Builder

	_, err := systemService.VolumesPrune(context.Background(), filters.Args{})
	if err != nil {
		output.WriteString(fmt.Sprintf("Failed to prune volumes: %s\n", err))
	} else {
		output.WriteString("Volumes pruned successfully\n")
	}

	args := filters.NewArgs()
	args.Add("dangling", "true")
	_, err = systemService.ImagesPrune(context.Background(), args)
	if err != nil {
		output.WriteString(fmt.Sprintf("Failed to prune dangling images: %s\n", err))
	} else {
		output.WriteString("Dangling images pruned successfully\n")
	}

	return output.String(), nil
}

// TotalCleanup prunes all unused containers, images, volumes, and networks.
func TotalCleanup() (string, error) {
	_, err := systemService.ContainersPrune(context.Background(), filters.Args{})
	if err != nil {
		return "", fmt.Errorf("failed to prune containers: %w", err)
	}

	_, err = systemService.NetworksPrune(context.Background(), filters.Args{})
	if err != nil {
		return "", fmt.Errorf("failed to prune networks: %w", err)
	}

	_, err = systemService.ImagesPrune(context.Background(), filters.Args{})
	if err != nil {
		return "", fmt.Errorf("failed to prune images: %w", err)
	}

	_, err = systemService.VolumesPrune(context.Background(), filters.Args{})
	if err != nil {
		return "", fmt.Errorf("failed to prune volumes: %w", err)
	}

	return "Total cleanup performed successfully", nil
}
