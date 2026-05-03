package tui

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBuildColumns_fillsWidthBudget(t *testing.T) {
	cols := BuildColumns(121, imageCols)

	assert.Equal(t, 121, paddedColumnWidth(cols))
}

func TestBuildColumns_fillsWidthBudgetForAllTableSpecs(t *testing.T) {
	for name, specs := range map[string][]Column{
		"containers": containerCols,
		"images":     imageCols,
		"volumes":    volumeCols,
		"networks":   networkCols,
	} {
		t.Run(name, func(t *testing.T) {
			cols := BuildColumns(140, specs)

			assert.Equal(t, 140, paddedColumnWidth(cols))
		})
	}
}

func TestTableColumns_preservesTitlesAndPositiveWidths(t *testing.T) {
	cols := tableColumns(120, volumeCols)

	require.Len(t, cols, len(volumeCols))
	for i, col := range cols {
		assert.Equal(t, volumeCols[i].Header, col.Title)
		assert.Positive(t, col.Width)
	}
}

func paddedColumnWidth(cols []Column) int {
	total := len(cols) * 2
	for _, col := range cols {
		total += col.Width
	}
	return total
}
