// Package tui provides the Terminal User Interface for Berth.
package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// colorizeLogLine applies color coding based on log level keywords.
func colorizeLogLine(line string) string {
	th := currentTheme
	lower := strings.ToLower(line)

	// Detect timestamp prefix (common formats: "2024-01-01" or "[12:34:56]")
	// and style it differently.
	switch {
	case strings.Contains(lower, "error") || strings.Contains(lower, "fatal") || strings.Contains(lower, "panic") || strings.Contains(lower, "critical"):
		return th.LogErrorStyle.Render(line)
	case strings.Contains(lower, "warn") || strings.Contains(lower, "warning"):
		return th.LogWarnStyle.Render(line)
	case strings.Contains(lower, "info") || strings.Contains(lower, "notice"):
		return th.LogInfoStyle.Render(line)
	case strings.Contains(lower, "debug") || strings.Contains(lower, "trace"):
		return th.LogDebugStyle.Render(line)
	default:
		return line
	}
}

// buildColorizedLogContent formats all log lines with colors and optional line numbers.
func buildColorizedLogContent(lines []string, showLineNumbers bool) string {
	if len(lines) == 0 {
		return lipgloss.NewStyle().
			Foreground(lipgloss.Color(colorMuted)).
			Render("  Waiting for log output...")
	}

	th := currentTheme
	var sb strings.Builder
	for i, line := range lines {
		if showLineNumbers {
			num := th.LogLineNumStyle.Render(fmt.Sprintf("%4d │ ", i+1))
			sb.WriteString(num + colorizeLogLine(line) + "\n")
		} else {
			sb.WriteString(colorizeLogLine(line) + "\n")
		}
	}
	return sb.String()
}

// View renders the main TUI view.
func (m Model) View() string {
	if m.err != nil {
		return currentTheme.ModalBoxStyle.Render(
			currentTheme.LogErrorStyle.Render("Error: "+m.err.Error()) +
				"\n\nPress q to quit.",
		)
	}

	header := m.renderHeader()
	tabBar := m.renderTabBar()

	// Help overlay takes full screen.
	if m.showHelp {
		return lipgloss.JoinVertical(lipgloss.Top, header, tabBar, m.renderHelp())
	}

	content := m.renderContent()

	// Modal overlaid on top of content.
	if m.modal != nil {
		content = m.renderModal(content)
	}

	footer := m.renderFooter()

	return lipgloss.JoinVertical(
		lipgloss.Top,
		header,
		tabBar,
		lipgloss.NewStyle().
			Width(m.width).
			Height(m.contentHeight()).
			Render(content),
		footer,
	)
}

// renderHeader renders the top bar with logo, view name, and engine badge.
func (m Model) renderHeader() string {
	th := currentTheme

	logo := th.HeaderLogoStyle.Render("berth")

	viewName := ""
	switch m.currentView {
	case InspectView:
		viewName = " › inspect " + m.currentInspectID
	case LogsView:
		viewName = " › logs " + m.currentLogContainerID
	case DetailsView:
		viewName = " › details " + m.currentDetailsID
	}

	left := lipgloss.NewStyle().
		Foreground(lipgloss.Color(colorText)).
		Render(logo + viewName)

	eng := strings.ToUpper(string(m.engineType))
	right := th.HeaderEngStyle.Render("⬡ " + eng)

	gap := m.width - lipgloss.Width(left) - lipgloss.Width(right) - 4
	if gap < 1 {
		gap = 1
	}
	spacer := strings.Repeat(" ", gap)

	return th.HeaderStyle.Width(m.width).Render(left + spacer + right)
}

// renderTabBar renders the tab navigation row.
func (m Model) renderTabBar() string {
	th := currentTheme

	type tabDef struct {
		label string
		view  ViewType
		count int
	}

	tabs := []tabDef{
		{"Containers", ContainersView, len(m.containers)},
		{"Images", ImagesView, len(m.images)},
		{"Volumes", VolumesView, len(m.volumes)},
		{"Networks", NetworksView, 0},
		{"System", SystemView, 0},
	}

	var rendered []string
	for _, tab := range tabs {
		label := tab.label
		if tab.count > 0 {
			label += " " + lipgloss.NewStyle().
				Foreground(lipgloss.Color(colorOverlay)).
				Render(fmt.Sprintf("%d", tab.count))
		}
		if m.currentView == tab.view {
			rendered = append(rendered, th.ActiveTabStyle.Render(label))
		} else {
			rendered = append(rendered, th.InactiveTabStyle.Render(label))
		}
	}

	bar := lipgloss.JoinHorizontal(lipgloss.Top, rendered...)
	// Fill remaining width with tab bar background.
	barW := lipgloss.Width(bar)
	if m.width > barW {
		fill := th.TabBarStyle.Width(m.width - barW).Render("")
		bar = lipgloss.JoinHorizontal(lipgloss.Top, bar, fill)
	}
	return bar
}

// renderContent returns the main body for the current view.
func (m Model) renderContent() string {
	switch m.currentView {
	case ContainersView:
		return m.renderContainerHeader() + "\n" + m.containerVP.View()
	case ImagesView:
		return m.imageTable.View()
	case VolumesView:
		return m.volumeTable.View()
	case NetworksView:
		return m.networkTable.View()
	case SystemView:
		return m.renderSystem()
	case InspectView:
		return m.inspectViewPort.View()
	case LogsView:
		return m.renderLogsView()
	case DetailsView:
		return m.detailsViewPort.View()
	}
	return ""
}

// renderLogsView renders the log viewport with a follow/pause indicator bar.
func (m Model) renderLogsView() string {
	th := currentTheme

	title := m.currentLogContainerID
	if m.currentLogGroupName != "" {
		title = m.currentLogGroupName
	}
	titleBar := lipgloss.NewStyle().
		Padding(0, 1).
		Bold(true).
		Render("Logs: " + title)

	var badge string
	if m.logFollowing {
		badge = th.LogFollowStyle.Render("▶ LIVE")
	} else {
		badge = th.LogPausedStyle.Render("⏸ PAUSED")
	}

	numBadge := ""
	if m.showLineNumbers {
		numBadge = " " + lipgloss.NewStyle().
			Foreground(lipgloss.Color(colorMuted)).
			Render("[line numbers on]")
	}

	indicator := lipgloss.NewStyle().
		Padding(0, 1).
		Render(badge + numBadge)

	return lipgloss.JoinVertical(lipgloss.Left, titleBar, indicator, m.logViewPort.View())
}

// renderSystem renders a dashboard-style system info view.
func (m Model) renderSystem() string {
	th := currentTheme
	si := m.systemInfo

	stat := func(label string, val int, style lipgloss.Style) string {
		return "  " + th.CardTitleStyle.Render(fmt.Sprintf("%-12s", label)) +
			style.Render(fmt.Sprintf("%d", val))
	}

	// ── Containers card ────────────────────────────────────────────────────
	containerCard := th.CardStyle.Render(
		th.SectionStyle.Render("▸ Containers") + "\n" +
			stat("Running", si.Running, th.BadgeRunningStyle) + "\n" +
			stat("Paused", si.Paused, th.BadgePausedStyle) + "\n" +
			stat("Stopped", si.Stopped, th.BadgeStoppedStyle) + "\n" +
			"  " + th.CardTitleStyle.Render(fmt.Sprintf("%-12s", "Total")) +
			th.CardValueStyle.Render(fmt.Sprintf("%d", si.Containers)),
	)

	// ── Resources card ─────────────────────────────────────────────────────
	resourceCard := th.CardStyle.Render(
		th.SectionStyle.Render("▸ Resources") + "\n" +
			"  " + th.CardTitleStyle.Render(fmt.Sprintf("%-12s", "Images")) +
			th.CardValueStyle.Render(fmt.Sprintf("%d", si.Images)) + "\n" +
			"  " + th.CardTitleStyle.Render(fmt.Sprintf("%-12s", "Volumes")) +
			th.CardValueStyle.Render(fmt.Sprintf("%d", si.Volumes)) + "\n" +
			"  " + th.CardTitleStyle.Render(fmt.Sprintf("%-12s", "Networks")) +
			th.CardValueStyle.Render(fmt.Sprintf("%d", si.Networks)) + "\n" +
			"  " + th.CardTitleStyle.Render(fmt.Sprintf("%-12s", "Disk Usage")) +
			th.CardValueStyle.Render(si.DiskUsage),
	)

	// ── Cleanup actions card ───────────────────────────────────────────────
	bBtn := th.ButtonSecondaryStyle.Render("b  Basic cleanup")
	aBtn := th.ButtonSecondaryStyle.Render("a  Advanced cleanup")
	tBtn := th.ButtonDangerStyle.Render("t  Total cleanup")
	actionsCard := th.CardStyle.Render(
		th.SectionStyle.Render("▸ Cleanup Actions") + "\n\n" +
			"  " + bBtn + "\n" +
			"  " + lipgloss.NewStyle().
			Foreground(lipgloss.Color(colorMuted)).
			Render("  Prune stopped containers, unused networks, dangling images") + "\n\n" +
			"  " + aBtn + "\n" +
			"  " + lipgloss.NewStyle().
			Foreground(lipgloss.Color(colorMuted)).
			Render("  + unused volumes") + "\n\n" +
			"  " + tBtn + "\n" +
			"  " + lipgloss.NewStyle().
			Foreground(lipgloss.Color(colorMuted)).
			Render("  Remove ALL unused resources"),
	)

	return lipgloss.JoinVertical(
		lipgloss.Left,
		lipgloss.JoinHorizontal(lipgloss.Top, containerCard, "  ", resourceCard),
		actionsCard,
	)
}

// BuildCommandPreview returns the docker/compose command equivalent for the current selection.
func (m Model) BuildCommandPreview() string {
	if m.currentView != ContainersView {
		return ""
	}
	if len(m.rows) == 0 {
		return " "
	}
	idx := m.containerCursor
	if idx < 0 || idx >= len(m.rows) {
		return " "
	}
	row := m.rows[idx]

	if row.Type == RowTypeGroup {
		project := row.GroupID
		switch m.lastActionKey {
		case "U":
			return fmt.Sprintf("docker compose -p %s up -d --build", project)
		case "R":
			return fmt.Sprintf("docker compose -p %s up -d --force-recreate", project)
		case "d":
			return fmt.Sprintf("docker compose -p %s down", project)
		case "p":
			return fmt.Sprintf("docker compose -p %s pull", project)
		case "b":
			return fmt.Sprintf("docker compose -p %s build", project)
		default:
			return fmt.Sprintf("docker compose -p %s up -d", project)
		}
	}

	name := row.Container.Names
	switch m.lastActionKey {
	case "s":
		return fmt.Sprintf("docker start %s", name)
	case "x":
		return fmt.Sprintf("docker stop %s", name)
	case "r":
		return fmt.Sprintf("docker restart %s", name)
	case "d":
		return fmt.Sprintf("docker rm %s", name)
	case "e":
		return fmt.Sprintf("docker exec -it %s sh", name)
	default:
		return fmt.Sprintf("docker logs -f %s", name)
	}
}

// renderFooter builds the bottom bar with status, hints, and optional overlays.
func (m Model) renderFooter() string {
	th := currentTheme
	var parts []string

	// Command preview (always first, ContainersView only)
	if preview := m.BuildCommandPreview(); preview != "" {
		prefix := th.FooterKeyStyle.Render("$ ")
		parts = append(parts, prefix+th.CommandPreviewStyle.Render(preview))
	}

	// Filter bar
	if m.filterActive {
		filterBar := th.FilterStyle.Render("/ " + m.filterInput.View())
		parts = append(parts, filterBar)
	}

	// Progress bar
	if m.progressVisible {
		parts = append(parts, th.StatusMessageStyle.Render(m.progressLabel))
		parts = append(parts, "  "+m.progressBar.ViewAs(m.progressBar.Percent()))
	}

	// Status / spinner
	if m.statusMessage != "" && !m.progressVisible {
		spinnerStr := ""
		if m.showSpinner {
			spinnerStr = th.SpinnerStyle.Render(m.spinner.View()) + " "
		}
		parts = append(parts, th.StatusMessageStyle.Render(spinnerStr+m.statusMessage))
	}

	// Key hints row
	parts = append(parts, m.renderKeyHints())

	return th.FooterStyle.Width(m.width).Render(strings.Join(parts, "\n"))
}

// renderKeyHints builds a styled key-hint line for the current view.
func (m Model) renderKeyHints() string {
	th := currentTheme
	type hint struct{ k, d string }

	global := []hint{{"?", "help"}, {"q", "quit"}}

	var viewHints []hint
	switch m.currentView {
	case ContainersView:
		idx := m.containerCursor
		if idx >= 0 && idx < len(m.rows) && m.rows[idx].Type == RowTypeGroup {
			viewHints = []hint{
				{"↑/↓", "move"}, {"→/←", "expand/collapse"},
				{"u", "up"}, {"U", "up+build"}, {"R", "recreate"},
				{"d", "down"}, {"p", "pull"}, {"b", "build"},
				{"/", "filter"},
			}
			break
		}
		viewHints = []hint{
			{"↑/↓", "move"}, {"enter", "details"}, {"l", "logs"},
			{"i", "inspect"}, {"s", "start"}, {"x", "stop"},
			{"r", "restart"}, {"d", "delete"}, {"e", "exec"}, {"/", "filter"},
		}
	case ImagesView:
		viewHints = []hint{{"d", "remove"}, {"P", "prune"}, {"/", "filter"}}
	case VolumesView:
		viewHints = []hint{{"d", "remove"}, {"/", "filter"}}
	case NetworksView:
		viewHints = []hint{{"i", "inspect"}}
	case SystemView:
		viewHints = []hint{{"b", "basic"}, {"a", "advanced"}, {"t", "total"}}
	case LogsView:
		viewHints = []hint{{"p", "pause"}, {"f", "follow"}, {"n", "line#"}, {"esc", "back"}}
		global = nil
	case InspectView, DetailsView:
		viewHints = []hint{{"↑/↓", "scroll"}, {"esc", "back"}}
		global = nil
	}

	var segments []string
	for _, h := range viewHints {
		segments = append(segments, FooterHint(h.k, h.d))
	}
	if len(global) > 0 && len(viewHints) > 0 {
		segments = append(segments, th.FooterDescStyle.Render("  •  "))
	}
	for _, h := range global {
		segments = append(segments, FooterHint(h.k, h.d))
	}

	return strings.Join(segments, th.FooterDescStyle.Render("  "))
}

// renderHelp renders a full-screen help overlay using bubbles/help.
func (m Model) renderHelp() string {
	th := currentTheme

	title := th.ModalTitleStyle.Render("Berth — Keyboard Reference (" + m.getViewName() + ")")

	// Use full-help mode (show all bindings).
	hm := m.helpModel
	hm.ShowAll = true
	helpContent := hm.View(m.currentKeyMap())

	hint := lipgloss.NewStyle().
		Foreground(lipgloss.Color(colorMuted)).
		MarginTop(1).
		Render("Press any key to close.")

	inner := lipgloss.JoinVertical(lipgloss.Left, title, "", helpContent, hint)

	boxW := m.width - 8
	if boxW < 50 {
		boxW = 50
	}

	box := th.ModalBoxStyle.Width(boxW).Render(inner)

	// Center horizontally.
	boxW2 := lipgloss.Width(box)
	leftPad := (m.width - boxW2) / 2
	if leftPad < 0 {
		leftPad = 0
	}
	return lipgloss.NewStyle().PaddingLeft(leftPad).Render(box)
}
