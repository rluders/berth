package tui

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"charm.land/bubbles/v2/help"
	"charm.land/bubbles/v2/progress"
	"charm.land/bubbles/v2/spinner"
	"charm.land/bubbles/v2/table"
	"charm.land/bubbles/v2/textinput"
	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"
	"github.com/charmbracelet/x/ansi"
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

	// Container viewport (replaces table.Model — bypasses double-styling that hides Status)
	containerVP     viewport.Model
	containerCursor int

	// Tables (images, volumes, networks still use bubbles/table)
	imageTable   table.Model
	volumeTable  table.Model
	networkTable table.Model

	// Raw data (for filtering / grouping)
	containers []controller.Container
	images     []controller.Image
	volumes    []controller.Volume

	// Container stats
	containerStats map[string]controller.ContainerStat

	// Accordion state
	collapsedGroups map[string]bool
	rows            []Row

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
	currentLogGroupName   string
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

	// Computed columns for the containers table (recomputed on resize).
	builtCols []Column

	// Last action key pressed in ContainersView (drives command preview).
	lastActionKey string

	// Compose streaming state
	composeOutput []string           // rolling 200-line buffer of streamed compose output
	composeCancel context.CancelFunc // cancels the running compose operation, nil when idle
}

// InitialModel returns an initialized Model with default values.
func InitialModel() Model {
	slog.Debug("InitialModel called")

	initCols := BuildColumns(116, containerCols) // 120-col default until first WindowSizeMsg

	imageTable := table.New(
		table.WithColumns(tableColumns(120, imageCols)),
		table.WithFocused(false),
		table.WithHeight(0),
	)

	volumeTable := table.New(
		table.WithColumns(tableColumns(120, volumeCols)),
		table.WithFocused(false),
		table.WithHeight(0),
	)

	networkTable := table.New(
		table.WithColumns(tableColumns(120, networkCols)),
		table.WithFocused(false),
		table.WithHeight(0),
	)

	s := tableStyles()
	imageTable.SetStyles(s)
	volumeTable.SetStyles(s)
	networkTable.SetStyles(s)

	fi := textinput.New()
	fi.Placeholder = "filter..."
	fi.CharLimit = 60

	return Model{
		engineType:      engine.DetectEngine(),
		currentView:     ContainersView,
		containerVP:     viewport.New(),
		builtCols:       initCols,
		imageTable:      imageTable,
		volumeTable:     volumeTable,
		networkTable:    networkTable,
		containerStats:  make(map[string]controller.ContainerStat),
		collapsedGroups: make(map[string]bool),
		systemInfo:      controller.SystemInfo{},
		inspectViewPort: viewport.New(),
		logViewPort:     viewport.New(),
		detailsViewPort: viewport.New(),
		logFollowing:    true,
		filterInput:     fi,
		spinner:         spinner.New(),
		helpModel:       help.New(),
		progressBar: progress.New(
			progress.WithDefaultBlend(),
			progress.WithoutPercentage(),
		),
	}
}

// Init initializes the Bubble Tea program.
func (m Model) Init() tea.Cmd {
	slog.Debug("Init called")
	return tea.Batch(fetchAllCmd(), m.spinner.Tick, statsTickCmd(), refreshTickCmd())
}

func tableStyles() table.Styles {
	s := table.DefaultStyles()
	s.Header = currentTheme.TableHeaderStyle.Padding(0, 1)
	s.Selected = currentTheme.TableSelectedStyle
	return s
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
		if m.currentLogGroupName != "" {
			return fmt.Sprintf("Logs  %s", m.currentLogGroupName)
		}
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
	h -= 2 // footer (key hints + base)
	if m.currentView == ContainersView {
		h -= 1 // command preview line
	}
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

// recomputeRows applies filter, rebuilds m.rows via BuildRows, and syncs the viewport.
func (m *Model) recomputeRows() {
	filter := strings.ToLower(m.filterInput.Value())
	var filtered []controller.Container
	for _, c := range m.containers {
		if filter != "" {
			haystack := strings.ToLower(c.Names + " " + c.Image + " " + c.Status + " " + c.State)
			if !strings.Contains(haystack, filter) {
				continue
			}
		}
		filtered = append(filtered, c)
	}
	m.rows = BuildRows(filtered, m.collapsedGroups)
	// Clamp cursor after filter may reduce row count.
	if len(m.rows) > 0 && m.containerCursor >= len(m.rows) {
		m.containerCursor = len(m.rows) - 1
	}
	m.syncContainerViewport()
}

// renderContainerHeader returns a styled header line for the containers viewport.
func (m Model) renderContainerHeader() string {
	cells := make([]string, len(m.builtCols))
	for i, col := range m.builtCols {
		cells[i] = renderCell(padHeader(col.Header, col.Width, col.Align), col.Width, col.Align)
	}
	return currentTheme.TableHeaderStyle.Width(m.width).Render(strings.Join(cells, " "))
}

// renderContainerViewRow renders one row as a full-width string with optional selection highlight.
func (m Model) renderContainerViewRow(row Row, selected bool) string {
	var values []string
	switch row.Type {
	case RowTypeGroup:
		running, total := groupAggStatus(row.Containers)
		label := GroupStatusColor(running, total)
		prefix := "▼ "
		if row.Collapsed {
			prefix = "▶ "
		}
		values = []string{currentTheme.GroupHeaderStyle.Render(prefix + row.GroupID), label, "", "", "", "", ""}

	case RowTypeContainer:
		c := row.Container
		cpuStr := "-"
		memStr := "-"
		if c.State == "running" {
			stat, ok := m.containerStats[c.ID]
			if !ok {
				cpuStr = "..."
				memStr = "..."
			} else {
				cpuStr = fmt.Sprintf("%.1f", stat.CPUPercent)
				if stat.MemLimit > 0 {
					memStr = utils.FormatBytes(stat.MemUsage)
				} else {
					memStr = "..."
				}
			}
		}
		name := c.Names
		if row.GroupID != "" {
			name = currentTheme.GroupChildStyle.Render("  › " + c.Names)
		}
		values = []string{
			name,
			FormatStatus(c.State),
			simplifyImage(c.Image),
			c.Ports,
			cpuStr,
			memStr,
			utils.FormatAge(c.CreatedAt),
		}
	}

	line := strings.Join(RenderRow(m.builtCols, values), " ")
	if selected {
		// Strip ANSI from pre-styled cells so selection background renders uniformly.
		line = currentTheme.TableSelectedStyle.Width(m.width).Render(ansi.Strip(line))
	}
	return line
}

// syncContainerViewport re-renders all rows into the viewport and scrolls to keep cursor visible.
func (m *Model) syncContainerViewport() {
	lines := make([]string, len(m.rows))
	for i, row := range m.rows {
		lines[i] = m.renderContainerViewRow(row, i == m.containerCursor)
	}
	m.containerVP.SetContent(strings.Join(lines, "\n"))
	// Ensure cursor is visible.
	if m.containerCursor < m.containerVP.YOffset() {
		m.containerVP.SetYOffset(m.containerCursor)
	} else if m.containerVP.Height() > 0 && m.containerCursor >= m.containerVP.YOffset()+m.containerVP.Height() {
		m.containerVP.SetYOffset(m.containerCursor - m.containerVP.Height() + 1)
	}
}

// moveContainerCursor moves the cursor by delta, clamped to valid row range.
func (m *Model) moveContainerCursor(delta int) {
	n := len(m.rows)
	if n == 0 {
		return
	}
	m.containerCursor = max(0, min(m.containerCursor+delta, n-1))
	m.syncContainerViewport()
}

// simplifyImage strips registry/org prefixes, returning only the last path
// segment (image name + tag). e.g. docker.io/library/postgres:16 → postgres:16
func simplifyImage(img string) string {
	parts := strings.Split(img, "/")
	return parts[len(parts)-1]
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
