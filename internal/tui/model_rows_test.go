package tui

import (
	"testing"

	"github.com/rluders/berth/internal/controller"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func makeContainer(id, name, image, state, project string) controller.Container {
	labels := map[string]string{}
	if project != "" {
		labels["com.docker.compose.project"] = project
	}
	return controller.Container{ID: id, Names: name, Image: image, State: state, Labels: labels}
}

// --- BuildRows ---

func TestBuildRows_flatListNoLabels(t *testing.T) {
	containers := []controller.Container{
		makeContainer("a", "nginx", "nginx:latest", "running", ""),
		makeContainer("b", "redis", "redis:7", "exited", ""),
	}

	rows := BuildRows(containers, nil)

	require.Len(t, rows, 2)
	assert.Equal(t, RowTypeContainer, rows[0].Type)
	assert.Equal(t, RowTypeContainer, rows[1].Type)
	assert.Empty(t, rows[0].GroupID)
}

func TestBuildRows_groupedByComposeProject(t *testing.T) {
	containers := []controller.Container{
		makeContainer("a", "web", "nginx", "running", "myapp"),
		makeContainer("b", "db", "postgres", "running", "myapp"),
	}

	rows := BuildRows(containers, nil)

	// 1 header + 2 container rows = 3
	require.Len(t, rows, 3)
	assert.Equal(t, RowTypeGroup, rows[0].Type)
	assert.Equal(t, "myapp", rows[0].GroupID)
	assert.Equal(t, RowTypeContainer, rows[1].Type)
	assert.Equal(t, "myapp", rows[1].GroupID)
	assert.Equal(t, RowTypeContainer, rows[2].Type)
}

func TestBuildRows_collapsedGroupHidesChildren(t *testing.T) {
	containers := []controller.Container{
		makeContainer("a", "web", "nginx", "running", "myapp"),
		makeContainer("b", "db", "postgres", "running", "myapp"),
	}
	collapsed := map[string]bool{"myapp": true}

	rows := BuildRows(containers, collapsed)

	// Only header, children hidden
	require.Len(t, rows, 1)
	assert.Equal(t, RowTypeGroup, rows[0].Type)
	assert.True(t, rows[0].Collapsed)
}

func TestBuildRows_mixedGroupAndStandalone(t *testing.T) {
	containers := []controller.Container{
		makeContainer("a", "web", "nginx", "running", "myapp"),
		makeContainer("b", "solo", "redis", "running", ""),
	}

	rows := BuildRows(containers, nil)

	// 1 group header + 1 group child + 1 standalone = 3
	require.Len(t, rows, 3)
	assert.Equal(t, RowTypeGroup, rows[0].Type)
	assert.Equal(t, RowTypeContainer, rows[1].Type)
	assert.Equal(t, RowTypeContainer, rows[2].Type)
	assert.Empty(t, rows[2].GroupID)
}

func TestBuildRows_emptyInput(t *testing.T) {
	rows := BuildRows(nil, nil)
	assert.Nil(t, rows)
}

// --- recomputeRows (filter logic) ---

func TestRecomputeRows_noFilterShowsAll(t *testing.T) {
	m := InitialModel()
	m.containers = []controller.Container{
		makeContainer("a", "nginx", "nginx:latest", "running", ""),
		makeContainer("b", "redis", "redis:7", "exited", ""),
	}

	m.recomputeRows()

	assert.Len(t, m.rows, 2)
}

func TestRecomputeRows_filterByName(t *testing.T) {
	m := InitialModel()
	m.containers = []controller.Container{
		makeContainer("a", "nginx", "nginx:latest", "running", ""),
		makeContainer("b", "redis", "redis:7", "running", ""),
	}
	m.filterInput.SetValue("nginx")

	m.recomputeRows()

	require.Len(t, m.rows, 1)
	assert.Equal(t, "nginx", m.rows[0].Container.Names)
}

func TestRecomputeRows_filterCaseInsensitive(t *testing.T) {
	m := InitialModel()
	m.containers = []controller.Container{
		makeContainer("a", "MyNginx", "nginx:latest", "running", ""),
	}
	m.filterInput.SetValue("mynginx")

	m.recomputeRows()

	assert.Len(t, m.rows, 1)
}

func TestRecomputeRows_filterNoMatch(t *testing.T) {
	m := InitialModel()
	m.containers = []controller.Container{
		makeContainer("a", "nginx", "nginx:latest", "running", ""),
	}
	m.filterInput.SetValue("postgres")

	m.recomputeRows()

	assert.Empty(t, m.rows)
}

func TestRecomputeRows_filterClearedShowsAll(t *testing.T) {
	m := InitialModel()
	m.containers = []controller.Container{
		makeContainer("a", "nginx", "nginx:latest", "running", ""),
		makeContainer("b", "redis", "redis:7", "running", ""),
	}
	m.filterInput.SetValue("nginx")
	m.recomputeRows()
	assert.Len(t, m.rows, 1)

	m.filterInput.SetValue("")
	m.recomputeRows()
	assert.Len(t, m.rows, 2)
}

func TestRecomputeRows_clampsCursorOnFilter(t *testing.T) {
	m := InitialModel()
	m.containers = []controller.Container{
		makeContainer("a", "nginx", "nginx:latest", "running", ""),
		makeContainer("b", "redis", "redis:7", "running", ""),
	}
	m.recomputeRows()
	m.containerCursor = 1 // point at redis

	m.filterInput.SetValue("nginx") // redis disappears
	m.recomputeRows()

	assert.Equal(t, 0, m.containerCursor)
}
