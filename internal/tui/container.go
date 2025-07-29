// Package tui provides the Terminal User Interface for Berth.
package tui

import (
	"fmt"
	"log/slog"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/rluders/berth/internal/controller"
)

// fetchContainersCmd is a Bubble Tea command that fetches a list of containers.
func fetchContainersCmd() tea.Cmd {
	return func() tea.Msg {
		slog.Debug("fetchContainersCmd: Calling controller.ListContainers...")
		containers, err := controller.ListContainers()
		if err != nil {
			slog.Error("fetchContainersCmd: Error listing containers", "error", err)
			return err
		}
		slog.Debug("fetchContainersCmd: Successfully listed containers.")
		return containers
	}
}

// startContainerCmd is a Bubble Tea command that starts a container.
func startContainerCmd(idOrName string) tea.Cmd {
	return func() tea.Msg {
		slog.Debug("startContainerCmd: Calling controller.StartContainer", "idOrName", idOrName)
		err := controller.StartContainer(idOrName)
		if err != nil {
			slog.Error("startContainerCmd: Error starting container", "idOrName", idOrName, "error", err)
			return err
		}
		slog.Debug("startContainerCmd: Successfully started container.", "idOrName", idOrName)
		return statusMsg(fmt.Sprintf("Container %s started.", idOrName))
	}
}

// stopContainerCmd is a Bubble Tea command that stops a container.
func stopContainerCmd(idOrName string) tea.Cmd {
	return func() tea.Msg {
		slog.Debug("stopContainerCmd: Calling controller.StopContainer", "idOrName", idOrName)
		err := controller.StopContainer(idOrName)
		if err != nil {
			slog.Error("stopContainerCmd: Error stopping container", "idOrName", idOrName, "error", err)
			return err
		}
		slog.Debug("stopContainerCmd: Successfully stopped container.", "idOrName", idOrName)
		return statusMsg(fmt.Sprintf("Container %s stopped.", idOrName))
	}
}

// removeContainerCmd is a Bubble Tea command that removes a container.
func removeContainerCmd(idOrName string) tea.Cmd {
	return func() tea.Msg {
		slog.Debug("removeContainerCmd: Calling controller.RemoveContainer", "idOrName", idOrName)
		err := controller.RemoveContainer(idOrName)
		if err != nil {
			slog.Error("removeContainerCmd: Error removing container", "idOrName", idOrName, "error", err)
			return err
		}
		slog.Debug("removeContainerCmd: Successfully removed container.", "idOrName", idOrName)
		return statusMsg(fmt.Sprintf("Container %s removed.", idOrName))
	}
}

// getLogsCmd is a Bubble Tea command that fetches logs for a container.
func getLogsCmd(idOrName string) tea.Cmd {
	return func() tea.Msg {
		slog.Debug("getLogsCmd: Calling controller.GetContainerLogs", "idOrName", idOrName)
		logs, err := controller.GetContainerLogs(idOrName)
		if err != nil {
			slog.Error("getLogsCmd: Error getting container logs", "idOrName", idOrName, "error", err)
			return err
		}
		slog.Debug("getLogsCmd: Successfully retrieved container logs.", "idOrName", idOrName)
		return logs
	}
}

// inspectContainerCmd is a Bubble Tea command that inspects a container.
func inspectContainerCmd(idOrName string) tea.Cmd {
	return func() tea.Msg {
		slog.Debug("inspectContainerCmd: Calling controller.InspectContainer", "idOrName", idOrName)
		output, err := controller.InspectContainer(idOrName)
		if err != nil {
			slog.Error("inspectContainerCmd: Error inspecting container", "idOrName", idOrName, "error", err)
			return err
		}
		slog.Debug("inspectContainerCmd: Successfully inspected container.", "idOrName", idOrName)
		return output
	}
}
