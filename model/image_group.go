package model

import (
	helpers "github.com/yang-zzhong/go-helpers"
	repo "github.com/yang-zzhong/go-model"
	"time"
)

type ImageGroup struct {
	Id     string `db:id char(36) pk`
	Name   string `db:name varchar(64)`
	Intro  string `db:intro varchar(512) nil`
	UserId string `db:user_id char(36)`
	TagIds string `db:tag_ids varchar(512) nil`

	CreatedAt time.Time `db:created_at datetime`
	UpdatedAt time.Time `db:updated_at datetime`
}

func (ig *ImageGroup) TableName() string {
	return "image_group"
}

func (ig *ImageGroup) PK() string {
	return "id"
}

func (ig *ImageGroup) NewId() string {
	return helpers.RandString(32)
}

func NewImageGroup() *ImageGroup {
	return CreateModel(new(ImageGroup)).(*ImageGroup)
}

func NewImageGroupRepo() (*repo.Repo, error) {
	return CreateRepo(new(ImageGroup))
}
