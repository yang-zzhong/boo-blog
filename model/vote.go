package model

import (
	"github.com/google/uuid"
	model "github.com/yang-zzhong/go-model"
	"time"
)

const (
	VOTE_BLOG    = "blog"
	VOTE_COMMENT = "comment"
)

type Vote struct {
	Id         uint32    `db:"id bigint pk"`
	UserId     uint32    `db:"user_id bigint"`
	TargetType string    `db:"target_type varchar(16)"`
	TargetId   uint32    `db:"target_id bigint"`
	Vote       int8      `db:"vote smallint"`
	VotedAt    time.Time `db:"voted_at datetime"`
	*model.Base
}

func (vote *Vote) TableName() string {
	return "vote"
}

func NewVote() *Vote {
	return model.NewModel(new(Vote)).(*Vote)
}

func (vote *Vote) Instance() *Vote {
	vote.Id = uuid.New().ID()

	return vote
}
