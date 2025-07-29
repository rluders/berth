// Package tui provides the Terminal User Interface for Berth.
package tui

import (
	"fmt"
	"log/slog"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/rluders/berth/internal/controller"
)

// fetchSystemInfoCmd is a Bubble Tea command that fetches system information.
func fetchSystemInfoCmd() tea.Cmd {
	return func() tea.Msg {
		slog.Debug("fetchSystemInfoCmd: Calling controller.GetSystemInfo...")
		info, err := controller.GetSystemInfo()
		if err != nil {
			slog.Error("fetchSystemInfoCmd: Error getting system info", "error", err)
			return err
		}
		slog.Debug("fetchSystemInfoCmd: Successfully retrieved system info.")
		return info
	}
}

// basicCleanupCmd is a Bubble Tea command that performs basic cleanup.
func basicCleanupCmd() tea.Cmd {
	return func() tea.Msg {
		slog.Debug("basicCleanupCmd: Calling controller.BasicCleanup...")
		output, err := controller.BasicCleanup()
		if err != nil {
			slog.Error("basicCleanupCmd: Error during basic cleanup", "error", err)
			return err
		}
		slog.Debug("basicCleanupCmd: Basic cleanup completed.", "output", output)
		return statusMsg(fmt.Sprintf("Basic cleanup: %s", output))
	}
}

// advancedCleanupCmd is a Bubble Tea command that performs advanced cleanup.
func advancedCleanupCmd() tea.Cmd {
	return func() tea.Msg {
		slog.Debug("advancedCleanupCmd: Calling controller.AdvancedCleanup...")
		output, err := controller.AdvancedCleanup()
		if err != nil {
			slog.Error("advancedCleanupCmd: Error during advanced cleanup", "error", err)
			return err
		}
		slog.Debug("advancedCleanupCmd: Advanced cleanup completed.", "output", output)
		return statusMsg(fmt.Sprintf("Advanced cleanup: %s", output))
	}
}

// totalCleanupCmd is a Bubble Tea command that performs total cleanup.
func totalCleanupCmd() tea.Cmd {
	return func() tea.Msg {
		slog.Debug("totalCleanupCmd: Calling controller.TotalCleanup...")
		output, err := controller.TotalCleanup()
		if err != nil {
			slog.Error("totalCleanupCmd: Error during total cleanup", "error", err)
			return err
		}
		slog.Debug("totalCleanupCmd: Total cleanup completed.", "output", output)
		return statusMsg(fmt.Sprintf("Total cleanup: %s", output))
	}
}
