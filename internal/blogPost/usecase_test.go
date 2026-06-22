package blogPost

import (
	"Gblog/pkg/blogstatus"
	"context"
	"errors"
	"testing"
)

type stubRepo struct {
	BlogPostRepository
	post *BlogPost
}

func (s *stubRepo) FindByID(_ context.Context, id string) (*BlogPost, error) {
	if s.post == nil || s.post.ID != id {
		return nil, errors.New("not found")
	}
	return s.post, nil
}

func (s *stubRepo) Update(_ context.Context, post *BlogPost) error {
	s.post = post
	return nil
}

func TestPublishPostUseCase_Execute(t *testing.T) {
	tests := []struct {
		name    string
		post    *BlogPost
		wantErr bool
	}{
		{"post not found", nil, true},
		{"already published", &BlogPost{ID: "1", Status: blogstatus.Published, Content: "long enough content here"}, true},
		{"content too short", &BlogPost{ID: "1", Status: blogstatus.Draft, Content: "short"}, true},
		{"valid publish", &BlogPost{ID: "1", Status: blogstatus.Draft, Content: "long enough content for publishing"}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &stubRepo{post: tt.post}
			uc := NewPublishPostUseCase(repo)
			err := uc.Execute(context.Background(), "1")
			if (err != nil) != tt.wantErr {
				t.Errorf("Execute() error = %v, wantErr = %v", err, tt.wantErr)
			}
		})
	}
}
