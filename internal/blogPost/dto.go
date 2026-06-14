// DTOs específicos por intenção
package blogPost

type CreatePostInputDTO struct {
	Title   string `json:"title"`
	Slug    string `json:"slug"`
	Content string `json:"content"`
}

type UpdatePostInputDTO struct {
	Title   string `json:"title"`
	Slug    string `json:"slug"`
	Content string `json:"content"`
}
