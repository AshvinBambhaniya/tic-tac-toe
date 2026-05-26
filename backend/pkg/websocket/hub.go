package websocket

import (
	"encoding/json"
	"sync"
	"time"

	"github.com/AshvinBambhaniya/tic-tac-toe/config"
	"github.com/AshvinBambhaniya/tic-tac-toe/pkg/structs"
	"github.com/AshvinBambhaniya/tic-tac-toe/services"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// Hub maintains the set of active clients and broadcasts messages to clients.
type Hub struct {
	// Registered clients.
	// gameID -> map[clientID]*Client
	rooms map[string]map[string]*Client

	// Register requests from the clients.
	Register chan *Client

	// Unregister requests from clients.
	Unregister chan *Client

	// gameID:userID -> timer
	disconnectTimers map[string]*time.Timer

	mu     sync.RWMutex
	logger *zap.Logger

	GameService *services.GameService
	config      config.AppConfig
}

func NewHub(logger *zap.Logger, gameService *services.GameService, config config.AppConfig) *Hub {
	return &Hub{
		rooms:            make(map[string]map[string]*Client),
		Register:         make(chan *Client),
		Unregister:       make(chan *Client),
		disconnectTimers: make(map[string]*time.Timer),
		logger:           logger,
		GameService:      gameService,
		config:           config,
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.mu.Lock()
			// Cancel disconnect timer if it exists
			timerKey := client.GameID + ":" + client.ID
			if timer, ok := h.disconnectTimers[timerKey]; ok {
				timer.Stop()
				delete(h.disconnectTimers, timerKey)
				h.logger.Info("Disconnection timer cancelled", zap.String("gameID", client.GameID), zap.String("clientID", client.ID))
			}

			if h.rooms[client.GameID] == nil {
				h.rooms[client.GameID] = make(map[string]*Client)
			}
			h.rooms[client.GameID][client.ID] = client
			h.mu.Unlock()
			h.logger.Info("Client registered", zap.String("gameID", client.GameID), zap.String("clientID", client.ID))

		case client := <-h.Unregister:
			h.mu.Lock()
			if room, ok := h.rooms[client.GameID]; ok {
				if _, ok := room[client.ID]; ok {
					delete(room, client.ID)
					close(client.Send)

					// If it's a game room (not lobby) and someone is still there, notify them
					if client.GameID != "lobby" && len(room) > 0 {
						// We need to broadcast, but we are holding the lock.
						// BroadcastToRoom also tries to take a lock.
						// To avoid deadlocks, we'll manually broadcast here or unlock first.
						go h.BroadcastToRoom(client.GameID, map[string]interface{}{
							"type":    "OPPONENT_LEFT",
							"payload": "Opponent disconnected",
						})

						// Start timeout timer
						timerKey := client.GameID + ":" + client.ID
						h.disconnectTimers[timerKey] = time.AfterFunc(time.Duration(h.config.GameDisconnectTimeout)*time.Second, func() {
							h.handleDisconnectTimeout(client.GameID, client.ID)
						})
						h.logger.Info("Disconnection timer started", zap.String("gameID", client.GameID), zap.String("clientID", client.ID), zap.Int("timeout", h.config.GameDisconnectTimeout))
					}

					if len(room) == 0 {
						delete(h.rooms, client.GameID)
					}
				}
			}
			h.mu.Unlock()
			h.logger.Info("Client unregistered", zap.String("gameID", client.GameID), zap.String("clientID", client.ID))
		}
	}
}

func (h *Hub) BroadcastToRoom(gameID string, message interface{}) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	room, ok := h.rooms[gameID]
	if !ok {
		return
	}

	payload, err := json.Marshal(message)
	if err != nil {
		h.logger.Error("Failed to marshal broadcast message", zap.Error(err))
		return
	}

	for _, client := range room {
		select {
		case client.Send <- payload:
		default:
			h.logger.Warn("Client send buffer full", zap.String("clientID", client.ID))
		}
	}
}

func (h *Hub) BroadcastToUser(roomID, userID string, message interface{}) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	room, ok := h.rooms[roomID]
	if !ok {
		return
	}

	client, ok := room[userID]
	if !ok {
		return
	}

	payload, err := json.Marshal(message)
	if err != nil {
		h.logger.Error("Failed to marshal user message", zap.Error(err))
		return
	}

	select {
	case client.Send <- payload:
	default:
		h.logger.Warn("Client send buffer full", zap.String("clientID", client.ID))
	}
}

func (h *Hub) handleDisconnectTimeout(gameIDStr, userIDStr string) {
	h.mu.Lock()

	timerKey := gameIDStr + ":" + userIDStr
	delete(h.disconnectTimers, timerKey)

	gameID, err := uuid.Parse(gameIDStr)
	if err != nil {
		h.mu.Unlock()
		return
	}
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		h.mu.Unlock()
		return
	}

	h.logger.Info("Disconnection timeout reached, forfeiting game", zap.String("gameID", gameIDStr), zap.String("userID", userIDStr))

	game, err := h.GameService.ForfeitGame(gameID, userID)
	if err != nil {
		h.logger.Error("Failed to forfeit game on timeout", zap.Error(err))
		h.mu.Unlock()
		return
	}

	// Broadcast the result
	h.mu.Unlock() // Unlock before broadcast to avoid deadlock
	h.BroadcastToRoom(gameIDStr, structs.WSMessage{
		Type:    structs.WSMessageTypeStateUpdate,
		Payload: game,
	})
}
