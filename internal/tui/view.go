package tui

import (
	"fmt"
	"github.com/charmbracelet/lipgloss"
	"strings"
)

func (m Model) View() string {
	if m.err != nil {
		return fmt.Sprintf("Error: %v", m.err)
	}

	header := currentTheme.HeaderStyle.Render(fmt.Sprintf("Berth - %s - %s Engine", m.getViewName(), strings.ToUpper(string(m.engineType))))

	content := ""
	switch m.currentView {
	case ContainersView:
		content = m.containerTable.View()
	case ImagesView:
		content = m.imageTable.View()
	case VolumesView:
		content = m.volumeTable.View()
	case NetworksView:
		content = m.networkTable.View()
	case SystemView:
		content = fmt.Sprintf("  Containers: %d (Running: %d, Paused: %d, Stopped: %d)\n", m.systemInfo.Containers, m.systemInfo.Running, m.systemInfo.Paused, m.systemInfo.Stopped) +
			fmt.Sprintf("  Images: %d\n", m.systemInfo.Images) +
			fmt.Sprintf("  Volumes: %d\n", m.systemInfo.Volumes) +
			fmt.Sprintf("  Networks: %d\n", m.systemInfo.Networks) +
			fmt.Sprintf("  Disk Usage: %s\n", m.systemInfo.DiskUsage)
	case InspectView:
		content = m.inspectViewPort.View()
	case LogsView:
		content = m.logViewPort.View()
	}

	footerContent := m.getFooterHelp()
	if m.statusMessage != "" {
		spinnerStr := ""
		if m.showSpinner {
			spinnerStr = m.spinner.View() + " "
		}
		footerContent = currentTheme.StatusMessageStyle.Render(spinnerStr+m.statusMessage) + "\n" + footerContent
	}

	return currentTheme.AppStyle.Render(
		lipgloss.JoinVertical(
			lipgloss.Top,
			header,
			lipgloss.NewStyle().Height(m.height-lipgloss.Height(header)-currentTheme.FooterStyle.GetVerticalPadding()-currentTheme.HeaderStyle.GetVerticalPadding()-currentTheme.AppStyle.GetVerticalPadding()*2).Render(content),
			currentTheme.FooterStyle.Render(footerContent),
		),
	)
}
