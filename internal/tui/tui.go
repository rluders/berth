package tui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/rluders/berth/internal/controller"
	"github.com/rluders/berth/internal/engine"
)

var (
	currentTheme = DefaultTheme()
)

type statusMsg string

type ViewType int

const (
	ContainersView ViewType = iota
	ImagesView
	VolumesView
	NetworksView
	SystemView
	InspectView
	LogsView
)

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

func InitialModel() Model {
	containerColumns := []table.Column{
		{Title: "ID", Width: 12},
		{Title: "Image", Width: 20},
		{Title: "Command", Width: 30},
		{Title: "Created", Width: 15},
		{Title: "Status", Width: 20},
		{Title: "Ports", Width: 20},
		{Title: "Names", Width: 20},
	}

	containerTable := table.New(
		table.WithColumns(containerColumns),
		table.WithFocused(true),
		table.WithHeight(10),
	)

	imageColumns := []table.Column{
		{Title: "ID", Width: 15},
		{Title: "Repository", Width: 30},
		{Title: "Tag", Width: 15},
		{Title: "Size", Width: 10},
		{Title: "Created", Width: 20},
	}

	imageTable := table.New(
		table.WithColumns(imageColumns),
		table.WithFocused(false),
		table.WithHeight(10),
	)

	volumeColumns := []table.Column{
		{Title: "Name", Width: 30},
		{Title: "Driver", Width: 15},
		{Title: "Scope", Width: 10},
		{Title: "Mountpoint", Width: 50},
	}

	volumeTable := table.New(
		table.WithColumns(volumeColumns),
		table.WithFocused(false),
		table.WithHeight(10),
	)

	networkColumns := []table.Column{
		{Title: "ID", Width: 15},
		{Title: "Name", Width: 30},
		{Title: "Driver", Width: 15},
		{Title: "Scope", Width: 10},
	}

	networkTable := table.New(
		table.WithColumns(networkColumns),
		table.WithFocused(false),
		table.WithHeight(10),
	)

	s := table.DefaultStyles()
	s.Header = currentTheme.TableHeaderStyle
	s.Selected = currentTheme.TableSelectedStyle
	containerTable.SetStyles(s)
	imageTable.SetStyles(s)
	volumeTable.SetStyles(s)
	networkTable.SetStyles(s)

	return Model{
		engineType:        engine.DetectEngine(),
		currentView:       ContainersView,
		containerTable:    containerTable,
		imageTable:        imageTable,
		volumeTable:       volumeTable,
		networkTable:      networkTable,
		systemInfo:        controller.SystemInfo{}, // Initialize with empty SystemInfo
		inspectViewPort:   viewport.New(0, 0),      // Initialize viewport for inspect
		inspectReady:      false,
		inspectRawContent: "",
		logViewPort:       viewport.New(0, 0), // Initialize viewport
		spinner:           spinner.New(),
	}
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
	switch m.currentView {
	case ContainersView:
		return "1:Containers • 2:Images • 3:Volumes • 4:Networks • 5:System • s:Start • x:Stop • d:Remove • l:Logs • i:Inspect • q:Quit"
	case ImagesView:
		return "1:Containers • 2:Images • 3:Volumes • 4:Networks • 5:System • d:Remove • q:Quit"
	case VolumesView:
		return "1:Containers • 2:Images • 3:Volumes • 4:Networks • 5:System • d:Remove • q:Quit"
	case NetworksView:
		return "1:Containers • 2:Images • 3:Volumes • 4:Networks • 5:System • i:Inspect • q:Quit"
	case SystemView:
		return "1:Containers • 2:Images • 3:Volumes • 4:Networks • 5:System • b:Basic Cleanup • a:Advanced Cleanup • t:Total Cleanup • q:Quit"
	case InspectView:
		return "q/esc:Return • ↑/↓:Scroll"
	case LogsView:
		return "q/esc:Return • ↑/↓:Scroll"
	}
	return "q:Quit"
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
		m.currentView = ContainersView // Fallback to ContainersView if stack is empty
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(fetchContainersCmd(), fetchImagesCmd(), fetchVolumesCmd(), fetchNetworksCmd(), fetchSystemInfoCmd(), m.spinner.Tick)
}
