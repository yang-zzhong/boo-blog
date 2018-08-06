package model

import (
	"github.com/google/uuid"
	model "github.com/yang-zzhong/go-model"
	"reflect"
	"time"
)

type UserImage struct {
	Id        uint32    `db:"id bigint pk"`
	UserId    uint32    `db:"user_id bigint"`
	Hash      string    `db:"hash char(32)"`
	CateId    uint32    `db:"cate_id bigint nil"`
	Tags      []string  `db:"tags varchar(32)[] nil"`
	CreatedAt time.Time `db:"created_at timestamp"`
	UpdatedAt time.Time `db:"updated_at timestamp"`
	*model.Base
}

func (this *UserImage) TableName() string {
	return "user_image"
}

func (this *UserImage) DBValue(colname string, value interface{}) interface{} {
	if colname == "tags" {
		return nullArrayDBValue(value)
	}
	return value
}

func (this *UserImage) Value(colname string, value interface{}) (result reflect.Value, catch bool) {
	if colname == "tags" {
		catch = true
		result = nullArrayValue(value)
		return
	}
	catch = false
	return
}

func NewUserImage() *UserImage {
	ui := model.NewModel(new(UserImage)).(*UserImage)
	ui.DeclareOne("image", new(Image), model.Nexus{
		"hash": "hash",
	})
	ui.DeclareOne("cate", new(Cate), model.Nexus{
		"id": "cate_id",
	})
	ui.OnUpdate(func(u interface{}) error {
		u.(*UserImage).UpdatedAt = time.Now()
		return nil
	})

	return ui
}

func (ui *UserImage) Instance() *UserImage {
	ui.Id = uuid.New().ID()
	ui.CreatedAt = time.Now()
	ui.UpdatedAt = time.Now()

	return ui
}
