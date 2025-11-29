package database

import (
	"github.com/antonchaban/articles-go/internal/entities"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// NewPostgresConnection initializes a new GORM DB connection to PostgreSQL
func NewPostgresConnection(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(&entities.Article{})
	if err != nil {
		return nil, err
	}

	return db, nil
}
