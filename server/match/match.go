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

	phaseType        PhaseType
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

	go m.offloadEvents()
	go m.run()

	return m
}

func (m Match) offloadEvents() {
	for {
		select {
		case playerID := <-m.playerConnect:
			log.Println("catchAllMessages playerConnect", playerID)
		case playerID := <-m.playerDisconnect:
			log.Println("catchAllMessages playerDisconnect", playerID)
		}
	}
}

func (m Match) PhaseType() PhaseType {
	return m.phaseType
}

func (m *Match) run() {
	for phase := JoinPhase; phase != nil; {
		phase = phase(m)
	}
}

func (m *Match) handlePlayerConnect(s *melody.Session) {
	playerID := s.Request.URL.Query().Get("player")
	m.Sessions[playerID] = s

	// TODO: cleanup and break out different messages responsiblities to funcs/goroutines
	state := Message{
		Type: MatchStateMsgType,
		MatchState: &MatchState{
			Phase: m.PhaseType(),
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
