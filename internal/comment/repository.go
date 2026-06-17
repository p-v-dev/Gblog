package comment

import (
	"context"

	"gorm.io/gorm"
)

type Repository interface {
	Create(ctx context.Context, comment *Comment) error
	FindByID(ctx context.Context, id string) (*Comment, error)
	FindByPostID(ctx context.Context, postID string) ([]Comment, error)
	Delete(ctx context.Context, id string) error
}

type commentRepositoryORM struct {
	db *gorm.DB
}

func NewCommentRepository(db *gorm.DB) Repository {
	return &commentRepositoryORM{db: db}
}

func (r *commentRepositoryORM) Create(ctx context.Context, comment *Comment) error {
	return r.db.WithContext(ctx).Create(comment).Error
}

func (r *commentRepositoryORM) FindByID(ctx context.Context, id string) (*Comment, error) {
	var c Comment
	err := r.db.WithContext(ctx).Where("id = ? AND is_active = ?", id, true).First(&c).Error
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *commentRepositoryORM) FindByPostID(ctx context.Context, postID string) ([]Comment, error) {
	var comments []Comment
	err := r.db.WithContext(ctx).
		Where("post_id = ? AND is_active = ?", postID, true).
		Order("created_at ASC").
		Find(&comments).Error
	return comments, err
}

func (r *commentRepositoryORM) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Model(&Comment{}).Where("id = ?", id).Update("is_active", false).Error
}
