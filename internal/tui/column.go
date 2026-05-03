package tui

import (
	"strings"

	"charm.land/bubbles/v2/table"
	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/x/ansi"
)

// AlignType controls horizontal text alignment within a column.
type AlignType int

const (
	AlignLeft AlignType = iota
	AlignRight
)

// Column describes a column's layout properties and its computed width.
// Width is zero in the package-level specs; set by BuildColumns.
type Column struct {
	Header   string
	MinWidth int // flexible: expands with terminal (Fixed must be 0)
	Fixed    int // fixed: exact char width (MinWidth must be 0)
	Align    AlignType
	Width    int // computed by BuildColumns
}

// containerCols defines the canonical column specs for the containers table.
var containerCols = []Column{
	{Header: "Name", MinWidth: 20, Align: AlignLeft},
	{Header: "Status", Fixed: 14, Align: AlignLeft},
	{Header: "Image", MinWidth: 20, Align: AlignLeft},
	{Header: "Ports", Fixed: 18, Align: AlignLeft},
	{Header: "CPU%", Fixed: 6, Align: AlignRight},
	{Header: "Mem", Fixed: 10, Align: AlignRight},
	{Header: "Age", Fixed: 6, Align: AlignRight},
}

var imageCols = []Column{
	{Header: "ID", Fixed: 14, Align: AlignLeft},
	{Header: "Repository", MinWidth: 30, Align: AlignLeft},
	{Header: "Tag", MinWidth: 20, Align: AlignLeft},
	{Header: "Size", Fixed: 12, Align: AlignRight},
	{Header: "Created", Fixed: 12, Align: AlignRight},
}

var volumeCols = []Column{
	{Header: "Name", MinWidth: 30, Align: AlignLeft},
	{Header: "Driver", Fixed: 12, Align: AlignLeft},
	{Header: "Scope", Fixed: 10, Align: AlignLeft},
	{Header: "Mountpoint", MinWidth: 60, Align: AlignLeft},
}

var networkCols = []Column{
	{Header: "ID", MinWidth: 20, Align: AlignLeft},
	{Header: "Name", MinWidth: 30, Align: AlignLeft},
	{Header: "Driver", Fixed: 12, Align: AlignLeft},
	{Header: "Scope", Fixed: 10, Align: AlignLeft},
}

// BuildColumns returns a copy of specs with Width computed for the given
// terminal width. Fixed columns keep their exact size; flexible columns share
// the remaining budget proportionally. Any rounding remainder goes to the last
// flexible column so the table always fills 100% of the available width.
func BuildColumns(width int, specs []Column) []Column {
	fixedSum, totalMinFlex := 0, 0
	for _, c := range specs {
		if c.Fixed > 0 {
			fixedSum += c.Fixed
		} else {
			totalMinFlex += c.MinWidth
		}
	}

	// Bubble Tea table adds 1-space padding on each side per column.
	paddingOverhead := len(specs) * 2
	flexBudget := width - fixedSum - paddingOverhead
	if flexBudget < totalMinFlex {
		flexBudget = totalMinFlex
	}

	cols := make([]Column, len(specs))
	copy(cols, specs)

	assigned := 0
	lastFlexIdx := -1
	for i, c := range cols {
		if c.Fixed > 0 {
			cols[i].Width = c.Fixed
		} else {
			w := flexBudget * c.MinWidth / totalMinFlex
			if w < c.MinWidth {
				w = c.MinWidth
			}
			cols[i].Width = w
			assigned += w
			lastFlexIdx = i
		}
	}
	// Give rounding remainder to last flex column → fills 100% width.
	if lastFlexIdx >= 0 {
		remainder := flexBudget - assigned
		if remainder > 0 {
			cols[lastFlexIdx].Width += remainder
		}
	}
	return cols
}

func tableColumns(width int, specs []Column) []table.Column {
	cols := BuildColumns(width, specs)
	tableCols := make([]table.Column, len(cols))
	for i, col := range cols {
		tableCols[i] = table.Column{Title: col.Header, Width: col.Width}
	}
	return tableCols
}

// renderCell truncates to width (ANSI-safe) then applies lipgloss padding/alignment.
// Truncate happens before styling so ANSI escape codes from styled values (e.g.
// FormatStatus) are measured correctly and not re-counted by the style engine.
func renderCell(value string, width int, align AlignType) string {
	value = ansi.Truncate(value, width, "…")

	style := lipgloss.NewStyle().Width(width)
	if align == AlignRight {
		style = style.Align(lipgloss.Right)
	} else {
		style = style.Align(lipgloss.Left)
	}
	return style.Render(value)
}

// RenderRow builds a row of pre-styled cells using renderCell for each column.
func RenderRow(cols []Column, values []string) []string {
	row := make([]string, len(cols))
	for i, v := range values {
		row[i] = renderCell(v, cols[i].Width, cols[i].Align)
	}
	return row
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
