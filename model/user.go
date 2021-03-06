package model

import (
	"crypto/md5"
	"encoding/hex"
	"github.com/google/uuid"
	helpers "github.com/yang-zzhong/go-helpers"
	. "github.com/yang-zzhong/go-model"
	"time"
)

type User struct {
	Id              uint32    `db:"id bigint pk"`
	Name            string    `db:"name varchar(128) uk"`
	NickName        string    `db:"nickname varchar(128) nil"`
	EmailAddr       string    `db:"email_addr varchar(128) nil"`
	PhoneNumber     string    `db:"phone_number varchar(128) nil"`
	PortraitImageId string    `db:"portrait_image_id varchar(32) nil"`
	Password        string    `db:"password varchar(128) protected"`
	Salt            string    `db:"salt char(8) protected"`
	CreatedAt       time.Time `db:"created_at datetime"`
	UpdatedAt       time.Time `db:"updated_at datetime"`
	*Base
}

func (user *User) TableName() string {
	return "users"
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
	user.Id = uuid.New().ID()
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	return user
}
