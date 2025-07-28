package controller

import (
	"fmt"
	"strings"

	"github.com/rluders/berth/internal/engine"
)

type Image struct {
	ID         string
	Repository string
	Tag        string
	Size       string
	Created    string
}

// ListImages lists all images.
func ListImages() ([]Image, error) {
	stdout, stderr, err := engine.RunEngineCommand("images", "--format", "{{.ID}}\t{{.Repository}}\t{{.Tag}}\t{{.Size}}\t{{.CreatedAt}}")
	if err != nil {
		return nil, fmt.Errorf("failed to list images: %s, %w", stderr, err)
	}

	lines := strings.Split(strings.TrimSpace(stdout), "\n")
	var images []Image
	for _, line := range lines {
		if line == "" {
			continue
		}
		fields := strings.Split(line, "\t")
		if len(fields) != 5 {
			// Log or handle malformed line
			continue
		}
		images = append(images, Image{
			ID:         fields[0],
			Repository: fields[1],
			Tag:        fields[2],
			Size:       fields[3],
			Created:    fields[4],
		})
	}
	return images, nil
}

// RemoveImage removes an image by its ID or name.
func RemoveImage(idOrName string) error {
	_, stderr, err := engine.RunEngineCommand("rmi", idOrName)
	if err != nil {
		return fmt.Errorf("failed to remove image %s: %s, %w", idOrName, stderr, err)
	}
	return nil
}
