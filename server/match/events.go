package match

type EventListener interface {
	OnPlayerConnect(id string)
	OnPlayerDisconnect(id string)
	OnPlayerGuess(guess PlayerGuess)
}
