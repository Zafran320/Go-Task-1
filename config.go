package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

var Config struct {
	DBUser     string
	DBPassword string
	DBHost     string
	DBPort     string
	DBName     string
}

func InitConfig() {

	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, relying on system environment variables")
	}

	viper.AutomaticEnv()

	Config.DBUser = viper.GetString("DB_USER")
	Config.DBPassword = viper.GetString("DB_PASSWORD")
	Config.DBHost = viper.GetString("DB_HOST")
	Config.DBPort = viper.GetString("DB_PORT")
	Config.DBName = viper.GetString("DB_NAME")

	log.Println("Configuration loaded successfully")
}
