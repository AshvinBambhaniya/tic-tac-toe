package services

import (
	"errors"
	"time"

	"github.com/AshvinBambhaniya/tic-tac-toe/config"
	"github.com/AshvinBambhaniya/tic-tac-toe/models"
	"github.com/AshvinBambhaniya/tic-tac-toe/pkg/structs"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type GameHub interface {
	BroadcastToUser(roomID, userID string, message interface{})
}

type GameService struct {
	gameModel  *models.GameModel
	gameEngine *GameEngine
	logger     *zap.Logger
	hub        GameHub
	config     config.AppConfig

	matchmakingQueue chan uuid.UUID
}

func NewGameService(gameModel *models.GameModel, gameEngine *GameEngine, logger *zap.Logger, config config.AppConfig) *GameService {
	s := &GameService{
		gameModel:        gameModel,
		gameEngine:       gameEngine,
		logger:           logger,
		config:           config,
		matchmakingQueue: make(chan uuid.UUID, 100),
	}
	go s.runMatchmaking()
	go s.runCleanupWorker()
	return s
}

func (s *GameService) runCleanupWorker() {
	// Run cleanup every hour
	ticker := time.NewTicker(1 * time.Hour)
	for {
		select {
		case <-ticker.C:
			s.logger.Info("Running stale game cleanup...")
			affected, err := s.gameModel.MarkStaleGamesAsFinished(time.Duration(s.config.GameCleanupThreshold) * time.Minute)
			if err != nil {
				s.logger.Error("Failed to cleanup stale games", zap.Error(err))
			} else if affected > 0 {
				s.logger.Info("Stale game cleanup complete", zap.Int64("games_cleaned", affected))
			}
		}
	}
}

func (s *GameService) SetHub(hub GameHub) {
	s.hub = hub
}

func (s *GameService) AddToMatchmaking(playerID uuid.UUID) {
	s.matchmakingQueue <- playerID
}

func (s *GameService) runMatchmaking() {
	var waitingPlayer *uuid.UUID

	for playerID := range s.matchmakingQueue {
		s.logger.Info("Received player in matchmaking queue", zap.String("playerID", playerID.String()))

		if waitingPlayer == nil {
			p := playerID
			waitingPlayer = &p
			s.logger.Info("Player added to waiting slot", zap.String("playerID", playerID.String()))
			continue
		}

		if *waitingPlayer == playerID {
			s.logger.Info("Same player tried to match again, skipping", zap.String("playerID", playerID.String()))
			continue
		}

		// Found two players!
		p1 := *waitingPlayer
		p2 := playerID
		s.logger.Info("Found a match!", zap.String("player1", p1.String()), zap.String("player2", p2.String()))

		game, err := s.CreateGame(p1)
		if err != nil {
			s.logger.Error("Failed to create match game", zap.Error(err))
			waitingPlayer = nil
			continue
		}

		_, err = s.JoinGame(game.ID, p2)
		if err != nil {
			s.logger.Error("Failed to join matched player", zap.Error(err))
		}
		// Notify both players if hub is set
		if s.hub != nil {
			msg := structs.WSMessage{
				Type: structs.WSMessageTypeMatchFound,
				Payload: map[string]string{
					"game_id": game.ID.String(),
				},
			}
			s.logger.Info("Notifying players of match", zap.String("gameID", game.ID.String()))

			// Try notifying immediately
			s.hub.BroadcastToUser("lobby", p1.String(), msg)
			s.hub.BroadcastToUser("lobby", p2.String(), msg)

			// Also notify after a short delay in case of slow websocket registration
			go func(id1, id2 uuid.UUID) {
				time.Sleep(1 * time.Second)
				s.hub.BroadcastToUser("lobby", id1.String(), msg)
				s.hub.BroadcastToUser("lobby", id2.String(), msg)
			}(p1, p2)
		} else {

			s.logger.Warn("Match found but Hub is not initialized in GameService")
		}

		waitingPlayer = nil
	}
}

func (s *GameService) CreateGame(playerXId uuid.UUID) (models.Game, error) {
	return s.gameModel.CreateGame(playerXId)
}

func (s *GameService) JoinGame(gameID uuid.UUID, playerOId uuid.UUID) (models.Game, error) {
	game, err := s.gameModel.GetGame(gameID)
	if err != nil {
		return game, err
	}

	if game.PlayerOId != nil && *game.PlayerOId != playerOId {
		return game, errors.New("game already has two players")
	}

	if game.PlayerXId == playerOId {
		return game, errors.New("cannot join your own game as opponent")
	}

	game.PlayerOId = &playerOId
	err = s.gameModel.UpdateGame(game)
	return game, err
}

func (s *GameService) GetFullGameState(gameID uuid.UUID) (interface{}, error) {
	game, err := s.gameModel.GetGame(gameID)
	if err != nil {
		return nil, err
	}

	moves, err := s.gameModel.GetMoves(gameID)
	if err != nil {
		return nil, err
	}

	results, err := s.gameModel.GetSubGridResults(gameID)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"game":    game,
		"moves":   moves,
		"results": results,
	}, nil
}

func (s *GameService) GetActiveGames(playerID uuid.UUID) ([]models.Game, error) {
	return s.gameModel.GetActiveGamesByPlayer(playerID)
}

func (s *GameService) GetPlayerProfile(playerID uuid.UUID) (structs.ResProfile, error) {
	history, err := s.gameModel.GetGameHistoryByPlayer(playerID)
	if err != nil {
		return structs.ResProfile{}, err
	}

	res := structs.ResProfile{
		History: history,
	}

	for _, g := range history {
		res.TotalGames++
		if g.Status == "draw" {
			res.Draws++
		} else if g.WinnerId != nil && *g.WinnerId == playerID {
			res.Wins++
		} else {
			res.Losses++
		}
	}

	if res.TotalGames > 0 {
		res.WinRate = (float64(res.Wins) / float64(res.TotalGames)) * 100
	}

	return res, nil
}

func (s *GameService) ForfeitGame(gameID uuid.UUID, forfeitingPlayerID uuid.UUID) (interface{}, error) {
	game, err := s.gameModel.GetGame(gameID)
	if err != nil {
		return nil, err
	}

	if game.Status != "ongoing" {
		return nil, errors.New("game is already finished")
	}

	// Determine the winner (the other player)
	var winnerID uuid.UUID
	if game.PlayerXId == forfeitingPlayerID {
		if game.PlayerOId == nil {
			// If there's no opponent yet, just finish the game without a winner
			game.Status = "finished"
		} else {
			winnerID = *game.PlayerOId
			game.WinnerId = &winnerID
			game.Status = "finished"
		}
	} else if game.PlayerOId != nil && *game.PlayerOId == forfeitingPlayerID {
		winnerID = game.PlayerXId
		game.WinnerId = &winnerID
		game.Status = "finished"
	} else {
		return nil, errors.New("player not part of this game")
	}

	err = s.gameModel.UpdateGame(game)
	if err != nil {
		return nil, err
	}

	return s.GetFullGameState(gameID)
}

func (s *GameService) ProcessMove(gameID uuid.UUID, playerID uuid.UUID, subGridIndex, cellIndex int) (interface{}, error) {
	game, err := s.gameModel.GetGame(gameID)
	if err != nil {
		return nil, err
	}

	moves, err := s.gameModel.GetMoves(gameID)
	if err != nil {
		return nil, err
	}

	results, err := s.gameModel.GetSubGridResults(gameID)
	if err != nil {
		return nil, err
	}

	err = s.gameEngine.ValidateMove(game, moves, results, playerID, subGridIndex, cellIndex)
	if err != nil {
		return nil, err
	}

	// 1. Create the move
	symbol := "X"
	if game.CurrentTurn == "O" {
		symbol = "O"
	}

	move := models.Move{
		GameID:       gameID,
		PlayerID:     playerID,
		SubGridIndex: int16(subGridIndex),
		CellIndex:    int16(cellIndex),
		Symbol:       symbol,
		MoveOrder:    len(moves) + 1,
	}

	move, err = s.gameModel.CreateMove(move)
	if err != nil {
		return nil, err
	}

	// 2. Check if sub-grid is won
	subGridMoves := make([]models.Move, 0)
	for _, m := range append(moves, move) {
		if int(m.SubGridIndex) == subGridIndex {
			subGridMoves = append(subGridMoves, m)
		}
	}

	var board [9]string
	for _, m := range subGridMoves {
		board[m.CellIndex] = m.Symbol
	}

	winner := s.gameEngine.IsWin(board)
	if winner != "" {
		res := models.SubGridResult{
			GameID:       gameID,
			GridIndex:    int16(subGridIndex),
			WinnerSymbol: winner,
			WonAtMoveId:  &move.ID,
		}
		_, err = s.gameModel.CreateSubGridResult(res)
		if err != nil {
			return nil, err
		}
		// Refresh results
		results, _ = s.gameModel.GetSubGridResults(gameID)
	} else if s.gameEngine.IsDraw(board) {
		res := models.SubGridResult{
			GameID:       gameID,
			GridIndex:    int16(subGridIndex),
			WinnerSymbol: "D",
			WonAtMoveId:  &move.ID,
		}
		_, err = s.gameModel.CreateSubGridResult(res)
		if err != nil {
			return nil, err
		}
		// Refresh results
		results, _ = s.gameModel.GetSubGridResults(gameID)
	}

	// 3. Update game state (turn and active sub-grid)
	game.CurrentTurn = "X"
	if symbol == "X" {
		game.CurrentTurn = "O"
	}

	// Determine next active sub-grid
	nextActiveSubGrid := int16(cellIndex)
	isNextSubGridWon := false
	for _, res := range results {
		if res.GridIndex == nextActiveSubGrid {
			isNextSubGridWon = true
			break
		}
	}

	if isNextSubGridWon {
		game.ActiveSubGrid = 9 // Play anywhere
	} else {
		game.ActiveSubGrid = nextActiveSubGrid
	}

	// 4. Check if entire game is won
	var metaBoard [9]string
	for _, res := range results {
		metaBoard[res.GridIndex] = res.WinnerSymbol
	}

	gameWinner := s.gameEngine.IsWin(metaBoard)
	if gameWinner != "" {
		game.Status = "finished"
		if gameWinner == "X" {
			game.WinnerId = &game.PlayerXId
		} else {
			game.WinnerId = game.PlayerOId
		}
	} else if s.gameEngine.IsDraw(metaBoard) {
		game.Status = "draw"
	}

	err = s.gameModel.UpdateGame(game)
	if err != nil {
		return nil, err
	}

	return s.GetFullGameState(gameID)
}
