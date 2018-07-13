package model

import (
	helpers "github.com/yang-zzhong/go-helpers"
	model "github.com/yang-zzhong/go-model"
	"time"
)

type UserLogin struct {
	Id       string         `db:"id varcher(128) pk"`
	UserId   string         `db:"user_id varcher(128)"`
	City     string         `db:"city varchar(64) nil"`
	LoginAt  time.Time      `db:"login_at datatime"`
	LogoutAt model.NullTime `db:"logout_at datetime nil"`
}

func (ul *UserLogin) TableName() string {
	return "user_login"
}

func (ul *UserLogin) PK() string {
	return "id"
}

func (ul *UserLogin) NewId() interface{} {
	return helpers.RandString(128)
}

func NewUserLogin() *UserLogin {
	return model.NewModel(new(UserLogin)).(*UserLogin)
}
