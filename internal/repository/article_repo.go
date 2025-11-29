package repository

import (
	"context"
	"errors"

	"github.com/antonchaban/articles-go/internal/entities"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

// PostgresRepo implements the ArticleRepository interface using PostgreSQL as the data store.
type PostgresRepo struct {
	db  *gorm.DB
	log *zap.Logger
}

func NewPostgresRepo(db *gorm.DB, logger *zap.Logger) *PostgresRepo {
	return &PostgresRepo{
		db:  db,
		log: logger.With(zap.String("layer", "repository")),
	}
}

// Create inserts a new article into the database.
func (r *PostgresRepo) Create(ctx context.Context, a *entities.Article) error {
	if err := r.db.WithContext(ctx).Create(a).Error; err != nil {
		r.log.Error("failed to create article", zap.Error(err))
		return err
	}
	return nil
}

// GetByID retrieves an article by its ID from the database.
func (r *PostgresRepo) GetByID(ctx context.Context, id uint) (*entities.Article, error) {
	var a entities.Article
	if err := r.db.WithContext(ctx).First(&a, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			r.log.Warn("article not found", zap.Int("id", int(id)))
			return nil, err
		}
		r.log.Error("database query failed", zap.Int("id", int(id)), zap.Error(err))
		return nil, err
	}
	return &a, nil
}
