package db

import (
	"_/pkg/chgkdb"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

func InitializeDatabase(dbPath string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	query := `
    CREATE TABLE IF NOT EXISTS questions (
        id INTEGER PRIMARY KEY,
        championship TEXT,
        tour TEXT,
		number INTEGER,
        question TEXT,
        answer TEXT,
        source TEXT,
        comment TEXT
    );`

	_, err = db.Exec(query)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func InsertQuestion(db *sql.DB, question chgkdb.Question) error {
	query := `
    INSERT INTO questions (championship, tour, number, question, answer, source, comment)
    VALUES (?, ?, ?, ?, ?, ?)`

	_, err := db.Exec(query, question.Championship, question.Tour, question.Number, question.Question, question.Answer, question.Source, question.Comment)
	return err
}

func InsertQuestions(db *sql.DB, questions []chgkdb.Question) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare(`
    INSERT INTO questions (championship, tour, number, question, answer, source, comment)
    VALUES (?, ?, ?, ?, ?, ?, ?)`)

	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, q := range questions {
		_, err := stmt.Exec(q.Championship, q.Tour, q.Number, q.Question, q.Answer, q.Source, q.Comment)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}
