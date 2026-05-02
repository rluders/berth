// Package controller provides the logic for interacting with container engines.
package controller

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/pkg/stdcopy"
	"time"

	"github.com/rluders/berth/internal/engine"
	"github.com/rluders/berth/internal/service"
)

var containerService service.ContainerService

func init() {
	cli, err := engine.NewClient()
	if err != nil {
		panic(fmt.Errorf("failed to create Docker client: %w", err))
	}
	containerService = service.NewContainerService(cli)
}

// Container represents a container's simplified information.
type Container struct {
	ID        string
	Image     string
	Command   string
	CreatedAt int64
	Status    string
	State     string
	Ports     string
	Names     string
	Labels    map[string]string
}

// ContainerDetails holds structured inspection data for the details view.
type ContainerDetails struct {
	ID       string
	Name     string
	Image    string
	Command  string
	Env      []string
	Ports    []PortBinding
	Mounts   []Mount
	Networks []NetworkEndpoint
	State    string
	Created  string
}

// PortBinding represents a single port mapping.
type PortBinding struct {
	ContainerPort string
	Protocol      string
	HostIP        string
	HostPort      string
}

// Mount represents a volume/bind mount.
type Mount struct {
	Type        string
	Source      string
	Destination string
	Mode        string
	RW          bool
}

// NetworkEndpoint represents a container's connection to a network.
type NetworkEndpoint struct {
	Name      string
	IPAddress string
	Gateway   string
}

// ContainerStat holds live resource usage for a container.
type ContainerStat struct {
	CPUPercent float64
	MemUsage   uint64
	MemLimit   uint64
}

// statsJSON is a minimal struct to decode Docker stats API response.
type statsJSON struct {
	CPUStats struct {
		CPUUsage struct {
			TotalUsage  uint64   `json:"total_usage"`
			PercpuUsage []uint64 `json:"percpu_usage"`
		} `json:"cpu_usage"`
		SystemCPUUsage uint64 `json:"system_cpu_usage"`
		OnlineCPUs     uint32 `json:"online_cpus"`
	} `json:"cpu_stats"`
	PreCPUStats struct {
		CPUUsage struct {
			TotalUsage uint64 `json:"total_usage"`
		} `json:"cpu_usage"`
		SystemCPUUsage uint64 `json:"system_cpu_usage"`
	} `json:"precpu_stats"`
	MemoryStats struct {
		Usage uint64            `json:"usage"`
		Limit uint64            `json:"limit"`
		Stats map[string]uint64 `json:"stats"`
	} `json:"memory_stats"`
}

// ListContainers lists all running and stopped containers.
func ListContainers() ([]Container, error) {
	containers, err := containerService.ListContainers(context.Background(), container.ListOptions{All: true})
	if err != nil {
		return nil, fmt.Errorf("failed to list containers: %w", err)
	}

	var result []Container
	for _, c := range containers {
		ports := formatPorts(c.Ports)
		result = append(result, Container{
			ID:        c.ID[:12],
			Image:     c.Image,
			Command:   c.Command,
			CreatedAt: c.Created,
			Status:    c.Status,
			State:     c.State,
			Ports:     ports,
			Names:     strings.TrimPrefix(strings.Join(c.Names, ","), "/"),
			Labels:    c.Labels,
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

// RestartContainer restarts a container by its ID or name.
func RestartContainer(idOrName string) error {
	return containerService.RestartContainer(context.Background(), idOrName, container.StopOptions{})
}

// RemoveContainer removes a container by its ID or name.
func RemoveContainer(idOrName string) error {
	return containerService.RemoveContainer(context.Background(), idOrName, container.RemoveOptions{Force: true})
}

// GetContainerLogs retrieves the logs of a container (one-shot).
func GetContainerLogs(idOrName string) (string, error) {
	out, err := containerService.ContainerLogs(context.Background(), idOrName, container.LogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Tail:       "500",
	})
	if err != nil {
		return "", fmt.Errorf("failed to get logs for container %s: %w", idOrName, err)
	}
	defer out.Close()

	var buf strings.Builder
	if _, err = stdcopy.StdCopy(&buf, &buf, out); err != nil {
		// Fallback for TTY containers (no multiplexing header)
		buf.Reset()
		if _, err2 := io.Copy(&buf, out); err2 != nil {
			return "", fmt.Errorf("failed to read logs: %w", err2)
		}
	}

	return buf.String(), nil
}

// StreamContainerLogs streams container logs line by line into ch, closing ch when done or ctx cancelled.
func StreamContainerLogs(ctx context.Context, idOrName string, ch chan<- string) {
	defer close(ch)

	out, err := containerService.ContainerLogs(ctx, idOrName, container.LogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Follow:     true,
		Tail:       "200",
	})
	if err != nil {
		return
	}
	defer out.Close()

	pr, pw := io.Pipe()
	go func() {
		defer pw.Close()
		if _, err := stdcopy.StdCopy(pw, pw, out); err != nil {
			// TTY container — fallback to raw copy
			if _, err2 := io.Copy(pw, out); err2 != nil {
				return
			}
		}
	}()

	scanner := bufio.NewScanner(pr)
	for scanner.Scan() {
		select {
		case <-ctx.Done():
			return
		case ch <- scanner.Text():
		}
	}
}

// InspectContainer inspects a container and returns raw JSON.
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

// GetContainerDetails returns structured inspection data for the details view.
func GetContainerDetails(idOrName string) (ContainerDetails, error) {
	inspect, err := containerService.ContainerInspect(context.Background(), idOrName)
	if err != nil {
		return ContainerDetails{}, fmt.Errorf("failed to inspect container %s: %w", idOrName, err)
	}

	name := inspect.Name
	if len(name) > 0 && name[0] == '/' {
		name = name[1:]
	}

	details := ContainerDetails{
		ID:      inspect.ID[:12],
		Name:    name,
		Image:   inspect.Config.Image,
		Command: strings.Join(inspect.Config.Cmd, " "),
		Env:     inspect.Config.Env,
		State:   inspect.State.Status,
		Created: formatCreated(inspect.Created),
	}

	for _, m := range inspect.Mounts {
		details.Mounts = append(details.Mounts, Mount{
			Type:        string(m.Type),
			Source:      m.Source,
			Destination: m.Destination,
			Mode:        m.Mode,
			RW:          m.RW,
		})
	}

	for netName, ep := range inspect.NetworkSettings.Networks {
		details.Networks = append(details.Networks, NetworkEndpoint{
			Name:      netName,
			IPAddress: ep.IPAddress,
			Gateway:   ep.Gateway,
		})
	}

	for portProto, bindings := range inspect.HostConfig.PortBindings {
		portStr := string(portProto)
		parts := strings.SplitN(portStr, "/", 2)
		containerPort := parts[0]
		protocol := ""
		if len(parts) > 1 {
			protocol = parts[1]
		}
		for _, b := range bindings {
			details.Ports = append(details.Ports, PortBinding{
				ContainerPort: containerPort,
				Protocol:      protocol,
				HostIP:        b.HostIP,
				HostPort:      b.HostPort,
			})
		}
	}

	return details, nil
}

// GetContainerStats returns one-shot CPU/memory stats for a container.
func GetContainerStats(idOrName string) (ContainerStat, error) {
	resp, err := containerService.ContainerStats(context.Background(), idOrName, false)
	if err != nil {
		return ContainerStat{}, err
	}
	defer resp.Body.Close()

	var s statsJSON
	if err := json.NewDecoder(resp.Body).Decode(&s); err != nil {
		return ContainerStat{}, err
	}

	cpuDelta := float64(s.CPUStats.CPUUsage.TotalUsage) - float64(s.PreCPUStats.CPUUsage.TotalUsage)
	systemDelta := float64(s.CPUStats.SystemCPUUsage) - float64(s.PreCPUStats.SystemCPUUsage)
	numCPUs := float64(s.CPUStats.OnlineCPUs)
	if numCPUs == 0 {
		numCPUs = float64(len(s.CPUStats.CPUUsage.PercpuUsage))
	}

	var cpuPercent float64
	if systemDelta > 0 && cpuDelta > 0 {
		cpuPercent = (cpuDelta / systemDelta) * numCPUs * 100.0
	}

	memUsage := s.MemoryStats.Usage
	if cache, ok := s.MemoryStats.Stats["cache"]; ok {
		memUsage -= cache
	}

	return ContainerStat{
		CPUPercent: cpuPercent,
		MemUsage:   memUsage,
		MemLimit:   s.MemoryStats.Limit,
	}, nil
}

// formatCreated parses a Docker RFC3339 created timestamp into a human-readable age.
func formatCreated(s string) string {
	t, err := time.Parse(time.RFC3339Nano, s)
	if err != nil {
		return s
	}
	d := time.Since(t)
	switch {
	case d < time.Minute:
		return fmt.Sprintf("%ds ago", int(d.Seconds()))
	case d < time.Hour:
		return fmt.Sprintf("%dm ago", int(d.Minutes()))
	case d < 24*time.Hour:
		return fmt.Sprintf("%dh ago", int(d.Hours()))
	default:
		return fmt.Sprintf("%dd ago", int(d.Hours()/24))
	}
}

// formatPorts converts Docker port list to a compact string.
func formatPorts(ports []container.Port) string {
	if len(ports) == 0 {
		return ""
	}
	seen := make(map[string]bool)
	var parts []string
	for _, p := range ports {
		var s string
		if p.PublicPort > 0 {
			s = fmt.Sprintf("%d->%d/%s", p.PublicPort, p.PrivatePort, p.Type)
		} else {
			s = fmt.Sprintf("%d/%s", p.PrivatePort, p.Type)
		}
		if !seen[s] {
			seen[s] = true
			parts = append(parts, s)
		}
	}
	return strings.Join(parts, ", ")
}
