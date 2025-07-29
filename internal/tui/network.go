// Package tui provides the Terminal User Interface for Berth.
package tui

import (
	"log/slog"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/rluders/berth/internal/controller"
)

// fetchNetworksCmd is a Bubble Tea command that fetches a list of networks.
func fetchNetworksCmd() tea.Cmd {
	return func() tea.Msg {
		slog.Debug("fetchNetworksCmd: Calling controller.ListNetworks...")
		networks, err := controller.ListNetworks()
		if err != nil {
			slog.Error("fetchNetworksCmd: Error listing networks", "error", err)
			return err
		}
		slog.Debug("fetchNetworksCmd: Successfully listed networks.")
		return networks
	}
}

// inspectNetworkCmd is a Bubble Tea command that inspects a network.
func inspectNetworkCmd(idOrName string) tea.Cmd {
	return func() tea.Msg {
		slog.Debug("inspectNetworkCmd: Calling controller.InspectNetwork", "idOrName", idOrName)
		output, err := controller.InspectNetwork(idOrName)
		if err != nil {
			slog.Error("inspectNetworkCmd: Error inspecting network", "idOrName", idOrName, "error", err)
			return err
		}
		slog.Debug("inspectNetworkCmd: Successfully inspected network.", "idOrName", idOrName)
		return output
	}
}
