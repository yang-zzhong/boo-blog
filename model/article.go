package model

import (
	helpers "github.com/yang-zzhong/go-helpers"
	. "github.com/yang-zzhong/go-model"
	"reflect"
	"strings"
	"time"
)

type Article struct {
	Id        string    `db:"id char(32) pk"`
	Title     string    `db:"title varchar(64)"`
	Content   string    `db:"content text"`
	UserId    string    `db:"user_id char(32)"`
	CateId    string    `db:"cate_id char(32) nil"`
	Tags      []string  `db:"tags varchar(256) nil"`
	CreatedAt time.Time `db:"created_at datetime"`
	UpdatedAt time.Time `db:"updated_at datetime"`
}

func (atl *Article) PK() string {
	return "id"
}

func (atl *Article) TableName() string {
	return "article"
}

func (atl *Article) NewId() interface{} {
	return helpers.RandString(32)
}

func (atl *Article) DBValue(fieldName string, value interface{}) interface{} {
	if fieldName == "tags" {
		result := strings.Join(value.([]string), ",")
		return result
	}
	return value
}

func (atl *Article) Value(fieldName string, value interface{}) (result reflect.Value, catched bool) {
	if fieldName == "tags" {
		catched = true
		val, _ := value.(string)
		if val != "" {
			result = reflect.ValueOf(strings.Split(val, ","))
			return
		}
		result = reflect.ValueOf([]string{})
		return
	}
	catched = false
	return
}

func NewArticle() *Article {
	return CreateModel(new(Article)).(*Article)
}

func NewArticleRepo() (*Repo, error) {
	return CreateRepo(new(Article))
}
