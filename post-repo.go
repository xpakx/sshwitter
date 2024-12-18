package main

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/charmbracelet/log"
)

type Post struct {
	id              int64
	userId          int64
	content         string
	likes           int
	replies         int
	username        string
	createdAt       time.Time
	liked           bool
	parentId        sql.NullInt64
}

func CreatePostTable(db *sql.DB) {
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS posts (
		id SERIAL PRIMARY KEY,
		content TEXT NOT NULL,
		user_id INTEGER REFERENCES users(id),
		likes INTEGER DEFAULT 0,
		replies INTEGER DEFAULT 0,
		created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
		parent_id INTEGER REFERENCES posts(id)
	);`

	_, err := db.Exec(createTableSQL)
	if err != nil {
		log.Fatalf("Failed to create table: %v", err)
	}

	log.Info("Table 'posts' created successfully!")
}

func SavePost(db *sql.DB, user SavedUser, content string) (int64, error) {
	log.Info("Saving post to db")
	var id int64
	query := `INSERT INTO posts (content, user_id)
        VALUES ($1, $2)
	RETURNING id`
	err := db.QueryRow(query, content, user.id).
		Scan(&id)

	if err != nil {
		log.Errorf("failed to insert post: %v", err)
		return 0, fmt.Errorf("failed to insert post: %v", err)
	}

	log.Info("Saved new post")
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

func FindUserPosts(db *sql.DB, user SavedUser, viewer SavedUser) ([]Post, error) {
	query := `
	SELECT p.id, p.content, p.user_id, p.created_at, u.username, p.likes, p.replies,
	       CASE WHEN l.user_id IS NOT NULL THEN TRUE ELSE FALSE END AS liked
	FROM posts p
	LEFT JOIN likes l ON p.id = l.post_id AND l.user_id = $2
	LEFT JOIN users u
	ON p.user_id = u.id
	WHERE p.user_id = $1
	ORDER BY p.created_at DESC`
	rows, err := db.Query(query, user.id, viewer.id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []Post
	for rows.Next() {
		var post Post
		if err := rows.Scan(&post.id, &post.content, &post.userId, &post.createdAt, &post.username, &post.likes, &post.replies, &post.liked); err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}

func FindAllPosts(db *sql.DB, viewer SavedUser) ([]Post, error) {
	query := `
	SELECT p.id, p.content, p.user_id, p.created_at, u.username, p.likes, p.replies,
	       CASE WHEN l.user_id IS NOT NULL THEN TRUE ELSE FALSE END AS liked
	FROM posts p
	LEFT JOIN likes l ON p.id = l.post_id AND l.user_id = $1
	LEFT JOIN users u
	ON p.user_id = u.id
	ORDER BY p.created_at DESC`
	rows, err := db.Query(query, viewer.id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []Post
	for rows.Next() {
		var post Post
		if err := rows.Scan(&post.id, &post.content, &post.userId, &post.createdAt, &post.username, &post.likes, &post.replies, &post.liked); err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}

func FindFollowedPosts(db *sql.DB, viewer SavedUser) ([]Post, error) {
	query := `
	SELECT p.id, p.content, p.user_id, p.created_at, u.username, p.likes, p.replies,
	       CASE WHEN l.user_id IS NOT NULL THEN TRUE ELSE FALSE END AS liked
	FROM posts p
	INNER JOIN follows f ON f.followed_id = p.user_id AND f.user_id = $1
	LEFT JOIN likes l ON p.id = l.post_id AND l.user_id = $1
	LEFT JOIN users u ON p.user_id = u.id
	WHERE f.user_id = $1
	ORDER BY p.created_at DESC`
	rows, err := db.Query(query, viewer.id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []Post
	for rows.Next() {
		var post Post
		if err := rows.Scan(&post.id, &post.content, &post.userId, &post.createdAt, &post.username, &post.likes, &post.replies, &post.liked); err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}

func FindLikedPosts(db *sql.DB, viewer SavedUser) ([]Post, error) {
	query := `
	SELECT p.id, p.content, p.user_id, p.created_at, u.username, p.likes, p.replies,
	       CASE WHEN l.user_id IS NOT NULL THEN TRUE ELSE FALSE END AS liked
	FROM posts p
	LEFT JOIN likes l ON p.id = l.post_id AND l.user_id = $1
	LEFT JOIN users u ON p.user_id = u.id
	WHERE l.user_id IS NOT NULL
	ORDER BY p.created_at DESC`
	rows, err := db.Query(query, viewer.id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []Post
	for rows.Next() {
		var post Post
		if err := rows.Scan(&post.id, &post.content, &post.userId, &post.createdAt, &post.username, &post.likes, &post.replies, &post.liked); err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}

func GetPostById(db *sql.DB, id int64, username string) (Post, bool) {
	var post Post
	log.Debug("Fetching post from db")
	query := `
	SELECT p.id, p.content, p.user_id, p.created_at, u.username, p.likes, p.replies,
	       CASE WHEN l.user_id IS NOT NULL THEN TRUE ELSE FALSE END AS liked,
	       p.parent_id
	FROM posts p
	LEFT JOIN likes l ON p.id = l.post_id 
	LEFT JOIN users u ON p.user_id = u.id
	WHERE u.username = $1
	AND p.id = $2`

	err := db.QueryRow(query, username, id).
		Scan(&post.id, &post.content, &post.userId, &post.createdAt, &post.username, &post.likes, &post.replies, &post.liked, &post.parentId)

	if err != nil {
		if err == sql.ErrNoRows {
			log.Debug("No post found")
			return Post{}, false
		}
		log.Errorf("Error while fetching post: %s", err)
		return Post{}, false
	}

	return post, true
}

func FindReplies(db *sql.DB, id int64, viewer SavedUser) ([]Post, error) {
	query := `
	SELECT p.id, p.content, p.user_id, p.created_at, u.username, p.likes, p.replies,
	       CASE WHEN l.user_id IS NOT NULL THEN TRUE ELSE FALSE END AS liked
	FROM posts p
	LEFT JOIN likes l ON p.id = l.post_id AND l.user_id = $1
	LEFT JOIN users u
	ON p.user_id = u.id
	WHERE p.parent_id = $2
	ORDER BY p.created_at DESC`
	rows, err := db.Query(query, viewer.id, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []Post
	for rows.Next() {
		var post Post
		if err := rows.Scan(&post.id, &post.content, &post.userId, &post.createdAt, &post.username, &post.likes, &post.replies, &post.liked); err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}

func ReplyToPost(db *sql.DB, user SavedUser, post Post, content string) (int64, error) {
	log.Info("Saving reply to db")
	var id int64
	query := `SELECT add_reply($1, $2, $3)`
	err := db.QueryRow(query, user.id, post.id, content).
		Scan(&id)

	if err != nil {
		log.Errorf("failed to insert post: %v", err)
		return 0, fmt.Errorf("failed to insert post: %v", err)
	}

	log.Info("Saved new reply")
	return id, nil
}

func CreateReplyFunction(db *sql.DB) {
	query := `
	CREATE OR REPLACE FUNCTION add_reply(user_id_param INTEGER, post_id_param INTEGER, content_param TEXT) RETURNS INTEGER
	LANGUAGE plpgsql
	AS $$
	DECLARE 
	        new_id INTEGER;
	BEGIN
		INSERT INTO posts (content, user_id, parent_id) 
		VALUES (content_param, user_id_param, post_id_param) RETURNING id INTO new_id;
		UPDATE posts SET replies = replies + 1 WHERE id = post_id_param;
	        RETURN new_id;
	END;
	$$;`

	_, err := db.Query(query)
	if err != nil {
		log.Fatalf("Failed to create procedure for inserting replies: %v", err)
	}

	log.Info("Procedure created successfully!")
}

func FindAllRepliesToUserPosts(db *sql.DB, viewer SavedUser) ([]Post, error) {
	query := `
	SELECT p.id, p.content, p.user_id, p.created_at, u.username, p.likes, p.replies,
	       CASE WHEN l.user_id IS NOT NULL THEN TRUE ELSE FALSE END AS liked
	FROM posts p
	LEFT JOIN likes l ON p.id = l.post_id AND l.user_id = $1
	LEFT JOIN users u ON p.user_id = u.id
	RIGHT JOIN posts pa ON p.parent_id = pa.id
	WHERE p.user_id = $1
	ORDER BY p.created_at DESC`
	rows, err := db.Query(query, viewer.id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []Post
	for rows.Next() {
		var post Post
		if err := rows.Scan(&post.id, &post.content, &post.userId, &post.createdAt, &post.username, &post.likes, &post.replies, &post.liked); err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}
