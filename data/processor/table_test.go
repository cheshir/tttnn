package main

import "testing"

func TestCellToIndex(t *testing.T) {
	tt := []struct {
		input    string
		expected int
		isError  bool
	}{
		{"a1", 0, false},
		{"h8", 112, false},
		{"a22", 0, true},
		{"z10", 0, true},
		{"h310", 0, true},
	}

	table := NewTable()

	for _, tc := range tt {
		index, err := table.cellToIndex(tc.input)
		if index != tc.expected {
			t.Errorf("expected index %d is not equal to actual %d", tc.expected, index)
		}

		// TODO Simplify me.
		if (err == nil && tc.isError) || (err != nil && !tc.isError) {
			t.Errorf("Expected error %v is not equal to actual %v", tc.isError, err)
		}
	}
}
