package services

import (
	"errors"

	"github.com/AshvinBambhaniya/tic-tac-toe/models"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type GameEngine struct {
	logger *zap.Logger
	// We might need models here to persist state
}

func NewGameEngine(logger *zap.Logger) *GameEngine {
	return &GameEngine{
		logger: logger,
	}
}

// IsWin checks if a 3x3 grid (represented as an array of 9 strings) has a winner.
func (ge *GameEngine) IsWin(board [9]string) string {
	winConditions := [8][3]int{
		{0, 1, 2}, {3, 4, 5}, {6, 7, 8}, // Rows
		{0, 3, 6}, {1, 4, 7}, {2, 5, 8}, // Cols
		{0, 4, 8}, {2, 4, 6},             // Diagonals
	}

	for _, condition := range winConditions {
		if board[condition[0]] != "" &&
			board[condition[0]] == board[condition[1]] &&
			board[condition[0]] == board[condition[2]] {
			return board[condition[0]]
		}
	}
	return ""
}

// IsDraw checks if a 3x3 grid is full without a winner.
func (ge *GameEngine) IsDraw(board [9]string) bool {
	for _, cell := range board {
		if cell == "" {
			return false
		}
	}
	return true
}

// ValidateMove checks if a move is valid based on the current game state.
func (ge *GameEngine) ValidateMove(game models.Game, moves []models.Move, subGridWinners []models.SubGridResult, playerID uuid.UUID, subGridIndex, cellIndex int) error {
	// 1. Check if it's the player's turn
	expectedSymbol := "X"
	if game.CurrentTurn == "O" {
		expectedSymbol = "O"
	}

	// Determine if the player is X or O
	var playerSymbol string
	if game.PlayerXId == playerID {
		playerSymbol = "X"
	} else if game.PlayerOId != nil && *game.PlayerOId == playerID {
		playerSymbol = "O"
	} else {
		return errors.New("player not part of this game")
	}

	if playerSymbol != expectedSymbol {
		return errors.New("not your turn")
	}

	// 2. Check if the game is already finished
	if game.Status != "ongoing" {
		return errors.New("game is already finished")
	}

	// 3. Check if the subGridIndex is valid (must be the active_sub_grid unless it's 9)
	if game.ActiveSubGrid != 9 && int(game.ActiveSubGrid) != subGridIndex {
		return errors.New("invalid sub-grid move")
	}

	// 4. Check if the sub-grid is already won or drawn
	for _, res := range subGridWinners {
		if int(res.GridIndex) == subGridIndex {
			return errors.New("sub-grid already completed")
		}
	}

	// 5. Check if the cell is already occupied
	for _, m := range moves {
		if int(m.SubGridIndex) == subGridIndex && int(m.CellIndex) == cellIndex {
			return errors.New("cell already occupied")
		}
	}

	return nil
}
