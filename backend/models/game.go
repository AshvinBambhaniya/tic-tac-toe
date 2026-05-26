package models

import (
	"time"

	"github.com/doug-martin/goqu/v9"
	"github.com/google/uuid"
)

const (
	GameTable           = "games"
	MoveTable           = "moves"
	SubGridResultsTable = "sub_grid_results"
)

type Game struct {
	ID             uuid.UUID  `json:"id" db:"id"`
	PlayerXId      uuid.UUID  `json:"player_x_id" db:"player_x_id"`
	PlayerOId      *uuid.UUID `json:"player_o_id" db:"player_o_id"`
	CurrentTurn    string     `json:"current_turn" db:"current_turn"`
	ActiveSubGrid  int16      `json:"active_sub_grid" db:"active_sub_grid"`
	WinnerId       *uuid.UUID `json:"winner_id" db:"winner_id"`
	Status         string     `json:"status" db:"status"`
	CreatedAt      time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at" db:"updated_at"`
}

type Move struct {
	ID            uuid.UUID `json:"id" db:"id"`
	GameID        uuid.UUID `json:"game_id" db:"game_id"`
	PlayerID      uuid.UUID `json:"player_id" db:"player_id"`
	SubGridIndex  int16     `json:"sub_grid_index" db:"sub_grid_index"`
	CellIndex     int16     `json:"cell_index" db:"cell_index"`
	Symbol        string    `json:"symbol" db:"symbol"`
	MoveOrder     int       `json:"move_order" db:"move_order"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
}

type SubGridResult struct {
	ID           uuid.UUID  `json:"id" db:"id"`
	GameID       uuid.UUID  `json:"game_id" db:"game_id"`
	GridIndex    int16      `json:"grid_index" db:"grid_index"`
	WinnerSymbol string     `json:"winner_symbol" db:"winner_symbol"`
	WonAtMoveId  *uuid.UUID `json:"won_at_move_id" db:"won_at_move_id"`
	CreatedAt    time.Time  `json:"created_at" db:"created_at"`
}

type GameHistory struct {
	Game
	PlayerXName string  `json:"player_x_name" db:"player_x_name"`
	PlayerOName *string `json:"player_o_name" db:"player_o_name"`
}

type GameModel struct {
	db *goqu.Database
}

func InitGameModel(goqu *goqu.Database) GameModel {
	return GameModel{
		db: goqu,
	}
}

func (m *GameModel) CreateGame(playerXId uuid.UUID) (Game, error) {
	var game Game
	_, err := m.db.Insert(GameTable).Rows(
		goqu.Record{
			"player_x_id": playerXId,
			"status":      "ongoing",
			"current_turn": "X",
			"active_sub_grid": 9,
		},
	).Returning("*").Executor().ScanStruct(&game)
	return game, err
}

func (m *GameModel) GetGame(id uuid.UUID) (Game, error) {
	var game Game
	_, err := m.db.From(GameTable).Where(goqu.Ex{"id": id}).ScanStruct(&game)
	return game, err
}

func (m *GameModel) UpdateGame(game Game) error {
	_, err := m.db.Update(GameTable).Where(goqu.Ex{"id": game.ID}).Set(
		goqu.Record{
			"player_o_id":     game.PlayerOId,
			"current_turn":    game.CurrentTurn,
			"active_sub_grid": game.ActiveSubGrid,
			"winner_id":       game.WinnerId,
			"status":          game.Status,
			"updated_at":      time.Now(),
		},
	).Executor().Exec()
	return err
}

func (m *GameModel) CreateMove(move Move) (Move, error) {
	var createdMove Move
	_, err := m.db.Insert(MoveTable).Rows(
		goqu.Record{
			"game_id":        move.GameID,
			"player_id":      move.PlayerID,
			"sub_grid_index": move.SubGridIndex,
			"cell_index":     move.CellIndex,
			"symbol":         move.Symbol,
			"move_order":     move.MoveOrder,
		},
	).Returning("*").Executor().ScanStruct(&createdMove)
	return createdMove, err
}

func (m *GameModel) GetMoves(gameID uuid.UUID) ([]Move, error) {
	var moves []Move
	err := m.db.From(MoveTable).Where(goqu.Ex{"game_id": gameID}).Order(goqu.I("move_order").Asc()).ScanStructs(&moves)
	return moves, err
}

func (m *GameModel) CreateSubGridResult(res SubGridResult) (SubGridResult, error) {
	var createdRes SubGridResult
	_, err := m.db.Insert(SubGridResultsTable).Rows(
		goqu.Record{
			"game_id":        res.GameID,
			"grid_index":     res.GridIndex,
			"winner_symbol":  res.WinnerSymbol,
			"won_at_move_id": res.WonAtMoveId,
		},
	).Returning("*").Executor().ScanStruct(&createdRes)
	return createdRes, err
}

func (m *GameModel) GetSubGridResults(gameID uuid.UUID) ([]SubGridResult, error) {
	var results []SubGridResult
	err := m.db.From(SubGridResultsTable).Where(goqu.Ex{"game_id": gameID}).ScanStructs(&results)
	return results, err
}

func (m *GameModel) GetActiveGamesByPlayer(playerID uuid.UUID) ([]Game, error) {
	var games []Game
	err := m.db.From(GameTable).Where(
		goqu.Ex{"status": "ongoing"},
		goqu.Or(
			goqu.Ex{"player_x_id": playerID},
			goqu.Ex{"player_o_id": playerID},
		),
	).ScanStructs(&games)
	return games, err
}

func (m *GameModel) MarkStaleGamesAsFinished(threshold time.Duration) (int64, error) {
	cutoff := time.Now().Add(-threshold)
	res, err := m.db.Update(GameTable).
		Set(goqu.Record{"status": "finished", "updated_at": time.Now()}).
		Where(
			goqu.Ex{"status": "ongoing"},
			goqu.I("updated_at").Lt(cutoff),
		).Executor().Exec()
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}
func (m *GameModel) GetGameHistoryByPlayer(playerID uuid.UUID) ([]GameHistory, error) {
	var history []GameHistory

	query := m.db.From(goqu.T(GameTable).As("g")).
		Select(
			"g.id",
			"g.player_x_id",
			"g.player_o_id",
			"g.current_turn",
			"g.active_sub_grid",
			"g.winner_id",
			"g.status",
			"g.created_at",
			"g.updated_at",
			goqu.I("ux.first_name").As("player_x_name"),
			goqu.I("uo.first_name").As("player_o_name"),
		).
		LeftJoin(goqu.T("users").As("ux"), goqu.On(goqu.I("g.player_x_id").Eq(goqu.I("ux.id")))).
		LeftJoin(goqu.T("users").As("uo"), goqu.On(goqu.I("g.player_o_id").Eq(goqu.I("uo.id")))).
		Where(
			goqu.Ex{"g.status": goqu.Op{"neq": "ongoing"}},
			goqu.Or(
				goqu.Ex{"g.player_x_id": playerID},
				goqu.Ex{"g.player_o_id": playerID},
			),
		).
		Order(goqu.I("g.updated_at").Desc())

	err := query.ScanStructs(&history)
	return history, err
}
