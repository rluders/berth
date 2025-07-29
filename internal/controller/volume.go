// Package controller provides the logic for interacting with container engines.
package controller

import (
	"context"
	"fmt"

	"github.com/docker/docker/api/types/volume"
	"github.com/rluders/berth/internal/engine"
	"github.com/rluders/berth/internal/service"
)

var volumeService service.VolumeService

func init() {
	cli, err := engine.NewClient()
	if err != nil {
		panic(fmt.Errorf("failed to create Docker client: %w", err))
	}
	volumeService = service.NewVolumeService(cli)
}

// Volume represents a volume's simplified information.
type Volume struct {
	Name       string
	Driver     string
	Scope      string
	Mountpoint string
}

// ListVolumes lists all volumes.
func ListVolumes() ([]Volume, error) {
	volumes, err := volumeService.VolumeList(context.Background(), volume.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list volumes: %w", err)
	}

	var result []Volume
	for _, v := range volumes.Volumes {
		result = append(result, Volume{
			Name:       v.Name,
			Driver:     v.Driver,
			Mountpoint: v.Mountpoint,
		})
	}

	return result, nil
}

// RemoveVolume removes a volume by its name.
func RemoveVolume(name string) error {
	return volumeService.VolumeRemove(context.Background(), name, false)
}
