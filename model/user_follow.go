package model

import (
	"github.com/google/uuid"
	model "github.com/yang-zzhong/go-model"
	"time"
)

type UserFollow struct {
	Id         uint32    `db:"id bigint pk"`
	UserId     uint32    `db:"user_id bigint"`
	Followed   uint32    `db:"followed bigint"`
	FollowedAt time.Time `db:"followed_at timestamp"`
	*model.Base
}

func (userFollow *UserFollow) TableName() string {
	return "user_follows"
}

func NewUserFollow() *UserFollow {
	userFollow := model.NewModel(new(UserFollow)).(*UserFollow)
	//
	// followed关联的用户被谁关注
	//
	userFollow.DeclareOne("followed", new(User), model.Nexus{
		"id": "user_id",
	})
	//
	// user_id关联的用户关注的用户
	//
	userFollow.DeclareOne("following", new(User), model.Nexus{
		"id": "followed",
	})

	return userFollow
}

func (userFollow *UserFollow) Instance() *UserFollow {
	userFollow.Id = uuid.New().ID()
	return userFollow
}
