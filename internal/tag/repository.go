package tag

import (
	"context"

	"gorm.io/gorm"
)

type Repository interface {
	Create(ctx context.Context, tag *Tag) error
	FindByID(ctx context.Context, id string) (*Tag, error)
	FindByName(ctx context.Context, name string) (*Tag, error)
	FetchAll(ctx context.Context) ([]Tag, error)
}

type tagRepositoryORM struct {
	db *gorm.DB
}

func NewTagRepository(db *gorm.DB) Repository {
	return &tagRepositoryORM{db: db}
}

func (r *tagRepositoryORM) Create(ctx context.Context, tag *Tag) error {
	return r.db.WithContext(ctx).Create(tag).Error
}

func (r *tagRepositoryORM) FindByID(ctx context.Context, id string) (*Tag, error) {
	var tag Tag
	err := r.db.WithContext(ctx).First(&tag, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &tag, nil
}

func (r *tagRepositoryORM) FindByName(ctx context.Context, name string) (*Tag, error) {
	var tag Tag
	err := r.db.WithContext(ctx).Where("name = ?", name).First(&tag).Error
	if err != nil {
		return nil, err
	}
	return &tag, nil
}

func (r *tagRepositoryORM) FetchAll(ctx context.Context) ([]Tag, error) {
	var tags []Tag
	err := r.db.WithContext(ctx).Order("name ASC").Find(&tags).Error
	return tags, err
}
