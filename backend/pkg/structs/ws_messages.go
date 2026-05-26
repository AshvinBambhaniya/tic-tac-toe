package structs

type WSMessageType string

const (
	WSMessageTypeMove        WSMessageType = "MOVE"
	WSMessageTypeStateUpdate WSMessageType = "STATE_UPDATE"
	WSMessageTypeError       WSMessageType = "ERROR"
	WSMessageTypeMatchFound  WSMessageType = "MATCH_FOUND"
	WSMessageTypeOpponentLeft WSMessageType = "OPPONENT_LEFT"
	WSMessageTypeForfeit     WSMessageType = "FORFEIT"
)

type WSMessage struct {
	Type    WSMessageType `json:"type"`
	Payload interface{}   `json:"payload"`
}

type WSMovePayload struct {
	SubGridIndex int `json:"subGridIndex"`
	CellIndex    int `json:"cellIndex"`
}

type WSStateUpdatePayload struct {
	Game   interface{} `json:"game"`
	Moves  interface{} `json:"moves"`
	Results interface{} `json:"results"`
}
