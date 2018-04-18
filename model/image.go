package model

import (
	helpers "github.com/yang-zzhong/go-helpers"
	. "github.com/yang-zzhong/go-model"
)

type Image struct {
	Id        string   `db:"id varchar(128) pk"`
	With      int      `db:"with int"`
	Height    int      `db:"height int"`
	MimeType  string   `db:"mime_type varchar(64)"`
	Size      int      `db:"size int"`
	GroupId   string   `db:"group_id varchar(128) nil"`
	CreatedAt NullTime `db:"created_at datetime"`
	UpdatedAt NullTime `db:"updated_at datetime"`
	DeletedAt NullTime `db:"deleted_at datetime nil"`
}

func (image *Image) TableName() string {
	return "image"
}

func (image *Image) PK() string {
	return "id"
}

func (image *Image) NewId() interface{} {
	return helpers.RandString(32)
}

func NewImage() *Image {
	return CreateModel(new(Image)).(*Image)
}

func NewImageRepo() (imageRepo *Repo, err error) {
	return CreateRepo(new(Image))
}
