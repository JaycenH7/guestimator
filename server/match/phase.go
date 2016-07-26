package match

import "time"

type Phase interface {
	EventListener
	Run(m *Match) <-chan struct{}
	Label() string
	TimeRemaining() time.Duration
}

var PhaseDuration = 2 * time.Second

// PhaseQueue is a queue of phases.
type PhaseQueue struct {
	phases []Phase
}

// NewPhaseQueue returns a new PhaseQueue populated with the provided phases.
func NewPhaseQueue(phases ...Phase) *PhaseQueue {
	return &PhaseQueue{
		phases: phases,
	}
}

// Size returns the number of phases remaining in the queue.
func (q PhaseQueue) Size() int {
	return len(q.phases)
}

// Next removes and returns the head of the queue or nil if queue is empty.
func (q *PhaseQueue) Next() Phase {
	len := len(q.phases)
	if len == 0 {
		return nil
	}

	p := q.phases[0]
	q.phases = q.phases[1:]

	return p
}

// Append adds a phase to the end of the queue.
func (q *PhaseQueue) Append(p Phase) {
	q.phases = append(q.phases, p)
}

// Prepend inserts a phase to the head of the queue.
func (q *PhaseQueue) Prepend(p Phase) {
	q.phases = append([]Phase{p}, q.phases...)
}
