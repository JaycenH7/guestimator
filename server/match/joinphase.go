package match

import "log"

type JoinPhase struct {
	onPlayerConnect chan string
}

func NewJoinPhase() *JoinPhase {
	p := &JoinPhase{
		onPlayerConnect: make(chan string),
	}
	return p
}

func (p JoinPhase) Label() string {
	return "Join"
}

func (p *JoinPhase) OnPlayerGuess(guess PlayerGuess) {
}

func (p *JoinPhase) OnPlayerConnect(id string) {
	log.Println("JoinPhase OnPlayerConnect", id)
	p.onPlayerConnect <- id
}

func (p *JoinPhase) OnPlayerDisconnect(id string) {
}

func (p *JoinPhase) Run(m *Match) <-chan struct{} {
	done := make(chan struct{})

	go func() {
		defer close(done)

		for {
			select {
			case playerID := <-p.onPlayerConnect:
				log.Println("JoinPhase player joined", playerID)
				if len(m.Sessions) == m.Capacity {
					log.Println("JoinPhase capacity reached")
					return
				}
			}
		}
	}()

	return done
}
