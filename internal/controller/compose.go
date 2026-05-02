package controller

import (
	"bufio"
	"context"
	"io"
	"os/exec"
)

// StreamCompose runs a compose command and fans stdout+stderr line-by-line into ch.
// ch is closed when the process exits or ctx is cancelled.
func StreamCompose(ctx context.Context, project, workDir string, ch chan<- string, args ...string) error {
	baseArgs := []string{"compose", "-p", project}
	cmd := exec.CommandContext(ctx, "docker", append(baseArgs, args...)...)
	if workDir != "" {
		cmd.Dir = workDir
	}

	pr, pw := io.Pipe()
	cmd.Stdout = pw
	cmd.Stderr = pw

	if err := cmd.Start(); err != nil {
		close(ch)
		return err
	}

	go func() {
		defer pw.Close()
		cmd.Wait() //nolint:errcheck
	}()

	go func() {
		defer close(ch)
		scanner := bufio.NewScanner(pr)
		for scanner.Scan() {
			select {
			case <-ctx.Done():
				return
			case ch <- scanner.Text():
			}
		}
	}()

	return nil
}

func ComposeUp(ctx context.Context, project, workDir string, ch chan<- string) error {
	return StreamCompose(ctx, project, workDir, ch, "up", "-d")
}

func ComposeUpBuild(ctx context.Context, project, workDir string, ch chan<- string) error {
	return StreamCompose(ctx, project, workDir, ch, "up", "-d", "--build")
}

func ComposeRecreate(ctx context.Context, project, workDir string, ch chan<- string) error {
	return StreamCompose(ctx, project, workDir, ch, "up", "-d", "--force-recreate")
}

func ComposeDown(ctx context.Context, project, workDir string, ch chan<- string) error {
	return StreamCompose(ctx, project, workDir, ch, "down")
}

func ComposePull(ctx context.Context, project, workDir string, ch chan<- string) error {
	return StreamCompose(ctx, project, workDir, ch, "pull")
}

func ComposeBuild(ctx context.Context, project, workDir string, ch chan<- string) error {
	return StreamCompose(ctx, project, workDir, ch, "build")
}
