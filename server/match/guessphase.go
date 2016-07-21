package match

type GuessPhase struct {
	onPlayerGuess chan PlayerGuess
	guesses       Guesses
}

func NewGuessPhase() *GuessPhase {
	p := &GuessPhase{
		guesses:       make(Guesses),
		onPlayerGuess: make(chan PlayerGuess),
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

func (p *GuessPhase) Run(m *Match) <-chan struct{} {
	done := make(chan struct{})

	go func() {
		defer close(done)

		for {
			select {
			case guess := <-p.onPlayerGuess:
				p.guesses[guess.PlayerID] = guess
				if len(p.guesses) == len(m.Sessions) {
					m.onRoundComplete <- Round{
						Guesses: p.guesses,
					}
					return
				}
			}
		}
	}()

	return done
}
