package tui

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBuildColumns_fillsWidthBudget(t *testing.T) {
	cols := BuildColumns(121, imageCols)

	total := len(cols) * 2
	for _, col := range cols {
		total += col.Width
	}

	assert.Equal(t, 121, total)
}

func TestTableColumns_preservesTitlesAndPositiveWidths(t *testing.T) {
	cols := tableColumns(120, volumeCols)

	require.Len(t, cols, len(volumeCols))
	for i, col := range cols {
		assert.Equal(t, volumeCols[i].Header, col.Title)
		assert.Positive(t, col.Width)
	}
}
