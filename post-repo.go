package main

import (
	"database/sql"
	"fmt"

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

func SavePost(db *sql.DB, user SavedUser, content string) (int64, error) {
	log.Info("Saving post to db")
	query := `INSERT INTO posts (content, userId)
			  VALUES ($1, $2)`
	result, err := db.Exec(query, content, user.id)

	if err != nil {
		log.Errorf("failed to insert post: %v", err)
		return 0, fmt.Errorf("failed to insert post: %v", err)
	}

	log.Info("Saved new post")
	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to retrieve last inserted ID: %v", err)
	}

	return id, nil
}
