package tui

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/rluders/berth/internal/controller"
)

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case spinner.TickMsg:
		var spinCmd tea.Cmd
		m.spinner, spinCmd = m.spinner.Update(msg)
		cmds = append(cmds, spinCmd)

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			if m.currentView == InspectView || m.currentView == LogsView {
				m.popView()
				m.logReady = false // Reset logReady when exiting logs view
				return m, nil
			}
			return m, tea.Quit
		case "1":
			m.currentView = ContainersView
			return m, nil
		case "2":
			m.currentView = ImagesView
			return m, nil
		case "3":
			m.currentView = VolumesView
			return m, nil
		case "4":
			m.currentView = NetworksView
			return m, nil
		case "5":
			m.currentView = SystemView
			return m, nil
		}

		if m.currentView == ContainersView {
			m.containerTable, cmd = m.containerTable.Update(msg)
			cmds = append(cmds, cmd)
			switch msg.String() {
			case "s": // Start container
				if len(m.containerTable.SelectedRow()) > 0 {
					containerID := m.containerTable.SelectedRow()[0]
					m.statusMessage = fmt.Sprintf("Starting container %s...", containerID)
					m.showSpinner = true
					cmds = append(cmds, startContainerCmd(containerID), m.spinner.Tick)
				}
			case "x": // Stop container
				if len(m.containerTable.SelectedRow()) > 0 {
					containerID := m.containerTable.SelectedRow()[0]
					m.statusMessage = fmt.Sprintf("Stopping container %s...", containerID)
					m.showSpinner = true
					cmds = append(cmds, stopContainerCmd(containerID), m.spinner.Tick)
				}
			case "d": // Remove container
				if len(m.containerTable.SelectedRow()) > 0 {
					containerID := m.containerTable.SelectedRow()[0]
					m.statusMessage = fmt.Sprintf("Removing container %s...", containerID)
					m.showSpinner = true
					cmds = append(cmds, removeContainerCmd(containerID), m.spinner.Tick)
				}
			case "l": // View logs
				if len(m.containerTable.SelectedRow()) > 0 {
					containerID := m.containerTable.SelectedRow()[0]
					m.pushView(LogsView)
					m.logReady = false // Reset logReady for new logs
					m.statusMessage = fmt.Sprintf("Fetching logs for %s...", containerID)
					m.showSpinner = true
					m.currentLogContainerID = containerID // Store the container ID
					cmds = append(cmds, getLogsCmd(containerID), m.spinner.Tick)
				}
			case "i": // Inspect container
				if len(m.containerTable.SelectedRow()) > 0 {
					containerID := m.containerTable.SelectedRow()[0]
					m.pushView(InspectView)
					m.currentInspectID = containerID
					m.statusMessage = fmt.Sprintf("Inspecting container %s...", containerID)
					m.showSpinner = true
					cmds = append(cmds, inspectContainerCmd(containerID), m.spinner.Tick)
				}
			}
		} else if m.currentView == ImagesView {
			m.imageTable, cmd = m.imageTable.Update(msg)
			cmds = append(cmds, cmd)
			switch msg.String() {
			case "d": // Remove image
				if len(m.imageTable.SelectedRow()) > 0 {
					imageID := m.imageTable.SelectedRow()[0]
					m.statusMessage = fmt.Sprintf("Removing image %s...", imageID)
					m.showSpinner = true
					cmds = append(cmds, removeImageCmd(imageID), m.spinner.Tick)
				}
			}
		} else if m.currentView == VolumesView {
			m.volumeTable, cmd = m.volumeTable.Update(msg)
			cmds = append(cmds, cmd)
			switch msg.String() {
			case "d": // Remove volume
				if len(m.volumeTable.SelectedRow()) > 0 {
					volumeName := m.volumeTable.SelectedRow()[0]
					m.statusMessage = fmt.Sprintf("Removing volume %s...", volumeName)
					m.showSpinner = true
					cmds = append(cmds, removeVolumeCmd(volumeName), m.spinner.Tick)
				}
			}
		} else if m.currentView == NetworksView {
			m.networkTable, cmd = m.networkTable.Update(msg)
			cmds = append(cmds, cmd)
			switch msg.String() {
			case "i": // Inspect network
				if len(m.networkTable.SelectedRow()) > 0 {
					networkID := m.networkTable.SelectedRow()[0]
					m.currentView = InspectView // Re-use inspect view for network inspect
					m.statusMessage = fmt.Sprintf("Inspecting network %s...", networkID)
					m.showSpinner = true
					cmds = append(cmds, inspectNetworkCmd(networkID), m.spinner.Tick)
				}
			}
		} else if m.currentView == SystemView {
			switch msg.String() {
			case "b": // Basic Cleanup
				m.statusMessage = "Performing basic cleanup..."
				m.showSpinner = true
				cmds = append(cmds, basicCleanupCmd(), m.spinner.Tick)
			case "a": // Advanced Cleanup
				m.statusMessage = "Performing advanced cleanup..."
				m.showSpinner = true
				cmds = append(cmds, advancedCleanupCmd(), m.spinner.Tick)
			case "t": // Total Cleanup
				m.statusMessage = "Performing total cleanup..."
				m.showSpinner = true
				cmds = append(cmds, totalCleanupCmd(), m.spinner.Tick)
			}
		} else if m.currentView == InspectView {
			// Delegate update to the inspect viewport
			newModel, cmd := m.inspectViewPort.Update(msg)
			m.inspectViewPort = newModel
			cmds = append(cmds, cmd)
		} else if m.currentView == LogsView {
			// Delegate update to the log viewport
			newModel, cmd := m.logViewPort.Update(msg)
			m.logViewPort = newModel
			cmds = append(cmds, cmd)
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		// Adjust table heights based on window size
		tableHeight := msg.Height - 10 // Header, footer, status, and some padding
		if tableHeight < 0 {
			tableHeight = 0
		}
		m.containerTable.SetHeight(tableHeight)
		m.imageTable.SetHeight(tableHeight)
		m.volumeTable.SetHeight(tableHeight)
		m.networkTable.SetHeight(tableHeight)

		// Adjust viewport sizes
		if m.currentView == InspectView && !m.inspectReady {
			// Pretty print JSON
			var prettyJSON bytes.Buffer
			if err := json.Indent(&prettyJSON, []byte(m.inspectRawContent), "", "  "); err != nil {
				m.inspectViewPort.SetContent(m.inspectRawContent + "\n\n(Error formatting JSON: " + err.Error() + ")")
			} else {
				m.inspectViewPort.SetContent(prettyJSON.String())
			}
			m.inspectReady = true
		}
		m.inspectViewPort.Width = msg.Width - 4
		m.inspectViewPort.Height = msg.Height - 6
		// No need to set content here, it's done when logs are fetched

	case []controller.Container:
		rows := make([]table.Row, len(msg))
		for i, c := range msg {
			rows[i] = table.Row{c.ID, c.Image, c.Command, c.Created, c.Status, c.Ports, c.Names}
		}
		m.containerTable.SetRows(rows)
		m.showSpinner = false
		m.statusMessage = ""
	case []controller.Image:
		rows := make([]table.Row, len(msg))
		for i, img := range msg {
			rows[i] = table.Row{img.ID, img.Repository, img.Tag, img.Size, img.Created}
		}
		m.imageTable.SetRows(rows)
		m.showSpinner = false
		m.statusMessage = ""
	case []controller.Volume:
		rows := make([]table.Row, len(msg))
		for i, vol := range msg {
			rows[i] = table.Row{vol.Name, vol.Driver, vol.Scope, vol.Mountpoint}
		}
		m.volumeTable.SetRows(rows)
		m.showSpinner = false
		m.statusMessage = ""
	case []controller.Network:
		rows := make([]table.Row, len(msg))
		for i, net := range msg {
			rows[i] = table.Row{net.ID, net.Name, net.Driver, net.Scope}
		}
		m.networkTable.SetRows(rows)
		m.showSpinner = false
		m.statusMessage = ""
	case controller.SystemInfo:
		m.systemInfo = msg
		m.showSpinner = false
		m.statusMessage = ""
	case string: // For logs or inspect output
		m.statusMessage = ""
		m.showSpinner = false
		if m.currentView == InspectView {

			// Pretty print JSON
			var prettyJSON bytes.Buffer
			if err := json.Indent(&prettyJSON, []byte(msg), "", "  "); err != nil {
				m.inspectViewPort.SetContent(msg + "\n\n(Error formatting JSON: " + err.Error() + ")")
			} else {
				m.inspectViewPort.SetContent(prettyJSON.String())
			}
			m.inspectReady = true
			// Manually send a WindowSizeMsg to the inspectViewPort to trigger content rendering
			cmds = append(cmds, func() tea.Msg {
				return tea.WindowSizeMsg{Width: m.width, Height: m.height}
			})
		} else if m.currentView == LogsView {
			m.logViewPort.SetContent(msg)
			m.logViewPort.GotoBottom()
			m.logReady = true
		}
	case error:
		m.err = msg
		m.showSpinner = false
		m.statusMessage = ""
	case statusMsg: // For status messages after actions
		m.statusMessage = string(msg)
		m.showSpinner = false
		cmds = append(cmds, tea.Batch(fetchContainersCmd(), fetchImagesCmd(), fetchVolumesCmd(), fetchNetworksCmd(), fetchSystemInfoCmd()))
	}

	return m, tea.Batch(cmds...)
}
