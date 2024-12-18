package main

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/charmbracelet/log"
)

type SavedUser struct {
	key             string
	verified        bool
	administrator   bool
	email           string
	username        string
	followers       int
	followed        int
	id              int64
	createdAt       time.Time
	description     sql.NullString
	location        sql.NullString
	birthDate       time.Time
}

func SaveUser(db *sql.DB, publicKey string, username string, email string, birthDate time.Time) (int64, error) {
	log.Info("Saving user to db")
	query := `INSERT INTO users (key, username, email, verified, administrator, birth_date)
			  VALUES ($1, $2, $3, $4, $5, $6)`
	result, err := db.Exec(query, publicKey, username, email, false, false, birthDate)

	if err != nil {
		log.Errorf("failed to insert user: %v", err)
		return 0, fmt.Errorf("failed to insert user: %v", err)
	}

	log.Info("Saved new user")
	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to retrieve last inserted ID: %v", err)
	}

	return id, nil
}

func AcceptUser(db *sql.DB, user SavedUser) error {
	query := `UPDATE users SET verified = true WHERE id = $1`

	result, err := db.Exec(query, user.id)
	if err != nil {
		log.Errorf("failed to update user verification status: %v", err)
		return fmt.Errorf("failed to update user verification status: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Errorf("failed to retrieve affected rows: %v", err)
		return fmt.Errorf("failed to retrieve affected rows: %v", err)
	}

	if rowsAffected == 0 {
		log.Errorf("no user found with username: %s", user.username)
		return fmt.Errorf("no user found with username: %s", user.username)
	}

	return nil
}

func DeleteUser(db *sql.DB, user SavedUser) error {
	query := `DELETE FROM users WHERE id = $1`
	result, err := db.Exec(query, user.id)
	if err != nil {
		log.Errorf("failed to delete user: %v", err)
		return fmt.Errorf("failed to delete user: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Errorf("failed to retrieve affected rows: %v", err)
		return fmt.Errorf("failed to retrieve affected rows: %v", err)
	}

	if rowsAffected == 0 {
		log.Errorf("no user found with username: %s", user.username)
		return fmt.Errorf("no user found with username: %s", user.username)
	}

	return nil
}

func GetUserByUsername(db *sql.DB, username string) (SavedUser, bool) {
	var user SavedUser
	log.Debug("Fetching user from db")
	query := `SELECT id, key, username, email, verified, administrator, followers, followed, created_at, description, location, birth_date FROM users WHERE username = $1`
	err := db.QueryRow(query, username).
		Scan(&user.id, &user.key, &user.username, &user.email, &user.verified, &user.administrator, &user.followers, &user.followed, &user.createdAt, &user.description, &user.location, &user.birthDate)

	if err != nil {
		if err == sql.ErrNoRows {
			log.Debug("No user found")
			return SavedUser{}, false
		}
		log.Errorf("Error while fetching user: %s", err)
		return SavedUser{}, false
	}

	return user, true
}

func CreateUserTable(db *sql.DB) {
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		key TEXT UNIQUE NOT NULL,
		username VARCHAR(50) UNIQUE NOT NULL,
		email VARCHAR(100) NOT NULL,
		verified BOOLEAN NOT NULL,
		administrator BOOLEAN NOT NULL,
		followers INTEGER DEFAULT 0,
		followed INTEGER DEFAULT 0,
		created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
		description VARCHAR(100),
		location VARCHAR(50),
		birth_date TIMESTAMP WITH TIME ZONE NOT NULL
	);`

	_, err := db.Exec(createTableSQL)
	if err != nil {
		log.Fatalf("Failed to create table: %v", err)
	}

	log.Info("Table 'users' created successfully!")
}

func findQuery(db *sql.DB, query string) ([]SavedUser, error) {
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []SavedUser
	for rows.Next() {
		var user SavedUser
		if err := rows.Scan(&user.id, &user.key, &user.username, &user.email, &user.verified, &user.administrator, &user.followers, &user.followed, &user.createdAt); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func GetAllUsers(db *sql.DB) ([]SavedUser, error) {
	query := `SELECT id, key, username, email, verified, administrator, followers, followed, created_at FROM users`
	return findQuery(db, query)
}

func GetUnverifiedUsers(db *sql.DB) ([]SavedUser, error) {
	query := `SELECT id, key, username, email, verified, administrator, followers, followed, created_at FROM users WHERE verified = false`
	return findQuery(db, query)
}

func UpdateUserData(db *sql.DB, user SavedUser, description string, location string) error {
	query := `UPDATE users SET description = $2, location = $3 WHERE id = $1`

	_, err := db.Exec(query, user.id, description, location)
	if err != nil {
		log.Errorf("failed to update user info status: %v", err)
		return fmt.Errorf("failed to update user info status: %v", err)
	}

	return nil
}

func SearchUsers(db *sql.DB, search string) ([]SavedUser, error) {
	query := `SELECT id, key, username, email, verified, administrator, followers, followed, created_at FROM users WHERE  username LIKE '%' || $1 || '%'`
	rows, err := db.Query(query, search)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []SavedUser
	for rows.Next() {
		var user SavedUser
		if err := rows.Scan(&user.id, &user.key, &user.username, &user.email, &user.verified, &user.administrator, &user.followers, &user.followed, &user.createdAt); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}
