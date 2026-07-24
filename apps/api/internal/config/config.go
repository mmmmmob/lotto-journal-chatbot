package config

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	DB_DSN                 string
	PORT                   string
	LineChannelSecret      string
	LineChannelAccessToken string
	APP_ENV                string
	R2AccountID            string
	R2AccessKeyID          string
	R2SecretAccessKey      string
	R2BucketName           string
	R2PublicURLPrefix      string
	OpenAIAPIKey           string
	CronSecret             string
	// DB connection pool tuning parameters
	DBMaxIdleConns    int
	DBMaxOpenConns    int
	DBConnMaxLifetime time.Duration
	DBConnMaxIdleTime time.Duration
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
	return &Config{
		DB_DSN:                 os.Getenv("DB_DSN"),
		PORT:                   os.Getenv("PORT"),
		LineChannelSecret:      os.Getenv("LINE_CHANNEL_SECRET"),
		LineChannelAccessToken: os.Getenv("LINE_CHANNEL_ACCESS_TOKEN"),
		APP_ENV:                appEnv,
		R2AccountID:            os.Getenv("R2_ACCOUNT_ID"),
		R2AccessKeyID:          os.Getenv("R2_ACCESS_KEY_ID"),
		R2SecretAccessKey:      os.Getenv("R2_SECRET_ACCESS_KEY"),
		R2BucketName:           os.Getenv("R2_BUCKET_NAME"),
		R2PublicURLPrefix:      os.Getenv("R2_PUBLIC_URL_PREFIX"),
		OpenAIAPIKey:           os.Getenv("OPENAI_API_KEY"),
		CronSecret:             os.Getenv("CRON_SECRET"),
		DBMaxIdleConns:         getEnvInt("DB_MAX_IDLE_CONNS", 0),
		DBMaxOpenConns:         getEnvInt("DB_MAX_OPEN_CONNS", 5),
		DBConnMaxLifetime:      getEnvDuration("DB_CONN_MAX_LIFETIME", 3*time.Minute),
		DBConnMaxIdleTime:      getEnvDuration("DB_CONN_MAX_IDLE_TIME", 1*time.Minute),
	}
}

func getEnvInt(key string, defaultVal int) int {
	valStr := os.Getenv(key)
	if valStr == "" {
		return defaultVal
	}
	val, err := strconv.Atoi(valStr)
	if err != nil {
		log.Printf("[config] invalid int value for %s: %s, using default: %d", key, valStr, defaultVal)
		return defaultVal
	}
	return val
}

func getEnvDuration(key string, defaultVal time.Duration) time.Duration {
	valStr := os.Getenv(key)
	if valStr == "" {
		return defaultVal
	}
	val, err := time.ParseDuration(valStr)
	if err != nil {
		log.Printf("[config] invalid duration value for %s: %s, using default: %v", key, valStr, defaultVal)
		return defaultVal
	}
	return val
}
