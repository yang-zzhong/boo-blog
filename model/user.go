package model

import (
	"crypto/md5"
	"encoding/hex"
	"github.com/google/uuid"
	helpers "github.com/yang-zzhong/go-helpers"
	model "github.com/yang-zzhong/go-model"
	"time"
)

type User struct {
	Id              uint32    `db:"id bigint pk"`
	Name            string    `db:"name varchar(128) uk"`
	NickName        string    `db:"nickname varchar(128) nil"`
	BlogName        string    `db:"blog_name varchar(64)"`
	EmailAddr       string    `db:"email_addr varchar(128) nil"`
	EmailAddrAuthed bool      `db:"email_addr_authed smallint"`
	PhoneNumber     string    `db:"phone_number varchar(128) nil"`
	PhoneAuthed     bool      `db:"phone_number_authed smallint"`
	PortraitImageId string    `db:"portrait_image_id varchar(32) nil"`
	Bio             string    `db:"bio varchar(512) nil"`
	ThemeId         uint32    `db:"theme_id bigint nil"`
	Password        string    `db:"password varchar(128) protected"`
	Salt            string    `db:"salt char(8) protected"`
	CreatedAt       time.Time `db:"created_at datetime"`
	UpdatedAt       time.Time `db:"updated_at datetime"`
	*model.Base
}

func (user *User) TableName() string {
	return "users"
}

func (user *User) Encrypt(str string) string {
	md5Sumb := md5.Sum(([]byte)(str + user.Salt))
	result := hex.EncodeToString(md5Sumb[:])
	return result
}

func NewUser() *User {
	user := model.NewModel(new(User)).(*User)
	user.DeclareMany("blogs", new(Blog), map[string]string{
		"id": "user_id",
	})
	user.DeclareMany("images", new(UserImage), map[string]string{
		"id": "user_id",
	})
	user.DeclareMany("cates", new(Cate), map[string]string{
		"id": "user_id",
	})
	user.DeclareOne("current_theme", new(Theme), map[string]string{
		"theme_id": "id",
	})
	user.DeclareMany("themes", new(Theme), map[string]string{
		"id": "user_id",
	})
	user.OnUpdate(func(u interface{}) error {
		u.(*User).UpdatedAt = time.Now()
		return nil
	})
	return user
}

func (user *User) Profile() map[string]interface{} {
	result := map[string]interface{}{
		"id":                user.Id,
		"name":              user.Name,
		"portrait_image_id": user.PortraitImageId,
	}
	if m, err := user.One("current_theme"); err != nil {
		return result
	} else if m != nil {
		theme := m.(*Theme)
		result["user_info_bg_image_id"] = theme.Content["user_info_bg_image_id"]
		result["user_info_bg_color"] = theme.Content["user_info_bg_color"]
		result["user_info_fg_color"] = theme.Content["user_info_fg_color"]
	}

	return result
}

func (user *User) Instance() *User {
	user.Salt = helpers.RandString(8)
	user.Id = uuid.New().ID()
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	return user
}
