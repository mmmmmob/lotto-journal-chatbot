package database

import (
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase(dsn string, maxIdle, maxOpen int, maxLifetime, maxIdleTime time.Duration) {
	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to database: " + err.Error())
	}

	sqlDB, err := DB.DB()
	if err != nil {
		panic("Failed to get sql.DB from gorm: " + err.Error())
	}

	// Validate and clamp parameters defensively inside ConnectDatabase to ensure safe limits
	if maxIdle < 0 {
		maxIdle = 0
	}
	if maxOpen < 1 {
		maxOpen = 1
	}
	if maxLifetime < 30*time.Second {
		maxLifetime = 3 * time.Minute // Safe default fallback to prevent infinite lifetime connection leaks
	}
	if maxIdleTime < 10*time.Second {
		maxIdleTime = 1 * time.Minute // Safe default fallback to prevent infinite idle connection leaks
	}

	// Optimize connection settings for serverless database (Neon) and scale-to-zero autoscaling (Fly.io)
	sqlDB.SetMaxIdleConns(maxIdle)
	sqlDB.SetMaxOpenConns(maxOpen)
	sqlDB.SetConnMaxLifetime(maxLifetime)
	sqlDB.SetConnMaxIdleTime(maxIdleTime)

	log.Printf("Database connection successful (pool config: maxOpen=%d, maxIdle=%d, maxLifetime=%v, maxIdleTime=%v)",
		maxOpen, maxIdle, maxLifetime, maxIdleTime)
}
