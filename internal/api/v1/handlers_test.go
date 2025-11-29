package v1

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/antonchaban/articles-go/internal/dto"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type MockArticleService struct {
	mock.Mock
}

func (m *MockArticleService) Create(ctx context.Context, req dto.CreateArticleRequest) (*dto.CreateArticleResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.CreateArticleResponse), args.Error(1)
}

func (m *MockArticleService) GetByID(ctx context.Context, id uint) (*dto.ArticleResponse, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.ArticleResponse), args.Error(1)
}

func setupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	return gin.New()
}

func TestCreateHandlerWithValidRequest(t *testing.T) {
	mockService := new(MockArticleService)
	logger := zap.NewNop()
	handler := NewArticleHandler(mockService, logger)

	router := setupTestRouter()
	router.POST("/articles", handler.Create)

	reqBody := dto.CreateArticleRequest{
		Title: "Test Article",
	}

	expectedResp := &dto.CreateArticleResponse{
		ID:        1,
		CreatedAt: time.Now().UTC(),
	}

	mockService.On("Create", mock.Anything, reqBody).Return(expectedResp, nil)

	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/articles", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var response dto.CreateArticleResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, expectedResp.ID, response.ID)
	mockService.AssertExpectations(t)
}

func TestCreateHandlerWithInvalidJSON(t *testing.T) {
	mockService := new(MockArticleService)
	logger := zap.NewNop()
	handler := NewArticleHandler(mockService, logger)

	router := setupTestRouter()
	router.POST("/articles", handler.Create)

	invalidJSON := []byte(`{"title": }`)
	req := httptest.NewRequest(http.MethodPost, "/articles", bytes.NewBuffer(invalidJSON))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	mockService.AssertNotCalled(t, "Create")
}

func TestCreateHandlerWithMissingRequiredField(t *testing.T) {
	mockService := new(MockArticleService)
	logger := zap.NewNop()
	handler := NewArticleHandler(mockService, logger)

	router := setupTestRouter()
	router.POST("/articles", handler.Create)

	reqBody := map[string]string{}
	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/articles", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	mockService.AssertNotCalled(t, "Create")
}

func TestCreateHandlerWithServiceError(t *testing.T) {
	mockService := new(MockArticleService)
	logger := zap.NewNop()
	handler := NewArticleHandler(mockService, logger)

	router := setupTestRouter()
	router.POST("/articles", handler.Create)

	reqBody := dto.CreateArticleRequest{
		Title: "Test Article",
	}

	expectedError := errors.New("service error")
	mockService.On("Create", mock.Anything, reqBody).Return(nil, expectedError)

	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/articles", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockService.AssertExpectations(t)
}

func TestGetHandlerReturnsArticleSuccessfully(t *testing.T) {
	mockService := new(MockArticleService)
	logger := zap.NewNop()
	handler := NewArticleHandler(mockService, logger)

	router := setupTestRouter()
	router.GET("/articles/:id", handler.Get)

	expectedResp := &dto.ArticleResponse{
		ID:        1,
		Title:     "Test Article",
		CreatedAt: time.Now().UTC(),
	}

	mockService.On("GetByID", mock.Anything, uint(1)).Return(expectedResp, nil)

	req := httptest.NewRequest(http.MethodGet, "/articles/1", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response dto.ArticleResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, expectedResp.ID, response.ID)
	assert.Equal(t, expectedResp.Title, response.Title)
	mockService.AssertExpectations(t)
}

func TestGetHandlerWithInvalidIDFormat(t *testing.T) {
	mockService := new(MockArticleService)
	logger := zap.NewNop()
	handler := NewArticleHandler(mockService, logger)

	router := setupTestRouter()
	router.GET("/articles/:id", handler.Get)

	req := httptest.NewRequest(http.MethodGet, "/articles/invalid", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	mockService.AssertNotCalled(t, "GetByID")
}

func TestGetHandlerWithNegativeID(t *testing.T) {
	mockService := new(MockArticleService)
	logger := zap.NewNop()
	handler := NewArticleHandler(mockService, logger)

	router := setupTestRouter()
	router.GET("/articles/:id", handler.Get)

	req := httptest.NewRequest(http.MethodGet, "/articles/-1", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	mockService.AssertNotCalled(t, "GetByID")
}

func TestGetHandlerWithNonExistentArticle(t *testing.T) {
	mockService := new(MockArticleService)
	logger := zap.NewNop()
	handler := NewArticleHandler(mockService, logger)

	router := setupTestRouter()
	router.GET("/articles/:id", handler.Get)

	mockService.On("GetByID", mock.Anything, uint(999)).Return(nil, gorm.ErrRecordNotFound)

	req := httptest.NewRequest(http.MethodGet, "/articles/999", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	mockService.AssertExpectations(t)
}

func TestGetHandlerWithServiceError(t *testing.T) {
	mockService := new(MockArticleService)
	logger := zap.NewNop()
	handler := NewArticleHandler(mockService, logger)

	router := setupTestRouter()
	router.GET("/articles/:id", handler.Get)

	expectedError := errors.New("service error")
	mockService.On("GetByID", mock.Anything, uint(1)).Return(nil, expectedError)

	req := httptest.NewRequest(http.MethodGet, "/articles/1", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	mockService.AssertExpectations(t)
}
