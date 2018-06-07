package model

import (
	"database/sql"
	model "github.com/yang-zzhong/go-model"
	"reflect"
	"strings"
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

func (ig *Cate) DBValue(fieldName string, val interface{}) interface{} {
	if fieldName == "tags" {
		return sql.NullString{strings.Join(val.([]string), ","), true}
	}
	return val
}

func (ig *Cate) Value(fieldName string, val interface{}) (result reflect.Value, catched bool) {
	if fieldName == "tags" {
		catched = true
		v := val.(sql.NullString)
		if v.Valid {
			str, _ := v.Value()
			result = reflect.ValueOf(strings.Split(str.(string), ","))
		} else {
			result = reflect.ValueOf([]string{})
		}
		return
	}
	catched = false
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
