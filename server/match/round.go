package match

type Guesses map[string]PlayerGuess

type Round struct {
	Guesses Guesses
}
