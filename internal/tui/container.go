// Package tui provides the Terminal User Interface for Berth.
package tui

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/rluders/berth/internal/controller"
)

// fetchContainersCmd is a Bubble Tea command that fetches a list of containers.
func fetchContainersCmd() tea.Cmd {
	return func() tea.Msg {
		containers, err := controller.ListContainers()
		if err != nil {
			return err
		}
		return containers
	}
}

// startContainerCmd is a Bubble Tea command that starts a container.
func startContainerCmd(idOrName string) tea.Cmd {
	return func() tea.Msg {
		err := controller.StartContainer(idOrName)
		if err != nil {
			return err
		}
		return statusMsg(fmt.Sprintf("Container %s started.", idOrName))
	}
}

// stopContainerCmd is a Bubble Tea command that stops a container.
func stopContainerCmd(idOrName string) tea.Cmd {
	return func() tea.Msg {
		err := controller.StopContainer(idOrName)
		if err != nil {
			return err
		}
		return statusMsg(fmt.Sprintf("Container %s stopped.", idOrName))
	}
}

// removeContainerCmd is a Bubble Tea command that removes a container.
func removeContainerCmd(idOrName string) tea.Cmd {
	return func() tea.Msg {
		err := controller.RemoveContainer(idOrName)
		if err != nil {
			return err
		}
		return statusMsg(fmt.Sprintf("Container %s removed.", idOrName))
	}
}

// getLogsCmd is a Bubble Tea command that fetches logs for a container.
func getLogsCmd(idOrName string) tea.Cmd {
	return func() tea.Msg {
		logs, err := controller.GetContainerLogs(idOrName)
		if err != nil {
			return err
		}
		return logs
	}
}

// inspectContainerCmd is a Bubble Tea command that inspects a container.
func inspectContainerCmd(idOrName string) tea.Cmd {
	return func() tea.Msg {
		output, err := controller.InspectContainer(idOrName)
		if err != nil {
			return err
		}
		return output
	}
}
