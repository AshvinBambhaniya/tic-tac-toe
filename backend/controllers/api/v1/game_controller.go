package v1

import (
	"net/http"
	"strings"

	"github.com/AshvinBambhaniya/tic-tac-toe/constants"
	"github.com/AshvinBambhaniya/tic-tac-toe/pkg/websocket"
	"github.com/AshvinBambhaniya/tic-tac-toe/services"
	"github.com/AshvinBambhaniya/tic-tac-toe/utils"
	fiber_ws "github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type GameController struct {
	gameService *services.GameService
	hub         *websocket.Hub
	logger      *zap.Logger
}

func NewGameController(gameService *services.GameService, hub *websocket.Hub, logger *zap.Logger) *GameController {
	return &GameController{
		gameService: gameService,
		hub:         hub,
		logger:      logger,
	}
}

func (ctrl *GameController) CreateGame(c *fiber.Ctx) error {
	userIDStr := c.Locals(constants.ContextUid).(string)
	userID, _ := uuid.Parse(userIDStr)

	var req struct {
		GameMode string `json:"game_mode"`
	}
	if err := c.BodyParser(&req); err != nil {
		req.GameMode = "normal"
	}
	if req.GameMode == "" {
		req.GameMode = "normal"
	}

	game, err := ctrl.gameService.CreateGame(userID, req.GameMode)
	if err != nil {
		ctrl.logger.Error("Failed to create game", zap.Error(err))
		return utils.JSONError(c, http.StatusInternalServerError, "Failed to create game")
	}

	return utils.JSONSuccess(c, http.StatusCreated, game)
}

func (ctrl *GameController) CreateAIGame(c *fiber.Ctx) error {
	userIDStr := c.Locals(constants.ContextUid).(string)
	userID, _ := uuid.Parse(userIDStr)

	var req struct {
		Difficulty int16  `json:"difficulty"`
		GameMode   string `json:"game_mode"`
	}
	if err := c.BodyParser(&req); err != nil {
		return utils.JSONFail(c, http.StatusBadRequest, "Invalid request body")
	}
	if req.GameMode == "" {
		req.GameMode = "normal"
	}

	game, err := ctrl.gameService.CreateAIGame(userID, req.Difficulty, req.GameMode)
	if err != nil {
		ctrl.logger.Error("Failed to create AI game", zap.Error(err))
		return utils.JSONError(c, http.StatusInternalServerError, "Failed to create AI game")
	}

	return utils.JSONSuccess(c, http.StatusCreated, game)
}

func (ctrl *GameController) JoinGame(c *fiber.Ctx) error {
	gameIDStr := c.Params("gameId")
	gameID, err := uuid.Parse(gameIDStr)
	if err != nil {
		return utils.JSONFail(c, http.StatusBadRequest, "Invalid game ID")
	}

	userIDStr := c.Locals(constants.ContextUid).(string)
	userID, _ := uuid.Parse(userIDStr)

	game, err := ctrl.gameService.JoinGame(gameID, userID)
	if err != nil {
		return utils.JSONFail(c, http.StatusBadRequest, err.Error())
	}

	return utils.JSONSuccess(c, http.StatusOK, game)
}

func (ctrl *GameController) Matchmake(c *fiber.Ctx) error {
	userIDStr, ok := c.Locals(constants.ContextUid).(string)
	if !ok || userIDStr == "" {
		ctrl.logger.Error("UserID not found in context for matchmaking")
		return utils.JSONFail(c, http.StatusUnauthorized, "Unauthenticated")
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		ctrl.logger.Error("Failed to parse userID", zap.String("userID", userIDStr), zap.Error(err))
		return utils.JSONFail(c, http.StatusBadRequest, "Invalid UserID")
	}

	ctrl.logger.Info("Received matchmaking request", zap.String("userID", userID.String()))
	ctrl.gameService.AddToMatchmaking(userID)

	return utils.JSONSuccess(c, http.StatusOK, map[string]string{"message": "Added to matchmaking queue"})
}

func (ctrl *GameController) GetActiveGames(c *fiber.Ctx) error {
	userIDStr := c.Locals(constants.ContextUid).(string)
	userID, _ := uuid.Parse(userIDStr)

	games, err := ctrl.gameService.GetActiveGames(userID)
	if err != nil {
		ctrl.logger.Error("Failed to fetch active games", zap.Error(err))
		return utils.JSONError(c, http.StatusInternalServerError, "Failed to fetch active games")
	}

	return utils.JSONSuccess(c, http.StatusOK, games)
}

func (ctrl *GameController) GetGameState(c *fiber.Ctx) error {
	gameIDStr := c.Params("gameId")
	gameID, err := uuid.Parse(gameIDStr)
	if err != nil {
		return utils.JSONFail(c, http.StatusBadRequest, "Invalid game ID")
	}

	state, err := ctrl.gameService.GetFullGameState(gameID)
	if err != nil {
		ctrl.logger.Error("Failed to fetch game state", zap.Error(err))
		return utils.JSONError(c, http.StatusInternalServerError, "Failed to fetch game state")
	}

	return utils.JSONSuccess(c, http.StatusOK, state)
}

func (ctrl *GameController) GetProfile(c *fiber.Ctx) error {
	userIDStr := c.Locals(constants.ContextUid).(string)
	userID, _ := uuid.Parse(userIDStr)

	profile, err := ctrl.gameService.GetPlayerProfile(userID)
	if err != nil {
		ctrl.logger.Error("Failed to fetch player profile", zap.Error(err))
		return utils.JSONError(c, http.StatusInternalServerError, "Failed to fetch player profile")
	}

	return utils.JSONSuccess(c, http.StatusOK, profile)
}

func (ctrl *GameController) HandleWebSocket(c *fiber.Ctx) error {
	userID := c.Locals(constants.ContextUid)
	if userID == nil {
		return fiber.ErrUnauthorized
	}

	if fiber_ws.IsWebSocketUpgrade(c) {
		return c.Next()
	}
	return fiber.ErrUpgradeRequired
}

func (ctrl *GameController) WebSocketHandler(c *fiber_ws.Conn) {
	gameID := c.Params("gameId")
	userID := strings.ToLower(c.Locals(constants.ContextUid).(string))

	client := &websocket.Client{
		Hub:    ctrl.hub,
		Conn:   c,
		Send:   make(chan []byte, 256),
		GameID: gameID,
		ID:     userID,
	}

	ctrl.hub.Register <- client

	// Send initial state to all players in the room (skip if lobby)
	if gameID != "lobby" {
		gameIDuuid, err := uuid.Parse(gameID)
		if err == nil {
			state, err := ctrl.gameService.GetFullGameState(gameIDuuid)
			if err == nil {
				// Broadcast state to EVERYONE in the room.
				// This ensures the other player's UI hides the "Opponent disconnected" banner.
				ctrl.hub.BroadcastToRoom(gameID, struct {
					Type    string      `json:"type"`
					Payload interface{} `json:"payload"`
				}{
					Type:    "STATE_UPDATE",
					Payload: state,
				})
			}
		}
	}

	go client.WritePump()
	client.ReadPump()
}
