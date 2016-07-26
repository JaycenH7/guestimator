package match

import "time"

type MatchResultPhase struct {
	endTime time.Time
}

func NewMatchResultPhase() *MatchResultPhase {
	p := &MatchResultPhase{}
	return p
}

func (p MatchResultPhase) Label() string {
	return "MatchResult"
}

func (p *MatchResultPhase) OnPlayerGuess(guess PlayerGuess) {
}

func (p *MatchResultPhase) OnPlayerConnect(id string) {
}

func (p *MatchResultPhase) OnPlayerDisconnect(id string) {
}

func (p MatchResultPhase) TimeRemaining() time.Duration {
	return p.endTime.Sub(time.Now())
}

func (p *MatchResultPhase) Run(m *Match) <-chan struct{} {
	done := make(chan struct{})
	p.endTime = time.Now().Add(PhaseDuration)

	go func() {
		defer func() {
			m.Hub.Close()
			close(done)
		}()

		select {
		case <-time.After(PhaseDuration):
			return
		}
	}()

	return done
}
