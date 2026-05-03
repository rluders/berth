package tui

import "github.com/charmbracelet/bubbles/key"

// GlobalKeys holds key bindings available in all views.
type GlobalKeys struct {
	Quit key.Binding
	Help key.Binding
	Back key.Binding
	Tab1 key.Binding
	Tab2 key.Binding
	Tab3 key.Binding
	Tab4 key.Binding
	Tab5    key.Binding
	TabNext key.Binding
	TabPrev key.Binding
}

// ContainerKeys holds key bindings for the containers view.
type ContainerKeys struct {
	Details key.Binding
	Start   key.Binding
	Stop     key.Binding
	Restart  key.Binding
	Delete   key.Binding
	Logs     key.Binding
	Inspect  key.Binding
	Exec     key.Binding
	Filter   key.Binding
	Group    key.Binding
	Expand   key.Binding
	Collapse key.Binding
}

// ComposeKeys holds key bindings for compose project-level actions.
type ComposeKeys struct {
	Up       key.Binding
	UpBuild  key.Binding
	Recreate key.Binding
	Down     key.Binding
	Pull     key.Binding
	Build    key.Binding
}

// ImageKeys holds key bindings for the images view.
type ImageKeys struct {
	Delete key.Binding
	Prune  key.Binding
	Filter key.Binding
}

// VolumeKeys holds key bindings for the volumes view.
type VolumeKeys struct {
	Delete key.Binding
	Filter key.Binding
}

// NetworkKeys holds key bindings for the networks view.
type NetworkKeys struct {
	Inspect key.Binding
}

// SystemKeys holds key bindings for the system view.
type SystemKeys struct {
	BasicCleanup    key.Binding
	AdvancedCleanup key.Binding
	TotalCleanup    key.Binding
}

// LogsKeys holds key bindings for the logs view.
type LogsKeys struct {
	Pause       key.Binding
	Follow      key.Binding
	LineNumbers key.Binding
}

// ConfirmKeys holds key bindings for the confirm dialog.
type ConfirmKeys struct {
	Yes key.Binding
}

// FilterKeys holds key bindings for the filter input.
type FilterKeys struct {
	Submit key.Binding
	Cancel key.Binding
}

// Keys is the global key binding registry.
var Keys = struct {
	Global    GlobalKeys
	Container ContainerKeys
	Compose   ComposeKeys
	Image     ImageKeys
	Volume    VolumeKeys
	Network   NetworkKeys
	System    SystemKeys
	Logs      LogsKeys
	Confirm   ConfirmKeys
	Filter    FilterKeys
}{
	Global: GlobalKeys{
		Quit: key.NewBinding(
			key.WithKeys("ctrl+c"),
			key.WithHelp("ctrl+c", "quit"),
		),
		Help: key.NewBinding(
			key.WithKeys("?"),
			key.WithHelp("?", "help"),
		),
		Back: key.NewBinding(
			key.WithKeys("q", "esc"),
			key.WithHelp("q/esc", "back/quit"),
		),
		Tab1: key.NewBinding(
			key.WithKeys("1"),
			key.WithHelp("1", "containers"),
		),
		Tab2: key.NewBinding(
			key.WithKeys("2"),
			key.WithHelp("2", "images"),
		),
		Tab3: key.NewBinding(
			key.WithKeys("3"),
			key.WithHelp("3", "volumes"),
		),
		Tab4: key.NewBinding(
			key.WithKeys("4"),
			key.WithHelp("4", "networks"),
		),
		Tab5: key.NewBinding(
			key.WithKeys("5"),
			key.WithHelp("5", "system"),
		),
		TabNext: key.NewBinding(
			key.WithKeys("tab"),
			key.WithHelp("tab", "next tab"),
		),
		TabPrev: key.NewBinding(
			key.WithKeys("shift+tab"),
			key.WithHelp("shift+tab", "prev tab"),
		),
	},
	Container: ContainerKeys{
		Details: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "details"),
		),
		Start: key.NewBinding(
			key.WithKeys("s"),
			key.WithHelp("s", "start"),
		),
		Stop: key.NewBinding(
			key.WithKeys("x"),
			key.WithHelp("x", "stop"),
		),
		Restart: key.NewBinding(
			key.WithKeys("r"),
			key.WithHelp("r", "restart"),
		),
		Delete: key.NewBinding(
			key.WithKeys("d"),
			key.WithHelp("d", "delete"),
		),
		Logs: key.NewBinding(
			key.WithKeys("l"),
			key.WithHelp("l", "logs"),
		),
		Inspect: key.NewBinding(
			key.WithKeys("i"),
			key.WithHelp("i", "inspect"),
		),
		Exec: key.NewBinding(
			key.WithKeys("e"),
			key.WithHelp("e", "exec shell"),
		),
		Filter: key.NewBinding(
			key.WithKeys("/"),
			key.WithHelp("/", "filter"),
		),
		Group: key.NewBinding(
			key.WithKeys("g"),
			key.WithHelp("g", "group by compose"),
		),
		Expand: key.NewBinding(
			key.WithKeys("right"),
			key.WithHelp("→", "expand group"),
		),
		Collapse: key.NewBinding(
			key.WithKeys("left"),
			key.WithHelp("←", "collapse group"),
		),
	},
	Compose: ComposeKeys{
		Up: key.NewBinding(
			key.WithKeys("u"),
			key.WithHelp("u", "compose up"),
		),
		UpBuild: key.NewBinding(
			key.WithKeys("U"),
			key.WithHelp("U", "compose up --build"),
		),
		Recreate: key.NewBinding(
			key.WithKeys("R"),
			key.WithHelp("R", "compose recreate"),
		),
		Down: key.NewBinding(
			key.WithKeys("d"),
			key.WithHelp("d", "compose down"),
		),
		Pull: key.NewBinding(
			key.WithKeys("p"),
			key.WithHelp("p", "compose pull"),
		),
		Build: key.NewBinding(
			key.WithKeys("b"),
			key.WithHelp("b", "compose build"),
		),
	},
	Image: ImageKeys{
		Delete: key.NewBinding(
			key.WithKeys("d"),
			key.WithHelp("d", "remove"),
		),
		Prune: key.NewBinding(
			key.WithKeys("P"),
			key.WithHelp("P", "prune dangling"),
		),
		Filter: key.NewBinding(
			key.WithKeys("/"),
			key.WithHelp("/", "filter"),
		),
	},
	Volume: VolumeKeys{
		Delete: key.NewBinding(
			key.WithKeys("d"),
			key.WithHelp("d", "remove"),
		),
		Filter: key.NewBinding(
			key.WithKeys("/"),
			key.WithHelp("/", "filter"),
		),
	},
	Network: NetworkKeys{
		Inspect: key.NewBinding(
			key.WithKeys("i"),
			key.WithHelp("i", "inspect"),
		),
	},
	System: SystemKeys{
		BasicCleanup: key.NewBinding(
			key.WithKeys("b"),
			key.WithHelp("b", "basic cleanup"),
		),
		AdvancedCleanup: key.NewBinding(
			key.WithKeys("a"),
			key.WithHelp("a", "advanced cleanup"),
		),
		TotalCleanup: key.NewBinding(
			key.WithKeys("t"),
			key.WithHelp("t", "total cleanup"),
		),
	},
	Logs: LogsKeys{
		Pause: key.NewBinding(
			key.WithKeys("p"),
			key.WithHelp("p", "pause"),
		),
		Follow: key.NewBinding(
			key.WithKeys("f"),
			key.WithHelp("f", "follow"),
		),
		LineNumbers: key.NewBinding(
			key.WithKeys("n"),
			key.WithHelp("n", "line numbers"),
		),
	},
	Confirm: ConfirmKeys{
		Yes: key.NewBinding(
			key.WithKeys("y", "Y"),
			key.WithHelp("y", "confirm"),
		),
	},
	Filter: FilterKeys{
		Submit: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "apply filter"),
		),
		Cancel: key.NewBinding(
			key.WithKeys("esc"),
			key.WithHelp("esc", "clear filter"),
		),
	},
}
