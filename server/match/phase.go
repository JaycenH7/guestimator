package match

type Phase interface {
	EventListener
	Run(m *Match) <-chan struct{}
	Label() string
}
