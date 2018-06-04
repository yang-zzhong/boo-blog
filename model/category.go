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
	*repo.Base
}

func (ig *Category) TableName() string {
	return "category"
}

func (ig *Category) PK() string {
	return "id"
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

func (ig *Category) One(name string) (interface{}, error) {
	return repo.One(ig.Base, ig, name)
}

func (ig *Category) Many(name string) (map[interface{}]interface{}, error) {
	return repo.Many(ig.Base, ig, name)
}

func NewCategory() *Category {
	cate := CreateModel(new(Category)).(*Category)
	cate.DeclareOne("user", new(User), map[string]string{
		"user_id": "id",
	})
	cate.DeclareMany("blogs", new(Article), map[string]string{
		"id": "cate_id",
	})

	return cate
}

func NewCategoryRepo() (*repo.Repo, error) {
	return CreateRepo(NewCategory())
}
