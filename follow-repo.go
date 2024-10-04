package main

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/charmbracelet/log"
	"github.com/lib/pq"
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
		followed_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
		CONSTRAINT unique_follow UNIQUE (user_id, followed_id)
	);`

	_, err := db.Exec(createTableSQL)
	if err != nil {
		log.Fatalf("Failed to create table: %v", err)
	}

	log.Info("Table 'follows' created successfully!")
}

func SaveFollow(db *sql.DB, user SavedUser, followed SavedUser) (int64, error) {
	log.Info("Saving follow to db")
	var id int64
	query := `INSERT INTO follows (user_id, followed_id)
        VALUES ($1, $2)
	RETURNING id`

	err := db.QueryRow(query, user.id, followed.id).
		Scan(&id)

	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			if pqErr.Code == "23505" {
				log.Warnf("Already followed: %v", err)
				return 0, fmt.Errorf("Already followed: %v", err)
			}
		}
		log.Errorf("failed to insert post: %v", err)
		return 0, fmt.Errorf("failed to insert follow: %v", err)
	}

	log.Info("Saved new follow")
	return id, nil
}
