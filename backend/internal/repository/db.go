package repository

import (
	"fmt"
	"log"
	"time"

	"github.com/TryHanger/digital_signage/backend/internal/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDB(cfg *config.Config) *gorm.DB {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBPort,
	)
	// Retry loop: DB container may need a few seconds to become ready.
	const maxAttempts = 20
	const delay = 2 * time.Second
	var db *gorm.DB
	var err error
	for attempt := 1; attempt <= maxAttempts; attempt++ {
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err == nil {
			return db
		}
		log.Printf("DB connect attempt %d/%d failed: %v. Retrying in %s...", attempt, maxAttempts, err, delay)
		time.Sleep(delay)
	}
	// final fatal if we couldn't connect after retries
	log.Fatal("Error connect to DB after retries:", err)
	return nil
}
