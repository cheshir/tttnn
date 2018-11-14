package ai

import "github.com/pkg/errors"

// Describes game board cell address like h8, j7 etc.
type TableCellIndex int

// Returns x and y position for this cell on the table.
func (index TableCellIndex) ToMatrixCoords() (int, int, error) {
	if index < 0 || index >= Width*Height {
		return 0, 0, errors.Errorf("table cell index %v is out of board. Available indexes range 0 - %v", index, Width*Height-1)
	}

	row := index / Width
	column := index % Width

	return int(row), int(column), nil
}
