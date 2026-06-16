// DTOs específicos por intenção
package blogPost

type CreatePostInputDTO struct {
	Title   string `json:"title"`
	Slug    string `json:"slug"`
	Content string `json:"content"`
	UserID  string `json:"user_id"`
}

type UpdatePostInputDTO struct {
	Title   string `json:"title"`
	Slug    string `json:"slug"`
	Content string `json:"content"`
}
