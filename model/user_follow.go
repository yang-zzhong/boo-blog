package model

import (
	"github.com/google/uuid"
	model "github.com/yang-zzhong/go-model"
	"time"
)

type UserFollow struct {
	Id         uint32    `db:"id bigint pk"`
	UserId     uint32    `db:"user_id bigint"`
	FollowedBy uint32    `db:"followed_by bigint"`
	FollowedAt time.Time `db:"followed_at datetime"`
	*model.Base
}

func (userFollow *UserFollow) TableName() string {
	return "user_follows"
}

func NewUserFollow() *UserFollow {
	userFollow := model.NewModel(new(UserFollow)).(*UserFollow)
	userFollow.DeclareOne("follow", new(User), model.Nexus{
		"id": "user_id",
	})

	return userFollow
}

func (userFollow *UserFollow) Instance() *UserFollow {
	userFollow.Id = uuid.New().ID()
	return userFollow
}
