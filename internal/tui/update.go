package tui

import (
	"fmt"
	"log/slog"
	"strings"

	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
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

	case tea.MouseMsg:
		newM, cmd := m.handleMouseMsg(msg)
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
		m.detailsViewPort.Width = viewW
		m.detailsViewPort.Height = contentH

		if m.currentView == InspectView && !m.inspectReady && m.inspectRawContent != "" {
			m.inspectViewPort.SetContent(prettyJSON(m.inspectRawContent))
			m.inspectReady = true
		}
		if m.currentView == DetailsView && !m.detailsReady && m.currentDetails.ID != "" {
			m.detailsViewPort.SetContent(renderDetailsContent(m.currentDetails))
			m.detailsReady = true
		}

	// ── Container list ────────────────────────────────────────────────────

	case containerListMsg:
		slog.Debug("containerListMsg", "count", len(msg))
		m.containers = []controller.Container(msg)
		m.containerTable.SetRows(m.buildContainerRows())
		m.showSpinner = false
		m.statusMessage = ""

	case imageListMsg:
		slog.Debug("imageListMsg", "count", len(msg))
		m.images = []controller.Image(msg)
		m.imageTable.SetRows(m.buildImageRows())
		m.showSpinner = false
		m.statusMessage = ""

	case volumeListMsg:
		slog.Debug("volumeListMsg", "count", len(msg))
		m.volumes = []controller.Volume(msg)
		m.volumeTable.SetRows(m.buildVolumeRows())
		m.showSpinner = false
		m.statusMessage = ""

	case networkListMsg:
		slog.Debug("networkListMsg", "count", len(msg))
		rows := make([]table.Row, len(msg))
		for i, net := range msg {
			rows[i] = table.Row{net.ID, net.Name, net.Driver, net.Scope}
		}
		m.networkTable.SetRows(rows)
		m.showSpinner = false
		m.statusMessage = ""

	case systemInfoMsg:
		m.systemInfo = controller.SystemInfo(msg)
		m.showSpinner = false
		m.statusMessage = ""

	// ── Stats ─────────────────────────────────────────────────────────────

	case containerStatsMsg:
		for id, stat := range msg {
			m.containerStats[id] = stat
		}
		m.containerTable.SetRows(m.buildContainerRows())

	case statsTickMsg:
		var ids []string
		for _, c := range m.containers {
			if strings.HasPrefix(c.Status, "Up") || c.Status == "running" {
				ids = append(ids, c.ID)
			}
		}
		if len(ids) > 0 {
			cmds = append(cmds, fetchStatsCmd(ids))
		}
		cmds = append(cmds, statsTickCmd())

	case refreshTickMsg:
		cmds = append(cmds, fetchContainersCmd(), refreshTickCmd())

	// ── Inspect ───────────────────────────────────────────────────────────

	case inspectMsg:
		m.showSpinner = false
		m.statusMessage = ""
		m.inspectRawContent = string(msg)
		if m.width > 0 {
			m.inspectViewPort.SetContent(prettyJSON(string(msg)))
			m.inspectReady = true
		}
		cmds = append(cmds, func() tea.Msg {
			return tea.WindowSizeMsg{Width: m.width, Height: m.height}
		})

	// ── Details ───────────────────────────────────────────────────────────

	case detailsMsg:
		m.showSpinner = false
		m.statusMessage = ""
		m.currentDetails = controller.ContainerDetails(msg)
		if m.width > 0 {
			m.detailsViewPort.SetContent(renderDetailsContent(m.currentDetails))
			m.detailsReady = true
		}
		cmds = append(cmds, func() tea.Msg {
			return tea.WindowSizeMsg{Width: m.width, Height: m.height}
		})

	// ── Logs ──────────────────────────────────────────────────────────────

	case logChunkMsg:
		m.logLines = append(m.logLines, string(msg))
		if len(m.logLines) > 10000 {
			m.logLines = m.logLines[len(m.logLines)-5000:]
		}
		m.logViewPort.SetContent(buildColorizedLogContent(m.logLines, m.showLineNumbers))
		if m.logFollowing {
			m.logViewPort.GotoBottom()
		}
		if m.logCh != nil {
			cmds = append(cmds, waitForLogLineCmd(m.logCh))
		}

	case logStreamDoneMsg:
		m.logCh = nil
		m.logCancel = nil

	// ── Progress bar ──────────────────────────────────────────────────────

	case progressMsg:
		if msg.done {
			cmds = append(cmds, m.progressBar.SetPercent(1.0))
			m.progressLabel = msg.label
			m.progressDone = true
			m.showSpinner = false
			// Emit statusMsg after bar reaches 100%.
			cmds = append(cmds, func() tea.Msg { return statusMsg(msg.label) })
		} else {
			m.progressVisible = true
			m.progressLabel = msg.label
			cmds = append(cmds, m.progressBar.SetPercent(msg.percent))
		}

	case progressTickMsg:
		if m.progressVisible && !m.progressDone {
			// Animate toward 0.85 while waiting for actual completion.
			next := m.progressBar.Percent() + 0.04
			if next > 0.85 {
				next = 0.85
			}
			cmds = append(cmds, m.progressBar.SetPercent(next), progressTickCmd())
		}

	case progress.FrameMsg:
		raw, cmd := m.progressBar.Update(msg)
		if pb, ok := raw.(progress.Model); ok {
			m.progressBar = pb
		}
		cmds = append(cmds, cmd)

	// ── Status / errors ───────────────────────────────────────────────────

	case statusMsg:
		slog.Debug("statusMsg", "msg", string(msg))
		m.statusMessage = string(msg)
		m.showSpinner = false
		m.progressVisible = false
		m.progressDone = false
		cmds = append(cmds, fetchAllCmd())

	case errMsg:
		slog.Error("errMsg", "error", msg.err)
		m.err = msg.err
		m.showSpinner = false
		m.statusMessage = ""
	}

	return m, tea.Batch(cmds...)
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
