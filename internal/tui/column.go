package tui

import (
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/lipgloss"
)

// AlignType controls horizontal text alignment within a column.
type AlignType int

const (
	AlignLeft  AlignType = iota
	AlignRight
)

// ColSpec describes a single column's layout properties.
type ColSpec struct {
	Header   string
	MinWidth int // flexible column: min chars, expands with terminal
	Fixed    int // fixed column: exact char width (MinWidth must be 0)
	Align    AlignType
}

// containerCols defines the canonical column specs for the containers table.
var containerCols = []ColSpec{
	{Header: "Name",   MinWidth: 20, Align: AlignLeft},
	{Header: "Status", Fixed: 14,   Align: AlignLeft},
	{Header: "Image",  MinWidth: 30, Align: AlignLeft},
	{Header: "Ports",  Fixed: 18,   Align: AlignLeft},
	{Header: "CPU%",   Fixed: 6,    Align: AlignRight},
	{Header: "Mem",    Fixed: 10,   Align: AlignRight},
	{Header: "Age",    Fixed: 6,    Align: AlignRight},
}

// computeWidths returns per-column pixel widths for the given available
// terminal width. Fixed columns keep their exact size; flexible columns
// (MinWidth > 0) share the remaining budget proportionally.
func computeWidths(tableWidth int, cols []ColSpec) []int {
	fixedSum := 0
	totalMinFlex := 0
	for _, c := range cols {
		if c.Fixed > 0 {
			fixedSum += c.Fixed
		} else {
			totalMinFlex += c.MinWidth
		}
	}

	// Bubble Tea table adds 1-space padding on each side per column.
	paddingOverhead := len(cols) * 2
	flexBudget := tableWidth - fixedSum - paddingOverhead
	if flexBudget < totalMinFlex {
		flexBudget = totalMinFlex
	}

	widths := make([]int, len(cols))
	for i, c := range cols {
		if c.Fixed > 0 {
			widths[i] = c.Fixed
		} else {
			w := flexBudget * c.MinWidth / totalMinFlex
			if w < c.MinWidth {
				w = c.MinWidth
			}
			widths[i] = w
		}
	}
	return widths
}

// renderCell formats value into exactly width printable characters with the
// requested alignment. ANSI escape codes in value are handled transparently.
func renderCell(value string, width int, align AlignType) string {
	style := lipgloss.NewStyle().Width(width).MaxWidth(width)
	if align == AlignRight {
		style = style.AlignHorizontal(lipgloss.Right)
	}
	return style.Render(value)
}

// RenderRow applies renderCell to each value using the matching ColSpec and width.
// Single source of truth for all table rows: container rows, group rows.
func RenderRow(cols []ColSpec, widths []int, values []string) table.Row {
	row := make(table.Row, len(cols))
	for i, v := range values {
		row[i] = renderCell(v, widths[i], cols[i].Align)
	}
	return row
}

// buildTableColumns builds Bubble Tea table columns whose header titles are
// pre-rendered to the same fixed width used for row cells.
func buildTableColumns(widths []int, cols []ColSpec) []table.Column {
	result := make([]table.Column, len(cols))
	for i, col := range cols {
		result[i] = table.Column{
			Title: renderCell(col.Header, widths[i], col.Align),
			Width: widths[i],
		}
	}
	return result
}
