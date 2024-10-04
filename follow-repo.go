package main

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/charmbracelet/log"
)

type Follow struct {
	id              int64
	userId          int64
	followedId      int64
	followedAt      time.Time
}

func CreateFollowTable(db *sql.DB) {
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS follows (
		id SERIAL PRIMARY KEY,
		user_id INTEGER REFERENCES users(id),
		followed_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
		followed_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
	);`

	_, err := db.Exec(createTableSQL)
	if err != nil {
		log.Fatalf("Failed to create table: %v", err)
	}

	log.Info("Table 'follows' created successfully!")
}
