// Package v1 contains the HTTP handlers for version 1 of the API.
package v1

import (
	"context"
	"net/http"
	"strconv"

	"github.com/antonchaban/articles-go/internal/dto"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// ArticleService defines the business logic operations for articles.
type ArticleService interface {
	// Create creates a new article and returns the created article details.
	Create(ctx context.Context, req dto.CreateArticleRequest) (*dto.CreateArticleResponse, error)
	// GetByID retrieves an article by its unique identifier.
	GetByID(ctx context.Context, id uint) (*dto.ArticleResponse, error)
}

// ArticleHandler handles HTTP requests related to articles.
type ArticleHandler struct {
	service ArticleService
	log     *zap.Logger
}

// NewArticleHandler creates a new ArticleHandler with the given service and logger.
func NewArticleHandler(s ArticleService, logger *zap.Logger) *ArticleHandler {
	return &ArticleHandler{
		service: s,
		log:     logger.With(zap.String("layer", "handler")),
	}
}

// Create handles POST requests to create a new article.
// It expects a JSON body conforming to dto.CreateArticleRequest.
// Returns 201 Created on success, 400 Bad Request for invalid input,
// or 500 Internal Server Error if article creation fails.
func (h *ArticleHandler) Create(c *gin.Context) {
	var req dto.CreateArticleRequest

	// Bind and validate JSON request body
	if err := c.ShouldBindJSON(&req); err != nil {
		h.log.Warn("invalid json request", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// call service layer to create the article
	resp, err := h.service.Create(c.Request.Context(), req)
	if err != nil {
		h.log.Error("failed to create article", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, resp)
}

// Get handles GET requests to retrieve an article by ID.
// The article ID should be provided as a URL parameter.
// Returns 200 OK with article data on success, 400 Bad Request for invalid ID format,
// or 404 Not Found if the article doesn't exist.
func (h *ArticleHandler) Get(c *gin.Context) {
	// Extract article ID from URL parameter
	idStr := c.Param("id")

	// Convert ID to integer and validate it's a positive number
	idInt, err := strconv.Atoi(idStr)
	if err != nil || idInt < 0 {
		h.log.Warn("invalid article id format", zap.String("id_param", idStr))
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID format; must be a positive integer"})
		return
	}

	idUint := uint(idInt)

	// Fetch article from service layer
	resp, err := h.service.GetByID(c.Request.Context(), idUint)
	if err != nil {
		h.log.Error("failed to fetch article", zap.Uint("id", idUint), zap.Error(err))
		c.JSON(http.StatusNotFound, gin.H{"error": "article not found"})
		return
	}

	c.JSON(http.StatusOK, resp)
}
