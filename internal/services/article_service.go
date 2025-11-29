package services

import (
	"context"
	"errors"
	"time"

	"github.com/antonchaban/articles-go/internal/dto"
	"github.com/antonchaban/articles-go/internal/entities"

	"go.uber.org/zap"
)

// ArticleRepository defines the methods that any
// data storage provider must implement to manage Articles.
type ArticleRepository interface {
	Create(ctx context.Context, article *entities.Article) error
	GetByID(ctx context.Context, id uint) (*entities.Article, error)
}

type ArticleService struct {
	repo ArticleRepository
	log  *zap.Logger
}

func NewArticleService(repo ArticleRepository, log *zap.Logger) *ArticleService {
	return &ArticleService{
		repo: repo,
		log:  log.With(zap.String("layer", "service")),
	}
}

// Create creates a new Article and returns its ID and creation timestamp
func (s *ArticleService) Create(ctx context.Context, req dto.CreateArticleRequest) (*dto.CreateArticleResponse, error) {
	if req.Title == "" {
		s.log.Warn("creation attempt with empty title")
		return nil, errors.New("title cannot be empty")
	}

	// prepare entity
	article := &entities.Article{
		Title:     req.Title,
		CreatedAt: time.Now().UTC(),
	}

	s.log.Info("creating new article", zap.String("title", req.Title))

	// repo cvall
	if err := s.repo.Create(ctx, article); err != nil {
		return nil, err
	}

	s.log.Info("article created successfully", zap.Uint("id", article.ID))

	// return response DTO
	return &dto.CreateArticleResponse{
		ID:        article.ID,
		CreatedAt: article.CreatedAt,
	}, nil
}

// GetByID retrieves an Article by its ID
func (s *ArticleService) GetByID(ctx context.Context, id uint) (*dto.ArticleResponse, error) {
	article, err := s.repo.GetByID(ctx, id)
	if err != nil {
		s.log.Warn("failed to retrieve article", zap.Uint("id", id), zap.Error(err))
		return nil, err
	}

	// return response DTO
	return &dto.ArticleResponse{
		ID:        article.ID,
		Title:     article.Title,
		CreatedAt: article.CreatedAt,
	}, nil
}
