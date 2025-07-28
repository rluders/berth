package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/rluders/container-tui/internal/controller"
	"github.com/rluders/container-tui/internal/engine"
)

type ViewType int

const (
	ContainersView ViewType = iota
	ImagesView
	VolumesView
	NetworksView
	SystemView
)

type model struct {
	engineType    engine.EngineType
	currentView   ViewType
	containerTable table.Model
	imageTable    table.Model
	volumeTable   table.Model
	networkTable  table.Model
	systemInfo    controller.SystemInfo
	err           error
	logs          string
	statusMessage string
}

func InitialModel() model {
	containerColumns := []table.Column{
		{Title: "ID", Width: 12},
		{Title: "Image", Width: 20},
		{Title: "Command", Width: 30},
		{Title: "Created", Width: 15},
		{Title: "Status", Width: 20},
		{Title: "Ports", Width: 20},
		{Title: "Names", Width: 20},
	}

	containerTable := table.New(
		table.WithColumns(containerColumns),
		table.WithFocused(true),
		table.WithHeight(10),
	)

	imageColumns := []table.Column{
		{Title: "ID", Width: 15},
		{Title: "Repository", Width: 30},
		{Title: "Tag", Width: 15},
		{Title: "Size", Width: 10},
		{Title: "Created", Width: 20},
	}

	imageTable := table.New(
		table.WithColumns(imageColumns),
		table.WithFocused(false),
		table.WithHeight(10),
	)

	volumeColumns := []table.Column{
		{Title: "Name", Width: 30},
		{Title: "Driver", Width: 15},
		{Title: "Scope", Width: 10},
		{Title: "Mountpoint", Width: 50},
	}

	volumeTable := table.New(
		table.WithColumns(volumeColumns),
		table.WithFocused(false),
		table.WithHeight(10),
	)

	networkColumns := []table.Column{
		{Title: "ID", Width: 15},
		{Title: "Name", Width: 30},
		{Title: "Driver", Width: 15},
		{Title: "Scope", Width: 10},
	}

	networkTable := table.New(
		table.WithColumns(networkColumns),
		table.WithFocused(false),
		table.WithHeight(10),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.BorderStyle(lipgloss.NormalBorder()).BorderForeground(lipgloss.Color("240")).BorderBottom(true).Bold(false)
	s.Selected = s.Selected.Foreground(lipgloss.Color("229")).Background(lipgloss.Color("57")).Bold(false)
	containerTable.SetStyles(s)
	imageTable.SetStyles(s)
	volumeTable.SetStyles(s)
	networkTable.SetStyles(s)

	return model{
		engineType:    engine.DetectEngine(),
		currentView:   ContainersView,
		containerTable: containerTable,
		imageTable:    imageTable,
		volumeTable:   volumeTable,
		networkTable:  networkTable,
		systemInfo:    controller.SystemInfo{}, // Initialize with empty SystemInfo
	}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(fetchContainersCmd(), fetchImagesCmd(), fetchVolumesCmd(), fetchNetworksCmd(), fetchSystemInfoCmd())
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "1":
			m.currentView = ContainersView
			return m, nil
		case "2":
			m.currentView = ImagesView
			return m, nil
		case "3":
			m.currentView = VolumesView
			return m, nil
		case "4":
			m.currentView = NetworksView
			return m, nil
		case "5":
            m.currentView = SystemView
            return m, nil
        }

        if m.currentView == ContainersView {
            switch msg.String() {
            case "s": // Start container
                if len(m.containerTable.SelectedRow()) > 0 {
                    containerID := m.containerTable.SelectedRow()[0]
                    return m, startContainerCmd(containerID)
                }
            case "x": // Stop container
                if len(m.containerTable.SelectedRow()) > 0 {
                    containerID := m.containerTable.SelectedRow()[0]
                    return m, stopContainerCmd(containerID)
                }
            case "d": // Remove container
                if len(m.containerTable.SelectedRow()) > 0 {
                    containerID := m.containerTable.SelectedRow()[0]
                    return m, removeContainerCmd(containerID)
                }
            case "l": // View logs
                if len(m.containerTable.SelectedRow()) > 0 {
                    containerID := m.containerTable.SelectedRow()[0]
                    return m, getLogsCmd(containerID)
                }
            }
        } else if m.currentView == ImagesView {
            switch msg.String() {
            case "d": // Remove image
                if len(m.imageTable.SelectedRow()) > 0 {
                    imageID := m.imageTable.SelectedRow()[0]
                    return m, removeImageCmd(imageID)
                }
            }
        } else if m.currentView == VolumesView {
            switch msg.String() {
            case "d": // Remove volume
                if len(m.volumeTable.SelectedRow()) > 0 {
                    volumeName := m.volumeTable.SelectedRow()[0]
                    return m, removeVolumeCmd(volumeName)
                }
            }
        } else if m.currentView == NetworksView {
            switch msg.String() {
            case "i": // Inspect network
                if len(m.networkTable.SelectedRow()) > 0 {
                    networkID := m.networkTable.SelectedRow()[0]
                    return m, inspectNetworkCmd(networkID)
                }
            }
        } else if m.currentView == SystemView {
            switch msg.String() {
            case "b": // Basic Cleanup
                return m, basicCleanupCmd()
            case "a": // Advanced Cleanup
                return m, advancedCleanupCmd()
            case "t": // Total Cleanup
                return m, totalCleanupCmd()
            }
        }

    case []controller.Container:
        rows := make([]table.Row, len(msg))
        for i, c := range msg {
            rows[i] = table.Row{c.ID, c.Image, c.Command, c.Created, c.Status, c.Ports, c.Names}
        }
        m.containerTable.SetRows(rows)
        m.statusMessage = ""
    case []controller.Image:
        rows := make([]table.Row, len(msg))
        for i, img := range msg {
            rows[i] = table.Row{img.ID, img.Repository, img.Tag, img.Size, img.Created}
        }
        m.imageTable.SetRows(rows)
        m.statusMessage = ""
    case []controller.Volume:
        rows := make([]table.Row, len(msg))
        for i, vol := range msg {
            rows[i] = table.Row{vol.Name, vol.Driver, vol.Scope, vol.Mountpoint}
        }
        m.volumeTable.SetRows(rows)
        m.statusMessage = ""
    case []controller.Network:
        rows := make([]table.Row, len(msg))
        for i, net := range msg {
            rows[i] = table.Row{net.ID, net.Name, net.Driver, net.Scope}
        }
        m.networkTable.SetRows(rows)
        m.statusMessage = ""
    case controller.SystemInfo:
        m.systemInfo = msg
        m.statusMessage = ""
    case string: // For logs or inspect output
        m.logs = msg
        m.statusMessage = ""
    case error:
        m.err = msg
        m.statusMessage = ""
    case statusMsg: // For status messages after actions
        m.statusMessage = string(msg)
        return m, tea.Batch(fetchContainersCmd(), fetchImagesCmd(), fetchVolumesCmd(), fetchNetworksCmd(), fetchSystemInfoCmd())
    }

    if m.currentView == ContainersView {
        m.containerTable, cmd = m.containerTable.Update(msg)
    } else if m.currentView == ImagesView {
        m.imageTable, cmd = m.imageTable.Update(msg)
    } else if m.currentView == VolumesView {
        m.volumeTable, cmd = m.volumeTable.Update(msg)
    } else if m.currentView == NetworksView {
        m.networkTable, cmd = m.networkTable.Update(msg)
    }
    return m, cmd
}

func (m model) View() string {
    if m.err != nil {
        return fmt.Sprintf("Error: %v", m.err)
    }

    sb := strings.Builder{}
    sb.WriteString(fmt.Sprintf("Detected Container Engine: %s\n\n", m.engineType))
    sb.WriteString("Press '1' for Containers, '2' for Images, '3' for Volumes, '4' for Networks, '5' for System Info\n\n")

    if m.currentView == ContainersView {
        sb.WriteString("Containers:\n")
        sb.WriteString(m.containerTable.View())
        sb.WriteString("\nPress 's' to start, 'x' to stop, 'd' to remove, 'l' for logs, 'q' to quit.")
    } else if m.currentView == ImagesView {
        sb.WriteString("Images:\n")
        sb.WriteString(m.imageTable.View())
        sb.WriteString("\nPress 'd' to remove, 'q' to quit.")
    } else if m.currentView == VolumesView {
        sb.WriteString("Volumes:\n")
        sb.WriteString(m.volumeTable.View())
        sb.WriteString("\nPress 'd' to remove, 'q' to quit.")
    } else if m.currentView == NetworksView {
        sb.WriteString("Networks:\n")
        sb.WriteString(m.networkTable.View())
        sb.WriteString("\nPress 'i' to inspect, 'q' to quit.")
    } else if m.currentView == SystemView {
        sb.WriteString("System Information:\n")
        sb.WriteString(fmt.Sprintf("  Containers: %d (Running: %d, Paused: %d, Stopped: %d)\n", m.systemInfo.Containers, m.systemInfo.Running, m.systemInfo.Paused, m.systemInfo.Stopped))
        sb.WriteString(fmt.Sprintf("  Images: %d\n", m.systemInfo.Images))
        sb.WriteString(fmt.Sprintf("  Volumes: %d\n", m.systemInfo.Volumes))
        sb.WriteString(fmt.Sprintf("  Networks: %d\n", m.systemInfo.Networks))
        sb.WriteString(fmt.Sprintf("  Disk Usage: %s\n", m.systemInfo.DiskUsage))
        sb.WriteString("\nPress 'q' to quit.")
    }

    if m.statusMessage != "" {
        sb.WriteString("\n" + m.statusMessage)
    }

    if m.logs != "" {
        sb.WriteString("\n\nLogs:\n" + m.logs)
    }

    return sb.String()
}

func fetchContainersCmd() tea.Cmd {
	return func() tea.Msg {
		containers, err := controller.ListContainers()
		if err != nil {
			return err
		}
		return containers
	}
}

func fetchImagesCmd() tea.Cmd {
	return func() tea.Msg {
		images, err := controller.ListImages()
		if err != nil {
			return err
		}
		return images
	}
}

func fetchVolumesCmd() tea.Cmd {
	return func() tea.Msg {
		volumes, err := controller.ListVolumes()
		if err != nil {
			return err
		}
		return volumes
	}
}

func fetchNetworksCmd() tea.Cmd {
	return func() tea.Msg {
		networks, err := controller.ListNetworks()
		if err != nil {
			return err
		}
		return networks
	}
}

func fetchSystemInfoCmd() tea.Cmd {
	return func() tea.Msg {
		systemInfo, err := controller.GetSystemInfo()
		if err != nil {
			return err
		}
		return systemInfo
	}
}

type statusMsg string

func startContainerCmd(idOrName string) tea.Cmd {
	return func() tea.Msg {
		err := controller.StartContainer(idOrName)
		if err != nil {
			return err
		}
		return statusMsg(fmt.Sprintf("Container %s started.", idOrName))
	}
}

func stopContainerCmd(idOrName string) tea.Cmd {
	return func() tea.Msg {
		err := controller.StopContainer(idOrName)
		if err != nil {
			return err
		}
		return statusMsg(fmt.Sprintf("Container %s stopped.", idOrName))
	}
}

func removeContainerCmd(idOrName string) tea.Cmd {
	return func() tea.Msg {
		err := controller.RemoveContainer(idOrName)
		if err != nil {
			return err
		}
		return statusMsg(fmt.Sprintf("Container %s removed.", idOrName))
	}
}

func getLogsCmd(idOrName string) tea.Cmd {
	return func() tea.Msg {
		logs, err := controller.GetContainerLogs(idOrName)
		if err != nil {
			return err
		}
		return logs
	}
}

func removeImageCmd(idOrName string) tea.Cmd {
	return func() tea.Msg {
		err := controller.RemoveImage(idOrName)
		if err != nil {
			return err
		}
		return statusMsg(fmt.Sprintf("Image %s removed.", idOrName))
	}
}

func removeVolumeCmd(name string) tea.Cmd {
	return func() tea.Msg {
		err := controller.RemoveVolume(name)
		if err != nil {
			return err
		}
		return statusMsg(fmt.Sprintf("Volume %s removed.", name))
	}
}

func inspectNetworkCmd(idOrName string) tea.Cmd {
	return func() tea.Msg {
		output, err := controller.InspectNetwork(idOrName)
		if err != nil {
			return err
		}
		return output
	}
}

func basicCleanupCmd() tea.Cmd {
	return func() tea.Msg {
		output, err := controller.BasicCleanup()
		if err != nil {
			return err
		}
		return statusMsg(fmt.Sprintf("Basic Cleanup completed:\n%s", output))
	}
}

func advancedCleanupCmd() tea.Cmd {
	return func() tea.Msg {
		output, err := controller.AdvancedCleanup()
		if err != nil {
			return err
		}
		return statusMsg(fmt.Sprintf("Advanced Cleanup completed:\n%s", output))
	}
}

func totalCleanupCmd() tea.Cmd {
	return func() tea.Msg {
		output, err := controller.TotalCleanup()
		if err != nil {
			return err
		}
		return statusMsg(fmt.Sprintf("Total Cleanup completed:\n%s", output))
	}
}

