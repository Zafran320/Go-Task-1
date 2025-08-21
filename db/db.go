package db

import (
	"backend-auth/models"
	"database/sql"
	"errors"
	"log"
	"time"

	"backend-auth/config"

	"github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
)

type DB struct {
	Db *sql.DB
}

func InitDB() *DB {
	cfg := mysql.Config{
		User:                 config.Config.DBUser,
		Passwd:               config.Config.DBPassword,
		Net:                  "tcp",
		Addr:                 config.Config.DBHost + ":" + config.Config.DBPort,
		DBName:               config.Config.DBName,
		AllowNativePasswords: true,
		ParseTime:            true,
	}

	var err error
	db, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal("sql.Open error:", err)
	}

	if err = db.Ping(); err != nil {
		log.Fatal("Database ping failed:", err)
	}

	log.Println("Connected to MySQL successfully")

	return &DB{Db: db}
}

// Sign up handler funtion
func HandleUserQuery(db *sql.DB, user models.User) error {
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
