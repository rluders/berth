// Package tui provides the Terminal User Interface for Berth.
package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/rluders/berth/internal/controller"
)

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
