package controller

import (
	"fmt"
	"strings"

	"github.com/rluders/container-tui/internal/engine"
)

type SystemInfo struct {
	Containers int
	Running    int
	Paused     int
	Stopped    int
	Images     int
	Volumes    int
	Networks   int
	DiskUsage  string
}

// GetSystemInfo retrieves system-wide information about containers, images, and volumes.
func GetSystemInfo() (SystemInfo, error) {
	var info SystemInfo

	// Get container stats
	stdout, stderr, err := engine.RunEngineCommand("info", "--format", "{{.Containers}}\t{{.ContainersRunning}}\t{{.ContainersPaused}}\t{{.ContainersStopped}}")
	if err != nil {
		return info, fmt.Errorf("failed to get container info: %s, %w", stderr, err)
	}
	fields := strings.Split(strings.TrimSpace(stdout), "\n")
	if len(fields) == 4 {
		fmt.Sscanf(fields[0], "%d", &info.Containers)
		fmt.Sscanf(fields[1], "%d", &info.Running)
		fmt.Sscanf(fields[2], "%d", &info.Paused)
		fmt.Sscanf(fields[3], "%d", &info.Stopped)
	}

	// Get image count
	stdout, stderr, err = engine.RunEngineCommand("images", "-q")
	if err != nil {
		return info, fmt.Errorf("failed to get image count: %s, %w", stderr, err)
	}
	info.Images = len(strings.Split(strings.TrimSpace(stdout), "\n"))

	// Get volume count
	stdout, stderr, err = engine.RunEngineCommand("volume", "ls", "-q")
	if err != nil {
		return info, fmt.Errorf("failed to get volume count: %s, %w", stderr, err)
	}
	info.Volumes = len(strings.Split(strings.TrimSpace(stdout), "\n"))

	// Get network count
	stdout, stderr, err = engine.RunEngineCommand("network", "ls", "-q")
	if err != nil {
		return info, fmt.Errorf("failed to get network count: %s, %w", stderr, err)
	}
	info.Networks = len(strings.Split(strings.TrimSpace(stdout), "\n"))

	// Get disk usage (simplified for now, can be improved)
	stdout, stderr, err = engine.RunEngineCommand("system", "df", "--format", "{{.Size}}")
	if err != nil {
		// This command might not be available in older versions or for Podman in the same way
		// Handle gracefully or provide a fallback
		info.DiskUsage = "N/A"
	} else {
		lines := strings.Split(strings.TrimSpace(stdout), "\n")
		if len(lines) > 0 {
			info.DiskUsage = lines[len(lines)-1] // Last line usually contains total size
		} else {
			info.DiskUsage = "N/A"
		}
	}

	return info, nil
}

// BasicCleanup removes stopped containers, unused networks, and unused images.
func BasicCleanup() (string, error) {
	var output strings.Builder

	// Prune containers
	stdout, stderr, err := engine.RunEngineCommand("container", "prune", "-f")
	if err != nil {
		output.WriteString(fmt.Sprintf("Failed to prune containers: %s\n", stderr))
	} else {
		output.WriteString(fmt.Sprintf("Container prune output:\n%s\n", stdout))
	}

	// Prune networks
	stdout, stderr, err = engine.RunEngineCommand("network", "prune", "-f")
	if err != nil {
		output.WriteString(fmt.Sprintf("Failed to prune networks: %s\n", stderr))
	} else {
		output.WriteString(fmt.Sprintf("Network prune output:\n%s\n", stdout))
	}

	// Prune images
	stdout, stderr, err = engine.RunEngineCommand("image", "prune", "-f")
	if err != nil {
		output.WriteString(fmt.Sprintf("Failed to prune images: %s\n", stderr))
	} else {
		output.WriteString(fmt.Sprintf("Image prune output:\n%s\n", stdout))
	}

	return output.String(), nil
}

// AdvancedCleanup removes dangling volumes and dangling images.
func AdvancedCleanup() (string, error) {
	var output strings.Builder

	// Prune volumes
	stdout, stderr, err := engine.RunEngineCommand("volume", "prune", "-f")
	if err != nil {
		output.WriteString(fmt.Sprintf("Failed to prune volumes: %s\n", stderr))
	} else {
		output.WriteString(fmt.Sprintf("Volume prune output:\n%s\n", stdout))
	}

	// Prune dangling images
	stdout, stderr, err = engine.RunEngineCommand("image", "prune", "-f", "--filter", "dangling=true")
	if err != nil {
		output.WriteString(fmt.Sprintf("Failed to prune dangling images: %s\n", stderr))
	} else {
		output.WriteString(fmt.Sprintf("Dangling image prune output:\n%s\n", stdout))
	}

	return output.String(), nil
}

// TotalCleanup prunes all unused containers, images, volumes, and networks.
func TotalCleanup() (string, error) {
	stdout, stderr, err := engine.RunEngineCommand("system", "prune", "-a", "--volumes", "-f")
	if err != nil {
		return "", fmt.Errorf("failed to perform total cleanup: %s, %w", stderr, err)
	}
	return stdout, nil
}