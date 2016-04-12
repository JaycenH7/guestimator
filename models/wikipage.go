package models

import (
	"net/url"
	"unicode"

	"gopkg.in/pg.v4"
)

const (
	BaseQueryUrl = "http://en.wikipedia.org/w/api.php?action=query&prop=extracts&format=json&exintro=&explaintext=&titles="
)

type Wikipage struct {
	ID        int        `json:"pageid"`
	Title     string     `json:"title"`
	Extract   string     `json:"extract,omitempty" sql:"-"`
	Questions []Question `json:"-"`
}

func WikipageUrl(title string) string {
	return BaseQueryUrl + url.QueryEscape(title)
}

func CreateWikipage(db *pg.DB, w *Wikipage) error {
	return db.Create(w)
}

func GetWikipage(db *pg.DB, id int) (*Wikipage, error) {
	w := Wikipage{}
	err := db.Model(&w).Where("id = ?", id).Select()
	return &w, err
}

func (w *Wikipage) ExtractQuestions() []Question {
	questions := make([]Question, 0)
	for _, sentence := range w.extractSentences() {
		if q := ParseQuestion(sentence); q != nil {
			q.WikipageID = w.ID
			questions = append(questions, *q)
		}
	}
	return questions
}

func (w Wikipage) extractSentences() []string {
	sentences := make([]string, 0)
	start := -1
	len := len(w.Extract)

	for i, c := range w.Extract {
		if start == -1 {
			if !unicode.IsSpace(c) {
				start = i
			}
		} else if c == '.' {
			if i+1 < len && !unicode.IsSpace(rune(w.Extract[i+1])) {
				continue
			} else if i+2 < len && !unicode.IsUpper(rune(w.Extract[i+2])) {
				continue
			}
			sentences = append(sentences, w.Extract[start:i+1])
			start = -1
		}
	}

	return sentences
}
