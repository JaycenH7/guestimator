//go:generate easyjson $GOFILE
package match

type MsgType int

const (
	MatchStateMsgType MsgType = iota
	PlayerJoinMsgType
)

//easyjson:json
type Message struct {
	Type       MsgType     `json:"type"`
	PlayerID   string      `json:"player_id,omitempty"`
	MatchState *MatchState `json:"match_state,omitempty"`
}

//easyjson:json
type MatchState struct {
	Phase PhaseType `json:"phase"`
}
