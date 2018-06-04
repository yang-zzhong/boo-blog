package model

import (
	"database/sql"
	. "github.com/yang-zzhong/go-model"
	"time"
)

type Tag struct {
	Name      string    `db:"title varchar(64) pk"`
	Intro     string    `db:"intro varchar(256) nil"`
	UserId    string    `db:"user_id char(32)"`
	IntroUrl  string    `db:"intro_url varchar(512) nil"`
	CreatedAt time.Time `db:"created_at datetime"`
	UpdatedAt time.Time `db:"updated_at datetime"`
	*Base
}

func (tag *Tag) PK() string {
	return "title"
}

func (tag *Tag) NewId() interface{} {
	return ""
}

func (tag *Tag) TableName() string {
	return "tag"
}

func (tag *Tag) One(name string) (interface{}, error) {
	return One(tag.Base, tag, name)
}

func NewTag() *Tag {
	tag := CreateModel(new(Tag)).(*Tag)
	tag.DeclareOne("user", new(User), map[string]string{
		"user_id": "id",
	})

	return tag
}

func NewTagRepo() (*Repo, error) {
	return CreateRepo(NewTag())
}
