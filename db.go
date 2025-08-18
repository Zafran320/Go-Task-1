package main

import (
	"database/sql"
	"log"

	"github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func InitDB() {
	cfg := mysql.Config{
		User:                 Config.DBUser,
		Passwd:               Config.DBPassword,
		Net:                  "tcp",
		Addr:                 Config.DBHost + ":" + Config.DBPort,
		DBName:               Config.DBName,
		AllowNativePasswords: true,
		ParseTime:            true,
	}

	var err error
	DB, err = sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal("sql.Open error:", err)
	}

	if err = DB.Ping(); err != nil {
		log.Fatal("Database ping failed:", err)
	}

	log.Println("Connected to MySQL successfully")
}
