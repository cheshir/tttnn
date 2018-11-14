package game

import (
	"strconv"

	"github.com/pkg/errors"
)

const (
	Empty Turn = iota
	Black
	White

	Undefined Result = iota
	BlackWin
	BlackLose
	Draw
	Invalid

	defaultWidth  = 15
	defaultHeight = 15

	neededForWin = 5
)

type Turn int

type Result int

func (r Result) String() string {
	switch r {
	case Undefined:
		return "Undefined"
	case BlackWin:
		return "Black win"
	case BlackLose:
		return "White win"
	case Draw:
		return "Draw"
	}

	return "Invalid"
}

var ErrFinished = errors.New("the game is already finished")
var ErrInvalidCell = errors.New("cell has an invalid formatted")
var ErrCellIsAlreadyFilled = errors.New("cell is already filled")

func New() *Table {
	return &Table{
		width:      defaultWidth,
		height:     defaultHeight,
		cells:      make([]Turn, defaultWidth*defaultHeight),
		free:       defaultWidth * defaultHeight,
		activeSide: Black,
		result:     Undefined,
	}
}

type Table struct {
	width      int
	height     int
	cells      []Turn
	free       int
	activeSide Turn
	result     Result
}

func (t *Table) Move(cell string) (Result, error) {
	if t.result != Undefined {
		return t.result, ErrFinished
	}

	index, err := t.cellToIndex(cell)
	if err != nil {
		return Undefined, err
	}

	if t.cells[index] != Empty {
		return Undefined, ErrCellIsAlreadyFilled
	}

	t.cells[index] = t.activeSide
	t.free--

	if t.isWinning(index) {
		if t.activeSide == Black {
			t.result = BlackWin
		} else {
			t.result = BlackLose
		}

		return t.result, nil
	}

	if t.free == 0 {
		t.result = Draw

		return t.result, nil
	}

	t.swapSides()

	return Undefined, nil
}

func (t *Table) cellToIndex(cell string) (int, error) {
	if len(cell) < 2 || len(cell) > 3 {
		return 0, ErrInvalidCell
	}

	row := cell[0] - 'a'
	column, err := strconv.Atoi(cell[1:])
	if err != nil || row < 0 || row >= 15 || column < 0 || column >= 15 {
		return 0, ErrInvalidCell
	}

	index := int(row)*t.width + column - 1

	return index, nil
}

func (t *Table) swapSides() {
	t.activeSide ^= 3 // 11 in binary. Works perfect with 1 (01) – black and 2 (10) – white.
}

func (t *Table) isWinning(currentMove int) bool {
	expected := t.cells[currentMove]

	horizontal := 1 + // Center.
		// Look left.
		t.calculateSiblings(
			expected,
			currentMove-1,
			func(i int) bool { return i >= (currentMove - currentMove%t.width) },
			func(i int) int { return i - 1 },
		) +
		// Look right.
		t.calculateSiblings(
			expected,
			currentMove+1,
			func(i int) bool { return i%t.width != 0 },
			func(i int) int { return i + 1 },
		)

	vertical := 1 + // Center.
		// Look top.
		t.calculateSiblings(
			expected,
			currentMove-t.width,
			func(i int) bool { return i >= 0 },
			func(i int) int { return i - t.width },
		) +
		// Look bottom.
		t.calculateSiblings(
			expected,
			currentMove+t.width,
			func(i int) bool { return i < len(t.cells) },
			func(i int) int { return i + t.width },
		)

	diagonal1 := 1 + // Center.
		// Look top.
		t.calculateSiblings(
			expected,
			currentMove-t.width-1,
			func(i int) bool { return i >= 0 && i%t.width != t.width-1 },
			func(i int) int { return i - t.width - 1 },
		) +
		// Look bottom.
		t.calculateSiblings(
			expected,
			currentMove+t.width+1,
			func(i int) bool { return i < len(t.cells) && i%t.width != 0 },
			func(i int) int { return i + t.width + 1 },
		)

	diagonal2 := 1 + // Center.
		// Look top.
		t.calculateSiblings(
			expected,
			currentMove-t.width+1,
			func(i int) bool { return i >= 0 && i%t.width != 0 },
			func(i int) int { return i - t.width + 1 },
		) +
		// Look bottom.
		t.calculateSiblings(
			expected,
			currentMove+t.width-1,
			func(i int) bool { return i < len(t.cells) && i%t.width != t.width-1 },
			func(i int) int { return i + t.width - 1 },
		)

	for _, dimensionSize := range []int{horizontal, vertical, diagonal1, diagonal2} {
		if dimensionSize >= neededForWin {
			return true
		}
	}

	return false
}

func (t *Table) calculateSiblings(expected Turn, start int, condition func(i int) bool, increment func(i int) int) int {
	count := 0

	for i := start; condition(i); i = increment(i) {
		if t.cells[i] != expected {
			break
		}
		count++
	}

	return count
}
