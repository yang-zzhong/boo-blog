package model

import (
	"github.com/google/uuid"
	model "github.com/yang-zzhong/go-model"
	"reflect"
	"time"
)

type Comment struct {
	Id             uint32    `db:"id bigint pk"`
	Content        string    `db:"content long_text"`
	UserId         uint32    `db:"user_id bigint"`
	Reply          uint32    `db:"reply bigint"`
	Ats            []string  `db:"ats text nil"`
	BlogId         uint32    `db:"blog_id bigint"`
	CommentId      uint32    `db:"comment_Id bigint nil"`
	CommentAllowed bool      `db:"comment_allowed int"`
	CommentedAt    time.Time `db:"commented_at datetime"`
	*model.Base
}

func NewComment() *Comment {
	comment := model.NewModel(new(Comment)).(*Comment)
	comment.DeclareOne("user", new(User), map[string]string{
		"user_id": "id",
	})
	comment.DeclareOne("blog", new(Blog), map[string]string{
		"blog_id": "id",
	})
	comment.DeclareOne("reply", new(User), map[string]string{
		"reply": "id",
	})
	return comment
}

func (comment *Comment) DBValue(colname string, value interface{}) interface{} {
	if colname == "ats" {
		return nullArrayDBValue(value)
	}
	return value
}

func (comment *Comment) Value(colname string, value interface{}) (result reflect.Value, catch bool) {
	if colname == "ats" {
		catch = true
		result = nullArrayValue(value)
		return
	}
	catch = false
	return
}

func (comment *Comment) Instance() *Comment {
	comment.Id = uuid.New().ID()
	comment.CommentedAt = time.Now()
	return comment
}
