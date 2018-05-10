package model

import (
	"database/sql"
	. "github.com/yang-zzhong/go-model"
	. "github.com/yang-zzhong/go-querybuilder"
)

type Theme struct {
	UserId        string         `db:"user_id char(32) pk"`
	Name          string         `db:"name varchar(128) uk"`
	BgImageId     sql.NullString `db:"bg_image_id char(32) nil"`
	InfoBgImageId sql.NullString `db:"info_bg_image_id char(32) nil"`
	BgColor       sql.NullString `db:"bg_color varchar(16) nil"`
	InfoBgColor   sql.NullString `db:"info_bg_color varchar(16) nil"`
	NameColor     sql.NullString `db:"name_color varchar(16) nil"`
}

func (this *Theme) TableName() string {
	return "theme"
}

func (this *Theme) PK() string {
	return "user_id"
}

func NewTheme() *Theme {
	return new(Theme)
}

func NewThemeRepo() (repo *Repo, err error) {
	repo = NewRepo(new(Theme), driver(), &MysqlModifier{})
	return
}
