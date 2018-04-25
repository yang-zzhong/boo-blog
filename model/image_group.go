package model

import (
	"database/sql"
	helpers "github.com/yang-zzhong/go-helpers"
	repo "github.com/yang-zzhong/go-model"
	"reflect"
	"strings"
	"time"
)

type ImageGroup struct {
	Id     string   `db:"id char(36) pk"`
	Name   string   `db:"name varchar(64)"`
	Intro  string   `db:"intro varchar(512) nil"`
	UserId string   `db:"user_id char(36)"`
	TagIds []string `db:"tag_ids varchar(512) nil"`

	CreatedAt time.Time `db:"created_at datetime"`
	UpdatedAt time.Time `db:"updated_at datetime"`
}

func (ig *ImageGroup) TableName() string {
	return "image_group"
}

func (ig *ImageGroup) PK() string {
	return "Id"
}

func (ig *ImageGroup) NewId() interface{} {
	return helpers.RandString(32)
}

func (ig *ImageGroup) DBValue(fieldName string, val interface{}) interface{} {
	if fieldName == "tag_ids" {
		return sql.NullString{strings.Join(val.([]string), ","), true}
	}
	return val
}

func (ig *ImageGroup) Value(fieldName string, val interface{}) (result reflect.Value, catched bool) {
	if fieldName == "tag_ids" {
		catched = true
		value, _ := val.(sql.NullString).Value()
		if value != nil {
			result = reflect.ValueOf(strings.Split(value.(string), ","))
			return
		}
		result = reflect.ValueOf([]string{})
		return
	}
	catched = false
	return
}

func NewImageGroup() *ImageGroup {
	return CreateModel(new(ImageGroup)).(*ImageGroup)
}

func NewImageGroupRepo() (*repo.Repo, error) {
	return CreateRepo(new(ImageGroup))
}
