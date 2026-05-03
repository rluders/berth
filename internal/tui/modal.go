package tui

import (
	"strings"

	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

// ButtonKind controls button visual style.
type ButtonKind int

const (
	ButtonKindPrimary ButtonKind = iota
	ButtonKindDanger
	ButtonKindSecondary
)

// ModalButton is one button in a modal dialog.
type ModalButton struct {
	Label string
	Kind  ButtonKind
	// Cmd is the tea.Cmd to dispatch when this button is activated.
	Cmd tea.Cmd
}

// Modal is a centered dialog with a title, body text, and focusable buttons.
type Modal struct {
	Title   string
	Body    string
	Buttons []ModalButton
	focused int
}

// NewConfirmModal creates a two-button danger confirm dialog.
func NewConfirmModal(title, body string, confirmCmd tea.Cmd) *Modal {
	return &Modal{
		Title: title,
		Body:  body,
		Buttons: []ModalButton{
			{Label: "Confirm", Kind: ButtonKindDanger, Cmd: confirmCmd},
			{Label: "Cancel", Kind: ButtonKindSecondary, Cmd: nil},
		},
		focused: 1, // default focus on Cancel (safer)
	}
}

// FocusNext moves focus to the next button (wraps).
func (m *Modal) FocusNext() {
	m.focused = (m.focused + 1) % len(m.Buttons)
}

// FocusPrev moves focus to the previous button (wraps).
func (m *Modal) FocusPrev() {
	m.focused = (m.focused - 1 + len(m.Buttons)) % len(m.Buttons)
}

// Activate returns the Cmd for the focused button (nil = cancel).
func (m *Modal) Activate() tea.Cmd {
	if m.focused >= 0 && m.focused < len(m.Buttons) {
		return m.Buttons[m.focused].Cmd
	}
	return nil
}

// ActivateAt returns the Cmd for the button at index i (nil = cancel).
func (m *Modal) ActivateAt(i int) tea.Cmd {
	if i >= 0 && i < len(m.Buttons) {
		return m.Buttons[i].Cmd
	}
	return nil
}

// View renders the modal box using the current theme.
func (m Modal) View(width int) string {
	th := currentTheme

	title := th.ModalTitleStyle.Render(m.Title)
	body := th.ModalBodyStyle.Render(m.Body)

	// Render buttons.
	var btns []string
	for i, b := range m.Buttons {
		var s lipgloss.Style
		if i == m.focused {
			s = th.ButtonFocusedStyle
		} else {
			switch b.Kind {
			case ButtonKindDanger:
				s = th.ButtonDangerStyle
			case ButtonKindPrimary:
				s = th.ButtonPrimaryStyle
			default:
				s = th.ButtonSecondaryStyle
			}
		}
		btns = append(btns, s.Render(b.Label))
	}
	buttonRow := lipgloss.JoinHorizontal(lipgloss.Left, btns...)
	hint := lipgloss.NewStyle().
		Foreground(lipgloss.Color(colorMuted)).
		Render("tab/← → to switch  enter to confirm  esc to cancel")

	inner := lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		body,
		"",
		buttonRow,
		hint,
	)

	// Constrain box width.
	boxW := width - 8
	if boxW < 40 {
		boxW = 40
	}

	box := th.ModalBoxStyle.Width(boxW).Render(inner)

	// Center horizontally.
	boxRenderedW := lipgloss.Width(box)
	leftPad := (width - boxRenderedW) / 2
	if leftPad < 0 {
		leftPad = 0
	}

	return lipgloss.NewStyle().
		PaddingLeft(leftPad).
		Render(box)
}

// modalKeys are the key bindings active while a modal is open.
var modalKeys = struct {
	FocusNext key.Binding
	FocusPrev key.Binding
	Confirm   key.Binding
	Cancel    key.Binding
	Left      key.Binding
	Right     key.Binding
}{
	FocusNext: key.NewBinding(key.WithKeys("tab")),
	FocusPrev: key.NewBinding(key.WithKeys("shift+tab")),
	Confirm:   key.NewBinding(key.WithKeys("enter")),
	Cancel:    key.NewBinding(key.WithKeys("esc")),
	Left:      key.NewBinding(key.WithKeys("left", "h")),
	Right:     key.NewBinding(key.WithKeys("right", "l")),
}

// handleModalKey processes key input when a modal is active.
func (m Model) handleModalKey(msg tea.KeyPressMsg) (Model, tea.Cmd) {
	modal := m.modal
	if modal == nil {
		return m, nil
	}

	switch {
	case key.Matches(msg, modalKeys.Cancel):
		m.modal = nil
		m.statusMessage = "Cancelled."
		return m, nil
	case key.Matches(msg, modalKeys.Confirm):
		cmd := modal.Activate()
		m.modal = nil
		if cmd == nil {
			m.statusMessage = "Cancelled."
			return m, nil
		}
		m.showSpinner = true
		return m, cmd
	case key.Matches(msg, modalKeys.FocusNext), key.Matches(msg, modalKeys.Right):
		modal.FocusNext()
		return m, nil
	case key.Matches(msg, modalKeys.FocusPrev), key.Matches(msg, modalKeys.Left):
		modal.FocusPrev()
		return m, nil
	}

	return m, nil
}

// renderModal overlays the modal centered on a background string.
func (m Model) renderModal(bg string) string {
	if m.modal == nil {
		return bg
	}
	modalView := m.modal.View(m.width)

	bgLines := strings.Split(bg, "\n")
	modalLines := strings.Split(modalView, "\n")

	// Center the modal vertically.
	bgH := len(bgLines)
	mH := len(modalLines)
	startY := (bgH - mH) / 2
	if startY < 0 {
		startY = 0
	}

	for i, line := range modalLines {
		idx := startY + i
		if idx < len(bgLines) {
			bgLines[idx] = line
		} else {
			bgLines = append(bgLines, line)
		}
	}

	return strings.Join(bgLines, "\n")
}
