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

func CreateFollowFunction(db *sql.DB) {
	query := `
	CREATE OR REPLACE PROCEDURE add_follow(user_id_param INTEGER, followed_id_param INTEGER)
	LANGUAGE plpgsql
	AS $$
	BEGIN
		INSERT INTO follows (user_id, followed_id) 
		VALUES (user_id_param, followed_id_param);

		UPDATE users SET followers = followers + 1 WHERE id = followed_id_param;
		UPDATE users SET followed = followed + 1 WHERE id = user_id_param;
	END;
	$$;`

	_, err := db.Exec(query)
	if err != nil {
		log.Fatalf("Failed to create procedure for inserting follows: %v", err)
	}

	log.Info("Procedure created successfully!")
}

func SaveFollow(db *sql.DB, user SavedUser, followed SavedUser) (int64, error) {
	log.Info("Saving follow to db")
	var id int64
	query := `CALL add_follow($1, $2)`

	err := db.QueryRow(query, user.id, followed.id).
		Scan(&id)

	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			if pqErr.Code == "23505" {
				log.Warnf("Already followed: %v", err)
				return 0, fmt.Errorf("Already followed: %v", err)
			}
		}
		log.Errorf("failed to insert follow: %v", err)
		return 0, fmt.Errorf("failed to insert follow: %v", err)
	}

	log.Info("Saved new follow")
	return id, nil
}

func CreateUnfollowFunction(db *sql.DB) {
	query := `
	CREATE OR REPLACE PROCEDURE delete_follow(user_id_param INTEGER, followed_id_param INTEGER)
	LANGUAGE plpgsql
	AS $$
	DECLARE 
	        is_deleted INTEGER;
	BEGIN
		DELETE FROM follows 
		WHERE user_id = user_id_param AND followed_id = followed_id_param;
		GET DIAGNOSTICS is_deleted = ROW_COUNT;

                IF is_deleted > 0 THEN
			UPDATE users 
			SET followers = CASE WHEN followers > 0 THEN followers - 1 ELSE followers END
			WHERE id = followed_id_param;

			UPDATE users 
			SET followed = CASE WHEN followed > 0 THEN followed - 1 ELSE followed END
			WHERE id = user_id_param;
		ELSE
		        RAISE EXCEPTION 'Not following' USING ERRCODE = 'P0001'; 
		END IF;
	END;
	$$;`

	_, err := db.Exec(query)
	if err != nil {
		log.Fatalf("Failed to create procedure for inserting follows: %v", err)
	}

	log.Info("Procedure created successfully!")
}

func DeleteFollow(db *sql.DB, user SavedUser, followed SavedUser) error {
	log.Info("Saving follow to db")
	var id int64
	query := `CALL delete_follow($1, $2)`

	err := db.QueryRow(query, user.id, followed.id).
		Scan(&id)

	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			if pqErr.Code == "P0001" {
				log.Warnf("Not followed: %v", err)
				return fmt.Errorf("Not followed: %v", err)
			}
		}
		log.Errorf("failed to delete follow: %v", err)
		return  fmt.Errorf("failed to delete follow: %v", err)
	}

	log.Info("Deleted follow")
	return nil
}

func CheckFollow(db *sql.DB, user SavedUser, other SavedUser) (bool, error) {
	var exists bool;
	query := `
	SELECT EXISTS (
	    SELECT 1 FROM follows WHERE user_id = $1 AND followed_id = $2
	);`

	err := db.QueryRow(query, user.id, other.id).
		Scan(&exists)

	if err != nil {
		log.Errorf("Error while fetching follows: %v", err)
		return  false, fmt.Errorf("Error while fetching follows: %v", err)
	}

	return exists, nil
}
