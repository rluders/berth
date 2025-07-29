// Package tui provides the Terminal User Interface for Berth.
package tui

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/rluders/berth/internal/controller"
)

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
