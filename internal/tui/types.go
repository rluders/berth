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
	DetailsView
)

// progressMsg drives the progress bar for long operations.
type progressMsg struct {
	percent float64
	label   string
	done    bool
}

// progressTickMsg animates the progress bar while an operation runs.
type progressTickMsg struct{}

// Typed message types for the Update dispatcher.
type (
	containerListMsg  []controller.Container
	imageListMsg      []controller.Image
	volumeListMsg     []controller.Volume
	networkListMsg    []controller.Network
	systemInfoMsg     controller.SystemInfo
	logsMsg           string
	logChunkMsg       string
	logStreamDoneMsg  struct{}
	inspectMsg        string
	detailsMsg        controller.ContainerDetails
	containerStatsMsg map[string]controller.ContainerStat
	statsTickMsg      struct{}
	refreshTickMsg    struct{}
	statusMsg         string
	errMsg            struct{ err error }
)

func (e errMsg) Error() string { return e.err.Error() }
