package blogPost

import (
	"Gblog/internal/tag"
	"Gblog/pkg/blogstatus"
	"bytes"
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/yuin/goldmark"
)

var slugRe = regexp.MustCompile(`[^a-z0-9]+`)

func (uc *createPostUseCase) resolveTags(ctx context.Context, names []string) ([]tag.Tag, error) {
	// ponytail: linear scan, fine for <100 tags per post
	var tags []tag.Tag
	for _, name := range names {
		t, err := uc.tagRepo.FindByName(ctx, name)
		if err != nil {
			t = &tag.Tag{Name: name}
			if err := uc.tagRepo.Create(ctx, t); err != nil {
				return nil, err
			}
		}
		tags = append(tags, *t)
	}
	return tags, nil
}

func (uc *updatePostUseCase) resolveTags(ctx context.Context, names []string) ([]tag.Tag, error) {
	var tags []tag.Tag
	for _, name := range names {
		t, err := uc.tagRepo.FindByName(ctx, name)
		if err != nil {
			t = &tag.Tag{Name: name}
			if err := uc.tagRepo.Create(ctx, t); err != nil {
				return nil, err
			}
		}
		tags = append(tags, *t)
	}
	return tags, nil
}

func slugify(s string) string {
	return strings.Trim(slugRe.ReplaceAllString(strings.ToLower(s), "-"), "-")
}

// ---------------------------------------------------------
// DTOs (Data Transfer Objects)
// ---------------------------------------------------------

// ponytail: minimal interface to decouple blogPost from user package
type UserExistenceChecker interface {
	UserExists(ctx context.Context, userID string) bool
}

// ---------------------------------------------------------
// CASO DE USO 1: Criar um Blog Post
// ---------------------------------------------------------

type CreatePostUseCase interface {
	Execute(ctx context.Context, input CreatePostInputDTO) (*BlogPost, error)
}

type createPostUseCase struct {
	repo      BlogPostRepository
	userCheck UserExistenceChecker
	tagRepo   tag.Repository
}

func NewCreatePostUseCase(repo BlogPostRepository, userCheck UserExistenceChecker, tagRepo tag.Repository) CreatePostUseCase {
	return &createPostUseCase{repo: repo, userCheck: userCheck, tagRepo: tagRepo}
}

func (uc *createPostUseCase) Execute(ctx context.Context, input CreatePostInputDTO) (*BlogPost, error) {
	statusDefault := blogstatus.Draft

	if input.Title == "" {
		return nil, errors.New("o título do post é obrigatório")
	}
	if input.UserID == "" {
		return nil, errors.New("o ID do usuário é obrigatório")
	}

	if !uc.userCheck.UserExists(ctx, input.UserID) {
		return nil, errors.New("usuário informado não existe")
	}

	tags, _ := uc.resolveTags(ctx, input.Tags)

	postEntity := &BlogPost{
		Title:   input.Title,
		Slug:    fmt.Sprintf("%s-%d", slugify(input.Title), time.Now().UnixMilli()),
		Content: input.Content,
		Status:  statusDefault,
		UserID:  input.UserID,
		Tags:    tags,
	}

	if err := uc.repo.Create(ctx, postEntity); err != nil {
		return nil, err
	}
	return postEntity, nil
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
	Execute(ctx context.Context, id, userID string, input UpdatePostInputDTO) error
}

type updatePostUseCase struct {
	repo    BlogPostRepository
	tagRepo tag.Repository
}

func NewUpdatePostUseCase(repo BlogPostRepository, tagRepo tag.Repository) UpdatePostUseCase {
	return &updatePostUseCase{repo: repo, tagRepo: tagRepo}
}

func (uc *updatePostUseCase) Execute(ctx context.Context, id, userID string, input UpdatePostInputDTO) error {
	post, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return errors.New("post não encontrado")
	}

	if !post.IsActive {
		return errors.New("não é possível editar um post inativo/deletado")
	}

	if post.UserID != userID {
		return errors.New("você não tem permissão para editar este post")
	}

	post.Title = input.Title
	post.Content = input.Content
	post.Slug = slugify(input.Title)

	tags, _ := uc.resolveTags(ctx, input.Tags)
	if input.Tags != nil {
		post.Tags = tags
	}

	return uc.repo.Update(ctx, post)
}

// ---------------------------------------------------------
// CASO DE USO 4: Deletar (Soft Delete)
// ---------------------------------------------------------

type DeletePostUseCase interface {
	Execute(ctx context.Context, id, userID string) error
}

type deletePostUseCase struct {
	repo BlogPostRepository
}

func NewDeletePostUseCase(repo BlogPostRepository) DeletePostUseCase {
	return &deletePostUseCase{repo: repo}
}

func (uc *deletePostUseCase) Execute(ctx context.Context, id, userID string) error {
	post, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return errors.New("post não encontrado")
	}

	if post.UserID != userID {
		return errors.New("você não tem permissão para deletar este post")
	}

	post.IsActive = false

	return uc.repo.Update(ctx, post)
}

func toPostOutput(p *BlogPost) PostOutput {
	tags := make([]tag.TagOutput, len(p.Tags))
	for i := range p.Tags {
		tags[i] = tag.TagOutput{ID: p.Tags[i].ID, Name: p.Tags[i].Name}
	}
	var buf bytes.Buffer
	if err := goldmark.Convert([]byte(p.Content), &buf); err != nil {
		buf.WriteString(p.Content)
	}
	return PostOutput{
		ID:          p.ID,
		Title:       p.Title,
		Content:     p.Content,
		ContentHTML: buf.String(),
		Status:      string(p.Status),
		Slug:        p.Slug,
		UserID:      p.UserID,
		Tags:        tags,
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
	}
}

type GetPostsUseCase interface {
	Execute(ctx context.Context, limit, offset int) ([]PostOutput, error)
}

type getPostsUseCase struct {
	repo BlogPostRepository
}

func NewGetPostsUseCase(repo BlogPostRepository) GetPostsUseCase {
	return &getPostsUseCase{repo: repo}
}

func (uc *getPostsUseCase) Execute(ctx context.Context, limit, offset int) ([]PostOutput, error) {
	posts, err := uc.repo.FetchAll(ctx, limit, offset)
	if err != nil {
		return nil, err
	}
	out := make([]PostOutput, len(posts))
	for i := range posts {
		out[i] = toPostOutput(&posts[i])
	}
	return out, nil
}

type GetPostBySlugUseCase interface {
	Execute(ctx context.Context, slug string) (*PostOutput, error)
}

type getPostBySlugUseCase struct {
	repo BlogPostRepository
}

func NewGetPostBySlugUseCase(repo BlogPostRepository) GetPostBySlugUseCase {
	return &getPostBySlugUseCase{repo: repo}
}

func (uc *getPostBySlugUseCase) Execute(ctx context.Context, slug string) (*PostOutput, error) {
	post, err := uc.repo.GetBySlug(ctx, slug)
	if err != nil {
		return nil, err
	}
	out := toPostOutput(post)
	return &out, nil
}
