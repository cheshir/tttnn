package ai

import (
	tf "github.com/tensorflow/tensorflow/tensorflow/go"
)

// TODO check different values.
const (
	Width  = 15
	Height = 15

	EmptyMove  float32 = 0
	PlayerMove float32 = 0.5 // The same value used in model during learning.
	AIMove     float32 = 1.0
)

// Describes game board.
type Table [Width][Height]float32

func (t *Table) ToTensor() (*tf.Tensor, error) {
	view := [1][15][15][1]float32{}

	for row, columns := range t {
		for column, value := range columns {
			view[0][row][column][0] = value
		}
	}

	return tf.NewTensor(view)
}

func (t *Table) IsEmptyAtIndex(index TableCellIndex) bool {
	row, column, err := index.ToMatrixCoords()
	if err != nil {
		return false
	}

	return t[row][column] == EmptyMove
}
