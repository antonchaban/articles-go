package api

import (
	"net/http"

	v1 "github.com/antonchaban/articles-go/internal/api/v1"
	"github.com/antonchaban/articles-go/internal/config"

	"github.com/gin-gonic/gin"
)

// NewServer initializes and configures the Gin HTTP engine with all routes and middlewares.
//
// The function performs the following setup:
//   - Configures Gin mode (Debug/Release) based on environment
//   - Initializes default middleware (Logger and Recovery)
//   - Registers a health check endpoint at /health
//   - Sets up API versioning with v1 routes at /api/v1
//
// Parameters:
//   - cfg: Application configuration containing environment settings
//   - articleHandler: Handler for article-related API endpoints (injected via DI)
//
// Returns:
//   - *gin.Engine: Configured Gin engine ready to serve HTTP requests
func NewServer(cfg *config.Config, articleHandler *v1.ArticleHandler) *gin.Engine {
	// Set Gin mode based on environment configuration
	// Production mode disables debug logging for better performance
	if cfg.AppEnv == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Initialize Gin engine with default middleware
	// Default includes Logger (request logging) and Recovery (panic recovery)
	r := gin.Default()

	// Apply additional recovery middleware for redundancy
	// Ensures graceful handling of panics in request handlers
	r.Use(gin.Recovery())

	// Register health check endpoint
	// Used for liveness probes
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "alive"})
	})

	// Create API v1 route group
	// All v1 endpoints will be prefixed with /api/v1
	apiV1 := r.Group("/api/v1")
	{
		// Register all article-related routes under /api/v1
		v1.RegisterRoutes(apiV1, articleHandler)
	}

	return r
}
