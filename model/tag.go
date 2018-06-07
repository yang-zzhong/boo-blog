package model

import (
	model "github.com/yang-zzhong/go-model"
	"time"
)

type Tag struct {
	Name      string    `db:"title varchar(64) pk"`
	Intro     string    `db:"intro varchar(256) nil"`
	UserId    string    `db:"user_id char(32)"`
	IntroUrl  string    `db:"intro_url varchar(512) nil"`
	CreatedAt time.Time `db:"created_at datetime"`
	UpdatedAt time.Time `db:"updated_at datetime"`
	*model.Base
}

func (tag *Tag) TableName() string {
	return "tags"
}

func NewTag() *Tag {
	tag := model.NewModel(new(Tag)).(*Tag)
	tag.DeclareOne("user", new(User), map[string]string{
		"user_id": "id",
	})

	return tag
}
