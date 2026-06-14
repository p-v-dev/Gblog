package blogPost

import "time"

// BlogPostInputDTO define o que o Use Case PRECISA RECEBER do mundo externo
type BlogPostInputDTO struct {
	Title   string `json:"title"`
	Slug    string `json:"slug"`
	Content string `json:"content"`
	Status  string `json:"status"`
}

// BlogPostOutputDTO define o que o Use Case VAI DEVOLVER para o mundo externo
type BlogPostOutputDTO struct {
	ID        uint      `json:"id"`
	Title     string    `json:"title"`
	Slug      string    `json:"slug"`
	Content   string    `json:"content"`
	Status    string    `json:"status"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
