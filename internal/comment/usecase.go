package comment

import (
	"context"
	"errors"
	"time"
)

// ponytail: minimal interface to decouple comment from blogPost package
type PostExistenceChecker interface {
	PostExists(ctx context.Context, id string) bool
}

type CreateCommentUseCase interface {
	Execute(ctx context.Context, input CreateCommentInput) (*CommentOutput, error)
}

type createCommentUseCase struct {
	repo       Repository
	postCheck  PostExistenceChecker
}

func NewCreateCommentUseCase(repo Repository, postCheck PostExistenceChecker) CreateCommentUseCase {
	return &createCommentUseCase{repo: repo, postCheck: postCheck}
}

func (uc *createCommentUseCase) Execute(ctx context.Context, input CreateCommentInput) (*CommentOutput, error) {
	if input.Content == "" {
		return nil, errors.New("o conteúdo do comentário é obrigatório")
	}
	if len(input.Content) < 2 {
		return nil, errors.New("o comentário deve ter pelo menos 2 caracteres")
	}
	if input.PostID == "" {
		return nil, errors.New("post não informado")
	}
	if input.UserID == "" {
		return nil, errors.New("usuário não informado")
	}

	if !uc.postCheck.PostExists(ctx, input.PostID) {
		return nil, errors.New("post não encontrado")
	}

	comment := &Comment{
		Content: input.Content,
		PostID:  input.PostID,
		UserID:  input.UserID,
	}

	if err := uc.repo.Create(ctx, comment); err != nil {
		return nil, err
	}

	return &CommentOutput{
		ID:        comment.ID,
		Content:   comment.Content,
		PostID:    comment.PostID,
		UserID:    comment.UserID,
		CreatedAt: time.Now(),
	}, nil
}

type ListCommentsUseCase interface {
	Execute(ctx context.Context, postID string) ([]CommentOutput, error)
}

type listCommentsUseCase struct {
	repo Repository
}

func NewListCommentsUseCase(repo Repository) ListCommentsUseCase {
	return &listCommentsUseCase{repo: repo}
}

func (uc *listCommentsUseCase) Execute(ctx context.Context, postID string) ([]CommentOutput, error) {
	if postID == "" {
		return nil, errors.New("post não informado")
	}

	comments, err := uc.repo.FindByPostID(ctx, postID)
	if err != nil {
		return nil, err
	}

	out := make([]CommentOutput, len(comments))
	for i := range comments {
		out[i] = CommentOutput{
			ID:        comments[i].ID,
			Content:   comments[i].Content,
			PostID:    comments[i].PostID,
			UserID:    comments[i].UserID,
			CreatedAt: comments[i].CreatedAt,
		}
	}
	return out, nil
}

type DeleteCommentUseCase interface {
	Execute(ctx context.Context, id, userID string) error
}

type deleteCommentUseCase struct {
	repo Repository
}

func NewDeleteCommentUseCase(repo Repository) DeleteCommentUseCase {
	return &deleteCommentUseCase{repo: repo}
}

func (uc *deleteCommentUseCase) Execute(ctx context.Context, id, userID string) error {
	comment, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return errors.New("comentário não encontrado")
	}

	if comment.UserID != userID {
		return errors.New("você não tem permissão para deletar este comentário")
	}

	return uc.repo.Delete(ctx, id)
}
