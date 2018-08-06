package model

import (
	"github.com/google/uuid"
	model "github.com/yang-zzhong/go-model"
	"reflect"
	"time"
)

type Cate struct {
	Id        uint32    `db:"id bigint pk"`
	Name      string    `db:"name varchar(64)"`
	Intro     string    `db:"intro varchar(512) nil"`
	UserId    uint32    `db:"user_id bigint"`
	Tags      []string  `db:"tags varchar(32)[] nil"`
	CreatedAt time.Time `db:"created_at timestamp"`
	UpdatedAt time.Time `db:"updated_at timestamp"`
	*model.Base
}

func (ig *Cate) TableName() string {
	return "cates"
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
	cate.DeclareOne("user", new(User), model.Nexus{
		"id": "user_id",
	})
	cate.DeclareMany("blogs", new(Blog), model.Nexus{
		"cate_id": "id",
	})
	cate.OnUpdate(func(i interface{}) error {
		i.(*Cate).UpdatedAt = time.Now()
		return nil
	})

	return cate
}

func (ig *Cate) Instance() *Cate {
	ig.Id = uuid.New().ID()
	ig.CreatedAt = time.Now()
	ig.UpdatedAt = time.Now()

	return ig
}
