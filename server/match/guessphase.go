package match

import (
	"time"

	"github.com/mrap/guestimator/models"
)

type GuessPhase struct {
	onPlayerGuess chan PlayerGuess
	guesses       Guesses
	question      models.Question
	endTime       time.Time
}

func NewGuessPhase(question models.Question) *GuessPhase {
	p := &GuessPhase{
		guesses:       make(Guesses),
		onPlayerGuess: make(chan PlayerGuess),
		question:      question,
	}
	return p
}

func (p GuessPhase) Label() string {
	return "Guess"
}

func (p *GuessPhase) OnPlayerGuess(guess PlayerGuess) {
	p.onPlayerGuess <- guess
}

func (p *GuessPhase) OnPlayerConnect(id string) {
}

func (p *GuessPhase) OnPlayerDisconnect(id string) {
}

func (p GuessPhase) TimeRemaining() time.Duration {
	return p.endTime.Sub(time.Now())
}

func (p *GuessPhase) Run(m *Match) <-chan struct{} {
	done := make(chan struct{})
	p.endTime = time.Now().Add(PhaseDuration)

	go func() {
		defer func() {
			m.onRoundComplete <- Round{
				Guesses: p.guesses,
			}
			close(done)
		}()

		for {
			select {
			case guess := <-p.onPlayerGuess:
				p.guesses[guess.PlayerID] = guess
				if len(p.guesses) == len(m.Sessions) {
					return
				}
			case <-time.After(PhaseDuration):
				return
			}
		}
	}()

	return done
}
