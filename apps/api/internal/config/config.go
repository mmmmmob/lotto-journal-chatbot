package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DB_DSN string
	PORT   string
}

func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, using system environment variables")
	}
	return &Config{
		DB_DSN: os.Getenv("DB_DSN"),
		PORT:   os.Getenv("PORT"),
	}
}
