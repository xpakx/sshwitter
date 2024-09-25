package main

import (
	"database/sql"
	"fmt"

	"github.com/charmbracelet/log"
)

type SavedUser struct {
	key             string
	verified        bool
	administrator   bool
	email           string
	username        string
	id              int64
}

// TODO: move to db
var users = map[string]SavedUser{
	"test" : { 
		key: "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAACAQDQ2EOmuzeJo6LN4jg7Vuvf0BSjrqDMSguc3Z6i0zkVURz5Bb61BZlG7PpyqOQ1aeISNfoMSpVwx1fQytIvQBfpF6OX+XMI/wzgOEvbPNeQ0RHsPJjq0x6wMLNEWpPl5f07pDgBlWB8IkBTKvSZQje/WsEwDvUnFRrWcC8PHs2H/WRpm+wagg9T5N6jDqlC711DJEWIyKwl744QHK4NBnyXHfK+0pW/JfhEelyQ+bTVfWNDu9V5uZI69hiKZNs4UANhAoUEhhZIy60ZHho6Zn8JkZkjORMwGi/hi8lUaIDYXXcqKGqKQdU2HU5NgpWVO3/w7KRQceegDiMO5Aa/yMEtdVi0B2NmUGVTZcCEwkqWbACqG5r23AmgrMX/Hh8L/9Z1nFwnxCY2bUd29DQI1q7GzTwYIxNi9y7/8H5+gmU6Yn3Wm5mUjpxWLF9QbU0fOFNZ/WO1h3rRYCwoouJ4ixWuCM6BLcBuuEutx24mjBaO3x0p68XJ8rxMuvS/n9TwTywPfeDS5Yft1hHovyRt1vSAHLxd8eSP65vJHJwsYAL8psGbm68CyYnzf8D4CPSJh4DSCQRzNnfFjYozX9QuXAhPtJkjPI7w6mJyPmjUaDB+sOkolIqIdF0jBXuaB/Hv/03H3ul5+SqpB0s37Wh0rwI2ORX0Ct45pYjj78WtAkukEQ==",
		verified: true,
		administrator: true,
		email: "",
		username: "test",
	},
}

func SaveUser(db *sql.DB, publicKey string, username string, email string) (int64, error) {
	query := `INSERT INTO users (key, username, email, verified, administrator)
			  VALUES ($1, $2, $3, $4, $5)`
	result, err := db.Exec(query, publicKey, username, email, false, false)

	if err != nil {
		return 0, fmt.Errorf("failed to insert user: %v", err)
	}

	log.Info("Saved new user")
	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to retrieve last inserted ID: %v", err)
	}

	return id, nil
}

func AcceptUser(user SavedUser) {
	user.verified = true;
	users[user.username] = user
}

func DeleteUser(user SavedUser) {
	delete(users, user.username)
}

func GetUserByUsername(db *sql.DB, username string) (SavedUser, bool) {
	var user SavedUser
	log.Debug("Fetching user from db")
	query := `SELECT id, key, username, email, verified, administrator FROM users WHERE username = $1`
	err := db.QueryRow(query, username).
		Scan(&user.id, &user.key, &user.username, &user.email, &user.verified, &user.administrator)

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
		key VARCHAR(50) UNIQUE NOT NULL,
		username VARCHAR(50) UNIQUE NOT NULL,
		email VARCHAR(100) NOT NULL,
		verified BOOLEAN NOT NULL,
		administrator BOOLEAN NOT NULL
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
		if err := rows.Scan(&user.id, &user.key, &user.username, &user.email, &user.verified, &user.administrator); err != nil {
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
	query := `SELECT id, key, username, email, verified, administrator FROM users`
	return findQuery(db, query)
}

func GetUnverifiedUsers(db *sql.DB) ([]SavedUser, error) {
	query := `SELECT id, key, username, email, verified, administrator FROM users WHERE verified = false`
	return findQuery(db, query)
}
