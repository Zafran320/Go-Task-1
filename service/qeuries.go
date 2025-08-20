package service

import (
	"backend-auth/models"
	"database/sql"
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func HandleUserQuery(db *sql.DB, user models.User) error {
	switch {
	default:
		var exists bool
		err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE username = ?)", user.Username).Scan(&exists)
		if err != nil {
			return errors.New("db_error")
		}
		if exists {
			return errors.New("user_exists")
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.PasswordHash), bcrypt.DefaultCost)
		if err != nil {
			return errors.New("hash_error")
		}

		_, err = db.Exec("INSERT INTO users (username, password_hash, time) VALUES (?, ?, ?)",
			user.Username, string(hashedPassword), time.Now())
		if err != nil {
			return errors.New("insert_error")
		}

		return nil
	}
}
