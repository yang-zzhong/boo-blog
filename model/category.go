package model

import (
	"database/sql"
	helpers "github.com/yang-zzhong/go-helpers"
	repo "github.com/yang-zzhong/go-model"
	"reflect"
	"strings"
	"time"
)

type Category struct {
	Id        string    `db:"id char(36) pk"`
	Name      string    `db:"name varchar(64)"`
	Intro     string    `db:"intro varchar(512) nil"`
	UserId    string    `db:"user_id char(36)"`
	Tags      []string  `db:"tags varchar(512) nil"`
	CreatedAt time.Time `db:"created_at datetime"`
	UpdatedAt time.Time `db:"updated_at datetime"`
}

func (ig *Category) TableName() string {
	return "category"
}

func (ig *Category) PK() string {
	return "Id"
}

func (ig *Category) NewId() interface{} {
	return helpers.RandString(32)
}

func (ig *Category) DBValue(fieldName string, val interface{}) interface{} {
	if fieldName == "tags" {
		return sql.NullString{strings.Join(val.([]string), ","), true}
	}
	return val
}

func (ig *Category) Value(fieldName string, val interface{}) (result reflect.Value, catched bool) {
	if fieldName == "tags" {
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

func NewCategory() *Category {
	return CreateModel(new(Category)).(*Category)
}

func NewCategoryRepo() (*repo.Repo, error) {
	return CreateRepo(new(Category))
}
