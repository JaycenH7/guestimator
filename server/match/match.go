package match

import (
	"log"
	"net/http"

	"github.com/olahol/melody"
)

const MIN_CAPACITY int = 2

type Match struct {
	ID       string
	Capacity int
	Hub      *melody.Melody
	Sessions map[string]*melody.Session

	CurrentPhase     Phase
	playerConnect    chan string
	playerDisconnect chan string
}

func NewMatch(id string, capacity int) *Match {
	if capacity < MIN_CAPACITY {
		capacity = MIN_CAPACITY
	}

	m := &Match{
		ID:               id,
		Capacity:         capacity,
		Hub:              melody.New(),
		Sessions:         make(map[string]*melody.Session),
		playerConnect:    make(chan string),
		playerDisconnect: make(chan string),
	}

	m.Hub.Upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	m.Hub.HandleConnect(m.handlePlayerConnect)

	m.Hub.HandleDisconnect(func(s *melody.Session) {
		playerID := s.Request.URL.Query().Get("player")
		m.playerDisconnect <- playerID
	})

	go m.run()

	return m
}

func (m *Match) run() {
	phases := []Phase{
		NewJoinPhase(),
		new(AnswerPhase),
	}

	for _, m.CurrentPhase = range phases {
		done := m.CurrentPhase.Run(m)

	PhaseLoop:
		for {
			select {
			case playerID := <-m.playerConnect:
				m.CurrentPhase.OnPlayerConnect(playerID)
			case playerID := <-m.playerDisconnect:
				m.CurrentPhase.OnPlayerDisconnect(playerID)
			case <-done:
				break PhaseLoop
			}
		}
	}
}

func (m *Match) handlePlayerConnect(s *melody.Session) {
	playerID := s.Request.URL.Query().Get("player")
	m.Sessions[playerID] = s

	// TODO: cleanup and break out different messages responsiblities to funcs/goroutines
	state := Message{
		Type: MatchStateMsgType,
		MatchState: &MatchState{
			Phase: m.CurrentPhase.Label(),
		},
	}

	msg, err := state.MarshalJSON()
	if err != nil {
		log.Println("Error marshaling match state", err)
	}
	s.Write(msg)

	pl := Message{
		Type:     PlayerJoinMsgType,
		PlayerID: playerID,
	}

	msg, err = pl.MarshalJSON()
	if err != nil {
		log.Println("Error marshaling player_connect message", err)
	}

	m.Hub.BroadcastOthers(msg, s)

	m.playerConnect <- playerID
}
