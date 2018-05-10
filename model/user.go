package model

import (
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	helpers "github.com/yang-zzhong/go-helpers"
	. "github.com/yang-zzhong/go-model"
	"time"
)

type User struct {
	Id          string         `db:"id char(32) pk"`
	Name        string         `db:"name varchar(128) uk"`
	NickName    sql.NullString `db:"nickname varchar(128) nil"`
	EmailAddr   sql.NullString `db:"email_addr varchar(128) nil"`
	PhoneNumber sql.NullString `db:"phone_number varchar(128) nil"`
	Password    string         `db:"password varchar(128)"`
	Salt        string         `db:"salt char(8)"`
	CreatedAt   time.Time      `db:"created_at datetime"`
	UpdatedAt   time.Time      `db:"updated_at datetime"`
}

func (user *User) TableName() string {
	return "user"
}

func (user *User) PK() string {
	return "id"
}

func (user *User) NewId() interface{} {
	return helpers.RandString(32)
}

func (user *User) Encrypt(str string) string {
	md5Sumb := md5.Sum(([]byte)(str + user.Salt))
	return hex.EncodeToString(md5Sumb[:])
}

func NewUser() *User {
	user := CreateModel(new(User)).(*User)
	user.Salt = helpers.RandString(8)
	return user
}

func NewUserRepo() (userRepo *Repo, err error) {
	return CreateRepo(new(User))
}
