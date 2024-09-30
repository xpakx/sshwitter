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

func DeletePost(db *sql.DB, post Post, user SavedUser) error {
	if user.id != post.userId {
		log.Errorf("Cannot delete, user %s tried to delete post %d", user.username, post.id)
		return fmt.Errorf("Cannot delete")
	}

	query := `DELETE FROM posts WHERE id = $1`
	result, err := db.Exec(query, post.id)
	if err != nil {
		log.Errorf("failed to delete post: %v", err)
		return fmt.Errorf("failed to delete post: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Errorf("failed to retrieve affected rows: %v", err)
		return fmt.Errorf("failed to retrieve affected rows: %v", err)
	}

	if rowsAffected == 0 {
		log.Errorf("no post found with id: %d", post.id)
		return fmt.Errorf("no post found with id: %d", post.id)
	}

	return nil
}
