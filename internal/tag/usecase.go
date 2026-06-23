package tag

import (
	"context"
	"errors"
)

type CreateTagUseCase interface {
	Execute(ctx context.Context, input CreateTagInput) (*TagOutput, error)
}

type createTagUseCase struct {
	repo Repository
}

func NewCreateTagUseCase(repo Repository) CreateTagUseCase {
	return &createTagUseCase{repo: repo}
}

func (uc *createTagUseCase) Execute(ctx context.Context, input CreateTagInput) (*TagOutput, error) {
	if input.Name == "" {
		return nil, errors.New("o nome da tag é obrigatório")
	}
	if len(input.Name) > 20 {
		return nil, errors.New("o nome da tag deve ter no máximo 20 caracteres")
	}

	existing, err := uc.repo.FindByName(ctx, input.Name)
	if err == nil {
		return &TagOutput{ID: existing.ID, Name: existing.Name}, nil
	}

	tag := &Tag{Name: input.Name}
	if err := uc.repo.Create(ctx, tag); err != nil {
		return nil, err
	}
	return &TagOutput{ID: tag.ID, Name: tag.Name}, nil
}

type ListTagsUseCase interface {
	Execute(ctx context.Context) ([]TagOutput, error)
}

type listTagsUseCase struct {
	repo Repository
}

func NewListTagsUseCase(repo Repository) ListTagsUseCase {
	return &listTagsUseCase{repo: repo}
}

func (uc *listTagsUseCase) Execute(ctx context.Context) ([]TagOutput, error) {
	tags, err := uc.repo.FetchAll(ctx)
	if err != nil {
		return nil, err
	}
	out := make([]TagOutput, len(tags))
	for i := range tags {
		out[i] = TagOutput{ID: tags[i].ID, Name: tags[i].Name}
	}
	return out, nil
}
