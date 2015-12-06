package models

import (
	"log"

	"github.com/mrap/stringutil"
	"gopkg.in/pg.v3"
)

const (
	BaseQueryUrl = "http://en.wikipedia.org/w/api.php?action=query&prop=extracts&format=json&exintro=&explaintext=&titles="
)

type WikiPage struct {
	PageID  int    `json:"pageid"`
	Title   string `json:"title"`
	Extract string `json:"extract"`
}

func WikiPageUrl(title string) string {
	encoded, err := stringutil.UrlEncoded(title)
	if err != nil {
		log.Println("Error making wiki page url", err)
	}
	return BaseQueryUrl + encoded
}

func CreateWikiPage(db *pg.DB, w *WikiPage) error {
	_, err := db.QueryOne(w, `
		INSERT INTO wiki_pages (page_id, title)
		VALUES (?page_id, ?title)
		RETURNING page_id
	`, w)

	return err
}

func GetWikiPage(db *pg.DB, page_id int) (*WikiPage, error) {
	w := &WikiPage{}
	_, err := db.QueryOne(w, `SELECT * FROM wiki_pages WHERE page_id = ?`, page_id)
	return w, err
}

func (w *WikiPage) ExtractQuestions() []Question {
	questions := make([]Question, 0)
	for _, sentence := range w.extractSentences() {
		if q := ParseQuestion(sentence); q != nil {
			q.PageID = w.PageID
			questions = append(questions, *q)
		}
	}
	return questions
}

func (w WikiPage) extractSentences() []string {
	sentences := make([]string, 0)

	i, start := 0, -1
	for _, c := range w.Extract {
		if start == -1 {
			if c != ' ' {
				start = i
			}
		} else if c == '.' {
			sentences = append(sentences, w.Extract[start:i+1])
			start = -1
		}
		i++
	}

	return sentences
}
