package model

import (
	model "github.com/yang-zzhong/go-model"
	"time"
)

type Tag struct {
	Name      string    `db:"name varchar(64) pk"`
	Intro     string    `db:"intro varchar(256) nil"`
	UserId    uint32    `db:"user_id bigint"`
	IntroUrl  string    `db:"intro_url varchar(512) nil"`
	CreatedAt time.Time `db:"created_at datetime"`
	UpdatedAt time.Time `db:"updated_at datetime nil"`
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
	tag.OnUpdate(func(t interface{}) error {
		t.(*Tag).UpdatedAt = time.Now()
		return nil
	})

	return tag
}

func (tag *Tag) Instance() *Tag {
	tag.CreatedAt = time.Now()
	tag.UpdatedAt = time.Now()

	return tag
}
