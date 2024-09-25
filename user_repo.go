package main

import (
	"database/sql"

	"github.com/charmbracelet/log"
)

type SavedUser struct {
	key             string
	verified        bool
	administrator   bool
	email           string
	username        string
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

func SaveUser(publicKey string, username string, email string) {
	_, exists := users[username]
	if !exists { // TODO: error; also, it's not thread safe, but map is only a temporary solution anyways
		val := SavedUser {
			key: publicKey,
			verified: false,
			administrator: false,
			username: username,
			email: email,
		}
		users[username] = val
		log.Info("Saved new user")

	} else {
		log.Error("User already exists")
	}
}

func AcceptUser(user SavedUser) {
	user.verified = true;
	users[user.username] = user
}

func DeleteUser(user SavedUser) {
	delete(users, user.username)
}

func GetUser(username string) (SavedUser, bool) {
	user, found := users[username]
	return user, found
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

func GetUnverifiedUsers() []SavedUser {
	var result []SavedUser = make([]SavedUser, 0)
	for _, user := range users {
		if !user.verified {
			result = append(result, user)
		}
	}
	return result
}

