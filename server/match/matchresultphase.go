package match

type MatchResultPhase struct {
}

func NewMatchResultPhase() *MatchResultPhase {
	p := new(MatchResultPhase)
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

func (p *MatchResultPhase) Run(m *Match) <-chan struct{} {
	done := make(chan struct{})

	go func() {
		defer close(done)

		m.Hub.Close()
	}()

	return done
}
