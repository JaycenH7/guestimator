package models

import (
	"bytes"
	"unicode"

	"gopkg.in/pg.v4"
)

type Question struct {
	Id        int
	FullText  string
	Positions []int `pg:",array"`
	PageID    int
}

func (q Question) String() string {
	return q.FullText
}

func (q Question) SansAnswers() string {
	if len(q.Positions) == 0 {
		return q.FullText
	}

	buf := new(bytes.Buffer)
	posI := 0
	pos := q.Positions[posI]

	for i, c := range q.FullText {
		if i == pos {
			buf.WriteRune('_')
			posI++
			if posI < len(q.Positions) {
				pos = q.Positions[posI]
			}
		} else {
			buf.WriteRune(c)
		}
	}

	return buf.String()
}

func CreateQuestion(db *pg.DB, q *Question) error {
	_, err := db.QueryOne(q, `
		INSERT INTO questions (page_id, full_text, positions)
		VALUES (?page_id, ?full_text, ?positions)
		RETURNING id
	`, q)

	return err
}

func GetQuestion(db *pg.DB, id int) (*Question, error) {
	q := &Question{}
	_, err := db.QueryOne(q, `SELECT * FROM questions WHERE id = ?`, id)
	return q, err
}

func ParseQuestion(s string) *Question {
	var isDigit bool
	positions := make([]int, 0)
	i, start := 0, -1

	for _, c := range s {
		isDigit = unicode.IsDigit(c)
		if start == -1 {
			if isDigit {
				start = i
			}
		} else if !isDigit {
			// Capture trailing % symbols and decimals too
			if c != '%' && !(c == '.' && i < len(s)-1 && unicode.IsDigit(rune(s[i+1]))) {
				positions = append(positions, start, i-1)
				start = -1
			}
		}
		i++
	}

	if len(positions) > 0 {
		return &Question{
			FullText:  s,
			Positions: positions,
		}
	} else {
		return nil
	}
}
