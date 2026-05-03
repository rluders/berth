package tui

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"time"

	tea "charm.land/bubbletea/v2"
	"github.com/rluders/berth/internal/controller"
)

// ── Fetch commands ────────────────────────────────────────────────────────────

func fetchContainersCmd() tea.Cmd {
	return func() tea.Msg {
		slog.Debug("fetchContainersCmd called")
		containers, err := controller.ListContainers()
		if err != nil {
			slog.Error("fetchContainersCmd error", "error", err)
			return errMsg{err}
		}
		return containerListMsg(containers)
	}
}

func fetchImagesCmd() tea.Cmd {
	return func() tea.Msg {
		slog.Debug("fetchImagesCmd called")
		images, err := controller.ListImages()
		if err != nil {
			slog.Error("fetchImagesCmd error", "error", err)
			return errMsg{err}
		}
		return imageListMsg(images)
	}
}

func fetchVolumesCmd() tea.Cmd {
	return func() tea.Msg {
		slog.Debug("fetchVolumesCmd called")
		volumes, err := controller.ListVolumes()
		if err != nil {
			slog.Error("fetchVolumesCmd error", "error", err)
			return errMsg{err}
		}
		return volumeListMsg(volumes)
	}
}

func fetchNetworksCmd() tea.Cmd {
	return func() tea.Msg {
		slog.Debug("fetchNetworksCmd called")
		networks, err := controller.ListNetworks()
		if err != nil {
			slog.Error("fetchNetworksCmd error", "error", err)
			return errMsg{err}
		}
		return networkListMsg(networks)
	}
}

func fetchSystemInfoCmd() tea.Cmd {
	return func() tea.Msg {
		slog.Debug("fetchSystemInfoCmd called")
		info, err := controller.GetSystemInfo()
		if err != nil {
			slog.Error("fetchSystemInfoCmd error", "error", err)
			return errMsg{err}
		}
		return systemInfoMsg(info)
	}
}

func fetchAllCmd() tea.Cmd {
	return tea.Batch(
		fetchContainersCmd(),
		fetchImagesCmd(),
		fetchVolumesCmd(),
		fetchNetworksCmd(),
		fetchSystemInfoCmd(),
	)
}

// ── Periodic tickers ──────────────────────────────────────────────────────────

func statsTickCmd() tea.Cmd {
	return tea.Tick(3*time.Second, func(time.Time) tea.Msg { return statsTickMsg{} })
}

func refreshTickCmd() tea.Cmd {
	return tea.Tick(5*time.Second, func(time.Time) tea.Msg { return refreshTickMsg{} })
}

// ── Container action commands ─────────────────────────────────────────────────

func startContainerCmd(idOrName string) tea.Cmd {
	return func() tea.Msg {
		slog.Debug("startContainerCmd", "id", idOrName)
		if err := controller.StartContainer(idOrName); err != nil {
			return errMsg{err}
		}
		return statusMsg(fmt.Sprintf("Container %s started.", idOrName))
	}
}

func stopContainerCmd(idOrName string) tea.Cmd {
	return func() tea.Msg {
		slog.Debug("stopContainerCmd", "id", idOrName)
		if err := controller.StopContainer(idOrName); err != nil {
			return errMsg{err}
		}
		return statusMsg(fmt.Sprintf("Container %s stopped.", idOrName))
	}
}

func restartContainerCmd(idOrName string) tea.Cmd {
	return func() tea.Msg {
		slog.Debug("restartContainerCmd", "id", idOrName)
		if err := controller.RestartContainer(idOrName); err != nil {
			return errMsg{err}
		}
		return statusMsg(fmt.Sprintf("Container %s restarted.", idOrName))
	}
}

func removeContainerCmd(idOrName string) tea.Cmd {
	return func() tea.Msg {
		slog.Debug("removeContainerCmd", "id", idOrName)
		if err := controller.RemoveContainer(idOrName); err != nil {
			return errMsg{err}
		}
		return statusMsg(fmt.Sprintf("Container %s removed.", idOrName))
	}
}

func fetchDetailsCmd(idOrName string) tea.Cmd {
	return func() tea.Msg {
		slog.Debug("fetchDetailsCmd", "id", idOrName)
		details, err := controller.GetContainerDetails(idOrName)
		if err != nil {
			return errMsg{err}
		}
		return detailsMsg(details)
	}
}

func inspectContainerCmd(idOrName string) tea.Cmd {
	return func() tea.Msg {
		slog.Debug("inspectContainerCmd", "id", idOrName)
		output, err := controller.InspectContainer(idOrName)
		if err != nil {
			return errMsg{err}
		}
		return inspectMsg(output)
	}
}

// ── Log streaming ─────────────────────────────────────────────────────────────

func startLogStreamCmd(id string) (chan string, context.CancelFunc, tea.Cmd) {
	ch := make(chan string, 500)
	ctx, cancel := context.WithCancel(context.Background())
	go controller.StreamContainerLogs(ctx, id, ch)
	return ch, cancel, waitForLogLineCmd(ch)
}

func startGroupLogStreamCmd(containers []controller.Container) (chan string, context.CancelFunc, tea.Cmd) {
	ch := make(chan string, 500)
	ctx, cancel := context.WithCancel(context.Background())
	entryCh := make(chan controller.LogEntry, 500)
	go controller.StreamMultiContainerLogs(ctx, containers, entryCh)
	go func() {
		defer close(ch)
		for entry := range entryCh {
			ch <- "[" + entry.ContainerName + "] " + entry.Line
		}
	}()
	return ch, cancel, waitForLogLineCmd(ch)
}

func waitForLogLineCmd(ch <-chan string) tea.Cmd {
	return func() tea.Msg {
		line, ok := <-ch
		if !ok {
			return logStreamDoneMsg{}
		}
		return logChunkMsg(line)
	}
}

// ── Stats ─────────────────────────────────────────────────────────────────────

func fetchStatsCmd(ids []string) tea.Cmd {
	return func() tea.Msg {
		result := make(map[string]controller.ContainerStat)
		for _, id := range ids {
			stat, err := controller.GetContainerStats(id)
			if err == nil {
				result[id] = stat
			}
		}
		return containerStatsMsg(result)
	}
}

// ── Image commands ────────────────────────────────────────────────────────────

func removeImageCmd(idOrName string) tea.Cmd {
	return func() tea.Msg {
		slog.Debug("removeImageCmd", "id", idOrName)
		if err := controller.RemoveImage(idOrName); err != nil {
			return errMsg{err}
		}
		return statusMsg(fmt.Sprintf("Image %s removed.", idOrName))
	}
}

func pruneImagesCmd() tea.Cmd {
	return func() tea.Msg {
		slog.Debug("pruneImagesCmd called")
		msg, err := controller.PruneImages()
		if err != nil {
			return errMsg{err}
		}
		return statusMsg(msg)
	}
}

// ── Volume commands ───────────────────────────────────────────────────────────

func removeVolumeCmd(name string) tea.Cmd {
	return func() tea.Msg {
		slog.Debug("removeVolumeCmd", "name", name)
		if err := controller.RemoveVolume(name); err != nil {
			return errMsg{err}
		}
		return statusMsg(fmt.Sprintf("Volume %s removed.", name))
	}
}

// ── Network commands ──────────────────────────────────────────────────────────

func inspectNetworkCmd(idOrName string) tea.Cmd {
	return func() tea.Msg {
		slog.Debug("inspectNetworkCmd", "id", idOrName)
		output, err := controller.InspectNetwork(idOrName)
		if err != nil {
			return errMsg{err}
		}
		return inspectMsg(output)
	}
}

// ── Progress tick ─────────────────────────────────────────────────────────────

func progressTickCmd() tea.Cmd {
	return tea.Tick(80*time.Millisecond, func(time.Time) tea.Msg { return progressTickMsg{} })
}

// ── System cleanup commands ───────────────────────────────────────────────────

func basicCleanupCmd() tea.Cmd {
	return func() tea.Msg {
		output, err := controller.BasicCleanup()
		if err != nil {
			return errMsg{err}
		}
		return progressMsg{percent: 1.0, label: "Basic cleanup: " + output, done: true}
	}
}

func advancedCleanupCmd() tea.Cmd {
	return func() tea.Msg {
		output, err := controller.AdvancedCleanup()
		if err != nil {
			return errMsg{err}
		}
		return progressMsg{percent: 1.0, label: "Advanced cleanup: " + output, done: true}
	}
}

func totalCleanupCmd() tea.Cmd {
	return func() tea.Msg {
		output, err := controller.TotalCleanup()
		if err != nil {
			return errMsg{err}
		}
		return progressMsg{percent: 1.0, label: "Total cleanup: " + output, done: true}
	}
}

// ── Compose commands ──────────────────────────────────────────────────────────

// composeStreamCmd starts a compose operation and returns the first streamed line as a message.
// Subsequent lines are self-scheduled via readNextComposeLineCmd.
func composeStreamCmd(ctx context.Context, project, workDir string, fn func(context.Context, string, string, chan<- string) error) tea.Cmd {
	return func() tea.Msg {
		ch := make(chan string, 64)
		if err := fn(ctx, project, workDir, ch); err != nil {
			return composeDoneMsg{project: project, err: err}
		}
		line, ok := <-ch
		if !ok {
			return composeDoneMsg{project: project}
		}
		return composeOutputMsg{project: project, line: line, ch: ch}
	}
}

// readNextComposeLineCmd reads the next line from an in-progress compose stream.
func readNextComposeLineCmd(ch <-chan string, project string) tea.Cmd {
	return func() tea.Msg {
		line, ok := <-ch
		if !ok {
			return composeDoneMsg{project: project}
		}
		return composeOutputMsg{project: project, line: line, ch: ch}
	}
}

func composeUpCmd(ctx context.Context, project, workDir string) tea.Cmd {
	return composeStreamCmd(ctx, project, workDir, controller.ComposeUp)
}

func composeUpBuildCmd(ctx context.Context, project, workDir string) tea.Cmd {
	return composeStreamCmd(ctx, project, workDir, controller.ComposeUpBuild)
}

func composeRecreateCmd(ctx context.Context, project, workDir string) tea.Cmd {
	return composeStreamCmd(ctx, project, workDir, controller.ComposeRecreate)
}

func composeDownCmd(ctx context.Context, project, workDir string) tea.Cmd {
	return composeStreamCmd(ctx, project, workDir, controller.ComposeDown)
}

func composePullCmd(ctx context.Context, project, workDir string) tea.Cmd {
	return composeStreamCmd(ctx, project, workDir, controller.ComposePull)
}

func composeBuildCmd(ctx context.Context, project, workDir string) tea.Cmd {
	return composeStreamCmd(ctx, project, workDir, controller.ComposeBuild)
}

// ── Exec shell ────────────────────────────────────────────────────────────────

func execShellCmd(containerID string) tea.Cmd {
	cmd := controller.ExecShell(containerID)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return tea.ExecProcess(cmd, func(err error) tea.Msg {
		if err != nil {
			return statusMsg("Exec ended: " + err.Error())
		}
		return statusMsg("Exec session ended.")
	})
}
