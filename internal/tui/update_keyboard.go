package tui

import (
	"bytes"
	"encoding/json"
	"log/slog"

	tea "github.com/charmbracelet/bubbletea"
)

// handleKeyMsg dispatches keyboard events to the appropriate handler.
func (m Model) handleKeyMsg(msg tea.KeyMsg) (Model, tea.Cmd) {
	slog.Debug("handleKeyMsg", "key", msg.String())

	// Global keys
	switch msg.String() {
	case "ctrl+c":
		return m, tea.Quit
	case "q", "esc":
		if m.currentView == InspectView || m.currentView == LogsView {
			m.popView()
			m.logReady = false
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

	// Per-view keys
	switch m.currentView {
	case ContainersView:
		return m.handleContainersKey(msg)
	case ImagesView:
		return m.handleImagesKey(msg)
	case VolumesView:
		return m.handleVolumesKey(msg)
	case NetworksView:
		return m.handleNetworksKey(msg)
	case SystemView:
		return m.handleSystemKey(msg)
	case InspectView:
		return m.handleInspectKey(msg)
	case LogsView:
		return m.handleLogsKey(msg)
	}

	return m, nil
}

func (m Model) handleContainersKey(msg tea.KeyMsg) (Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd
	m.containerTable, cmd = m.containerTable.Update(msg)
	cmds = append(cmds, cmd)

	if len(m.containerTable.SelectedRow()) == 0 {
		return m, tea.Batch(cmds...)
	}
	id := m.containerTable.SelectedRow()[0]

	switch msg.String() {
	case "s":
		m.statusMessage = "Starting container " + id + "..."
		m.showSpinner = true
		cmds = append(cmds, startContainerCmd(id), m.spinner.Tick)
	case "x":
		m.statusMessage = "Stopping container " + id + "..."
		m.showSpinner = true
		cmds = append(cmds, stopContainerCmd(id), m.spinner.Tick)
	case "d":
		m.statusMessage = "Removing container " + id + "..."
		m.showSpinner = true
		cmds = append(cmds, removeContainerCmd(id), m.spinner.Tick)
	case "l":
		m.pushView(LogsView)
		m.logReady = false
		m.currentLogContainerID = id
		m.statusMessage = "Fetching logs for " + id + "..."
		m.showSpinner = true
		cmds = append(cmds, getLogsCmd(id), m.spinner.Tick)
	case "i":
		m.pushView(InspectView)
		m.currentInspectID = id
		m.inspectReady = false
		m.statusMessage = "Inspecting container " + id + "..."
		m.showSpinner = true
		cmds = append(cmds, inspectContainerCmd(id), m.spinner.Tick)
	}

	return m, tea.Batch(cmds...)
}

func (m Model) handleImagesKey(msg tea.KeyMsg) (Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd
	m.imageTable, cmd = m.imageTable.Update(msg)
	cmds = append(cmds, cmd)

	if msg.String() == "d" && len(m.imageTable.SelectedRow()) > 0 {
		id := m.imageTable.SelectedRow()[0]
		m.statusMessage = "Removing image " + id + "..."
		m.showSpinner = true
		cmds = append(cmds, removeImageCmd(id), m.spinner.Tick)
	}

	return m, tea.Batch(cmds...)
}

func (m Model) handleVolumesKey(msg tea.KeyMsg) (Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd
	m.volumeTable, cmd = m.volumeTable.Update(msg)
	cmds = append(cmds, cmd)

	if msg.String() == "d" && len(m.volumeTable.SelectedRow()) > 0 {
		name := m.volumeTable.SelectedRow()[0]
		m.statusMessage = "Removing volume " + name + "..."
		m.showSpinner = true
		cmds = append(cmds, removeVolumeCmd(name), m.spinner.Tick)
	}

	return m, tea.Batch(cmds...)
}

func (m Model) handleNetworksKey(msg tea.KeyMsg) (Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd
	m.networkTable, cmd = m.networkTable.Update(msg)
	cmds = append(cmds, cmd)

	if msg.String() == "i" && len(m.networkTable.SelectedRow()) > 0 {
		id := m.networkTable.SelectedRow()[0]
		m.pushView(InspectView)
		m.currentInspectID = id
		m.inspectReady = false
		m.statusMessage = "Inspecting network " + id + "..."
		m.showSpinner = true
		cmds = append(cmds, inspectNetworkCmd(id), m.spinner.Tick)
	}

	return m, tea.Batch(cmds...)
}

func (m Model) handleSystemKey(msg tea.KeyMsg) (Model, tea.Cmd) {
	switch msg.String() {
	case "b":
		m.statusMessage = "Performing basic cleanup..."
		m.showSpinner = true
		return m, tea.Batch(basicCleanupCmd(), m.spinner.Tick)
	case "a":
		m.statusMessage = "Performing advanced cleanup..."
		m.showSpinner = true
		return m, tea.Batch(advancedCleanupCmd(), m.spinner.Tick)
	case "t":
		m.statusMessage = "Performing total cleanup..."
		m.showSpinner = true
		return m, tea.Batch(totalCleanupCmd(), m.spinner.Tick)
	}
	return m, nil
}

func (m Model) handleInspectKey(msg tea.KeyMsg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	m.inspectViewPort, cmd = m.inspectViewPort.Update(msg)
	return m, cmd
}

func (m Model) handleLogsKey(msg tea.KeyMsg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	m.logViewPort, cmd = m.logViewPort.Update(msg)
	return m, cmd
}

// prettyJSON formats raw JSON content, falling back to raw on error.
func prettyJSON(raw string) string {
	var buf bytes.Buffer
	if err := json.Indent(&buf, []byte(raw), "", "  "); err != nil {
		return raw + "\n\n(Error formatting JSON: " + err.Error() + ")"
	}
	return buf.String()
}
