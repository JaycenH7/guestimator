package match

type AnswerPhase struct {
	onPlayerConnect <-chan string
}

func (p AnswerPhase) Label() string {
	return "Answer"
}

func (p *AnswerPhase) OnPlayerConnect(id string) {
}

func (p *AnswerPhase) OnPlayerDisconnect(id string) {
}

func (p *AnswerPhase) Run(m *Match) <-chan struct{} {
	done := make(chan struct{})
	go func() {
		// TODO: implement
		close(done)
	}()

	return done
}
