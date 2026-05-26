package structs

import "github.com/AshvinBambhaniya/tic-tac-toe/models"

// All response sturcts
// Response struct have Res prefix

type ResProfile struct {
	TotalGames int                  `json:"total_games"`
	Wins       int                  `json:"wins"`
	Losses     int                  `json:"losses"`
	Draws      int                  `json:"draws"`
	WinRate    float64              `json:"win_rate"`
	History    []models.GameHistory `json:"history"`
}
