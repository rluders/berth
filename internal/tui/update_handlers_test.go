package tui

import (
	"errors"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
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
