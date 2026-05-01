package tui

import (
	"fmt"
	"log/slog"
	"strings"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/rluders/berth/internal/controller"
	"github.com/rluders/berth/internal/engine"
)


var currentTheme = DefaultTheme()

// Model represents the main application model.
type Model struct {
	engineType            engine.EngineType
	currentView           ViewType
	containerTable        table.Model
	imageTable            table.Model
	volumeTable           table.Model
	networkTable          table.Model
	systemInfo            controller.SystemInfo
	inspectViewPort       viewport.Model
	inspectReady          bool
	inspectRawContent     string
	logViewPort           viewport.Model
	logReady              bool
	err                   error
	statusMessage         string
	showSpinner           bool
	spinner               spinner.Model
	width                 int
	height                int
	currentLogContainerID string
	currentInspectID      string
	viewStack             []ViewType
}

// InitialModel returns an initialized Model with default values.
func InitialModel() Model {
	slog.Debug("InitialModel called")

	containerTable := table.New(
		table.WithColumns([]table.Column{
			{Title: "ID", Width: 12},
			{Title: "Image", Width: 20},
			{Title: "Command", Width: 30},
			{Title: "Created", Width: 15},
			{Title: "Status", Width: 20},
			{Title: "Ports", Width: 20},
			{Title: "Names", Width: 20},
		}),
		table.WithFocused(true),
		table.WithHeight(0),
	)

	imageTable := table.New(
		table.WithColumns([]table.Column{
			{Title: "ID", Width: 15},
			{Title: "Repository", Width: 30},
			{Title: "Tag", Width: 15},
			{Title: "Size", Width: 10},
			{Title: "Created", Width: 20},
		}),
		table.WithFocused(false),
		table.WithHeight(0),
	)

	volumeTable := table.New(
		table.WithColumns([]table.Column{
			{Title: "Name", Width: 30},
			{Title: "Driver", Width: 15},
			{Title: "Scope", Width: 10},
			{Title: "Mountpoint", Width: 50},
		}),
		table.WithFocused(false),
		table.WithHeight(0),
	)

	networkTable := table.New(
		table.WithColumns([]table.Column{
			{Title: "ID", Width: 15},
			{Title: "Name", Width: 30},
			{Title: "Driver", Width: 15},
			{Title: "Scope", Width: 10},
		}),
		table.WithFocused(false),
		table.WithHeight(0),
	)

	s := table.DefaultStyles()
	s.Header = currentTheme.TableHeaderStyle
	s.Selected = currentTheme.TableSelectedStyle
	containerTable.SetStyles(s)
	imageTable.SetStyles(s)
	volumeTable.SetStyles(s)
	networkTable.SetStyles(s)

	return Model{
		engineType:      engine.DetectEngine(),
		currentView:     ContainersView,
		containerTable:  containerTable,
		imageTable:      imageTable,
		volumeTable:     volumeTable,
		networkTable:    networkTable,
		systemInfo:      controller.SystemInfo{},
		inspectViewPort: viewport.New(0, 0),
		logViewPort:     viewport.New(0, 0),
		spinner:         spinner.New(),
	}
}

// Init initializes the Bubble Tea program.
func (m Model) Init() tea.Cmd {
	slog.Debug("Init called")
	return tea.Batch(fetchAllCmd(), m.spinner.Tick)
}

func (m Model) getViewName() string {
	switch m.currentView {
	case ContainersView:
		return "Containers"
	case ImagesView:
		return "Images"
	case VolumesView:
		return "Volumes"
	case NetworksView:
		return "Networks"
	case SystemView:
		return "System"
	case InspectView:
		return fmt.Sprintf("Inspect %s", m.currentInspectID)
	case LogsView:
		return fmt.Sprintf("Logs for %s", m.currentLogContainerID)
	}
	return "Unknown"
}

func (m Model) getFooterHelp() string {
	nav := "1:Containers • 2:Images • 3:Volumes • 4:Networks • 5:System"
	switch m.currentView {
	case ContainersView:
		return nav + " • s:Start • x:Stop • d:Remove • l:Logs • i:Inspect • q:Quit"
	case ImagesView:
		return nav + " • d:Remove • q:Quit"
	case VolumesView:
		return nav + " • d:Remove • q:Quit"
	case NetworksView:
		return nav + " • i:Inspect • q:Quit"
	case SystemView:
		return nav + " • b:Basic Cleanup • a:Advanced Cleanup • t:Total Cleanup • q:Quit"
	case InspectView, LogsView:
		return "q/esc:Return • ↑/↓:Scroll"
	}
	return "q:Quit"
}

func (m Model) headerText() string {
	return fmt.Sprintf("Berth - %s - %s Engine", m.getViewName(), strings.ToUpper(string(m.engineType)))
}

func (m *Model) pushView(view ViewType) {
	m.viewStack = append(m.viewStack, m.currentView)
	m.currentView = view
}

func (m *Model) popView() {
	if len(m.viewStack) > 0 {
		m.currentView = m.viewStack[len(m.viewStack)-1]
		m.viewStack = m.viewStack[:len(m.viewStack)-1]
	} else {
		m.currentView = ContainersView
	}
}

// contentHeight calculates available height for the main content area
// following Golden Rule #1: subtract header + footer + app padding.
func (m Model) contentHeight() int {
	h := m.height
	h -= lipgloss.Height(currentTheme.HeaderStyle.Render(m.headerText()))
	h -= currentTheme.FooterStyle.GetVerticalFrameSize()
	h -= currentTheme.AppStyle.GetVerticalFrameSize()
	// status message adds an extra line above the footer help
	if m.statusMessage != "" {
		h -= 1
	}
	if h < 0 {
		return 0
	}
	return h
}
