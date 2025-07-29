// Package tui provides the Terminal User Interface for Berth.
package tui

import (
	"fmt"
	"log/slog"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/rluders/berth/internal/controller"
)

// fetchImagesCmd is a Bubble Tea command that fetches a list of images.
func fetchImagesCmd() tea.Cmd {
	return func() tea.Msg {
		slog.Debug("fetchImagesCmd: Calling controller.ListImages...")
		images, err := controller.ListImages()
		if err != nil {
			slog.Error("fetchImagesCmd: Error listing images", "error", err)
			return err
		}
		slog.Debug("fetchImagesCmd: Successfully listed images.")
		return images
	}
}

// removeImageCmd is a Bubble Tea command that removes an image.
func removeImageCmd(idOrName string) tea.Cmd {
	return func() tea.Msg {
		slog.Debug("removeImageCmd: Calling controller.RemoveImage", "idOrName", idOrName)
		err := controller.RemoveImage(idOrName)
		if err != nil {
			slog.Error("removeImageCmd: Error removing image", "idOrName", idOrName, "error", err)
			return err
		}
		slog.Debug("removeImageCmd: Successfully removed image.", "idOrName", idOrName)
		return statusMsg(fmt.Sprintf("Image %s removed.", idOrName))
	}
}
