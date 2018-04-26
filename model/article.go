package model

import (
	helpers "github.com/yang-zzhong/go-helpers"
	"strings"
	"time"
)

type Article struct {
	Id      string   `db:"id char(32) pk"`
	Title   string   `db:"title varchar(64)"`
	UserId  string   `db:"user_id char(32)"`
	GroupId string   `db:"group_id char(32) nil"`
	TagIds  []string `db:"tag_ids varchar(256)"`

	PublishedAt time.Time `db:"pushlished_at datetime"`
	ModifiedAt  time.Time `db:"modified_at datetime"`
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
	if fieldName == "tag_ids" {
		result, _ := strings.Join(value.([]string), ",")
		return result
	}
	return value
}

func (atl *Article) Value(fieldName string, value interface{}) (result interface{}, catched bool) {
	if fieldName == "tag_ids" {
		catched = true
		val, _ := val.(sql.NullString).Value()
		if val != nil {
			result = reflect.ValueOf(strings.Split(val.(string), ","))
			return
		}
		result = reflect.ValueOf([]string{})
		return
	}
	catched = false
	return
}
