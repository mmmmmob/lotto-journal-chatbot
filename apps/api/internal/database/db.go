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

	// Optimize connection settings for serverless database (Neon) and scale-to-zero autoscaling (Fly.io)
	sqlDB.SetMaxIdleConns(maxIdle)
	sqlDB.SetMaxOpenConns(maxOpen)
	sqlDB.SetConnMaxLifetime(maxLifetime)
	sqlDB.SetConnMaxIdleTime(maxIdleTime)

	log.Printf("Database connection successful (pool config: maxOpen=%d, maxIdle=%d, maxLifetime=%v, maxIdleTime=%v)",
		maxOpen, maxIdle, maxLifetime, maxIdleTime)
}
