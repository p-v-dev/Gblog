package user

import (
	"context"
	"errors"
	"net/mail"

	"golang.org/x/crypto/bcrypt"
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

type GetUserUseCase interface {
	Execute(ctx context.Context, id string) (*UserOutput, error)
}

type getUserUseCase struct {
	repo Repository
}

func NewGetUserUseCase(repo Repository) GetUserUseCase {
	return &getUserUseCase{repo: repo}
}

func (uc *getUserUseCase) Execute(ctx context.Context, id string) (*UserOutput, error) {
	user, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return nil, errors.New("usuário não encontrado")
	}
	return &UserOutput{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	}, nil
}

type LoginUseCase interface {
	Execute(ctx context.Context, email, password string) (*UserOutput, error)
}

type loginUseCase struct {
	repo Repository
}

func NewLoginUseCase(repo Repository) LoginUseCase {
	return &loginUseCase{repo: repo}
}

func (uc *loginUseCase) Execute(ctx context.Context, email, password string) (*UserOutput, error) {
	if email == "" {
		return nil, errors.New("o email é obrigatório")
	}
	if password == "" {
		return nil, errors.New("a senha é obrigatória")
	}

	user, err := uc.repo.FindByEmail(ctx, email)
	if err != nil {
		return nil, errors.New("email ou senha inválidos")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, errors.New("email ou senha inválidos")
	}

	return &UserOutput{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	}, nil
}

type UpdateUserUseCase interface {
	Execute(ctx context.Context, id string, input UpdateUserInput) (*UserOutput, error)
}

type updateUserUseCase struct {
	repo Repository
}

func NewUpdateUserUseCase(repo Repository) UpdateUserUseCase {
	return &updateUserUseCase{repo: repo}
}

func (uc *updateUserUseCase) Execute(ctx context.Context, id string, input UpdateUserInput) (*UserOutput, error) {
	user, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return nil, errors.New("usuário não encontrado")
	}

	if input.Name == "" && input.Email == "" {
		return nil, errors.New("pelo menos um campo (name ou email) deve ser enviado")
	}

	if input.Name != "" {
		user.Name = input.Name
	}
	if input.Email != "" {
		if _, err := mail.ParseAddress(input.Email); err != nil {
			return nil, errors.New("email inválido")
		}
		user.Email = input.Email
	}

	if err := uc.repo.Update(ctx, user); err != nil {
		return nil, err
	}

	return &UserOutput{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	}, nil
}

type DeleteUserUseCase interface {
	Execute(ctx context.Context, id string) error
}

type deleteUserUseCase struct {
	repo Repository
}

func NewDeleteUserUseCase(repo Repository) DeleteUserUseCase {
	return &deleteUserUseCase{repo: repo}
}

func (uc *deleteUserUseCase) Execute(ctx context.Context, id string) error {
	_, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return errors.New("usuário não encontrado")
	}
	return uc.repo.Delete(ctx, id)
}

func (uc *createUserUseCase) Execute(ctx context.Context, input CreateUserInput) (*UserOutput, error) {
	if input.Name == "" {
		return nil, errors.New("o nome é obrigatório")
	}
	if input.Email == "" {
		return nil, errors.New("o email é obrigatório")
	}
	if _, err := mail.ParseAddress(input.Email); err != nil {
		return nil, errors.New("email inválido")
	}
	if input.Password == "" {
		return nil, errors.New("a senha é obrigatória")
	}
	if len(input.Password) < 8 {
		return nil, errors.New("a senha deve ter no mínimo 8 caracteres")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("erro ao processar senha")
	}

	user := &User{
		Name:     input.Name,
		Email:    input.Email,
		Password: string(hash),
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
