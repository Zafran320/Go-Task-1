package service

import (
	"backend-auth/db"
	"backend-auth/middleware"
	"database/sql"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

func SignIn(dbConn *sql.DB, username, password string) (string, error) {
	if username == "" || password == "" {
		return "", errors.New("username and password are required")
	}

	storedHash, err := db.GetPassword(dbConn, username)
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	err = bcrypt.CompareHashAndPassword([]byte(storedHash), []byte(password))
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	token, err := middleware.CreateToken(username)
	if err != nil {
		return "", errors.New("could not create token")
	}

	return token, nil
}
