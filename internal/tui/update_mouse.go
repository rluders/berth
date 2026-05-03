package tui

import (
	"fmt"

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
		m.containerVP.ScrollUp(3)
		return m, nil
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
		m.inspectViewPort.ScrollUp(3)
	case LogsView:
		m.logFollowing = false
		m.logViewPort.ScrollUp(3)
	case DetailsView:
		m.detailsViewPort.ScrollUp(3)
	}
	return m, nil
}

func (m Model) handleScrollDown() (Model, tea.Cmd) {
	switch m.currentView {
	case ContainersView:
		m.containerVP.ScrollDown(3)
		return m, nil
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
		m.inspectViewPort.ScrollDown(3)
	case LogsView:
		m.logViewPort.ScrollDown(3)
	case DetailsView:
		m.detailsViewPort.ScrollDown(3)
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
// Replicates renderTabBar() label+style logic to get accurate widths via lipgloss.Width().
func (m Model) handleTabClick(x int) (Model, tea.Cmd) {
	th := currentTheme
	tabs := []struct {
		label string
		count int
		view  ViewType
	}{
		{"Containers", len(m.containers), ContainersView},
		{"Images", len(m.images), ImagesView},
		{"Volumes", len(m.volumes), VolumesView},
		{"Networks", 0, NetworksView},
		{"System", 0, SystemView},
	}

	cursor := 0
	for _, tab := range tabs {
		label := tab.label
		if tab.count > 0 {
			label += " " + lipgloss.NewStyle().
				Foreground(lipgloss.Color(colorOverlay)).
				Render(fmt.Sprintf("%d", tab.count))
		}
		var rendered string
		if m.currentView == tab.view {
			rendered = th.ActiveTabStyle.Render(label)
		} else {
			rendered = th.InactiveTabStyle.Render(label)
		}
		w := lipgloss.Width(rendered)
		if x >= cursor && x < cursor+w {
			m.currentView = tab.view
			return m, nil
		}
		cursor += w
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
		// relY==0 is the header line; data rows start at relY==1.
		dataIdx := rowIndex + m.containerVP.YOffset
		if dataIdx >= 0 && dataIdx < len(m.rows) {
			m.containerCursor = dataIdx
			m.syncContainerViewport()
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
