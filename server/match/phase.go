package match

import "log"

type PhaseType int

const (
	JoinPhaseType PhaseType = iota
	AnswerPhaseType
)

type Phase func(*Match) Phase

func JoinPhase(m *Match) Phase {
	m.currentPhaseType = JoinPhaseType

	for {
		select {
		case playerID := <-m.playerConnect:
			log.Println("JoinPhase player joined", playerID)
			if len(m.Sessions) == m.Capacity {
				log.Println("JoinPhase capacity reached")
				return AnswerPhase
			}
		}
	}
}

func AnswerPhase(m *Match) Phase {
	m.phaseType = AnswerPhaseType

	for {
		select {
		default:
		}
	}

	return nil
}
