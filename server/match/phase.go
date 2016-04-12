package match

import "log"

type Phase func(*Match) Phase

func JoinPhase(m *Match) Phase {
	for {
		select {
		case id := <-m.playerConnect:
			log.Println("JoinPhase player joined", id)
			if len(m.Sessions) != m.Capacity {
				log.Printf("Match #%s still waiting for %d players to connect.", m.ID, m.Capacity-len(m.Sessions))
			} else {
				log.Printf("Match #%s capacity reached.", m.ID)
			}
		}
	}
}
