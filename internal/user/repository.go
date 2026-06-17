package user

import (
	"context"

	"gorm.io/gorm"
)

type Repository interface {
	Create(ctx context.Context, user *User) error
	FindByID(ctx context.Context, id string) (*User, error)
	FindByEmail(ctx context.Context, email string) (*User, error)
	FetchAll(ctx context.Context, limit, offset int) ([]User, error)
	Update(ctx context.Context, user *User) error
	Delete(ctx context.Context, id string) error
}

type userRepositoryORM struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) Repository {
	return &userRepositoryORM{db: db}
}

func (r *userRepositoryORM) Create(ctx context.Context, user *User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

func (r *userRepositoryORM) FindByID(ctx context.Context, id string) (*User, error) {
	var user User
	err := r.db.WithContext(ctx).Where("id = ? AND is_active = ?", id, true).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepositoryORM) FindByEmail(ctx context.Context, email string) (*User, error) {
	var user User
	err := r.db.WithContext(ctx).Where("email = ? AND is_active = ?", email, true).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepositoryORM) FetchAll(ctx context.Context, limit, offset int) ([]User, error) {
	var users []User
	err := r.db.WithContext(ctx).
		Where("is_active = ?", true).
		Limit(limit).
		Offset(offset).
		Order("created_at DESC").
		Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (r *userRepositoryORM) Update(ctx context.Context, user *User) error {
	return r.db.WithContext(ctx).Save(user).Error
}

func (r *userRepositoryORM) Delete(ctx context.Context, id string) error {
	err := r.db.WithContext(ctx).Model(&User{}).Where("id = ?", id).Update("is_active", false).Error
	if err != nil {
		return err
	}
	return r.db.WithContext(ctx).Delete(&User{}, id).Error
}
