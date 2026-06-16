package blogPost

import (
	"Gblog/pkg/blogstatus"
	"context"
	"errors"
)

// ---------------------------------------------------------
// DTOs (Data Transfer Objects)
// ---------------------------------------------------------

// ---------------------------------------------------------
// CASO DE USO 1: Criar um Blog Post
// ---------------------------------------------------------

type CreatePostUseCase interface {
	Execute(ctx context.Context, input CreatePostInputDTO) error
}

type createPostUseCase struct {
	repo BlogPostRepository
}

func NewCreatePostUseCase(repo BlogPostRepository) CreatePostUseCase {
	return &createPostUseCase{repo: repo}
}

func (uc *createPostUseCase) Execute(ctx context.Context, input CreatePostInputDTO) error {
	statusDefault := blogstatus.Draft

	if input.Title == "" {
		return errors.New("o título do post é obrigatório")
	}
	if input.UserID == "" {
		return errors.New("o ID do usuário é obrigatório")
	}

	postEntity := &BlogPost{
		Title:   input.Title,
		Slug:    input.Slug,
		Content: input.Content,
		Status:  statusDefault,
		UserID:  input.UserID,
	}

	return uc.repo.Create(ctx, postEntity)
}

// ---------------------------------------------------------
// CASO DE USO 2: Publicar um Blog Post
// ---------------------------------------------------------

type PublishPostUseCase interface {
	Execute(ctx context.Context, id string) error
}

type publishPostUseCase struct {
	repo BlogPostRepository
}

func NewPublishPostUseCase(repo BlogPostRepository) PublishPostUseCase {
	return &publishPostUseCase{repo: repo}
}

func (uc *publishPostUseCase) Execute(ctx context.Context, id string) error {
	post, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return errors.New("post não encontrado")
	}

	if post.Status == "published" {
		return errors.New("este post já está publicado")
	}

	if len(post.Content) < 10 {
		return errors.New("conteúdo muito curto para ser publicado")
	}

	post.Status = blogstatus.Published

	return uc.repo.Update(ctx, post)
}

// ---------------------------------------------------------
// CASO DE USO 3: Editar Blog Post
// ---------------------------------------------------------

type UpdatePostUseCase interface {
	Execute(ctx context.Context, id string, input UpdatePostInputDTO) error
}

type updatePostUseCase struct {
	repo BlogPostRepository
}

func NewUpdatePostUseCase(repo BlogPostRepository) UpdatePostUseCase {
	return &updatePostUseCase{repo: repo}
}

func (uc *updatePostUseCase) Execute(ctx context.Context, id string, input UpdatePostInputDTO) error {
	post, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return errors.New("post não encontrado")
	}

	if !post.IsActive {
		return errors.New("não é possível editar um post inativo/deletado")
	}

	post.Title = input.Title
	post.Content = input.Content
	post.Slug = input.Slug

	return uc.repo.Update(ctx, post)
}

// ---------------------------------------------------------
// CASO DE USO 4: Deletar (Soft Delete)
// ---------------------------------------------------------

type DeletePostUseCase interface {
	Execute(ctx context.Context, id string) error
}

type deletePostUseCase struct {
	repo BlogPostRepository
}

func NewDeletePostUseCase(repo BlogPostRepository) DeletePostUseCase {
	return &deletePostUseCase{repo: repo}
}

func (uc *deletePostUseCase) Execute(ctx context.Context, id string) error {
	post, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return errors.New("post não encontrado")
	}

	post.IsActive = false

	return uc.repo.Update(ctx, post)
}
