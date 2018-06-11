package model

import (
	model "github.com/yang-zzhong/go-model"
	"time"
)

type Comment struct {
	Id          string    `db:"id char(36) pk"`
	Content     string    `db:"content long_text"`
	UserId      string    `db:"user_id char(36)"`
	Ats         []string  `db:"ats text nil"`
	CommentedAt time.Time `db:"commented_at datetime"`
	*model.Base
}

func NewComment() *Comment {
	return model.NewModel(new(Comment)).(*Comment)
}

func (comment *Comment) Instance() *Comment {
	comment.DeclareOne("user", new(User), map[string]string{
		"user_id": "id",
	})
	return comment
}
