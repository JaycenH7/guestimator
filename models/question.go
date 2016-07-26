//go:generate easyjson $GOFILE
package models

import (
	"bytes"
	"errors"
	"fmt"
	"strconv"
	"unicode"

	"gopkg.in/pg.v4"
)

//easyjson:json
type Question struct {
	Id         int    `json:"id"`
	FullText   string `json:"full_text"`
	Positions  []int  `json:"pos" pg:",array"`
	WikipageID int
	Wikipage   *Wikipage `json:"wikipage,omitempty"`
	Answer     *Answer   `json:"answer,omitempty"`
}

//easyjson:json
type Answer struct {
	Exact     float64 `json:"exact"`
	Formatted string  `json:"formatted"`
}

func (q Question) String() string {
	return q.FullText
}

func (q *Question) PopulateAnswer() error {
	answer, err := q.FirstAnswer()
	if err != nil {
		return err
	}
	q.Answer = &answer
	return nil
}

func (q Question) AnswerAt(pos int) (Answer, error) {
	endPos := -1
	for i, p := range q.Positions {
		if p == pos && i < len(q.Positions)-1 {
			endPos = q.Positions[i+1]
			break
		}
	}

	if endPos == -1 {
		return Answer{}, errors.New(fmt.Sprintf("Could not find answer at position: %d", pos))
	}

	buf := new(bytes.Buffer)
	i := -1
	for _, c := range q.FullText {
		i++
		if i >= pos && i <= endPos {
			buf.WriteRune(c)
		}
		if i == endPos {
			break
		}
	}

	formatted := buf.String()
	exact, err := strconv.ParseFloat(formatted, 64)
	answer := Answer{
		Exact:     exact,
		Formatted: formatted,
	}
	return answer, err
}

func (q Question) FirstAnswer() (Answer, error) {
	return q.AnswerAt(q.Positions[0])
}

func (q Question) SansAnswers() Question {
	return q.SansAnswersAt(-1)
}

func (q Question) SansAnswersAt(answerPos int) Question {
	q.FullText = q.FullTextSansAnswersAt(answerPos)
	q.Answer = nil

	wp := *q.Wikipage
	wp.Extract = ""
	wp.Questions = nil
	q.Wikipage = &wp

	return q
}

func (q Question) FullTextSansAnswers() string {
	return q.FullTextSansAnswersAt(-1)
}

func (q Question) FullTextSansAnswersAt(answerPos int) string {
	if len(q.Positions) < 2 {
		return q.FullText
	}

	i := -1
	buf := new(bytes.Buffer)
	posI := 0
	startPos := q.Positions[posI]
	endPos := q.Positions[posI+1]

	for _, c := range q.FullText {
		i++
		if i >= startPos && i <= endPos && (answerPos == -1 || startPos == answerPos) {
			buf.WriteRune('_')
			if i == endPos {
				posI += 2
				if posI < len(q.Positions) {
					startPos = q.Positions[posI]
					endPos = q.Positions[posI+1]
				}
			}
		} else {
			buf.WriteRune(c)
		}
	}

	return buf.String()
}

func CreateQuestion(db *pg.DB, q *Question) error {
	return db.Create(q)
}

func GetQuestion(db *pg.DB, id int) (*Question, error) {
	q := Question{}
	err := db.Model(&q).Where("id = ?", id).Select()
	return &q, err
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
