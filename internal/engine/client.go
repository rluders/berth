// Package engine provides functionality for creating a Docker client.
package engine

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/docker/docker/client"
)

// NewClient creates a Docker/Podman client using the detected engine.
// For Podman, uses the user socket path when DOCKER_HOST is not set.
func NewClient() (*client.Client, error) {
	opts := []client.Opt{client.WithAPIVersionNegotiation()}

	if detectedEngine == Podman && os.Getenv("DOCKER_HOST") == "" {
		socketPath := podmanSocketPath()
		if socketPath != "" {
			opts = append(opts, client.WithHost("unix://"+socketPath))
		}
	} else {
		opts = append(opts, client.FromEnv)
	}

	cli, err := client.NewClientWithOpts(opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create client for %s: %w", detectedEngine, err)
	}
	return cli, nil
}

// podmanSocketPath returns the Podman socket path for the current user.
func podmanSocketPath() string {
	// Rootless Podman (preferred)
	if xdg := os.Getenv("XDG_RUNTIME_DIR"); xdg != "" {
		p := filepath.Join(xdg, "podman", "podman.sock")
		if _, err := os.Stat(p); err == nil {
			return p
		}
	}
	// Fallback: /run/user/<uid>/podman/podman.sock
	uid := fmt.Sprintf("%d", os.Getuid())
	p := filepath.Join("/run/user", uid, "podman", "podman.sock")
	if _, err := os.Stat(p); err == nil {
		return p
	}
	// Root Podman
	if _, err := os.Stat("/run/podman/podman.sock"); err == nil {
		return "/run/podman/podman.sock"
	}
	return ""
}
