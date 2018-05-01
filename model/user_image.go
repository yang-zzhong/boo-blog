package model

import (
	. "github.com/yang-zzhong/go-model"
	"time"
)

type UserImage struct {
	Id      string `db:"id char(32) pk"`
	UserId  string `db:"user_id char(32)"`
	Hash    string `db:"hash char(32)"`
	GroupId string `db:"group_id char(32)"`
	Tags    string `db:"varchar(256) nil"`

	CreatedAt time.Time `db:"created_at datatime"`
	UpdatedAt time.Time `db:"updated_at datatime"`
}

func (this *UserImage) TableName() string {
	return "user_image"
}

func (this *UserImage) PK() string {
	return "id"
}

func NewUserImage() *UserImage {
	return CreateModel(new(UserImage)).(*UserImage)
}

func NewUserImageRepo() (*Repo, error) {
	return CreateRepo(new(UserImage))
}
