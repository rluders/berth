package tui

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/rluders/berth/internal/controller"
	"github.com/rluders/berth/internal/engine"
	"github.com/rluders/berth/internal/utils"
)

var currentTheme = DefaultTheme()

// Model represents the main application model.
type Model struct {
	engineType  engine.EngineType
	currentView ViewType
	viewStack   []ViewType

	// Tables
	containerTable table.Model
	imageTable     table.Model
	volumeTable    table.Model
	networkTable   table.Model

	// Raw data (for filtering / grouping)
	containers []controller.Container
	images     []controller.Image
	volumes    []controller.Volume

	// Container stats
	containerStats map[string]controller.ContainerStat

	// Compose group toggle
	groupByCompose bool

	// Accordion state (only active when groupByCompose == true)
	collapsedGroups      map[string]bool
	containerVisibleRows []containerRowMeta

	// System info
	systemInfo controller.SystemInfo

	// Inspect view
	inspectViewPort   viewport.Model
	inspectReady      bool
	inspectRawContent string
	currentInspectID  string

	// Logs view
	logViewPort           viewport.Model
	logReady              bool
	logLines              []string
	logFollowing          bool
	logCh                 chan string
	logCancel             context.CancelFunc
	currentLogContainerID string
	showLineNumbers       bool

	// Details view
	detailsViewPort  viewport.Model
	detailsReady     bool
	currentDetailsID string
	currentDetails   controller.ContainerDetails

	// Search / filter
	filterInput  textinput.Model
	filterActive bool

	// Modal dialog (replaces old confirmAction)
	modal *Modal

	// Help overlay
	showHelp  bool
	helpModel help.Model

	// Progress bar (cleanup / prune operations)
	progressBar     progress.Model
	progressVisible bool
	progressLabel   string
	progressDone    bool

	// Status
	err           error
	statusMessage string
	showSpinner   bool
	spinner       spinner.Model

	// Window
	width  int
	height int

	// Column widths for the containers table (recomputed on resize).
	containerColWidths []int
}

// InitialModel returns an initialized Model with default values.
func InitialModel() Model {
	slog.Debug("InitialModel called")

	initWidths := computeWidths(116, containerCols) // 120-col default until first WindowSizeMsg
	containerTable := table.New(
		table.WithColumns(buildTableColumns(initWidths, containerCols)),
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

	fi := textinput.New()
	fi.Placeholder = "filter..."
	fi.CharLimit = 60

	return Model{
		engineType:          engine.DetectEngine(),
		currentView:         ContainersView,
		containerTable:      containerTable,
		containerColWidths:  initWidths,
		imageTable:     imageTable,
		volumeTable:    volumeTable,
		networkTable:   networkTable,
		containerStats:  make(map[string]controller.ContainerStat),
		collapsedGroups: make(map[string]bool),
		systemInfo:     controller.SystemInfo{},
		inspectViewPort: viewport.New(0, 0),
		logViewPort:     viewport.New(0, 0),
		detailsViewPort: viewport.New(0, 0),
		logFollowing: true,
		filterInput:  fi,
		spinner:      spinner.New(),
		helpModel:    help.New(),
		progressBar: progress.New(
			progress.WithDefaultGradient(),
			progress.WithoutPercentage(),
		),
	}
}

// Init initializes the Bubble Tea program.
func (m Model) Init() tea.Cmd {
	slog.Debug("Init called")
	return tea.Batch(fetchAllCmd(), m.spinner.Tick, statsTickCmd(), refreshTickCmd())
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
		return fmt.Sprintf("Logs  %s", m.currentLogContainerID)
	case DetailsView:
		return fmt.Sprintf("Details  %s", m.currentDetailsID)
	}
	return "Unknown"
}


func (m Model) headerText() string {
	eng := strings.ToUpper(string(m.engineType))
	view := m.getViewName()
	extra := ""
	if m.currentView == LogsView {
		mode := "follow"
		if !m.logFollowing {
			mode = "paused"
		}
		extra = fmt.Sprintf(" [%s]", mode)
	}
	if m.groupByCompose && m.currentView == ContainersView {
		extra = " [grouped]"
	}
	return fmt.Sprintf("Berth  %s  %s Engine%s", view, eng, extra)
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

// contentHeight calculates available height for the main content area.
func (m Model) contentHeight() int {
	h := m.height
	h -= 1 // header row
	h -= 1 // tab bar row
	h -= 2 // footer (key hints + engine line)
	if m.statusMessage != "" || m.showSpinner {
		h -= 1
	}
	if m.filterActive {
		h -= 1
	}
	if m.progressVisible {
		h -= 2 // label + bar
	}
	if h < 0 {
		return 0
	}
	return h
}

// buildContainerRows produces filtered, optionally compose-grouped table rows
// and a parallel metadata slice. Both slices always have the same length.
func (m Model) buildContainerRows() ([]table.Row, []containerRowMeta) {
	filter := strings.ToLower(m.filterInput.Value())

	// Apply filter to raw containers.
	var filtered []controller.Container
	for _, c := range m.containers {
		if filter != "" {
			haystack := strings.ToLower(c.Names + " " + c.Image + " " + c.Status)
			if !strings.Contains(haystack, filter) {
				continue
			}
		}
		filtered = append(filtered, c)
	}

	var rows []table.Row
	var metas []containerRowMeta

	containerRow := func(c controller.Container, namePrefix string) table.Row {
		stat := m.containerStats[c.ID]
		cpuStr := fmt.Sprintf("%.1f", stat.CPUPercent)
		memStr := ""
		if stat.MemLimit > 0 {
			memStr = utils.FormatBytes(stat.MemUsage)
		}
		values := []string{
			namePrefix + c.Names,
			StatusColor(c.Status),
			c.Image,
			c.Ports,
			cpuStr,
			memStr,
			utils.FormatAge(c.CreatedAt),
		}
		row := make(table.Row, len(containerCols))
		for i, v := range values {
			row[i] = renderCell(v, m.containerColWidths[i], containerCols[i].Align)
		}
		return row
	}

	if !m.groupByCompose {
		for _, c := range filtered {
			rows = append(rows, containerRow(c, ""))
			metas = append(metas, containerRowMeta{
				kind:          rowKindContainer,
				containerID:   c.ID,
				containerName: c.Names,
			})
		}
		return rows, metas
	}

	// Grouped mode: split into compose groups and standalone.
	groups, standalone := buildComposeGroups(filtered)

	for _, g := range groups {
		running, total := groupAggStatus(g.containers)
		label := aggStatusLabel(running, total)
		collapsed := m.collapsedGroups[g.project]

		prefix := "▼ "
		if collapsed {
			prefix = "▶ "
		}

		groupValues := []string{prefix + g.project, label, "", "", "", "", ""}
		groupRow := make(table.Row, len(containerCols))
		for i, v := range groupValues {
			groupRow[i] = renderCell(v, m.containerColWidths[i], containerCols[i].Align)
		}
		rows = append(rows, groupRow)
		metas = append(metas, containerRowMeta{kind: rowKindGroup, groupName: g.project})

		if !collapsed {
			for _, c := range g.containers {
				rows = append(rows, containerRow(c, "  › "))
				metas = append(metas, containerRowMeta{
					kind:          rowKindContainer,
					groupName:     g.project,
					containerID:   c.ID,
					containerName: c.Names,
				})
			}
		}
	}

	// Standalone containers (no compose label) appended without a group header.
	for _, c := range standalone {
		rows = append(rows, containerRow(c, ""))
		metas = append(metas, containerRowMeta{
			kind:          rowKindContainer,
			containerID:   c.ID,
			containerName: c.Names,
		})
	}

	return rows, metas
}

// buildImageRows produces filtered image rows.
func (m Model) buildImageRows() []table.Row {
	filter := strings.ToLower(m.filterInput.Value())
	var rows []table.Row
	for _, img := range m.images {
		if filter != "" {
			if !strings.Contains(strings.ToLower(img.Repository+" "+img.Tag), filter) {
				continue
			}
		}
		rows = append(rows, table.Row{img.ID, img.Repository, img.Tag, img.Size, img.Created})
	}
	return rows
}

// buildVolumeRows produces filtered volume rows.
func (m Model) buildVolumeRows() []table.Row {
	filter := strings.ToLower(m.filterInput.Value())
	var rows []table.Row
	for _, v := range m.volumes {
		if filter != "" {
			if !strings.Contains(strings.ToLower(v.Name), filter) {
				continue
			}
		}
		rows = append(rows, table.Row{v.Name, v.Driver, v.Scope, v.Mountpoint})
	}
	return rows
}
