package tui

import (
	"github.com/rluders/berth/internal/controller"
)

// ViewType represents the different views in the TUI.
type ViewType int

const (
	ContainersView ViewType = iota
	ImagesView
	VolumesView
	NetworksView
	SystemView
	InspectView
	LogsView
)

// Typed message types for the Update dispatcher.
type (
	containerListMsg []controller.Container
	imageListMsg     []controller.Image
	volumeListMsg    []controller.Volume
	networkListMsg   []controller.Network
	systemInfoMsg    controller.SystemInfo
	logsMsg          string
	inspectMsg       string
	statusMsg        string
	errMsg           struct{ err error }
)

func (e errMsg) Error() string { return e.err.Error() }
