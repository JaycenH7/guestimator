//go:generate easyjson $GOFILE
package event

type Type int

const (
	TypePlayerJoin Type = iota
)

//easyjson:json
type Event struct {
	Type     Type   `json:"type"`
	PlayerID string `json:"player_id"`
}
