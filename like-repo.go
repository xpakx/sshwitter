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

func CreateLikeFunction(db *sql.DB) {
	query := `
	CREATE OR REPLACE PROCEDURE add_like(user_id_param INTEGER, post_id_param INTEGER)
	LANGUAGE plpgsql
	AS $$
	BEGIN
		INSERT INTO likes (user_id, post_id) 
		VALUES (user_id_param, post_id_param);

		UPDATE posts SET likes = likes + 1 WHERE id = post_id_param;
	END;
	$$;`

	_, err := db.Exec(query)
	if err != nil {
		log.Fatalf("Failed to create procedure for inserting likes: %v", err)
	}

	log.Info("Procedure created successfully!")
}

func SaveLike(db *sql.DB, user SavedUser, post Post) error {
	log.Info("Saving like to db")
	query := `CALL add_like($1, $2)`

	_, err := db.Exec(query, user.id, post.id)

	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			if pqErr.Code == "23505" {
				log.Warnf("Already liked: %v", err)
				return fmt.Errorf("Already liked: %v", err)
			}
		}
		log.Errorf("failed to insert like: %v", err)
		return fmt.Errorf("failed to insert like: %v", err)
	}

	log.Info("Saved new like")
	return nil
}
