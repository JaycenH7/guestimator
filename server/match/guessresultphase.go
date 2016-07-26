package match

import "time"

type GuessResultPhase struct {
	endTime time.Time
}

func NewGuessResultPhase() *GuessResultPhase {
	p := &GuessResultPhase{}
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

func (p GuessResultPhase) TimeRemaining() time.Duration {
	return p.endTime.Sub(time.Now())
}

func (p *GuessResultPhase) Run(m *Match) <-chan struct{} {
	done := make(chan struct{})
	p.endTime = time.Now().Add(PhaseDuration)

	go func() {
		defer func() {
			close(done)
		}()

		// TODO: We should wait until all players are ready for the next phase
		select {
		case <-time.After(PhaseDuration):
			return
		}
	}()

	return done
}
