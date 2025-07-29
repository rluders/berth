// Package tui provides the Terminal User Interface for Berth.
package tui

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/rluders/berth/internal/controller"
)

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
