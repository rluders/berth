package controller

import (
	"fmt"
	"os/exec"
)

func runCompose(project, workDir string, args ...string) error {
	baseArgs := []string{"compose", "-p", project}
	cmd := exec.Command("docker", append(baseArgs, args...)...)
	if workDir != "" {
		cmd.Dir = workDir
	}
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%w\n%s", err, string(out))
	}
	return nil
}

func ComposeUp(project, workDir string) error {
	return runCompose(project, workDir, "up", "-d")
}

func ComposeUpBuild(project, workDir string) error {
	return runCompose(project, workDir, "up", "-d", "--build")
}

func ComposeRecreate(project, workDir string) error {
	return runCompose(project, workDir, "up", "-d", "--force-recreate")
}

func ComposeDown(project, workDir string) error {
	return runCompose(project, workDir, "down")
}

func ComposePull(project, workDir string) error {
	return runCompose(project, workDir, "pull")
}

func ComposeBuild(project, workDir string) error {
	return runCompose(project, workDir, "build")
}
