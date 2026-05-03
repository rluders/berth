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
	{Header: "Name",   MinWidth: 20, Align: AlignLeft},
	{Header: "Status", Fixed:    12, Align: AlignLeft},
	{Header: "Image",  MinWidth: 30, Align: AlignLeft},
	{Header: "Ports",  Fixed:    18, Align: AlignLeft},
	{Header: "CPU%",   Fixed:     6, Align: AlignRight},
	{Header: "Mem",    Fixed:    10, Align: AlignRight},
	{Header: "Age",    Fixed:     6, Align: AlignRight},
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
func RenderRow(cols []Column, values []string) table.Row {
	row := make(table.Row, len(cols))
	for i, v := range values {
		w := cols[i].Width
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
func buildTableColumns(cols []Column) []table.Column {
	result := make([]table.Column, len(cols))
	for i, col := range cols {
		result[i] = table.Column{
			Title: padHeader(col.Header, col.Width, col.Align),
			Width: col.Width,
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
