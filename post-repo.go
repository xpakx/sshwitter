package main

import (
	"database/sql"

	"github.com/charmbracelet/log"
)

type Post struct {
	id              int64
	userId          int64
	content         string
}

func CreatePostTable(db *sql.DB) {
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS posts (
		id SERIAL PRIMARY KEY,
		content TEXT NOT NULL,
		user_id INTEGER REFERENCES users(id)
	);`

	_, err := db.Exec(createTableSQL)
	if err != nil {
		log.Fatalf("Failed to create table: %v", err)
	}

	log.Info("Table 'posts' created successfully!")
}
