package match

import (
	"log"
	"time"
)

type JoinPhase struct {
	onPlayerConnect chan string
	endTime         time.Time
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

func (p JoinPhase) TimeRemaining() time.Duration {
	return p.endTime.Sub(time.Now())
}

func (p *JoinPhase) Run(m *Match) <-chan struct{} {
	done := make(chan struct{})
	p.endTime = time.Now().Add(PhaseDuration)

	go func() {
		defer func() {
			close(done)
		}()

		for {
			select {
			case playerID := <-p.onPlayerConnect:
				log.Println("JoinPhase player joined", playerID)
				if len(m.Sessions) == m.Capacity {
					log.Println("JoinPhase capacity reached")
					return
				}
			case <-time.After(PhaseDuration):
				return
			}
		}
	}()

	return done
}
