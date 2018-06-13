package model

import (
	. "github.com/yang-zzhong/go-model"
)

type Theme struct {
	UserId        uint32 `db:"user_id bigint pk"`
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

func NewTheme() *Theme {
	theme := NewModel(new(Theme)).(*Theme)
	theme.DeclareOne("user", new(User), map[string]string{
		"user_id": "id",
	})

	return theme
}
