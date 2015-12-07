package models

import "gopkg.in/pg.v3"

type Question struct {
	Id        int
	FullText  string
	Positions []int
	PageID    int
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
		isDigit = runeIsDigit(c)
		if start == -1 {
			if isDigit {
				start = i
			}
		} else if !isDigit {
			// Capture trailing % symbols and decimals too
			if c != '%' && !(c == '.' && i < len(s)-1 && runeIsDigit(rune(s[i+1]))) {
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

const (
	digit0 = '0'
	digit1 = '1'
	digit2 = '2'
	digit3 = '3'
	digit4 = '4'
	digit5 = '5'
	digit6 = '6'
	digit7 = '7'
	digit8 = '8'
	digit9 = '9'
)

func runeIsDigit(r rune) bool {
	return r == digit0 ||
		r == digit1 ||
		r == digit2 ||
		r == digit3 ||
		r == digit4 ||
		r == digit5 ||
		r == digit6 ||
		r == digit7 ||
		r == digit8 ||
		r == digit9
}
