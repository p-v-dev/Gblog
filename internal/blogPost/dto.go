package blogPost

import (
	"Gblog/internal/tag"
	"time"
)

type CreatePostInputDTO struct {
	Title   string   `json:"title"`
	Content string   `json:"content"`
	UserID  string   `json:"user_id"`
	Tags    []string `json:"tags"`
}

type UpdatePostInputDTO struct {
	Title   string   `json:"title"`
	Content string   `json:"content"`
	Tags    []string `json:"tags"`
}

type PostOutput struct {
	ID          string          `json:"id"`
	Title       string          `json:"title"`
	Content     string          `json:"content"`
	ContentHTML string          `json:"content_html"`
	Status      string          `json:"status"`
	Slug        string          `json:"slug"`
	UserID      string          `json:"user_id"`
	Tags        []tag.TagOutput `json:"tags"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
}
