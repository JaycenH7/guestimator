package fixtures

import "github.com/mrap/guestimator/models"

var basicQuestion = models.Question{
	FullText:  "He will turn 26 in 2016.",
	Positions: []int{13, 14, 19, 22},
}

var basicWikipage = models.Wikipage{
	ID:      1,
	Title:   "My Wikipage",
	Extract: "Extract text",
	Questions: []models.Question{
		basicQuestion,
	},
}

func Question() models.Question {
	q := basicQuestion

	wp := Wikipage()
	q.Wikipage = &wp

	return q
}

func Wikipage() models.Wikipage {
	return basicWikipage
}
