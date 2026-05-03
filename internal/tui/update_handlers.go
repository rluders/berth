package tui

import (
	"fmt"
	"log/slog"
	"strings"

	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/rluders/berth/internal/controller"
)

func (m Model) handleWindowSizeMsg(msg tea.WindowSizeMsg) (Model, tea.Cmd) {
	m.width = msg.Width
	m.height = msg.Height

	m.builtCols = BuildColumns(msg.Width-4, containerCols)

	contentH := m.contentHeight()
	m.containerVP.Width = msg.Width
	// Table header renders 2 lines (text + BorderBottom). Always compute using
	// ContainersView content height (which subtracts the command-preview footer line).
	containerContentH := contentH
	if m.currentView != ContainersView {
		containerContentH-- // ContainersView subtracts extra line for command preview
	}
	m.containerVP.Height = containerContentH - 2 // -2: table header text + border bottom
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
	m.detailsViewPort.Width = viewW
	m.detailsViewPort.Height = contentH

	m.syncContainerViewport()

	if m.currentView == InspectView && !m.inspectReady && m.inspectRawContent != "" {
		m.inspectViewPort.SetContent(prettyJSON(m.inspectRawContent))
		m.inspectReady = true
	}
	if m.currentView == DetailsView && !m.detailsReady && m.currentDetails.ID != "" {
		m.detailsViewPort.SetContent(renderDetailsContent(m.currentDetails))
		m.detailsReady = true
	}

	return m, nil
}

func (m Model) handleContainerListMsg(msg containerListMsg) (Model, tea.Cmd) {
	slog.Debug("containerListMsg", "count", len(msg))
	m.containers = []controller.Container(msg)
	m.recomputeRows()
	m.showSpinner = false
	m.statusMessage = ""
	return m, nil
}

func (m Model) handleImageListMsg(msg imageListMsg) (Model, tea.Cmd) {
	slog.Debug("imageListMsg", "count", len(msg))
	m.images = []controller.Image(msg)
	m.imageTable.SetRows(m.buildImageRows())
	m.showSpinner = false
	m.statusMessage = ""
	return m, nil
}

func (m Model) handleVolumeListMsg(msg volumeListMsg) (Model, tea.Cmd) {
	slog.Debug("volumeListMsg", "count", len(msg))
	m.volumes = []controller.Volume(msg)
	m.volumeTable.SetRows(m.buildVolumeRows())
	m.showSpinner = false
	m.statusMessage = ""
	return m, nil
}

func (m Model) handleNetworkListMsg(msg networkListMsg) (Model, tea.Cmd) {
	slog.Debug("networkListMsg", "count", len(msg))
	rows := make([]table.Row, len(msg))
	for i, net := range msg {
		rows[i] = table.Row{net.ID, net.Name, net.Driver, net.Scope}
	}
	m.networkTable.SetRows(rows)
	m.showSpinner = false
	m.statusMessage = ""
	return m, nil
}

func (m Model) handleSystemInfoMsg(msg systemInfoMsg) (Model, tea.Cmd) {
	m.systemInfo = controller.SystemInfo(msg)
	m.showSpinner = false
	m.statusMessage = ""
	return m, nil
}

func (m Model) handleContainerStatsMsg(msg containerStatsMsg) (Model, tea.Cmd) {
	for id, stat := range msg {
		m.containerStats[id] = stat
	}
	running := map[string]bool{}
	for _, c := range m.containers {
		if c.State == "running" {
			running[c.ID] = true
		}
	}
	for id := range m.containerStats {
		if !running[id] {
			delete(m.containerStats, id)
		}
	}
	m.recomputeRows()
	return m, nil
}

func (m Model) handleStatsTickMsg() (Model, tea.Cmd) {
	var ids []string
	for _, c := range m.containers {
		if strings.HasPrefix(c.Status, "Up") || c.Status == "running" {
			ids = append(ids, c.ID)
		}
	}
	var cmds []tea.Cmd
	if len(ids) > 0 {
		cmds = append(cmds, fetchStatsCmd(ids))
	}
	cmds = append(cmds, statsTickCmd())
	return m, tea.Batch(cmds...)
}

func (m Model) handleRefreshTickMsg() (Model, tea.Cmd) {
	return m, tea.Batch(fetchContainersCmd(), refreshTickCmd())
}

func (m Model) handleInspectMsg(msg inspectMsg) (Model, tea.Cmd) {
	m.showSpinner = false
	m.statusMessage = ""
	m.inspectRawContent = string(msg)
	if m.width > 0 {
		m.inspectViewPort.SetContent(prettyJSON(string(msg)))
		m.inspectReady = true
	}
	return m, func() tea.Msg {
		return tea.WindowSizeMsg{Width: m.width, Height: m.height}
	}
}

func (m Model) handleDetailsMsg(msg detailsMsg) (Model, tea.Cmd) {
	m.showSpinner = false
	m.statusMessage = ""
	m.currentDetails = controller.ContainerDetails(msg)
	if m.width > 0 {
		m.detailsViewPort.SetContent(renderDetailsContent(m.currentDetails))
		m.detailsReady = true
	}
	return m, func() tea.Msg {
		return tea.WindowSizeMsg{Width: m.width, Height: m.height}
	}
}

func (m Model) handleLogChunkMsg(msg logChunkMsg) (Model, tea.Cmd) {
	m.logLines = append(m.logLines, string(msg))
	if len(m.logLines) > 10000 {
		m.logLines = m.logLines[len(m.logLines)-5000:]
	}
	m.logViewPort.SetContent(buildColorizedLogContent(m.logLines, m.showLineNumbers))
	if m.logFollowing {
		m.logViewPort.GotoBottom()
	}
	if m.logCh != nil {
		return m, waitForLogLineCmd(m.logCh)
	}
	return m, nil
}

func (m Model) handleLogStreamDoneMsg() (Model, tea.Cmd) {
	m.logCh = nil
	m.logCancel = nil
	return m, nil
}

func (m Model) handleProgressMsg(msg progressMsg) (Model, tea.Cmd) {
	if msg.done {
		m.progressLabel = msg.label
		m.progressDone = true
		m.showSpinner = false
		return m, tea.Batch(
			m.progressBar.SetPercent(1.0),
			func() tea.Msg { return statusMsg(msg.label) },
		)
	}
	m.progressVisible = true
	m.progressLabel = msg.label
	return m, m.progressBar.SetPercent(msg.percent)
}

func (m Model) handleProgressTickMsg() (Model, tea.Cmd) {
	if !m.progressVisible || m.progressDone {
		return m, nil
	}
	// Animate toward 0.85 while waiting for actual completion.
	next := m.progressBar.Percent() + 0.04
	if next > 0.85 {
		next = 0.85
	}
	return m, tea.Batch(m.progressBar.SetPercent(next), progressTickCmd())
}

func (m Model) handleProgressFrameMsg(msg progress.FrameMsg) (Model, tea.Cmd) {
	raw, cmd := m.progressBar.Update(msg)
	if pb, ok := raw.(progress.Model); ok {
		m.progressBar = pb
	}
	return m, cmd
}

func (m Model) handleStatusMsg(msg statusMsg) (Model, tea.Cmd) {
	slog.Debug("statusMsg", "msg", string(msg))
	m.statusMessage = string(msg)
	m.showSpinner = false
	m.progressVisible = false
	m.progressDone = false
	return m, fetchAllCmd()
}

func (m Model) handleComposeOutputMsg(msg composeOutputMsg) (Model, tea.Cmd) {
	m.composeOutput = append(m.composeOutput, msg.line)
	if len(m.composeOutput) > 200 {
		m.composeOutput = m.composeOutput[1:]
	}
	m.statusMessage = msg.line
	return m, readNextComposeLineCmd(msg.ch, msg.project)
}

func (m Model) handleComposeDoneMsg(msg composeDoneMsg) (Model, tea.Cmd) {
	m.showSpinner = false
	m.composeCancel = nil
	if msg.err != nil {
		m.statusMessage = fmt.Sprintf("[%s] compose failed: %v", msg.project, msg.err)
	} else {
		m.statusMessage = fmt.Sprintf("[%s] compose done.", msg.project)
	}
	return m, fetchContainersCmd()
}

func (m Model) handleErrMsg(msg errMsg) (Model, tea.Cmd) {
	slog.Error("errMsg", "error", msg.err)
	m.err = msg.err
	m.showSpinner = false
	m.statusMessage = ""
	return m, nil
}

// renderDetailsContent formats ContainerDetails into a card-based scrollable view.
func renderDetailsContent(d controller.ContainerDetails) string {
	th := currentTheme

	field := func(label, value string) string {
		l := th.CardTitleStyle.Render(fmt.Sprintf("%-12s", label))
		v := th.CardValueStyle.Render(value)
		return "  " + l + "  " + v
	}

	section := func(title string, lines []string) string {
		header := th.SectionStyle.Render("▸ " + title)
		body := strings.Join(lines, "\n")
		return th.CardStyle.Render(header + "\n" + body)
	}

	// ── Container section ──────────────────────────────────────────────────
	stateBadge := StatusBadge(d.State)
	infoLines := []string{
		field("ID", d.ID),
		field("Name", d.Name),
		field("Image", d.Image),
		field("Command", d.Command),
		field("State", stateBadge),
		field("Created", d.Created),
	}

	// ── Environment section ────────────────────────────────────────────────
	envLines := []string{}
	if len(d.Env) == 0 {
		envLines = append(envLines, "  "+th.CardValueStyle.Foreground(lipgloss.Color(colorMuted)).Render("(none)"))
	} else {
		for _, e := range d.Env {
			parts := strings.SplitN(e, "=", 2)
			if len(parts) == 2 {
				masked := strings.Repeat("•", len(parts[1]))
				if len(masked) > 12 {
					masked = masked[:12]
				}
				envLines = append(envLines, field(parts[0], masked))
			} else {
				envLines = append(envLines, "  "+e)
			}
		}
	}

	// ── Ports section ──────────────────────────────────────────────────────
	portLines := []string{}
	if len(d.Ports) == 0 {
		portLines = append(portLines, "  "+th.CardValueStyle.Foreground(lipgloss.Color(colorMuted)).Render("(none)"))
	} else {
		for _, p := range d.Ports {
			hostIP := p.HostIP
			if hostIP == "" {
				hostIP = "0.0.0.0"
			}
			line := fmt.Sprintf("%s/%s → %s:%s", p.ContainerPort, p.Protocol, hostIP, p.HostPort)
			portLines = append(portLines, "  "+th.CardValueStyle.Render(line))
		}
	}

	// ── Mounts section ─────────────────────────────────────────────────────
	mountLines := []string{}
	if len(d.Mounts) == 0 {
		mountLines = append(mountLines, "  "+th.CardValueStyle.Foreground(lipgloss.Color(colorMuted)).Render("(none)"))
	} else {
		for _, mt := range d.Mounts {
			rw := th.LogDebugStyle.Render("ro")
			if mt.RW {
				rw = th.LogInfoStyle.Render("rw")
			}
			line := fmt.Sprintf("[%s] %s → %s (%s)", mt.Type, mt.Source, mt.Destination, rw)
			mountLines = append(mountLines, "  "+th.CardValueStyle.Render(line))
		}
	}

	// ── Networks section ───────────────────────────────────────────────────
	netLines := []string{}
	if len(d.Networks) == 0 {
		netLines = append(netLines, "  "+th.CardValueStyle.Foreground(lipgloss.Color(colorMuted)).Render("(none)"))
	} else {
		for _, n := range d.Networks {
			line := fmt.Sprintf("%s  IP: %s  GW: %s", n.Name, n.IPAddress, n.Gateway)
			netLines = append(netLines, "  "+th.CardValueStyle.Render(line))
		}
	}

	return strings.Join([]string{
		section("Container", infoLines),
		section("Environment", envLines),
		section("Ports", portLines),
		section("Mounts", mountLines),
		section("Networks", netLines),
	}, "\n")
}
