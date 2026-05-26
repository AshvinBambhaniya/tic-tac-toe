package websocket

import (
	"encoding/json"
	"log"
	"time"

	"github.com/AshvinBambhaniya/tic-tac-toe/pkg/structs"
	fiber_ws "github.com/gofiber/contrib/websocket"
	"github.com/google/uuid"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

// Client is a middleman between the websocket connection and the hub.
type Client struct {
	Hub *Hub

	// The websocket connection.
	Conn *fiber_ws.Conn

	// Buffered channel of outbound messages.
	Send chan []byte

	// Game and User Info
	GameID string
	ID     string // UserID
}

// ReadPump pumps messages from the websocket connection to the hub.
func (c *Client) ReadPump() {
	defer func() {
		c.Hub.Unregister <- c
		c.Conn.Close()
	}()
	c.Conn.SetReadLimit(maxMessageSize)
	c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	c.Conn.SetPongHandler(func(string) error { c.Conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			if fiber_ws.IsUnexpectedCloseError(err, fiber_ws.CloseGoingAway, fiber_ws.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		c.handleMessage(message)
	}
}

// WritePump pumps messages from the hub to the websocket connection.
func (c *Client) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.Send:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				c.Conn.WriteMessage(fiber_ws.CloseMessage, []byte{})
				return
			}

			w, err := c.Conn.NextWriter(fiber_ws.BinaryMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// Add queued chat messages to the current websocket message.
			n := len(c.Send)
			for i := 0; i < n; i++ {
				w.Write(<-c.Send)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.Conn.WriteMessage(fiber_ws.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func (c *Client) handleMessage(message []byte) {
	var wsMsg structs.WSMessage
	err := json.Unmarshal(message, &wsMsg)
	if err != nil {
		log.Printf("Error unmarshaling message: %v", err)
		return
	}

	switch wsMsg.Type {
	case structs.WSMessageTypeMove:
		payloadBytes, _ := json.Marshal(wsMsg.Payload)
		var movePayload structs.WSMovePayload
		json.Unmarshal(payloadBytes, &movePayload)
		log.Printf("Received move: SubGridIndex=%d, CellIndex=%d", movePayload.SubGridIndex, movePayload.CellIndex)

		gameID, err := uuid.Parse(c.GameID)
		if err != nil {
			return
		}
		playerID, err := uuid.Parse(c.ID)
		if err != nil {
			return
		}

		log.Printf("Processing move for game %s, player %s", c.GameID, c.ID)
		game, err := c.Hub.GameService.ProcessMove(gameID, playerID, movePayload.SubGridIndex, movePayload.CellIndex)
		if err != nil {
			c.Send <- c.errorPayload(err.Error())
			return
		}

		// Broadcast new state to all players in the room
		c.Hub.BroadcastToRoom(c.GameID, structs.WSMessage{
			Type:    structs.WSMessageTypeStateUpdate,
			Payload: game,
		})
	case structs.WSMessageTypeForfeit:
		gameID, err := uuid.Parse(c.GameID)
		if err != nil {
			return
		}
		playerID, err := uuid.Parse(c.ID)
		if err != nil {
			return
		}

		game, err := c.Hub.GameService.ForfeitGame(gameID, playerID)
		if err != nil {
			c.Send <- c.errorPayload(err.Error())
			return
		}

		// Broadcast new state to all players in the room
		c.Hub.BroadcastToRoom(c.GameID, structs.WSMessage{
			Type:    structs.WSMessageTypeStateUpdate,
			Payload: game,
		})
	}
}

func (c *Client) errorPayload(msg string) []byte {
	errPayload, _ := json.Marshal(structs.WSMessage{
		Type:    structs.WSMessageTypeError,
		Payload: msg,
	})
	return errPayload
}
