package services

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/antonchaban/articles-go/internal/dto"
	"github.com/antonchaban/articles-go/internal/entities"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type MockArticleRepository struct {
	mock.Mock
}

func (m *MockArticleRepository) Create(ctx context.Context, article *entities.Article) error {
	args := m.Called(ctx, article)
	return args.Error(0)
}

func (m *MockArticleRepository) GetByID(ctx context.Context, id uint) (*entities.Article, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Article), args.Error(1)
}

func TestCreateArticleWithValidTitle(t *testing.T) {
	mockRepo := new(MockArticleRepository)
	logger := zap.NewNop()
	service := NewArticleService(mockRepo, logger)

	req := dto.CreateArticleRequest{
		Title: "Valid Article Title",
	}

	mockRepo.On("Create", mock.Anything, mock.MatchedBy(func(a *entities.Article) bool {
		return a.Title == "Valid Article Title"
	})).Run(func(args mock.Arguments) {
		article := args.Get(1).(*entities.Article)
		article.ID = 1
	}).Return(nil)

	resp, err := service.Create(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, uint(1), resp.ID)
	assert.NotZero(t, resp.CreatedAt)
	mockRepo.AssertExpectations(t)
}

func TestCreateArticleWithEmptyTitle(t *testing.T) {
	mockRepo := new(MockArticleRepository)
	logger := zap.NewNop()
	service := NewArticleService(mockRepo, logger)

	req := dto.CreateArticleRequest{
		Title: "",
	}

	resp, err := service.Create(context.Background(), req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Equal(t, "title cannot be empty", err.Error())
	mockRepo.AssertNotCalled(t, "Create")
}

func TestCreateArticleWithRepositoryError(t *testing.T) {
	mockRepo := new(MockArticleRepository)
	logger := zap.NewNop()
	service := NewArticleService(mockRepo, logger)

	req := dto.CreateArticleRequest{
		Title: "Valid Article Title",
	}

	expectedError := errors.New("database error")
	mockRepo.On("Create", mock.Anything, mock.Anything).Return(expectedError)

	resp, err := service.Create(context.Background(), req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Equal(t, expectedError, err)
	mockRepo.AssertExpectations(t)
}

func TestGetByIDReturnsArticle(t *testing.T) {
	mockRepo := new(MockArticleRepository)
	logger := zap.NewNop()
	service := NewArticleService(mockRepo, logger)

	expectedArticle := &entities.Article{
		ID:        1,
		Title:     "Test Article",
		CreatedAt: time.Now().UTC(),
	}

	mockRepo.On("GetByID", mock.Anything, uint(1)).Return(expectedArticle, nil)

	resp, err := service.GetByID(context.Background(), 1)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, expectedArticle.ID, resp.ID)
	assert.Equal(t, expectedArticle.Title, resp.Title)
	assert.Equal(t, expectedArticle.CreatedAt, resp.CreatedAt)
	mockRepo.AssertExpectations(t)
}

func TestGetByIDWithNonExistentArticle(t *testing.T) {
	mockRepo := new(MockArticleRepository)
	logger := zap.NewNop()
	service := NewArticleService(mockRepo, logger)

	mockRepo.On("GetByID", mock.Anything, uint(999)).Return(nil, gorm.ErrRecordNotFound)

	resp, err := service.GetByID(context.Background(), 999)

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
	mockRepo.AssertExpectations(t)
}

func TestGetByIDWithRepositoryError(t *testing.T) {
	mockRepo := new(MockArticleRepository)
	logger := zap.NewNop()
	service := NewArticleService(mockRepo, logger)

	expectedError := errors.New("database connection error")
	mockRepo.On("GetByID", mock.Anything, uint(1)).Return(nil, expectedError)

	resp, err := service.GetByID(context.Background(), 1)

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Equal(t, expectedError, err)
	mockRepo.AssertExpectations(t)
}
