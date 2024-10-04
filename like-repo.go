package main

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/charmbracelet/log"
	"github.com/lib/pq"
)

type Like struct {
	id              int64
	userId          int64
	postId          int64
	likedAt         time.Time
}

func CreateLikeTable(db *sql.DB) {
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS likes (
		id SERIAL PRIMARY KEY,
		user_id INTEGER REFERENCES users(id),
		post_id INTEGER REFERENCES posts(id) ON DELETE CASCADE,
		liked_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
		CONSTRAINT unique_like UNIQUE (user_id, post_id)
	);`

	_, err := db.Exec(createTableSQL)
	if err != nil {
		log.Fatalf("Failed to create table: %v", err)
	}

	log.Info("Table 'likes' created successfully!")
}
