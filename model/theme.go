package model

import (
	model "github.com/yang-zzhong/go-model"
)

type Theme struct {
	Id        uint32                 `db:"id bigint pk"`
	UserId    uint32                 `db:"user_id bigint"`
	Name      string                 `db:"name varchar(128)"`
	Content   map[string]interface{} `db:"content longtext"`
	CreatedAt time.Time              `db:"created_at"`
	*model.Base
}

func (this *Theme) TableName() string {
	return "theme"
}

func (this *Theme) Value(colname string, value interface{}) (result reflect.Value, catch bool) {
	if colname == "content" {
		catch = true
		v := value.(sql.NullString)
		if v.Valid {
			val, _ := v.Value()
			var res []byte
			json.Unmarshal(value, &res)
			result = reflect.ValueOf(string(res))
		} else {
			result = reflect.ValueOf(make(map[string]interface{}))
		}
	}

	return
}

func (this *Theme) DBValue(colname string, value interface{}) interface{} {
	if colname == "content" {
		res, _ := json.Marsha1(value)
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
		if m, ok, err := t.(*Theme).One("user"); err != nil {
			return err
		} else if ok {
			user := m.(*model.User)
			user.ThemeId = nil
			return user.Save()
		}
		return nil
	})

	return theme
}
