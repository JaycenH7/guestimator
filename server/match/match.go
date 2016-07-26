package match

import (
	"log"
	"math"
	"net/http"
	"sort"
	"time"

	"github.com/mrap/guestimator/models"
	"github.com/olahol/melody"
)

const MIN_CAPACITY int = 2

type playerSession struct {
	playerID string
	session  *melody.Session
}

type Match struct {
	ID        string
	Capacity  int
	Hub       *melody.Melody
	Questions []models.Question
	Sessions  map[string]*melody.Session

	CurrentPhase     Phase
	Scores           Scores
	CurrentQuestion  models.Question
	playerConnect    chan playerSession
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
		playerConnect:    make(chan playerSession),
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
		done := m.CurrentPhase.Run(m)

		m.broadcastMatchState()
		phaseTicker := time.NewTicker(300 * time.Millisecond)

	PhaseLoop:
		for {
			select {
			case <-phaseTicker.C:
				m.broadcastMatchState()
			case ps := <-m.playerConnect:
				m.Sessions[ps.playerID] = ps.session
				m.CurrentPhase.OnPlayerConnect(ps.playerID)
			case playerID := <-m.playerDisconnect:
				m.CurrentPhase.OnPlayerDisconnect(playerID)
			case guess := <-m.playerGuess:
				m.CurrentPhase.OnPlayerGuess(guess)
			case round := <-m.onRoundComplete:
				phases.Prepend(NewGuessResultPhase())

				distPlayerMap := make(map[float64][]string)
				dists := make([]float64, 0)

				answer := m.CurrentQuestion.Answer
				var dist float64
				for playerID, playerGuess := range round.Guesses {
					if playerGuess.Guess.Min > answer.Exact || playerGuess.Guess.Max < answer.Exact {
						dist = math.MaxFloat64
					} else {
						dist = math.Abs(answer.Exact-playerGuess.Guess.Min) + math.Abs(answer.Exact-playerGuess.Guess.Max)
					}
					distPlayerMap[dist] = append(distPlayerMap[dist], playerID)
					dists = append(dists, dist)
				}

				sort.Float64s(dists)

				for i, score := range dists {
					for _, playerID := range distPlayerMap[score] {
						if score < math.MaxFloat64 {
							m.Scores[playerID] += len(round.Guesses) - i
						} else {
							m.Scores[playerID] += 0
						}
					}
				}

			case <-done:
				break PhaseLoop
			}
		}

		phaseTicker.Stop()

		if phases.Size() == 0 {
			if len(m.Questions) == 0 {
				phases.Append(NewMatchResultPhase())
			} else {
				m.CurrentQuestion = m.Questions[0]
				if err := m.CurrentQuestion.PopulateAnswer(); err != nil {
					// TODO: we probably don't want to use questions unless we have an exact answer
					log.Println("Error getting answer from question.", err)
				}
				phases.Append(NewGuessPhase(m.CurrentQuestion))
				m.Questions = m.Questions[1:]
			}
		}
	}
}

func (m *Match) handlePlayerConnect(s *melody.Session) {
	playerID := s.Request.URL.Query().Get("player")

	pl := Message{
		Type:     PlayerJoinMsgType,
		PlayerID: playerID,
	}
	msg, err := pl.MarshalJSON()
	if err != nil {
		log.Println("Error marshaling player_connect message", err)
	}
	m.Hub.BroadcastOthers(msg, s)

	m.broadcastMatchState(s)

	m.playerConnect <- playerSession{
		playerID: playerID,
		session:  s,
	}
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

func (m *Match) broadcastMatchState(sessions ...*melody.Session) {
	state := MatchState{
		Phase: m.CurrentPhase.Label(),
	}

	if m.CurrentPhase.TimeRemaining() > 0 {
		state.PhaseMsLeft = int(m.CurrentPhase.TimeRemaining() / time.Millisecond)
	}

	switch m.CurrentPhase.(type) {
	case *GuessPhase:
		question := m.CurrentQuestion.SansAnswersAt(m.CurrentQuestion.Positions[0])
		state.Question = &question
	case *GuessResultPhase:
		state.Question = &m.CurrentQuestion
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

	if len(sessions) == 0 {
		m.Hub.Broadcast(msgJson)
	} else {
		for _, s := range sessions {
			s.Write(msgJson)
		}
	}
}
