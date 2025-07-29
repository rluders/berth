// Package tui provides the Terminal User Interface for Berth.
package tui

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/rluders/berth/internal/controller"
)

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