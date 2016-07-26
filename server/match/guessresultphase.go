package match

import "time"

type GuessResultPhase struct {
}

func NewGuessResultPhase() *GuessResultPhase {
	p := new(GuessResultPhase)
	return p
}

func (p GuessResultPhase) Label() string {
	return "GuessResult"
}

func (p *GuessResultPhase) OnPlayerGuess(guess PlayerGuess) {
}

func (p *GuessResultPhase) OnPlayerConnect(id string) {
}

func (p *GuessResultPhase) OnPlayerDisconnect(id string) {
}

func (p *GuessResultPhase) Run(m *Match) <-chan struct{} {
	done := make(chan struct{})

	go func() {
		defer close(done)

		// TODO: We should wait until all players are ready for the next phase
		time.Sleep(500 * time.Millisecond)
	}()

	return done
}
