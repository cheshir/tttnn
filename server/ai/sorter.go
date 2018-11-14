package ai

import "sort"

func NewSorter(list []float64) *sorter {
	s := &sorter{Interface: sort.Float64Slice(list), indexes: make([]int, len(list))}
	for i := range s.indexes {
		s.indexes[i] = i
	}

	return s
}

// Sort values and their indexes.
type sorter struct {
	sort.Interface
	indexes []int
}

func (s sorter) Swap(i, j int) {
	s.Interface.Swap(i, j)
	s.indexes[i], s.indexes[j] = s.indexes[j], s.indexes[i]
}
func (s sorter) Indexes() []int {
	return s.indexes
}
