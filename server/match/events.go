//go:generate ffjson $GOFILE
package match

type EventPayload struct {
	Type string      `json:"type"`
	Body interface{} `json:"body"`
}

type PlayerJoinEvent struct {
	PlayerID string `json:"player_id"`
}
