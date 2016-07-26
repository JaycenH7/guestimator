//go:generate easyjson $GOFILE
package match

import "github.com/mrap/guestimator/models"

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
	Phase       string           `json:"phase"`
	PhaseMsLeft int              `json:"ms_left,omitempty"`
	Question    *models.Question `json:"question,omitempty"`
	Scores      Scores           `json:"scores,omitempty"`
}

//easyjson:json
type Guess struct {
	Min float64 `json:"min"`
	Max float64 `json:"max"`
}

type PlayerGuess struct {
	PlayerID string
	Guess    Guess
}

type Scores map[string]int
