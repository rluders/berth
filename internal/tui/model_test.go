package tui

import (
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
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

func TestSimplifyImage(t *testing.T) {
	cases := []struct{ input, want string }{
		{"docker.io/library/postgres:16", "postgres:16"},
		{"nginx:latest", "nginx:latest"},
		{"nginx", "nginx"},
		{"<none>", "<none>"},
		{"gcr.io/google-containers/pause:3.1", "pause:3.1"},
		{"sha256:a1b2c3d4e5f6789012345678901234567890123456789012345678901234", "sha256:a1b2c3d4e5f6"},
		{"sha256:abc", "sha256:abc"},
	}
	for _, c := range cases {
		assert.Equal(t, c.want, simplifyImage(c.input), "input: %s", c.input)
	}
}

func windowSize(width, height int) tea.WindowSizeMsg {
	return tea.WindowSizeMsg{Width: width, Height: height}
}
