//go:generate easyjson $GOFILE
package match

type MsgType int

const (
	EmptyMsgType MsgType = iota
	MatchStateMsgType
	PlayerJoinMsgType
	GuessMsgType
)

//easyjson:json
type Message struct {
	Type       MsgType     `json:"type"`
	PlayerID   string      `json:"player_id,omitempty"`
	MatchState *MatchState `json:"match_state,omitempty"`
	Guess      *Guess      `json:"guess,omitempty"`
}

//easyjson:json
type MatchState struct {
	Phase string `json:"phase"`
}

//easyjson:json
type Guess struct {
	Min int `json:"min"`
	Max int `json:"max"`
}

type PlayerGuess struct {
	PlayerID string
	Guess    Guess
}
