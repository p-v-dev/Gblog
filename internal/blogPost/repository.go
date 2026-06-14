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
	Delete(ctx context.Context, id uint) error
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
	return r.db.WithContext(ctx).Create(post).Error
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

// Delete remove um post
func (r *blogPostRepositoryORM) Delete(ctx context.Context, id uint) error {
	// 1. Atualiza o campo IsActive para false antes de deletar logicamente
	err := r.db.WithContext(ctx).Model(&BlogPost{}).Where("id = ?", id).Update("is_active", false).Error
	if err != nil {
		return err
	}

	// 2. Executa o Soft Delete do GORM (preenche a coluna deleted_at)
	return r.db.WithContext(ctx).Delete(&BlogPost{}, id).Error
}
