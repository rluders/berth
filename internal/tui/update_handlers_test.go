package tui

import (
	"errors"
	"strings"
	"testing"

	"charm.land/bubbles/v2/table"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/rluders/berth/internal/controller"
	"github.com/stretchr/testify/assert"
)

func TestHandleContainerListMsg_setsContainers(t *testing.T) {
	m := InitialModel()
	containers := []controller.Container{{ID: "abc123", Names: "test", State: "running"}}

	result, cmd := updateModel(t, m, containerListMsg(containers))

	assert.Len(t, result.containers, 1)
	assert.Equal(t, "abc123", result.containers[0].ID)
	assert.False(t, result.showSpinner)
	assert.Empty(t, result.statusMessage)
	assert.Nil(t, cmd)
}

func TestHandleContainerListMsg_empty(t *testing.T) {
	m := InitialModel()
	m.showSpinner = true

	result, _ := updateModel(t, m, containerListMsg(nil))

	assert.Empty(t, result.containers)
	assert.False(t, result.showSpinner)
}

func TestHandleImageListMsg_setsImages(t *testing.T) {
	m := InitialModel()
	images := []controller.Image{{ID: "img1", Repository: "nginx", Tag: "latest"}}

	result, cmd := updateModel(t, m, imageListMsg(images))

	assert.Len(t, result.images, 1)
	assert.False(t, result.showSpinner)
	assert.Nil(t, cmd)
}

func TestHandleErrMsg_setsError(t *testing.T) {
	m := InitialModel()
	m.showSpinner = true
	m.statusMessage = "doing stuff"
	testErr := errors.New("something broke")

	result, cmd := updateModel(t, m, errMsg{err: testErr})

	assert.Equal(t, testErr, result.err)
	assert.False(t, result.showSpinner)
	assert.Empty(t, result.statusMessage)
	assert.Nil(t, cmd)
}

func TestHandleComposeOutputMsg_appendsToBuffer(t *testing.T) {
	m := InitialModel()
	ch := make(chan string, 1)
	close(ch)

	result, _ := updateModel(t, m, composeOutputMsg{project: "myapp", line: "Pulling image...", ch: ch})

	assert.Len(t, result.composeOutput, 1)
	assert.Equal(t, "Pulling image...", result.composeOutput[0])
	assert.Equal(t, "Pulling image...", result.statusMessage)
}

func TestHandleComposeOutputMsg_rollingBuffer(t *testing.T) {
	m := InitialModel()
	m.composeOutput = make([]string, 200)
	ch := make(chan string, 1)
	close(ch)

	result, _ := updateModel(t, m, composeOutputMsg{project: "myapp", line: "new line", ch: ch})

	assert.Len(t, result.composeOutput, 200, "buffer must not exceed 200 entries")
	assert.Equal(t, "new line", result.composeOutput[199])
}

func TestHandleComposeDoneMsg_success(t *testing.T) {
	m := InitialModel()
	m.showSpinner = true
	cancelCalled := false
	m.composeCancel = func() { cancelCalled = true }

	result, cmd := updateModel(t, m, composeDoneMsg{project: "myapp"})

	assert.False(t, result.showSpinner)
	assert.Nil(t, result.composeCancel)
	assert.Equal(t, "[myapp] compose done.", result.statusMessage)
	assert.NotNil(t, cmd)
	_ = cancelCalled
}

func TestHandleComposeDoneMsg_withError(t *testing.T) {
	m := InitialModel()
	m.showSpinner = true

	result, _ := updateModel(t, m, composeDoneMsg{project: "myapp", err: errors.New("exit 1")})

	assert.False(t, result.showSpinner)
	assert.Contains(t, result.statusMessage, "[myapp] compose failed")
	assert.Contains(t, result.statusMessage, "exit 1")
}

func TestHandleWindowSizeMsg_setsWidthHeight(t *testing.T) {
	m := InitialModel()

	result, cmd := updateModel(t, m, tea.WindowSizeMsg{Width: 120, Height: 40})

	assert.Equal(t, 120, result.width)
	assert.Equal(t, 40, result.height)
	assert.Nil(t, cmd)
}

func TestHandleWindowSizeMsg_setsListTableWidths(t *testing.T) {
	m := InitialModel()

	result, _ := updateModel(t, m, tea.WindowSizeMsg{Width: 120, Height: 40})

	assert.Equal(t, 120, result.imageTable.Width())
	assert.Equal(t, 120, result.volumeTable.Width())
	assert.Equal(t, 120, result.networkTable.Width())
}

func TestHandleWindowSizeMsg_setsBubblesTableHeaderWidths(t *testing.T) {
	m := InitialModel()

	result, _ := updateModel(t, m, tea.WindowSizeMsg{Width: 120, Height: 40})

	for name, view := range map[string]string{
		"images":   result.imageTable.View(),
		"volumes":  result.volumeTable.View(),
		"networks": result.networkTable.View(),
	} {
		t.Run(name, func(t *testing.T) {
			lines := strings.Split(view, "\n")

			assert.GreaterOrEqual(t, len(lines), 2)
			assert.Equal(t, 120, lipgloss.Width(lines[0]))
			assert.Equal(t, 120, lipgloss.Width(lines[1]))
		})
	}
}

func TestHandleWindowSizeMsg_growsFlexibleColumns(t *testing.T) {
	m := InitialModel()

	narrow, _ := updateModel(t, m, tea.WindowSizeMsg{Width: 120, Height: 40})
	wide, _ := updateModel(t, narrow, tea.WindowSizeMsg{Width: 180, Height: 40})

	assert.Greater(t, columnWidthByTitle(t, wide.imageTable.Columns(), "Repository"), columnWidthByTitle(t, narrow.imageTable.Columns(), "Repository"))
	assert.Greater(t, columnWidthByTitle(t, wide.volumeTable.Columns(), "Mountpoint"), columnWidthByTitle(t, narrow.volumeTable.Columns(), "Mountpoint"))
	assert.Greater(t, columnWidthByTitle(t, wide.networkTable.Columns(), "Name"), columnWidthByTitle(t, narrow.networkTable.Columns(), "Name"))
}

func TestHandleWindowSizeMsg_narrowTablesKeepPositiveColumns(t *testing.T) {
	m := InitialModel()

	result, _ := updateModel(t, m, tea.WindowSizeMsg{Width: 40, Height: 20})

	for _, cols := range [][]table.Column{
		result.imageTable.Columns(),
		result.volumeTable.Columns(),
		result.networkTable.Columns(),
	} {
		for _, col := range cols {
			assert.Positive(t, col.Width)
		}
	}
}

func columnWidthByTitle(t *testing.T, cols []table.Column, title string) int {
	t.Helper()
	for _, col := range cols {
		if col.Title == title {
			return col.Width
		}
	}
	t.Fatalf("column %q not found", title)
	return 0
}

func TestHandleLogChunkMsg_appendsLine(t *testing.T) {
	m := InitialModel()

	result, cmd := updateModel(t, m, logChunkMsg("2024-01-01 INFO started"))

	assert.Len(t, result.logLines, 1)
	assert.Equal(t, "2024-01-01 INFO started", result.logLines[0])
	assert.Nil(t, cmd)
}

func TestHandleLogStreamDoneMsg_clearsChannel(t *testing.T) {
	m := InitialModel()
	m.logCh = make(chan string)
	cancelCalled := false
	m.logCancel = func() { cancelCalled = true }

	result, cmd := updateModel(t, m, logStreamDoneMsg{})

	assert.Nil(t, result.logCh)
	assert.Nil(t, result.logCancel)
	assert.Nil(t, cmd)
	_ = cancelCalled
}
