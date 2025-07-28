package tui

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/rluders/berth/internal/controller"
)

func fetchContainersCmd() tea.Cmd {
	return func() tea.Msg {
		containers, err := controller.ListContainers()
		if err != nil {
			return err
		}
		return containers
	}
}

func fetchImagesCmd() tea.Cmd {
	return func() tea.Msg {
		images, err := controller.ListImages()
		if err != nil {
			return err
		}
		return images
	}
}

func fetchVolumesCmd() tea.Cmd {
	return func() tea.Msg {
		volumes, err := controller.ListVolumes()
		if err != nil {
			return err
		}
		return volumes
	}
}

func fetchNetworksCmd() tea.Cmd {
	return func() tea.Msg {
		networks, err := controller.ListNetworks()
		if err != nil {
			return err
		}
		return networks
	}
}

func fetchSystemInfoCmd() tea.Cmd {
	return func() tea.Msg {
		systemInfo, err := controller.GetSystemInfo()
		if err != nil {
			return err
		}
		return systemInfo
	}
}

func startContainerCmd(idOrName string) tea.Cmd {
	return func() tea.Msg {
		err := controller.StartContainer(idOrName)
		if err != nil {
			return err
		}
		return statusMsg(fmt.Sprintf("Container %s started.", idOrName))
	}
}

func stopContainerCmd(idOrName string) tea.Cmd {
	return func() tea.Msg {
		err := controller.StopContainer(idOrName)
		if err != nil {
			return err
		}
		return statusMsg(fmt.Sprintf("Container %s stopped.", idOrName))
	}
}

func removeContainerCmd(idOrName string) tea.Cmd {
	return func() tea.Msg {
		err := controller.RemoveContainer(idOrName)
		if err != nil {
			return err
		}
		return statusMsg(fmt.Sprintf("Container %s removed.", idOrName))
	}
}

func getLogsCmd(idOrName string) tea.Cmd {
	return func() tea.Msg {
		logs, err := controller.GetContainerLogs(idOrName)
		if err != nil {
			return err
		}
		return logs
	}
}

func inspectContainerCmd(idOrName string) tea.Cmd {
	return func() tea.Msg {
		output, err := controller.InspectContainer(idOrName)
		if err != nil {
			return err
		}
		return output
	}
}

func removeImageCmd(idOrName string) tea.Cmd {
	return func() tea.Msg {
		err := controller.RemoveImage(idOrName)
		if err != nil {
			return err
		}
		return statusMsg(fmt.Sprintf("Image %s removed.", idOrName))
	}
}

func removeVolumeCmd(name string) tea.Cmd {
	return func() tea.Msg {
		err := controller.RemoveVolume(name)
		if err != nil {
			return err
		}
		return statusMsg(fmt.Sprintf("Volume %s removed.", name))
	}
}

func inspectNetworkCmd(idOrName string) tea.Cmd {
	return func() tea.Msg {
		output, err := controller.InspectNetwork(idOrName)
		if err != nil {
			return err
		}
		return output
	}
}

func basicCleanupCmd() tea.Cmd {
	return func() tea.Msg {
		output, err := controller.BasicCleanup()
		if err != nil {
			return err
		}
		return statusMsg(fmt.Sprintf("Basic Cleanup completed:\n%s", output))
	}
}

func advancedCleanupCmd() tea.Cmd {
	return func() tea.Msg {
		output, err := controller.AdvancedCleanup()
		if err != nil {
			return err
		}
		return statusMsg(fmt.Sprintf("Advanced Cleanup completed:\n%s", output))
	}
}

func totalCleanupCmd() tea.Cmd {
	return func() tea.Msg {
		output, err := controller.TotalCleanup()
		if err != nil {
			return err
		}
		return statusMsg(fmt.Sprintf("Total Cleanup completed:\n%s", output))
	}
}
