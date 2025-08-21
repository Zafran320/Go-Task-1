package db

import (
	"database/sql"
	"errors"
)

// sign in handler function
func GetPassword(db *sql.DB, username string) (string, error) {
	var storedHash string
	err := db.QueryRow("SELECT password_hash FROM users WHERE username = ?", username).Scan(&storedHash)
	if err == sql.ErrNoRows {
		return "", errors.New("invalid_credentials")
	}
	if err != nil {
		return "", err
	}
	return storedHash, nil
}
