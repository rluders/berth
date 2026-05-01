package tui

import (
	"fmt"
	"log/slog"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/rluders/berth/internal/controller"
)

func fetchContainersCmd() tea.Cmd {
	return func() tea.Msg {
		slog.Debug("fetchContainersCmd called")
		containers, err := controller.ListContainers()
		if err != nil {
			slog.Error("fetchContainersCmd error", "error", err)
			return errMsg{err}
		}
		return containerListMsg(containers)
	}
}

func fetchImagesCmd() tea.Cmd {
	return func() tea.Msg {
		slog.Debug("fetchImagesCmd called")
		images, err := controller.ListImages()
		if err != nil {
			slog.Error("fetchImagesCmd error", "error", err)
			return errMsg{err}
		}
		return imageListMsg(images)
	}
}

func fetchVolumesCmd() tea.Cmd {
	return func() tea.Msg {
		slog.Debug("fetchVolumesCmd called")
		volumes, err := controller.ListVolumes()
		if err != nil {
			slog.Error("fetchVolumesCmd error", "error", err)
			return errMsg{err}
		}
		return volumeListMsg(volumes)
	}
}

func fetchNetworksCmd() tea.Cmd {
	return func() tea.Msg {
		slog.Debug("fetchNetworksCmd called")
		networks, err := controller.ListNetworks()
		if err != nil {
			slog.Error("fetchNetworksCmd error", "error", err)
			return errMsg{err}
		}
		return networkListMsg(networks)
	}
}

func fetchSystemInfoCmd() tea.Cmd {
	return func() tea.Msg {
		slog.Debug("fetchSystemInfoCmd called")
		info, err := controller.GetSystemInfo()
		if err != nil {
			slog.Error("fetchSystemInfoCmd error", "error", err)
			return errMsg{err}
		}
		return systemInfoMsg(info)
	}
}

func fetchAllCmd() tea.Cmd {
	return tea.Batch(
		fetchContainersCmd(),
		fetchImagesCmd(),
		fetchVolumesCmd(),
		fetchNetworksCmd(),
		fetchSystemInfoCmd(),
	)
}

func startContainerCmd(idOrName string) tea.Cmd {
	return func() tea.Msg {
		slog.Debug("startContainerCmd called", "id", idOrName)
		if err := controller.StartContainer(idOrName); err != nil {
			slog.Error("startContainerCmd error", "id", idOrName, "error", err)
			return errMsg{err}
		}
		return statusMsg(fmt.Sprintf("Container %s started.", idOrName))
	}
}

func stopContainerCmd(idOrName string) tea.Cmd {
	return func() tea.Msg {
		slog.Debug("stopContainerCmd called", "id", idOrName)
		if err := controller.StopContainer(idOrName); err != nil {
			slog.Error("stopContainerCmd error", "id", idOrName, "error", err)
			return errMsg{err}
		}
		return statusMsg(fmt.Sprintf("Container %s stopped.", idOrName))
	}
}

func removeContainerCmd(idOrName string) tea.Cmd {
	return func() tea.Msg {
		slog.Debug("removeContainerCmd called", "id", idOrName)
		if err := controller.RemoveContainer(idOrName); err != nil {
			slog.Error("removeContainerCmd error", "id", idOrName, "error", err)
			return errMsg{err}
		}
		return statusMsg(fmt.Sprintf("Container %s removed.", idOrName))
	}
}

func getLogsCmd(idOrName string) tea.Cmd {
	return func() tea.Msg {
		slog.Debug("getLogsCmd called", "id", idOrName)
		logs, err := controller.GetContainerLogs(idOrName)
		if err != nil {
			slog.Error("getLogsCmd error", "id", idOrName, "error", err)
			return errMsg{err}
		}
		return logsMsg(logs)
	}
}

func inspectContainerCmd(idOrName string) tea.Cmd {
	return func() tea.Msg {
		slog.Debug("inspectContainerCmd called", "id", idOrName)
		output, err := controller.InspectContainer(idOrName)
		if err != nil {
			slog.Error("inspectContainerCmd error", "id", idOrName, "error", err)
			return errMsg{err}
		}
		return inspectMsg(output)
	}
}

func removeImageCmd(idOrName string) tea.Cmd {
	return func() tea.Msg {
		slog.Debug("removeImageCmd called", "id", idOrName)
		if err := controller.RemoveImage(idOrName); err != nil {
			slog.Error("removeImageCmd error", "id", idOrName, "error", err)
			return errMsg{err}
		}
		return statusMsg(fmt.Sprintf("Image %s removed.", idOrName))
	}
}

func removeVolumeCmd(name string) tea.Cmd {
	return func() tea.Msg {
		slog.Debug("removeVolumeCmd called", "name", name)
		if err := controller.RemoveVolume(name); err != nil {
			slog.Error("removeVolumeCmd error", "name", name, "error", err)
			return errMsg{err}
		}
		return statusMsg(fmt.Sprintf("Volume %s removed.", name))
	}
}

func inspectNetworkCmd(idOrName string) tea.Cmd {
	return func() tea.Msg {
		slog.Debug("inspectNetworkCmd called", "id", idOrName)
		output, err := controller.InspectNetwork(idOrName)
		if err != nil {
			slog.Error("inspectNetworkCmd error", "id", idOrName, "error", err)
			return errMsg{err}
		}
		return inspectMsg(output)
	}
}

func basicCleanupCmd() tea.Cmd {
	return func() tea.Msg {
		slog.Debug("basicCleanupCmd called")
		output, err := controller.BasicCleanup()
		if err != nil {
			slog.Error("basicCleanupCmd error", "error", err)
			return errMsg{err}
		}
		return statusMsg(fmt.Sprintf("Basic cleanup: %s", output))
	}
}

func advancedCleanupCmd() tea.Cmd {
	return func() tea.Msg {
		slog.Debug("advancedCleanupCmd called")
		output, err := controller.AdvancedCleanup()
		if err != nil {
			slog.Error("advancedCleanupCmd error", "error", err)
			return errMsg{err}
		}
		return statusMsg(fmt.Sprintf("Advanced cleanup: %s", output))
	}
}

func totalCleanupCmd() tea.Cmd {
	return func() tea.Msg {
		slog.Debug("totalCleanupCmd called")
		output, err := controller.TotalCleanup()
		if err != nil {
			slog.Error("totalCleanupCmd error", "error", err)
			return errMsg{err}
		}
		return statusMsg(fmt.Sprintf("Total cleanup: %s", output))
	}
}
