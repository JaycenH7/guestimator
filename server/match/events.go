//go:generate easyjson $GOFILE
package match

//easyjson:json
type EventPayload struct {
	Type string      `json:"type"`
	Body interface{} `json:"body"`
}

//easyjson:json
type PlayerJoinEvent struct {
	PlayerID string `json:"player_id"`
}
