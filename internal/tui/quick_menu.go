package tui

import (
	"fmt"
	"strings"

	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

// QuickMenuItem is one selectable action in the quick actions overlay.
type QuickMenuItem struct {
	Label  string
	Key    string
	Action func(m Model) (Model, tea.Cmd)
}

// QuickMenu is a vertical action selection overlay for a container.
type QuickMenu struct {
	Title   string
	Items   []QuickMenuItem
	focused int
}

// FocusNext moves selection down (wraps).
func (q *QuickMenu) FocusNext() {
	q.focused = (q.focused + 1) % len(q.Items)
}

// FocusPrev moves selection up (wraps).
func (q *QuickMenu) FocusPrev() {
	q.focused = (q.focused - 1 + len(q.Items)) % len(q.Items)
}

// Activate calls the focused item's action.
func (q *QuickMenu) Activate(m Model) (Model, tea.Cmd) {
	if q.focused >= 0 && q.focused < len(q.Items) {
		return q.Items[q.focused].Action(m)
	}
	return m, nil
}

// View renders the quick menu box.
func (q QuickMenu) View(width int) string {
	th := currentTheme

	title := th.ModalTitleStyle.Render(q.Title)

	var lines []string
	for i, item := range q.Items {
		keyBadge := lipgloss.NewStyle().
			Foreground(lipgloss.Color(colorMuted)).
			Render(fmt.Sprintf("[%s]", item.Key))

		label := item.Label
		if i == q.focused {
			row := lipgloss.NewStyle().
				Foreground(lipgloss.Color(colorBase)).
				Background(lipgloss.Color(colorMauve)).
				Bold(true).
				Padding(0, 1).
				Render(fmt.Sprintf("%s %s", keyBadge, label))
			lines = append(lines, row)
		} else {
			row := lipgloss.NewStyle().
				Padding(0, 1).
				Render(keyBadge + " " + label)
			lines = append(lines, row)
		}
	}

	hint := lipgloss.NewStyle().
		Foreground(lipgloss.Color(colorMuted)).
		Render("↑/↓ select  enter run  esc close")

	inner := lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		"",
		strings.Join(lines, "\n"),
		"",
		hint,
	)

	boxW := 36
	if width-8 > boxW {
		boxW = width - 8
	}
	if boxW > 48 {
		boxW = 48
	}

	box := th.ModalBoxStyle.Width(boxW).Render(inner)

	boxRenderedW := lipgloss.Width(box)
	leftPad := (width - boxRenderedW) / 2
	if leftPad < 0 {
		leftPad = 0
	}

	return lipgloss.NewStyle().PaddingLeft(leftPad).Render(box)
}

var quickMenuKeys = struct {
	Up     key.Binding
	Down   key.Binding
	Confirm key.Binding
	Cancel  key.Binding
}{
	Up:      key.NewBinding(key.WithKeys("up", "k")),
	Down:    key.NewBinding(key.WithKeys("down", "j")),
	Confirm: key.NewBinding(key.WithKeys("enter")),
	Cancel:  key.NewBinding(key.WithKeys("esc")),
}

// handleQuickMenuKey processes key input when the quick menu is open.
func (m Model) handleQuickMenuKey(msg tea.KeyPressMsg) (Model, tea.Cmd) {
	qm := m.quickMenu
	if qm == nil {
		return m, nil
	}

	switch {
	case key.Matches(msg, quickMenuKeys.Cancel):
		m.quickMenu = nil
		return m, nil
	case key.Matches(msg, quickMenuKeys.Up):
		qm.FocusPrev()
		return m, nil
	case key.Matches(msg, quickMenuKeys.Down):
		qm.FocusNext()
		return m, nil
	case key.Matches(msg, quickMenuKeys.Confirm):
		m.quickMenu = nil
		return qm.Activate(m)
	}

	return m, nil
}

// renderQuickMenu overlays the quick menu centered on a background string.
func (m Model) renderQuickMenu(bg string) string {
	if m.quickMenu == nil {
		return bg
	}
	menuView := m.quickMenu.View(m.width)

	bgLines := strings.Split(bg, "\n")
	menuLines := strings.Split(menuView, "\n")

	bgH := len(bgLines)
	mH := len(menuLines)
	startY := (bgH - mH) / 2
	if startY < 0 {
		startY = 0
	}

	for i, line := range menuLines {
		idx := startY + i
		if idx < len(bgLines) {
			bgLines[idx] = line
		} else {
			bgLines = append(bgLines, line)
		}
	}

	return strings.Join(bgLines, "\n")
}

// NewContainerQuickMenu builds a quick actions menu for a container row.
func NewContainerQuickMenu(id, name string) *QuickMenu {
	return &QuickMenu{
		Title: fmt.Sprintf("Actions  %s", name),
		Items: []QuickMenuItem{
			{
				Label: "Logs",
				Key:   "l",
				Action: func(m Model) (Model, tea.Cmd) {
					m.stopLogStream()
					m.logLines = nil
					m.logFollowing = true
					m.currentLogContainerID = id
					m.currentLogGroupName = ""
					m.pushView(LogsView)
					m.logReady = true
					ch, cancel, waitCmd := startLogStreamCmd(id)
					m.logCh = ch
					m.logCancel = cancel
					return m, waitCmd
				},
			},
			{
				Label: "Exec shell",
				Key:   "e",
				Action: func(m Model) (Model, tea.Cmd) {
					return m, execShellCmd(id)
				},
			},
			{
				Label: "Restart",
				Key:   "r",
				Action: func(m Model) (Model, tea.Cmd) {
					m.statusMessage = fmt.Sprintf("docker restart %s", name)
					m.showSpinner = true
					return m, tea.Batch(restartContainerCmd(id), m.spinner.Tick)
				},
			},
			{
				Label: "Stop",
				Key:   "x",
				Action: func(m Model) (Model, tea.Cmd) {
					m.statusMessage = fmt.Sprintf("docker stop %s", name)
					m.showSpinner = true
					return m, tea.Batch(stopContainerCmd(id), m.spinner.Tick)
				},
			},
			{
				Label: "Delete",
				Key:   "d",
				Action: func(m Model) (Model, tea.Cmd) {
					m.modal = NewConfirmModal(
						"Delete Container",
						fmt.Sprintf("Delete container %s?\nThis action cannot be undone.", name),
						tea.Batch(removeContainerCmd(id), m.spinner.Tick),
					)
					m.showSpinner = false
					return m, nil
				},
			},
		},
	}
}
