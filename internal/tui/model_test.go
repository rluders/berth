package tui

import (
	"testing"

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
