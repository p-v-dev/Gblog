package comment

import "time"

type CreateCommentInput struct {
	Content string `json:"content"`
	PostID  string `json:"-"`
	UserID  string `json:"-"`
}

type CommentOutput struct {
	ID        string    `json:"id"`
	Content   string    `json:"content"`
	PostID    string    `json:"post_id"`
	UserID    string    `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
}
