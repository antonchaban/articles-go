package database

import (
	"fmt"
	"log"
	"time"

	"github.com/antonchaban/articles-go/internal/entities"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// NewPostgresConnection initializes a new GORM DB connection to PostgreSQL
func NewPostgresConnection(dsn string) (*gorm.DB, error) {
	var db *gorm.DB
	var err error

	// Retry logic: Try 5 times with a 2-second delay
	// This handles the race condition where the App starts faster than Postgres
	maxRetries := 5
	for i := 1; i <= maxRetries; i++ {
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err == nil {
			break
		}

		log.Printf("Failed to connect to DB (attempt %d/%d). Retrying in 2s...", i, maxRetries)
		time.Sleep(2 * time.Second)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to connect to database after %d attempts: %w", maxRetries, err)
	}

	err = db.AutoMigrate(&entities.Article{})
	if err != nil {
		return nil, err
	}

	return db, nil
}
