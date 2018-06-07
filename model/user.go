package model

import (
	"crypto/md5"
	"encoding/hex"
	helpers "github.com/yang-zzhong/go-helpers"
	. "github.com/yang-zzhong/go-model"
	"time"
)

type User struct {
	Id          string    `db:"id char(32) pk"`
	Name        string    `db:"name varchar(128) uk"`
	NickName    string    `db:"nickname varchar(128) nil"`
	EmailAddr   string    `db:"email_addr varchar(128) nil"`
	PhoneNumber string    `db:"phone_number varchar(128) nil"`
	Password    string    `db:"password varchar(128)"`
	Salt        string    `db:"salt char(8)"`
	CreatedAt   time.Time `db:"created_at datetime"`
	UpdatedAt   time.Time `db:"updated_at datetime"`
	*Base
}

func (user *User) TableName() string {
	return "user"
}

func (user *User) Encrypt(str string) string {
	md5Sumb := md5.Sum(([]byte)(str + user.Salt))
	return hex.EncodeToString(md5Sumb[:])
}

func NewUser() *User {
	user := NewModel(new(User)).(*User)
	user.DeclareMany("blogs", new(Blog), map[string]string{
		"id": "user_id",
	})
	user.DeclareMany("images", new(UserImage), map[string]string{
		"id": "user_id",
	})
	user.DeclareMany("cates", new(Cate), map[string]string{
		"id": "user_id",
	})
	user.DeclareOne("theme", new(Theme), map[string]string{
		"id": "user_id",
	})
	return user
}

func (user *User) Instance() *User {
	user.Salt = helpers.RandString(8)
	return user
}
