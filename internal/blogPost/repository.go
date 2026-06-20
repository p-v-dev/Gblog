package blogPost

import (
	"context"

	"gorm.io/gorm"
)

// BlogPostRepository define o contrato (interface) das operações de banco
type BlogPostRepository interface {
	Create(ctx context.Context, post *BlogPost) error // <-- Sem o "entity."
	GetBySlug(ctx context.Context, slug string) (*BlogPost, error)
	FetchAll(ctx context.Context, limit, offset int) ([]BlogPost, error)
	Update(ctx context.Context, post *BlogPost) error
	FindByID(ctx context.Context, id string) (*BlogPost, error)
}

// blogPostRepositoryORM é a implementação concreta usando o GORM
type blogPostRepositoryORM struct {
	db *gorm.DB
}

// NewBlogPostRepository é o construtor que devolve a interface
func NewBlogPostRepository(db *gorm.DB) BlogPostRepository {
	return &blogPostRepositoryORM{db: db}
}

// Create insere um novo post no banco de dados
func (r *blogPostRepositoryORM) Create(ctx context.Context, post *BlogPost) error {
	if err := r.db.WithContext(ctx).Omit("Tags").Create(post).Error; err != nil {
		return err
	}
	for _, tag := range post.Tags {
		if err := r.db.WithContext(ctx).Exec(
			"INSERT INTO post_tags (blog_post_id, tag_id) VALUES (?, ?)",
			post.ID, tag.ID,
		).Error; err != nil {
			return err
		}
	}
	return nil
}

// GetBySlug busca um post específico pela URL amigável (slug)
func (r *blogPostRepositoryORM) GetBySlug(ctx context.Context, slug string) (*BlogPost, error) {
	var post BlogPost
	// Adicionamos a condição is_active = true
	err := r.db.WithContext(ctx).Where("slug = ? AND is_active = ?", slug, true).First(&post).Error
	if err != nil {
		return nil, err
	}
	return &post, nil
}
func (r *blogPostRepositoryORM) FindByID(ctx context.Context, id string) (*BlogPost, error) {
	var post BlogPost
	// Adicionamos a condição is_active = true
	err := r.db.WithContext(ctx).Where("id = ? AND is_active = ?", id, true).First(&post).Error
	if err != nil {
		return nil, err
	}
	return &post, nil
}

// FetchAll lista os posts com paginação (limit e offset)
func (r *blogPostRepositoryORM) FetchAll(ctx context.Context, limit, offset int) ([]BlogPost, error) {
	var posts []BlogPost
	// Adicionamos a condição is_active = true
	err := r.db.WithContext(ctx).
		Where("is_active = ?", true).
		Limit(limit).
		Offset(offset).
		Order("created_at DESC").
		Find(&posts).Error

	if err != nil {
		return nil, err
	}
	return posts, nil
}

// Update salva as alterações de um post existente
func (r *blogPostRepositoryORM) Update(ctx context.Context, post *BlogPost) error {
	return r.db.WithContext(ctx).Save(post).Error
}
