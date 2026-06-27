package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DB_DSN                 string
	PORT                   string
	LineChannelSecret      string
	LineChannelAccessToken string
	APP_ENV                string
	CronSyncSchedule       string
	CronVerifySchedule     string
}

func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, using system environment variables")
	}
	appEnv := os.Getenv("APP_ENV")
	if appEnv == "" {
		appEnv = "development"
	}
	cronSync := os.Getenv("CRON_SYNC_SCHEDULE")
	if cronSync == "" {
		cronSync = "0 3 * * *"
	}
	cronVerify := os.Getenv("CRON_VERIFY_SCHEDULE")
	if cronVerify == "" {
		cronVerify = "*/5 16-23 * * *"
	}
	return &Config{
		DB_DSN:                 os.Getenv("DB_DSN"),
		PORT:                   os.Getenv("PORT"),
		LineChannelSecret:      os.Getenv("LINE_CHANNEL_SECRET"),
		LineChannelAccessToken: os.Getenv("LINE_CHANNEL_ACCESS_TOKEN"),
		APP_ENV:                appEnv,
		CronSyncSchedule:       cronSync,
		CronVerifySchedule:     cronVerify,
	}
}
