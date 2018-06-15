package model

import (
	"encoding/json"
	"github.com/google/uuid"
	model "github.com/yang-zzhong/go-model"
	"reflect"
	"time"
)

type Theme struct {
	Id        uint32            `db:"id bigint pk"`
	UserId    uint32            `db:"user_id bigint"`
	Name      string            `db:"name varchar(128)"`
	Content   map[string]string `db:"content longtext"`
	CreatedAt time.Time         `db:"created_at datetime"`
	*model.Base
}

func (this *Theme) TableName() string {
	return "theme"
}

func (this *Theme) Value(colname string, value interface{}) (result reflect.Value, catch bool) {
	if colname == "content" {
		catch = true
		val := value.(string)
		var res map[string]string
		json.Unmarshal([]byte(val), &res)
		result = reflect.ValueOf(res)
	}

	return
}

func (this *Theme) DBValue(colname string, value interface{}) interface{} {
	if colname == "content" {
		res, _ := json.Marshal(value)
		return string(res)
	}

	return value
}

func (this *Theme) Instance() *Theme {
	this.Id = uuid.New().ID()
	this.CreatedAt = time.Now()
	return this
}

func NewTheme() *Theme {
	theme := model.NewModel(new(Theme)).(*Theme)
	theme.DeclareOne("user", new(User), map[string]string{
		"user_id": "id",
	})
	theme.OnDelete(func(t interface{}) error {
		if m, err := t.(*Theme).One("user"); err != nil {
			return err
		} else if m != nil {
			user := m.(*User)
			user.ThemeId = 0
			return user.Save()
		}
		return nil
	})

	return theme
}
