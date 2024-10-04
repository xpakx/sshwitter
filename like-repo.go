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

func CreateUnlikeFunction(db *sql.DB) {
	query := `
	CREATE OR REPLACE PROCEDURE delete_like(user_id_param INTEGER, post_id_param INTEGER)
	LANGUAGE plpgsql
	AS $$
	DECLARE 
	        is_deleted INTEGER;
	BEGIN
		DELETE FROM likes 
		WHERE user_id = user_id_param AND post_id = post_id_param;
		GET DIAGNOSTICS is_deleted = ROW_COUNT;

                IF is_deleted > 0 THEN
			UPDATE posts 
			SET likes = CASE WHEN likes > 0 THEN likes - 1 ELSE likes END
			WHERE id = post_id_param;

		ELSE
		        RAISE EXCEPTION 'Not liking' USING ERRCODE = 'P0001'; 
		END IF;
	END;
	$$;`

	_, err := db.Exec(query)
	if err != nil {
		log.Fatalf("Failed to create procedure for inserting follows: %v", err)
	}

	log.Info("Procedure created successfully!")
}

func DeleteLike(db *sql.DB, user SavedUser, post Post) error {
	log.Info("Deleting like in db")
	query := `CALL delete_like($1, $2)`

	_, err := db.Exec(query, user.id, post.id)

	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			if pqErr.Code == "P0001" {
				log.Warnf("Not liking: %v", err)
				return fmt.Errorf("Not liking: %v", err)
			}
		}
		log.Errorf("failed to delete like: %v", err)
		return  fmt.Errorf("failed to delete like: %v", err)
	}

	log.Info("Deleted like")
	return nil
}
