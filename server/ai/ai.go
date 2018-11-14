package ai

import (
	"sort"

	"github.com/pkg/errors"

	tf "github.com/tensorflow/tensorflow/tensorflow/go"
)

var ErrNoAvailableMoves = errors.New("no available moves – the board is full")

type Config struct {
	ModelDir    string `envconfig:"model_dir"`
	InputLayer  string `envconfig:"input_layer"`
	OutputLayer string `envconfig:"output_layer"`
	Tags        []string
}

type AI struct {
	model       *tf.SavedModel
	inputLayer  string
	outputLayer string
}

func New(config Config) (*AI, error) {
	model, err := tf.LoadSavedModel(config.ModelDir, config.Tags, nil)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to load saved model")
	}

	ai := &AI{
		model:       model,
		inputLayer:  config.InputLayer,
		outputLayer: config.OutputLayer,
	}

	return ai, nil
}

// Returns predicted move index, it probability or error.
func (ai *AI) Predict(table *Table) (int, float64, error) {
	input, err := table.ToTensor()
	if err != nil {
		return 0, 0, errors.WithMessage(err, "failed to convert table to tensor")
	}

	output, err := ai.model.Session.Run(
		map[tf.Output]*tf.Tensor{
			ai.model.Graph.Operation(ai.inputLayer).Output(0): input,
		},
		[]tf.Output{
			ai.model.Graph.Operation(ai.outputLayer).Output(0),
		},
		nil,
	)

	if err != nil {
		return 0, 0, errors.Wrap(err, "failed to predict move")
	}

	return ai.getNextMove(table, output)
}

// Returns predicted move index, it probability or error.
func (ai *AI) getNextMove(table *Table, output []*tf.Tensor) (int, float64, error) {
	outputValue := output[0].Value()
	prediction := outputValue.([][]float32)[0] // Result is a vector with dimension 1 x 225.
	moves, probabilities := ai.sortMovesByProbabilities(prediction)

	for i, move := range moves {
		if table.IsEmptyAtIndex(TableCellIndex(move)) {
			return move, probabilities[i], nil
		}
	}

	return 0, 0, ErrNoAvailableMoves
}

// Sort probabilities and their indexes.
func (ai *AI) sortMovesByProbabilities(prediction []float32) ([]int, []float64) {
	// Sorter requires float64 values.
	probabilities := ai.convertSliceFromFloat32ToFloat64(prediction)
	sorter := NewSorter(probabilities)
	sort.Sort(sort.Reverse(sorter))

	return sorter.Indexes(), probabilities
}

func (ai *AI) convertSliceFromFloat32ToFloat64(list []float32) []float64 {
	result := make([]float64, len(list))

	for i, value := range list {
		result[i] = float64(value)
	}

	return result
}

func (ai *AI) Close() error {
	return ai.model.Session.Close()
}
