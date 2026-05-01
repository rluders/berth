package tui

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

// handleKeyMsg dispatches keyboard events to the appropriate handler.
func (m Model) handleKeyMsg(msg tea.KeyMsg) (Model, tea.Cmd) {
	slog.Debug("handleKeyMsg", "key", msg.String())

	// Modal dialog intercepts all keys.
	if m.modal != nil {
		return m.handleModalKey(msg)
	}

	// Help overlay: any key dismisses it.
	if m.showHelp {
		m.showHelp = false
		return m, nil
	}

	// Filter input intercepts typing when active.
	if m.filterActive {
		return m.handleFilterKey(msg)
	}

	// Global keys.
	switch {
	case key.Matches(msg, Keys.Global.Quit):
		return m, tea.Quit
	case key.Matches(msg, Keys.Global.Help):
		m.showHelp = true
		return m, nil
	case key.Matches(msg, Keys.Global.Back):
		switch m.currentView {
		case InspectView, DetailsView:
			m.popView()
			return m, nil
		case LogsView:
			m.stopLogStream()
			m.popView()
			m.logReady = false
			return m, nil
		}
		return m, tea.Quit
	case key.Matches(msg, Keys.Global.Tab1):
		m.currentView = ContainersView
		return m, nil
	case key.Matches(msg, Keys.Global.Tab2):
		m.currentView = ImagesView
		return m, nil
	case key.Matches(msg, Keys.Global.Tab3):
		m.currentView = VolumesView
		return m, nil
	case key.Matches(msg, Keys.Global.Tab4):
		m.currentView = NetworksView
		return m, nil
	case key.Matches(msg, Keys.Global.Tab5):
		m.currentView = SystemView
		return m, nil
	}

	// Per-view keys.
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
	case DetailsView:
		return m.handleDetailsKey(msg)
	}

	return m, nil
}

func (m Model) handleFilterKey(msg tea.KeyMsg) (Model, tea.Cmd) {
	switch {
	case key.Matches(msg, Keys.Filter.Cancel), key.Matches(msg, Keys.Filter.Submit):
		m.filterActive = false
		m.filterInput.Blur()
		m.rebuildFilteredTables()
		return m, nil
	default:
		var cmd tea.Cmd
		m.filterInput, cmd = m.filterInput.Update(msg)
		m.rebuildFilteredTables()
		return m, cmd
	}
}

func (m *Model) rebuildFilteredTables() {
	switch m.currentView {
	case ContainersView:
		rows, metas := m.buildContainerRows()
		m.containerTable.SetRows(rows)
		m.containerVisibleRows = metas
	case ImagesView:
		m.imageTable.SetRows(m.buildImageRows())
	case VolumesView:
		m.volumeTable.SetRows(m.buildVolumeRows())
	}
}

func (m Model) handleContainersKey(msg tea.KeyMsg) (Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd
	m.containerTable, cmd = m.containerTable.Update(msg)
	cmds = append(cmds, cmd)

	switch {
	case key.Matches(msg, Keys.Container.Filter):
		m.filterActive = true
		m.filterInput.Focus()
		return m, tea.Batch(cmds...)
	case key.Matches(msg, Keys.Container.Group):
		m.groupByCompose = !m.groupByCompose
		rows, metas := m.buildContainerRows()
		m.containerTable.SetRows(rows)
		m.containerVisibleRows = metas
		return m, tea.Batch(cmds...)
	}

	if len(m.containerTable.SelectedRow()) == 0 {
		return m, tea.Batch(cmds...)
	}

	// In grouped mode dispatch via metadata to avoid parsing styled strings.
	if m.groupByCompose {
		idx := m.containerTable.Cursor()
		if idx < 0 || idx >= len(m.containerVisibleRows) {
			return m, tea.Batch(cmds...)
		}
		meta := m.containerVisibleRows[idx]

		rebuildRows := func() {
			rows, metas := m.buildContainerRows()
			m.containerTable.SetRows(rows)
			m.containerVisibleRows = metas
		}

		switch meta.kind {
		case rowKindGroup:
			switch {
			case key.Matches(msg, Keys.Container.Expand),
				key.Matches(msg, Keys.Container.Details):
				delete(m.collapsedGroups, meta.groupName)
				rebuildRows()
			case key.Matches(msg, Keys.Container.Collapse):
				m.collapsedGroups[meta.groupName] = true
				rebuildRows()
			}
			return m, tea.Batch(cmds...)

		case rowKindContainer:
			if key.Matches(msg, Keys.Container.Collapse) && meta.groupName != "" {
				m.collapsedGroups[meta.groupName] = true
				rebuildRows()
				return m, tea.Batch(cmds...)
			}
			// Fall through to action dispatch using metadata IDs.
			return m.dispatchContainerAction(msg, meta.containerID, meta.containerName, cmds)
		}

		return m, tea.Batch(cmds...)
	}

	// Ungrouped mode: resolve via display name.
	name := m.containerTable.SelectedRow()[0]
	id := m.resolveContainerID(name)
	if id == "" {
		id = name
	}
	return m.dispatchContainerAction(msg, id, name, cmds)
}

// dispatchContainerAction handles action keys (Details, Start, Stop, etc.) for a resolved container.
func (m Model) dispatchContainerAction(msg tea.KeyMsg, id, name string, cmds []tea.Cmd) (Model, tea.Cmd) {
	switch {
	case key.Matches(msg, Keys.Container.Details):
		m.pushView(DetailsView)
		m.currentDetailsID = id
		m.detailsReady = false
		m.statusMessage = "Loading details..."
		m.showSpinner = true
		cmds = append(cmds, fetchDetailsCmd(id), m.spinner.Tick)
	case key.Matches(msg, Keys.Container.Start):
		m.statusMessage = fmt.Sprintf("docker start %s", name)
		m.showSpinner = true
		cmds = append(cmds, startContainerCmd(id), m.spinner.Tick)
	case key.Matches(msg, Keys.Container.Stop):
		m.statusMessage = fmt.Sprintf("docker stop %s", name)
		m.showSpinner = true
		cmds = append(cmds, stopContainerCmd(id), m.spinner.Tick)
	case key.Matches(msg, Keys.Container.Restart):
		m.statusMessage = fmt.Sprintf("docker restart %s", name)
		m.showSpinner = true
		cmds = append(cmds, restartContainerCmd(id), m.spinner.Tick)
	case key.Matches(msg, Keys.Container.Delete):
		m.modal = NewConfirmModal(
			"Delete Container",
			fmt.Sprintf("Delete container %s?\nThis action cannot be undone.", name),
			tea.Batch(removeContainerCmd(id), m.spinner.Tick),
		)
		m.showSpinner = false
	case key.Matches(msg, Keys.Container.Logs):
		m.stopLogStream()
		m.logLines = nil
		m.logFollowing = true
		m.currentLogContainerID = id
		m.pushView(LogsView)
		m.logReady = true
		ch, cancel, waitCmd := startLogStreamCmd(id)
		m.logCh = ch
		m.logCancel = cancel
		cmds = append(cmds, waitCmd)
	case key.Matches(msg, Keys.Container.Inspect):
		m.pushView(InspectView)
		m.currentInspectID = id
		m.inspectReady = false
		m.statusMessage = fmt.Sprintf("docker inspect %s", name)
		m.showSpinner = true
		cmds = append(cmds, inspectContainerCmd(id), m.spinner.Tick)
	case key.Matches(msg, Keys.Container.Exec):
		cmds = append(cmds, execShellCmd(id))
	}
	return m, tea.Batch(cmds...)
}

// resolveContainerID finds a container's full ID by display name.
func (m Model) resolveContainerID(name string) string {
	for _, c := range m.containers {
		if c.Names == name {
			return c.ID
		}
	}
	return ""
}

func (m Model) handleImagesKey(msg tea.KeyMsg) (Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd
	m.imageTable, cmd = m.imageTable.Update(msg)
	cmds = append(cmds, cmd)

	switch {
	case key.Matches(msg, Keys.Image.Filter):
		m.filterActive = true
		m.filterInput.Focus()
	case key.Matches(msg, Keys.Image.Delete):
		if len(m.imageTable.SelectedRow()) > 0 {
			id := m.imageTable.SelectedRow()[0]
			m.modal = NewConfirmModal(
				"Remove Image",
				fmt.Sprintf("Remove image %s?\nThis action cannot be undone.", id),
				tea.Batch(removeImageCmd(id), m.spinner.Tick),
			)
		}
	case key.Matches(msg, Keys.Image.Prune):
		m.modal = NewConfirmModal(
			"Prune Images",
			"Remove all dangling images?\nThis action cannot be undone.",
			tea.Batch(pruneImagesCmd(), m.spinner.Tick),
		)
	}

	return m, tea.Batch(cmds...)
}

func (m Model) handleVolumesKey(msg tea.KeyMsg) (Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd
	m.volumeTable, cmd = m.volumeTable.Update(msg)
	cmds = append(cmds, cmd)

	switch {
	case key.Matches(msg, Keys.Volume.Filter):
		m.filterActive = true
		m.filterInput.Focus()
	case key.Matches(msg, Keys.Volume.Delete):
		if len(m.volumeTable.SelectedRow()) > 0 {
			name := m.volumeTable.SelectedRow()[0]
			m.modal = NewConfirmModal(
				"Remove Volume",
				fmt.Sprintf("Remove volume %s?\nThis action cannot be undone.", name),
				tea.Batch(removeVolumeCmd(name), m.spinner.Tick),
			)
		}
	}

	return m, tea.Batch(cmds...)
}

func (m Model) handleNetworksKey(msg tea.KeyMsg) (Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd
	m.networkTable, cmd = m.networkTable.Update(msg)
	cmds = append(cmds, cmd)

	if key.Matches(msg, Keys.Network.Inspect) && len(m.networkTable.SelectedRow()) > 0 {
		id := m.networkTable.SelectedRow()[0]
		m.pushView(InspectView)
		m.currentInspectID = id
		m.inspectReady = false
		m.statusMessage = fmt.Sprintf("docker network inspect %s", id)
		m.showSpinner = true
		cmds = append(cmds, inspectNetworkCmd(id), m.spinner.Tick)
	}

	return m, tea.Batch(cmds...)
}

func (m Model) handleSystemKey(msg tea.KeyMsg) (Model, tea.Cmd) {
	switch {
	case key.Matches(msg, Keys.System.BasicCleanup):
		m.modal = NewConfirmModal(
			"Basic Cleanup",
			"Prune stopped containers, unused networks, and dangling images.",
			tea.Batch(
				basicCleanupCmd(),
				func() tea.Msg { return progressMsg{percent: 0.05, label: "Running basic cleanup...", done: false} },
				progressTickCmd(),
			),
		)
	case key.Matches(msg, Keys.System.AdvancedCleanup):
		m.modal = NewConfirmModal(
			"Advanced Cleanup",
			"Prune everything in basic cleanup plus unused volumes.",
			tea.Batch(
				advancedCleanupCmd(),
				func() tea.Msg { return progressMsg{percent: 0.05, label: "Running advanced cleanup...", done: false} },
				progressTickCmd(),
			),
		)
	case key.Matches(msg, Keys.System.TotalCleanup):
		m.modal = NewConfirmModal(
			"Total Cleanup",
			"Remove ALL unused resources including volumes.\nThis action cannot be undone.",
			tea.Batch(
				totalCleanupCmd(),
				func() tea.Msg { return progressMsg{percent: 0.05, label: "Running total cleanup...", done: false} },
				progressTickCmd(),
			),
		)
	}
	return m, nil
}

func (m Model) handleInspectKey(msg tea.KeyMsg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	m.inspectViewPort, cmd = m.inspectViewPort.Update(msg)
	return m, cmd
}

func (m Model) handleLogsKey(msg tea.KeyMsg) (Model, tea.Cmd) {
	switch {
	case key.Matches(msg, Keys.Logs.Pause):
		m.logFollowing = false
		return m, nil
	case key.Matches(msg, Keys.Logs.Follow):
		m.logFollowing = true
		m.logViewPort.GotoBottom()
		return m, nil
	case key.Matches(msg, Keys.Logs.LineNumbers):
		m.showLineNumbers = !m.showLineNumbers
		m.logViewPort.SetContent(buildColorizedLogContent(m.logLines, m.showLineNumbers))
		return m, nil
	}
	var cmd tea.Cmd
	m.logViewPort, cmd = m.logViewPort.Update(msg)
	return m, cmd
}

func (m Model) handleDetailsKey(msg tea.KeyMsg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	m.detailsViewPort, cmd = m.detailsViewPort.Update(msg)
	return m, cmd
}

// stopLogStream cancels and cleans up the log stream goroutine.
func (m *Model) stopLogStream() {
	if m.logCancel != nil {
		m.logCancel()
		m.logCancel = nil
	}
	m.logCh = nil
}

// prettyJSON formats raw JSON content, falling back to raw on error.
func prettyJSON(raw string) string {
	var buf bytes.Buffer
	if err := json.Indent(&buf, []byte(raw), "", "  "); err != nil {
		return raw + "\n\n(Error formatting JSON: " + err.Error() + ")"
	}
	return buf.String()
}
