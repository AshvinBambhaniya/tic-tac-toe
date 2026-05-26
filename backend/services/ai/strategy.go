package ai

import (
	"github.com/AshvinBambhaniya/tic-tac-toe/models"
)

// Move represents a coordinate on the board
type Move struct {
	SubGridIndex int
	CellIndex    int
}

// Strategy defines the interface for different AI implementations
type Strategy interface {
	GetBestMove(game models.Game, moves []models.Move, results []models.SubGridResult) (Move, error)
}
