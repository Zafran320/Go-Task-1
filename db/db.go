package db

import (
	"database/sql"
	"log"

	"backend-auth/config"

	"github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func InitDB() {
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
	DB, err = sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal("sql.Open error:", err)
	}

	if err = DB.Ping(); err != nil {
		log.Fatal("Database ping failed:", err)
	}

	log.Println("Connected to MySQL successfully")
}
