package tui

import (
	"log/slog"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/rluders/berth/internal/controller"
)

// Update handles incoming messages and updates the model.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		cmds = append(cmds, cmd)

	case tea.KeyMsg:
		newM, cmd := m.handleKeyMsg(msg)
		return newM, cmd

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		contentH := m.contentHeight()
		m.containerTable.SetHeight(contentH)
		m.imageTable.SetHeight(contentH)
		m.volumeTable.SetHeight(contentH)
		m.networkTable.SetHeight(contentH)

		viewW := msg.Width - currentTheme.AppStyle.GetHorizontalFrameSize() - 4
		if viewW < 0 {
			viewW = 0
		}
		m.inspectViewPort.Width = viewW
		m.inspectViewPort.Height = contentH
		m.logViewPort.Width = viewW
		m.logViewPort.Height = contentH

		// If inspect content arrived before window size, render it now.
		if m.currentView == InspectView && !m.inspectReady && m.inspectRawContent != "" {
			m.inspectViewPort.SetContent(prettyJSON(m.inspectRawContent))
			m.inspectReady = true
		}

	case containerListMsg:
		slog.Debug("containerListMsg received", "count", len(msg))
		rows := make([]table.Row, len(msg))
		for i, c := range msg {
			rows[i] = table.Row{c.ID, c.Image, c.Command, c.Created, c.Status, c.Ports, c.Names}
		}
		m.containerTable.SetRows(rows)
		m.showSpinner = false
		m.statusMessage = ""

	case imageListMsg:
		slog.Debug("imageListMsg received", "count", len(msg))
		rows := make([]table.Row, len(msg))
		for i, img := range msg {
			rows[i] = table.Row{img.ID, img.Repository, img.Tag, img.Size, img.Created}
		}
		m.imageTable.SetRows(rows)
		m.showSpinner = false
		m.statusMessage = ""

	case volumeListMsg:
		slog.Debug("volumeListMsg received", "count", len(msg))
		rows := make([]table.Row, len(msg))
		for i, vol := range msg {
			rows[i] = table.Row{vol.Name, vol.Driver, vol.Scope, vol.Mountpoint}
		}
		m.volumeTable.SetRows(rows)
		m.showSpinner = false
		m.statusMessage = ""

	case networkListMsg:
		slog.Debug("networkListMsg received", "count", len(msg))
		rows := make([]table.Row, len(msg))
		for i, net := range msg {
			rows[i] = table.Row{net.ID, net.Name, net.Driver, net.Scope}
		}
		m.networkTable.SetRows(rows)
		m.showSpinner = false
		m.statusMessage = ""

	case systemInfoMsg:
		slog.Debug("systemInfoMsg received")
		m.systemInfo = controller.SystemInfo(msg)
		m.showSpinner = false
		m.statusMessage = ""

	case inspectMsg:
		slog.Debug("inspectMsg received")
		m.showSpinner = false
		m.statusMessage = ""
		m.inspectRawContent = string(msg)
		if m.width > 0 {
			m.inspectViewPort.SetContent(prettyJSON(string(msg)))
			m.inspectReady = true
		}
		// Trigger viewport to re-measure in case window size already arrived.
		cmds = append(cmds, func() tea.Msg {
			return tea.WindowSizeMsg{Width: m.width, Height: m.height}
		})

	case logsMsg:
		slog.Debug("logsMsg received")
		m.showSpinner = false
		m.statusMessage = ""
		m.logViewPort.SetContent(string(msg))
		m.logViewPort.GotoBottom()
		m.logReady = true

	case statusMsg:
		slog.Debug("statusMsg received", "msg", string(msg))
		m.statusMessage = string(msg)
		m.showSpinner = false
		cmds = append(cmds, fetchAllCmd())

	case errMsg:
		slog.Error("errMsg received", "error", msg.err)
		m.err = msg.err
		m.showSpinner = false
		m.statusMessage = ""
	}

	return m, tea.Batch(cmds...)
}
