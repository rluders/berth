package tui

import (
	"github.com/rluders/berth/internal/controller"
)

type RowType int

const (
	RowTypeGroup     RowType = iota // compose project header row
	RowTypeContainer                // individual container row
)

// Row is the canonical unit of the visible containers list.
type Row struct {
	Type       RowType
	GroupID    string               // project name; empty for standalone containers
	Name       string
	Collapsed  bool                 // group rows: current collapse state
	Containers []controller.Container // group rows: member containers
	Container  *controller.Container  // container rows: the container
}

type composeGroup struct {
	project    string
	containers []controller.Container
}

// BuildRows computes the flat visible row list from containers and collapse state.
func BuildRows(containers []controller.Container, collapsed map[string]bool) []Row {
	groups, standalone := buildComposeGroups(containers)

	var rows []Row
	for _, g := range groups {
		isCollapsed := collapsed[g.project]
		rows = append(rows, Row{
			Type:       RowTypeGroup,
			GroupID:    g.project,
			Name:       g.project,
			Collapsed:  isCollapsed,
			Containers: g.containers,
		})
		if !isCollapsed {
			for _, c := range g.containers {
				c := c
				rows = append(rows, Row{
					Type:      RowTypeContainer,
					GroupID:   g.project,
					Name:      c.Names,
					Container: &c,
				})
			}
		}
	}
	for _, c := range standalone {
		c := c
		rows = append(rows, Row{
			Type:      RowTypeContainer,
			GroupID:   "",
			Name:      c.Names,
			Container: &c,
		})
	}
	return rows
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
