package model

import (
	. "github.com/yang-zzhong/go-model"
)

type Theme struct {
	UserId          uint32 `db:"user_id bigint pk"`
	Name            string `db:"name varchar(128) uk"`
	HeaderBgImageId string `db:"header_bg_image_id char(32) nil"`
	InfoBgImageId   string `db:"info_bg_image_id char(32) nil"`
	BgColor         string `db:"bg_color varchar(16) nil"`
	FgColor         string `db:"fg_color varchar(16) nil"`
	TagFgColor      string `db:"tag_fg_color varchar(16) nil"`
	TagBgColor      string `db:"tag_bg_color varchar(16) nil"`
	HeaderBgColor   string `db:"header_bg_color varchar(16) nil"`
	HeaderFgColor   string `db:"header_fg_color varchar(16) nil"`
	PaperBgColor    string `db:"paper_bg_color varchar(16) nil"`
	PaperFgColor    string `db:"paper_fg_color varchar(16) nil"`
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
