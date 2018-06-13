package model

import (
	"github.com/google/uuid"
	. "github.com/yang-zzhong/go-model"
	"reflect"
	"time"
)

type UserImage struct {
	Id        uint32    `db:"id bigint pk"`
	UserId    uint32    `db:"user_id bigint"`
	Hash      string    `db:"hash char(32)"`
	GroupId   string    `db:"group_id char(32) nil"`
	Tags      []string  `db:"tags varchar(256) nil"`
	CreatedAt time.Time `db:"created_at datetime"`
	UpdatedAt time.Time `db:"updated_at datetime"`
	*Base
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
	ui := NewModel(new(UserImage)).(*UserImage)
	ui.DeclareOne("image", new(Image), map[string]string{
		"hash": "hash",
	})
	ui.DeclareOne("cate", new(Cate), map[string]string{
		"group_id": "id",
	})

	return ui
}

func (ui *UserImage) Instance() *UserImage {
	ui.Id = uuid.New().ID()
	ui.CreatedAt = time.Now()
	ui.UpdatedAt = time.Now()

	return ui
}
