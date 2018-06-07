package model

import (
	. "github.com/yang-zzhong/go-model"
	"time"
)

type UserImage struct {
	Id        string    `db:"id char(32) pk"`
	UserId    string    `db:"user_id char(32)"`
	Hash      string    `db:"hash char(32)"`
	GroupId   string    `db:"group_id char(32) nil"`
	Tags      string    `db:"tags varchar(256) nil"`
	CreatedAt time.Time `db:"created_at datetime"`
	UpdatedAt time.Time `db:"updated_at datetime"`
	*Base
}

func (this *UserImage) TableName() string {
	return "user_image"
}

func NewUserImage() *UserImage {
	ui := NewModel(new(UserImage)).(*UserImage)
	ui.DeclareOne("image", new(Image), map[string]string{
		"hash": "hash",
	})
	ui.DeclareOne("cate", new(Cate), map[string]string{
		"group_id": "id",
	})

	return ui
}
