package match

import (
	"log"
	"math"
	"net/http"
	"sort"

	"github.com/mrap/guestimator/models"
	"github.com/olahol/melody"
)

const MIN_CAPACITY int = 2

type Match struct {
	ID        string
	Capacity  int
	Hub       *melody.Melody
	Questions []models.Question
	Sessions  map[string]*melody.Session

	CurrentPhase     Phase
	Scores           Scores
	CurrentQuestion  models.Question
	playerConnect    chan string
	playerDisconnect chan string
	playerGuess      chan PlayerGuess
	onRoundComplete  chan Round
}

func NewMatch(id string, capacity int, questions []models.Question) *Match {
	if capacity < MIN_CAPACITY {
		capacity = MIN_CAPACITY
	}

	m := &Match{
		ID:               id,
		Capacity:         capacity,
		Hub:              melody.New(),
		Questions:        questions,
		Sessions:         make(map[string]*melody.Session),
		Scores:           make(Scores),
		playerConnect:    make(chan string),
		playerGuess:      make(chan PlayerGuess),
		playerDisconnect: make(chan string),
		onRoundComplete:  make(chan Round),
	}

	m.Hub.Upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	m.Hub.HandleMessage(m.handlePlayerMessage)

	m.Hub.HandleConnect(m.handlePlayerConnect)

	m.Hub.HandleDisconnect(func(s *melody.Session) {
		playerID := s.Request.URL.Query().Get("player")
		m.playerDisconnect <- playerID
	})

	go m.run()

	return m
}

func (m *Match) run() {
	phases := NewPhaseQueue(
		NewJoinPhase(),
	)

	for phases.Size() > 0 {
		m.CurrentPhase = phases.Next()
		m.broadcastMatchState()
		done := m.CurrentPhase.Run(m)

	PhaseLoop:
		for {
			select {
			case playerID := <-m.playerConnect:
				m.CurrentPhase.OnPlayerConnect(playerID)
			case playerID := <-m.playerDisconnect:
				m.CurrentPhase.OnPlayerDisconnect(playerID)
			case guess := <-m.playerGuess:
				m.CurrentPhase.OnPlayerGuess(guess)
			case round := <-m.onRoundComplete:
				phases.Prepend(NewGuessResultPhase())

				exact, err := m.CurrentQuestion.FirstAnswer()
				if err != nil {
					log.Println("Error getting answer from question.", err)
				}

				distPlayerMap := make(map[float64][]string)
				dists := make([]float64, 0)

				var dist float64
				for playerID, playerGuess := range round.Guesses {
					if playerGuess.Guess.Min > exact || playerGuess.Guess.Max < exact {
						dist = math.MaxFloat64
					} else {
						dist = math.Abs(exact-playerGuess.Guess.Min) + math.Abs(exact-playerGuess.Guess.Max)
					}
					distPlayerMap[dist] = append(distPlayerMap[dist], playerID)
					dists = append(dists, dist)
				}

				sort.Float64s(dists)

				for i, score := range dists {
					for _, playerID := range distPlayerMap[score] {
						if score < math.MaxFloat64 {
							m.Scores[playerID] += m.Capacity - i
						} else {
							m.Scores[playerID] += 0
						}
					}
				}

			case <-done:
				break PhaseLoop
			}
		}

		if phases.Size() == 0 {
			if len(m.Questions) == 0 {
				phases.Append(NewMatchResultPhase())
			} else {
				m.CurrentQuestion = m.Questions[0]
				phases.Append(NewGuessPhase(m.CurrentQuestion))
				m.Questions = m.Questions[1:]
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

func (m *Match) handlePlayerMessage(s *melody.Session, inMsg []byte) {
	msg := &Message{}

	err := msg.UnmarshalJSON(inMsg)
	if err != nil {
		log.Println("Error unmarshaling player message", err)
		return
	}

	if msg.Type == GuessMsgType {
		guess := PlayerGuess{
			PlayerID: s.Request.URL.Query().Get("player"),
			Guess:    *msg.Guess,
		}
		m.playerGuess <- guess
	}
}

func (m *Match) broadcastMatchState() {
	state := MatchState{
		Phase: m.CurrentPhase.Label(),
	}

	switch m.CurrentPhase.(type) {
	case *GuessPhase:
		question := m.CurrentQuestion.SansAnswers()
		state.Question = &question
	case *GuessResultPhase:
		state.Scores = m.Scores
	}

	msg := Message{
		Type:       MatchStateMsgType,
		MatchState: &state,
	}

	msgJson, err := msg.MarshalJSON()
	if err != nil {
		log.Println("Error marshaling match state", err)
	}
	m.Hub.Broadcast(msgJson)
}
