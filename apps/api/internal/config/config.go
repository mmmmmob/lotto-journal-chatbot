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
	cfg := &Config{
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
	}

	// Validate and clamp DB connection pool configurations to prevent connection leaks
	maxIdle := getEnvInt("DB_MAX_IDLE_CONNS", 0)
	if maxIdle < 0 {
		log.Printf("[config] Warning: DB_MAX_IDLE_CONNS (%d) cannot be negative. Clamping to 0.", maxIdle)
		maxIdle = 0
	}
	cfg.DBMaxIdleConns = maxIdle

	maxOpen := getEnvInt("DB_MAX_OPEN_CONNS", 5)
	if maxOpen < 1 {
		log.Printf("[config] Warning: DB_MAX_OPEN_CONNS (%d) must be at least 1 to prevent unlimited connection leaks. Clamping to default: 5.", maxOpen)
		maxOpen = 5
	}
	cfg.DBMaxOpenConns = maxOpen

	maxLifetime := getEnvDuration("DB_CONN_MAX_LIFETIME", 3*time.Minute)
	if maxLifetime < 30*time.Second {
		log.Printf("[config] Warning: DB_CONN_MAX_LIFETIME (%v) is below safe minimum (30s). Clamping to default: 3m.", maxLifetime)
		maxLifetime = 3 * time.Minute
	}
	cfg.DBConnMaxLifetime = maxLifetime

	maxIdleTime := getEnvDuration("DB_CONN_MAX_IDLE_TIME", 1*time.Minute)
	if maxIdleTime < 10*time.Second {
		log.Printf("[config] Warning: DB_CONN_MAX_IDLE_TIME (%v) is below safe minimum (10s). Clamping to default: 1m.", maxIdleTime)
		maxIdleTime = 1 * time.Minute
	}
	cfg.DBConnMaxIdleTime = maxIdleTime

	return cfg
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
