package tui

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
)

var mainTabs = []ViewType{ContainersView, ImagesView, VolumesView, NetworksView, SystemView}

// handleKeyMsg dispatches keyboard events to the appropriate handler.
func (m Model) handleKeyMsg(msg tea.KeyPressMsg) (Model, tea.Cmd) {
	slog.Debug("handleKeyMsg", "key", msg.String())

	// Quick menu intercepts all keys when open.
	if m.quickMenu != nil {
		return m.handleQuickMenuKey(msg)
	}

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
			m.currentLogGroupName = ""
			return m, nil
		}
		return m, tea.Quit
	case key.Matches(msg, Keys.Global.Tab1):
		m.leaveLogView()
		m.currentView = ContainersView
		return m, nil
	case key.Matches(msg, Keys.Global.Tab2):
		m.leaveLogView()
		m.currentView = ImagesView
		return m, nil
	case key.Matches(msg, Keys.Global.Tab3):
		m.leaveLogView()
		m.currentView = VolumesView
		return m, nil
	case key.Matches(msg, Keys.Global.Tab4):
		m.leaveLogView()
		m.currentView = NetworksView
		return m, nil
	case key.Matches(msg, Keys.Global.Tab5):
		m.leaveLogView()
		m.currentView = SystemView
		return m, nil
	case key.Matches(msg, Keys.Global.TabNext):
		m.leaveLogView()
		return m.cycleTab(+1), nil
	case key.Matches(msg, Keys.Global.TabPrev):
		m.leaveLogView()
		return m.cycleTab(-1), nil
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

func (m Model) handleFilterKey(msg tea.KeyPressMsg) (Model, tea.Cmd) {
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
		m.recomputeRows()
	case ImagesView:
		m.imageTable.SetRows(m.buildImageRows())
	case VolumesView:
		m.volumeTable.SetRows(m.buildVolumeRows())
	}
}

func (m Model) handleContainersKey(msg tea.KeyPressMsg) (Model, tea.Cmd) {
	var cmds []tea.Cmd

	// Movement keys — handled manually since we no longer use bubbles/table.
	switch msg.String() {
	case "up", "k":
		m.moveContainerCursor(-1)
		m.lastActionKey = ""
		return m, nil
	case "down", "j":
		m.moveContainerCursor(+1)
		m.lastActionKey = ""
		return m, nil
	case "home":
		if len(m.rows) > 0 {
			m.containerCursor = 0
			m.syncContainerViewport()
		}
		m.lastActionKey = ""
		return m, nil
	case "G", "end":
		if len(m.rows) > 0 {
			m.containerCursor = len(m.rows) - 1
			m.syncContainerViewport()
		}
		m.lastActionKey = ""
		return m, nil
	case "pgup", "ctrl+b":
		step := max(1, m.containerVP.Height())
		m.moveContainerCursor(-step)
		m.lastActionKey = ""
		return m, nil
	case "pgdown", "ctrl+f":
		step := max(1, m.containerVP.Height())
		m.moveContainerCursor(+step)
		m.lastActionKey = ""
		return m, nil
	}

	// Track last action key for command preview; movement keys reset to default.
	switch msg.String() {
	case "s", "x", "r", "d", "l", "i", "e", "u", "U", "R", "p", "b":
		m.lastActionKey = msg.String()
	default:
		m.lastActionKey = ""
	}

	if key.Matches(msg, Keys.Container.Filter) {
		m.filterActive = true
		m.filterInput.Focus()
		return m, tea.Batch(cmds...)
	}

	if len(m.rows) == 0 {
		return m, tea.Batch(cmds...)
	}

	idx := m.containerCursor
	if idx < 0 || idx >= len(m.rows) {
		return m, tea.Batch(cmds...)
	}
	row := m.rows[idx]

	switch row.Type {
	case RowTypeGroup:
		switch {
		case key.Matches(msg, Keys.Container.Details):
			m.collapsedGroups[row.GroupID] = !m.collapsedGroups[row.GroupID]
			m.recomputeRows()
			saveState(persistedState{CollapsedGroups: m.collapsedGroups})
		case key.Matches(msg, Keys.Container.Expand):
			delete(m.collapsedGroups, row.GroupID)
			m.recomputeRows()
			saveState(persistedState{CollapsedGroups: m.collapsedGroups})
		case key.Matches(msg, Keys.Container.Collapse):
			m.collapsedGroups[row.GroupID] = true
			m.recomputeRows()
			saveState(persistedState{CollapsedGroups: m.collapsedGroups})
		case key.Matches(msg, Keys.Container.Start):
			m.statusMessage = fmt.Sprintf("docker start [%s]", row.GroupID)
			m.showSpinner = true
			cmds = append(cmds, startGroupContainersCmd(row.Containers), m.spinner.Tick)
		case key.Matches(msg, Keys.Container.Stop):
			m.statusMessage = fmt.Sprintf("docker stop [%s]", row.GroupID)
			m.showSpinner = true
			cmds = append(cmds, stopGroupContainersCmd(row.Containers), m.spinner.Tick)
		case key.Matches(msg, Keys.Container.Restart):
			m.statusMessage = fmt.Sprintf("docker restart [%s]", row.GroupID)
			m.showSpinner = true
			cmds = append(cmds, restartGroupContainersCmd(row.Containers), m.spinner.Tick)
		case key.Matches(msg, Keys.Container.Delete):
			var removeCmds []tea.Cmd
			for _, c := range row.Containers {
				removeCmds = append(removeCmds, removeContainerCmd(c.ID))
			}
			removeCmds = append(removeCmds, m.spinner.Tick)
			m.modal = NewConfirmModal(
				"Delete Group",
				fmt.Sprintf("Delete all containers in %s?\nThis action cannot be undone.", row.GroupID),
				tea.Batch(removeCmds...),
			)
			m.showSpinner = false
		case key.Matches(msg, Keys.Container.Logs):
			m.stopLogStream()
			m.logLines = nil
			m.logFollowing = true
			m.currentLogGroupName = row.GroupID
			m.currentLogContainerID = ""
			group := findGroupContainers(m.containers, row.GroupID)
			m.pushView(LogsView)
			m.logReady = true
			ch, cancel, waitCmd := startGroupLogStreamCmd(group)
			m.logCh = ch
			m.logCancel = cancel
			cmds = append(cmds, waitCmd)
		case key.Matches(msg, Keys.Container.QuickActions):
			m.statusMessage = "Group: use s/x/r/d to start/stop/restart/delete all containers"
		default:
			workDir := m.composeWorkDir(row.GroupID)
			return m.dispatchComposeAction(msg, row.GroupID, workDir, cmds)
		}

	case RowTypeContainer:
		if key.Matches(msg, Keys.Container.Collapse) && row.GroupID != "" {
			m.collapsedGroups[row.GroupID] = true
			m.recomputeRows()
			return m, tea.Batch(cmds...)
		}
		if key.Matches(msg, Keys.Container.QuickActions) {
			m.quickMenu = NewContainerQuickMenu(row.Container.ID, row.Container.Names)
			return m, nil
		}
		return m.dispatchContainerAction(msg, row.Container.ID, row.Container.Names, cmds)
	}

	return m, tea.Batch(cmds...)
}

// dispatchContainerAction handles action keys (Details, Start, Stop, etc.) for a resolved container.
func (m Model) dispatchContainerAction(msg tea.KeyPressMsg, id, name string, cmds []tea.Cmd) (Model, tea.Cmd) {
	switch {
	case key.Matches(msg, Keys.Container.Details):
		m.pushView(DetailsView)
		m.currentDetailsID = id
		m.detailsReady = false
		m.statusMessage = fmt.Sprintf("Loading details %s...", name)
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

// composeWorkDir returns the working directory for a compose project by checking container labels.
func (m Model) composeWorkDir(project string) string {
	for _, c := range m.containers {
		if c.Labels["com.docker.compose.project"] == project {
			if wd := c.Labels["com.docker.compose.project.working_dir"]; wd != "" {
				return wd
			}
		}
	}
	return ""
}

// startComposeOp cancels any running compose op, creates a fresh context, and returns the cmd.
func (m *Model) startComposeOp(project, workDir string, fn func(context.Context, string, string) tea.Cmd) tea.Cmd {
	if m.composeCancel != nil {
		m.composeCancel()
	}
	ctx, cancel := context.WithCancel(context.Background())
	m.composeCancel = cancel
	m.showSpinner = true
	m.composeOutput = nil
	return fn(ctx, project, workDir)
}

// dispatchComposeAction handles compose project-level action keys when a group row is selected.
func (m Model) dispatchComposeAction(msg tea.KeyPressMsg, project, workDir string, cmds []tea.Cmd) (Model, tea.Cmd) {
	switch {
	case key.Matches(msg, Keys.Compose.Up):
		m.statusMessage = fmt.Sprintf("docker compose up -d  [%s]", project)
		cmds = append(cmds, m.startComposeOp(project, workDir, composeUpCmd), m.spinner.Tick)
	case key.Matches(msg, Keys.Compose.UpBuild):
		m.statusMessage = fmt.Sprintf("docker compose up -d --build  [%s]", project)
		cmds = append(cmds, m.startComposeOp(project, workDir, composeUpBuildCmd), m.spinner.Tick)
	case key.Matches(msg, Keys.Compose.Recreate):
		ctx, cancel := context.WithCancel(context.Background())
		if m.composeCancel != nil {
			m.composeCancel()
		}
		m.composeCancel = cancel
		m.composeOutput = nil
		m.modal = NewConfirmModal(
			"Force Recreate",
			fmt.Sprintf("docker compose up -d --force-recreate\nProject: %s", project),
			tea.Batch(composeRecreateCmd(ctx, project, workDir), m.spinner.Tick),
		)
	case key.Matches(msg, Keys.Compose.Down):
		ctx, cancel := context.WithCancel(context.Background())
		if m.composeCancel != nil {
			m.composeCancel()
		}
		m.composeCancel = cancel
		m.composeOutput = nil
		m.modal = NewConfirmModal(
			"Compose Down",
			fmt.Sprintf("docker compose down\nProject: %s", project),
			tea.Batch(composeDownCmd(ctx, project, workDir), m.spinner.Tick),
		)
	case key.Matches(msg, Keys.Compose.Pull):
		m.statusMessage = fmt.Sprintf("docker compose pull  [%s]", project)
		cmds = append(cmds, m.startComposeOp(project, workDir, composePullCmd), m.spinner.Tick)
	case key.Matches(msg, Keys.Compose.Build):
		m.statusMessage = fmt.Sprintf("docker compose build  [%s]", project)
		cmds = append(cmds, m.startComposeOp(project, workDir, composeBuildCmd), m.spinner.Tick)
	}
	return m, tea.Batch(cmds...)
}

func (m Model) handleImagesKey(msg tea.KeyPressMsg) (Model, tea.Cmd) {
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

func (m Model) handleVolumesKey(msg tea.KeyPressMsg) (Model, tea.Cmd) {
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

func (m Model) handleNetworksKey(msg tea.KeyPressMsg) (Model, tea.Cmd) {
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

func (m Model) handleSystemKey(msg tea.KeyPressMsg) (Model, tea.Cmd) {
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

func (m Model) handleInspectKey(msg tea.KeyPressMsg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	m.inspectViewPort, cmd = m.inspectViewPort.Update(msg)
	return m, cmd
}

func (m Model) handleLogsKey(msg tea.KeyPressMsg) (Model, tea.Cmd) {
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

func (m Model) handleDetailsKey(msg tea.KeyPressMsg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	m.detailsViewPort, cmd = m.detailsViewPort.Update(msg)
	return m, cmd
}

func (m Model) cycleTab(delta int) Model {
	for i, v := range mainTabs {
		if v == m.currentView {
			m.currentView = mainTabs[(i+delta+len(mainTabs))%len(mainTabs)]
			return m
		}
	}
	return m
}

// leaveLogView stops the log stream and resets log state when navigating away via tab switch.
func (m *Model) leaveLogView() {
	if m.currentView != LogsView {
		return
	}
	m.stopLogStream()
	m.logReady = false
	m.currentLogGroupName = ""
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
