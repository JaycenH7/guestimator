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
