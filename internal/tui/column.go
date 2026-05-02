package tui

import (
	"strings"

	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/x/ansi"
)

// AlignType controls horizontal text alignment within a column.
type AlignType int

const (
	AlignLeft AlignType = iota
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
	{Header: "Name", MinWidth: 20, Align: AlignLeft},
	{Header: "Status", Fixed: 14, Align: AlignLeft},
	{Header: "Image", MinWidth: 30, Align: AlignLeft},
	{Header: "Ports", Fixed: 18, Align: AlignLeft},
	{Header: "CPU%", Fixed: 6, Align: AlignRight},
	{Header: "Mem", Fixed: 10, Align: AlignRight},
	{Header: "Age", Fixed: 6, Align: AlignRight},
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

// RenderRow builds a table row with ANSI-safe truncation and alignment.
//
// The bubbles table internally applies runewidth.Truncate to each cell value.
// runewidth counts ANSI escape sequences as visible characters, so pre-rendering
// with lipgloss Width inflates the runewidth beyond the column width and causes
// the table to truncate — destroying the ANSI prefix and making styled text invisible.
//
// Fix: pass values without lipgloss pre-padding. The table's own lipgloss.Width
// (which is charmbracelet/x/ansi-aware) handles padding correctly after truncation.
// We only apply ANSI-safe ansi.Truncate here to cap content before the table sees it.
// Right-aligned columns get plain-text padding (no ANSI) so runewidth stays accurate.
func RenderRow(cols []ColSpec, widths []int, values []string) table.Row {
	row := make(table.Row, len(cols))
	for i, v := range values {
		w := widths[i]
		v = ansi.Truncate(v, w, "…")
		if cols[i].Align == AlignRight {
			vw := ansi.StringWidth(v)
			if vw < w {
				v = strings.Repeat(" ", w-vw) + v
			}
		}
		row[i] = v
	}
	return row
}

// buildTableColumns builds Bubble Tea table columns with plain-padded headers.
// Headers are ASCII so plain string padding avoids ANSI inflation in runewidth.
func buildTableColumns(widths []int, cols []ColSpec) []table.Column {
	result := make([]table.Column, len(cols))
	for i, col := range cols {
		result[i] = table.Column{
			Title: padHeader(col.Header, widths[i], col.Align),
			Width: widths[i],
		}
	}
	return result
}

// padHeader returns a plain-text header string padded to exactly width chars.
// Headers are always ASCII so len() == visible width.
func padHeader(value string, width int, align AlignType) string {
	vw := len(value)
	if vw >= width {
		return value[:width]
	}
	pad := strings.Repeat(" ", width-vw)
	if align == AlignRight {
		return pad + value
	}
	return value + pad
}
