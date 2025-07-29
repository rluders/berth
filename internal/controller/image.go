// Package controller provides the logic for interacting with container engines.
package controller

import (
	"context"
	"fmt"

	dockerImageTypes "github.com/docker/docker/api/types/image"
	"github.com/rluders/berth/internal/engine"
	"github.com/rluders/berth/internal/service"
)

var imageService service.ImageService

func init() {
	cli, err := engine.NewClient()
	if err != nil {
		panic(fmt.Errorf("failed to create Docker client: %w", err))
	}
	imageService = service.NewImageService(cli)
}

// Image represents an image's simplified information.
type Image struct {
	ID         string
	Repository string
	Tag        string
	Size       string
	Created    string
}

// ListImages lists all images.
func ListImages() ([]Image, error) {
	images, err := imageService.ImageList(context.Background(), dockerImageTypes.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list images: %w", err)
	}

	var result []Image
	for _, i := range images {
		result = append(result, Image{
			ID:         i.ID[7:19],
			Repository: i.RepoTags[0],
			Tag:        i.RepoTags[0],
			Size:       fmt.Sprintf("%d", i.Size),
			Created:    fmt.Sprintf("%d", i.Created),
		})
	}

	return result, nil
}

// RemoveImage removes an image by its ID or name.
func RemoveImage(idOrName string) error {
	_, err := imageService.ImageRemove(context.Background(), idOrName, dockerImageTypes.RemoveOptions{})
	return err
}
