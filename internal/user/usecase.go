package user

import (
	"context"
	"errors"
)

type CreateUserUseCase interface {
	Execute(ctx context.Context, input CreateUserInput) (*UserOutput, error)
}

type createUserUseCase struct {
	repo Repository
}

func NewCreateUserUseCase(repo Repository) CreateUserUseCase {
	return &createUserUseCase{repo: repo}
}

func (uc *createUserUseCase) Execute(ctx context.Context, input CreateUserInput) (*UserOutput, error) {
	if input.Name == "" {
		return nil, errors.New("o nome é obrigatório")
	}
	if input.Email == "" {
		return nil, errors.New("o email é obrigatório")
	}

	user := &User{
		Name:  input.Name,
		Email: input.Email,
	}

	if err := uc.repo.Create(ctx, user); err != nil {
		return nil, err
	}

	return &UserOutput{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	}, nil
}
