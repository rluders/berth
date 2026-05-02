package tui

import (
	"github.com/rluders/berth/internal/controller"
)

type rowKind int

const (
	rowKindGroup     rowKind = iota // compose project header row
	rowKindContainer                // individual container row
)

// containerRowMeta is a parallel entry for every visible table row.
// It carries enough context to dispatch keyboard actions without re-parsing display strings.
type containerRowMeta struct {
	kind          rowKind
	groupName     string // group rows: project name; container rows: parent project (empty if standalone)
	containerID   string // container rows only
	containerName string // container rows only
}

type composeGroup struct {
	project    string
	containers []controller.Container
}

// buildComposeGroups partitions containers into ordered compose groups and a
// standalone slice (containers without the compose project label).
func buildComposeGroups(containers []controller.Container) (groups []composeGroup, standalone []controller.Container) {
	projectIndex := map[string]int{}
	for _, c := range containers {
		project := c.Labels["com.docker.compose.project"]
		if project == "" {
			standalone = append(standalone, c)
			continue
		}
		if idx, ok := projectIndex[project]; ok {
			groups[idx].containers = append(groups[idx].containers, c)
		} else {
			projectIndex[project] = len(groups)
			groups = append(groups, composeGroup{project: project, containers: []controller.Container{c}})
		}
	}
	return
}

// groupAggStatus counts running vs total containers in a group.
func groupAggStatus(containers []controller.Container) (running, total int) {
	for _, c := range containers {
		total++
		if c.State == "running" {
			running++
		}
	}
	return
}

// findGroupContainers returns all containers belonging to the given compose project.
func findGroupContainers(containers []controller.Container, project string) []controller.Container {
	var result []controller.Container
	for _, c := range containers {
		if c.Labels["com.docker.compose.project"] == project {
			result = append(result, c)
		}
	}
	return result
}
