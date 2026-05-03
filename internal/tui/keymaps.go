package tui

import "charm.land/bubbles/v2/key"

// containersKeyMap implements help.KeyMap for the containers view.
type containersKeyMap struct{}

func (containersKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{
		Keys.Container.Details,
		Keys.Container.Logs,
		Keys.Container.Start,
		Keys.Container.Stop,
		Keys.Container.Delete,
		Keys.Global.Help,
		Keys.Global.Quit,
	}
}

func (containersKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{Keys.Container.Details, Keys.Container.Logs, Keys.Container.Inspect, Keys.Container.Exec},
		{Keys.Container.Start, Keys.Container.Stop, Keys.Container.Restart, Keys.Container.Delete},
		{Keys.Container.Filter, Keys.Container.Group, Keys.Container.Expand, Keys.Container.Collapse},
		{Keys.Compose.Up, Keys.Compose.UpBuild, Keys.Compose.Recreate, Keys.Compose.Down},
		{Keys.Compose.Pull, Keys.Compose.Build},
		{Keys.Global.Tab1, Keys.Global.Tab2, Keys.Global.Tab3, Keys.Global.Tab4, Keys.Global.Tab5},
		{Keys.Global.Help, Keys.Global.Back},
	}
}

// imagesKeyMap implements help.KeyMap for the images view.
type imagesKeyMap struct{}

func (imagesKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{Keys.Image.Delete, Keys.Image.Prune, Keys.Image.Filter, Keys.Global.Help, Keys.Global.Quit}
}

func (imagesKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{Keys.Image.Delete, Keys.Image.Prune, Keys.Image.Filter},
		{Keys.Global.Tab1, Keys.Global.Tab2, Keys.Global.Tab3, Keys.Global.Tab4, Keys.Global.Tab5},
		{Keys.Global.Help, Keys.Global.Quit},
	}
}

// volumesKeyMap implements help.KeyMap for the volumes view.
type volumesKeyMap struct{}

func (volumesKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{Keys.Volume.Delete, Keys.Volume.Filter, Keys.Global.Help, Keys.Global.Quit}
}

func (volumesKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{Keys.Volume.Delete, Keys.Volume.Filter},
		{Keys.Global.Tab1, Keys.Global.Tab2, Keys.Global.Tab3, Keys.Global.Tab4, Keys.Global.Tab5},
		{Keys.Global.Help, Keys.Global.Quit},
	}
}

// networksKeyMap implements help.KeyMap for the networks view.
type networksKeyMap struct{}

func (networksKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{Keys.Network.Inspect, Keys.Global.Help, Keys.Global.Quit}
}

func (networksKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{Keys.Network.Inspect},
		{Keys.Global.Tab1, Keys.Global.Tab2, Keys.Global.Tab3, Keys.Global.Tab4, Keys.Global.Tab5},
		{Keys.Global.Help, Keys.Global.Quit},
	}
}

// systemKeyMap implements help.KeyMap for the system view.
type systemKeyMap struct{}

func (systemKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{Keys.System.BasicCleanup, Keys.System.AdvancedCleanup, Keys.System.TotalCleanup, Keys.Global.Quit}
}

func (systemKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{Keys.System.BasicCleanup, Keys.System.AdvancedCleanup, Keys.System.TotalCleanup},
		{Keys.Global.Help, Keys.Global.Quit},
	}
}

// logsKeyMap implements help.KeyMap for the logs view.
type logsKeyMap struct{}

func (logsKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{Keys.Logs.Pause, Keys.Logs.Follow, Keys.Global.Back}
}

func (logsKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{Keys.Logs.Pause, Keys.Logs.Follow},
		{Keys.Global.Back, Keys.Global.Help},
	}
}

// viewportKeyMap implements help.KeyMap for inspect/details views.
type viewportKeyMap struct{}

func (viewportKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{Keys.Global.Back, Keys.Global.Help}
}

func (viewportKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{Keys.Global.Back, Keys.Global.Help},
	}
}

// currentKeyMap returns the help.KeyMap for the active view.
func (m Model) currentKeyMap() interface {
	ShortHelp() []key.Binding
	FullHelp() [][]key.Binding
} {
	switch m.currentView {
	case ContainersView:
		return containersKeyMap{}
	case ImagesView:
		return imagesKeyMap{}
	case VolumesView:
		return volumesKeyMap{}
	case NetworksView:
		return networksKeyMap{}
	case SystemView:
		return systemKeyMap{}
	case LogsView:
		return logsKeyMap{}
	case InspectView, DetailsView:
		return viewportKeyMap{}
	}
	return containersKeyMap{}
}
