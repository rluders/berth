package tui

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
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

func TestUpdate_unknownMsgPreservesModel(t *testing.T) {
	m := InitialModel()

	result, cmd := updateModel(t, m, struct{}{})

	assert.Equal(t, m.currentView, result.currentView)
	assert.Nil(t, cmd)
}
