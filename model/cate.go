package model

import (
	model "github.com/yang-zzhong/go-model"
	"reflect"
	"time"
)

type Cate struct {
	Id        string    `db:"id char(36) pk"`
	Name      string    `db:"name varchar(64)"`
	Intro     string    `db:"intro varchar(512) nil"`
	UserId    string    `db:"user_id char(36)"`
	Tags      []string  `db:"tags varchar(512) nil"`
	CreatedAt time.Time `db:"created_at datetime"`
	UpdatedAt time.Time `db:"updated_at datetime"`
	*model.Base
}

func (ig *Cate) TableName() string {
	return "category"
}

func (ig *Cate) DBValue(colname string, value interface{}) interface{} {
	if colname == "tags" {
		return nullArrayDBValue(value)
	}
	return value
}

func (ig *Cate) Value(colname string, value interface{}) (result reflect.Value, catch bool) {
	if colname == "tags" {
		catch = true
		result = nullArrayValue(value)
		return
	}
	catch = false
	return
}

func NewCate() *Cate {
	cate := model.NewModel(new(Cate)).(*Cate)
	cate.DeclareOne("user", new(User), map[string]string{
		"user_id": "id",
	})
	cate.DeclareMany("blogs", new(Blog), map[string]string{
		"id": "cate_id",
	})

	return cate
}

func (ig *Cate) Instance() *Cate {
	return Instance(ig).(*Cate)
}
