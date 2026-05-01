package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// handleMouseMsg dispatches mouse events to the appropriate handler.
func (m Model) handleMouseMsg(msg tea.MouseMsg) (Model, tea.Cmd) {
	// Ignore mouse while modal or filter input is active.
	if m.modal != nil || m.filterActive {
		return m, nil
	}

	switch msg.Button {
	case tea.MouseButtonWheelUp:
		return m.handleScrollUp()
	case tea.MouseButtonWheelDown:
		return m.handleScrollDown()
	case tea.MouseButtonLeft:
		return m.handleLeftClick(msg)
	}

	return m, nil
}

func (m Model) handleScrollUp() (Model, tea.Cmd) {
	switch m.currentView {
	case ContainersView:
		var cmd tea.Cmd
		m.containerTable, cmd = m.containerTable.Update(tea.KeyMsg{Type: tea.KeyUp})
		return m, cmd
	case ImagesView:
		var cmd tea.Cmd
		m.imageTable, cmd = m.imageTable.Update(tea.KeyMsg{Type: tea.KeyUp})
		return m, cmd
	case VolumesView:
		var cmd tea.Cmd
		m.volumeTable, cmd = m.volumeTable.Update(tea.KeyMsg{Type: tea.KeyUp})
		return m, cmd
	case NetworksView:
		var cmd tea.Cmd
		m.networkTable, cmd = m.networkTable.Update(tea.KeyMsg{Type: tea.KeyUp})
		return m, cmd
	case InspectView:
		m.inspectViewPort.LineUp(3)
	case LogsView:
		m.logFollowing = false
		m.logViewPort.LineUp(3)
	case DetailsView:
		m.detailsViewPort.LineUp(3)
	}
	return m, nil
}

func (m Model) handleScrollDown() (Model, tea.Cmd) {
	switch m.currentView {
	case ContainersView:
		var cmd tea.Cmd
		m.containerTable, cmd = m.containerTable.Update(tea.KeyMsg{Type: tea.KeyDown})
		return m, cmd
	case ImagesView:
		var cmd tea.Cmd
		m.imageTable, cmd = m.imageTable.Update(tea.KeyMsg{Type: tea.KeyDown})
		return m, cmd
	case VolumesView:
		var cmd tea.Cmd
		m.volumeTable, cmd = m.volumeTable.Update(tea.KeyMsg{Type: tea.KeyDown})
		return m, cmd
	case NetworksView:
		var cmd tea.Cmd
		m.networkTable, cmd = m.networkTable.Update(tea.KeyMsg{Type: tea.KeyDown})
		return m, cmd
	case InspectView:
		m.inspectViewPort.LineDown(3)
	case LogsView:
		m.logViewPort.LineDown(3)
	case DetailsView:
		m.detailsViewPort.LineDown(3)
	}
	return m, nil
}

func (m Model) handleLeftClick(msg tea.MouseMsg) (Model, tea.Cmd) {
	// Calculate header height to determine if click landed on tab bar.
	headerH := lipgloss.Height(currentTheme.HeaderStyle.Render(m.headerText()))
	tabBarH := 1 // tab bar is 1 line, rendered after header in Task 8

	// Click on header area: check for tab bar clicks.
	if msg.Y >= headerH && msg.Y < headerH+tabBarH {
		return m.handleTabClick(msg.X)
	}

	// Click in content area: handle table row selection.
	contentStartY := headerH + tabBarH
	if msg.Y >= contentStartY {
		return m.handleTableClick(msg.Y-contentStartY, msg.X)
	}

	return m, nil
}

// handleTabClick maps an X coordinate to a tab and switches views.
// Tab positions are computed to match the renderTabBar() layout in view.go.
func (m Model) handleTabClick(x int) (Model, tea.Cmd) {
	tabs := []struct {
		label string
		view  ViewType
	}{
		{"Containers", ContainersView},
		{"Images", ImagesView},
		{"Volumes", VolumesView},
		{"Networks", NetworksView},
		{"System", SystemView},
	}

	// Each tab is rendered as " <label> " (2 padding each side).
	cursor := 0
	for _, tab := range tabs {
		width := len(tab.label) + 4 // 2 padding left + 2 right
		if x >= cursor && x < cursor+width {
			m.currentView = tab.view
			return m, nil
		}
		cursor += width
	}

	return m, nil
}

// handleTableClick maps a Y offset (relative to table start) to a row selection.
func (m Model) handleTableClick(relY, _ int) (Model, tea.Cmd) {
	// Row 0 = table header; data rows start at relY == 1.
	if relY <= 0 {
		return m, nil
	}
	rowIndex := relY - 1 // 0-based data row index

	switch m.currentView {
	case ContainersView:
		rows := m.containerTable.Rows()
		if rowIndex < len(rows) {
			for i := 0; i < rowIndex; i++ {
				var cmd tea.Cmd
				m.containerTable, cmd = m.containerTable.Update(tea.KeyMsg{Type: tea.KeyDown})
				_ = cmd
			}
			// Reset to top first, then navigate to target row.
			m.containerTable.GotoTop()
			for i := 0; i < rowIndex; i++ {
				m.containerTable, _ = m.containerTable.Update(tea.KeyMsg{Type: tea.KeyDown})
			}
		}
	case ImagesView:
		rows := m.imageTable.Rows()
		if rowIndex < len(rows) {
			m.imageTable.GotoTop()
			for i := 0; i < rowIndex; i++ {
				m.imageTable, _ = m.imageTable.Update(tea.KeyMsg{Type: tea.KeyDown})
			}
		}
	case VolumesView:
		rows := m.volumeTable.Rows()
		if rowIndex < len(rows) {
			m.volumeTable.GotoTop()
			for i := 0; i < rowIndex; i++ {
				m.volumeTable, _ = m.volumeTable.Update(tea.KeyMsg{Type: tea.KeyDown})
			}
		}
	case NetworksView:
		rows := m.networkTable.Rows()
		if rowIndex < len(rows) {
			m.networkTable.GotoTop()
			for i := 0; i < rowIndex; i++ {
				m.networkTable, _ = m.networkTable.Update(tea.KeyMsg{Type: tea.KeyDown})
			}
		}
	}

	return m, nil
}
