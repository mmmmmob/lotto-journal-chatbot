package database

import (
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase(dsn string) {
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
	sqlDB.SetMaxIdleConns(0)
	sqlDB.SetMaxOpenConns(5)
	sqlDB.SetConnMaxLifetime(3 * time.Minute)

	log.Println("Database connection successful")
}
