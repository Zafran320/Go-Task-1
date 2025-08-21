package db

import (
	"database/sql"
	"log"

	"backend-auth/config"

	"github.com/go-sql-driver/mysql"
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
