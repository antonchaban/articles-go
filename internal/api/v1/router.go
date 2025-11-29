package v1

import (
	"github.com/gin-gonic/gin"
)

// RegisterRoutes sets up the routing for the Article feature.
// It accepts a RouterGroup so we can version the API (e.g., /api/v1) easily.
// Routes registered:
//   - POST /articles - Create a new article
//   - GET  /articles/:id - Get an article by ID
func RegisterRoutes(router *gin.RouterGroup, handler *ArticleHandler) {
	// Group routes under /articles
	articles := router.Group("/articles")
	{
		articles.POST("", handler.Create)
		articles.GET("/:id", handler.Get)
	}
}
