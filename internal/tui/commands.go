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

// fetchImagesCmd is a Bubble Tea command that fetches a list of images.
func fetchImagesCmd() tea.Cmd {
	return func() tea.Msg {
		images, err := controller.ListImages()
		if err != nil {
			return err
		}
		return images
	}
}

// fetchVolumesCmd is a Bubble Tea command that fetches a list of volumes.
func fetchVolumesCmd() tea.Cmd {
	return func() tea.Msg {
		volumes, err := controller.ListVolumes()
		if err != nil {
			return err
		}
		return volumes
	}
}

// fetchNetworksCmd is a Bubble Tea command that fetches a list of networks.
func fetchNetworksCmd() tea.Cmd {
	return func() tea.Msg {
		networks, err := controller.ListNetworks()
		if err != nil {
			return err
		}
		return networks
	}
}

// fetchSystemInfoCmd is a Bubble Tea command that fetches system information.
func fetchSystemInfoCmd() tea.Cmd {
	return func() tea.Msg {
		systemInfo, err := controller.GetSystemInfo()
		if err != nil {
			return err
		}
		return systemInfo
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

// removeImageCmd is a Bubble Tea command that removes an image.
func removeImageCmd(idOrName string) tea.Cmd {
	return func() tea.Msg {
		err := controller.RemoveImage(idOrName)
		if err != nil {
			return err
		}
		return statusMsg(fmt.Sprintf("Image %s removed.", idOrName))
	}
}

// removeVolumeCmd is a Bubble Tea command that removes a volume.
func removeVolumeCmd(name string) tea.Cmd {
	return func() tea.Msg {
		err := controller.RemoveVolume(name)
		if err != nil {
			return err
		}
		return statusMsg(fmt.Sprintf("Volume %s removed.", name))
	}
}

// inspectNetworkCmd is a Bubble Tea command that inspects a network.
func inspectNetworkCmd(idOrName string) tea.Cmd {
	return func() tea.Msg {
		output, err := controller.InspectNetwork(idOrName)
		if err != nil {
			return err
		}
		return output
	}
}

// basicCleanupCmd is a Bubble Tea command that performs basic cleanup.
func basicCleanupCmd() tea.Cmd {
	return func() tea.Msg {
		output, err := controller.BasicCleanup()
		if err != nil {
			return err
		}
		return statusMsg(fmt.Sprintf("Basic Cleanup completed:\n%s", output))
	}
}

// advancedCleanupCmd is a Bubble Tea command that performs advanced cleanup.
func advancedCleanupCmd() tea.Cmd {
	return func() tea.Msg {
		output, err := controller.AdvancedCleanup()
		if err != nil {
			return err
		}
		return statusMsg(fmt.Sprintf("Advanced Cleanup completed:\n%s", output))
	}
}

// totalCleanupCmd is a Bubble Tea command that performs total cleanup.
func totalCleanupCmd() tea.Cmd {
	return func() tea.Msg {
		output, err := controller.TotalCleanup()
		if err != nil {
			return err
		}
		return statusMsg(fmt.Sprintf("Total Cleanup completed:\n%s", output))
	}
}

