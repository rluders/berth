package tui

import (
	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
)

// Update dispatches incoming messages to the appropriate handler.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd

	case tea.KeyMsg:
		return m.handleKeyMsg(msg)

	case tea.MouseMsg:
		return m.handleMouseMsg(msg)

	case tea.WindowSizeMsg:
		return m.handleWindowSizeMsg(msg)

	case containerListMsg:
		return m.handleContainerListMsg(msg)

	case imageListMsg:
		return m.handleImageListMsg(msg)

	case volumeListMsg:
		return m.handleVolumeListMsg(msg)

	case networkListMsg:
		return m.handleNetworkListMsg(msg)

	case systemInfoMsg:
		return m.handleSystemInfoMsg(msg)

	case containerStatsMsg:
		return m.handleContainerStatsMsg(msg)

	case statsTickMsg:
		return m.handleStatsTickMsg()

	case refreshTickMsg:
		return m.handleRefreshTickMsg()

	case inspectMsg:
		return m.handleInspectMsg(msg)

	case detailsMsg:
		return m.handleDetailsMsg(msg)

	case logChunkMsg:
		return m.handleLogChunkMsg(msg)

	case logStreamDoneMsg:
		return m.handleLogStreamDoneMsg()

	case progressMsg:
		return m.handleProgressMsg(msg)

	case progressTickMsg:
		return m.handleProgressTickMsg()

	case progress.FrameMsg:
		return m.handleProgressFrameMsg(msg)

	case statusMsg:
		return m.handleStatusMsg(msg)

	case composeOutputMsg:
		return m.handleComposeOutputMsg(msg)

	case composeDoneMsg:
		return m.handleComposeDoneMsg(msg)

	case errMsg:
		return m.handleErrMsg(msg)
	}

	return m, nil
}
