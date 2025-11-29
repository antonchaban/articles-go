package dto

import "time"

type CreateArticleRequest struct {
	Title string `json:"title" binding:"required"`
}

type CreateArticleResponse struct {
	ID        uint      `json:"id"`
	CreatedAt time.Time `json:"created_at"`
}

type ArticleResponse struct {
	ID        uint      `json:"id"`
	Title     string    `json:"title"`
	CreatedAt time.Time `json:"created_at"`
}
