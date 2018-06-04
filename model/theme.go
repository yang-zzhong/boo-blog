package model

import (
	"database/sql"
	. "github.com/yang-zzhong/go-model"
	. "github.com/yang-zzhong/go-querybuilder"
)

type Theme struct {
	UserId        string `db:"user_id char(32) pk"`
	Name          string `db:"name varchar(128) uk"`
	BgImageId     string `db:"bg_image_id char(32) nil"`
	InfoBgImageId string `db:"info_bg_image_id char(32) nil"`
	BgColor       string `db:"bg_color varchar(16) nil"`
	InfoBgColor   string `db:"info_bg_color varchar(16) nil"`
	NameColor     string `db:"name_color varchar(16) nil"`
	*Base
}

func (this *Theme) TableName() string {
	return "theme"
}

func (this *Theme) PK() string {
	return "user_id"
}

func (this *Theme) One(name string) (interface{}, error) {
	return One(this.Base, this, name)
}

func NewTheme() *Theme {
	theme := new(Theme)
	theme.DeclareOne("user", new(User), map[string]string{
		"user_id": "id",
	})

	return theme
}

func NewThemeRepo() (repo *Repo, err error) {
	repo = NewRepo(NewTheme(), driver(), &MysqlModifier{})
	return
}
