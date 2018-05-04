package model

import (
	"database/sql"
	. "github.com/yang-zzhong/go-model"
	"time"
)

type Tag struct {
	Name      string         `db:"title varchar(64) pk"`
	Intro     sql.NullString `db:"intro varchar(256) nil"`
	UserId    string         `db:"user_id char(32)"`
	IntroUrl  sql.NullString `db:"intro_url varchar(512) nil"`
	CreatedAt time.Time      `db:"created_at datetime"`
	UpdatedAt time.Time      `db:"updated_at datetime"`
}

func (tag *Tag) PK() string {
	return "name"
}

func (tag *Tag) TableName() string {
	return "tag"
}

func NewTag() *Tag {
	return CreateModel(new(Tag)).(*Tag)
}

func NewTagRepo() (*Repo, error) {
	return CreateRepo(new(Tag))
}
