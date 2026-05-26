package ai

import (
	"math"
	"math/rand"
	"time"

	"github.com/AshvinBambhaniya/tic-tac-toe/models"
)

type MinimaxStrategy struct {
	maxDepth int
}

func NewMinimaxStrategy(depth int) *MinimaxStrategy {
	rand.Seed(time.Now().UnixNano())
	return &MinimaxStrategy{maxDepth: depth}
}

func (s *MinimaxStrategy) GetBestMove(game models.Game, moves []models.Move, results []models.SubGridResult) (Move, error) {
	state := s.toBoardState(game, moves, results)
	
	// If it's easy mode (depth 0), pick a random move
	if s.maxDepth == 0 {
		validMoves := s.getValidMoves(state)
		if len(validMoves) == 0 {
			return Move{}, nil
		}
		return validMoves[rand.Intn(len(validMoves))], nil
	}

	bestMove, _ := s.minimax(state, s.maxDepth, math.Inf(-1), math.Inf(1), true)
	return bestMove, nil
}

type boardState struct {
	grids          [9][9]string
	subGridWinners [9]string
	activeSubGrid  int // 0-8, or 9 for any
	currentTurn    string
	botSymbol      string
}

func (s *MinimaxStrategy) toBoardState(game models.Game, moves []models.Move, results []models.SubGridResult) boardState {
	state := boardState{
		activeSubGrid: int(game.ActiveSubGrid),
		currentTurn:   game.CurrentTurn,
		botSymbol:     "O", // By convention in our matchmaking
	}

	for _, m := range moves {
		state.grids[m.SubGridIndex][m.CellIndex] = m.Symbol
	}

	for _, r := range results {
		state.subGridWinners[r.GridIndex] = r.WinnerSymbol
	}

	return state
}

func (s *MinimaxStrategy) getValidMoves(state boardState) []Move {
	var moves []Move

	if state.activeSubGrid == 9 {
		for g := 0; g < 9; g++ {
			if state.subGridWinners[g] == "" {
				for c := 0; c < 9; c++ {
					if state.grids[g][c] == "" {
						moves = append(moves, Move{SubGridIndex: g, CellIndex: c})
					}
				}
			}
		}
	} else {
		g := state.activeSubGrid
		for c := 0; c < 9; c++ {
			if state.grids[g][c] == "" {
				moves = append(moves, Move{SubGridIndex: g, CellIndex: c})
			}
		}
	}

	return moves
}

func (s *MinimaxStrategy) minimax(state boardState, depth int, alpha, beta float64, isMaximizing bool) (Move, float64) {
	// Check terminal states
	winner := s.checkWinner(state.subGridWinners)
	if winner == state.botSymbol {
		return Move{}, 1000 + float64(depth)
	}
	if winner != "" && winner != "D" {
		return Move{}, -1000 - float64(depth)
	}
	if depth == 0 {
		return Move{}, s.evaluate(state)
	}

	validMoves := s.getValidMoves(state)
	if len(validMoves) == 0 {
		return Move{}, 0 // Draw
	}

	var bestMove Move
	if isMaximizing {
		maxEval := math.Inf(-1)
		for _, move := range validMoves {
			newState := s.applyMove(state, move)
			_, eval := s.minimax(newState, depth-1, alpha, beta, false)
			if eval > maxEval {
				maxEval = eval
				bestMove = move
			}
			alpha = math.Max(alpha, eval)
			if beta <= alpha {
				break
			}
		}
		return bestMove, maxEval
	} else {
		minEval := math.Inf(1)
		for _, move := range validMoves {
			newState := s.applyMove(state, move)
			_, eval := s.minimax(newState, depth-1, alpha, beta, true)
			if eval < minEval {
				minEval = eval
				bestMove = move
			}
			beta = math.Min(beta, eval)
			if beta <= alpha {
				break
			}
		}
		return bestMove, minEval
	}
}

func (s *MinimaxStrategy) applyMove(state boardState, move Move) boardState {
	newState := state
	newState.grids[move.SubGridIndex][move.CellIndex] = state.currentTurn
	
	// Update sub-grid winner
	if newState.subGridWinners[move.SubGridIndex] == "" {
		winner := s.checkWinner(newState.grids[move.SubGridIndex])
		if winner != "" {
			newState.subGridWinners[move.SubGridIndex] = winner
		} else if s.isGridFull(newState.grids[move.SubGridIndex]) {
			newState.subGridWinners[move.SubGridIndex] = "D"
		}
	}

	// Update turn
	if state.currentTurn == "X" {
		newState.currentTurn = "O"
	} else {
		newState.currentTurn = "X"
	}

	// Update active sub-grid
	nextGrid := move.CellIndex
	if newState.subGridWinners[nextGrid] != "" {
		newState.activeSubGrid = 9
	} else {
		newState.activeSubGrid = int(nextGrid)
	}

	return newState
}

func (s *MinimaxStrategy) checkWinner(board [9]string) string {
	winConditions := [8][3]int{
		{0, 1, 2}, {3, 4, 5}, {6, 7, 8},
		{0, 3, 6}, {1, 4, 7}, {2, 5, 8},
		{0, 4, 8}, {2, 4, 6},
	}
	for _, c := range winConditions {
		if board[c[0]] != "" && board[c[0]] != "D" && board[c[0]] == board[c[1]] && board[c[0]] == board[c[2]] {
			return board[c[0]]
		}
	}
	return ""
}

func (s *MinimaxStrategy) isGridFull(grid [9]string) bool {
	for _, cell := range grid {
		if cell == "" {
			return false
		}
	}
	return true
}

func (s *MinimaxStrategy) evaluate(state boardState) float64 {
	score := 0.0
	
	// Meta-board score
	score += s.evaluateGrid(state.subGridWinners, state.botSymbol) * 50

	// Sub-grid scores
	for i := 0; i < 9; i++ {
		if state.subGridWinners[i] == "" {
			score += s.evaluateGrid(state.grids[i], state.botSymbol)
		}
	}

	return score
}

func (s *MinimaxStrategy) evaluateGrid(grid [9]string, symbol string) float64 {
	opponent := "X"
	if symbol == "X" {
		opponent = "O"
	}

	score := 0.0
	winConditions := [8][3]int{
		{0, 1, 2}, {3, 4, 5}, {6, 7, 8},
		{0, 3, 6}, {1, 4, 7}, {2, 5, 8},
		{0, 4, 8}, {2, 4, 6},
	}

	for _, c := range winConditions {
		myCount := 0
		oppCount := 0
		for i := 0; i < 3; i++ {
			if grid[c[i]] == symbol {
				myCount++
			} else if grid[c[i]] == opponent {
				oppCount++
			}
		}

		if myCount == 2 && oppCount == 0 {
			score += 10
		} else if myCount == 1 && oppCount == 0 {
			score += 1
		} else if oppCount == 2 && myCount == 0 {
			score -= 10
		} else if oppCount == 1 && myCount == 0 {
			score -= 1
		}
	}

	// Center advantage
	if grid[4] == symbol {
		score += 3
	} else if grid[4] == opponent {
		score -= 3
	}

	return score
}
