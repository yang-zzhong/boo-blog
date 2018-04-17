package model

import (
	helpers "github.com/yang-zzhong/go-helpers"
	. "github.com/yang-zzhong/go-model"
	"time"
)

type User struct {
	Id          string    `db:"id varchar(128) pk"`
	Name        string    `db:"name varchar(128) uk"`
	NickName    string    `db:"nickname varchar(128) nil"`
	EmailAddr   string    `db:"email_addr varchar(128) nil"`
	PhoneNumber string    `db:"phone_number varchar(128) nil"`
	CreatedAt   time.Time `db:"created_at datetime"`
	UpdatedAt   time.Time `db:"updated_at datetime"`
	DeletedAt   time.Time `db:"deleted_at datetime nil"`
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

func NewUser() *User {
	return CreateModel(new(User)).(*User)
}

func NewUserRepo() (userRepo *Repo, err error) {
	return CreateRepo(new(User))
}
