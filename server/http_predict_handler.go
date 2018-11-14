package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/cheshir/tttnn/server/ai"
)

// TODO add sessions.
// TODO add logger.
type predictHandler struct {
	*baseHandler
	ai *ai.AI
}

type predictRequest []int

type predictResponse struct {
	Move        int     `json:"move"`
	Probability float64 `json:"probability"`
}

func newPredictHandler(ai *ai.AI) *predictHandler {
	return &predictHandler{
		ai: ai,
	}
}

func (handler *predictHandler) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		log.Printf("[ERROR] Method '%v' is not allowed\n", request.Method)
		handler.sendError(response, http.StatusBadRequest)
		return
	}

	moves := predictRequest{}
	err := json.NewDecoder(request.Body).Decode(&moves)
	if err != nil {
		log.Printf("[ERROR] Failed to read and decode request data: %v\n", err)
		handler.sendError(response, http.StatusBadRequest)
		return
	}
	request.Body.Close()

	table, err := handler.createTable(moves)
	if err != nil {
		log.Printf("[ERROR] Failed to create table from moves: %v\n", err)
		handler.sendError(response, http.StatusBadRequest)
		return
	}

	nextMove, probability, err := handler.ai.Predict(table)
	if err != nil {
		log.Printf("[ERROR] Failed to predict next move: %v\n", err)
		handler.sendError(response, http.StatusInternalServerError)
		return
	}

	responseData := predictResponse{
		Move:        nextMove,
		Probability: probability,
	}

	if err := json.NewEncoder(response).Encode(responseData); err != nil {
		log.Printf("[ERROR] Failed to encode response: %v. Data: %#v\n", err, responseData)
		handler.sendError(response, http.StatusInternalServerError)
		return
	}
}

func (handler *predictHandler) createTable(moves []int) (*ai.Table, error) {
	var table ai.Table
	var blackMove, whiteMove float32

	if len(moves)%2 == 0 {
		blackMove = ai.PlayerMove
		whiteMove = ai.AIMove
	} else {
		blackMove = ai.AIMove
		whiteMove = ai.PlayerMove
	}

	nextMove := blackMove

	for _, move := range moves {
		row, column, err := ai.TableCellIndex(move).ToMatrixCoords()
		if err != nil {
			return nil, err
		}

		table[row][column] = nextMove

		if nextMove == blackMove {
			nextMove = whiteMove
		} else {
			nextMove = blackMove
		}
	}

	return &table, nil
}
