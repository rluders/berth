// Package engine provides functionality for creating a Docker client.
package engine

import (
	"github.com/docker/docker/client"
)

// NewClient creates a new Docker client.
func NewClient() (*client.Client, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, err
	}
	return cli, nil
}
