package entities

import (
	"time"
)

// Article represents simple article entity.
type Article struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Title     string    `gorm:"not null" json:"title"`
	CreatedAt time.Time `json:"created_at"`
}
