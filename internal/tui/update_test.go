package tui

import (
	"errors"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/rluders/berth/internal/controller"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// updateModel is a test helper that calls Update and returns the concrete Model.
func updateModel(t *testing.T, m Model, msg tea.Msg) (Model, tea.Cmd) {
	t.Helper()
	newM, cmd := m.Update(msg)
	result, ok := newM.(Model)
	require.True(t, ok, "Update must return Model")
	return result, cmd
}

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
	// Spinner is zero value until Tick fires; just ensure no panic on Init.
	cmd := m.Init()
	assert.NotNil(t, cmd)
}

func TestUpdate_containerListMsg_setsContainers(t *testing.T) {
	m := InitialModel()
	containers := []controller.Container{{ID: "abc123", Names: "test", State: "running"}}

	result, cmd := updateModel(t, m, containerListMsg(containers))

	assert.Len(t, result.containers, 1)
	assert.Equal(t, "abc123", result.containers[0].ID)
	assert.False(t, result.showSpinner)
	assert.Empty(t, result.statusMessage)
	assert.Nil(t, cmd)
}

func TestUpdate_containerListMsg_empty(t *testing.T) {
	m := InitialModel()
	m.showSpinner = true

	result, _ := updateModel(t, m, containerListMsg(nil))

	assert.Empty(t, result.containers)
	assert.False(t, result.showSpinner)
}

func TestUpdate_imageListMsg_setsImages(t *testing.T) {
	m := InitialModel()
	images := []controller.Image{{ID: "img1", Repository: "nginx", Tag: "latest"}}

	result, cmd := updateModel(t, m, imageListMsg(images))

	assert.Len(t, result.images, 1)
	assert.False(t, result.showSpinner)
	assert.Nil(t, cmd)
}

func TestUpdate_errMsg_setsError(t *testing.T) {
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

func TestUpdate_composeOutputMsg_appendsToBuffer(t *testing.T) {
	m := InitialModel()
	ch := make(chan string, 1)
	close(ch)

	result, _ := updateModel(t, m, composeOutputMsg{project: "myapp", line: "Pulling image...", ch: ch})

	assert.Len(t, result.composeOutput, 1)
	assert.Equal(t, "Pulling image...", result.composeOutput[0])
	assert.Equal(t, "Pulling image...", result.statusMessage)
}

func TestUpdate_composeOutputMsg_rollingBuffer(t *testing.T) {
	m := InitialModel()
	// Pre-fill buffer to exactly 200 entries.
	m.composeOutput = make([]string, 200)
	ch := make(chan string, 1)
	close(ch)

	result, _ := updateModel(t, m, composeOutputMsg{project: "myapp", line: "new line", ch: ch})

	assert.Len(t, result.composeOutput, 200, "buffer must not exceed 200 entries")
	assert.Equal(t, "new line", result.composeOutput[199])
}

func TestUpdate_composeDoneMsg_success(t *testing.T) {
	m := InitialModel()
	m.showSpinner = true
	cancelCalled := false
	m.composeCancel = func() { cancelCalled = true }

	result, cmd := updateModel(t, m, composeDoneMsg{project: "myapp"})

	assert.False(t, result.showSpinner)
	assert.Nil(t, result.composeCancel)
	assert.Equal(t, "[myapp] compose done.", result.statusMessage)
	assert.NotNil(t, cmd) // fetchContainersCmd
	_ = cancelCalled      // cancel not called by handler — caller manages lifecycle
}

func TestUpdate_composeDoneMsg_withError(t *testing.T) {
	m := InitialModel()
	m.showSpinner = true

	result, _ := updateModel(t, m, composeDoneMsg{project: "myapp", err: errors.New("exit 1")})

	assert.False(t, result.showSpinner)
	assert.Contains(t, result.statusMessage, "[myapp] compose failed")
	assert.Contains(t, result.statusMessage, "exit 1")
}

func TestUpdate_windowSizeMsg_setsWidthHeight(t *testing.T) {
	m := InitialModel()

	result, cmd := updateModel(t, m, tea.WindowSizeMsg{Width: 120, Height: 40})

	assert.Equal(t, 120, result.width)
	assert.Equal(t, 40, result.height)
	assert.Nil(t, cmd)
}

func TestUpdate_logChunkMsg_appendsLine(t *testing.T) {
	m := InitialModel()

	result, cmd := updateModel(t, m, logChunkMsg("2024-01-01 INFO started"))

	assert.Len(t, result.logLines, 1)
	assert.Equal(t, "2024-01-01 INFO started", result.logLines[0])
	assert.Nil(t, cmd) // logCh is nil so no waitForLogLineCmd
}

func TestUpdate_logStreamDoneMsg_clearsChannel(t *testing.T) {
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
