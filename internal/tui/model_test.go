package tui

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/stretchr/testify/assert"
)

func TestInitialModel_defaultView(t *testing.T) {
	m := InitialModel()
	assert.Equal(t, ContainersView, m.currentView)
}

func TestInitialModel_mapsInitialized(t *testing.T) {
	m := InitialModel()
	assert.NotNil(t, m.containerStats)
	assert.NotNil(t, m.collapsedGroups)
}

func TestInitialModel_spinnerReady(t *testing.T) {
	m := InitialModel()
	cmd := m.Init()
	assert.NotNil(t, cmd)
}

func TestRenderContainerHeader_fillsModelWidth(t *testing.T) {
	m := InitialModel()
	result, _ := updateModel(t, m, windowSize(120, 40))

	for _, line := range strings.Split(result.renderContainerHeader(), "\n") {
		assert.Equal(t, result.width, lipgloss.Width(line))
	}
}

func TestRenderContainerSelectedRow_fillsModelWidth(t *testing.T) {
	m := InitialModel()
	result, _ := updateModel(t, m, windowSize(120, 40))
	row := Row{Type: RowTypeGroup, GroupID: "project"}

	line := result.renderContainerViewRow(row, true)

	assert.Equal(t, result.width, lipgloss.Width(line))
}

func windowSize(width, height int) tea.WindowSizeMsg {
	return tea.WindowSizeMsg{Width: width, Height: height}
}
