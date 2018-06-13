package model

import (
	"github.com/google/uuid"
	model "github.com/yang-zzhong/go-model"
	"time"
)

type Comment struct {
	Id          uint32    `db:"id bigint pk"`
	Content     string    `db:"content long_text"`
	UserId      uint32    `db:"user_id bigint"`
	Ats         []string  `db:"ats text nil"`
	CommentedAt time.Time `db:"commented_at datetime"`
	*model.Base
}

func NewComment() *Comment {
	comment := model.NewModel(new(Comment)).(*Comment)
	comment.DeclareOne("user", new(User), map[string]string{
		"user_id": "id",
	})
	return comment
}

func (comment *Comment) Instance() *Comment {
	comment.Id = uuid.New().ID()
	comment.CommentedAt = time.Now()
	return comment
}
