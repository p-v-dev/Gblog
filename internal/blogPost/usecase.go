package blogPost

import (
	"Gblog/pkg"
	"context"
	"errors"
)

// BlogPostUseCase agora expõe apenas DTOs para o mundo externo
type BlogPostUseCase interface {
	Create(ctx context.Context, input BlogPostInputDTO) error
	GetBySlug(ctx context.Context, slug string) (*BlogPostOutputDTO, error)
	FetchAll(ctx context.Context, limit, offset int) ([]BlogPostOutputDTO, error)
	Update(ctx context.Context, id uint, input BlogPostInputDTO) error
	Delete(ctx context.Context, id uint) error
}

type blogPostUseCase struct {
	repo BlogPostRepository
}

func NewBlogPostUseCase(repo BlogPostRepository) BlogPostUseCase {
	return &blogPostUseCase{repo: repo}
}

// Create mapeia o Input DTO para a Entity antes de salvar
func (uc *blogPostUseCase) Create(ctx context.Context, input BlogPostInputDTO) error {
	statusConvertido := pkg.BlogStatus(input.Status)

	// Valida se o status enviado é aceito pelo seu sistema
	if !statusConvertido.IsValid() {
		return errors.New("status inválido: escolha entre draft, published ou archived")
	}

	postEntity := &BlogPost{
		Title:   input.Title,
		Slug:    input.Slug,
		Content: input.Content,
		Status:  statusConvertido,
	}

	return uc.repo.Create(ctx, postEntity)
}

// GetBySlug busca a Entity e mapeia para o Output DTO
func (uc *blogPostUseCase) GetBySlug(ctx context.Context, slug string) (*BlogPostOutputDTO, error) {
	post, err := uc.repo.GetBySlug(ctx, slug)
	if err != nil {
		return nil, err
	}

	return &BlogPostOutputDTO{
		ID:        post.ID,
		Title:     post.Title,
		Slug:      post.Slug,
		Content:   post.Content,
		Status:    string(post.Status),
		IsActive:  post.IsActive,
		CreatedAt: post.CreatedAt,
		UpdatedAt: post.UpdatedAt,
	}, nil
}

// FetchAll busca a lista de Entities e transforma em uma lista de Output DTOs
func (uc *blogPostUseCase) FetchAll(ctx context.Context, limit, offset int) ([]BlogPostOutputDTO, error) {
	posts, err := uc.repo.FetchAll(ctx, limit, offset)
	if err != nil {
		return nil, err
	}

	dtos := make([]BlogPostOutputDTO, len(posts))
	for i, post := range posts {
		dtos[i] = BlogPostOutputDTO{
			ID:        post.ID,
			Title:     post.Title,
			Slug:      post.Slug,
			Content:   post.Content,
			Status:    string(post.Status),
			IsActive:  post.IsActive,
			CreatedAt: post.CreatedAt,
			UpdatedAt: post.UpdatedAt,
		}
	}

	return dtos, nil
}

// Update recebe o ID e os novos dados, busca o registro atual e atualiza a Entity
func (uc *blogPostUseCase) Update(ctx context.Context, id uint, input BlogPostInputDTO) error {
	// Uma boa prática é verificar se o post existe antes de atualizar via ID
	// Aqui assumimos uma atualização direta baseada no ID recebido
	statusDoPacote := pkg.BlogStatus(input.Status)
	postEntity := &BlogPost{
		Title:   input.Title,
		Slug:    input.Slug,
		Content: input.Content,
		Status:  statusDoPacote,
	}
	postEntity.ID = id // Injeta o ID vindo da rota/parâmetro na Entity do GORM

	return uc.repo.Update(ctx, postEntity)
}

func (uc *blogPostUseCase) Delete(ctx context.Context, id uint) error {
	return uc.repo.Delete(ctx, id)
}
